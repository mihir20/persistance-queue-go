// Code generated by MockGen. DO NOT EDIT.
// Source: dao.go

// Package mocks is a generated GoMock package.
package mocks

import (
	eventbus "persistent-queue/api/eventbus"
	taskqueue "persistent-queue/api/taskqueue"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
)

// MockEventsDao is a mock of EventsDao interface.
type MockEventsDao struct {
	ctrl     *gomock.Controller
	recorder *MockEventsDaoMockRecorder
}

// MockEventsDaoMockRecorder is the mock recorder for MockEventsDao.
type MockEventsDaoMockRecorder struct {
	mock *MockEventsDao
}

// NewMockEventsDao creates a new mock instance.
func NewMockEventsDao(ctrl *gomock.Controller) *MockEventsDao {
	mock := &MockEventsDao{ctrl: ctrl}
	mock.recorder = &MockEventsDaoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEventsDao) EXPECT() *MockEventsDaoMockRecorder {
	return m.recorder
}

// CountEventsInQueue mocks base method.
func (m *MockEventsDao) CountEventsInQueue(taskQueues []taskqueue.TaskQueue) (map[taskqueue.TaskQueue]int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountEventsInQueue", taskQueues)
	ret0, _ := ret[0].(map[taskqueue.TaskQueue]int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CountEventsInQueue indicates an expected call of CountEventsInQueue.
func (mr *MockEventsDaoMockRecorder) CountEventsInQueue(taskQueues interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountEventsInQueue", reflect.TypeOf((*MockEventsDao)(nil).CountEventsInQueue), taskQueues)
}

// CreateEvent mocks base method.
func (m *MockEventsDao) CreateEvent(event *eventbus.PassengerEvent, taskQueues []taskqueue.TaskQueue) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateEvent", event, taskQueues)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateEvent indicates an expected call of CreateEvent.
func (mr *MockEventsDaoMockRecorder) CreateEvent(event, taskQueues interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateEvent", reflect.TypeOf((*MockEventsDao)(nil).CreateEvent), event, taskQueues)
}

// DeleteEvent mocks base method.
func (m *MockEventsDao) DeleteEvent(taskQueue taskqueue.TaskQueue, event *eventbus.PassengerEvent) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteEvent", taskQueue, event)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteEvent indicates an expected call of DeleteEvent.
func (mr *MockEventsDaoMockRecorder) DeleteEvent(taskQueue, event interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteEvent", reflect.TypeOf((*MockEventsDao)(nil).DeleteEvent), taskQueue, event)
}

// GetEvents mocks base method.
func (m *MockEventsDao) GetEvents(taskQueue taskqueue.TaskQueue, cutOffTime time.Time, countOfEvents int64) ([]*eventbus.PassengerEvent, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEvents", taskQueue, cutOffTime, countOfEvents)
	ret0, _ := ret[0].([]*eventbus.PassengerEvent)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEvents indicates an expected call of GetEvents.
func (mr *MockEventsDaoMockRecorder) GetEvents(taskQueue, cutOffTime, countOfEvents interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEvents", reflect.TypeOf((*MockEventsDao)(nil).GetEvents), taskQueue, cutOffTime, countOfEvents)
}

// UpdateEvent mocks base method.
func (m *MockEventsDao) UpdateEvent(taskQueue taskqueue.TaskQueue, oldPassenger, updatedPassengerEvent *eventbus.PassengerEvent, nextExecutionTime time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateEvent", taskQueue, oldPassenger, updatedPassengerEvent, nextExecutionTime)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateEvent indicates an expected call of UpdateEvent.
func (mr *MockEventsDaoMockRecorder) UpdateEvent(taskQueue, oldPassenger, updatedPassengerEvent, nextExecutionTime interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateEvent", reflect.TypeOf((*MockEventsDao)(nil).UpdateEvent), taskQueue, oldPassenger, updatedPassengerEvent, nextExecutionTime)
}
