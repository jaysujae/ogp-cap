package models

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PullRequestInc/go-gpt3"
	"github.com/jinzhu/gorm"
	"math"
	"math/rand"
	"sort"
)

type CustomFloats struct {
	Floats []float64
}

func (c *CustomFloats) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &c.Floats)
}

func (c CustomFloats) Value() (driver.Value, error) {
	return json.Marshal(c.Floats)
}

// User defines the shape of the user table in the database
type User struct {
	gorm.Model
	Nickname          string
	Introduction      string `gorm:"not null"`
	Image             string
	CustomFloatsValue CustomFloats `gorm:"type:json"`
	Chats             []Chat
	Comments          []Comment
}

// UserService defines all methods of the user service
type UserService interface {
	Create(user *User) error
	ByID(token string) (*User, error)
	GetGroupUsersByID(id uint) (*[]User, error)
	Delete(id uint) error
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
	var currentUser User
	if err := ug.db.Where("id = ?", id).First(&currentUser).Error; err != nil {
		return nil, err
	}

	users := []User{}
	if err := ug.db.Find(&users).Error; err != nil {
		return nil, err
	}

	type UserWithValue struct {
		User
		Value float64
	}

	usersWithValue := make([]UserWithValue, len(users))
	for i, user := range users {
		val, _ := Cosine(user.CustomFloatsValue.Floats, currentUser.CustomFloatsValue.Floats)
		usersWithValue[i] = UserWithValue{User: user, Value: val}
	}

	sort.Slice(usersWithValue, func(i, j int) bool {
		return usersWithValue[i].Value > usersWithValue[j].Value
	})
	for i := range users {
		users[i] = usersWithValue[i].User
	}

	return &users, nil
}

func Cosine(a []float64, b []float64) (cosine float64, err error) {
	count := 0
	length_a := len(a)
	length_b := len(b)
	if length_a > length_b {
		count = length_a
	} else {
		count = length_b
	}
	sumA := 0.0
	s1 := 0.0
	s2 := 0.0
	for k := 0; k < count; k++ {
		if k >= length_a {
			s2 += math.Pow(b[k], 2)
			continue
		}
		if k >= length_b {
			s1 += math.Pow(a[k], 2)
			continue
		}
		sumA += a[k] * b[k]
		s1 += math.Pow(a[k], 2)
		s2 += math.Pow(b[k], 2)
	}
	if s1 == 0 || s2 == 0 {
		return 0.0, errors.New("Vectors should not be null (all zeros)")
	}
	return sumA / (math.Sqrt(s1) * math.Sqrt(s2)), nil
}

func (ug *UserGorm) Create(user *User) error {
	if len(user.Introduction) < 10 {
		return errors.New("Please write more than 10 letters.")
	}
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

	embedding, _ := Client.Embeddings(context.TODO(), gpt3.EmbeddingsRequest{
		Input: []string{user.Introduction},
		Model: "text-embedding-ada-002",
		User:  fmt.Sprint(user.ID),
	})
	user.CustomFloatsValue = CustomFloats{Floats: embedding.Data[0].Embedding}

	return ug.db.Create(&user).Error
}

func (ug *UserGorm) ByID(id string) (*User, error) {
	user := &User{}
	if err := ug.db.First(user, "id = ?", id).Error; err != nil {
		return nil, errors.New("not found")
	}
	return user, nil
}

func (ug *UserGorm) Delete(id uint) error {
	user := &User{}
	if err := ug.db.First(user, "id = ?", id).Error; err != nil {
		errors.New("not found")
	}
	ug.db.Delete(&user) // Delete the user
	return nil
}
