package models

import (
	"errors"
	"github.com/jinzhu/gorm"
)

// User defines the shape of the user table in the database
type User struct {
	gorm.Model
	Introduction string `gorm:"not null"`
}

// UserDB defines all methods of the user service
type UserDB interface {
	Create(user *User) error
	ByID(token string) (*User, error)
	GetGroupUsersByID(id uint) (*[]User, error)
}

var _ UserDB = &userGorm{}

type userGorm struct {
	db *gorm.DB
}

func (ug *userGorm) GetGroupUsersByID(id uint) (*[]User, error) {
	users := &[]User{}
	if err := ug.db.Where("id <= ? AND id >= ?", id+5, id-5).Find(users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func newUserGorm(db *gorm.DB) *userGorm {
	return &userGorm{
		db: db,
	}
}

// UserService defines all methods of the user service
type UserService interface {
	UserDB
}

type userService struct {
	UserDB
}

// NewUserService returns the userservice struct
func NewUserService(db *gorm.DB) UserService {
	ug := newUserGorm(db)
	return &userService{
		UserDB: ug,
	}
}

func (ug *userGorm) Create(user *User) error {
	return ug.db.Create(&user).Error
}

func (ug *userGorm) ByID(id string) (*User, error) {
	user := &User{}
	if err := ug.db.First(user, "id = ?", id).Error; err != nil {
		return nil, errors.New("not found")
	}
	return user, nil
}
