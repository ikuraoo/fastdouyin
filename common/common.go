package common

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func SendError(c *gin.Context, message string) {
	c.JSON(http.StatusOK, Response{StatusCode: RESP_MISTAKE, StatusMsg: message})
	return
}

var (
	WRONG_ID        int64  = 0
	MISTAKE         int64  = -1
	RESP_MISTAKE    int32  = 1
	RESP_SUCCESS    int32  = 0
	USER_NOT_EXIT   string = "用户不存在"
	PUBLISH_COMMENT int64  = 1
	DELETE_COMMENT  int64  = 2
	FOLLOW          int64  = 1
	CANCEL          int64  = 2
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type VideoMessage struct {
	Id            int64       `json:"id,omitempty"`
	Author        UserMessage `json:"author"`
	PlayUrl       string      `json:"play_url" json:"play_url,omitempty"`
	CoverUrl      string      `json:"cover_url"`
	FavoriteCount int64       `json:"favorite_count"`
	CommentCount  int64       `json:"comment_count"`
	IsFavorite    bool        `json:"is_favorite"`
	Title         string      `json:"title"`
}
type Comment struct {
	Id         int64       `json:"id,omitempty"`
	User       UserMessage `json:"user"`
	Content    string      `json:"content,omitempty"`
	CreateDate string      `json:"create_date,omitempty"`
}

type UserMessage struct {
	//gorm.Model
	Id             int64  `json:"id,omitempty"  gorm:"primaryKey"`
	Name           string `json:"name,omitempty"`
	FollowCount    int64  `json:"follow_count"`
	FollowerCount  int64  `json:"follower_count"`
	IsFollow       bool   `json:"is_follow"`
	TotalFavorited int64  `json:"total_favorited"`
	WorkCount      int64  `json:"work_count"`
	FavoriteCount  int64  `json:"favorite_count"`
}

type Message struct {
	Id         int64  `json:"id,omitempty"`
	Content    string `json:"content,omitempty"`
	CreateTime string `json:"create_time,omitempty"`
}

type MessageSendEvent struct {
	UserId     int64  `json:"user_id,omitempty"`
	ToUserId   int64  `json:"to_user_id,omitempty"`
	MsgContent string `json:"msg_content,omitempty"`
}

type MessagePushEvent struct {
	FromUserId int64  `json:"user_id,omitempty"`
	MsgContent string `json:"msg_content,omitempty"`
}
