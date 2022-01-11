package dao

import (
	"github.com/aaronchen2k/deeptest/internal/comm/consts"
	"github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	"github.com/aaronchen2k/deeptest/internal/server/config"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
	"path/filepath"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db *gorm.DB
)

// GetDB 数据库单例
func GetDB() *gorm.DB {
	conn := DBFile()
	dialector := sqlite.Open(conn)

	db, err := gorm.Open(dialector, &gorm.Config{
		SkipDefaultTransaction: false,
		Logger:                 logger.Default.LogMode(logger.Info),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "",
			SingularTable: false,
		},
	})

	if err != nil {
		logUtils.Info(err.Error())
	}

	_ = db.Use(
		dbresolver.Register(
			dbresolver.Config{ /* xxx */ }).
			SetConnMaxIdleTime(time.Hour).
			SetConnMaxLifetime(24 * time.Hour).
			SetMaxIdleConns(100).
			SetMaxOpenConns(200),
	)

	db.Session(&gorm.Session{FullSaveAssociations: true, AllowGlobalUpdate: false})

	err = db.AutoMigrate(
		model.Models...,
	)
	if err != nil {
		logUtils.Info(err.Error())
	}

	return db
}

func DBFile() string {
	path := filepath.Join(serverConfig.CONFIG.System.WorkDir, commConsts.App+".db")
	return path
}
