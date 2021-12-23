package models

import (
	"math/rand"
	"strings"
	"text/template"
	"time"
)

// handlers related templates.
const (
	TplWelcomeAnon = "欢迎新来的小可爱们进入直播间呀~"
	TplThankFollow = "谢谢 @%s 的关注呀~"
	TplThankGift   = "感谢 @%s 赠送的%d个%s~"
	TplSthWrong    = "不好意思呀，出现了一些状况"
)

// commands related templates.
const (
	TplModeSwitch = "模式启停成功"
	TplRankEmpty  = "今天的榜单好像空空的~"
	TplSongAdd    = "点歌 %s 成功啦~"
	TplSongDone   = "完成了一首歌曲~"
)

// modules related templates.
const (
	TplSignDuplicate = "已经签到过啦\n\n您已经连续签到%s天\n\n%s"
	TplSignSuccess   = "签到成功啦，已经连续签到%d天~\n\n%s\n\n%s"
	TplStarFortune   = "%s今日幸运值：%s\n\n%s"
	TplPiaStart      = "pia戏模式启动成功啦～\n本文角色：%s\n请注意，如果输出文本时没有出现结果，很可能是存在违禁词，请使用重定位命令 r 定位到需要的位置，再通过 s 指令进行单条输出～"
	TplPiaEmpty      = "pia戏模式还没有启动哦，输入“帮助”来获取帮助吧~"
	TplPiaDone       = "本篇文章已经结束啦～"
	TplPiaOutBounds  = "选择的位置超过范围啦"
	TplPiaRelocate   = "成功重定位啦～"
	TplPiaStop       = "pia戏模式结束成功啦～"
)

// games related templates.
const (
	TplGameNull         = "当前还没有进行中的游戏哦～"
	TplGameExists       = "当前游戏未结束，请先结束"
	TplGameCreate       = "游戏创建成功啦～"
	TplGameJoin         = "玩家 @%s 加入成功～"
	TplGameJoinDup      = "加入过的玩家不要再加入啦"
	TplGameStart        = "游戏开始啦～"
	TplGameNext         = "请玩家 @%s 进行操作"
	TplGameStop         = "停止游戏成功啦～"
	TplGameNoPlayer     = "当前还没有玩家哦～"
	TplGameNotEnough    = "人数不足，无法开始游戏～"
	TplGameInputIllegal = "拜托，输入正确的内容好嘛"
	TplGameBomb         = "@%s 踩到了炸弹～"
	TplGameBombRange    = "炸弹在 [%d, %d] 区间内（包含边界）哦～"
	TplGameParcelNumber = "需要输入的数字：%s"
	TplGameParcelOver   = "击鼓传花结束啦，最后拿到的玩家是 @%s 小同学～"
	TplGameGuessWord    = "本轮游戏的词汇：%s"
	TplGameGuessStart   = "请该玩家查看私信，给出第一轮提示哦～"
	TplGameGuessTip     = "请该玩家给出下一轮提示哦～"
	TplGameGuessSuccess = "玩家 @%s 猜对啦～"
	TplGameGuessFailed  = "不好意思猜错啦"
	TplGameRankEmpty    = "当前游戏榜单为空"
)

const (
	TmplBasePath = "templates/"
	TmplStartUp  = "startup.tmpl"
	TmplHelper   = "helper.tmpl"
	TmplModes    = "mode.tmpl"
	TmplRoomInfo = "room_info.tmpl"
	TmplPlaylist = "playlist.tmpl"
	TmplWelcome  = "welcome.tmpl"
)

const TplDefaultPoem = "孜孜不倦，不易乎世。"

var _replies = [...]string{
	"我来啦我来啦 ο(=•ω＜=)ρ⌒☆",
	"要抱抱 *(੭*ˊᵕˋ)੭*ଘ",
	"ヾ(≧▽≦*)o 是要和我一起玩嘛",
	"(｡･∀･)ﾉﾞ嗨",
	"(❤️´艸｀❤️)",
}

var _welcomes = [...]string{
	"(*ෆ´ ˘ `ෆ*)♡ 希望我们的相遇是彼此的幸运呀～",
	"要好好相处 ヾ(≧O≦)〃嗷~",
	"(`･ω･′)ゞ respect!",
	"ε≡٩(๑>₃<)۶ 我刚跑了三公里来欢迎你哦～",
	"( ｡ớ ₃ờ)ھ 可以为了我不要离开嘛～",
	"一起来聊天叭 (๑•ᴗ•๑)♡～",
}

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

