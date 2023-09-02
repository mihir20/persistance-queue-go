package frontend

import "net/http"

type IService interface {
	PublishEvent(w http.ResponseWriter, r *http.Request)
	HealthCheck(w http.ResponseWriter, r *http.Request)
}

type PublishEventRequest struct {
	Name string `json:"name"`
}
