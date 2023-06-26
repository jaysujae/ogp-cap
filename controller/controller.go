package controller

import (
	"fmt"
	"github.com/gorilla/mux"
	appcontext "hackathon/context"
	"hackathon/models"
	"hackathon/views"
	"net/http"
	"strconv"
)

// Controller defines the shape of the static struct
type Controller struct {
	HomeView *views.Views

	us          models.UserService
	chatService models.ChatService
	cs          models.CommentService
}

type commentForm struct {
	Content string `schema:"content"`
}

// New returns the static struct
func New(us models.UserService, ps models.ChatService, cs models.CommentService) *Controller {
	return &Controller{
		HomeView: views.NewView("bootstrap", "static/home", "user/signup"),

		us:          us,
		chatService: ps,
		cs:          cs,
	}
}

// Home handles the / GET
func (c *Controller) Home(w http.ResponseWriter, r *http.Request) {
	user := appcontext.GetUserFromContext(r)
	if user == nil {
		c.HomeView.Render(w, r, nil)
		return
	}
	userID := user.ID
	chats, err := c.chatService.FindByUserID(userID)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	companions, err := c.us.GetGroupUsersByID(userID)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	fmt.Println(user.Image)
	data := struct {
		Chats      *[]models.Chat
		Companions *[]models.User
		Current    *models.User
	}{
		Chats:      chats,
		Companions: companions,
		Current:    user,
	}
	c.HomeView.Render(w, r, data)
}

// UserPage shows user's chats
func (c *Controller) UserPage(w http.ResponseWriter, r *http.Request) {
	user := appcontext.GetUserFromContext(r)
	if user == nil {
		c.HomeView.Render(w, r, nil)
		return
	}
	userID := user.ID

	vars := mux.Vars(r)
	id := vars["id"]
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return
	}
	uID := uint(idInt)
	chats, err := c.chatService.FindByUserID(uID)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	companions, err := c.us.GetGroupUsersByID(userID)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	data := struct {
		Chats      *[]models.Chat
		Companions *[]models.User
		Current    *models.User
	}{
		Chats:      chats,
		Companions: companions,
		Current:    user,
	}
	c.HomeView.Render(w, r, data)
}

// Comment handles commenting on a chat
func (c *Controller) Comment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return
	}
	uID := uint(idInt)
	user := appcontext.GetUserFromContext(r)
	userID := user.ID
	form := &commentForm{}
	ParseForm(r, form)
	comment := &models.Comment{
		UserID:  userID,
		ChatID:  uID,
		Content: form.Content,
	}

	if err := c.cs.Create(comment); err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	referer := r.Header.Get("Referer")
	http.Redirect(w, r, referer, http.StatusFound)

}
