package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/ikuraoo/fastdouyin/common"
	"github.com/ikuraoo/fastdouyin/service"
	"net/http"
	"strconv"
)

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(c *gin.Context) {
	//解析参数
	rawUserId, _ := c.Get("userId")
	userId, ok := rawUserId.(int64)
	if !ok {
		common.SendError(c, "token解析失败")
	}

	rawVideoId := c.Query("video_id")
	videoId, err := strconv.ParseInt(rawVideoId, 10, 64)
	if err != nil {
		common.SendError(c, err.Error())
	}

	rawActionType := c.Query("action_type")
	actionType, err := strconv.ParseInt(rawActionType, 10, 64)
	if err != nil {
		common.SendError(c, err.Error())
	}

	//调用service层
	err = service.UserFavoriteAction(userId, videoId, actionType)
	if err != nil {
		common.SendError(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, common.Response{StatusCode: 0})

}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	//解析参数
	rawUserId, _ := c.Get("userId")
	userId, ok := rawUserId.(int64)
	if !ok {
		common.SendError(c, "token解析失败")
	}
	id := c.Query("user_id")
	targetId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		common.SendError(c, err.Error())
	}

	//调用service
	uidFavorVideoList, err := service.QueryUserFavorList(userId, targetId)
	if err != nil {
		common.SendError(c, err.Error())
	}
	c.JSON(http.StatusOK, VideoListResponse{
		Response: common.Response{
			StatusCode: 0,
		},
		VideoList: uidFavorVideoList,
	})
}
