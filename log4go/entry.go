package log4go

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"sync"
	"time"
)

var (
	bufferPool *sync.Pool

	callDepth int

	ErrorKey = "error"
)

const (
	defaultCallDepth int = 4

	_log4goPackage = "github.com/bbqbyte/go-farm/log4go"
)

func init() {
	bufferPool = &sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}

	callDepth = defaultCallDepth
}

type Entry struct {
	Logger *Logger

	Data  Fields
	Time  time.Time
}

func NewEntry(logger *Logger) *Entry {
	return &Entry{
		Logger: logger,
		Data:   make(Fields, 6),
	}
}

func (entry *Entry) setLogger(logger *Logger) {
	entry.Logger = logger
}

func (entry *Entry) WithFields(fields Fields) *Entry {
	data := make(Fields, len(entry.Data)+len(fields))
	for k, v := range entry.Data {
		data[k] = v
	}
	fieldErr := entry.err
	for k, v := range fields {
		isErrField := false
		if t := reflect.TypeOf(v); t != nil {
			switch t.Kind() {
			case reflect.Func:
				isErrField = true
			case reflect.Ptr:
				isErrField = t.Elem().Kind() == reflect.Func
			}
		}
		if isErrField {
			tmp := fmt.Sprintf("can not add field %q", k)
			if fieldErr != "" {
				fieldErr = entry.err + ", " + tmp
			} else {
				fieldErr = tmp
			}
		} else {
			data[k] = v
		}
	}
	return &Entry{Logger: entry.Logger, Data: data, Time: entry.Time}
}

func (entry *Entry) WithError(err error) *Entry {
	return entry.WithFields(Fields{ErrorKey: err})
}

func (entry Entry) log1(level Level, msg string) {
	var buffer *bytes.Buffer

	entry.Time = time.Now()

	entry.Level = level
	entry.Message = msg
	entry.Logger.mu.Lock()
	if entry.Logger.ReportCaller {
		entry.Caller = getCaller()
	}
	entry.Logger.mu.Unlock()

	buffer = bufferPool.Get().(*bytes.Buffer)
	buffer.Reset()
	defer bufferPool.Put(buffer)
	entry.Buffer = buffer

	entry.write()

	entry.Buffer = nil
}

func (entry *Entry) write() {
	entry.Logger.mu.Lock()
	defer entry.Logger.mu.Unlock()
	serialized, err := entry.Logger.Formatter.Format(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to obtain reader, %v\n", err)
		return
	}
	if _, err = entry.Logger.Out.Write(serialized); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write to log, %v\n", err)
	}
}

func (entry *Entry) log(level Level, args ...interface{}) {
	if entry.Logger.IsLevelEnabled(level) {
		//entry.log(level, fmt.Sprint(args...))
	}
}

func (entry *Entry) Print(args ...interface{}) {
	entry.Info(args...)
}

func (entry *Entry) Debug(args ...interface{}) {
	entry.log(LevelDebug, args...)
}

func (entry *Entry) Info(args ...interface{}) {
	entry.log(LevelInfo, args...)
}

func (entry *Entry) Warn(args ...interface{}) {
	entry.log(LevelWarn, args...)
}

func (entry *Entry) Error(args ...interface{}) {
	entry.log(LevelError, args...)
}

func (entry *Entry) Fatal(args ...interface{}) {
	entry.log(LevelFatal, args...)
	entry.Logger.Exit(1)
}

func (entry *Entry) Panic(args ...interface{}) {
	entry.log(LevelPanic, args...)
	panic(fmt.Sprint(args...))
}

// Entry Printf family functions

func (entry *Entry) logf(level Level, format string, args ...interface{}) {
	if entry.Logger.IsLevelEnabled(level) {
		entry.log(level, fmt.Sprintf(format, args...))
	}
}

func (entry *Entry) Printf(format string, args ...interface{}) {
	entry.Infof(format, args...)
}

func (entry *Entry) Debugf(format string, args ...interface{}) {
	entry.logf(LevelDebug, format, args...)
}

func (entry *Entry) Infof(format string, args ...interface{}) {
	entry.logf(LevelInfo, format, args...)
}

func (entry *Entry) Warnf(format string, args ...interface{}) {
	entry.logf(LevelWarn, format, args...)
}

func (entry *Entry) Errorf(format string, args ...interface{}) {
	entry.logf(LevelError, format, args...)
}

func (entry *Entry) Fatalf(format string, args ...interface{}) {
	entry.logf(LevelFatal, format, args...)
	entry.Logger.Exit(1)
}

func (entry *Entry) Panicf(format string, args ...interface{}) {
	entry.logf(LevelPanic, format, args...)
}

// Entry Println family functions

func (entry *Entry) logln(level Level, args ...interface{}) {
	if entry.Logger.IsLevelEnabled(level) {
		entry.log(level, entry.sprintlnn(args...))
	}
}

func (entry *Entry) Println(args ...interface{}) {
	entry.Infoln(args...)
}

func (entry *Entry) Debugln(args ...interface{}) {
	entry.logln(LevelDebug, args...)
}

func (entry *Entry) Infoln(args ...interface{}) {
	entry.logln(LevelInfo, args...)
}

func (entry *Entry) Warnln(args ...interface{}) {
	entry.logln(LevelWarn, args...)
}

func (entry *Entry) Errorln(args ...interface{}) {
	entry.logln(LevelError, args...)
}

func (entry *Entry) Fatalln(args ...interface{}) {
	entry.logln(LevelFatal, args...)
	entry.Logger.Exit(1)
}

func (entry *Entry) Panicln(args ...interface{}) {
	entry.logln(LevelPanic, args...)
}

func (entry *Entry) sprintlnn(args ...interface{}) string {
	msg := fmt.Sprintln(args...)
	return msg[:len(msg)-1]
}
