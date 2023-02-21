package service

import (
	"errors"
	"github.com/ikuraoo/fastdouyin/configure"

	"github.com/ikuraoo/fastdouyin/common"
	"github.com/ikuraoo/fastdouyin/entity"
)

type FollowActionMessage struct {
	userId     int64
	followId   int64
	actionType int64
}

func NewFollowAction(userId, followId, action int64) *FollowActionMessage {
	return &FollowActionMessage{
		userId:     userId,
		followId:   followId,
		actionType: action,
	}

}

func UserFollowAction(userId, followId, actionType int64) error {
	return NewFollowAction(userId, followId, actionType).Do()
}

func (f *FollowActionMessage) Do() error {
	var err error
	if err = f.checkNum(); err != nil {
		return err
	}

	switch f.actionType {
	case PLUS:
		err = f.PlusOperation()
	case MINUS:
		err = f.CancelOperation()
	default:
		return errors.New("未定义的行为")
	}
	return err
}

func (f *FollowActionMessage) checkNum() error {
	exist := entity.NewUserDaoInstance().IsUserExistById(f.userId)
	if !exist {
		return errors.New("用户不存在")
	}
	exist = entity.NewUserDaoInstance().IsUserExistById(f.followId)
	if !exist {
		return errors.New("用户不存在")
	}

	if f.actionType != common.FOLLOW && f.actionType != common.CANCEL {
		return errors.New("未定义的行为")
	}
	return nil
}

func (f *FollowActionMessage) PlusOperation() error {
	//查询是否关注
	if f.userId == f.followId {
		return errors.New("自己不能关注自己")
	}
	isfollow := configure.NewProxyIndexMap().GetAFollowB(f.userId, f.followId)
	//isfollow, _ := entity.NewFollowDaoInstance().QueryIsFollow(f.userId, f.followId)
	if !isfollow {
		err := entity.NewFollowDaoInstance().AddUserFollow(f.userId, f.followId)
		//同时插入redis缓存
		configure.NewProxyIndexMap().SetAFollowB(f.userId, f.followId, true)
		if err != nil {
			return err
		}
	} else {
		return errors.New("已关注该用户")
	}
	return nil
}

func (f *FollowActionMessage) CancelOperation() error {
	if f.userId == f.followId {
		return errors.New("")
	}
	//查询是否关注,查缓存
	isfollow := configure.NewProxyIndexMap().GetAFollowB(f.userId, f.followId)
	//isfollow, _ := entity.NewFollowDaoInstance().QueryIsFollow(f.userId, f.followId)
	if isfollow {
		err := entity.NewFollowDaoInstance().CancelUserFollow(f.userId, f.followId)
		//同时更新redis缓存
		configure.NewProxyIndexMap().SetAFollowB(f.userId, f.followId, false)
		if err != nil {
			return err
		}
	} else {
		return errors.New("没有关注该用户")
	}
	return nil
}

func QueryFollowList(userId, targetId int64) ([]*common.UserResponse, error) {
	//检查参数
	_, err := entity.NewVideoDaoInstance().QueryAuther(userId)
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	_, err = entity.NewVideoDaoInstance().QueryAuther(targetId)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	var userList []*entity.User
	err = entity.NewUserDaoInstance().GetFollowListByUserId(targetId, &userList)
	if err != nil {
		return nil, err
	}

	userlist, err := ConvertUsersType(userId, userList)
	return userlist, nil
}

func QueryFollowerList(userId, targetId int64) ([]*common.UserResponse, error) {
	//检查参数
	_, err := entity.NewVideoDaoInstance().QueryAuther(userId)
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	_, err = entity.NewVideoDaoInstance().QueryAuther(targetId)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	var userList []*entity.User
	err = entity.NewUserDaoInstance().GetFollowerListByUserId(targetId, &userList)
	if err != nil {
		return nil, err
	}

	userlist, err := ConvertUsersType(userId, userList)
	return userlist, nil
}

func ConvertUsersType(userId int64, userList []*entity.User) (userlist []*common.UserResponse, err error) {
	for _, user := range userList {
		userResponse, err := ConvertToUserResponse(user, userId)
		if err != nil {
			return nil, errors.New("类型转换失败")
		}
		userlist = append(userlist, userResponse)
	}
	return userlist, nil
}
