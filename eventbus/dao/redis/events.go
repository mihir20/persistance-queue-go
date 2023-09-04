package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"persistent-queue/api/eventbus"
	"persistent-queue/api/taskqueue"
	"time"
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

func (c *EventsRedisCache) GetEvent(taskQueue taskqueue.TaskQueue) (*eventbus.PassengerEvent, int64, error) {
	eventStr, err := c.redisClient.ZRangeWithScores(context.Background(), string(taskQueue), 0, 0).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, 0, fmt.Errorf("error while fetching event from queue, %w", err)
	}
	if errors.Is(err, redis.Nil) || len(eventStr) == 0 {
		return nil, 0, nil
	}
	event := &eventbus.PassengerEvent{}
	err = json.Unmarshal([]byte(eventStr[0].Member.(string)), event)
	if err != nil {
		return nil, 0, fmt.Errorf("error while unmarshalling, eventStr: %v ,err: %w", eventStr, err)
	}
	return event, int64(eventStr[0].Score), nil
}

func (c *EventsRedisCache) UpdateEvent(taskQueue taskqueue.TaskQueue, oldPassenger, updatedPassengerEvent *eventbus.PassengerEvent, nextExecutionTime time.Time) error {
	ctx := context.Background()
	oldBytes, err := json.Marshal(oldPassenger)
	if err != nil {
		return err
	}
	newBytes, err := json.Marshal(updatedPassengerEvent)
	pipeline := c.redisClient.Pipeline()
	pipeline.ZRem(ctx, string(taskQueue), oldBytes)
	pipeline.ZAdd(ctx, string(taskQueue), redis.Z{
		Score:  float64(nextExecutionTime.Unix()),
		Member: newBytes,
	})
	_, err = pipeline.Exec(ctx)
	if err != nil {
		return fmt.Errorf("error inserting multiple keys, err%w", err)
	}
	return nil
}

func (c *EventsRedisCache) DeleteEvent(taskQueue taskqueue.TaskQueue, passenger *eventbus.PassengerEvent) error {
	ctx := context.Background()
	bytes, err := json.Marshal(passenger)
	if err != nil {
		return err
	}
	err = c.redisClient.ZRem(ctx, string(taskQueue), bytes).Err()
	if err != nil {
		return err
	}
	return nil
}
