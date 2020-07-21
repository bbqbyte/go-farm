package log

import (
	"keywea.com/cloud/pblib/pbconfig"
	"keywea.com/cloud/pblib/pb/log"
	"time"
)

const (
	AdapterConsole   = "console"
	AdapterFile      = "file"
	AdapterRemote    = "remote"
)

type newLoggerFunc func() ILogger

type ILogger interface {
	Init(configor pbconfig.Configor) error
	SetLogLevel(level log.Level)
	WriteLog(logname, msg string, level log.Level, when time.Time, context, fields []log.Field)
	Destroy()
}

var (
	adapters = make(map[string]newLoggerFunc)
	LevelPrefix = [log.LevelFatal + 1]string{"[D]", "[I]", "[W]", "[E]", "[P]", "[F]"}
)

// 注册Adapter，如console/file/remote
func Register(adaptername string, adapterFn newLoggerFunc) {
	if adapterFn == nil {
		panic("pblog: Register provider is nil")
	}
	if _, dup := adapters[adaptername]; dup {
		panic("pblog: Register duplicate for provider " + adaptername)
	}
	adapters[adaptername] = adapterFn
}

type LogWrap struct {
	Name string // logger name
	OutWriter string // 指定输出
	Level log.Level
	Msg   string
	Context  []log.Field
	Fields  []log.Field
	When  time.Time
}