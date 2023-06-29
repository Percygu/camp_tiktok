package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
	"time"
	"videosvr/config"
	"videosvr/log"
)

var (
	db     *gorm.DB
	dbOnce sync.Once
)

// openDB 连接db
func openDB() {
	dbConfig := config.GetGlobalConfig().DbConfig
	connArgs := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", dbConfig.Username,
		dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Database)
	// log.Info("mdb addr:" + connArgs)

	var err error
	db, err = gorm.Open(mysql.Open(connArgs), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic("fetch db connection err:" + err.Error())
	}

	sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConn)                                        // 设置最大空闲连接
	sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConn)                                        // 设置最大打开的连接
	sqlDB.SetConnMaxLifetime(time.Duration(dbConfig.MaxIdleTime * int64(time.Second))) // 设置空闲时间为(s)
}

// GetDB 获取数据库连接
func GetDB() *gorm.DB {
	dbOnce.Do(openDB)
	return db
}

func CloseDB() {
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			log.Errorf("close db err:%v", err)
		}
		sqlDB.Close()
	}
}
