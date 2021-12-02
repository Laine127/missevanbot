package models

// Handler related
const (
	TplBotStart    = "芝士机器人在线了，可以在直播间输入 “帮助” 或者 @我 来获取支持哦～" // 机器人启动
	TplWelcome     = "欢迎 @%s 进入直播间~"                        // 欢迎用户
	TplWelcomeAnon = "欢迎新同学进入直播间~"                          // 欢迎匿名用户
	TplThankFollow = "感谢 @%s 的关注~"                          // 感谢关注
	TplThankGift   = "感谢 @%s 赠送的%d个%s~"                     // 感谢礼物
)

// Command related
const (
	TplRoomInfo  = "当前直播间信息：\n- 房间名：%s\n- 公告：%s\n- 主播：%s\n- 直播平台：%s\n- 当前在线：%d\n- 累计人数：%d\n- 管理员：\n" // 房间信息模板
	TplRankEmpty = "今天的榜单好像空空的~"                                                                     // 榜单为空
	TplBaitStop  = "我突然有点困了"                                                                         // 关闭演员模式
	TplMusicAdd  = "点歌 %s 成功啦~"
	TplMusicNone = "当前还没有人点歌哦~"
	TplMusicDone = "完成了一首歌曲~"
)

// Modules related
const (
	TplSignDuplicate = "已经签到过啦\n\n您已经连续签到%s天\n\n%s"
	TplSignSuccess   = "签到成功啦，已经连续签到%d天~\n\n%s\n\n%s"
)

const TplDefaultPoem = "孜孜不倦，不易乎世。"

var LuckPool = [...]string{
	"大凶", "大吉", "不详",
	"吉兆", "无事", "隆运",
	"破财", "稳步", "跌宕",
	"未可知",
}

var ChatPool = [...]string{
	"你好呀～",
	"QAQ",
	"有什么事情嘛～",
	"不好意思呀我暂时没办法理你～",
	"抱歉现在有点忙～",
}
