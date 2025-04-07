package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/julianshen/nebula/data"
	"github.com/julianshen/nebula/handlers/triggers"
)

// This is a simple utility to store a trigger in etcd
// It creates a sample trigger and stores it in etcd

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: store_trigger <namespace> <name>")
		os.Exit(1)
	}

	namespace := os.Args[1]
	name := os.Args[2]

	// Create a test trigger
	trigger := createHighValueOrderTrigger()

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

	err = store.SaveTrigger(ctx, namespace, name, trigger)
	if err != nil {
		log.Fatalf("Failed to save trigger: %v", err)
	}

	fmt.Printf("Successfully stored trigger %s/%s in etcd\n", namespace, name)
}

// createHighValueOrderTrigger creates a test trigger for high-value orders
func createHighValueOrderTrigger() *data.Trigger {
	return &data.Trigger{
		ID:         "high-value-order",
		Name:       "High Value Order",
		Namespace:  "sales",
		ObjectType: "order",
		EventType:  "created",
		Enabled:    true,
		Criteria:   `event.payload.after.amount > 1000 && event.payload.after.status == "confirmed" && (event.payload.after.region == "US" || event.payload.after.region == "EU")`,
	}
}
