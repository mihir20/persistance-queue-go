package main

import (
	"github.com/redis/go-redis/v9"
	"persistent-queue/eventbus"
	eventsCache "persistent-queue/eventbus/dao/redis"
	"persistent-queue/fileconsumer"
	"persistent-queue/pkg/retrystrategy"
	"persistent-queue/pkg/subscriber"
	taskqueueNs "persistent-queue/pkg/taskqueue"
	"time"
)

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	eventbusService := eventbus.NewService(eventsCache.NewEventsRedisCache(redisClient))
	consumerService := fileconsumer.NewService()
	retryStrategy := retrystrategy.NewExponentialBackOffRetryStrategy(3*time.Second, 3)
	newSubscriber := subscriber.NewSubscriber(2, taskqueueNs.FileConsumerTaskQueue,
		eventbusService, retryStrategy, consumerService.ConsumeEvent)
	newSubscriber.StartWorker()
}
