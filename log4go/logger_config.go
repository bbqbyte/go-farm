package log4go

type LoggerOption func(*Op)

type Op struct {
	addtivity bool
	level     Level
	appenders []string
}

var defaultOption = Op{
	level:     LevelTrace,
	addtivity: false,
}

func (op *Op) applyOpts(opts []LoggerOption) {
	for _, opt := range opts {
		opt(op)
	}
}

func WithLevel(level Level) LoggerOption {
	return func(op *Op) { op.level = level }
}

func WithAddtivity() LoggerOption {
	return func(op *Op) { op.addtivity = true }
}