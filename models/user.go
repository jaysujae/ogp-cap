package models

import (
	"errors"
	"github.com/jinzhu/gorm"
	"math/rand"
)

// User defines the shape of the user table in the database
type User struct {
	gorm.Model
	Nickname     string
	Introduction string `gorm:"not null"`
	Image        string
	Chats        []Chat
	Comments     []Comment
}

// UserService defines all methods of the user service
type UserService interface {
	Create(user *User) error
	ByID(token string) (*User, error)
	GetGroupUsersByID(id uint) (*[]User, error)
}

type UserGorm struct {
	db *gorm.DB
}

func NewUserGorm(db *gorm.DB) *UserGorm {
	return &UserGorm{
		db: db,
	}
}
func (ug *UserGorm) GetGroupUsersByID(id uint) (*[]User, error) {
	users := &[]User{}
	if err := ug.db.Where("id <= ? AND id >= ?", id+5, id-5).Find(users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (ug *UserGorm) Create(user *User) error {
	adjectives := []string{
		"Happy",
		"Bright",
		"Jolly",
		"Kind",
		"Brave",
		"Cool",
		"Wise",
		"Calm",
		"Bold",
		"Fit",
	}
	animals := []string{
		"Bear",
		"Deer",
		"Donkey",
		"Elephant",
		"Fox",
		"Monkey",
		"Panda",
		"Rabbit",
		"Squirrel",
		"Zebra",
	}
	adjRand := rand.Intn(10)
	nameRand := rand.Intn(10)
	user.Nickname = adjectives[adjRand] + " " + animals[nameRand]
	user.Image = animals[nameRand] + ".jpg"
	return ug.db.Create(&user).Error
}

func (ug *UserGorm) ByID(id string) (*User, error) {
	user := &User{}
	if err := ug.db.First(user, "id = ?", id).Error; err != nil {
		return nil, errors.New("not found")
	}
	return user, nil
}
