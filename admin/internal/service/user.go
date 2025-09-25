package service

import (
	"github.com/puoxiu/discron/common/models"
	"github.com/puoxiu/discron/common/pkg/dbclient"
)

type UserService struct {
}

var DefaultUserService = new(UserService)

func (us *UserService) Login(username, password string) (u *models.User, err error) {
	err = dbclient.GetMysqlDB().Table(models.CronixUserTableName).Where("username = ? And password = ?", username, password).Find(u).Error
	return
}
