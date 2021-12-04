package models

import (
	"math/rand"
	"time"
)

// Handler related
const (
	TplBotStart    = "芝士机器人在线啦，可以在直播间输入 “帮助” 或者 @我 来获取支持哦～" // 机器人启动
	TplWelcome     = "欢迎 @%s 进入直播间~"                        // 欢迎用户
	TplWelcomeAnon = "欢迎新同学进入直播间~"                          // 欢迎匿名用户
	TplThankFollow = "感谢 @%s 的关注~"                          // 感谢关注
	TplThankGift   = "感谢 @%s 赠送的%d个%s~"                     // 感谢礼物
)

// Command related
const (
	TplRoomInfo = `当前直播间信息：
- 房间名：%s
- 主播：%s
- 直播平台：%s
- 当前在线：%d
- 累计人数：%d
- 管理员：`  // 直播间信息模板
	TplRankEmpty = "今天的榜单好像空空的~" // 榜单为空
	TplBaitStop  = "我突然有点困了"     // 关闭演员模式
	TplMusicAdd  = "点歌 %s 成功啦~"
	TplMusicNone = "当前还没有人点歌哦~"
	TplMusicDone = "完成了一首歌曲~"
)

// Modules related
const (
	TplSignDuplicate = "已经签到过啦\n\n您已经连续签到%s天\n\n%s"
	TplSignSuccess   = "签到成功啦，已经连续签到%d天~\n\n%s\n\n%s"
	TplStarFortune   = "%s今日幸运值：%s\n\n%s"
	TplPiaStart      = "pia戏模式启动成功啦～"
	TplPiaStop       = "pia戏模式结束成功啦～"
)

const TplDefaultPoem = "孜孜不倦，不易乎世。"

var _lucks = [...]string{
	"连理之木", "景星庆云", "有凤来仪",
	"避凶趋吉", "逢凶化吉", "时来运转",
	"大凶", "不祥", "趋凶", "平平",
	"大吉", "兆祥", "趋吉", "平平",
	"稳步", "跌宕", "生财", "破财",
	"未可知", "不可说", "玄之又玄",
	"怎么说呢还不错", "怎么说呢就很好", "怎么说呢一般般",
	"我的天呐运气爆棚", "我的天呐不是很好", "我的天呐蛮不错的",
	"或许有财运", "或许有桃花运", "或许有狗屎运",
	"不告诉你", "就不告诉你", "不想告诉你",
	"还是不要知道为好", "还是少打听比较好", "还是不让你知道比较好",
}

var _chats = [...]string{
	"你好呀～",
	"QAQ",
	"有什么事情嘛～",
	"不好意思呀我暂时没办法理你～",
	"抱歉现在有点忙～",
}

// LuckString 返回运势字符串
func LuckString() string {
	result := "今日运势："
	rand.Seed(time.Now().UnixNano())
	return result + _lucks[rand.Intn(len(_lucks))]
}

// ChatString 返回聊天字符串
func ChatString() string {
	rand.Seed(time.Now().UnixNano())
	return _chats[rand.Intn(len(_chats))]
}
