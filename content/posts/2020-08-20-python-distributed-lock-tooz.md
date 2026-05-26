---
title: 用tooz实现Python分布式锁：从选型到实现
date: 2020-08-20
tags: [分布式锁, tooz, Redis]
categories: [Python]
excerpt: 在Slark项目中，多个Agent节点需要串行执行Ansible命令，借助OpenStack的分布式协调库tooz实现基于Redis的分布式锁，通过装饰器优雅地控制并发。
---

## 现象

在 Slark API 向 Worker 发送执行 Ansible 命令的请求时，希望这些任务不要并行执行。最终执行 Ansible 的节点为 Agent 节点，而 Agent 可能被部署多个实例，因此需要使用**分布式锁**来控制执行规则。

## 技术选型

分布式锁的实现选用 **tooz**，它是 OpenStack 的分布式管理工具库。

### tooz 组件功能

**Tooz 项目旨在通过提供协调 API 帮助开发人员构建分布式应用程序**，集中最常见的分布式原语，如组成员协议、锁定服务和领导者选举。

Tooz 是一个 Python 库，提供了标准的 coordination API。最初由 eNovance 的几位工程师编写，其主要目标是解决分布式系统的通用问题，比如节点管理、主节点选举以及分布式锁等。Tooz 抽象了高级接口，支持对接十多种 DLM 驱动，比如 ZooKeeper、Redis、MySQL、Etcd、Consul 等。

在 Slark 项目中主要使用 tooz 提供的分布式锁，从而允许分布式节点获取和释放锁来实现同步，以此解决多个 Agent 进程的并发问题。

## 锁的粒度

在 Slark 项目中，锁的最大粒度应该是**集群**。不同集群之间不会相互干扰。

而对于 Agent 进程，其天然最大粒度就是集群，因此只需要控制好其内部函数锁的粒度即可。

在需要控制的函数中，其操作的资源主要有：

- **Slurm 配置文件**
- **Linux 用户**

这两个资源之间没有很强的关联，可以并行执行。

需要注意的是，Slurm 的配置文件需要被控制成串行——同一时刻只能有一个刷配置的任务在执行。此外，用户的创建、删除和修改密码也需要加锁串行执行，以防止操作顺序被打乱，保证执行顺序的绝对正确。

## 实现

当前 Agent 中 Ansible 部分的代码结构如下：

```python
class AnsibleProvision(base.BaseProvision):
    def __init__(self, *args, **kwargs):
        self.extraVars = {}
        self.playbooks_dir = os.path.join(
            os.path.abspath(__file__).rsplit(os.sep, 1)[0],
            "playbooks"
        )
        if not os.path.exists(self.playbooks_dir):
            os.makedirs(self.playbooks_dir)
        self.inventory = os.path.join(self.playbooks_dir, "inventory")
        super(AnsibleProvision, self).__init__(*args, **kwargs)

    def _run(self, extra_vars=None, **kwargs):
        if "inventory" not in kwargs:
            kwargs["inventory"] = self.inventory
        _extra_vars = self.extraVars.copy()
        if extra_vars:
            _extra_vars.update(extra_vars)
        response = ansible_runner.run(
            extravars=_extra_vars,
            **kwargs
        )

        response.failed_reasons = collect_response_failed_list(response)
        return response

    def refresh_inventory(self, inventory):
        provision_utils.refresh_inventory(self.inventory, inventory)

    def run_playbook(self, playbook, extra_vars, **kwargs):
        LOG.debug("run playbook: {}, extra_vars: {}".format(
            playbook,
            extra_vars
        ))
        response = self._run(
            extra_vars,
            playbook=playbook,
            **kwargs
        )

        if response.status == "failed":
            LOG.error(collect_response_failed_info(response))
        return response

    def check(self):
        LOG.debug("[ansible] check node status")
        response = self._run(module="ping", host_pattern="all")
        return response

    def load_config(self, cluster_name, nodes, partitions):
        LOG.debug("ansible driver load_config :\n"
                  "cluster_name: {}\n"
                  "nodes: {}\n"
                  "partitions: {}".format(cluster_name, nodes, partitions))

        extra_vars = {
            "cluster_name": cluster_name,
            "nodes": nodes,
            "partitions": partitions
        }

        _playbook = os.path.join(self.playbooks_dir,
                                 "load_slurm_config.yaml")

        response = self.run_playbook(_playbook, extra_vars)
        return response

    def create_users(self, users):
        LOG.debug("[ansible] create user")
        response = self.run_playbook(
            playbook=os.path.join(self.playbooks_dir, "create_user.yaml"),
            extra_vars={
                "users": users
            }
        )
        return response

    def delete_users(self, users):
        LOG.debug("[ansible] delete user")
        response = self.run_playbook(
            playbook=os.path.join(self.playbooks_dir, "delete_user.yaml"),
            extra_vars={
                "users": users
            }
        )
        return response

    def set_user_password(self, username, password):
        LOG.debug("[ansible] set user password")
        response = self.run_playbook(
            playbook=os.path.join(self.playbooks_dir, "set_user_password.yaml"),
            extra_vars={
                "username": username,
                "password": password
            }
        )
        return response
```

