package service

import (
	"errors"
	"github.com/ikuraoo/fastdouyin/common"
	"github.com/ikuraoo/fastdouyin/entity"
	"time"
)

const (
	PLUS  = 1
	MINUS = 2
)

type FavoriteActionMessage struct {
	userId     int64
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
func UserFavorite(userId, videoId, actionType int64) error {
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
	if _, err := entity.NewUserDaoInstance().QueryById(f.userId); err != nil {
		return errors.New("用户不存在")
	}
	if f.actionType != PLUS && f.actionType != MINUS {
		return errors.New("未定义的行为")
	}
	return nil
}

func (f *FavoriteActionMessage) PlusOperation() error {
	//视频点赞数目+1
	err := entity.NewVideoDaoInstance().PlusFavorite(f.videoId)
	if err != nil {
		return errors.New("不要重复点赞")
	}
	//对应的用户是否点赞的映射状态更新
	favorite := &entity.Favourite{
		//Id:          ,
		UId:         f.userId,
		VId:         f.videoId,
		IsFavourite: true,
		CreateTime:  time.Now(),
		UpdateTime:  time.Now(),
	}
	err = entity.NewFavoriteDaoInstance().CreateFavourite(favorite)
	if err != nil {
		return err
	}
	return nil
}

func (f *FavoriteActionMessage) MinusOperation() error {
	//视频点赞数目-1
	err := entity.NewVideoDaoInstance().MinusFavorite(f.videoId)
	if err != nil {
		return errors.New("点赞数目已经为0")
	}
	//对应的用户是否点赞的映射状态更新
	err = entity.NewFavoriteDaoInstance().DeleteFavouriteByUidVid(f.userId, f.videoId)
	return nil
}

func QueryUserFavorList(userId int64) ([]*common.VideoMessage, error) {
	///var videos *[]entity.Video
	videos, err := entity.NewVideoDaoInstance().QueryFavorVideoListByUserId(userId)
	if err != nil {
		return nil, err
	}
	videosWithUsers, _ := CombinationVideosAndUsers(userId, videos)
	return videosWithUsers, nil
}
