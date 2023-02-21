package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ikuraoo/fastdouyin/common"
	"github.com/ikuraoo/fastdouyin/service"
	"github.com/ikuraoo/fastdouyin/util"
	"net/http"
	"path/filepath"
	"strconv"
)

type VideoListResponse struct {
	common.Response
	VideoList []*common.VideoResponse `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	//解析参数
	rawUserId, ok := c.Get("userId")
	if !ok {
		common.SendError(c, "解析token失败")
	}
	userId := rawUserId.(int64)
	title := c.PostForm("title")
	data, err := c.FormFile("data")
	if err != nil {
		common.SendError(c, err.Error())
	}

	//视频路径
	filename := filepath.Base(data.Filename)
	finalName := fmt.Sprintf("%d_%s", userId, filename)
	saveFile := filepath.Join("./public/videos/", finalName)
	if err = c.SaveUploadedFile(data, saveFile); err != nil {
		common.SendError(c, err.Error())
	}

	//封面路径
	coverName := util.NewFileCoverName(userId)
	saveCoverFile := filepath.Join("./public/covers/", coverName)
	err = util.SnapShotFromVideo(saveFile, saveCoverFile, 1)
	if err != nil {
		common.SendError(c, err.Error())
	}

	if err = service.VideoPublish(userId, title, finalName, coverName); err != nil {
		common.SendError(c, err.Error())
	}
	c.JSON(http.StatusOK, common.Response{
		StatusCode: 0,
		StatusMsg:  finalName + " uploaded successfully",
	})

}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	//解析参数
	rawUserId, ok := c.Get("userId")
	if !ok {
		common.SendError(c, "解析token失败")
	}
	userId := rawUserId.(int64)
	rawAuthorId := c.Query("user_id")
	authorId, err := strconv.ParseInt(rawAuthorId, 10, 64)

	//service
	videoFeed, err := service.PublishList(userId, authorId)
	if err != nil {
		common.SendError(c, err.Error())
	}
	c.JSON(http.StatusOK, VideoListResponse{
		Response: common.Response{
			StatusCode: 0,
		},
		VideoList: videoFeed,
	})
}
