package component

import "keywea.com/cloud/pblib/pb/log"

var (
	plog = log.New("[Component].pbcomponent")
)

func SetLogLevel(level log.Level) {
	plog.SetLevel(level)
}
