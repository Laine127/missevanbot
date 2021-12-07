package models

import (
	"math/rand"
	"time"
)

// handlers related
const (
	TplBotStart    = "芝士机器人在线啦，可以在直播间输入 “帮助” 或者 @我 来获取支持哦～" // 机器人启动
	TplWelcome     = "欢迎 @%s 同学进入直播间~"                      // 欢迎用户
	TplWelcomeAnon = "欢迎新来的小可爱们进入直播间呀~"                     // 欢迎匿名用户
	TplThankFollow = "谢谢 @%s 的关注呀~"                         // 感谢关注
	TplThankGift   = "感谢 @%s 赠送的%d个%s~"                     // 感谢礼物
)

// commands related
const (
	TplRoomInfo = `当前直播间信息：
- 房间名：%s
- 主播：%s
- 粉丝数：%d
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

// modules related
const (
	TplSignDuplicate = "已经签到过啦\n\n您已经连续签到%s天\n\n%s"
	TplSignSuccess   = "签到成功啦，已经连续签到%d天~\n\n%s\n\n%s"
	TplStarFortune   = "%s今日幸运值：%s\n\n%s"
	TplPiaStart      = "pia戏模式启动成功啦～\n本文角色：%s\n请注意，如果输出文本时没有出现结果，很可能是存在违禁词，请使用重定位命令 r 定位到需要的位置，再通过 s 指令进行单条输出~"
	TplPiaEmpty      = "pia戏模式还没有启动哦，输入“帮助”来获取帮助吧~"
	TplPiaDone       = "本篇文章已经结束啦～"
	TplPiaOutBounds  = "选择的位置超过范围啦"
	TplPiaRelocate   = "成功重定位啦~"
	TplPiaStop       = "pia戏模式结束成功啦～"
)

const (
	TplGameExists       = "当前游戏未结束，请先结束"
	TplGameCreate       = "游戏创建成功啦~"
	TplGameJoin         = "玩家 @%s 加入成功~"
	TplGameJoinDup      = "加入过的玩家不要再加入啦"
	TplGameStart        = "游戏开始啦~"
	TplGameNext         = "请玩家 @%s 进行操作"
	TplGameStop         = "终止游戏成功啦~"
	TplGameNotEnough    = "人数不足，无法开始游戏~"
	TplGameInputIllegal = "拜托，输入正确的内容好嘛"
	TplGameBomb         = "@%s 踩到了炸弹~"
	TplGameBombRange    = "炸弹在 [%d, %d] 区间内（包含边界）哦~"
)

const TplDefaultPoem = "孜孜不倦，不易乎世。"

// _lucks 运势词汇池
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

// _chats 对话语句池
var _chats = [...]string{
	"你好呀～",
	"QAQ",
	"有什么事情嘛～",
	"不好意思呀我暂时没办法理你～",
	"抱歉现在有点忙～",
}

// _comforts 宽慰语句池
var _comforts = [...]string{
	"不可以emo哦~",
	"怎么啦，为什么要说emo呢",
	"可以为了我不要伤心嘛",
	"有什么难过的事情可以告诉主播呀",
	"不开心的事就忘光光吧",
	"我不知道怎么安慰你才好，但是生活总得继续",
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

// ComfortString 返回宽慰语句字符串
func ComfortString() string {
	rand.Seed(time.Now().UnixNano())
	return _comforts[rand.Intn(len(_comforts))]
}
