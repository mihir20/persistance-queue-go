package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
	"net/http"
	"os"
	"persistent-queue/eventbus"
	eventsCache "persistent-queue/eventbus/dao/redis"
	"persistent-queue/frontend"
)

func main() {
	defer func() {
		fmt.Println("Shutting server down...")
	}()
	router := mux.NewRouter()
	http.Handle("/", router)

	redisHost := ""
	port := ""
	if val, ok := os.LookupEnv("DB_HOST"); ok {
		redisHost = val
	}
	if val, ok := os.LookupEnv("DB_PORT"); ok {
		port = val
	}
	// todo: add wire impl
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	frontendService := frontend.NewService(eventbus.NewService(eventsCache.NewEventsRedisCache(redisClient)))

	router.HandleFunc("/publish", frontendService.PublishEvent).Methods("POST")

	router.HandleFunc("/healthcheck", frontendService.HealthCheck).Methods("GET")

	fmt.Println("starting server on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("error starting the server, Err:%s", err.Error())
		return
	}
}
