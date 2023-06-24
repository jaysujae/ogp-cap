package controller

import (
	"github.com/gorilla/mux"
	appcontext "hackathon/context"
	"hackathon/models"
	"hackathon/views"
	"net/http"
	"strconv"
)

// Controller defines the shape of the static struct
type Controller struct {
	HomeView  *views.Views
	GroupView *views.Views

	us models.UserService
	ps models.ChatService
	cs models.CommentService
}

type commentForm struct {
	Content string `schema:"content"`
}

// New returns the static struct
func New(us models.UserService, ps models.ChatService, cs models.CommentService) *Controller {
	return &Controller{
		HomeView:  views.NewView("bootstrap", "static/home", "user/signup"),
		GroupView: views.NewView("bootstrap", "user/group"),

		us: us,
		ps: ps,
		cs: cs,
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
	posts, err := c.ps.FindByUserID(userID)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	c.HomeView.Render(w, r, posts)
}

// Group shows a group of people similar to the user
func (c *Controller) Group(w http.ResponseWriter, r *http.Request) {
	userID := appcontext.GetUserFromContext(r).ID
	users, err := c.us.GetGroupUsersByID(userID)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	c.GroupView.Render(w, r, users)
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