Agent 本身作用的范围只是自己所管理的单个集群。基于此，可以这样编写一个 lock 装饰器，将需要互斥的函数变为串行执行：

```python
from oslo_utils import reflection
from tooz import coordination

coordinator = coordination.get_coordinator(
    CONF.coordination.backend_url,
    uuidutils.generate_uuid().encode('ascii'))
coordinator.start(start_heart=True)

def lock(name):

    def wrap(f):

        @functools.wraps(f)
        def inner(*args, **kwargs):
            _lock = coordinator.get_lock(name)
            t1 = timeutils.now()
            t2 = None
            try:
                with _lock.acquire(blocking=False):
                    t2 = timeutils.now()
                    LOG.debug('Lock "%(name)s" acquired by "%(function)s" :: '
                              'waited %(wait_secs)0.3fs',
                              {'name': name,
                               'function': reflection.get_callable_name(f),
                               'wait_secs': (t2 - t1)})
                    return f(*args, **kwargs)
            finally:
                t3 = timeutils.now()
                if t2 is None:
                    held_secs = "N/A"
                else:
                    held_secs = "%0.3fs" % (t3 - t2)
                LOG.debug('Lock "%(name)s" released by "%(function)s" :: held '
                          '%(held_secs)s',
                          {'name': name,
                           'function': reflection.get_callable_name(f),
                           'held_secs': held_secs})
        return inner

    return wrap
```

## 使用

在需要加锁的方法上添加装饰器即可，锁的名称区分不同资源：

```python
class AnsibleProvision(base.BaseProvision):
    ......

    @lock("config")
    def load_config(self, cluster_name, nodes, partitions):
        LOG.debug("ansible driver load_config :\n"
                  "cluster_name: {}\n"
                  "nodes: {}\n"
                  "partitions: {}".format(cluster_name, nodes, partitions))

        extra_vars = {
            "cluster_name": cluster_name,
            "nodes": nodes,
            "partitions": partitions
        }

        _playbook = os.path.join(self.playbooks_dir,
                                 "load_slurm_config.yaml")

        response = self.run_playbook(_playbook, extra_vars)
        return response

    @lock("user")
    def create_users(self, users):
        LOG.debug("[ansible] create user")
        response = self.run_playbook(
            playbook=os.path.join(self.playbooks_dir, "create_user.yaml"),
            extra_vars={
                "users": users
            }
        )
        return response

    @lock("user")
    def delete_users(self, users):
        LOG.debug("[ansible] delete user")
        response = self.run_playbook(
            playbook=os.path.join(self.playbooks_dir, "delete_user.yaml"),
            extra_vars={
                "users": users
            }
        )
        return response

    @lock("user")
    def set_user_password(self, username, password):
        LOG.debug("[ansible] set user password")
        response = self.run_playbook(
            playbook=os.path.join(self.playbooks_dir, "set_user_password.yaml"),
            extra_vars={
                "username": username,
                "password": password
            }
        )
        return response
```

通过 `config` 和 `user` 两个不同的锁名，Slurm 配置刷新和用户管理各自串行，同时两者之间互不干扰，可以并行执行。
