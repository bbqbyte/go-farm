package rsa

import (
	"github.com/bbqbyte/go-farm/logger"
)

var (
	plog = log4go.New("[RSA]")
)

func SetLogLevel(level log4go.Level) {
	plog.SetLevel(level)
}
