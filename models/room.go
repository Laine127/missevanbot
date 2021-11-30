package models

import (
	"time"

	"missevan-fm/config"
)

// Room 直播间
type Room struct {
	*config.RoomConfig             // 当前房间的配置
	Count              int         // 统计进入的数量
	Online             int         // 记录当前直播间在线人数
	Bait               bool        // 是否开启演员模式
	Timer              *time.Timer // 定时任务计时器
}
