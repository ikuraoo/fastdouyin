package entity

import "time"

type Follow struct {
	ID         uint      `gorm:"primarykey"`
	MyId       int64     `gorm:"column:my_uid"`
	HisId      int64     `gorm:"column:his_uid"`
	IsFollow   bool      `gorm:"column:is_follow"`
	CreateTime time.Time `gorm:"column:create_time"`
	UpdateTime time.Time `gorm:"column:update_time"`
}
