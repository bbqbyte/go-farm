package log

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"sync"
	"time"

	"keywea.com/cloud/pblib/pb/log"
)

type ioWriter struct {
	sync.Mutex
	buf bytes.Buffer
	writer io.Writer
}

func NewIoLogWriter(wr io.Writer) *ioWriter {
	return &ioWriter{writer: wr}
}

func (lg *ioWriter) println(logname, msg string, level log.Level, when time.Time, context, fields []log.Field) {
	lg.Lock()
	lg.buf.Reset()
	FormatHeader(&lg.buf, logname, when)
	lg.buf.WriteString(" ")
	lg.buf.WriteString(LevelPrefix[level])
	lg.writer.Write(lg.buf.Bytes())
	if len(msg) > 0 {
		lg.writer.Write([]byte{' '})
		lg.writer.Write([]byte(msg))
	}

	wr := ioEncoder{lg.writer}
	for _, f := range context {
		lg.writer.Write([]byte{' '})
		f.Encode(wr)
	}
	for _, f := range fields {
		lg.writer.Write([]byte{' '})
		f.Encode(wr)
	}
	lg.writer.Write([]byte{'\n'})
	lg.Unlock()
}


type ioEncoder struct {
	io.Writer
}

func (e ioEncoder) EncodeBool(key string, val bool) {
	fmt.Fprintf(e, "%s=%t", key, val)
}

func (e ioEncoder) EncodeFloat64(key string, val float64) {
	fmt.Fprintf(e, "%s=%f", key, val)
}

func (e ioEncoder) EncodeInt(key string, val int) {
	fmt.Fprintf(e, "%s=%d", key, val)
}

func (e ioEncoder) EncodeInt64(key string, val int64) {
	fmt.Fprintf(e, "%s=%d", key, val)
}

func (e ioEncoder) EncodeDuration(key string, val time.Duration) {
	fmt.Fprintf(e, "%s=%s", key, val)
}

func (e ioEncoder) EncodeUint(key string, val uint) {
	fmt.Fprintf(e, "%s=%d", key, val)
}

func (e ioEncoder) EncodeUint64(key string, val uint64) {
	fmt.Fprintf(e, "%s=%d", key, val)
}

func (e ioEncoder) EncodeString(key string, val string) {
	fmt.Fprintf(e, "%s=%q", key, val)
}

func (e ioEncoder) EncodeObject(key string, val interface{}) {
	fmt.Fprintf(e, "%s=%q", key, val)
}

func (e ioEncoder) EncodeType(key string, val reflect.Type) {
	fmt.Fprintf(e, "%s=%v", key, val)
}

