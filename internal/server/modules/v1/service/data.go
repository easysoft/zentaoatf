package service

import (
	"errors"
	"fmt"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	"github.com/aaronchen2k/deeptest/internal/server/config"
	"github.com/aaronchen2k/deeptest/internal/server/consts"
	"github.com/aaronchen2k/deeptest/internal/server/core/cache"
	"github.com/aaronchen2k/deeptest/internal/server/core/module"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/model"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/repo"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/source"
	"golang.org/x/crypto/bcrypt"

	"github.com/snowlyg/helper/str"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	ErrViperEmpty = errors.New("配置服务未初始化")
)

type DataService struct {
	DataRepo   *repo.DataRepo     `inject:""`
	UserRepo   *repo.UserRepo     `inject:""`
	UserSource *source.UserSource `inject:""`
	RoleSource *source.RoleSource `inject:""`
	PermSource *source.PermSource `inject:""`
}

func NewDataService() *DataService {
	return &DataService{}
}

// writeConfig 写入配置文件
func (s *DataService) writeConfig(viper *viper.Viper, conf serverConfig.Config) error {
	cs := str.StructToMap(serverConsts.CONFIG)
	for k, v := range cs {
		viper.Set(k, v)
	}
	return viper.WriteConfig()
}

// 回滚配置
func (s *DataService) refreshConfig(viper *viper.Viper, conf serverConfig.Config) error {
	err := s.writeConfig(viper, conf)
	if err != nil {
		logUtils.Errorf("还原配置文件设置错误", zap.String("refreshConfig(consts.VIPER)", err.Error()))
		return err
	}
	return nil
}

// InitDB 创建数据库并初始化
func (s *DataService) InitDB(req serverDomain.DataRequest) error {
	defaultConfig := serverConsts.CONFIG
	if serverConsts.VIPER == nil {
		logUtils.Errorf("初始化错误", zap.String("InitDB", ErrViperEmpty.Error()))
		return ErrViperEmpty
	}

	level := req.Level
	if level == "" {
		level = "debug"
	}
	addr := req.Addr
	if addr == "" {
		addr = "127.0.0.1:8085"
	}

	serverConsts.CONFIG.System.CacheType = req.CacheType
	serverConsts.CONFIG.System.Level = level
	serverConsts.CONFIG.System.Addr = addr
	serverConsts.CONFIG.System.DbType = req.SqlType

	if serverConsts.CONFIG.System.CacheType == "redis" {
		serverConsts.CONFIG.Redis = serverConfig.Redis{
			DB:       req.Cache.DB,
			Addr:     fmt.Sprintf("%s:%s", req.Cache.Host, req.Cache.Port),
			Password: req.Cache.Password,
		}
		err := cache.Init() // redis缓存
		if err != nil {
			logUtils.Errorf("认证驱动初始化错误", zap.String("cache.Init() ", err.Error()))
			return err
		}
	}

	if req.Db.Host == "" {
		req.Db.Host = "127.0.0.1"
	}

	if req.Db.Port == "" {
		req.Db.Port = "3306"
	}

	if err := s.DataRepo.CreateTable(req.Db); err != nil {
		return err
	}

	logUtils.Infof("新建数据库", zap.String("库名", req.Db.DBName))

	serverConsts.CONFIG.Mysql.Path = fmt.Sprintf("%s:%s", req.Db.Host, req.Db.Port)
	serverConsts.CONFIG.Mysql.Dbname = req.Db.DBName
	serverConsts.CONFIG.Mysql.Username = req.Db.UserName
	serverConsts.CONFIG.Mysql.Password = req.Db.Password
	serverConsts.CONFIG.Mysql.LogMode = req.Db.LogMode

	m := serverConsts.CONFIG.Mysql
	if m.Dbname == "" {
		logUtils.Errorf("缺少数据库参数")
		return errors.New("缺少数据库参数")
	}

	if err := s.writeConfig(serverConsts.VIPER, serverConsts.CONFIG); err != nil {
		logUtils.Errorf("更新配置文件错误", zap.String("writeConfig(consts.VIPER)", err.Error()))
	}

	if s.DataRepo.DB == nil {
		logUtils.Error("数据库初始化错误")
		s.refreshConfig(serverConsts.VIPER, defaultConfig)
		return errors.New("数据库初始化错误")
	}

	err := s.DataRepo.DB.AutoMigrate(model.Models...)
	if err != nil {
		logUtils.Errorf("迁移数据表错误", zap.String("错误:", err.Error()))
		s.refreshConfig(serverConsts.VIPER, defaultConfig)
		return err
	}

	if req.ClearData {
		err = s.initData(
			s.PermSource,
			s.RoleSource,
			s.UserSource,
		)
		if err != nil {
			logUtils.Errorf("填充数据错误", zap.String("错误:", err.Error()))
			s.refreshConfig(serverConsts.VIPER, defaultConfig)
			return err
		}

		// update password
		if req.Sys.AdminPassword != "" {
			hash, err := bcrypt.GenerateFromPassword([]byte(req.Sys.AdminPassword), bcrypt.DefaultCost)
			if err != nil {
				logUtils.Errorf("密码加密错误", zap.String("错误:", err.Error()))
				return nil
			}
			req.Sys.AdminPassword = string(hash)
			s.UserRepo.UpdatePasswordByName(serverConsts.AdminUserName, req.Sys.AdminPassword)
		}
	}

	return nil
}

// initDB 初始化数据
func (s *DataService) initData(InitDBFunctions ...module.InitDBFunc) error {
	for _, v := range InitDBFunctions {
		err := v.Init()
		if err != nil {
			return err
		}
	}
	return nil
}
