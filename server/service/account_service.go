// Author: xufei
// Date: 2019-09-05 11:06

package service

import (
	"errors"
	"fmt"
	"gim/internal/lg"
	"gim/model"
	"strconv"
	"strings"

	"go.uber.org/zap"

	"github.com/go-redis/redis"
)

const (
	AccountPrefix = "gim-account:"
)

var (
	errAccountStore      = errors.New("用户存储失败")
	errAccountRegistered = errors.New("用户已注册")
)

type AccountServiceInterface interface {
	Register(user model.User) (model.User, error)
}

type accountService struct {
	client *redis.Client
}

func (s *accountService) Register(user model.User) (model.User, error) {
	key := fmt.Sprintf("%s%d", AccountPrefix, user.UserID)

	val := s.client.Get(user.UserName).Val()
	if val == "" {
		err := s.client.Set(key, user.UserName, -1).Err()
		if err != nil {
			lg.Logger().Error("redis 存储失败", zap.Error(err))
			return user, errAccountStore
		}

		err = s.client.Set(user.UserName, key, -1).Err()
		if err != nil {
			lg.Logger().Error("redis 存储失败", zap.Error(err))
			return user, errAccountStore
		}
		return user, nil
	}

	// 返回已注册 UserID
	id := strings.Split(":", val)[1]
	user.UserID, _ = strconv.ParseInt(id, 10, 64)

	return user, errAccountRegistered
}

func NewAccountService(client *redis.Client) *accountService {
	return &accountService{
		client: client,
	}
}
