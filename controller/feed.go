package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/ikuraoo/fastdouyin/common"
	"github.com/ikuraoo/fastdouyin/middleware"
	"github.com/ikuraoo/fastdouyin/service"
	"net/http"
	"strconv"
	"time"
)

type FeedResponse struct {
	common.Response
	VideoList []*common.VideoMessage `json:"video_list,omitempty"`
	NextTime  int64                  `json:"next_time,omitempty"`
}

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	var latestTime time.Time

	rawTimestamp, ok := c.GetQuery("latest_time")
	if ok {
		intTime, _ := strconv.ParseInt(rawTimestamp, 10, 64)
		latestTime = time.Unix(intTime, 0)
	} else {
		latestTime = time.Now()
	}

	token, ok := c.GetQuery("token")
	var videoList []*common.VideoMessage
	var err error
	if !ok {
		videoList, err = DoNoToken(latestTime)
		if err != nil {
			common.SendError(c, err.Error())
		}
	} else {
		videoList, err = DoHasToken(token, latestTime)
		if err != nil {
			common.SendError(c, err.Error())
		}
	}

	c.JSON(http.StatusOK, FeedResponse{
		Response:  common.Response{StatusCode: 0},
		VideoList: videoList,
		NextTime:  time.Now().Unix(),
	})

}

func DoNoToken(latestTime time.Time) ([]*common.VideoMessage, error) {
	videoList, err := service.VideoFeed(0, latestTime)
	if err != nil {
		return nil, err
	}
	return videoList, nil
}

func DoHasToken(token string, latestTime time.Time) ([]*common.VideoMessage, error) {
	//解析成功
	if _, claim, ok := middleware.ParseToken(token); ok == nil {
		//token超时
		if time.Now().Unix() > claim.ExpiresAt {
			return nil, errors.New("token超时")
		}

		//调用service层接口
		videoList, err := service.VideoFeed(claim.UserId, latestTime)
		if err != nil {
			return nil, err
		}
		return videoList, nil
	}
	//解析失败
	return nil, errors.New("token不正确")
}
