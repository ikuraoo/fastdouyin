package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/ikuraoo/fastdouyin/common"
	"github.com/ikuraoo/fastdouyin/service"
	"net/http"
	"strconv"
)

type CommentActionResponse struct {
	common.Response
	Comment common.Comment `json:"comment,omitempty"`
}

type CommentListResponse struct {
	common.Response
	CommentList []*common.Comment `json:"comment_list,omitempty"`
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {

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
	actionType, err := strconv.ParseInt(rawActionType, 10, 32)
	if err != nil {
		common.SendError(c, err.Error())
	}

	//调用逻辑
	var comment *common.Comment
	switch actionType {
	case common.PUBLISH_COMMENT:
		commentText := c.Query("comment_text")
		comment, err = service.CreateComment(userId, videoId, commentText)

	case common.DELETE_COMMENT:
		commentid := c.Query("comment_id")
		commentId, err := strconv.ParseInt(commentid, 10, 64)
		if err != nil {
			common.SendError(c, err.Error())
		}
		comment, err = service.DeleteComment(commentId, videoId)
		if err != nil {
			common.SendError(c, err.Error())
		}
	}
	c.JSON(http.StatusOK, CommentActionResponse{Response: common.Response{StatusCode: 0},
		Comment: *comment,
	})

}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {

	rawUserId, _ := c.Get("userId")
	userId, _ := rawUserId.(int64)

	rawVideoId := c.Query("video_id")
	videoId, _ := strconv.ParseInt(rawVideoId, 10, 64)

	commentList, err := service.QueryCommentList(userId, videoId)
	commentListwithuser, err := service.CombinationCommentsAndUsers(commentList)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{StatusCode: 1, StatusMsg: err.Error()})
	}
	c.JSON(http.StatusOK, CommentListResponse{
		Response:    common.Response{StatusCode: 0},
		CommentList: commentListwithuser,
	})
}
