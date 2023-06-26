package models

import (
	"github.com/jinzhu/gorm"
)

// Comment defines the shape of the comment model
type Comment struct {
	gorm.Model
	UserID  uint   `gorm:"not null"`
	User    User   `gorm:"foreignKey:UserID"`
	ChatID  uint   `gorm:"not null"`
	Content string `gorm:"not null"`
}

// CommentService interface
type CommentService interface {
	commentDB
}

type commentDB interface {
	Create(comment *Comment) error
}

type commentService struct {
	commentDB
}

type commentGorm struct {
	db *gorm.DB
}

var _ commentDB = &commentGorm{}
var _ CommentService = &commentService{}

func newCommentGorm(db *gorm.DB) *commentGorm {
	return &commentGorm{
		db: db,
	}
}

// NewCommentService returns the CommentService interface
func NewCommentService(db *gorm.DB) CommentService {
	gb := newCommentGorm(db)
	return &commentService{
		commentDB: gb,
	}
}

func (pg *commentGorm) Create(comment *Comment) error {
	err := pg.db.Create(comment).Error
	if err != nil {
		return err
	}
	return nil
}

/*
func (pg *commentGorm) FindByUserID(id uint) (*[]Chat, error) {
	posts := &[]Chat{}
	if err := pg.db.Find(posts, "user_id = ?", id).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

*/
