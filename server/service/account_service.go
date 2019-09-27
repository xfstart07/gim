// Author: xufei
// Date: 2019-09-05 11:06

package service

import (
	"fmt"
	"gim/model"
	"gim/server/constant"
	"strconv"
	"strings"
	"sync"

	"github.com/pkg/errors"

	"github.com/go-redis/redis"
)

const (
	AccountPrefix = "gim-account:"
)

type AccountServiceInterface interface {
	Register(user model.User) (model.User, error)

	SaveSession(int64, string)
	GetSessionByUserID(int64) model.User
	RemoveSession(userID int64)
	SaveAndCheckLogin(userID int64) bool

	StoreServerChannelInfo(userID int64, channelInfo model.UserChannelInfo) error
	ServerChannelInfo(userID int64) model.UserChannelInfo
	GetAllServerChannelInfo() []model.UserChannelInfo
	GetAllOnlineUsers() []model.User
}

var (
	accountSrv  *accountService
	accountOnly sync.Once
)

type accountService struct {
	client   *redis.Client
	sessions sync.Map
}

func (s *accountService) Register(user model.User) (model.User, error) {
	key := accountRegisterKey(user.UserID)

	val := s.client.Get(user.UserName).Val()
	if val == "" {
		if err := s.client.Set(key, user.UserName, -1).Err(); err != nil {
			return user, errors.WithStack(err)
		}

		if err := s.client.Set(user.UserName, key, -1).Err(); err != nil {
			return user, errors.WithStack(err)
		}
		return user, nil
	}

	// 返回已注册 UserID
	id := strings.Split(":", val)[1]
	user.UserID, _ = strconv.ParseInt(id, 10, 64)

	return user, constant.ErrAccountRegistered
}

func GetAccountService(client *redis.Client) *accountService {
	accountOnly.Do(func() {
		accountSrv = &accountService{
			client: client,
		}
	})

	return accountSrv
}

func accountRegisterKey(userID int64) string {
	return fmt.Sprintf("%s%d", AccountPrefix, userID)
}
