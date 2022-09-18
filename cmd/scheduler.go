package main

import (
	"log"
	"net/http"

	config "github.com/AleksandrCherepanov/go_telegram/pkg/telegram/config"
	"github.com/AleksandrCherepanov/tg-scheduler/internal/middleware"
	"github.com/AleksandrCherepanov/tg-scheduler/internal/scheduler"
	"github.com/AleksandrCherepanov/tg-scheduler/internal/server"
	"github.com/gorilla/mux"
)

func main() {
	_, err := config.GetConfig()
	if err != nil {
		log.Fatalf("Can't intitalize config. %v", err.Error())
	}

	server := server.NewRouter()
	router := mux.NewRouter()

	router.HandleFunc("/schedule", server.Resolve).Methods("POST", "GET")

	loggedRouter := middleware.Logging(router)
	panicRecoveryRouter := middleware.PanicRecovery(loggedRouter)
	
	go scheduler.Notify()
	log.Fatal(http.ListenAndServe(":4000", panicRecoveryRouter))
}
