# Nebula: Event Trigger System

Nebula is a Go-based event trigger system that uses etcd for configuration storage with dynamic reloading capabilities. It provides a robust, centralized way to manage trigger configurations across multiple instances of the trigger service, with real-time updates when configurations change.

## Features

- **etcd-based Trigger Storage**: Store trigger definitions in etcd for centralized management
- **Dynamic Reloading**: Automatically reload trigger definitions when they change in etcd
- **Flexible Trigger Conditions**: Define complex trigger conditions with support for logical operators
- **Event Storage**: Store events in MongoDB for historical analysis
- **NATS Integration**: Use NATS for event distribution and processing
- **gRPC API**: Manage triggers via a gRPC API
- **Docker Support**: Run the complete system with Docker Compose

## Architecture

Nebula consists of the following components:

1. **Trigger Service (`triggerd`)**: Evaluates events against trigger definitions
2. **Event Store Service (`eventstore`)**: Stores events in MongoDB
3. **etcd**: Stores trigger definitions
4. **NATS**: Distributes events between services
5. **MongoDB**: Stores event history

## Getting Started

### Prerequisites

- Go 1.16 or later
- Docker and Docker Compose
- etcd
- NATS
- MongoDB

### Running with Docker Compose

The easiest way to run Nebula is with Docker Compose:

```bash
docker-compose up -d
```

This will start all required services: etcd, NATS, MongoDB, triggerd, and eventstore.

### Running Manually

1. Start the required services (etcd, NATS, MongoDB)
2. Start the eventstore service:

```bash
go run services/eventstore/main.go
```

3. Start the triggerd service:

```bash
go run services/triggerd/main.go
```

## Managing Triggers

### Using the gRPC Client

You can manage triggers using the provided gRPC client utility:

```bash
# List all triggers in a namespace
go run utils/grpc_client/main.go --cmd list --namespace sales

# Add a new trigger
go run utils/grpc_client/main.go --cmd add --namespace sales --name high-value-order --field1 payload.after.amount --op1 gt --value1 1000 --field2 payload.after.region --op2 eq --value2 US

# Update an existing trigger
go run utils/grpc_client/main.go --cmd update --namespace sales --id high-value-order --name high-value-order --field1 payload.after.amount --op1 gt --value1 2000 --field2 payload.after.region --op2 eq --value2 US

# Remove a trigger
go run utils/grpc_client/main.go --cmd remove --namespace sales --id high-value-order
```

### Using the etcd Utility

You can also create triggers directly in etcd using the provided utility:

```bash
go run utils/simple_trigger/main.go
```

This will create a simple trigger that matches orders with amount > 1000 and region = "US".

## Emitting Events

You can emit test events using the provided utility:

```bash
go run utils/emit_event/main.go --amount 1500 --region US
```

## Checking Stored Events

You can check the events stored in MongoDB using the provided utility:

```bash
go run utils/check_mongo/main.go
```

## End-to-End Testing

For a complete end-to-end test, follow the instructions in [utils/end_to_end_test.md](utils/end_to_end_test.md).

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

Similar projects that inspired Nebula:

- [rynbrd/sentinel](https://github.com/rynbrd/sentinel): Triggered templating and command execution for etcd
- [sheldonh/etcd-trigger](https://github.com/sheldonh/etcd-trigger): Send values from etcd to an HTTP endpoint on change
