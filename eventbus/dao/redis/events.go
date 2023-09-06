package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"persistent-queue/api/eventbus"
	"persistent-queue/api/taskqueue"
	"strconv"
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

func (c *EventsRedisCache) GetEvents(taskQueue taskqueue.TaskQueue, cutOffTime time.Time, countOfEvents int64) ([]*eventbus.PassengerEvent, error) {
	eventStr, err := c.redisClient.ZRangeByScoreWithScores(context.Background(), string(taskQueue), &redis.ZRangeBy{
		Min:   "0",
		Max:   strconv.FormatInt(cutOffTime.Unix(), 10),
		Count: countOfEvents,
	}).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("error while fetching event from queue, %w", err)
	}
	if errors.Is(err, redis.Nil) || len(eventStr) == 0 {
		return nil, nil
	}
	events := make([]*eventbus.PassengerEvent, len(eventStr))
	for i, element := range eventStr {
		event := &eventbus.PassengerEvent{}
		err = json.Unmarshal([]byte(element.Member.(string)), event)
		if err != nil {
			return nil, fmt.Errorf("error while unmarshalling, eventStr: %v ,err: %w", eventStr, err)
		}
		events[i] = event
	}
	return events, nil
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

func (c *EventsRedisCache) CountEventsInQueue(taskQueues []taskqueue.TaskQueue) (map[taskqueue.TaskQueue]int64, error) {
	ctx := context.Background()
	mp := make(map[taskqueue.TaskQueue]int64)
	for _, queue := range taskQueues {
		count, err := c.redisClient.ZCard(ctx, string(queue)).Result()
		if err != nil {
			return nil, fmt.Errorf("error getting count for key:%s, err%w", queue, err)
		}
		mp[queue] = count
	}
	return mp, nil
}
