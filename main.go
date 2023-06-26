package main

import (
	"bufio"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"hackathon/controller"
	"hackathon/middleware"
	"hackathon/models"
	"net/http"
	"os"
)

const serverPort = "3000"

type Result struct {
	Body   string `json:"body"`
	Status int    `json:"status"`
}

func getStdinLine() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return line, nil
}
func main() {
	godotenv.Load()
	psqlInfo := os.Getenv("DATABASE_URL")

	svc, err := models.NewServices(psqlInfo)
	requireUserMW := middleware.NewRequireUser(svc.User)
	userMW := middleware.NewUser(svc.User)
	if err != nil {
		panic(err)
	}
	defer svc.Close()
	svc.AutoMigrate()
	// reset db
	svc.DestroyAndCreate()

	defaultController := controller.New(svc.User, svc.Chat, svc.Comment)
	userC := controller.NewUser(svc.User, svc.Chat)
	postC := controller.NewPost(svc.Chat)

	r := mux.NewRouter()
	r.HandleFunc("/", defaultController.Home)
	r.HandleFunc("/group", requireUserMW.RequireUserMiddleWare(defaultController.Group)).Methods("GET")

	r.HandleFunc("/signup", userC.Register).Methods("POST")
	r.HandleFunc("/logout", userC.LogOut).Methods("GET")

	r.HandleFunc("/post", requireUserMW.RequireUserMiddleWare(postC.HandlePost)).Methods("POST")
	r.HandleFunc("/post", requireUserMW.RequireUserMiddleWare(postC.HandlePost)).Methods("POST")
	r.HandleFunc("/delete/{id}", requireUserMW.RequireUserMiddleWare(postC.HandleDelete)).Methods("POST")

	r.HandleFunc("/user/{id}", requireUserMW.RequireUserMiddleWare(postC.ListPage)).Methods("GET")
	r.HandleFunc("/comment/{id}", requireUserMW.RequireUserMiddleWare(defaultController.Comment)).Methods("POST")

	fmt.Printf("Listening at port %s", serverPort)
	_ = http.ListenAndServe(":"+serverPort, userMW.UserMiddleWareFn(r))
}
