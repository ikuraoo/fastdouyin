package entity

import (
	"sync"
	"time"
)

type Favourite struct {
	Id          int64 `gorm:"column:id"`
	UId         int64 `gorm:"column:uid"`
	VId         int64 `gorm:"column:vid"`
	IsFavourite bool
	CreateTime  time.Time
	UpdateTime  time.Time
}

type FavoriteDao struct {
}

var favoriteDao *FavoriteDao
var favoriteOnce sync.Once

func NewFavoriteDaoInstance() *FavoriteDao {
	favoriteOnce.Do(
		func() {
			favoriteDao = &FavoriteDao{}
		})
	return favoriteDao
}

func (*FavoriteDao) CreateFavourite(favourite *Favourite) error {
	err := db.Create(favourite).Error
	if err != nil {
		return err
	}
	return nil
}

func (*FavoriteDao) DeleteFavouriteByUidVid(uid, vid int64) error {
	err := db.Where("uid = ? and vid = ?", uid, vid).Delete(Favourite{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (*FavoriteDao) GetVideoFavorState(uid, vid int64) (bool, error) {
	err := db.Where("uid = ? and vid = ?", uid, vid).First(Favourite{}).Error
	if err != nil {
		return false, err
	}
	return true, nil
}
