package entity

import (
	"errors"
	"github.com/ikuraoo/fastdouyin/configure"
	"gorm.io/gorm"
	"log"
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

func (c *CommentDao) QueryCommentById(commentId int64) (*Comment, error) {
	var comment Comment
	err := configure.Db.Where("id = ?", commentId).First(&comment).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (c *CommentDao) IsCommenExistById(id int64) bool {
	var comment Comment
	err := configure.Db.Where("id = ?", id).First(&comment).Error
	if err != nil {
		log.Println(err)
	}
	if comment.Id == 0 {
		return false
	}
	return true
}
func (c *CommentDao) AddCommentAndUpdateCount(comment *Comment) error {

	//执行事务
	return configure.Db.Transaction(func(tx *gorm.DB) error {
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
	return configure.Db.Transaction(func(tx *gorm.DB) error {
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
	if err := configure.Db.Model(&Comment{}).Where("vid=?", videoId).Find(comments).Error; err != nil {
		return err
	}
	return nil
}
