# gocron - 定时任务管理系统

> 本项目基于 [ouqiang/gocron](https://github.com/ouqiang/gocron) 修改开发  
> 原作者：[@ouqiang](https://github.com/ouqiang)  
> License：MIT License  
> Copyright (c) 2017 qiang.ou  
> Modified by Blazing (2026)

---

[![Downloads](https://img.shields.io/github/downloads/HyhBlazing/gocron/total.svg)](https://github.com/HyhBlazing/gocron/releases)
[![License](https://img.shields.io/github/license/HyhBlazing/gocron.svg)](https://github.com/HyhBlazing/gocron/blob/master/LICENSE)
[![Release](https://img.shields.io/github/release/HyhBlazing/gocron.svg?label=Release)](https://github.com/HyhBlazing/gocron/releases)

---

# 项目简介

使用 Go 语言开发的轻量级定时任务集中调度与管理系统，用于替代 Linux Crontab。

原项目中的延迟任务已拆分为独立项目：[delay-queue](https://github.com/ouqiang/delay-queue)


---

## 功能特性

### 核心能力

- Web 界面集中管理定时任务
- 支持 Crontab 表达式（精确到秒）
- 任务失败自动重试
- 任务超时强制终止
- 任务依赖配置（A 任务完成后执行 B 任务）
- 账户权限控制
- 任务执行日志查询
- 执行结果通知（支持邮件、Webhook）

---

### 任务类型

#### Shell 任务

> 在任务节点上执行 Shell 命令  
> 支持多个节点同时执行

#### HTTP 任务

> 由调度器直接发起 HTTP 请求  
> 不依赖任务节点

---

### 截图

### 调度流程

![流程图](https://raw.githubusercontent.com/ouqiang/gocron/master/assets/screenshot/scheduler.png)

---

### 所有定时列表
<img width="2561" height="1919" alt="image" src="https://github.com/user-attachments/assets/a369bdce-dab7-4bd1-b81e-15cacbfb3b32" />

### 定时详情

![定时详情](https://github.com/user-attachments/assets/904bf962-8bd8-42a1-bd7c-7b961bcad242)

定时列表日志
<img width="1534" height="971" alt="image" src="https://github.com/user-attachments/assets/91d83cc2-7bd4-4031-a175-8118c429f214" />

定时编辑
<img width="2561" height="1919" alt="image" src="https://github.com/user-attachments/assets/f6ca99ab-19bc-4d24-a50b-cbad2d041617" />

### 定时日志列表

![日志列表](https://github.com/user-attachments/assets/4059b097-56ff-4ecb-b3f9-04f6c102d501)

### 定时日志详情

<img width="1205" height="526" alt="image" src="https://github.com/user-attachments/assets/e3933132-5ec1-42b9-a1c4-b904a92851cd" />

### 定时日志执行结果
<img width="1222" height="540" alt="image" src="https://github.com/user-attachments/assets/d8e0ba88-cc92-4f21-9fb9-b527a2475be2" />

### 统计模块

![统计](https://github.com/user-attachments/assets/fdb31d76-affd-48f6-a85c-585aa2d6e6f8)

---

## 支持平台

- Linux

---

## 环境要求

- MySQL

---

## 下载

[Releases](https://github.com/HyhBlazing/gocron/releases)

---

## 安装说明

### 二进制安装

1. 解压压缩包  
2. 进入解压目录：

```bash
cd gocron
```
3. 启动服务
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
