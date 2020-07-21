package plugin

import "keywea.com/cloud/pblib/pb/log"

var (
	plog = log.New("[PLUGIN.pbactor]")
)

func SetLogLevel(level log.Level) {
	plog.SetLevel(level)
}
