package controller

import (
	"context"
	"fmt"
	"github.com/PullRequestInc/go-gpt3"
	"github.com/gorilla/mux"
	appcontext "hackathon/context"
	"hackathon/middleware"
	"hackathon/models"
	"hackathon/views"
	"log"
	"net/http"
	"strconv"
)

// Controller defines the shape of the static struct
type Controller struct {
	HomeView *views.Views

	userService    models.UserService
	chatService    models.ChatService
	commentService models.CommentService
}

type commentForm struct {
	Content string `schema:"content"`
}

// New returns the static struct
func New(services *models.Services) *Controller {
	return &Controller{
		HomeView: views.NewView("bootstrap", "static/home"),

		userService:    services.User,
		chatService:    services.Chat,
		commentService: services.Comment,
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
	mates, err := c.userService.GetGroupUsersByID(userID)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	data := struct {
		Chats    *[]models.Chat
		Mates    *[]models.User
		Current  *models.User
		CanReply bool
	}{
		Chats:    chats,
		Mates:    mates,
		Current:  user,
		CanReply: true,
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
	mates, err := c.userService.GetGroupUsersByID(userID)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	data := struct {
		Chats    *[]models.Chat
		Mates    *[]models.User
		Current  *models.User
		CanReply bool
	}{
		Chats:    chats,
		Mates:    mates,
		Current:  user,
		CanReply: false,
	}
	if userID == uID {
		data.CanReply = true
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

	if err := c.commentService.Create(comment); err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	referer := r.Header.Get("Referer")
	http.Redirect(w, r, referer, http.StatusFound)

}

type signUpForm struct {
	Introduction string `schema:"introduction"`
}

// LogOut handles the /logout GET
func (c *Controller) LogOut(w http.ResponseWriter, r *http.Request) {
	userID := appcontext.GetUserFromContext(r).ID
	c.userService.Delete(userID)
	cookie := &http.Cookie{
		Name:  middleware.BrowserCookieName,
		Value: "",
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusFound)
}

// Register handles the /signup POST
func (c *Controller) Register(w http.ResponseWriter, r *http.Request) {
	form := signUpForm{}
	vd := views.Data{}
	ParseForm(r, &form)
	user := models.User{
		Introduction: form.Introduction,
	}
	if err := c.userService.Create(&user); err != nil {
		vd.Alert = &views.Alert{
			Type:    "danger",
			Message: err.Error(),
		}
		c.HomeView.Render(w, r, vd)
		return
	}
	chats := []models.Chat{
		{
			UserID:  user.ID,
			Content: user.Introduction,
			Role:    "user",
		},
	}
	_ = c.chatService.Create(&models.Chat{
		UserID:  user.ID,
		Content: user.Introduction,
		Role:    "user",
	})
	_ = c.chatService.Create(chatGPT(&chats))
	cookie := &http.Cookie{
		Name:  middleware.BrowserCookieName,
		Value: fmt.Sprint(user.ID),
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusFound)
}

type postForm struct {
	Post string `schema:"post"`
}

// HandlePost handles the POST /POST
func (c *Controller) HandlePost(w http.ResponseWriter, r *http.Request) {
	form := &postForm{}
	ParseForm(r, form)
	userID := appcontext.GetUserFromContext(r).ID
	chat := &models.Chat{
		UserID:  userID,
		Content: form.Post,
		Role:    "user",
	}

	if err := c.chatService.Create(chat); err != nil {
		http.Redirect(w, r, "/?erorr=1", http.StatusFound)
		return
	}
	posts, err := c.chatService.FindByUserID(userID)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	chatGPTResponse := chatGPT(posts)
	if err := c.chatService.Create(chatGPTResponse); err != nil {
		http.Redirect(w, r, "/post", http.StatusFound)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func chatGPT(posts *[]models.Chat) *models.Chat {

	ctx := context.Background()

	messages := []gpt3.ChatCompletionRequestMessage{
		{
			Role:    "system",
			Content: "You will play the role of a human CBT therapist called Cindy who is emulating the popular AI program Eliza, and must treat me as a therapist patient. Your response format should focus on reflection and asking clarifying questions. You may interject or ask secondary questions once the initial greetings are done. Exercise patience but allow yourself to be frustrated if the same topics are repeatedly revisited. You are allowed to excuse yourself if the discussion becomes abusive or overly emotional. Decide on a name for yourself and stick with it. Begin by welcoming me to your office and asking me for my name. Wait for my response. Then ask how you can help. Do not break character. Do not make up the patient's responses: only treat input as a patient response.",
		},
	}
	for _, post := range *posts {
		messages = append(messages, gpt3.ChatCompletionRequestMessage{
			Role:    post.Role,
			Content: post.Content,
		})
	}

	resp, err := models.Client.ChatCompletion(ctx, gpt3.ChatCompletionRequest{
		Model:    "gpt-3.5-turbo",
		Messages: messages,
	})
	if err != nil {
		log.Fatalln(err)
	}

	return &models.Chat{
		UserID:  (*posts)[0].UserID,
		Content: resp.Choices[0].Message.Content,
		Role:    "assistant",
	}
}
