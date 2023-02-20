package service

import (
	"github.com/ikuraoo/fastdouyin/entity"
	"time"
)

type CommentActionMessage struct {
	userId      int64
	videoId     int64
	commentId   int64
	actionType  int64
	commentText string

	comment *entity.Comment
}

type CommentResponse struct {
	Id         int64
	User       *entity.User
	Content    string
	CreateDate time.Time
}

func CreateComment(userId, videoId int64, commentText string) (*entity.Comment, error) {
	comment := &entity.Comment{
		//Id:         0,
		VId:        videoId,
		UId:        userId,
		Content:    commentText,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
		IsDeleted:  false,
	}
	err := entity.NewCommentDaoInstance().AddCommentAndUpdateCount(comment)
	if err != nil {
		return comment, err
	}

	return comment, nil
}

func DeleteComment(commentId, videoId int64) (*entity.Comment, error) {
	var comment entity.Comment
	err := entity.NewCommentDaoInstance().QueryCommentById(commentId, &comment)
	if err != nil {
		return nil, err
	}
	err = entity.NewCommentDaoInstance().DeleteCommentAndUpdateCountById(commentId, videoId)
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func QueryCommentList(userId, videoId int64) ([]*entity.Comment, error) {
	var comments []*entity.Comment
	err := entity.NewCommentDaoInstance().QueryCommentListByVideoId(videoId, &comments)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func CombinationCommentsAndUsers(comments []*entity.Comment) ([]*CommentResponse, error) {
	var commenresponses []*CommentResponse
	for _, comment := range comments {
		commentUser, _ := entity.NewUserDaoInstance().QueryById(comment.UId)
		commentWithuse := CommentResponse{
			Id:         comment.Id,
			User:       commentUser,
			Content:    comment.Content,
			CreateDate: comment.CreateTime,
		}
		commenresponses = append(commenresponses, &commentWithuse)
	}
	return commenresponses, nil
}
