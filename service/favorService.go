package service

import (
	"errors"
	"github.com/ikuraoo/fastdouyin/common"
	"github.com/ikuraoo/fastdouyin/entity"
)

const (
	PLUS  = 1
	MINUS = 2
)

type FavoriteActionMessage struct {
	userId     int64
	authorId   int64
	videoId    int64
	actionType int64
}

func NewFavoriteAction(userId, videoId, action int64) *FavoriteActionMessage {
	return &FavoriteActionMessage{
		userId:     userId,
		videoId:    videoId,
		actionType: action,
	}

}
func UserFavoriteAction(userId, videoId, actionType int64) error {
	return NewFavoriteAction(userId, videoId, actionType).Do()
}

func (f *FavoriteActionMessage) Do() error {
	var err error
	if err = f.checkNum(); err != nil {
		return err
	}

	switch f.actionType {
	case PLUS:
		err = f.PlusOperation()
	case MINUS:
		err = f.MinusOperation()
	default:
		return errors.New("未定义的行为")
	}
	return err
}

func (f *FavoriteActionMessage) checkNum() error {
	authorId, err := entity.NewVideoDaoInstance().QueryAuther(f.videoId)
	if err != nil {
		return err
	}
	f.authorId = authorId

	if err != nil {
		return errors.New("用户不存在")
	}

	if f.actionType != PLUS && f.actionType != MINUS {
		return errors.New("未定义的行为")
	}
	return nil
}

func (f *FavoriteActionMessage) PlusOperation() error {
	//视频点赞数目+1
	exist, _ := entity.NewFavoriteDaoInstance().GetVideoFavorState(f.userId, f.videoId)
	if !exist {
		err := entity.NewVideoDaoInstance().PlusFavorite(f.userId, f.authorId, f.videoId)
		if err != nil {
			return err
		}
	} else {
		return errors.New("已经点过赞")
	}

	return nil
}

func (f *FavoriteActionMessage) MinusOperation() error {
	//视频点赞数目-1
	exist, _ := entity.NewFavoriteDaoInstance().GetVideoFavorState(f.userId, f.videoId)
	if exist {
		err := entity.NewVideoDaoInstance().MinusFavorite(f.userId, f.authorId, f.videoId)
		if err != nil {
			return err
		}
	} else {
		return errors.New("未点赞")
	}
	return nil
}

func QueryUserFavorList(userId, targetId int64) ([]*common.VideoResponse, error) {
	//检查user是否存在
	_, err := entity.NewUserDaoInstance().QueryById(userId)
	if err != nil {
		return nil, errors.New("user don't exist")
	}
	_, err = entity.NewUserDaoInstance().QueryById(targetId)
	if err != nil {
		return nil, errors.New("user don't exist")
	}

	//查询点赞视频
	videos, err := entity.NewVideoDaoInstance().QueryFavorVideoListByUserId(targetId)
	if err != nil {
		return nil, err
	}

	//视频格式封装
	videosWithUsers, err := CombinationVideosAndUser(userId, videos)
	if err != nil {
		return nil, err
	}
	return videosWithUsers, nil
}
