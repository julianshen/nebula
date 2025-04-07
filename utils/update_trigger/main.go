package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/julianshen/nebula/data"
	"github.com/julianshen/nebula/handlers/triggers"
)

// This utility updates the simple trigger to change its conditions
// It modifies the trigger to match orders with amount > 500 and region = "EU"

func main() {
	// Create the updated trigger
	trigger := createUpdatedTrigger()

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

	fmt.Println("Successfully updated simple trigger in etcd")
	fmt.Println("New trigger conditions:")
	fmt.Println("- Amount > 500")
	fmt.Println("- Region = EU")
}

// createUpdatedTrigger creates an updated version of the simple trigger
func createUpdatedTrigger() *data.Trigger {
	return &data.Trigger{
		ID:         "simple-trigger",
		Name:       "Simple Trigger (Updated)",
		Namespace:  "sales",
		ObjectType: "order",
		EventType:  "created",
		Enabled:    true,
		Criteria:   `event.payload.after.amount > 500 && event.payload.after.region == "EU"`,
	}
}
