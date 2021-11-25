# MissEvan-FM

猫耳 FM 直播间机器人 Go 语言实现，仅作为娱乐用途。未来看心情更新~

## 功能

- 新观众欢迎
- 新关注感谢
- 礼物感谢
- 直播间在线人数查看（可能对手机端比较有用？）
- 签到和签到排行榜
- 定时发送彩虹屁
- 监听直播间状态（如开播/下播）
- 消息推送

## 如何使用？

编译本项目，在可执行文件同目录下创建 _config.yaml_ 文件，填入配置信息，执行可执行文件即可。

```yaml
cookie: ".cookie" # 存储 Cookie 的文件路径 
redis: # Redis 相关配置
  host: ""
  passwd:
  db: 1
push: # 各类推送服务密钥
  bark: "" # Bark App 推送通知
```

```shell
# .\missevan-fm.exe [直播间 ID]
.\missevan-fm.exe 6666666
```

## 机制

- 基于 Websocket 连接监听直播间消息
- 每隔 30 秒发送心跳
- 使用 Cookie 设置登录态，在直播间发送消息

## 目录结构

- cache：Redis 连接相关模块
- config：配置文件读取模块
- module：各独立模块
    - push：消息推送模块
    - room：直播间相关模块
    - send：发送消息模块
    - sign：签到模块
- connect.go：处理 Websocket 连接
- main.go：程序入口
- message.go：直播间消息 JSON 结构体
- message_handler.go：处理房间各类消息