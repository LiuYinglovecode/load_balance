package dao

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

// Config mysql config
type Config struct {
	DSN       string
	Active    int       // pool
	Idle      int       //pool
	IdleTime  time.Time // connect max life time
	LogEnable bool      // enable log sql query
	Logger    logrus.FieldLogger
}

// NewMysql from config
func NewMysql(c *Config) (db *gorm.DB) {
	db, err := Open(c)
	if err != nil {
		panic(err)
	}
	return
}

// Open mysql from config
func Open(c *Config) (db *gorm.DB, err error) {
	db, err = gorm.Open("mysql", c.DSN)
	if err != nil {
		c.Logger.Errorf("sql.Open() error(%v)", err)
		return nil, err
	}

	db.SetLogger(c.Logger)
	db.LogMode(c.LogEnable)
	db.DB().SetMaxOpenConns(c.Active)
	db.DB().SetMaxIdleConns(c.Idle)
	db.DB().SetConnMaxLifetime(time.Hour)
	return db, nil
}
