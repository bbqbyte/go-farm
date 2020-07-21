package log

const (
	LevelDebug = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelPanic
	LevelFatal
	LevelOff // 日志关闭
)

type (
	Level int

	writeLogFunc func(l *Logger, level Level, msg string, fields []Field)
)

var (
	DefaultLogLevel Level = LevelDebug
	DefaultCallDepth = 4
	logSwitch int32 = 1 // 总开关，1-开，0-关

	loggers []*Logger

	logFunc writeLogFunc = func(l *Logger, level Level, msg string, fields []Field) {
	}
)

type Logger struct {
	name string
	outWriter string // 指定输出
	level Level
	context  []Field
}


func New(name string, context ...Field) *Logger {
	r := &Logger{name: name, outWriter: "", level: DefaultLogLevel, context: context}
	loggers = append(loggers, r)
	return r
}

func (l *Logger) Name() string {
	return l.name
}

func (l *Logger) OutWriter() string {
	return l.outWriter
}

func (l *Logger) Context() []Field {
	return l.context
}

func (l *Logger) SetOutWriter(writername string) {
	l.outWriter = writername
}

func (l *Logger) Level() Level {
	return l.level
}

func (l *Logger) SetLevel(level Level) {
	l.level = level
}

// 设置所有Logger的level
func SetLevel(level Level) {
	for _, logger := range loggers {
		logger.SetLevel(level)
	}
}

func SetCallDepth(depth int) {
	DefaultCallDepth = depth
}

// 总开关
func SwitchLog(f int32) {
	logSwitch = f
}

func SetLogFunc(f writeLogFunc) {
	logFunc = f
}

func (l *Logger) Fatal(msg string, fields ...Field) {
	if l.Level() <= LevelFatal && logSwitch != 0 {
		logFunc(l, LevelFatal, msg, fields)
	}
}

func (l *Logger) Panic(msg string, fields ...Field) {
	if l.Level() <= LevelPanic && logSwitch != 0  {
		logFunc(l, LevelPanic, msg, fields)
	}
}

func (l *Logger) Error(msg string, fields ...Field) {
	if l.Level() <= LevelError && logSwitch != 0  {
		logFunc(l, LevelError, msg, fields)
	}
}

func (l *Logger) Warn(msg string, fields ...Field) {
	if l.Level() <= LevelWarn && logSwitch != 0  {
		logFunc(l, LevelWarn, msg, fields)
	}
}

func (l *Logger) Info(msg string, fields ...Field) {
	if l.Level() <= LevelInfo && logSwitch != 0  {
		logFunc(l, LevelInfo, msg, fields)
	}
}

func (l *Logger) Debug(msg string, fields ...Field) {
	if l.Level() <= LevelDebug && logSwitch != 0  {
		logFunc(l, LevelDebug,  msg, fields)
	}
}
