package eventbus

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"persistent-queue/api/event"
	"persistent-queue/api/eventbus"
	"persistent-queue/api/taskqueue"
	taskqueueNs "persistent-queue/pkg/taskqueue"
	"reflect"
	"testing"
	"time"
)

var (
	samplePassengerEvent = &eventbus.PassengerEvent{
		Event:         &event.Event{EventId: uuid.NewString()},
		RetryAttempts: 1,
		EventTime:     time.Now(),
	}
)

func TestService_CountEventsInQueue(t *testing.T) {
	ctr := gomock.NewController(t)
	tests := []struct {
		name      string
		setupMock func(dependencies *mockDependencies)
		want      map[taskqueue.TaskQueue]int64
		wantErr   bool
	}{
		{
			name: "success",
			setupMock: func(dependencies *mockDependencies) {
				dependencies.eventsDao.EXPECT().CountEventsInQueue([]taskqueue.TaskQueue{
					taskqueueNs.SnowflakeConsumerTaskQueue,
					taskqueueNs.FileConsumerTaskQueue,
					taskqueueNs.VendorApiConsumerTaskQueue,
				}).Return(map[taskqueue.TaskQueue]int64{
					taskqueueNs.SnowflakeConsumerTaskQueue: 4,
					taskqueueNs.FileConsumerTaskQueue:      5,
					taskqueueNs.VendorApiConsumerTaskQueue: 6,
				}, nil)
			},
			want: map[taskqueue.TaskQueue]int64{
				taskqueueNs.SnowflakeConsumerTaskQueue: 4,
				taskqueueNs.FileConsumerTaskQueue:      5,
				taskqueueNs.VendorApiConsumerTaskQueue: 6,
			},
			wantErr: false,
		},
		{
			name: "db call failed",
			setupMock: func(dependencies *mockDependencies) {
				dependencies.eventsDao.EXPECT().CountEventsInQueue([]taskqueue.TaskQueue{
					taskqueueNs.SnowflakeConsumerTaskQueue,
					taskqueueNs.FileConsumerTaskQueue,
					taskqueueNs.VendorApiConsumerTaskQueue,
				}).Return(nil, errors.New("failed to connect"))
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, md := newMockEventBusService(ctr)
			tt.setupMock(md)
			got, err := s.CountEventsInQueue()
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

func TestService_DequeueEventFromTaskQueue(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		taskQueue      taskqueue.TaskQueue
		passengerEvent *eventbus.PassengerEvent
	}
	tests := []struct {
		name      string
		args      args
		setupMock func(dependencies *mockDependencies)
		wantErr   bool
	}{
		{
			name: "success dequeue",
			args: args{
				taskQueue:      taskqueueNs.VendorApiConsumerTaskQueue,
				passengerEvent: samplePassengerEvent,
			},
			setupMock: func(dependencies *mockDependencies) {
				dependencies.eventsDao.EXPECT().
					DeleteEvent(taskqueueNs.VendorApiConsumerTaskQueue, samplePassengerEvent).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "failed to dequeue",
			args: args{
				taskQueue:      taskqueueNs.VendorApiConsumerTaskQueue,
				passengerEvent: samplePassengerEvent,
			},
			setupMock: func(dependencies *mockDependencies) {
				dependencies.eventsDao.EXPECT().
					DeleteEvent(taskqueueNs.VendorApiConsumerTaskQueue, samplePassengerEvent).
					Return(errors.New("failed to connect"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, md := newMockEventBusService(ctrl)
			tt.setupMock(md)
			if err := s.DequeueEventFromTaskQueue(tt.args.taskQueue, tt.args.passengerEvent); (err != nil) != tt.wantErr {
				t.Errorf("DequeueEventFromTaskQueue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_GetEventToProcess(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		taskQueue taskqueue.TaskQueue
	}
	tests := []struct {
		name      string
		setupMock func(dependencies *mockDependencies)
		args      args
		want      *eventbus.PassengerEvent
		wantErr   bool
	}{
		{
			name: "success get",
			args: args{
				taskQueue: taskqueueNs.VendorApiConsumerTaskQueue,
			},
			setupMock: func(dependencies *mockDependencies) {
				dependencies.eventsDao.EXPECT().
					GetEvent(taskqueueNs.VendorApiConsumerTaskQueue).
					Return(samplePassengerEvent, time.Now().Unix(), nil)
			},
			want:    samplePassengerEvent,
			wantErr: false,
		},
		{
			name: "failed to dequeue",
			args: args{
				taskQueue: taskqueueNs.VendorApiConsumerTaskQueue,
			},
			setupMock: func(dependencies *mockDependencies) {
				dependencies.eventsDao.EXPECT().
					GetEvent(taskqueueNs.VendorApiConsumerTaskQueue).
					Return(nil, int64(0), errors.New("failed to connect"))
			},
			wantErr: true,
		},
		{
			name: "found non processable data in db",
			args: args{
				taskQueue: taskqueueNs.VendorApiConsumerTaskQueue,
			},
			setupMock: func(dependencies *mockDependencies) {
				dependencies.eventsDao.EXPECT().
					GetEvent(taskqueueNs.VendorApiConsumerTaskQueue).
					Return(samplePassengerEvent, time.Now().Add(2*time.Hour).Unix(), nil)
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, md := newMockEventBusService(ctrl)
			tt.setupMock(md)
			got, err := s.GetEventToProcess(tt.args.taskQueue)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEventToProcess() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetEventToProcess() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_UpdatePassengerEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		taskQueue         taskqueue.TaskQueue
		oldPassenger      *eventbus.PassengerEvent
		newPassengerEvent *eventbus.PassengerEvent
		nextExecutionTime time.Time
	}
	tests := []struct {
		name      string
		setupMock func(dependencies *mockDependencies)
		args      args
		wantErr   bool
	}{
		{
			name: "success dequeue",
			args: args{
				taskQueue:         taskqueueNs.VendorApiConsumerTaskQueue,
				oldPassenger:      samplePassengerEvent,
				newPassengerEvent: samplePassengerEvent,
				nextExecutionTime: time.Unix(16000000, 0),
			},
			setupMock: func(dependencies *mockDependencies) {
				dependencies.eventsDao.EXPECT().
					UpdateEvent(taskqueueNs.VendorApiConsumerTaskQueue, samplePassengerEvent, samplePassengerEvent, time.Unix(16000000, 0)).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "failed to dequeue",
			args: args{
				taskQueue:         taskqueueNs.VendorApiConsumerTaskQueue,
				oldPassenger:      samplePassengerEvent,
				newPassengerEvent: samplePassengerEvent,
				nextExecutionTime: time.Unix(16000000, 0),
			},
			setupMock: func(dependencies *mockDependencies) {
				dependencies.eventsDao.EXPECT().
					UpdateEvent(taskqueueNs.VendorApiConsumerTaskQueue, samplePassengerEvent, samplePassengerEvent, time.Unix(16000000, 0)).
					Return(errors.New("failed to connect"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, md := newMockEventBusService(ctrl)
			tt.setupMock(md)
			if err := s.UpdatePassengerEvent(tt.args.taskQueue, tt.args.oldPassenger, tt.args.newPassengerEvent, tt.args.nextExecutionTime); (err != nil) != tt.wantErr {
				t.Errorf("UpdatePassengerEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
