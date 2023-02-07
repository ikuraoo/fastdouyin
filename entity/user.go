package entity

import (
	"github.com/ikuraoo/fastdouyin/dao"
	"sync"
	"time"
)

type User struct {
	Id            int64 `gorm:"column:id"`
	Name          string
	Password      string
	Token         string
	FollowCount   int64
	FollowerCount int64
	CreateTime    time.Time
	UpdateTime    time.Time
	IsDeleted     bool
}

type UserDao struct {
}

var userDao *UserDao
var userOnce sync.Once

func NewUserDaoInstance() *UserDao {
	userOnce.Do(
		func() {
			userDao = &UserDao{}
		})
	return userDao
}

func (*UserDao) CreateUser(user *User) error {
	result := dao.Db.Create(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (*UserDao) FindByToken(name string) (User, bool) {
	var user User
	result := dao.Db.Where("token = ?", name).First(&user)
	if result.RowsAffected != 0 {
		return user, true
	}
	return user, false
}
