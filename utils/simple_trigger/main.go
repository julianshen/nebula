package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/julianshen/nebula/data"
	"github.com/julianshen/nebula/handlers/triggers"
)

// This utility creates a simple trigger and stores it in etcd
// It creates a trigger that matches orders with amount > 1000 and region = "US"

func main() {
	// Create a simple trigger
	trigger := createSimpleTrigger()

	// Connect to etcd
	etcdEndpoints := []string{"localhost:2379"}
	store, err := triggers.NewEtcdStore(etcdEndpoints, "/triggers/")
	if err != nil {
		log.Fatalf("Failed to create etcd store: %v", err)
	}
	defer store.Close()

	// Save the trigger
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = store.SaveTrigger(ctx, "sales", "simple-trigger", trigger)
	if err != nil {
		log.Fatalf("Failed to save trigger: %v", err)
	}

	fmt.Println("Successfully stored simple trigger in etcd")
	fmt.Println("Trigger conditions:")
	fmt.Println("- Amount > 1000")
	fmt.Println("- Region = US")
}

// createSimpleTrigger creates a simple trigger for testing
func createSimpleTrigger() *data.Trigger {
	return &data.Trigger{
		ID:         "simple-trigger",
		Name:       "Simple Trigger",
		Namespace:  "sales",
		ObjectType: "order",
		EventType:  "created",
		Enabled:    true,
		Criteria:   `event.payload.after.amount > 1000 && event.payload.after.region == "US"`,
	}
}
