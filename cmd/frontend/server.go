package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"persistent-queue/eventbus"
	"persistent-queue/frontend"
)

func main() {
	defer func() {
		fmt.Println("Shutting server down...")
	}()
	router := mux.NewRouter()
	http.Handle("/", router)

	// todo: add wire impl
	frontendService := frontend.NewService(eventbus.NewService(nil))

	router.HandleFunc("/publish", frontendService.PublishEvent).Methods("POST")

	router.HandleFunc("/healthcheck", frontendService.HealthCheck).Methods("GET")

	fmt.Println("starting server on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("error starting the server, Err:%s", err.Error())
		return
	}
}
