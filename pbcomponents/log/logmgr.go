package log

import (
	"fmt"
	"keywea.com/cloud/pblib/pbconfig"
	"keywea.com/cloud/pblib/pb/events"
	"keywea.com/cloud/pblib/pb/log"
	"sync"
	"time"
)

var (
	msgObjPool *sync.Pool

	ilog *logPane
	ionce sync.Once
)

type logPane struct {
	lock           sync.Mutex
	inited         bool
	closed 		   bool
	logChanLen     int64
	logChan        chan *LogWrap
	signalChan     chan string
	wg             sync.WaitGroup
	defaultWriter  string // 默认log writer
	writers    	   map[string]*PBLogWriter // logname : writer
}

type PBLogWriter struct {
	name string
	adapter string
	level log.Level
	writer ILogger

	wmu sync.Mutex
	closed bool
}

func NewLogWriter(logWriterName string, configor pbconfig.Configor) (*PBLogWriter, error) {
	ionce.Do(func() {
		if ilog == nil {
			ilog = &logPane{}
			ilog.InitLog(configor)
			events.AddShutdownHook(func() error {
				ilog.Close()
				return nil
			}, events.SHUTDOWN_INDEX_LOG)
			log.SetLogFunc(ilog.Publish)
		}
	})

	return ilog.SetLogWriter(logWriterName, configor)
}

func (sl *logPane) InitLog(configor pbconfig.Configor) {
	sl.lock.Lock()
	defer sl.lock.Unlock()
	if sl.inited {
		return
	}
	sl.closed = false

	sl.logChanLen, _ = configor.GetInt64("chanlen", 1e3)
	sl.logChan = make(chan *LogWrap, sl.logChanLen)
	msgObjPool = &sync.Pool{
		New: func() interface{} {
			return &LogWrap{}
		},
	}
	sl.writers = make(map[string]*PBLogWriter)

	sl.signalChan = make(chan string, 1)

	sl.wg.Add(1)
	go sl.startLog()

	sl.inited = true
}

func (sl *logPane) startLog() {
	closed := false

	for {
		select {
		case log := <-sl.logChan:
			// send to logger instance
			sl.pushToWriters(log.Name, log.OutWriter, log.Msg, log.Level, log.When, log.Context, log.Fields)
			// return to pool
			msgObjPool.Put(log)
		case sig := <-sl.signalChan:
			sl.flush()
			if sig == "close" {
				for _, l := range sl.writers {
					l.Destroy()
				}
				sl.writers = nil
				closed = true
				sl.wg.Done()
			}
		}
		if closed {
			break
		}
	}
}

func (sl *logPane) SetLoggerLevel(writerName string, level log.Level) {
	if conf, ok := sl.writers[writerName]; ok {
		conf.level = level
		sl.writers[writerName].SetLogLevel(level)
	}
}

func (sl *logPane) SetLogWriter(logWriterName string, configor pbconfig.Configor) (*PBLogWriter, error) {
	sl.lock.Lock()
	defer sl.lock.Unlock()

	if _, ok := sl.writers[logWriterName]; ok {
		return nil, fmt.Errorf("pblog: duplicate log writer name %q (you have set this log writer before)", logWriterName)
	}

	adapterName := configor.GetString("adapter", AdapterConsole)
	logAdapter, ok := adapters[adapterName]
	if !ok {
		return nil, fmt.Errorf("pblog: unknown adaptername %q (forgotten Register?)", adapterName)
	}
	logInst := logAdapter()
	err := logInst.Init(configor)
	if err != nil {
		return nil, err
	}

	level, _ := configor.GetInt("level", log.LevelDebug)
	isDefaultWriter, _ := configor.GetBool("default")

	sl.writers[logWriterName] = &PBLogWriter{
		name: logWriterName,
		adapter: adapterName,
		level: log.Level(level),
		writer: logInst,
	}

	if isDefaultWriter {
		sl.defaultWriter = logWriterName
	}
	return sl.writers[logWriterName], nil
}

func (sl *logPane) Publish(l *log.Logger, level log.Level, msg string, fields []log.Field) {
	if !sl.inited || sl.closed {
		return
	}
	log := msgObjPool.Get().(*LogWrap)
	log.When=  time.Now()
	log.Level = level
	log.Name = l.Name()
	log.OutWriter = l.OutWriter()
	log.Msg = msg
	log.Context = l.Context()
	log.Fields = fields
	sl.logChan <- log
}

func (sl *logPane) pushToWriters(logname, outWriter, msg string, level log.Level, when time.Time, context, fields []log.Field) {
	if outWriter == "" {
		outWriter = sl.defaultWriter
	}
	if outWriter != "" && sl.writers[outWriter] != nil {
		sl.writers[outWriter].WriteLog(logname, msg, level, when, context, fields)
	}
}

func (sl *logPane) flush() {
	for {
		if len(sl.logChan) > 0 {
			log := <-sl.logChan
			sl.pushToWriters(log.Name, log.OutWriter, log.Msg, log.Level, log.When, log.Context, log.Fields)
			msgObjPool.Put(log)
			continue
		}
		break
	}
}

func (sl *logPane) Close() {
	if !sl.inited || sl.closed {
		return
	}
	sl.closed = true
	sl.signalChan <- "close"
	sl.wg.Wait()
	close(sl.logChan)

	close(sl.signalChan)
}

// log writer
func (lw *PBLogWriter) SetLogLevel(level log.Level) {
	lw.wmu.Lock()
	defer lw.wmu.Unlock()
	lw.level = level
}

func (lw *PBLogWriter) Name() string {
	return lw.name
}

func (lw *PBLogWriter) Adapter() string {
	return lw.adapter
}

func (lw *PBLogWriter) LogLevel() log.Level {
	return lw.level
}

func (lw *PBLogWriter) Destroy() {
	lw.wmu.Lock()
	defer lw.wmu.Unlock()
	if lw.closed {
		return
	}
	lw.closed = true
	lw.writer.Destroy()
}

func (lw *PBLogWriter) WriteLog(logname, msg string, level log.Level, when time.Time, context, fields []log.Field) {
	lw.writer.WriteLog(logname, msg, level, when, context, fields)
}
