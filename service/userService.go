package service

import (
	"errors"
	"github.com/ikuraoo/fastdouyin/common"
	"github.com/ikuraoo/fastdouyin/configure"
	"github.com/ikuraoo/fastdouyin/entity"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserLoginMessage struct {
	name     string
	password string
}

func UserLogin(name string, password string) (int64, error) {
	var user *UserLoginMessage = &UserLoginMessage{
		name:     name,
		password: password,
	}
	return user.Do()
}

func (u *UserLoginMessage) Do() (int64, error) {
	if err := u.check(); err != nil {
		return 0, errors.New("请检查输入")
	}
	id, err := u.login()
	if err != nil {
		return common.WRONG_ID, err
	}
	return id, nil
}

func (u *UserLoginMessage) check() error {
	//登录不需要检查
	return nil
}

func (u *UserLoginMessage) login() (int64, error) {
	user, err := entity.NewUserDaoInstance().QueryByName(u.name)

	if err != nil || user.Id == 0 {
		return common.WRONG_ID, errors.New("用户不存在")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(u.password))
	if err != nil {
		return common.WRONG_ID, errors.New("密码错误")
	}
	return user.Id, nil
}

type UserRegisterMessage struct {
	name     string
	password string
}

func UserRegister(username string, password string) (int64, error) {
	var userRegisterMessage = &UserRegisterMessage{
		name:     username,
		password: password,
	}
	return userRegisterMessage.Do()
}

func (u *UserRegisterMessage) Do() (int64, error) {
	//检查参数
	if err := u.check(); err != nil {
		return common.WRONG_ID, err
	}

	id, err := u.register()
	if err != nil {
		return common.MISTAKE, err
	}
	return id, nil
}

func (u *UserRegisterMessage) check() error {
	exist := entity.NewUserDaoInstance().IsUserExistByName(u.name)
	if exist {
		return errors.New("用户已存在")
	}
	return nil
}

func (u *UserRegisterMessage) register() (int64, error) {
	password, err := bcrypt.GenerateFromPassword([]byte(u.password), bcrypt.DefaultCost)
	if err != nil {
		return common.MISTAKE, errors.New("密码加密出错")
	}

	user := &entity.User{
		//Id:            0,
		Name:           u.name,
		Password:       string(password),
		FollowCount:    0,
		FollowerCount:  0,
		TotalFavorited: 0,
		WorkCount:      0,
		FavoriteCount:  0,
		CreateTime:     time.Now(),
		UpdateTime:     time.Now(),
		IsDeleted:      false,
	}
	err = entity.NewUserDaoInstance().CreateUser(user)
	if err != nil {
		return common.MISTAKE, err
	}
	return user.Id, nil
}

func UserInfo(userId, targetId int64) (*common.UserResponse, error) {
	//检查参数
	exist := entity.NewUserDaoInstance().IsUserExistById(userId)
	if !exist {
		return nil, errors.New(common.USER_NOT_EXIT)
	}
	user, err := entity.NewUserDaoInstance().QueryById(targetId)
	if err != nil || user.Id == 0 {
		return nil, errors.New(common.USER_NOT_EXIT)
	}

	//转为输出格式
	userResponse, err := ConvertToUserResponse(user, userId)
	if err != nil {
		return nil, errors.New("格式转换失败")
	}
	return userResponse, err
}

func ConvertToUserResponse(author *entity.User, userId int64) (*common.UserResponse, error) {
	var isfollow bool
	if userId == 0 {
		isfollow = false
	} else if userId == author.Id {
		isfollow = true
	} else {
		isfollow = configure.NewProxyIndexMap().GetAFollowB(userId, author.Id)
		//isfollow, _ = entity.NewFollowDaoInstance().QueryIsFollow(userId, author.Id)
	}

	userResponse := common.UserResponse{
		Id:              author.Id,
		Name:            author.Name,
		FollowCount:     author.FollowCount,
		FollowerCount:   author.FollowerCount,
		IsFollow:        isfollow,
		Avatar:          author.Avatar,
		BackgroundImage: author.BackgroundImage,
		Signature:       author.Signature,
		TotalFavorited:  author.TotalFavorited,
		WorkCount:       author.WorkCount,
		FavoriteCount:   author.FavoriteCount,
	}
	return &userResponse, nil
}
