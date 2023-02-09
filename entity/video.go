package entity

import "time"

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
