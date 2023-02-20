package service

import (
	"github.com/ikuraoo/fastdouyin/common"
	"github.com/ikuraoo/fastdouyin/entity"
	"github.com/spf13/viper"
	"strconv"
	"time"
)

func VideoFeed(userId int64, latestTime time.Time) ([]*common.VideoMessage, error) {
	MaxNum := viper.GetString("video.maxNum")
	maxNum, err := strconv.ParseInt(MaxNum, 10, 64)
	if err != nil {
		return nil, err
	}

	videos, err := GetVideos(maxNum, latestTime)
	videosWithUsers, err := CombinationVideosAndUsers(userId, videos)
	return videosWithUsers, err
}

func GetVideos(maxNum int64, latestTime time.Time) (*[]entity.Video, error) {
	videos, err := entity.NewVideoDaoInstance().QueryVideos(maxNum, latestTime)
	if err != nil {
		return nil, err
	}
	return videos, nil
}

func CombinationVideosAndUsers(UserId int64, videos *[]entity.Video) ([]*common.VideoMessage, error) {
	videoPath := viper.GetString("video.videoPath")
	coverPath := viper.GetString("video.coverPath")
	var video entity.Video
	var author *entity.User
	var err error
	var videosWithUsers = make([]*common.VideoMessage, 0)
	for _, video = range *videos {
		author, err = entity.NewUserDaoInstance().QueryById(video.UId)
		if err != nil {
			continue
		}
		isfavorite, _ := entity.NewFavoriteDaoInstance().GetVideoFavorState(UserId, video.Id)
		authorMessage, _ := ConvertToUserMessage(author, UserId)
		videoMessage := common.VideoMessage{
			Id:            video.Id,
			Author:        *authorMessage,
			PlayUrl:       videoPath + video.PlayUrl,
			CoverUrl:      coverPath + video.CoverUrl,
			CommentCount:  video.CommentCount,
			FavoriteCount: video.FavouriteCount,
			IsFavorite:    isfavorite,
			Title:         video.Title,
		}
		videosWithUsers = append(videosWithUsers, &videoMessage)
	}
	return videosWithUsers, nil
}

func VideoPublish(uid int64, title, filename, covername string) error {

	video := &entity.Video{
		//Id:             0,
		UId:            uid,
		PlayUrl:        filename,
		CoverUrl:       covername,
		CommentCount:   0,
		FavouriteCount: 0,
		Title:          title,
		CreateTime:     time.Now(),
		UpdateTime:     time.Now(),
		IsDeleted:      false,
	}
	err := entity.NewVideoDaoInstance().CreateVideo(video, uid)

	if err != nil {
		return err
	}
	return nil
}

func PublishList(userId int64) ([]*common.VideoMessage, error) {
	//检查用户是否存在
	_, err := entity.NewUserDaoInstance().QueryById(userId)
	if err != nil {
		return nil, err
	}

	//查询用户发布视频
	videos, err := entity.NewVideoDaoInstance().QueryByAuther(userId)
	if err != nil {
		return nil, err
	}

	//将视频与用户信息连接，转为可输出格式
	videosWithUsers, err := CombinationVideosAndUsers(userId, videos)
	if err != nil {
		return nil, err
	}
	return videosWithUsers, nil
}
