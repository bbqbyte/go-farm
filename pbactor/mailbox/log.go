package mailbox

import "keywea.com/cloud/pblib/pb/log"

var (
	plog = log.New("[MAILBOX.pbactor]")
)

func SetLogLevel(level log.Level) {
	plog.SetLevel(level)
}
