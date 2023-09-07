package redis

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"persistent-queue/api/event"
	"persistent-queue/api/eventbus"
	"persistent-queue/api/taskqueue"
	taskqueueNs "persistent-queue/pkg/taskqueue"
	"reflect"
	"testing"
	"time"
)

var (
	redisConfig = &redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	}

	samplePassengerEvent1 = &eventbus.PassengerEvent{
		Event:         &event.Event{EventId: "id-1"},
		RetryAttempts: 1,
		EventTime:     time.Unix(100000, 0),
	}
	samplePassengerEvent2 = &eventbus.PassengerEvent{
		Event:         &event.Event{EventId: "id-2"},
		RetryAttempts: 1,
		EventTime:     time.Unix(100002, 0),
	}
	samplePassengerEvent3 = &eventbus.PassengerEvent{
		Event:         &event.Event{EventId: "id-3"},
		RetryAttempts: 1,
		EventTime:     time.Unix(100003, 0),
	}
	samplePassengerEvent4 = &eventbus.PassengerEvent{
		Event:         &event.Event{EventId: "id-4"},
		RetryAttempts: 1,
		EventTime:     time.Unix(100004, 0),
	}
)

func closeRedisClient(redisClient *redis.Client) {
	_ = redisClient.Close()
}

// to flush all the keys
func clearCache(r *redis.Client) {
	r.FlushDB(context.Background())
}

// populate the cache with keys
func populateCache(r *redis.Client) error {
	bytes1, _ := json.Marshal(samplePassengerEvent1)
	bytes2, _ := json.Marshal(samplePassengerEvent2)
	bytes3, _ := json.Marshal(samplePassengerEvent3)
	bytes4, _ := json.Marshal(samplePassengerEvent4)
	err := r.ZAdd(context.Background(), string(taskqueueNs.SnowflakeConsumerTaskQueue), redis.Z{
		Score:  100000,
		Member: bytes1,
	}, redis.Z{
		Score:  100002,
		Member: bytes2,
	}, redis.Z{
		Score:  100003,
		Member: bytes3,
	}, redis.Z{
		Score:  100004,
		Member: bytes4,
	}).Err()
	if err != nil {
		return err
	}
	err = r.ZAdd(context.Background(), string(taskqueueNs.VendorApiConsumerTaskQueue), redis.Z{
		Score:  100000,
		Member: bytes1,
	}, redis.Z{
		Score:  100002,
		Member: bytes2,
	}, redis.Z{
		Score:  100003,
		Member: bytes3,
	}).Err()
	if err != nil {
		return err
	}
	return nil
}

