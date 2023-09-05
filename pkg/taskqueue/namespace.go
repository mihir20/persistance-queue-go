package taskqueue

import "persistent-queue/api/taskqueue"

const (
	SnowflakeConsumerTaskQueue taskqueue.TaskQueue = "snowflakeconsumer-task-queue"
	VendorApiConsumerTaskQueue taskqueue.TaskQueue = "vendorapiconsumer-task-queue"
	FileConsumerTaskQueue      taskqueue.TaskQueue = "vendorapiconsumer-task-queue"
)
