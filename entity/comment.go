package entity

import (
	"time"
)

type Comment struct {
	Id         int64 `gorm:"column:id"`
	VId        int64 `gorm:"column:vid"`
	UId        int64 `gorm:"column:uid"`
	Content    string
	CreateTime time.Time
	UpdateTime time.Time
	IsDeleted  bool
}
