package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ikuraoo/fastdouyin/common"
	"github.com/ikuraoo/fastdouyin/middleware"
	"github.com/ikuraoo/fastdouyin/service"
	"net/http"
)

type UserLoginResponse struct {
	common.Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	common.Response
	common.UserMessage
}

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	id, err := service.UserRegister(username, password)
	fmt.Println(id)
	if err != nil {
		common.SendError(c, err.Error())
	}

	token, err := middleware.CreateToken(id)
	if err != nil {
		common.SendError(c, err.Error())
	}

	c.JSON(http.StatusOK, UserLoginResponse{
		Response: common.Response{StatusCode: common.RESP_SUCCESS},
		UserId:   id,
		Token:    token,
	})
	return
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	id, err := service.UserLogin(username, password)
	//登录失败
	if err != nil {
		common.SendError(c, err.Error())
	}

	token, err := middleware.CreateToken(id)
	if err != nil {
		common.SendError(c, err.Error())
	}

	c.JSON(http.StatusOK, UserLoginResponse{
		Response: common.Response{StatusCode: 0},
		UserId:   id,
		Token:    token,
	})
	return
}

func UserInfo(c *gin.Context) {

	rawUserId, ok := c.Get("userId")
	if !ok {
		common.SendError(c, "解析token失败")
	}

	userId := rawUserId.(int64)
	user, err := service.UserInfo(userId)
	if err != nil {
		common.SendError(c, err.Error())
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response:    common.Response{StatusCode: 0},
			UserMessage: *user,
		})
	}
}
