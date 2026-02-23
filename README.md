本项目基于 [gocron](https://github.com/ouqiang/gocron) 项目修改

原作者：[ouqiang](https://github.com/ouqiang)

采用 MIT License

Copyright (c) 2017 qiang.ou

Modified by Blazing 2026

# gocron - 定时任务管理系统

[![Downloads](https://img.shields.io/github/downloads/ouqiang/gocron/total.svg)](https://github.com/HyhBlazing/gocron/releases)
[![license](https://img.shields.io/github/license/mashape/apistatus.svg?maxAge=2592000)](https://github.com/HyhBlazing/gocron/blob/master/LICENSE)
[![Release](https://img.shields.io/github/release/ouqiang/gocron.svg?label=Release)](https://github.com/HyhBlazing/gocron/releases)

# 项目简介

使用Go语言开发的轻量级定时任务集中调度和管理系统, 用于替代Linux-crontab

原有的延时任务拆分为独立项目[延迟队列](https://github.com/ouqiang/delay-queue)

## 功能特性

- Web界面管理定时任务
- crontab时间表达式, 精确到秒
- 任务执行失败可重试
- 任务执行超时, 强制结束
- 任务依赖配置, A任务完成后再执行B任务
- 账户权限控制
- 任务类型
  - shell任务
    > 在任务节点上执行shell命令, 支持任务同时在多个节点上运行
  - HTTP任务
    > 访问指定的URL地址, 由调度器直接执行, 不依赖任务节点
- 查看任务执行结果日志
- 任务执行结果通知, 支持邮件、Webhook

### 截图

![流程图](https://raw.githubusercontent.com/ouqiang/gocron/master/assets/screenshot/scheduler.png)
所有定时列表
<img width="2561" height="1919" alt="image" src="https://github.com/user-attachments/assets/a369bdce-dab7-4bd1-b81e-15cacbfb3b32" />
定时详情
<img width="1332" height="568" alt="image" src="https://github.com/user-attachments/assets/904bf962-8bd8-42a1-bd7c-7b961bcad242" />
定时列表日志
<img width="1534" height="971" alt="image" src="https://github.com/user-attachments/assets/91d83cc2-7bd4-4031-a175-8118c429f214" />
定时编辑
<img width="2561" height="1919" alt="image" src="https://github.com/user-attachments/assets/f6ca99ab-19bc-4d24-a50b-cbad2d041617" />
定时日志列表
<img width="2561" height="1919" alt="image" src="https://github.com/user-attachments/assets/4059b097-56ff-4ecb-b3f9-04f6c102d501" />
定时日志详情
<img width="1205" height="526" alt="image" src="https://github.com/user-attachments/assets/e3933132-5ec1-42b9-a1c4-b904a92851cd" />
定时日志执行结果
<img width="1222" height="540" alt="image" src="https://github.com/user-attachments/assets/d8e0ba88-cc92-4f21-9fb9-b527a2475be2" />
统计部分
<img width="2561" height="1919" alt="image" src="https://github.com/user-attachments/assets/fdb31d76-affd-48f6-a85c-585aa2d6e6f8" />

### 支持平台

> Linux

### 环境要求

> MySQL

## 下载

[releases](https://github.com/HyhBlazing/gocron/releases)

## 安装

### 二进制安装

1. 解压压缩包  
2. `cd 解压目录`
3. 启动

- 调度器启动
  - Linux、Mac OS: `./gocron web`
- 任务节点启动, 默认监听0.0.0.0:5921
  - Linux、Mac OS: `./gocron-node`

4. 浏览器访问 http://localhost:5920

## 反馈

提交[issue](https://github.com/HyhBlazing/gocron/issues/new)

## ChangeLog

## v1.5.6
- 所有定时
  - 定时名称：定时名称、定时ID、标签、创建时间，放在了一列。
  - cron表达式：支持了查看未来10次执行时间，鼠标悬浮可以复制
  - 筛选项
    - 名称ID筛选：支持了定时ID和定时名称的筛选
    - 支持标签列表下拉，支持搜索
  - 详情部分修改为按钮模态窗打开
  - 增加日志当前页面模态窗打开
- 定时日志
  - 筛选项
    - 支持定时ID、定时名称筛选
    - 支持标签列表下拉，支持搜索
  - 在列表增加了任务对应标签的查看
  - 日志详情通过模态窗打开
  - 日志结果通过模态窗打开
- 通知配置
  - 去掉了Slack的配置，我没用过，所以隐藏掉了
- 节点管理
  - 记忆中好像没做更改
- 用户管理
  - 增加了角色筛选
  - 增加邮箱搜索
- 登录日志
  - 增加了用户名、登录IP的搜索
- 统计
  - 查询所有定时总数、日志总数、用户总数、登录日志总数
  - 也有部分分析的，可以看图。

## v1.5

- 前端使用Vue+ElementUI重构
- 任务通知
  - 新增WebHook通知
  - 自定义通知模板
  - 匹配任务执行结果关键字发送通知
- 任务列表页显示任务下次执行时间

## v1.4

- HTTP任务支持POST请求
- 后台手动停止运行中的shell任务
- 任务执行失败重试间隔时间支持用户自定义
- 修复API接口调用报403错误

## v1.3

- 支持多用户登录
- 增加用户权限控制

## v1.2.2

- 用户登录页增加图形验证码
- 支持从旧版本升级
- 任务批量开启、关闭、删除
- 调度器与任务节点支持HTTPS双向认证
- 修复任务列表页总记录数显示错误

## v1.1

- 任务可同时在多个节点上运行
- \*nix平台默认禁止以root用户运行任务节点
- 子任务命令中增加预定义占位符, 子任务可根据主任务运行结果执行相应操作
- 删除守护进程模块
- Web访问日志输出到终端
