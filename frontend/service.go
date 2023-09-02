package frontend

import (
	"fmt"
	"net/http"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Println("received HealthCheck request")
	_, err := fmt.Fprintf(w, "health persistent queue service")
	if err != nil {
		return
	}
}

func PublishEvent(w http.ResponseWriter, r *http.Request) {
	fmt.Println("received PublishEvent request")
	_, err := fmt.Fprintf(w, "Received Your Request")
	if err != nil {
		return
	}
}
