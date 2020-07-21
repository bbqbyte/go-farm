package log

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type outputMode int

// DiscardNonColorEscSeq supports the divided color escape sequence.
// But non-color escape sequence is not output.
// Please use the OutputNonColorEscSeq If you want to output a non-color
// escape sequences such as ncurses. However, it does not support the divided
// color escape sequence.
const (
	_ outputMode = iota
	DiscardNonColorEscSeq
	OutputNonColorEscSeq
)

// NewAnsiColorWriter creates and initializes a new ansiColorWriter
// using io.Writer w as its initial contents.
// In the console of Windows, which change the foreground and background
// colors of the text by the escape sequence.
// In the console of other systems, which writes to w all text.
func NewAnsiColorWriter(w io.Writer) io.Writer {
	return NewModeAnsiColorWriter(w, DiscardNonColorEscSeq)
}

// NewModeAnsiColorWriter create and initializes a new ansiColorWriter
// by specifying the outputMode.
func NewModeAnsiColorWriter(w io.Writer, mode outputMode) io.Writer {
	if _, ok := w.(*ansiColorWriter); !ok {
		return &ansiColorWriter{
			w:    w,
			mode: mode,
		}
	}
	return w
}

type ansiColorWriter struct {
	w    io.Writer
	mode outputMode
}

func (cw *ansiColorWriter) Write(p []byte) (int, error) {
	return cw.w.Write(p)
}

var (
	green   = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	white   = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	yellow  = string([]byte{27, 91, 57, 55, 59, 52, 51, 109})
	red     = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	blue    = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	magenta = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	cyan    = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})

	w32Green   = string([]byte{27, 91, 52, 50, 109})
	w32White   = string([]byte{27, 91, 52, 55, 109})
	w32Yellow  = string([]byte{27, 91, 52, 51, 109})
	w32Red     = string([]byte{27, 91, 52, 49, 109})
	w32Blue    = string([]byte{27, 91, 52, 52, 109})
	w32Magenta = string([]byte{27, 91, 52, 53, 109})
	w32Cyan    = string([]byte{27, 91, 52, 54, 109})

	reset = string([]byte{27, 91, 48, 109})
)

// ColorByStatus return color by http code
// 2xx return Green
// 3xx return White
// 4xx return Yellow
// 5xx return Red
func ColorByStatus(cond bool, code int) string {
	switch {
	case code >= 200 && code < 300:
		return map[bool]string{true: green, false: w32Green}[cond]
	case code >= 300 && code < 400:
		return map[bool]string{true: white, false: w32White}[cond]
	case code >= 400 && code < 500:
		return map[bool]string{true: yellow, false: w32Yellow}[cond]
	default:
		return map[bool]string{true: red, false: w32Red}[cond]
	}
}

// ColorByMethod return color by http code
// GET return Blue
// POST return Cyan
// PUT return Yellow
// DELETE return Red
// PATCH return Green
// HEAD return Magenta
// OPTIONS return WHITE
func ColorByMethod(cond bool, method string) string {
	switch method {
	case "GET":
		return map[bool]string{true: blue, false: w32Blue}[cond]
	case "POST":
		return map[bool]string{true: cyan, false: w32Cyan}[cond]
	case "PUT":
		return map[bool]string{true: yellow, false: w32Yellow}[cond]
	case "DELETE":
		return map[bool]string{true: red, false: w32Red}[cond]
	case "PATCH":
		return map[bool]string{true: green, false: w32Green}[cond]
	case "HEAD":
		return map[bool]string{true: magenta, false: w32Magenta}[cond]
	case "OPTIONS":
		return map[bool]string{true: white, false: w32White}[cond]
	default:
		return reset
	}
}

// Guard Mutex to guarantee atomic of W32Debug(string) function
var mu sync.Mutex

// W32Debug Helper method to output colored logs in Windows terminals
func W32Debug(msg string) {
	mu.Lock()
	defer mu.Unlock()

	current := time.Now()
	w := NewAnsiColorWriter(os.Stdout)

	fmt.Fprintf(w, "[beego] %v %s\n", current.Format("2006/01/02 - 15:04:05"), msg)
}