package log4go

import "time"

const (
	defaultTimestampFormat = time.RFC3339
	FieldKeyMsg            = "msg"
	FieldKeyLevel          = "level"
	FieldKeyTime           = "time"
	FieldKeyError          = "err"
	FieldKeyFunc           = "func"
	FieldKeyFile           = "file"
)

type Layout interface {
	Format(*Entry) ([]byte, error)
}
