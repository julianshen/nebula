# End-to-End Test for Nebula Event Trigger System

This document outlines the steps to perform an end-to-end test of the Nebula Event Trigger System.

## Prerequisites

Ensure that the following services are running:

```bash
# Start all required services (MongoDB, NATS, etcd)
docker-compose up -d
```

## Test Steps

### 1. Create a Simple Trigger

First, create a simple trigger that matches orders with amount > 1000 and region = "US":

```bash
go run utils/simple_trigger/main.go
```

### 2. Start the Services

Start the triggerd service in one terminal:

```bash
go run services/triggerd/main.go
```

Start the eventstore service in another terminal:

```bash
go run services/eventstore/main.go
```

### 3. Emit a Matching Event

Emit an event that matches the trigger (amount > 1000, region = "US"):

```bash
go run utils/emit_event/main.go --amount 1500 --region US
```

Check the triggerd logs to verify that the trigger was matched.

### 4. Emit a Non-Matching Event

Emit an event that doesn't match the trigger:

```bash
go run utils/emit_event/main.go --amount 500 --region US
```

Check the triggerd logs to verify that the trigger was not matched.

### 5. Check MongoDB for Stored Events

Verify that both events were stored in MongoDB:

```bash
go run utils/check_mongo/main.go
```

### 6. Update the Trigger

Update the trigger to match orders with amount > 500 and region = "EU":

```bash
go run utils/update_trigger/main.go
```

### 7. Test the Updated Trigger

Emit an event that matches the updated trigger:

```bash
go run utils/emit_event/main.go --amount 1000 --region EU
```

Check the triggerd logs to verify that the trigger was matched.

Emit an event that doesn't match the updated trigger:

```bash
go run utils/emit_event/main.go --amount 1000 --region US
```

Check the triggerd logs to verify that the trigger was not matched.

### 8. Check MongoDB Again

Verify that all events were stored in MongoDB:

```bash
go run utils/check_mongo/main.go
```

## Expected Results

1. The triggerd service should load the trigger from etcd on startup.
2. When a matching event is emitted, the triggerd service should log that the trigger was matched.
3. When a non-matching event is emitted, the triggerd service should not log a match.
4. The eventstore service should store all events in MongoDB.
5. When the trigger is updated, the triggerd service should automatically reload it.
6. The updated trigger should match different events based on the new conditions.
