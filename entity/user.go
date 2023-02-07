package entity

import (
	"errors"
	"github.com/ikuraoo/fastdouyin/constant"
	"sync"
	"time"
)

type User struct {
	Id            int64 `gorm:"column:id"`
	Name          string
	Password      string
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

func (*UserDao) QueryById(id int64) (*User, error) {
	var user User
	err := db.Where("id = ?", id).Find(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (*UserDao) QueryByName(name string) (*User, error) {
	var user User
	err := db.Where("name = ?", name).Find(&user).Error
	if err != nil {
		return nil, errors.New("查询出错")
	}
	if user.Id == constant.WRONG_ID {
		return nil, errors.New(constant.USER_NOT_EXIT)
	}
	return &user, nil
}

func (*UserDao) CreateUser(user *User) error {
	err := db.Create(user).Error
	if err != nil {
		return err
	}
	return nil
}
