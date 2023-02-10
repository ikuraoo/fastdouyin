package service

import (
	"fmt"
	"github.com/ikuraoo/fastdouyin/constant"
	"github.com/ikuraoo/fastdouyin/entity"
	"github.com/spf13/viper"
	"strconv"
)

type AuthorUser struct {
	Id            int64
	Name          string
	FollowCount   int64
	FollowerCount int64
	IsFollow      bool
}

type VideoWithUser struct {
	Id             int64
	Author         *AuthorUser
	PlayUrl        string `json:"play_url" json:"play_url,omitempty"`
	CoverUrl       string `json:"cover_url,omitempty"`
	CommentCount   int64
	FavouriteCount int64
	Title          string
}

func VideoFeed(myUId int64, LatestTime int64) ([]*VideoWithUser, error) {
	MaxNum := viper.GetString("video.maxNum")
	maxNum, err := strconv.ParseInt(MaxNum, 10, 64)
	if err != nil {
		return nil, err
	}

	videos, err := GetVideos(maxNum)
	videosWithUsers, err := CombinationVideosAndUsers(myUId, videos)
	return videosWithUsers, err
}

func GetVideos(maxNum int64) (*[]entity.Video, error) {
	videos, err := entity.NewVideoDaoInstance().QueryVideos(maxNum)
	if err != nil {
		return nil, err
	}
	return videos, nil
}

func CombinationVideosAndUsers(myUId int64, videos *[]entity.Video) ([]*VideoWithUser, error) {
	videoPath := viper.GetString("video.videoPath")
	coverPath := viper.GetString("video.coverPath")
	var video entity.Video
	var author *entity.User
	var err error
	var videosWithUsers = make([]*VideoWithUser, 0)
	var isFollow bool
	for _, video = range *videos {
		author, err = entity.NewUserDaoInstance().QueryById(video.UId)
		if err != nil {
			continue
		}
		if myUId != constant.WRONG_ID {
			isFollow, err = entity.NewFollowDaoInstance().QueryIsFollow(myUId, author.Id)
			if err != nil {
				fmt.Println(err)
			}
		}
		videoWithUser := VideoWithUser{
			Id: video.Id,
			Author: &AuthorUser{
				Id:            author.Id,
				Name:          author.Name,
				FollowCount:   author.FollowCount,
				FollowerCount: author.FollowerCount,
				IsFollow:      isFollow,
			},
			PlayUrl:        videoPath + video.PlayUrl,
			CoverUrl:       coverPath + video.CoverUrl,
			CommentCount:   video.CommentCount,
			FavouriteCount: video.FavouriteCount,
			Title:          video.Title,
		}
		videosWithUsers = append(videosWithUsers, &videoWithUser)
	}
	return videosWithUsers, nil
}

func PublishList(uid int64) ([]*VideoWithUser, error) {
	videos, err := entity.NewVideoDaoInstance().QueryByAuther(uid)
	if err != nil {
		return nil, err
	}
	videosWithUsers, err := CombinationVideosAndUsers(uid, videos)
	return videosWithUsers, nil
}
