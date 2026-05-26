---
title: ansible-runner private_data_dir引发的幽灵变量问题
date: 2020-07-10
tags: [Ansible, ansible-runner, 排错]
categories: [Python]
excerpt: 使用ansible-runner执行备份时，命令行正常但Agent执行报错，最终定位到private_data_dir目录下env/extravars文件不会被更新，导致已被删除的变量像幽灵一样残留在后续执行中。
---

## 现象

使用命令行执行备份操作一切正常，但通过 Agent 执行备份时却报错——在新建 MariaDB 容器时提示无法找到镜像。

## 分析

经询问得知，版本指定为 `2.2.0`，而报错信息却说找不到 tag 为 `train` 的镜像。

查看执行的 Playbook，定位到镜像名生成代码逻辑，在 `roles/mariadb/defaults/main.yml` 中找到这段代码：

```yaml
mariabackup_image: "{{ docker_registry ~ '/' if docker_registry else '' }}{{ docker_namespace }}/{{ kolla_base_distro }}-{{ mariadb_install_type }}-mariadb"
mariabackup_tag: "{{ openstack_release }}"
mariabackup_image_full: "{{ mariabackup_image }}:{{ mariabackup_tag }}"
```

可以发现 `openstack_release` 决定了镜像的版本。继续往上追溯，在 Playbook 入口处发现其引用了同目录下 `detect-release.yml`，其中对 `openstack_release` 有相关操作。通过 debug 对比发现：**命令行执行其值为 `2.2.0`，而 Agent 执行其值为 `rocky`。**

加入 debug 语句打印变量：

```yaml
msg: "{{ mariabackup_image_full }}, {{ openstack_release }}, {{ mariadb_tag }}, {{ release }}"
```

两组结果分别为：

```
"msg": "172.17.20.5:4000/uas/centos-binary-mariadb:2.2.0, 2.2.0, 2.2.0, 2.2.0"
"msg": "172.17.20.5:4000/uas/centos-binary-mariadb:train, rocky, train, 2.2.0"
```

Agent 执行时 `openstack_release` 的值是 `rocky`，且 `rocky` 与最终结果 `train` 也不一致，说明变量链路中有更复杂的覆盖关系。

## 定位 private_data_dir

在排查测试的过程中发现，将调用 ansible 时给 `ansible_runner.run` 函数**不传递** `private_data_dir` 参数，程序就能正常运行。

进一步测试：将程序设置的 `private_data_dir` 所在目录下的 `artifacts` 和 `env` 目录删除后，程序也可以正常运行，而且后续无论是否在 `ansible_runner` 中传递 `private_data_dir` 参数，程序都能正常运行。

这强烈暗示问题出在 `private_data_dir`。查阅 ansible-runner 官方文档（https://ansible-runner.readthedocs.io/en/1.4.7/ansible_runner.html#ansible_runner.interface.run），关于该参数的描述只有一段：

> The directory containing all runner metadata needed to invoke the runner module. Output artifacts will also be stored here for later consumption.
>
> 包含调用运行器模块所需的所有运行器元数据的目录。输出工件也将存储在这里供以后使用。

因为是调用时变量不正确，所以问题很可能出在 **metadata 的获取**上。

## 复现问题

由于线上环境自从删除 `private_data_dir` 目录后异常再无出现，于是在本地创建一个小 Demo 来复现。

传入变量后，发现变量内容被写入 `env/extravars` 文件。如果修改传入的变量内容，该文件**不会更新**，但 Playbook 打印的变量值却是正确的新值。

但如果按照实际的逻辑，除了传入的数据外，还有从 YAML 文件读取的数据。在 YAML 文件中写入：

```yaml
release: "rocky"
```

清空 `private_data_dir` 目录后执行，`release` 变量被正常引用。但如果将 `release: "rocky"` 注释掉，此时 `env/extravars` 中仍然记录着 `release: "rocky"`，执行后发现程序在调用 `ansible_runner` 时并没有传入 `release`，但 Playbook 执行后仍然能打印出 `release: "rocky"`。

**也就是说，因为 `env/extravars` 不会被更新，导致在取消一个变量后，该变量在执行 Playbook 时仍然存在——就像幽灵一样。**

## 源码印证

阅读 `ansible_runner` 的源码，在 `ansible_runner.utils.py` 下的 `dump_artifacts` 函数中找到了问题所在：

`dump_artifact` 函数负责将数据写入磁盘，可以看到**只有在目标文件不存在时才会写入**，不具备更新功能。已存在的 `env/extravars` 文件永远不会被覆盖。

## 还原故障全貌

结合线上环境的情况，有理由相信事情是这样发生的：

最开始 `global.yaml` 中配置了：

```yaml
openstack_release: "rocky"
openstack_aggressive_release: "train"
```

执行了一遍程序后，`private_data_dir` 目录下的 `env/extravars` 文件中记录了：

```yaml
openstack_release: "rocky"
openstack_aggressive_release: "train"
```

后续修改了 `global.yaml` 为：

```yaml
release: "2.2.0"
# openstack_aggressive_release: "2.2.0"
# ryze_tag: "train"
```

再执行备份时，由于 `env/extravars` 不会被更新，以下幽灵变量依然存在：

```yaml
openstack_release: "rocky"
openstack_aggressive_release: "train"
```

又被加载进了执行上下文。而 `/etc/maine-swallow/group_vars/maine_all.yml` 下有这样的逻辑：

```yaml
mariadb_tag: "{{ openstack_aggressive_release | default('train') }}"
mariabackup_tag: "{{ mariadb_tag }}"
```

`/usr/share/kolla-ansible/ansible/roles/mariadb/defaults/main.yml` 中则有：

```yaml
####################
# Backups
####################
mariabackup_image: "{{ docker_registry ~ '/' if docker_registry else '' }}{{ docker_namespace }}/{{ kolla_base_distro }}-{{ mariadb_install_type }}-mariadb"
mariabackup_tag: "{{ openstack_release }}"
mariabackup_image_full: "{{ mariabackup_image }}:{{ mariabackup_tag }}"

mariadb_backup_host: "{{ groups['mariadb'][0] }}"
mariadb_backup_database_schema: "PERCONA_SCHEMA"
mariadb_backup_database_user: "backup"
mariadb_backup_database_address: "{{ database_address }}"
mariadb_backup_type: "full"
```

残留的 `openstack_release: "rocky"` 和 `openstack_aggressive_release: "train"` 变量被渲染后，最终导致 `mariabackup_tag` 变成了 `train`，镜像拉取失败。

## 解决方案

虽然删除错误的 `env/extravars` 文件后可以临时解决问题，但 `private_data_dir` 的使用并不是很必要，而且会导致这类不易排查的 Bug。建议**取消该参数的使用**，避免变量残留问题再次发生。
