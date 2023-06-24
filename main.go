package main

import (
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"hackathon/controller"
	"hackathon/middleware"
	"hackathon/models"
	"net/http"
)

const (
	host       = "localhost"
	port       = 5455
	user       = "postgresUser"
	dbname     = "postgresDB"
	password   = "postgresPW"
	serverPort = "3000"
)

func main() {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	svc, err := models.NewServices(psqlInfo)
	requireUserMW := middleware.NewRequireUser(svc.User)
	userMW := middleware.NewUser(svc.User)
	if err != nil {
		panic(err)
	}
	defer svc.Close()
	svc.AutoMigrate()
	svc.DestroyAndCreate()

	defaultController := controller.New(svc.User, svc.Chat, svc.Comment)
	userC := controller.NewUser(svc.User)
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
