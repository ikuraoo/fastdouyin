package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/ikuraoo/fastdouyin/common"
	"github.com/ikuraoo/fastdouyin/service"
	"net/http"
	"strconv"
)

type UserListResponse struct {
	common.Response
	UserList []*common.UserResponse `json:"user_list"`
}

// RelationAction no practical effect, just check if token is valid
func RelationAction(c *gin.Context) {
	//解析参数
	rawUserId, _ := c.Get("userId")
	userId, ok := rawUserId.(int64)
	if !ok {
		common.SendError(c, "token解析失败")
		return
	}

	rawFollowId := c.Query("to_user_id")
	followId, err := strconv.ParseInt(rawFollowId, 10, 64)
	if err != nil {
		common.SendError(c, err.Error())
		return
	}

	rawActionType := c.Query("action_type")
	actionType, err := strconv.ParseInt(rawActionType, 10, 32)
	if err != nil {
		common.SendError(c, err.Error())
		return
	}

	//service
	err = service.UserFollowAction(userId, followId, actionType)
	if err != nil {
		common.SendError(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, common.Response{StatusCode: 0})

}

// FollowList all users have same follow list
func FollowList(c *gin.Context) {
	//解析参数
	rawUserId, _ := c.Get("userId")
	userId, ok := rawUserId.(int64)
	if !ok {
		common.SendError(c, "token解析失败")
		return
	}

	rawtargetId := c.Query("user_id")
	targetId, err := strconv.ParseInt(rawtargetId, 10, 64)
	if err != nil {
		common.SendError(c, err.Error())
		return
	}

	list, err := service.QueryFollowList(userId, targetId)
	if err != nil {
		common.SendError(c, err.Error())
	} else {
		c.JSON(http.StatusOK, UserListResponse{
			Response: common.Response{
				StatusCode: 0,
			},
			UserList: list,
		})
	}

}

// FollowerList all users have same follower list
func FollowerList(c *gin.Context) {
	//解析参数
	rawUserId, _ := c.Get("userId")
	userId, ok := rawUserId.(int64)
	if !ok {
		common.SendError(c, "token解析失败")
		return
	}
	rawtargetId := c.Query("user_id")
	targetId, err := strconv.ParseInt(rawtargetId, 10, 64)
	if err != nil {
		common.SendError(c, err.Error())
		return
	}

	//service层
	list, err := service.QueryFollowerList(userId, targetId)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{StatusCode: 1, StatusMsg: err.Error()})
	} else {
		c.JSON(http.StatusOK, UserListResponse{
			Response: common.Response{
				StatusCode: 0,
			},
			UserList: list,
		})
	}

}
