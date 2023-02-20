package entity

import (
	"errors"
	"gorm.io/gorm"
	"sync"
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

type CommentDao struct {
}

var commentDao *CommentDao
var commentOnce sync.Once

func NewCommentDaoInstance() *CommentDao {
	commentOnce.Do(
		func() {
			commentDao = &CommentDao{}
		})
	return commentDao
}

func (c *CommentDao) QueryCommentById(commentId int64, comment *Comment) error {
	if comment == nil {
		return errors.New("QueryCommentById comment 空指针")
	}
	return db.Where("id = ?", commentId).First(comment).Error
}
func (c *CommentDao) AddCommentAndUpdateCount(comment *Comment) error {

	//执行事务
	return db.Transaction(func(tx *gorm.DB) error {
		//添加评论数据
		if err := tx.Create(comment).Error; err != nil {
			// 返回任何错误都会回滚事务
			return err
		}
		//增加count
		if err := tx.Exec("UPDATE videos v SET v.comment_count = v.comment_count+1 WHERE v.id=?", comment.VId).Error; err != nil {
			return err
		}
		// 返回 nil 提交事务
		return nil
	})
}

func (c *CommentDao) DeleteCommentAndUpdateCountById(commentId, videoId int64) error {
	//执行事务
	return db.Transaction(func(tx *gorm.DB) error {
		//删除评论
		if err := tx.Exec("DELETE FROM comments WHERE id = ?", commentId).Error; err != nil {
			// 返回任何错误都会回滚事务
			return err
		}
		//减少count
		if err := tx.Exec("UPDATE videos v SET v.comment_count = v.comment_count-1 WHERE v.id=? AND v.comment_count>0", videoId).Error; err != nil {
			return err
		}
		// 返回 nil 提交事务
		return nil
	})
}

func (c *CommentDao) QueryCommentListByVideoId(videoId int64, comments *[]*Comment) error {
	if comments == nil {
		return errors.New("QueryCommentListByVideoId comments空指针")
	}
	if err := db.Model(&Comment{}).Where("vid=?", videoId).Find(comments).Error; err != nil {
		return err
	}
	return nil
}
