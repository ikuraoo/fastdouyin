package service

import (
	"github.com/ikuraoo/fastdouyin/common"
	"github.com/ikuraoo/fastdouyin/entity"
	"time"
)

type PublishListMessage struct {
	userId   int64
	authorId int64
}

func VideoPublish(uid int64, title, filename, covername string) error {

	video := &entity.Video{
		//Id:             0,
		UId:           uid,
		PlayUrl:       filename,
		CoverUrl:      covername,
		CommentCount:  0,
		FavoriteCount: 0,
		Title:         title,
		CreateTime:    time.Now(),
		UpdateTime:    time.Now(),
		IsDeleted:     false,
	}
	err := entity.NewVideoDaoInstance().CreateVideo(video, uid)

	if err != nil {
		return err
	}
	return nil
}

func PublishList(userId, authorId int64) ([]*common.VideoResponse, error) {
	var publishListMessage = &PublishListMessage{
		userId:   userId,
		authorId: authorId,
	}
	return publishListMessage.Do()
}

func (p *PublishListMessage) Do() ([]*common.VideoResponse, error) {
	if err := p.check(); err != nil {
		return nil, err
	}

	//查询用户发布视频
	videos, err := entity.NewVideoDaoInstance().QueryByAuther(p.authorId)
	if err != nil {
		return nil, err
	}

	//将视频与用户信息连接，转为可输出格式

	if err != nil {
		return nil, err
	}
	videosWithUsers, err := CombinationVideosAndUser(p.userId, videos)
	if err != nil {
		return nil, err
	}
	return videosWithUsers, nil
}

func (p PublishListMessage) check() error {
	_, err := entity.NewUserDaoInstance().QueryById(p.userId)
	if err != nil {
		return err
	}
	_, err = entity.NewUserDaoInstance().QueryById(p.authorId)
	if err != nil {
		return err
	}
	return nil
}
