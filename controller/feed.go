package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ikuraoo/fastdouyin/constant"
	"github.com/ikuraoo/fastdouyin/service"
	"net/http"
	"strconv"
	"time"
)

type FeedResponse struct {
	Response
	VideoList []*service.VideoWithUser `json:"video_list,omitempty"`
	NextTime  int64                    `json:"next_time,omitempty"`
}

// Feed same demo video list for every request
func Feed(c *gin.Context) {

	uid := c.Query("uid")
	fmt.Println(uid)
	if uid == "" {
		uid = "0"
	}
	myUId, _ := strconv.ParseInt(uid, 10, 64)
	videoFeed, err := service.VideoFeed(myUId, 1)
	if err != nil {
		c.JSON(http.StatusOK, FeedResponse{
			Response: Response{
				StatusCode: constant.RESP_MISTAKE,
				StatusMsg:  err.Error(),
			},
		})
		return
	}
	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: videoFeed,
		NextTime:  time.Now().Unix(),
	})
}
