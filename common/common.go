package common

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func SendError(c *gin.Context, message string) {
	c.JSON(http.StatusOK, Response{
		StatusCode: RESP_MISTAKE,
		StatusMsg:  message})
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

type VideoResponse struct {
	Id            int64        `json:"id,omitempty"`
	Author        UserResponse `json:"author"`
	PlayUrl       string       `json:"play_url" json:"play_url,omitempty"`
	CoverUrl      string       `json:"cover_url"`
	FavoriteCount int64        `json:"favorite_count"`
	CommentCount  int64        `json:"comment_count"`
	IsFavorite    bool         `json:"is_favorite"`
	Title         string       `json:"title"`
}
type Comment struct {
	Id         int64        `json:"id,omitempty"`
	User       UserResponse `json:"user"`
	Content    string       `json:"content,omitempty"`
	CreateDate string       `json:"create_date,omitempty"`
}

type UserResponse struct {
	//gorm.Model
	Id              int64  `json:"id,omitempty"  gorm:"primaryKey"`
	Name            string `json:"name,omitempty"`
	FollowCount     int64  `json:"follow_count"`
	FollowerCount   int64  `json:"follower_count"`
	IsFollow        bool   `json:"is_follow"`
	Avatar          string `json:"avatar,omitempty"`
	BackgroundImage string `json:"background_image,omitempty"`
	Signature       string `json:"signature,omitempty"`
	TotalFavorited  int64  `json:"total_favorited"`
	WorkCount       int64  `json:"work_count"`
	FavoriteCount   int64  `json:"favorite_count"`
}
