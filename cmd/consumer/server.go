package main

import (
	"github.com/go-co-op/gocron"
	"github.com/redis/go-redis/v9"
	"log"
	"persistent-queue/consumer"
	"persistent-queue/eventbus"
	eventsCache "persistent-queue/eventbus/dao/redis"
	"time"
)

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	consumerService := consumer.NewService(eventbus.NewService(eventsCache.NewEventsRedisCache(redisClient)))
	s := gocron.NewScheduler(time.UTC)

	// Every starts the job immediately and then runs at the
	// specified interval
	_, err := s.Every(5).Seconds().Do(func() {
		err := consumerService.PollEventsQueue()
		if err != nil {
			log.Printf("error performing polling on consumer, err: %s\n", err.Error())
			return
		}
	})
	if err != nil {
		log.Printf("error setting up job, err %s\n", err.Error())
	}
	s.StartBlocking()
}
