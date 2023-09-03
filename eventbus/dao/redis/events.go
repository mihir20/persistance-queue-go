package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	eventModel "persistent-queue/api/event"
	"persistent-queue/api/taskqueue"
	"time"
)

type EventsRedisCache struct {
	redisClient *redis.Client
}

func NewEventsRedisCache(redisClient *redis.Client) *EventsRedisCache {
	return &EventsRedisCache{redisClient: redisClient}
}

func (c *EventsRedisCache) CreateEvent(event *eventModel.Event, taskQueues []taskqueue.TaskQueue) error {
	ctx := context.Background()
	bytes, err := json.Marshal(event)
	if err != nil {
		return err
	}
	pipeline := c.redisClient.Pipeline()
	for _, queue := range taskQueues {
		pipeline.LPush(ctx, string(queue), string(bytes))
	}
	_, err = pipeline.Exec(ctx)
	if err != nil {
		return fmt.Errorf("error inserting multiple keys, err%w", err)
	}
	return nil
}

func (c *EventsRedisCache) GetEvent(taskQueue taskqueue.TaskQueue) (*eventModel.Event, error) {
	eventStr, err := c.redisClient.RPop(context.Background(), string(taskQueue)).Result()
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

func (c *EventsRedisCache) UpdateEvent(taskQueue taskqueue.TaskQueue, event *eventModel.Event) error {
	bytes, err := json.Marshal(event)
	if err != nil {
		return err
	}
	err = c.redisClient.RPush(context.Background(), string(taskQueue), string(bytes), 5*time.Minute).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *EventsRedisCache) DeleteEvent(event *eventModel.Event) error {
	return nil
}
