package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	eventModel "persistent-queue/api/event"
	"time"
)

type EventsRedisCache struct {
	redisClient *redis.Client
}

func NewEventsRedisCache(redisClient *redis.Client) *EventsRedisCache {
	return &EventsRedisCache{redisClient: redisClient}
}

func (c *EventsRedisCache) CreateEvent(event *eventModel.Event) error {
	bytes, err := json.Marshal(event)
	if err != nil {
		return err
	}
	err = c.redisClient.LPush(context.Background(), "queue1", string(bytes)).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *EventsRedisCache) GetEvent(eventName string) (*eventModel.Event, error) {
	eventStr, err := c.redisClient.RPop(context.Background(), "queue1").Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("error while fetching event from queue, %w", err)
	}
	if errors.Is(err, redis.Nil) || len(eventStr) == 0 {
		return nil, nil
	}
	event := &eventModel.Event{}
	fmt.Println(eventStr)
	err = json.Unmarshal([]byte(eventStr), event)
	if err != nil {
		return nil, fmt.Errorf("error while unmarshalling, eventStr: %s ,err: %w", eventStr, err)
	}
	return event, nil
}

func (c *EventsRedisCache) UpdateEvent(event *eventModel.Event) error {
	bytes, err := json.Marshal(event)
	if err != nil {
		return err
	}
	err = c.redisClient.RPush(context.Background(), "queue1", string(bytes), 5*time.Minute).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *EventsRedisCache) DeleteEvent(event *eventModel.Event) error {
	return nil
}
