package models

import (
	"github.com/jinzhu/gorm"
)

// Services struct defines all services
type Services struct {
	db      *gorm.DB
	User    UserService
	Chat    ChatService
	Comment CommentService
}

// NewServices returns the services struct
func NewServices(connectionString string) (*Services, error) {
	db, err := gorm.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	return &Services{
		User:    NewUserService(db),
		Chat:    NewChatService(db),
		Comment: NewCommentService(db),
		db:      db,
	}, nil
}

// AutoMigrate automatically creates the table in the database
func (s *Services) AutoMigrate() error {
	return s.db.AutoMigrate(&User{}, &Chat{}, &Comment{}).Error
}

// DestroyAndCreate drops all tables and recreates
func (s *Services) DestroyAndCreate() error {
	if err := s.db.DropTableIfExists(&User{}, &Chat{}, &Comment{}).Error; err != nil {
		return err
	}
	return s.AutoMigrate()
}

// Close closes connection to the database
func (s *Services) Close() error {
	return s.db.Close()
}
