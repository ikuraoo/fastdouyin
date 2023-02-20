package entity

import (
	"errors"
	"gorm.io/gorm"
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

func (*VideoDao) QueryVideos(maxNum int64, latestTime time.Time) (*[]Video, error) {
	var videos []Video
	err := db.Where("create_time<=?", latestTime).Order("create_time DESC").Limit(int(maxNum)).Find(&videos).Error
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

func (*VideoDao) CreateVideo(video *Video, uid int64) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(video).Error; err != nil {
			return err
		}
		if err := tx.Exec("UPDATE users SET work_count=work_count+1 WHERE id = ?", uid).Error; err != nil {
			return err
		}
		return nil
	})
	err := db.Create(video).Error
	if err != nil {
		return err
	}
	return nil
}

func (*VideoDao) PlusFavorite(videoId int64) error {
	err := db.Model(&Video{}).Where("id = ?", videoId).UpdateColumn("favourite_count", gorm.Expr("favourite_count + ?", 1)).Error
	if err != nil {
		return err
	}
	return nil
}

func (*VideoDao) MinusFavorite(videoId int64) error {
	err := db.Model(&Video{}).Where("id = ? and favourite_count > 0", videoId).UpdateColumn("favourite_count", gorm.Expr("favourite_count - ?", 1)).Error
	if err != nil {
		return err
	}
	return nil
}

func (*VideoDao) QueryFavorVideoListByUserId(userId int64) (*[]Video, error) {
	var videoList []Video
	if err := db.Raw("SELECT v.* FROM videos v , favourites u WHERE u.uid = ? AND u.vid = v.id", userId).Scan(&videoList).Error; err != nil {
		return nil, err
	}
	return &videoList, nil
}

func (*VideoDao) QueryVideoCountByUserId(userId int64, count *int64) error {
	if count == nil {
		return errors.New("QueryVideoCountByUserId count 空指针")
	}
	return db.Model(&Video{}).Where("uid=?", userId).Count(count).Error
}
