package service

import (
	"github.com/ikuraoo/fastdouyin/common"
	"github.com/ikuraoo/fastdouyin/configure"
	"github.com/ikuraoo/fastdouyin/entity"
	"github.com/spf13/viper"
	"strconv"
	"time"
)

func VideoFeed(userId int64, latestTime time.Time) ([]*common.VideoResponse, error) {
	MaxNum := viper.GetString("video.maxNum")
	maxNum, err := strconv.ParseInt(MaxNum, 10, 64)
	if err != nil {
		return nil, err
	}

	videos, err := GetVideos(maxNum, latestTime)
	videosWithUsers, err := CombinationVideosAndUser(userId, videos)
	return videosWithUsers, err
}

func GetVideos(maxNum int64, latestTime time.Time) (*[]entity.Video, error) {
	videos, err := entity.NewVideoDaoInstance().QueryVideos(maxNum, latestTime)
	if err != nil {
		return nil, err
	}
	return videos, nil
}

func CombinationVideosAndUser(UserId int64, videos *[]entity.Video) ([]*common.VideoResponse, error) {
	videoPath := viper.GetString("video.videoPath")
	coverPath := viper.GetString("video.coverPath")
	var video entity.Video
	var author *entity.User
	var err error
	var isfavorite bool

	var videosWithUsers = make([]*common.VideoResponse, 0)
	for _, video = range *videos {
		author, err = entity.NewUserDaoInstance().QueryById(video.UId)
		if err != nil {
			continue
		}
		if UserId == 0 {
			isfavorite = false
		} else {
			isfavorite = configure.NewProxyIndexMap().GetVideoFavor(UserId, video.Id)
		}

		//isfavorite, _ := entity.NewFavoriteDaoInstance().GetVideoFavorState(UserId, video.Id)
		authorMessage, _ := ConvertToUserResponse(author, UserId)
		videoResponse := common.VideoResponse{
			Id:            video.Id,
			Author:        *authorMessage,
			PlayUrl:       videoPath + video.PlayUrl,
			CoverUrl:      coverPath + video.CoverUrl,
			CommentCount:  video.CommentCount,
			FavoriteCount: video.FavoriteCount,
			IsFavorite:    isfavorite,
			Title:         video.Title,
		}
		videosWithUsers = append(videosWithUsers, &videoResponse)
	}
	return videosWithUsers, nil
}
