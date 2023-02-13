package entity

import (
	"sync"
	"time"
)

type Video struct {
	Id             int64     `gorm:"column:id"`
	UId            int64     `gorm:"column:uid"`
	PlayUrl        string    `gorm:"column:play_url"`
	CoverUrl       string    `gorm:"column:cover_url"`
	CommentCount   int64     `gorm:"column:comment_count"`
	FavouriteCount int64     `gorm:"column:favourite_count"`
	Title          string    `gorm:"column:title"`
	CreateTime     time.Time `gorm:"column:create_time"`
	UpdateTime     time.Time `gorm:"column:update_time"`
	IsDeleted      bool      `gorm:"column:is_deleted"`
}

type VideoDao struct {
}

var videoDao *VideoDao
var videoOnce sync.Once

func NewVideoDaoInstance() *VideoDao {
	videoOnce.Do(
		func() {
			videoDao = &VideoDao{}
		})
	return videoDao
}

func (*VideoDao) QueryVideos(maxNum int64) (*[]Video, error) {
	var videos []Video
	err := db.Limit(int(maxNum)).Find(&videos).Error
	if err != nil {
		return nil, err
	}
	return &videos, nil
}

func (*VideoDao) QueryByAuther(uid int64) (*[]Video, error) {
	var videos []Video
	err := db.Where("uid = ?", uid).Find(&videos).Error
	if err != nil {
		return nil, err
	}
	return &videos, nil
}

func (*VideoDao) CreateVideo(video *Video) error {
	err := db.Create(video).Error
	if err != nil {
		return err
	}
	return nil
}
