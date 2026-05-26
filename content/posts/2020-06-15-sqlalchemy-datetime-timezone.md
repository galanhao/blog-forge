---
title: 当SQLAlchemy遇上时区：datetime naive vs aware踩坑记
date: 2020-06-15
tags: [SQLAlchemy, MySQL, 时区]
categories: [Python]
excerpt: 使用SQLAlchemy更新记录时遇到"can't compare offset-naive and offset-aware datetimes"报错，深入排查后发现是数据库datetime类型不支持时区信息，导致读出的时间从aware变成了naive。
---

## 现象

在归档系统中，通过运营平台对"系统备份"下的"备份位置"进行任意修改操作时，API 会抛出 500 错误，而且信息根本无法被修改。

## 问题分析

查看日志可以发现，报错信息为：

> *can't compare offset-naive and offset-aware datetimes*

查阅 Python 官方文档（https://docs.python.org/3/library/datetime.html）可知，Date 和 Time 根据其是否携带时区信息（timezone）可以分为 **Aware** 和 **Naive** 两种：

- **Aware**：携带时区信息
- **Naive**：不携带时区信息，其使用的具体是 UTC、本地时区还是其他时区，由调用程序自行约定

根据报错信息可以确定，问题的根源在于 `updated_at` 字段赋值时，原先的值和后续赋予的值一个是 Naive（`datetime.datetime(2022, 5, 18, 7, 39, 15)`），另一个是 Aware（`datetime.datetime(2022, 5, 18, 7, 48, 11, tzinfo=<UTC>)`）。

通过打点调试发现，`updated_at` 原先的值为 Naive，而后续赋予的值为 Aware。SQLAlchemy 在执行 Update 操作时，会比较原先的值和新赋予的值，如果不一致则更新（通过阅读报错对应的源码可以印证这一逻辑）。

奇怪的是，阅读代码逻辑后发现，`updated_at` 从创建到更新使用的都是同一个方法生成的 **Aware** 时间，但再次用 SQLAlchemy 读出时，却变成了 Naive 时间。

## 追根溯源

在 SQLAlchemy 的官方 Issue（https://github.com/sqlalchemy/sqlalchemy/issues/1462）中找到了类似的问题，作者也提到并不是所有数据库都能支持 Aware 时间的存储。

进一步查阅 MariaDB 官方文档（https://mariadb.com/kb/en/datetime/#time-zones），其中有明确说明：

> If a column uses the `DATETIME` data type, then any inserted values are stored as-is, so no automatic time zone conversions are performed.

也就是说，MariaDB 的 `DATETIME` 类型是**不支持** Aware 时间格式的——时区信息在写入时会被丢弃。

## 解决方案

既然数据库层面的 `DATETIME` 类型不支持时区，那么时区问题就需要在代码层面解决：

1. **统一返回 Naive 时间**：因为数据库获取当前时间统一使用 `lich.utils.tz.py` 下的 `utc_now` 函数，直接让该函数返回 `datetime.datetime.utcnow()` 即可，这样返回的是不带时区信息的 Datetime 对象，数据库统一存储 UTC 时间，不会再出现 naive 与 aware 混用的问题。

2. **使用更标准的方法**：基于上面的思路，也可以选用更具规范性的 `oslo_utils.timeutils.utcnow()`，同样能达到效果。

## 扩展

后续在定义数据库 Model 时，不必在每个 Model 中手动添加 `updated_at` 和 `created_at`，可以直接继承 `oslo_db.sqlalchemy.models.TimestampMixin`，省去重复劳动。
