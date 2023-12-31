package frontend

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"net/http"
	eventModel "persistent-queue/api/event"
	"persistent-queue/api/eventbus"
	"persistent-queue/api/frontend"
	"time"
)

type Service struct {
	eventBusService eventbus.IService
}

func NewService(eventBusService eventbus.IService) *Service {
	return &Service{
		eventBusService: eventBusService,
	}
}

func (s *Service) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Println("received HealthCheck request")
	mp, err := s.eventBusService.CountEventsInQueue()
	if err != nil {
		log.Printf("failed to get count from queue %s", err.Error())
		sendJsonResponse(w, nil, http.StatusInternalServerError)
		return
	}
	sendJsonResponse(w, mp, http.StatusOK)
}

func (s *Service) PublishEvent(w http.ResponseWriter, r *http.Request) {
	fmt.Println("received PublishEvent request")
	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("failed to read req %s", err.Error())
		sendJsonResponse(w, nil, http.StatusInternalServerError)
		return
	}
	publishReq := &frontend.PublishEventRequest{}
	err = json.Unmarshal(req, publishReq)
	if err != nil {
		log.Printf("failed to unmarshal req %s", err.Error())
		sendJsonResponse(w, nil, http.StatusInternalServerError)
		return
	}
	err = s.eventBusService.EnqueueEvent(&eventModel.Event{
		EventId:     uuid.NewString(),
		UserID:      publishReq.UserID,
		Payload:     publishReq.Payload,
		PublishedAt: time.Now(),
	})
	if err != nil {
		log.Printf("failed to enqueue publishReq %s", err.Error())
		sendJsonResponse(w, nil, http.StatusInternalServerError)
		return
	}
	sendJsonResponse(w, "published event", http.StatusOK)
}

// sendJsonResponse updates the header with a httpStatusCode and response body
func sendJsonResponse(w http.ResponseWriter, resp interface{}, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	if resp != nil {
		data, err := json.Marshal(resp)
		if err != nil {
			return
		}
		w.Write(data)
	}
	return
}
