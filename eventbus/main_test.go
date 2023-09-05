package eventbus

import (
	"github.com/golang/mock/gomock"
	"persistent-queue/api/taskqueue"
	"persistent-queue/eventbus/dao/mocks"
	taskqueueNs "persistent-queue/pkg/taskqueue"
)

type mockDependencies struct {
	eventsDao            *mocks.MockEventsDao
	registeredTaskQueues []taskqueue.TaskQueue
}

func newMockEventBusService(ctr *gomock.Controller) (*Service, *mockDependencies) {
	mockEventsDao := mocks.NewMockEventsDao(ctr)
	return &Service{
			eventsDao: mockEventsDao,
			registeredTaskQueues: []taskqueue.TaskQueue{
				taskqueueNs.SnowflakeConsumerTaskQueue,
				taskqueueNs.FileConsumerTaskQueue, taskqueueNs.VendorApiConsumerTaskQueue,
			},
		}, &mockDependencies{eventsDao: mockEventsDao, registeredTaskQueues: []taskqueue.TaskQueue{
			taskqueueNs.SnowflakeConsumerTaskQueue,
			taskqueueNs.FileConsumerTaskQueue, taskqueueNs.VendorApiConsumerTaskQueue,
		}}
}
