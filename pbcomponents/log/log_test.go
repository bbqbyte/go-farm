package log

import (
	"fmt"
	"keywea.com/cloud/pblib/pbapp"
	"keywea.com/cloud/pblib/pbconfig"
	"keywea.com/cloud/pblib/pb/log"
	"testing"
)

func TestLog(t *testing.T) {
	s := "{\"level\": 0, \"default\":true, \"adapter\":\"file\", \"filename\":\"/volumes/D/x.log\"}"
	c, err := pbconfig.NewConfigData("json", []byte(s))
	if err != nil {
		t.Fatal("parse json error")
	} else {
		l, e := NewLogWriter("default", c)
		if e != nil {
			fmt.Sprint(e)
		}
		logf := log.New("[testlog]")
		logf.Info("show me the money", log.String("x", "aa"), log.Stack())
		logf.Error("I'm dying", log.String("x", "bb"), log.Stack())
		logf.Fatal("I'm died", log.String("x", "cc"), log.Stack())
		pbapp.Wait()
		l.Destroy()
	}
}
