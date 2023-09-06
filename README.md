# Persistent-Queue
A high-performance, distributed, and persistent queue library for managing messages or tasks in your applications.
 

## Introduction
A distributed persistent queue is a data structure used in distributed computing and distributed systems to manage and store messages or tasks in a reliable and scalable manner. It allows multiple producers to enqueue (add) messages or tasks and multiple consumers to dequeue (remove and process) them while ensuring durability and fault tolerance.

### Features
- Distributed: Seamlessly integrates into distributed systems.
- Persistence: Ensures messages are stored persistently, preventing data loss.
- Reliability: Built with reliability and fault tolerance in mind.
- Scalability: Easily scales horizontally as your application grows.
- Message Ordering: Supports FIFO.

## Getting Started

---
### Pre-requisites
Need to have go1.18+ and docker installed on your system

### Setup
1. Running Frontend Server
```shell
docker-compose up frontend
```
this will start frontend server at `localhost` port `8080`
2. Running Consumers
For demo purposes we have setup `3` demo consumers
   1. Snowflake Consumer (All the events will be successfully consumed by this consumer)
        ```shell
        docker-compose up snowflake-api-consumer
      ```
   2. File Consumer (All events will be permanently getting failed by this consumer)
      ```shell
      docker-compose up file-consumer
      ```
   3. Vendor API Consumer (All the events will be transiently getting failed in this consumer)
        ```shell
        docker-compose up vendor-api-consumer
        ```
## APIs 

---
### 1. Healthcheck
```http
GET /healthcheck
```

## Responses

Will return number of messages available in each task queue

```javascript
{
    "fileconsumer-task-queue": 4, 
    "snowflakeconsumer-task-queue": 5,
    "vendorapiconsumer-task-queue": 7
}
```
---

### 2. Publish

```http
POST /publish
```

## Request Body
```javascript
{
    "userid": "uuid",
    "payload": "string"
}
```

## Responses

Will return plain text response, with http status code 
```javascript
"published event"
```

## Status Codes

Gophish returns the following status codes in its API:

| Status Code | Description |
| :--- | :--- |
| 200 | `OK` |
| 500 | `INTERNAL SERVER ERROR` |

