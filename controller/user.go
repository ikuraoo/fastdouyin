package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ikuraoo/fastdouyin/common"
	"github.com/ikuraoo/fastdouyin/middleware"
	"github.com/ikuraoo/fastdouyin/service"
	"net/http"
	"strconv"
)

type UserLoginResponse struct {
	common.Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	common.Response
	*common.UserResponse `json:"user"`
}

func Register(c *gin.Context) {
	//解析参数
	username := c.Query("username")
	password := c.Query("password")

	//调用service层
	id, err := service.UserRegister(username, password)
	fmt.Println(id)
	if err != nil {
		common.SendError(c, err.Error())
		return
	}

	//颁发token
	token, err := middleware.CreateToken(id)
	if err != nil {
		common.SendError(c, err.Error())
	}

	c.JSON(http.StatusOK, UserLoginResponse{
		Response: common.Response{StatusCode: common.RESP_SUCCESS},
		UserId:   id,
		Token:    token,
	})
}

func Login(c *gin.Context) {
	//解析参数
	username := c.Query("username")
	password := c.Query("password")

	//service层
	id, err := service.UserLogin(username, password)
	//登录失败
	if err != nil {
		common.SendError(c, err.Error())
		return
	}

	token, err := middleware.CreateToken(id)
	if err != nil {
		common.SendError(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, UserLoginResponse{
		Response: common.Response{StatusCode: 0},
		UserId:   id,
		Token:    token,
	})
}

func UserInfo(c *gin.Context) {
	//解析参数
	rawUserId, ok := c.Get("userId")
	if !ok {
		common.SendError(c, "解析token失败")
		return
	}
	userId := rawUserId.(int64)

	rawtargetId := c.Query("user_id")
	targetId, err := strconv.ParseInt(rawtargetId, 10, 64)
	if err != nil {
		common.SendError(c, err.Error())
		return
	}
	fmt.Println(targetId)
	//service
	user, err := service.UserInfo(userId, targetId)
	if err != nil {
		common.SendError(c, err.Error())
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response:     common.Response{StatusCode: 0},
			UserResponse: user,
		})
	}
}
