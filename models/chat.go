package models

import (
	"errors"
	"github.com/jinzhu/gorm"
)

var (
	// ErrPostMissing is returned when post field is missing
	ErrPostMissing = errors.New("models: Provide a Chat")
	// ErrInvalidID is returned when an invalid ID is passed
	ErrInvalidID = errors.New("models: Chat not found")
)

// Chat defines the shape of the post model
type Chat struct {
	gorm.Model
	UserID   uint   `gorm:"not null"`
	User     User   `gorm:"foreignKey:UserID"`
	Content  string `gorm:"not null"`
	Role     string `gorm:"not null"`
	Comments []Comment
}

// ChatService interface
type ChatService interface {
	postDB
}

type postDB interface {
	Create(post *Chat) error
	Delete(post *Chat) error
	FindByUserID(id uint) (*[]Chat, error)
	FindPostByID(id uint) (*Chat, error)
}

type postService struct {
	postDB
}

type postVal struct {
	postDB
}

type postGorm struct {
	db *gorm.DB
}

var _ postDB = &postGorm{}
var _ ChatService = &postService{}

func newPostGorm(db *gorm.DB) *postGorm {
	return &postGorm{
		db: db,
	}
}

func newPostVal(pg *postGorm) *postVal {
	return &postVal{
		postDB: pg,
	}
}

// NewChatService returns the ChatService interface
func NewChatService(db *gorm.DB) ChatService {
	pg := newPostGorm(db)
	pv := newPostVal(pg)
	return &postService{
		postDB: pv,
	}
}

type postValFn func(post *Chat) error

func runPostValFns(post *Chat, fns ...postValFn) error {
	for _, fn := range fns {
		if err := fn(post); err != nil {
			return err
		}
	}
	return nil
}

func (pv *postVal) checkForPost(post *Chat) error {
	if post.Content == "" {
		return ErrPostMissing
	}
	return nil
}
func (pv *postVal) checkID(post *Chat) error {
	if post.ID == 0 {
		return ErrInvalidID
	}
	return nil
}

func (pv *postVal) Create(post *Chat) error {
	if err := runPostValFns(post); err != nil {
		return err
	}
	return pv.postDB.Create(post)
}

func (pg *postGorm) Create(post *Chat) error {
	return pg.db.Create(post).Error
}

func (pg *postGorm) FindByUserID(id uint) (*[]Chat, error) {
	var chats []Chat
	result := pg.db.Preload("Comments").Preload("User").Where("user_id = ?", id).Find(&chats)

	if result.Error != nil {
		return nil, result.Error
	}

	return &chats, nil
}

func (pg *postGorm) FindPostByID(id uint) (*Chat, error) {
	post := &Chat{}
	if err := pg.db.First(post, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return post, nil
}

func (pv *postVal) Delete(post *Chat) error {
	if err := runPostValFns(post, pv.checkID); err != nil {
		return err
	}
	return pv.postDB.Delete(post)
}

func (pg *postGorm) Delete(post *Chat) error {
	return pg.db.Delete(post, "id = ?", post.ID).Error
}
