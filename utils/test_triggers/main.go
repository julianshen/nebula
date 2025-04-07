package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/julianshen/nebula/data"
	"github.com/julianshen/nebula/handlers/triggers"
)

// This is a simple test script to demonstrate the trigger system
// It creates some test triggers, simulates events, and tests matching

func main() {
	// Create a test trigger
	trigger1 := createHighValueOrderTrigger()
	triggerYAML, err := trigger1.ToYAML()
	if err != nil {
		log.Fatalf("Failed to marshal trigger to YAML: %v", err)
	}

	// Print the trigger YAML
	fmt.Println("=== Trigger YAML ===")
	fmt.Println(string(triggerYAML))
	fmt.Println("===================")

	// Test loading the trigger from YAML
	reader := strings.NewReader(string(triggerYAML))
	loadedTrigger, err := triggers.LoadTrigger(reader)
	if err != nil {
		log.Fatalf("Failed to load trigger from YAML: %v", err)
	}

	fmt.Printf("Loaded trigger: %s - %s\n", loadedTrigger.ID, loadedTrigger.Name)

	// Create some test events
	events := []data.Event{
		createOrderEvent("order1", 500, "confirmed", "US"),  // Not high value, but confirmed and US
		createOrderEvent("order2", 1500, "confirmed", "US"), // High value, confirmed, and US - should match
		createOrderEvent("order3", 1500, "pending", "US"),   // High value, but pending
		createOrderEvent("order4", 1500, "confirmed", "CA"), // High value, confirmed, but not US or EU
		createOrderEvent("order5", 1500, "confirmed", "EU"), // High value, confirmed, and EU - should match
	}

	// Test matching each event against the trigger
	fmt.Println("\n=== Testing Trigger Matching ===")
	for i, event := range events {
		matched, err := triggers.MatchTrigger(loadedTrigger, &events[i])
		if err != nil {
			log.Fatalf("Error matching trigger: %v", err)
		}
		fmt.Printf("Event %s (amount: %v, status: %s, region: %s): %v\n",
			event.ID, event.Payload.After["amount"], event.Payload.After["status"],
			event.Payload.After["region"], matched)
	}

	// Test in-memory trigger store
	fmt.Println("\n=== Testing In-Memory Trigger Store ===")
	testInMemoryTriggerStore(trigger1, events)

	fmt.Println("\nAll tests completed successfully!")
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

// createOrderEvent creates a test order event
func createOrderEvent(id string, amount int, status, region string) data.Event {
	event := data.Event{
		ID:           id,
		EventType:    "created",
		EventVersion: "1.0",
		Namespace:    "sales",
		ObjectType:   "order",
		ObjectID:     id,
		Timestamp:    time.Now(),
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

// testInMemoryTriggerStore tests the in-memory trigger store
func testInMemoryTriggerStore(trigger *data.Trigger, events []data.Event) {
	// Create a simple in-memory store
	store := &InMemoryTriggerStore{
		triggers: make(map[string]map[string]*data.Trigger),
	}

	// Add the trigger to the store
	ctx := context.Background()
	err := store.SaveTrigger(ctx, trigger.Namespace, trigger.ID, trigger)
	if err != nil {
		log.Fatalf("Failed to save trigger: %v", err)
	}

	// Get all triggers
	allTriggers := store.GetAllTriggers()
	fmt.Printf("Store contains %d triggers\n", len(allTriggers))

	// Test matching events against all triggers in the store
	for i, event := range events {
		fmt.Printf("Testing event %s against all triggers in store...\n", event.ID)
		for _, t := range allTriggers {
			matched, err := triggers.MatchTrigger(t, &events[i])
			if err != nil {
				log.Fatalf("Error matching trigger: %v", err)
			}
			if matched {
				fmt.Printf("  Matched trigger: %s - %s\n", t.ID, t.Name)
			}
		}
	}

	// Test deleting a trigger
	err = store.DeleteTrigger(ctx, trigger.Namespace, trigger.ID)
	if err != nil {
		log.Fatalf("Failed to delete trigger: %v", err)
	}

	// Verify it was deleted
	allTriggers = store.GetAllTriggers()
	fmt.Printf("Store contains %d triggers after deletion\n", len(allTriggers))
}

// InMemoryTriggerStore is a simple in-memory implementation of TriggerStore for testing
type InMemoryTriggerStore struct {
	triggers map[string]map[string]*data.Trigger
}

func (s *InMemoryTriggerStore) LoadAll(ctx context.Context) error {
	return nil
}

func (s *InMemoryTriggerStore) Watch(ctx context.Context) {
	// No-op for in-memory store
}

func (s *InMemoryTriggerStore) GetTriggers(namespace string) []*data.Trigger {
	namespaceTriggers, ok := s.triggers[namespace]
	if !ok {
		return nil
	}

	triggers := make([]*data.Trigger, 0, len(namespaceTriggers))
	for _, trigger := range namespaceTriggers {
		triggers = append(triggers, trigger)
	}

	return triggers
}

func (s *InMemoryTriggerStore) GetAllTriggers() []*data.Trigger {
	var allTriggers []*data.Trigger
	for _, namespaceTriggers := range s.triggers {
		for _, trigger := range namespaceTriggers {
			allTriggers = append(allTriggers, trigger)
		}
	}

	return allTriggers
}

func (s *InMemoryTriggerStore) SaveTrigger(ctx context.Context, namespace, name string, trigger *data.Trigger) error {
	if _, ok := s.triggers[namespace]; !ok {
		s.triggers[namespace] = make(map[string]*data.Trigger)
	}
	s.triggers[namespace][name] = trigger
	return nil
}

func (s *InMemoryTriggerStore) DeleteTrigger(ctx context.Context, namespace, name string) error {
	if namespaceTriggers, ok := s.triggers[namespace]; ok {
		delete(namespaceTriggers, name)
	}
	return nil
}

func (s *InMemoryTriggerStore) Close() error {
	return nil
}
