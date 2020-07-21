package log

import (
	"keywea.com/cloud/pblib/pbconfig"
	"keywea.com/cloud/pblib/pb/log"
	"os"
	"runtime"
	"time"
)

// brush is a color join function
type brush func(string) string

// newBrush return a fix color Brush
func newBrush(color string) brush {
	pre := "\033["
	reset := "\033[0m"
	return func(text string) string {
		return pre + color + "m" + text + reset
	}
}

var colors = []brush{
	newBrush("1;37"), // Fatal          white
	newBrush("1;35"), // Panic           magenta
	newBrush("1;31"), // Error              red
	newBrush("1;33"), // Warning            yellow
	newBrush("1;34"), // Informational      blue
	newBrush("1;44"), // Debug              Background blue
}

// consoleWriter implements LoggerInterface and writes messages to terminal.
type consoleWriter struct {
	lg       *ioWriter
	level    log.Level
	Colorful bool `json:"color"` //this filed is useful only when system's terminal supports color
}

// NewConsole create ConsoleWriter returning as LoggerInterface.
func NewConsole() ILogger {
	cw := &consoleWriter{
		lg:       NewIoLogWriter(os.Stdout),
		level:	  log.LevelDebug,
		Colorful: runtime.GOOS != "windows",
	}
	return cw
}

// Init init console logger.
// jsonConfig like '{"level":LevelTrace}'.
func (c *consoleWriter) Init(configor pbconfig.Configor) error {
	if configor != nil {
		c.Colorful, _ = configor.GetBool("color")
		level, _ := configor.GetInt("level", log.LevelDebug)
		c.level = log.Level(level)
	}
	if runtime.GOOS == "windows" {
		c.Colorful = false
	}

	return nil
}

func (c *consoleWriter) SetLogLevel(level log.Level) {
	c.level = level
}

// WriteMsg write message in console.
func (c *consoleWriter) WriteLog(logname, msg string, level log.Level, when time.Time, context, fields []log.Field) {
	if c.level > level {
		return
	}
	if c.Colorful {
		msg = colors[level](msg)
	}
	c.lg.println(logname, msg, level, when, context, fields)
}

// Destroy implementing method. empty.
func (c *consoleWriter) Destroy() {

}

func init() {
	Register(AdapterConsole, NewConsole)
}