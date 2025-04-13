# Event System Specification (v1.3)

## Overview

This document specifies the event system architecture and data formats for the Nebula Event Trigger System. It builds upon the previous v1.2 specification with enhancements to the trigger evaluation system and event processing pipeline.

## Event Structure

Events represent state changes in the system and follow this structure:

```json
{
  "event_id": "unique-event-identifier",
  "event_type": "object.action",
  "event_version": "1.3",
  "namespace": "domain-specific-namespace",
  "object_type": "type-of-object",
  "object_id": "identifier-of-affected-object",
  "timestamp": "2025-04-12T10:15:00Z",
  "actor": {
    "type": "user|system|service",
    "id": "actor-identifier"
  },
  "context": {
    "request_id": "original-request-identifier",
    "trace_id": "distributed-tracing-identifier"
  },
  "payload": {
    "before": {
      // Object state before the change (optional)
    },
    "after": {
      // Object state after the change (optional)
    }
  },
  "nats_meta": {
    "stream": "stream-name",
    "sequence": 12345,
    "received_at": "2025-04-12T10:15:01Z"
  }
}
```

### Field Descriptions

- **event_id**: Unique identifier for the event
- **event_type**: Dot-notation type of the event (e.g., "user.created", "order.updated")
- **event_version**: Version of the event schema (currently "1.3")
- **namespace**: Domain-specific namespace for the event
- **object_type**: Type of object affected by the event
- **object_id**: Identifier of the specific object affected
- **timestamp**: ISO 8601 timestamp when the event occurred
- **actor**: Entity that caused the event
  - **type**: Type of actor (user, system, service)
  - **id**: Identifier of the actor
- **context**: Additional contextual information
  - **request_id**: Original request identifier
  - **trace_id**: Distributed tracing identifier
- **payload**: Event-specific data
  - **before**: Object state before the change (optional)
  - **after**: Object state after the change (optional)
- **nats_meta**: NATS-specific metadata
  - **stream**: NATS stream name
  - **sequence**: Sequence number in the stream
  - **received_at**: Timestamp when the event was received by NATS

## Trigger Definition

Triggers define conditions for reacting to events and follow this structure:

```json
{
  "id": "trigger-identifier",
  "name": "human-readable-name",
  "namespace": "domain-specific-namespace",
  "object_type": "type-of-object",
  "event_type": "object.action",
  "criteria": "expression-based-condition",
  "description": "human-readable-description",
  "enabled": true
}
```

### Field Descriptions

- **id**: Unique identifier for the trigger
- **name**: Human-readable name for the trigger
- **namespace**: Domain-specific namespace for the trigger
- **object_type**: Type of object to match (optional, if empty matches all)
- **event_type**: Type of event to match (optional, if empty matches all)
- **criteria**: Expression-based condition for matching events
- **description**: Human-readable description of the trigger
- **enabled**: Whether the trigger is active

## Expression Language

The criteria field in trigger definitions uses the expr language (https://github.com/expr-lang/expr) to define conditions that are evaluated against events.

### Available Variables

- **event**: The full event object with all fields

### Examples

Match events where the amount is greater than 1000:
```
event.payload.after.amount > 1000
```

Match events where the region is "US":
```
event.payload.after.region == "US"
```

Match events with a combination of conditions:
```
event.event_type == "order.created" && event.payload.after.amount > 1000 && event.payload.after.region == "US"
```

### Custom Functions

- **has(obj, path)**: Checks if a nested path exists in an object
  - Example: `has(event.payload.after, "user.role")`

## Event Processing Pipeline

1. Events are published to NATS on subject patterns like `event.{namespace}.{object_type}.{event_type}`
2. The eventstore service subscribes to these events and stores them in MongoDB
3. The triggerd service also subscribes to these events and evaluates them against registered triggers
4. When a trigger matches an event, the corresponding action is executed

### Batch Processing

The eventstore service uses batch processing for efficiency:
- Default batch size: 1000 events
- Default batch timeout: 5 seconds

## Dynamic Trigger Management

Triggers are stored in etcd with the following features:
- Centralized storage accessible by all triggerd instances
- Dynamic reloading when triggers change
- Namespace-based organization
- gRPC API for trigger management

## Changes from v1.2

1. Enhanced expression language support with custom functions
2. Improved namespace-based event processing
3. Dynamic trigger reloading from etcd
4. Batch processing optimizations
5. Added support for NATS metadata
