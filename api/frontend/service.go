package frontend

import "net/http"

type IService interface {
	// PublishEvent publish a new event in the queue
	PublishEvent(w http.ResponseWriter, r *http.Request)
	// HealthCheck returns health information of a queue
	HealthCheck(w http.ResponseWriter, r *http.Request)
}

type PublishEventRequest struct {
	UserID  string `json:"userid"`
	Payload string `json:"payload"`
}
