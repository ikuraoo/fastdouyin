package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/ikuraoo/fastdouyin/common"
	"github.com/ikuraoo/fastdouyin/entity"
	"github.com/ikuraoo/fastdouyin/service"
	"net/http"
	"strconv"
	"time"
)

type CommentResponse struct {
	Id         int64
	User       *entity.User
	Content    string
	CreateDate time.Time
}

type CommentActionResponse struct {
	common.Response
	Comment CommentResponse `json:"comment,omitempty"`
}

type CommentListResponse struct {
	common.Response
	CommentList []*service.CommentResponse `json:"comment_list,omitempty"`
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {
	//token := c.Query("token")
	//actionType := c.Query("action_type")

	rawUserId, _ := c.Get("my_uid")
	userId, _ := rawUserId.(int64)

	rawVideoId := c.Query("video_id")
	videoId, _ := strconv.ParseInt(rawVideoId, 10, 64)

	rawActionType := c.Query("action_type")
	actionType, _ := strconv.ParseInt(rawActionType, 10, 32)

	commentText := c.Query("comment_text")

	commentid := c.Query("comment_id")
	commentId, _ := strconv.ParseInt(commentid, 10, 64)

	var err error
	var comment *entity.Comment
	switch actionType {
	case common.PUBLISH_COMMENT:
		comment, err = service.CreateComment(userId, videoId, commentText)
	case common.DELETE_COMMENT:
		comment, err = service.DeleteComment(commentId, videoId)
	}

	user, err := entity.NewUserDaoInstance().QueryById(userId)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{StatusCode: 1, StatusMsg: err.Error()})
	}
	c.JSON(http.StatusOK, CommentActionResponse{Response: common.Response{StatusCode: 0},
		Comment: CommentResponse{
			Id:         comment.Id,
			User:       user,
			Content:    commentText,
			CreateDate: comment.CreateTime,
		}})

}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {

	rawUserId, _ := c.Get("my_uid")
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
