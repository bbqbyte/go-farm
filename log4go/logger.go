package log4go

import (
	"github.com/bbqbyte/go-farm/events"
	"sync/atomic"
)

type Logger struct {
	Name string

	option Op
}

func (logger *Logger) newEntry() *Entry {
	entry, ok := ilog.entryPool.Get().(*Entry)
	if ok {
		entry.setLogger(logger)
		return entry
	}
	return NewEntry(logger)
}

func (logger *Logger) releaseEntry(entry *Entry) {
	entry.Data = map[string]interface{}{}
	ilog.entryPool.Put(entry)
}

func (logger *Logger) WithFields(fields Fields) *Entry {
	entry := logger.newEntry()
	defer logger.releaseEntry(entry)
	return entry.WithFields(fields)
}

func (logger *Logger) WithError(err error) *Entry {
	entry := logger.newEntry()
	defer logger.releaseEntry(entry)
	return entry.WithError(err)
}

func (logger *Logger) log(level Level, args ...interface{}) {
	if logger.IsLevelEnabled(level) {
		entry := logger.newEntry()
		entry.log(level, args...)
		logger.releaseEntry(entry)
	}
}

func (logger *Logger) Print(args ...interface{}) {
	entry := logger.newEntry()
	entry.Print(args...)
	logger.releaseEntry(entry)
}

func (logger *Logger) Debug(args ...interface{}) {
	logger.log(LevelDebug, args...)
}

func (logger *Logger) Info(args ...interface{}) {
	logger.log(LevelInfo, args...)
}

func (logger *Logger) Warn(args ...interface{}) {
	logger.log(LevelWarn, args...)
}

func (logger *Logger) Error(args ...interface{}) {
	logger.log(LevelError, args...)
}

func (logger *Logger) Fatal(args ...interface{}) {
	logger.log(LevelFatal, args...)
	logger.Exit(1)
}

func (logger *Logger) Panic(args ...interface{}) {
	logger.log(LevelPanic, args...)
}

func (logger *Logger) logf(level Level, format string, args ...interface{}) {
	if logger.IsLevelEnabled(level) {
		entry := logger.newEntry()
		entry.logf(level, format, args...)
		logger.releaseEntry(entry)
	}
}

func (logger *Logger) Printf(format string, args ...interface{}) {
	entry := logger.newEntry()
	entry.Printf(format, args...)
	logger.releaseEntry(entry)
}

func (logger *Logger) Debugf(format string, args ...interface{}) {
	logger.logf(LevelDebug, format, args...)
}

func (logger *Logger) Infof(format string, args ...interface{}) {
	logger.logf(LevelInfo, format, args...)
}

func (logger *Logger) Warnf(format string, args ...interface{}) {
	logger.logf(LevelWarn, format, args...)
}

func (logger *Logger) Errorf(format string, args ...interface{}) {
	logger.logf(LevelError, format, args...)
}

func (logger *Logger) Fatalf(format string, args ...interface{}) {
	logger.logf(LevelFatal, format, args...)
	logger.Exit(1)
}

func (logger *Logger) Panicf(format string, args ...interface{}) {
	logger.logf(LevelPanic, format, args...)
}

func (logger *Logger) logln(level Level, args ...interface{}) {
	if logger.IsLevelEnabled(level) {
		entry := logger.newEntry()
		entry.logln(level, args...)
		logger.releaseEntry(entry)
	}
}

func (logger *Logger) Println(args ...interface{}) {
	entry := logger.newEntry()
	entry.Println(args...)
	logger.releaseEntry(entry)
}

func (logger *Logger) Debugln(args ...interface{}) {
	logger.logln(LevelDebug, args...)
}

func (logger *Logger) Infoln(args ...interface{}) {
	logger.logln(LevelInfo, args...)
}

func (logger *Logger) Warnln(args ...interface{}) {
	logger.logln(LevelWarn, args...)
}

func (logger *Logger) Errorln(args ...interface{}) {
	logger.logln(LevelError, args...)
}

func (logger *Logger) Fatalln(args ...interface{}) {
	logger.logln(LevelFatal, args...)
	logger.Exit(1)
}

func (logger *Logger) Panicln(args ...interface{}) {
	logger.logln(LevelPanic, args...)
}

// operation

func (logger *Logger) Exit(code int) {
	events.Exit(code)
}

func (logger *Logger) level() Level {
	return Level(atomic.LoadUint32((*uint32)(&logger.Level)))
}

func (logger *Logger) SetLevel(level Level) {
	atomic.StoreUint32((*uint32)(&logger.Level), uint32(level))
}

func (logger *Logger) GetLevel() Level {
	return logger.level()
}

func (logger *Logger) IsLevelEnabled(level Level) bool {
	return logger.level() <= level
}

func (logger *Logger) SetAddtivity() {
	atomic.StoreUint32(&logger.Addtivity, 1)
}

func (logger *Logger) IsAddtivity() bool {
	return atomic.LoadUint32(&logger.Addtivity) == 1
}
