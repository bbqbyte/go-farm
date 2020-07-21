package pbapp

import "keywea.com/cloud/pblib/pb/log"

var (
	plog = log.New("[APP].pbapp")
)

func SetLogLevel(level log.Level) {
	plog.SetLevel(level)
}
