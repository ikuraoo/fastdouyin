package service

import (
	"errors"
	"fmt"
	"github.com/ikuraoo/fastdouyin/constant"
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
		return constant.WRONG_ID, err
	}
	return id, nil
}

func (u *UserLoginMessage) check() error {
	//实现参数的检查、防止sql注入
	return nil
}

func (u *UserLoginMessage) login() (int64, error) {
	user, err := entity.NewUserDaoInstance().QueryByName(u.name)

	if err != nil || user.Id == 0 {
		return constant.WRONG_ID, errors.New("用户不存在")
	}

	fmt.Println("库中密码 :" + user.Password)
	fmt.Println("输入密码：" + u.password)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(u.password))
	if err != nil {
		return constant.WRONG_ID, errors.New("密码错误")
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
	if err := u.check(); err != nil {
		return constant.WRONG_ID, errors.New("请检查输入")
	}
	id, err := u.register()
	if err != nil {
		return constant.MISTAKE, err
	}
	return id, nil
}

func (u *UserRegisterMessage) check() error {
	//实现参数的检查、防止sql注入
	return nil
}

func (u *UserRegisterMessage) register() (int64, error) {
	password, err := bcrypt.GenerateFromPassword([]byte(u.password), bcrypt.DefaultCost)
	if err != nil {
		return constant.MISTAKE, errors.New("密码加密出错")
	}
	_, err = entity.NewUserDaoInstance().QueryByName(u.name)
	if err == nil {
		return constant.WRONG_ID, errors.New("用户已存在")
	}

	user := &entity.User{
		Id:            0,
		Name:          u.name,
		Password:      string(password),
		FollowCount:   0,
		FollowerCount: 0,
		CreateTime:    time.Now(),
		UpdateTime:    time.Now(),
		IsDeleted:     false,
	}
	err = entity.NewUserDaoInstance().CreateUser(user)
	if err != nil {
		return constant.MISTAKE, err
	}
	return user.Id, nil
}
