package taskqueue

import "persistent-queue/api/taskqueue"

// different task queues for consuming events
const (
	SnowflakeConsumerTaskQueue taskqueue.TaskQueue = "snowflakeconsumer-task-queue"
	VendorApiConsumerTaskQueue taskqueue.TaskQueue = "vendorapiconsumer-task-queue"
	FileConsumerTaskQueue      taskqueue.TaskQueue = "fileconsumer-task-queue"
)
