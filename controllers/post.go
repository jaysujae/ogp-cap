package controllers

import (
	"context"
	"github.com/PullRequestInc/go-gpt3"
	"github.com/joho/godotenv"
	appcontext "hackathon/context"
	"hackathon/models"
	"hackathon/views"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

// Post defines the shape of the struct
type Post struct {
	ps       models.PostService
	NewView  *views.Views
	ListView *views.Views
}

type postForm struct {
	Post string `schema:"post"`
}

// NewPost returns the post struct
func NewPost(ps models.PostService) *Post {
	return &Post{
		ps:       ps,
		NewView:  views.NewView("bootstrap", "post/new"),
		ListView: views.NewView("bootstrap", "post/list"),
	}
}

// PostPage responds to the GET /POST route
func (p *Post) PostPage(w http.ResponseWriter, r *http.Request) {
	p.NewView.Render(w, r, nil)
}

// ListPage list all users posts
func (p *Post) ListPage(w http.ResponseWriter, r *http.Request) {
	userID := appcontext.GetUserFromContext(r).ID
	posts, err := p.ps.FindByUserID(userID)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	p.ListView.Render(w, r, posts)
}

// HandlePost handles the POST /POST
func (p *Post) HandlePost(w http.ResponseWriter, r *http.Request) {
	form := &postForm{}
	ParseForm(r, form)
	userID := appcontext.GetUserFromContext(r).ID
	post := &models.Post{
		UserID: userID,
		Post:   form.Post,
	}
	posts, err := p.ps.FindByUserID(userID)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	if err := p.ps.Create(post); err != nil {
		http.Redirect(w, r, "/post", http.StatusFound)
		return
	}

	chatGPTResponse := chatGPT(posts)
	if err := p.ps.Create(chatGPTResponse); err != nil {
		http.Redirect(w, r, "/post", http.StatusFound)
		return
	}
	http.Redirect(w, r, "/list", http.StatusFound)
}

func chatGPT(posts *[]models.Post) *models.Post {
	godotenv.Load()

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatalln("Missing API KEY")
	}

	ctx := context.Background()
	client := gpt3.NewClient(apiKey)

	messages := []gpt3.ChatCompletionRequestMessage{
		{
			Role:    "system",
			Content: "You will play the role of a human CBT therapist called Cindy who is emulating the popular AI program Eliza, and must treat me as a therapist patient. Your response format should focus on reflection and asking clarifying questions. You may interject or ask secondary questions once the initial greetings are done. Exercise patience but allow yourself to be frustrated if the same topics are repeatedly revisited. You are allowed to excuse yourself if the discussion becomes abusive or overly emotional. Decide on a name for yourself and stick with it. Begin by welcoming me to your office and asking me for my name. Wait for my response. Then ask how you can help. Do not break character. Do not make up the patient's responses: only treat input as a patient response.",
		},
	}
	for _, post := range *posts {
		messages = append(messages, gpt3.ChatCompletionRequestMessage{
			Role:    "user",
			Content: post.Post,
		})
	}

	resp, err := client.ChatCompletion(ctx, gpt3.ChatCompletionRequest{
		Model:    "gpt-3.5-turbo",
		Messages: messages,
	})
	if err != nil {
		log.Fatalln(err)
	}

	return &models.Post{
		UserID: 1,
		Post:   resp.Choices[0].Message.Content,
	}
}

// HandleDelete deletes a post from the database
func (p *Post) HandleDelete(w http.ResponseWriter, r *http.Request) {
	pd := &views.Data{}
	vars := mux.Vars(r)
	id := vars["id"]
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return
	}
	uID := uint(idInt)
	user := appcontext.GetUserFromContext(r)
	post, err := p.ps.FindPostByID(uID)
	if err != nil {
		pd.Alert = &views.Alert{
			Type:    "danger",
			Message: err.Error(),
		}
		p.ListView.Render(w, r, pd)
		return
	}
	if post.UserID != user.ID {
		pd.Alert = &views.Alert{
			Type:    "danger",
			Message: models.ErrPostMissing.Error(),
		}
		p.ListView.Render(w, r, pd)
		return
	}
	if err := p.ps.Delete(post); err != nil {
		pd.Alert = &views.Alert{
			Type:    "danger",
			Message: err.Error(),
		}
		return
	}

	http.Redirect(w, r, "/list", http.StatusFound)
}
