package log4go

import (
	"sync"
	"sync/atomic"
	"time"
)

var (
	ilog  *loggerManager
	ionce sync.Once
)

type exitFunc func(int)

func init() {
	ilog = &loggerManager{
		closed:    true,
		logIdSeed: 1000,
		appenders: make(map[string]IAppender),
		rootLogger: &Logger{
		},
		logIds: make(map[string]uint32),
		loggs:  make(map[uint32]ILogger),
	}
}

func GetLogger(name string, opts ...LoggerOption) ILogger {
	return ilog.getLogger(name, opts...)
}

func StartLogger() {
	ionce.Do(func() {
	})
}

type loggerManager struct {
	lock   sync.RWMutex
	closed bool

	logChanLen int64
	logChan    chan *Entry
	signalChan chan string

	entryPool  sync.Pool
	logIdSeed  uint32
	rootLogger *Logger
	logIds     map[string]uint32
	loggs      map[uint32]ILogger
	appenders  map[string]IAppender

	exitFn exitFunc
}

func (a *loggerManager) getLogger(name string, opts ...LoggerOption) ILogger {
	a.lock.RLock()
	l, ok := a.logIds[name]
	if ok {
		a.lock.RUnlock()
		return a.loggs[l]
	} else {
		a.lock.RUnlock()
		a.lock.Lock()
		if _, ok := a.logIds[name]; !ok {
			a.logIds[name] = atomic.AddUint32(&a.logIdSeed, 1)
			a.loggs[a.logIds[name]] = &Logger{
				Name:  name,
			}
			l = a.logIds[name]
		}
		a.lock.Unlock()
		return a.loggs[l]
	}
}

type logEvent struct {
	Logger  *Logger
	Data    Fields
	Time    time.Time
	Level   Level
	Message string
}
