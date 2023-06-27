package main

import (
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

	defaultController := controller.New(svc)

	r := mux.NewRouter()
	r.HandleFunc("/", defaultController.Home)
	r.HandleFunc("/user/{id}", requireUserMW.RequireUserMiddleWare(defaultController.UserPage)).Methods("GET")
	r.HandleFunc("/comment/{id}", requireUserMW.RequireUserMiddleWare(defaultController.Comment)).Methods("POST")

	r.HandleFunc("/signup", defaultController.Register).Methods("POST")
	r.HandleFunc("/logout", defaultController.LogOut).Methods("GET")

	r.HandleFunc("/post", requireUserMW.RequireUserMiddleWare(defaultController.HandlePost)).Methods("POST")

	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))
	fmt.Printf("Listening at port %s", serverPort)
	_ = http.ListenAndServe(":"+serverPort, userMW.UserMiddleWareFn(r))
}
