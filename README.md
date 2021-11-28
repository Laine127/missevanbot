# MissEvan-FM

猫耳 FM 直播间机器人 Go 语言实现，仅作为娱乐用途。未来看心情更新~

## 功能

- 新观众欢迎
- 新关注感谢
- 礼物感谢
- 直播间在线人数查看（可能对手机端比较有用？）
- 签到和签到排行榜
- 定时发送彩虹屁（
- 监听直播间状态（如开播/下播）
- 消息推送
- 汉语注音功能（提供给中文学习者）
- 判断不同用户（机器人管理员、主播、房管、普通用户）
- 输出日志

## 如何使用？

编译本项目，在可执行文件同目录下创建 _config.yaml_ 文件，填入配置信息，执行可执行文件即可。

```yaml
name: "芝士Bot" # Bot 昵称，必须与帐号昵称完全相同
cookie: ".cookie" # 存储 Cookie 的文件路径 
level: "debug" # 日志等级
redis: # Redis 相关配置
  host: ""
  passwd:
  db: 1
push: # 各类推送服务密钥
  bark: "" # Bark App 推送通知
rooms: # 需要启用的直播间
  - id: 111111111
    name: "主播一号" # 主播昵称，可以随意自定义，暂时没有用处
    pinyin: false # 是否开启用户名注音功能
    rainbow_max_interval: 10 # 彩虹屁发送的最大时间间隔，单位：分钟
  - id: 222222222
    name: "主播二号"
```

```shell
# Windows
go build
.\missevan-fm.exe

# Linux
go build
./missevan-fm
```

## 目录结构

- bot
    - config.go：配置文件读取模块
    - redis.go：Redis 连接相关模块
- handler：处理房间各类消息的模块
    - chat.go：处理聊天信息
    - command.go：命令消息处理
    - const.go：直播间消息 JSON 结构体定义
    - handler.go：消息处理入口
- module：各独立模块
    - message.go：处理消息发送
    - module.go：初始化
    - praise.go：彩虹屁模块
    - push.go：消息推送模块
    - room.go：直播间相关模块
    - sign.go：签到模块
- connect.go：处理 Websocket 连接
- main.go：程序入口
