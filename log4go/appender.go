package log4go

type IAppender interface {
	Write()
	Close()
}

type Appender struct {
	Name string
	Layout Layout
}
