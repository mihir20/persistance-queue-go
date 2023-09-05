package main

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"os"
	"persistent-queue/eventbus"
	eventsCache "persistent-queue/eventbus/dao/redis"
	"persistent-queue/pkg/retrystrategy"
	"persistent-queue/pkg/subscriber"
	taskqueueNs "persistent-queue/pkg/taskqueue"
	"persistent-queue/snowflakeconsumer"
	"time"
)

func main() {
	redisHost := ""
	port := ""
	if val, ok := os.LookupEnv("DB_HOST"); ok {
		redisHost = val
	}
	if val, ok := os.LookupEnv("DB_PORT"); ok {
		port = val
	}
	// todo: add wire impl
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	eventbusService := eventbus.NewService(eventsCache.NewEventsRedisCache(redisClient))
	consumerService := snowflakeconsumer.NewService(eventbusService)
	retryStrategy := retrystrategy.NewExponentialBackOffRetryStrategy(3*time.Second, 3)
	newSubscriber := subscriber.NewSubscriber(2, taskqueueNs.SnowflakeConsumerTaskQueue,
		eventbusService, retryStrategy, consumerService.ConsumeEvent)
	newSubscriber.StartWorker()
}
