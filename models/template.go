package models

const TplDefaultPoem = "孜孜不倦，不易乎世。"

// Handlers related templates.
const (
	TplWelcomeAnon = "欢迎新来的小可爱们进入直播间呀～"
	TplThankFollow = "谢谢 @%s 的关注呀～"
	TplThankGift   = "感谢 @%s 赠送的%d个%s～"
	TplSthWrong    = "不好意思呀，出现了一些状况～"
)

// Commands related templates.
const (
	TplIllegal    = "输入的参数不正确哦～"
	TplModeSwitch = "模式启停成功～"
	TplRankEmpty  = "今天的榜单好像空空的～"
	TplSongReq    = "点歌 %s 成功啦～"
	TplSongDone   = "完成了一首歌曲～"
)

// Modules related templates.
const (
	TplSignDuplicate = "已经签到过啦\n\n您已经连续签到%s天～\n\n%s"
	TplSignSuccess   = "签到成功啦，已经连续签到%d天～\n\n%s\n\n%s"
	TplStarFortune   = "%s今日幸运值：%s\n\n%s"
	TplPiaStart      = "pia戏模式启动成功啦～\n本文角色：%s\n请注意，如果输出文本时没有出现结果，很可能是存在违禁词，请使用重定位命令 r 定位到需要的位置，再通过 s 指令进行单条输出～"
	TplPiaEmpty      = "pia戏模式还没有启动哦，输入“帮助”来获取帮助吧～"
	TplPiaDone       = "本篇文章已经结束啦～"
	TplPiaOutBounds  = "选择的位置超过范围啦～"
	TplPiaRelocate   = "成功重定位啦～"
	TplPiaStop       = "pia戏模式结束成功啦～"
)

// Games related templates.
const (
	TplGameNull         = "当前还没有进行中的游戏哦～"
	TplGameExists       = "当前游戏未结束，请先结束～"
	TplGameCreate       = "游戏创建成功啦～"
	TplGameJoin         = "玩家 @%s 加入成功～"
	TplGameJoinDup      = "加入过的玩家不要再加入啦～"
	TplGameStart        = "游戏开始啦～"
	TplGameNext         = "请玩家 @%s 进行操作～"
	TplGameStop         = "停止游戏成功啦～"
	TplGameNoPlayer     = "当前还没有玩家哦～"
	TplGameNotEnough    = "人数不足，无法开始游戏～"
	TplGameInputIllegal = "拜托，输入正确的内容好嘛～"
	TplGameBomb         = "@%s 踩到了炸弹～"
	TplGameBombRange    = "炸弹在 [%d, %d] 区间内（包含边界）哦～"
	TplGameParcelNumber = "需要输入的数字：%s"
	TplGameParcelOver   = "击鼓传花结束啦，最后拿到的玩家是 @%s 小同学～"
	TplGameGuessWord    = "本轮游戏的词汇：%s"
	TplGameGuessStart   = "请该玩家查看私信，给出第一轮提示哦～"
	TplGameGuessTip     = "请该玩家给出下一轮提示哦～"
	TplGameGuessSuccess = "玩家 @%s 猜对啦～"
	TplGameGuessFailed  = "不好意思猜错啦～"
	TplGameRankEmpty    = "当前游戏榜单为空～"
)
