package db

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"keywea.com/cloud/pblib/pbconfig"
	"keywea.com/cloud/pblib/pb/events"
	"sync"
	"time"
)

var (
	errNotFoundDatabaseSource = errors.New("database source Not Found")

	pbdb *dss
	dbmu sync.RWMutex
)

type DsConf struct {
	DriverName     	  string
	DataSourceName 	  string
	MaxOpenConns      int
	MaxIdleConns      int
	ConnMaxLifetime	  time.Duration // idletimeout
}

type dss struct {
	datasources map[string]*DS
	dsConfs map[string]DsConf

	mu sync.Mutex
}

type DS struct {
	dsname string
	db *sqlx.DB

	closed bool
	dsmu sync.Mutex
}

func NewDB(name string, configor pbconfig.Configor) (*DS, error) {
	dbmu.Lock()
	if pbdb == nil {
		pbdb = &dss{
			datasources: make(map[string]*DS),
			dsConfs: make(map[string]DsConf),
		}
		events.AddShutdownHook(func() error {
			pbdb.Destroy()
			return nil
		}, events.SHUTDOWN_INDEX_DB)
	}
	dbmu.Unlock()
	return pbdb.createDS(name, configor)
}

func (dm *dss) createDS(name string, configor pbconfig.Configor) (*DS, error) {
	dm.mu.Lock()
	defer dm.mu.Unlock()
	if configor == nil {
		return nil, fmt.Errorf("database source=%s create Error on nil configor", name)
	}
	config := dm.parseConfig(configor)
	if db, ok := dm.datasources[name]; ok { // datasource exists
		return db.UpdatePool(name, config)
	}
	return dm.create(name, config)
}

func (dm *dss) create(dsname string, conf DsConf) (*DS, error) {
	db, err := sqlx.Connect(conf.DriverName, conf.DataSourceName)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(conf.MaxOpenConns)
	db.SetMaxIdleConns(conf.MaxIdleConns)
	if conf.ConnMaxLifetime > 0 {
		db.SetConnMaxLifetime(conf.ConnMaxLifetime)
	}
	dm.datasources[dsname] = &DS{
		dsname: dsname,
		db: db,
	}
	dm.dsConfs[dsname] = conf
	return dm.datasources[dsname], nil
}

func (dm *dss) parseConfig(configor pbconfig.Configor) DsConf {
	driverName := configor.GetString("driverName", "mysql")
	dataSourceName := configor.GetString("dataSourceName", "")
	if driverName == "" || dataSourceName == "" {
		panic("parse database conf fatal")
	}
	maxOpenConns, _ := configor.GetInt("maxOpenConns", 50)
	maxIdleConns, _ := configor.GetInt("maxIdleConns", 50)
	connMaxLifetime, _ := configor.GetInt("connMaxLifetime", 0)

	return DsConf{
		DriverName: driverName,
		DataSourceName: dataSourceName,
		MaxOpenConns: maxOpenConns,
		MaxIdleConns: maxIdleConns,
		ConnMaxLifetime: time.Second*time.Duration(connMaxLifetime),
	}
}

func (dm *dss) Get(name string) *DS {
	p, ok := dm.datasources[name]
	if !ok {
		panic(errNotFoundDatabaseSource)
	}
	return p
}

func (dm *dss) Destroy() {
	for _, v := range dm.datasources {
		err := v.Destroy()
		if err != nil {
		}
	}
}

// db
func (db *DS) UpdatePool(dsname string, conf DsConf) (*DS, error) {
	db.dsmu.Lock()
	defer db.dsmu.Unlock()

	oldConfig := pbdb.dsConfs[dsname]

	if oldConfig.DriverName != conf.DriverName || oldConfig.DataSourceName != conf.DataSourceName {
		err := db.Destroy()
		if err != nil {
		}
		delete(pbdb.datasources, dsname)
		delete(pbdb.dsConfs, dsname)
		return pbdb.create(dsname, conf)
	}

	if conf.MaxOpenConns > 0 && oldConfig.MaxOpenConns != conf.MaxOpenConns {
		db.db.SetMaxOpenConns(conf.MaxOpenConns)
	}
	if conf.MaxIdleConns > 0 && oldConfig.MaxIdleConns != conf.MaxIdleConns {
		db.db.SetMaxIdleConns(conf.MaxIdleConns)
	}
	if conf.ConnMaxLifetime > 0 && oldConfig.ConnMaxLifetime != conf.ConnMaxLifetime {
		db.db.SetConnMaxLifetime(conf.ConnMaxLifetime)
	}

	delete(pbdb.dsConfs, dsname)
	pbdb.dsConfs[dsname] = conf

	return db, nil
}

func (db *DS) ParseConfig(configor *pbconfig.Configor) DsConf {
	return pbdb.parseConfig(*configor)
}

func (db *DS) Ds() *sqlx.DB {
	return db.db
}

func (db *DS) Destroy() error {
	db.dsmu.Lock()
	defer db.dsmu.Unlock()

	if db.closed {
		return nil
	}
	db.closed = true
	return db.db.Close()
}