// _chats contains simple chat text.
var _chats = [...]string{
	"你好呀～",
	"QAQ",
	"有什么事情嘛～",
	"不好意思呀我暂时没办法理你～",
	"抱歉现在有点忙～",
}

// _comforts contains comfort sentences.
var _comforts = [...]string{
	"不可以emo哦~",
	"怎么啦，为什么要说emo呢",
	"可以为了我不要伤心嘛",
	"有什么难过的事情可以告诉主播呀",
	"不开心的事就忘光光吧",
	"我不知道怎么安慰你才好，但是生活总得继续",
}

// _words contains words use for the game guess-word.
var _words = [...]string{
	"亡羊补牢", "颠三倒四", "拔苗助长", "画蛇添足", "顺手牵羊",
	"三长二短", "抱头鼠窜", "鸡鸣狗盗", "头破血流", "坐井观天",
	"眼高手低", "目瞪口呆", "胸无点墨", "鸡飞狗跳", "鼠目寸光",
	"盲人摸象", "画蛇添足", "画龙点睛", "抱头鼠窜", "狗急跳墙",
	"虎背熊腰", "守株待兔", "亡羊补牢", "对牛弹琴", "如鱼得水",
	"打草惊蛇", "打草惊蛇", "走马观花", "三头六臂", "丢三落四",
	"红杏出墙", "三头六臂", "对牛弹琴", "如鱼得水", "画龙点睛",
	"声东击西", "鸡飞狗跳", "鹿死谁手", "掩耳盗铃", "对牛弹琴",
	"藕断丝连", "可歌可泣", "眉飞色舞", "连蹦带跳", "左顾右盼",
	"嬉皮笑脸", "愁眉苦脸", "东倒西歪", "蹑手蹑脚", "喜出望外",
	"垂头丧气", "暴跳如雷", "狼吞虎咽", "见钱眼开", "摇头晃脑",
	"昂首挺胸", "捧腹大笑", "幸灾乐祸", "贼眉鼠眼", "牛头马面",
	"虎头蛇尾", "兔死狐悲", "龙腾虎跃", "狗急跳墙", "号啕大哭",
	"自言自语", "摇头晃脑", "撒腿就跑", "垂头丧气", "昂首挺胸",
	"手舞足蹈", "张牙舞爪", "眉开眼笑", "大惊小怪", "从容不迫",
	"目瞪口呆", "兴高采烈", "呆若木鸡", "幸灾乐祸", "神气十足",
	"唉声叹气", "哭笑不得", "捧腹大笑", "指手划脚", "东张西望",
	"一瘸一拐", "挤眉弄眼", "蹑手蹑脚", "废寝忘食", "闻鸡起舞",
	"守株待兔", "掩耳盗铃", "长吁短叹", "盲人摸象", "狼吞虎咽",
	"抓耳挠腮", "哈哈大笑", "恍然大悟", "一五一十", "三心二意",
	"争先恐后", "一刀两断", "丢三落四", "坐井观天", "快马加鞭",
}

func ReplyString() string {
	return _replies[randIndex(len(_replies))]
}

// WelcomeString return a welcome word in string type which chosen from _welcomes.
func WelcomeString() string {
	return _welcomes[randIndex(len(_welcomes))]
}

// LuckString return a luck word in string type which chosen from _lucks.
func LuckString() string {
	result := "今日运势："
	return result + _lucks[randIndex(len(_lucks))]
}

// ChatString return a sentence in string type which chosen from _chats.
func ChatString() string {
	return _chats[randIndex(len(_chats))]
}

// ComfortString return a sentence in string type which chosen from _comforts.
func ComfortString() string {
	return _comforts[randIndex(len(_comforts))]
}

// WordString return a word in string type which chosen from _words.
func WordString() string {
	return _words[randIndex(len(_words))]
}

func randIndex(length int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(length)
}

func NewTemplate(name string, data interface{}) (string, error) {
	path := TmplBasePath + name
	text := new(strings.Builder)
	funcs := template.FuncMap{"add": add}
	tmpl := template.Must(template.New(name).Funcs(funcs).ParseFiles(path)) // hello.tmpl
	if err := tmpl.Execute(text, data); err != nil {
		return "", err
	}
	return text.String(), nil
}

func add(x, y int) int {
	return x + y
}
