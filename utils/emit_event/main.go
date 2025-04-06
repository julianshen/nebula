package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/julianshen/nebula/data"
	"github.com/nats-io/nats.go"
)

// This utility emits test events to NATS for testing the trigger system

func main() {
	// Parse command line flags
	var (
		natsURL    = flag.String("nats", "nats://localhost:4222", "NATS server URL")
		subject    = flag.String("subject", "event.test", "NATS subject")
		eventType  = flag.String("type", "created", "Event type")
		namespace  = flag.String("namespace", "sales", "Namespace")
		objectType = flag.String("object-type", "order", "Object type")
		objectID   = flag.String("id", "", "Object ID (defaults to timestamp)")
		amount     = flag.Int("amount", 1500, "Order amount")
		status     = flag.String("status", "confirmed", "Order status")
		region     = flag.String("region", "US", "Order region")
	)

	flag.Parse()

	// Generate a default object ID if not provided
	if *objectID == "" {
		*objectID = fmt.Sprintf("order-%d", time.Now().Unix())
	}

	// Create the event
	event := createOrderEvent(*objectID, *eventType, *namespace, *objectType, *amount, *status, *region)

	// Marshal the event to JSON
	eventJSON, err := json.MarshalIndent(event, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal event to JSON: %v", err)
	}

	// Print the event
	fmt.Println("Emitting event:")
	fmt.Println(string(eventJSON))

	// Connect to NATS
	nc, err := nats.Connect(*natsURL)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer nc.Close()

	// Publish the event
	err = nc.Publish(*subject, eventJSON)
	if err != nil {
		log.Fatalf("Failed to publish event: %v", err)
	}

	// Ensure the event is delivered
	err = nc.Flush()
	if err != nil {
		log.Fatalf("Failed to flush NATS connection: %v", err)
	}

	fmt.Printf("Event published to %s\n", *subject)
}

// createOrderEvent creates a test order event
func createOrderEvent(id, eventType, namespace, objectType string, amount int, status, region string) data.Event {
	event := data.Event{
		ID:           id,
		EventType:    eventType,
		EventVersion: "1.0",
		Namespace:    namespace,
		ObjectType:   objectType,
		ObjectID:     id,
		Timestamp:    time.Now(),
		Actor: struct {
			Type string `json:"type"`
			ID   string `json:"id"`
		}{
			Type: "user",
			ID:   "test-user",
		},
		Context: struct {
			RequestID string `json:"request_id"`
			TraceID   string `json:"trace_id"`
		}{
			RequestID: "req-" + strconv.FormatInt(time.Now().UnixNano(), 10),
			TraceID:   "trace-" + strconv.FormatInt(time.Now().UnixNano(), 10),
		},
	}

	// Initialize payload
	event.Payload.Before = make(map[string]interface{})
	event.Payload.After = make(map[string]interface{})

	// Set test values
	event.Payload.After["amount"] = amount
	event.Payload.After["status"] = status
	event.Payload.After["region"] = region

	return event
}
