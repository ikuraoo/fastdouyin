package entity

import (
	"errors"
	"github.com/ikuraoo/fastdouyin/common"
	"sync"
	"time"
)

type User struct {
	Id             int64 `gorm:"column:id"`
	Name           string
	Password       string `json:"Password,omitempty"`
	FollowCount    int64
	FollowerCount  int64
	TotalFavorited int64
	WorkCount      int64
	FavoriteCount  int64
	CreateTime     time.Time `json:"CreateTime,omitempty"`
	UpdateTime     time.Time `json:"UpdateTime,omitempty"`
	IsDeleted      bool
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
	if user.Id == common.WRONG_ID {
		return nil, errors.New(common.USER_NOT_EXIT)
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

func (*UserDao) GetFollowListByUserId(userId int64, userList *[]*User) error {
	if userList == nil {
		return errors.New("空指针错误")
	}
	var err error
	if err = db.Raw("SELECT u.* FROM users u, follows f WHERE f.my_uid = ? AND f.his_uid = u.id", userId).Scan(&userList).Error; err != nil {
		return err
	}
	return nil
}

func (*UserDao) GetFollowerListByUserId(userId int64, userList *[]*User) error {
	if userList == nil {
		return errors.New("空指针错误")
	}
	var err error
	if err = db.Raw("SELECT u.* FROM users u, follows f WHERE f.his_uid = ? AND f.my_uid = u.id", userId).Scan(&userList).Error; err != nil {
		return err
	}

	return nil
}
