package service

import (
	"errors"

	"github.com/ikuraoo/fastdouyin/common"
	"github.com/ikuraoo/fastdouyin/entity"
)

type User struct {
	//gorm.Model
	Id            int64  `json:"id,omitempty"  gorm:"primaryKey"`
	Name          string `json:"name,omitempty"`
	FollowCount   int64  `json:"follow_count"`
	FollowerCount int64  `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

func FollowAction(userId, followId, actionType int64) error {
	var err error
	switch actionType {
	case common.FOLLOW:
		err = entity.NewFollowDaoInstance().AddUserFollow(userId, followId)
	case common.CANCEL:
		err = entity.NewFollowDaoInstance().CancelUserFollow(userId, followId)
	default:
		return errors.New("未定义操作")

	}
	return err
}

func QueryFollowList(userId int64) ([]*User, error) {
	var userList []*entity.User
	err := entity.NewUserDaoInstance().GetFollowListByUserId(userId, &userList)
	if err != nil {
		return nil, err
	}
	var userlist []*User
	userlist, err = UserTypeConvert(userList)
	return userlist, nil
}

func QueryFollowerList(userId int64) ([]*User, error) {
	var userList []*entity.User
	err := entity.NewUserDaoInstance().GetFollowerListByUserId(userId, &userList)
	if err != nil {
		return nil, err
	}
	var userlist []*User
	userlist, err = UserTypeConvert(userList)
	return userlist, nil
}

func UserTypeConvert(entityusers []*entity.User) ([]*User, error) {
	var userlist []*User
	for _, entityuser := range entityusers {
		user := User{
			Id:            entityuser.Id,
			Name:          entityuser.Name,
			FollowCount:   entityuser.FollowCount,
			FollowerCount: entityuser.FollowerCount,
			IsFollow:      false,
		}
		userlist = append(userlist, &user)
	}
	return userlist, nil
}
