package controller

import (
	"github.com/ikuraoo/fastdouyin/entity"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin

var usersLoginInfo = map[string]User{
	"zhangleidouyin": {
		Id:            1,
		Name:          "zhanglei",
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      true,
	},
}

//var userIdSequence = int64(0)

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User entity.User `json:"user"`
}

func Register(c *gin.Context) {

	username := c.Query("username")
	password := c.Query("password")
	token := username + password
	//if _, exist := usersLoginInfo[token]; exist {
	userdao := entity.NewUserDaoInstance()
	if _, exist := userdao.FindByToken(token); exist {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
	} else {
		//atomic.AddInt64(&userIdSequence, 1)
		newUser := entity.User{
			//Id:            userIdSequence,
			Name:          username,
			Password:      password,
			Token:         token,
			FollowCount:   0,
			FollowerCount: 0,
			CreateTime:    time.Now(),
			UpdateTime:    time.Now(),
			IsDeleted:     false,
		}
		//usersLoginInfo[token] = newUser
		userdao.CreateUser(&newUser)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   newUser.Id,
			Token:    token,
		})
	}

}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	token := username + password
	userdao := entity.NewUserDaoInstance()
	if user, exist := userdao.FindByToken(token); exist {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   user.Id,
			Token:    token,
		})
	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}
}

func UserInfo(c *gin.Context) {
	token := c.Query("token")
	userdao := entity.NewUserDaoInstance()
	if user, exist := userdao.FindByToken(token); exist {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0},
			User:     user,
		})
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}
}
