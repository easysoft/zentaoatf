package repo

import (
	"database/sql"
	"fmt"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"gorm.io/gorm"
)

type DataRepo struct {
	DB *gorm.DB `inject:""`
}

func NewDataRepo() *DataRepo {
	return &DataRepo{}
}

// CreateTable 创建数据库(mysql)
func (s *DataRepo) CreateTable(reqSql serverDomain.DataDb) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/", reqSql.UserName, reqSql.Password, reqSql.Host, reqSql.Port)
	createSql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_general_ci;", reqSql.DBName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)
	if err = db.Ping(); err != nil {
		return err
	}
	_, err = db.Exec(createSql)
	return err
}
