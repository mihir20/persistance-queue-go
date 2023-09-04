package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"persistent-queue/api/eventbus"
	"persistent-queue/api/taskqueue"
)

type EventsRedisCache struct {
	redisClient *redis.Client
}

func NewEventsRedisCache(redisClient *redis.Client) *EventsRedisCache {
	return &EventsRedisCache{redisClient: redisClient}
}

func (c *EventsRedisCache) CreateEvent(passenger *eventbus.PassengerEvent, taskQueues []taskqueue.TaskQueue) error {
	ctx := context.Background()
	bytes, err := json.Marshal(passenger)
	if err != nil {
		return err
	}
	pipeline := c.redisClient.Pipeline()
	for _, queue := range taskQueues {
		pipeline.ZAdd(ctx, string(queue), redis.Z{
			Score:  float64(passenger.EventTime.Unix()),
			Member: bytes,
		})
	}
	_, err = pipeline.Exec(ctx)
	if err != nil {
		return fmt.Errorf("error inserting multiple keys, err%w", err)
	}
	return nil
}

func (c *EventsRedisCache) GetEvent(taskQueue taskqueue.TaskQueue) (*eventbus.PassengerEvent, error) {
	eventStr, err := c.redisClient.ZRange(context.Background(), string(taskQueue), 0, 0).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("error while fetching event from queue, %w", err)
	}
	if errors.Is(err, redis.Nil) || len(eventStr) == 0 {
		return nil, nil
	}
	event := &eventbus.PassengerEvent{}
	fmt.Println(eventStr)
	err = json.Unmarshal([]byte(eventStr[0]), event)
	if err != nil {
		return nil, fmt.Errorf("error while unmarshalling, eventStr: %s ,err: %w", eventStr, err)
	}
	return event, nil
}

func (c *EventsRedisCache) UpdateEvent(taskQueue taskqueue.TaskQueue, event *eventbus.PassengerEvent) error {
	bytes, err := json.Marshal(event)
	if err != nil {
		return err
	}
	err = c.redisClient.RPush(context.Background(), string(taskQueue), string(bytes)).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *EventsRedisCache) DeleteEvent(event *eventbus.PassengerEvent) error {
	return nil
}
