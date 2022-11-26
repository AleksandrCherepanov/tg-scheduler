package main

import (
	"log"
	"net/http"

	"github.com/AleksandrCherepanov/tg-scheduler/internal/api"
	"github.com/AleksandrCherepanov/tg-scheduler/internal/middleware"
	"github.com/AleksandrCherepanov/tg-scheduler/internal/scheduler"
	"github.com/AleksandrCherepanov/tg-scheduler/internal/server"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Can't intitalize config. %v", err.Error())
	}

	server := server.NewRouter()
	router := mux.NewRouter()

	router.HandleFunc("/schedule", server.Resolve).Methods("POST", "GET")
	api.RegisterRoutes(router)

	loggedRouter := middleware.Logging(router)
	basicAuthRouter := middleware.BasicAuth(loggedRouter)
	panicRecoveryRouter := middleware.PanicRecovery(basicAuthRouter)

	notificator := scheduler.GetNotificator()
	go notificator.StartNotification()
	log.Println(http.ListenAndServe(":4000", panicRecoveryRouter))
}
