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
	UserList []*service.User `json:"user_list"`
}

// RelationAction no practical effect, just check if token is valid
func RelationAction(c *gin.Context) {
	rawUserId, _ := c.Get("my_uid")
	userId, _ := rawUserId.(int64)

	rawFollowId := c.Query("to_user_id")
	followId, _ := strconv.ParseInt(rawFollowId, 10, 64)

	rawActionType := c.Query("action_type")
	actionType, _ := strconv.ParseInt(rawActionType, 10, 32)

	err := service.FollowAction(userId, followId, actionType)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
	c.JSON(http.StatusOK, common.Response{StatusCode: 0})

}

// FollowList all users have same follow list
func FollowList(c *gin.Context) {
	rawUserId, _ := c.Get("my_uid")
	userId, _ := rawUserId.(int64)

	list, err := service.QueryFollowList(userId)
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

// FollowerList all users have same follower list
func FollowerList(c *gin.Context) {
	rawUserId, _ := c.Get("my_uid")
	userId, _ := rawUserId.(int64)

	list, err := service.QueryFollowerList(userId)
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

// FriendList all users have same friend list
func FriendList(c *gin.Context) {
	c.JSON(http.StatusOK, UserListResponse{
		Response: common.Response{
			StatusCode: 0,
		},
		//UserList: []User{DemoUser},
	})
}
