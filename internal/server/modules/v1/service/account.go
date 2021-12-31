package service

import (
	"errors"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/repo"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"time"

	"github.com/snowlyg/multi"
	"go.uber.org/zap"
)

var (
	ErrUserNameOrPassword = errors.New("用户名或密码错误")
)

type AccountService struct {
	UserRepo *repo.UserRepo `inject:""`
}

func NewAuthService() *AccountService {
	return &AccountService{}
}

// GetAccessToken 登录
func (s *AccountService) GetAccessToken(req serverDomain.LoginRequest) (string, error) {
	admin, err := s.UserRepo.FindPasswordByUserName(req.Username)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(req.Password)); err != nil {
		logUtils.Errorf("用户名或密码错误", zap.String("密码:", req.Password), zap.String("hash:", admin.Password), zap.String("bcrypt.CompareHashAndPassword()", err.Error()))
		return "", ErrUserNameOrPassword
	}

	claims := &multi.CustomClaims{
		ID:            strconv.FormatUint(uint64(admin.Id), 10),
		Username:      req.Username,
		AuthorityId:   "",
		AuthorityType: multi.AdminAuthority,
		LoginType:     multi.LoginTypeWeb,
		AuthType:      multi.AuthPwd,
		CreationDate:  time.Now().Local().Unix(),
		ExpiresIn:     multi.RedisSessionTimeoutWeb.Milliseconds(),
	}
	token, _, err := multi.AuthDriver.GenerateToken(claims)
	if err != nil {
		return "", err
	}

	return token, nil
}
