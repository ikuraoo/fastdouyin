package entity

import "gorm.io/gorm"

type Video struct {
	gorm.Model
	UId            int64  `gorm:"column:uid"`
	PlayUrl        string `gorm:"column:play_url"`
	CoverUrl       string `gorm:"column:cover_url"`
	CommentCount   int64  `gorm:"column:comment_count"`
	FavouriteCount int64  `gorm:"column:favourite_count"`
	Title          string
}
