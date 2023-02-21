package entity

import (
	"errors"
	"github.com/ikuraoo/fastdouyin/configure"
	"gorm.io/gorm"
	"log"
	"sync"
	"time"
)

type Video struct {
	Id            int64     `gorm:"column:id"`
	UId           int64     `gorm:"column:uid"`
	PlayUrl       string    `gorm:"column:play_url"`
	CoverUrl      string    `gorm:"column:cover_url"`
	CommentCount  int64     `gorm:"column:comment_count"`
	FavoriteCount int64     `gorm:"column:favorite_count"`
	Title         string    `gorm:"column:title"`
	CreateTime    time.Time `gorm:"column:create_time"`
	UpdateTime    time.Time `gorm:"column:update_time"`
	IsDeleted     bool      `gorm:"column:is_deleted"`
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
	err := configure.Db.Where("create_time<=?", latestTime).Order("create_time DESC").Limit(int(maxNum)).Find(&videos).Error
	if err != nil {
		return nil, err
	}
	return &videos, nil
}

func (*VideoDao) QueryAuther(vid int64) (uid int64, err error) {
	var video Video
	err = configure.Db.Where("id = ?", vid).Find(&video).Error
	if err != nil {
		return 0, err
	}
	return video.UId, nil
}

func (*VideoDao) QueryByVid(vid int64) (*Video, error) {
	var video Video
	err := configure.Db.Where("id = ?", vid).Find(&video).Error
	if err != nil {
		return nil, err
	}
	return &video, nil
}

func (*VideoDao) IsVideoExistById(id int64) bool {
	var video Video
	err := configure.Db.Where("id = ?", id).First(&video).Error
	if err != nil {
		log.Println(err)
	}
	if video.Id == 0 {
		return false
	}
	return true
}

func (*VideoDao) QueryByAuther(uid int64) (*[]Video, error) {
	var videos []Video
	err := configure.Db.Where("uid = ?", uid).Find(&videos).Error
	if err != nil {
		return nil, errors.New("查询视频失败")
	}
	return &videos, nil
}

func (*VideoDao) CreateVideo(video *Video, uid int64) error {
	return configure.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(video).Error; err != nil {
			return err
		}
		if err := tx.Exec("UPDATE users SET work_count=work_count+1 WHERE id = ?", uid).Error; err != nil {
			return err
		}
		return nil
	})
	err := configure.Db.Create(video).Error
	if err != nil {
		return err
	}
	return nil
}

func (*VideoDao) PlusFavorite(userId, authorId, videoId int64) error {
	return configure.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("UPDATE videos SET favorite_count=favorite_count+1 WHERE id = ?", videoId).Error; err != nil {
			return err
		}
		if err := tx.Exec("UPDATE users SET total_favorited=total_favorited+1 WHERE id = ?", authorId).Error; err != nil {
			return err
		}
		if err := tx.Exec("UPDATE users SET favorite_count=favorite_count+1 WHERE id = ?", userId).Error; err != nil {
			return err
		}
		//if err := tx.Exec("INSERT INTO `favourites` (`uid`,`vid`,'is_favourite','create_time', 'update_time') VALUES (?,?, true, time.Now(), time.Now())", userId, videoId).Error; err != nil {
		if err := tx.Exec("INSERT INTO `favourites` (`uid`,`vid`) VALUES (?,?)", userId, videoId).Error; err != nil {
			return err
		}
		return nil
	})
}

func (*VideoDao) MinusFavorite(userId, authorId, videoId int64) error {

	return configure.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("UPDATE videos SET favorite_count=favorite_count-1 WHERE id = ? and favorite_count > 0", videoId).Error; err != nil {
			return err
		}
		if err := tx.Exec("UPDATE users SET total_favorited=total_favorited-1 WHERE id = ? and favorite_count > 0", authorId).Error; err != nil {
			return err
		}
		if err := tx.Exec("UPDATE users SET favorite_count=favorite_count-1 WHERE id = ? and favorite_count > 0", userId).Error; err != nil {
			return err
		}
		//if err := tx.Exec("INSERT INTO `favourites` (`uid`,`vid`,'is_favourite','create_time', 'update_time') VALUES (?,?, true, time.Now(), time.Now())", userId, videoId).Error; err != nil {
		if err := tx.Exec("DELETE FROM `favourites` WHERE `uid` = ? AND `vid` = ?", userId, videoId).Error; err != nil {
			return err
		}
		return nil
	})
}

func (*VideoDao) QueryFavorVideoListByUserId(userId int64) (*[]Video, error) {
	var videoList []Video
	if err := configure.Db.Raw("SELECT v.* FROM videos v , favourites u WHERE u.uid = ? AND u.vid = v.id", userId).Scan(&videoList).Error; err != nil {
		return nil, err
	}
	return &videoList, nil
}

func (*VideoDao) QueryVideoCountByUserId(userId int64, count *int64) error {
	if count == nil {
		return errors.New("QueryVideoCountByUserId count 空指针")
	}
	return configure.Db.Model(&Video{}).Where("uid=?", userId).Count(count).Error
}
