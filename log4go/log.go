package log4go

const (
	LevelTrace = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
	LevelPanic
	LevelOff // 日志关闭
)

type (
	Level uint32

	Fields map[string]interface{}
)

type ILogBase interface {
	Printf(format string, args ...interface{})
	Print(args ...interface{})
	Println(args ...interface{})

	Debugf(format string, args ...interface{})
	Debug(args ...interface{})
	Debugln(args ...interface{})

	Infof(format string, args ...interface{})
	Info(args ...interface{})
	Infoln(args ...interface{})

	Warnf(format string, args ...interface{})
	Warn(args ...interface{})
	Warnln(args ...interface{})

	Errorf(format string, args ...interface{})
	Error(args ...interface{})
	Errorln(args ...interface{})

	Fatalf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalln(args ...interface{})

	Panicf(format string, args ...interface{})
	Panic(args ...interface{})
	Panicln(args ...interface{})
}

type ILogger interface {
	ILogBase

	WithFields(fields Fields) *Entry
	WithError(err error) *Entry
}
