package controller

import (
	"fmt"
	"hackathon/middleware"
	"hackathon/models"
	"hackathon/views"
	"net/http"
)

// User defines the shape of a user
type User struct {
	us         models.UserService
	cs         models.ChatService
	SignUpView *views.Views
}

type signUpForm struct {
	Introduction string `schema:"introduction"`
}

// NewUser returns the User struct
func NewUser(us models.UserService, cs models.ChatService) *User {
	return &User{
		us:         us,
		cs:         cs,
		SignUpView: views.NewView("bootstrap", "user/signup"),
	}
}

// LogOut handles the /logout GET
func (u *User) LogOut(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:  middleware.BrowserCookieName,
		Value: "",
	}
	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/", http.StatusFound)
}

// SignUp handles the /signup GET
func (u *User) SignUp(w http.ResponseWriter, r *http.Request) {
	u.SignUpView.Render(w, r, nil)
}

// Register handles the /signup POST
func (u *User) Register(w http.ResponseWriter, r *http.Request) {
	form := signUpForm{}
	vd := views.Data{}
	ParseForm(r, &form)
	user := models.User{
		Introduction: form.Introduction,
	}
	if err := u.us.Create(&user); err != nil {
		vd.Alert = &views.Alert{
			Type:    "danger",
			Message: err.Error(),
		}
		u.SignUpView.Render(w, r, vd)
		return
	}
	chats := []models.Chat{
		{
			UserID:  user.ID,
			Content: user.Introduction,
			Role:    "user",
		},
	}
	_ = u.cs.Create(&models.Chat{
		UserID:  user.ID,
		Content: user.Introduction,
		Role:    "user",
	})
	_ = u.cs.Create(chatGPT(&chats))
	cookie := &http.Cookie{
		Name:  middleware.BrowserCookieName,
		Value: fmt.Sprint(user.ID),
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusFound)
}

// CookieTest is used to test the cookie
func (u *User) CookieTest(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(middleware.BrowserCookieName)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	user, err := u.us.ByID(cookie.Value)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	fmt.Fprintf(w, "%+v", user)
}

// Protected is a test route for middleware
func (u *User) Protected(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "User is logged in")
}