func TestEventsRedisCache_CountEventsInQueue(t *testing.T) {
	redisClient := redis.NewClient(redisConfig)
	defer closeRedisClient(redisClient)
	type args struct {
		taskQueues []taskqueue.TaskQueue
	}
	tests := []struct {
		name    string
		args    args
		want    map[taskqueue.TaskQueue]int64
		wantErr bool
	}{
		{
			name: "successful fetch",
			args: args{
				taskQueues: []taskqueue.TaskQueue{
					taskqueueNs.SnowflakeConsumerTaskQueue,
					taskqueueNs.VendorApiConsumerTaskQueue,
				},
			},
			want: map[taskqueue.TaskQueue]int64{
				taskqueueNs.SnowflakeConsumerTaskQueue: 4,
				taskqueueNs.VendorApiConsumerTaskQueue: 3,
			},
			wantErr: false,
		},
		{
			name: "successful fetch 3 task queues",
			args: args{
				taskQueues: []taskqueue.TaskQueue{
					taskqueueNs.SnowflakeConsumerTaskQueue,
					taskqueueNs.VendorApiConsumerTaskQueue,
					taskqueueNs.FileConsumerTaskQueue,
				},
			},
			want: map[taskqueue.TaskQueue]int64{
				taskqueueNs.SnowflakeConsumerTaskQueue: 4,
				taskqueueNs.VendorApiConsumerTaskQueue: 3,
				taskqueueNs.FileConsumerTaskQueue:      0,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearCache(redisClient)
			populateErr := populateCache(redisClient)
			if populateErr != nil {
				t.Errorf("error while populating cache , %v", populateErr)
				return
			}
			c := &EventsRedisCache{
				redisClient: redisClient,
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
	redisClient := redis.NewClient(redisConfig)
	defer closeRedisClient(redisClient)
	type args struct {
		passenger  *eventbus.PassengerEvent
		taskQueues []taskqueue.TaskQueue
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "create successful",
			args: args{
				passenger: &eventbus.PassengerEvent{
					Event:         &event.Event{EventId: "new-event-id"},
					RetryAttempts: 1,
					EventTime:     time.Unix(100000, 0),
				},
				taskQueues: []taskqueue.TaskQueue{
					taskqueueNs.SnowflakeConsumerTaskQueue,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearCache(redisClient)
			populateErr := populateCache(redisClient)
			if populateErr != nil {
				t.Errorf("error while populating cache , %v", populateErr)
				return
			}
			c := &EventsRedisCache{
				redisClient: redisClient,
			}
			err := c.CreateEvent(tt.args.passenger, tt.args.taskQueues)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				return
			}
			for _, taskQueue := range tt.args.taskQueues {
				resArr, err := redisClient.ZRangeArgs(context.Background(), redis.ZRangeArgs{
					Key:     string(taskQueue),
					Start:   tt.args.passenger.EventTime.UnixNano(),
					Stop:    tt.args.passenger.EventTime.UnixNano(),
					ByScore: true,
				}).Result()
				if err != nil {
					t.Errorf("error while fetching all records, error = %v", err)
				}
				notFound := true
				for _, s := range resArr {
					e := &eventbus.PassengerEvent{}
					err = json.Unmarshal([]byte(s), e)
					if e.Event.EventId == tt.args.passenger.Event.EventId {
						notFound = false
					}
				}
				if notFound {
					t.Errorf("UpdateEvent() failed as different event id is present at given score")
				}
			}
		})
	}
}

func TestEventsRedisCache_DeleteEvent(t *testing.T) {
	redisClient := redis.NewClient(redisConfig)
	defer closeRedisClient(redisClient)
	type fields struct {
		redisClient *redis.Client
	}
	type args struct {
		taskQueue taskqueue.TaskQueue
		passenger *eventbus.PassengerEvent
	}
	tests := []struct {
		name    string
		args    args
		fields  fields
		wantErr bool
	}{
		{
			name:   "delete successful",
			fields: fields{redisClient: redisClient},
			args: args{
				taskQueue: taskqueueNs.SnowflakeConsumerTaskQueue,
				passenger: samplePassengerEvent1,
			},
			wantErr: false,
		},
		{
			name:   "delete failed",
			fields: fields{redisClient: redis.NewClient(&redis.Options{Addr: "dffa"})},
			args: args{
				taskQueue: taskqueueNs.FileConsumerTaskQueue,
				passenger: samplePassengerEvent1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearCache(redisClient)
			populateErr := populateCache(redisClient)
			if populateErr != nil {
				t.Errorf("error while populating cache , %v", populateErr)
				return
			}
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
	redisClient := redis.NewClient(redisConfig)
	defer closeRedisClient(redisClient)
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
		{
			name: "successful get all till cutoff",
			fields: fields{
				redisClient: redisClient,
			},
			args: args{
				taskQueue:     taskqueueNs.SnowflakeConsumerTaskQueue,
				cutOffTime:    time.Unix(0, 100003),
				countOfEvents: 10,
			},
			want:    []*eventbus.PassengerEvent{samplePassengerEvent1, samplePassengerEvent2, samplePassengerEvent3},
			wantErr: false,
		},
		{
			name: "successful get 2 till cutoff",
			fields: fields{
				redisClient: redisClient,
			},
			args: args{
				taskQueue:     taskqueueNs.SnowflakeConsumerTaskQueue,
				cutOffTime:    time.Unix(0, 100003),
				countOfEvents: 2,
			},
			want:    []*eventbus.PassengerEvent{samplePassengerEvent1, samplePassengerEvent2},
			wantErr: false,
		},
		{
			name: "successful get 5 till cutoff",
			fields: fields{
				redisClient: redisClient,
			},
			args: args{
				taskQueue:     taskqueueNs.SnowflakeConsumerTaskQueue,
				cutOffTime:    time.Unix(0, 100005),
				countOfEvents: 5,
			},
			want:    []*eventbus.PassengerEvent{samplePassengerEvent1, samplePassengerEvent2, samplePassengerEvent3, samplePassengerEvent4},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearCache(redisClient)
			populateErr := populateCache(redisClient)
			if populateErr != nil {
				t.Errorf("error while populating cache , %v", populateErr)
				return
			}
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
	redisClient := redis.NewClient(redisConfig)
	defer closeRedisClient(redisClient)
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
		{
			name: "successful update",
			fields: fields{
				redisClient: redisClient,
			},
			args: args{
				taskQueue:    taskqueueNs.SnowflakeConsumerTaskQueue,
				oldPassenger: samplePassengerEvent1,
				updatedPassengerEvent: func() *eventbus.PassengerEvent {
					return &eventbus.PassengerEvent{
						Event:         samplePassengerEvent1.Event,
						RetryAttempts: samplePassengerEvent1.RetryAttempts + 1,
						EventTime:     samplePassengerEvent1.EventTime,
					}
				}(),
				nextExecutionTime: time.Unix(0, 1000001),
			},
			wantErr: false,
		},
		{
			name: "failed update",
			fields: fields{
				redisClient: redis.NewClient(&redis.Options{Addr: "asdf"}),
			},
			args: args{
				taskQueue:             taskqueueNs.SnowflakeConsumerTaskQueue,
				oldPassenger:          samplePassengerEvent1,
				updatedPassengerEvent: samplePassengerEvent4,
				nextExecutionTime:     time.Unix(1000001, 0),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearCache(redisClient)
			populateErr := populateCache(redisClient)
			if populateErr != nil {
				t.Errorf("error while populating cache , %v", populateErr)
				return
			}
			c := &EventsRedisCache{
				redisClient: tt.fields.redisClient,
			}
			err := c.UpdateEvent(tt.args.taskQueue, tt.args.oldPassenger, tt.args.updatedPassengerEvent, tt.args.nextExecutionTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				return
			}
			resArr, err := redisClient.ZRangeArgs(context.Background(), redis.ZRangeArgs{
				Key:     string(tt.args.taskQueue),
				Start:   tt.args.nextExecutionTime.UnixNano(),
				Stop:    tt.args.nextExecutionTime.UnixNano(),
				ByScore: true,
			}).Result()
			if err != nil {
				t.Errorf("error while fetching all records, error = %v", err)
			}
			notFound := true
			for _, s := range resArr {
				e := &eventbus.PassengerEvent{}
				err = json.Unmarshal([]byte(s), e)
				if e.Event.EventId == tt.args.oldPassenger.Event.EventId {
					notFound = false
				}
			}
			if notFound {
				t.Errorf("UpdateEvent() failed as different event id is present at given score")
			}
		})
	}
}
