package entity

import "time"

type Favourite struct {
	Id          int64 `gorm:"column:id"`
	UId         int64 `gorm:"column:uid"`
	VId         int64 `gorm:"column:vid"`
	IsFavourite bool
	CreateTime  time.Time
	UpdateTime  time.Time
}
