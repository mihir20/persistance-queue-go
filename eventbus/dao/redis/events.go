package redis

import (
	"context"
	"encoding/json"
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
	err = c.redisClient.LPush(context.Background(), "queue1", string(bytes), 5*time.Minute).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *EventsRedisCache) GetEvent(eventName string) (*eventModel.Event, error) {
	eventStr, err := c.redisClient.RPop(context.Background(), "queue1").Result()
	if err != nil {
		return nil, err
	}
	event := &eventModel.Event{}
	err = json.Unmarshal([]byte(eventStr), event)
	if err != nil {
		return nil, err
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
