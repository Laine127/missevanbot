package util

import (
	"fmt"
	"os"
	"time"
)

type Logger struct {
	level int32 // 日志等级
	room  int   // 房间ID
	file  *os.File
}

const (
	LevelFatal int32 = iota
	LevelError
	LevelWarning
	LevelInfo
	LevelDebug
)

var _levelMap = map[string]int32{
	"fatal": LevelFatal,
	"error": LevelError,
	"warn":  LevelWarning,
	"info":  LevelInfo,
	"debug": LevelDebug,
}

// NewLogger 返回一个 Logger 对象
func NewLogger(str string, roomID int) (*Logger, error) {
	_ = os.Mkdir("logs", 0666) // 创建文件夹
	file, err := os.OpenFile(fmt.Sprintf("logs/%d.log", roomID), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	return &Logger{
		level: _levelMap[str],
		room:  roomID,
		file:  file,
	}, nil
}

// logPrint 根据日志等级判断是否输出日志到文件和标准输出
func (ll *Logger) logPrint(level int32, text string) {
	if level == -1 || level <= ll.level {
		_, _ = ll.file.WriteString(fmt.Sprintf("%s\n", text))
		fmt.Printf("[%d] %s\n", ll.room, text)
	}
}

// Print 直接输出日志，mode 为日志类型
func (ll *Logger) Print(mode, text string) {
	ll.logPrint(-1, fmt.Sprintf("[%s] [%s] %s", timeNow(), mode, text))
}

func (ll *Logger) Fatal(err error) {
	ll.logPrint(LevelFatal, fmt.Sprintf("[%s] [%s] %s", timeNow(), "FATAL", err.Error()))
	os.Exit(1)
}

func (ll *Logger) Error(err error) {
	ll.logPrint(LevelError, fmt.Sprintf("[%s] [%s] %s", timeNow(), "ERROR", err.Error()))
}

func (ll *Logger) Warning(err error) {
	ll.logPrint(LevelWarning, fmt.Sprintf("[%s] [%s] %s", timeNow(), "WARNING", err.Error()))
}

func (ll *Logger) Info(err error) {
	ll.logPrint(LevelInfo, fmt.Sprintf("[%s] [%s] %s", timeNow(), "INFO", err.Error()))
}

func (ll *Logger) Debug(err error) {
	ll.logPrint(LevelDebug, fmt.Sprintf("[%s] [%s] %s", timeNow(), "DEBUG", err.Error()))
}

func timeNow() string {
	return time.Now().Format("15:04:05")
}
