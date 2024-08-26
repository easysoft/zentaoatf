package dao

import (
	"path/filepath"
	"time"

	"github.com/fatih/color"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"

	"github.com/easysoft/zentaoatf/internal/server/modules/v1/model"

	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	logUtils "github.com/easysoft/zentaoatf/pkg/lib/log"
)

var (
	db *gorm.DB
)

// GetDB 数据库单例
func GetDB() *gorm.DB {
	if db != nil {
		return db
	}

	conn := DBFile()
	dialector := sqlite.Open(conn)

	var err error
	db, err = gorm.Open(dialector, &gorm.Config{
		SkipDefaultTransaction: false,
		Logger:                 logger.Default.LogMode(logger.Info),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "",
			SingularTable: false,
		},
	})

	if err != nil {
		logUtils.Infof(color.RedString("open db failed, error: %s.", err.Error()))
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
		logUtils.Infof(color.RedString("migrate models failed, error: %s.", err.Error()))
	}

	return db
}

func DBFile() string {
	path := filepath.Join(commConsts.WorkDir, commConsts.App+".db")
	return path
}
