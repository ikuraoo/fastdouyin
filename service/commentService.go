package service

import (
	"errors"
	"github.com/ikuraoo/fastdouyin/common"
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

func CreateComment(userId, videoId int64, commentText string) (*common.Comment, error) {
	//检查参数
	user, err := entity.NewUserDaoInstance().QueryById(userId)
	if err != nil {
		return nil, errors.New("user don't exist")
	}

	//类型转换
	userResponse, err := ConvertToUserResponse(user, 0)
	if err != nil {
		return nil, errors.New("类型转换失败")
	}
	exist := entity.NewVideoDaoInstance().IsVideoExistById(videoId)
	if !exist {
		return nil, errors.New("video don't exist")
	}

	comment := &entity.Comment{
		//Id:         0,
		VId:        videoId,
		UId:        userId,
		Content:    commentText,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
		IsDeleted:  false,
	}
	err = entity.NewCommentDaoInstance().AddCommentAndUpdateCount(comment)
	if err != nil {
		return nil, err
	}

	commentResponse, err := CombinationCommentAndAuthor(comment, userResponse)
	if err != nil {
		return commentResponse, err
	}
	return commentResponse, nil
}

func DeleteComment(commentId, videoId int64) (*common.Comment, error) {
	//检查参数
	exist := entity.NewCommentDaoInstance().IsCommenExistById(commentId)
	if !exist {
		return nil, errors.New("user don't exist")
	}
	exist = entity.NewVideoDaoInstance().IsVideoExistById(videoId)
	if !exist {
		return nil, errors.New("video don't exist")
	}

	var comment *entity.Comment
	comment, err := entity.NewCommentDaoInstance().QueryCommentById(commentId)
	if err != nil {
		return nil, err
	}
	err = entity.NewCommentDaoInstance().DeleteCommentAndUpdateCountById(commentId, videoId)
	author, err := entity.NewUserDaoInstance().QueryById(comment.UId)
	authorMessage, err := ConvertToUserResponse(author, 0)

	if err != nil {
		return nil, err
	}
	commentResponse, err := CombinationCommentAndAuthor(comment, authorMessage)

	if err != nil {
		return nil, err
	}
	return commentResponse, nil
}

func QueryCommentList(videoId int64) ([]*entity.Comment, error) {
	var comments []*entity.Comment
	err := entity.NewCommentDaoInstance().QueryCommentListByVideoId(videoId, &comments)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func CombinationCommentAndAuthor(comment *entity.Comment, user *common.UserResponse) (*common.Comment, error) {
	commentResponse := common.Comment{
		Id:         comment.Id,
		User:       *user,
		Content:    comment.Content,
		CreateDate: comment.CreateTime.Format("1-2"),
	}
	return &commentResponse, nil
}
func CombinationCommentsAndUsers(comments []*entity.Comment) ([]*common.Comment, error) {
	var commenresponses []*common.Comment
	for _, comment := range comments {
		commentUser, _ := entity.NewUserDaoInstance().QueryById(comment.UId)
		commontUserMassage, _ := ConvertToUserResponse(commentUser, 0)
		commentWithuse := common.Comment{
			Id:         comment.Id,
			User:       *commontUserMassage,
			Content:    comment.Content,
			CreateDate: comment.CreateTime.Format("1-2"),
		}
		commenresponses = append(commenresponses, &commentWithuse)
	}
	return commenresponses, nil
}
