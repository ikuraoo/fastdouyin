package entity

import (
	"github.com/ikuraoo/fastdouyin/common"
	"gorm.io/gorm"
	"strconv"
	"sync"
	"time"
)

type Follow struct {
	ID         int64     `gorm:"primarykey"`
	MyId       int64     `gorm:"column:my_uid"`
	HisId      int64     `gorm:"column:his_uid"`
	IsFollow   bool      `gorm:"column:is_follow"`
	CreateTime time.Time `gorm:"column:create_time"`
	UpdateTime time.Time `gorm:"column:update_time"`
}

type FollowDao struct {
}

var followDao *FollowDao
var followOnce sync.Once

func NewFollowDaoInstance() *FollowDao {
	followOnce.Do(
		func() {
			followDao = &FollowDao{}
		})
	return followDao
}

func (*FollowDao) QueryIsFollow(myUId int64, hisUId int64) (bool, error) {
	var follow Follow
	err := db.Where("my_uid = ?", strconv.Itoa(int(myUId))).Where("his_uid = ?", strconv.Itoa(int(hisUId))).Find(&follow).Error
	if err != nil {
		return false, err
	}
	if follow.ID == common.WRONG_ID {
		return false, nil
	}
	return true, nil
}

func (*FollowDao) AddUserFollow(userId, followId int64) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("UPDATE users SET follow_count=follow_count+1 WHERE id = ?", userId).Error; err != nil {
			return err
		}
		if err := tx.Exec("UPDATE users SET follower_count=follower_count+1 WHERE id = ?", followId).Error; err != nil {
			return err
		}
		if err := tx.Exec("INSERT INTO `follows` (`my_uid`,`his_uid`) VALUES (?,?)", userId, followId).Error; err != nil {
			return err
		}
		return nil
	})
}
func (*FollowDao) CancelUserFollow(userId, followId int64) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("UPDATE users SET follow_count=follow_count-1 WHERE id = ? AND follow_count>0", userId).Error; err != nil {
			return err
		}
		if err := tx.Exec("UPDATE users SET follower_count=follower_count-1 WHERE id = ? AND follower_count>0", followId).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM `follows` WHERE my_uid=? AND his_uid=?", userId, followId).Error; err != nil {
			return err
		}
		return nil
	})
}
