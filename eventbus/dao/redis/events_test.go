package redis

import (
	"github.com/redis/go-redis/v9"
	"persistent-queue/api/eventbus"
	"persistent-queue/api/taskqueue"
	"reflect"
	"testing"
	"time"
)

func TestEventsRedisCache_CountEventsInQueue(t *testing.T) {
	type fields struct {
		redisClient *redis.Client
	}
	type args struct {
		taskQueues []taskqueue.TaskQueue
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[taskqueue.TaskQueue]int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &EventsRedisCache{
				redisClient: tt.fields.redisClient,
			}
			got, err := c.CountEventsInQueue(tt.args.taskQueues)
			if (err != nil) != tt.wantErr {
				t.Errorf("CountEventsInQueue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CountEventsInQueue() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEventsRedisCache_CreateEvent(t *testing.T) {
	type fields struct {
		redisClient *redis.Client
	}
	type args struct {
		passenger  *eventbus.PassengerEvent
		taskQueues []taskqueue.TaskQueue
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &EventsRedisCache{
				redisClient: tt.fields.redisClient,
			}
			if err := c.CreateEvent(tt.args.passenger, tt.args.taskQueues); (err != nil) != tt.wantErr {
				t.Errorf("CreateEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEventsRedisCache_DeleteEvent(t *testing.T) {
	type fields struct {
		redisClient *redis.Client
	}
	type args struct {
		taskQueue taskqueue.TaskQueue
		passenger *eventbus.PassengerEvent
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &EventsRedisCache{
				redisClient: tt.fields.redisClient,
			}
			if err := c.DeleteEvent(tt.args.taskQueue, tt.args.passenger); (err != nil) != tt.wantErr {
				t.Errorf("DeleteEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEventsRedisCache_GetEvents(t *testing.T) {
	type fields struct {
		redisClient *redis.Client
	}
	type args struct {
		taskQueue     taskqueue.TaskQueue
		cutOffTime    time.Time
		countOfEvents int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*eventbus.PassengerEvent
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &EventsRedisCache{
				redisClient: tt.fields.redisClient,
			}
			got, err := c.GetEvents(tt.args.taskQueue, tt.args.cutOffTime, tt.args.countOfEvents)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEvents() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetEvents() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEventsRedisCache_UpdateEvent(t *testing.T) {
	type fields struct {
		redisClient *redis.Client
	}
	type args struct {
		taskQueue             taskqueue.TaskQueue
		oldPassenger          *eventbus.PassengerEvent
		updatedPassengerEvent *eventbus.PassengerEvent
		nextExecutionTime     time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &EventsRedisCache{
				redisClient: tt.fields.redisClient,
			}
			if err := c.UpdateEvent(tt.args.taskQueue, tt.args.oldPassenger, tt.args.updatedPassengerEvent, tt.args.nextExecutionTime); (err != nil) != tt.wantErr {
				t.Errorf("UpdateEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
