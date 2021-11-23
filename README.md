# MissEvan-FM

猫耳 FM 直播间机器人 Go 语言实现，模块化封装扩展功能。未来看心情更新~。

## 功能

- 新观众欢迎
- 新关注感谢
- 礼物感谢
- 直播间人数查看（可能对手机端比较有用？）

## 如何使用？

编译本项目，在可执行文件同目录下创建 _.cookie_ 文件，填入当前账号的 Cookie，执行可执行文件即可。

```shell
# .\missevan-fm.exe [直播间 ID]
.\missevan-fm.exe 6666666
```

## 机制

- 基于 Websocket 连接监听直播间消息
- 每隔 30 秒发送心跳
- 使用 Cookie 设置登录态，在直播间发送消息

## 目录结构

- module：各独立模块
    - send：发送消息模块
- connect.go：处理 Websocket 连接
- handler.go：处理房间各类消息
- main.go：程序入口
- message.go：直播间消息 JSON 结构体