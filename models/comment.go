package models

import (
	"context"
	"errors"
	"github.com/PullRequestInc/go-gpt3"
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
	resp, err := Client.Moderation(context.TODO(), gpt3.ModerationRequest{
		Input: comment.Content,
		Model: "text-moderation-latest",
	})
	if err != nil {
		return err
	}
	if len(resp.Results) == 0 {
		return errors.New("nil resp")
	}
	if resp.Results[0].Flagged {
		return errors.New("violent")
	}
	err = pg.db.Create(comment).Error
	if err != nil {
		return err
	}
	return nil
}
