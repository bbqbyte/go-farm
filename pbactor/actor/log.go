package actor

import "keywea.com/cloud/pblib/pb/log"

var (
	plog = log.New("[ACTOR].pbactor")
)

func SetLogLevel(level log.Level) {
	plog.SetLevel(level)
}
