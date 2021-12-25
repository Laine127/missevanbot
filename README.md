# missevanbot

猫耳 FM 直播间机器人（MissEvan Bot）Go 语言实现，仅作为娱乐用途。未来看心情更新~

## 功能

- 观众欢迎
- 用户名汉语注音（提供给中文学习者）
- 关注感谢
- 礼物感谢
- 直播间数据查看
- 直播间签到及签到排行
- 星座运势（全是好运）
- 指定城市天气查询
- 用户点歌记录
- 随机彩虹屁（
- Pia 戏戏文显示（支持防屏蔽）
- 小游戏
    - 数字炸弹
    - 击鼓传花
    - 你说我猜
- 小游戏战绩榜单
- 用户权限划分（机器人管理员、主播、房管、普通用户）
- 直播间状态监控（如开播/下播）
- 消息推送
- 运行日志

## 如何使用？

编译本项目，在可执行文件同目录下创建 _config.yaml_ 文件，填入配置信息，执行可执行文件即可。

```yaml
nickname: "知世" # 机器人昵称
log: "data/missevan.log" # 日志文件路径
level: "info" # 日志输出等级
redis: # Redis 相关配置
  host: ""
  passwd: ""
  db: 0
push: # 各类推送服务密钥
  bark: "" # Bark App 推送通知
admin: 11111 # 管理员 ID
rooms: # 需要启用的直播间
  - id: 111111111
    name: "主播一号" # 主播昵称，可以随意自定义，暂时没有用处
    enable: true # 是否为当前直播间启用机器人
    watch: true # 是否监控开播/下播
  - id: 222222222
    name: "主播二号"
    watch: false
```

```shell
# Windows
go build
.\missevanbot.exe

# Linux
go build
./missevanbot
```

### Docker

1. 创建 *~/missevan* 目录，放入 *config.yaml*。
2. 将 *config.yaml* 中的 `redis.host` 改为 `"redis:6379"`
3. 构建镜像：
    ```shell
    cd missevanbot/
    docker build -t missevan-bot .
    ```
4. 启动容器：
    ```shell
    export REDIS_PASSWORD=mypassword # 设置 Redis 访问密码
    docker-compose up
    ```

## 目录结构

- cmd：主函数入口
- config：需要初始化的模块
    - config.go：配置文件
    - redis.go：Redis 客户端
- core
    - connect.go：Websocket 连接处理，获取消息
    - cron.go：定时任务
    - match.go：处理消息
    - send.go：发送消息
- handlers：处理房间各类消息的模块
    - game：游戏相关模块
    - chat.go：处理聊天信息
    - command.go：命令消息处理
    - message.go：消息处理入口
    - keyword.go：关键词消息处理
- models：结构体模型
    - command.go：命令相关
    - game.go：游戏相关结构
    - message.go：直播间消息相关
    - room.go：直播间实例
    - template.go：消息模板
- modules：各独立模块
    - thirdparty：第三方组件
    - checkin.go：签到模块
    - fm.go：猫耳 FM 相关模块
    - http.go：HTTP 请求模块
    - game.go：游戏相关模块
    - mode.go：模式相关模块
    - push.go：消息推送模块
    - tasks.go：定时任务模块
    - template.go：消息模板模块
    - word.go：文本模板模块
- utils：辅助工具
    - logger：日志组件