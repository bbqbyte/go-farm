package pbapp

import (
	"keywea.com/cloud/pblib/pb/component"
	"keywea.com/cloud/pblib/pbcomponents/log"
	"keywea.com/cloud/pblib/pbcomponents/storage/redis"
	"keywea.com/cloud/pblib/pbcomponents/storage/db"
)

var (
	app *component.PBC
)

func Create() *component.PBC {
	app = component.NewPBC()
	return app
}

func GetApp() *component.PBC {
	return app
}

func GetRedis(poolname string) *redis.RPool {
	return app.Instance(poolname).(*redis.RPool)
}

func GetLogWriter(name string) *log.PBLogWriter {
	return app.Instance(name).(*log.PBLogWriter)
}

func GetDatabaseSource(dsname string) *db.DS {
	return app.Instance(dsname).(*db.DS)
}