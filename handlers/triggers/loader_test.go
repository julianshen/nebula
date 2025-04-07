package triggers

import (
	"strings"
	"testing"
)

func TestLoadTrigger(t *testing.T) {
	// Test YAML content
	yamlContent := `
id: trigger-123
name: High Value Order
namespace: sales
object_type: order
event_type: created
enabled: true
criteria: event.payload.after.amount > 1000 && event.payload.after.status == "confirmed" && (event.payload.after.region == "US" || event.payload.after.region == "EU")
`

	// Create a reader from the YAML content
	reader := strings.NewReader(yamlContent)

	// Load the trigger
	trigger, err := LoadTrigger(reader)
	if err != nil {
		t.Fatalf("Failed to load trigger: %v", err)
	}

	// Verify the trigger properties
	if trigger.ID != "trigger-123" {
		t.Errorf("Expected ID 'trigger-123', got '%s'", trigger.ID)
	}
	if trigger.Name != "High Value Order" {
		t.Errorf("Expected Name 'High Value Order', got '%s'", trigger.Name)
	}
	if trigger.Namespace != "sales" {
		t.Errorf("Expected Namespace 'sales', got '%s'", trigger.Namespace)
	}
	if trigger.ObjectType != "order" {
		t.Errorf("Expected ObjectType 'order', got '%s'", trigger.ObjectType)
	}
	if trigger.EventType != "created" {
		t.Errorf("Expected EventType 'created', got '%s'", trigger.EventType)
	}
	if !trigger.Enabled {
		t.Error("Expected Enabled to be true")
	}

	// Verify criteria
	expectedCriteria := `event.payload.after.amount > 1000 && event.payload.after.status == "confirmed" && (event.payload.after.region == "US" || event.payload.after.region == "EU")`
	if trigger.Criteria != expectedCriteria {
		t.Errorf("Expected criteria '%s', got '%s'", expectedCriteria, trigger.Criteria)
	}
}

func TestLoadTrigger_InvalidYAML(t *testing.T) {
	// Invalid YAML content
	yamlContent := `
id: trigger-123
name: Invalid Trigger
invalid yaml format
`

	// Create a reader from the YAML content
	reader := strings.NewReader(yamlContent)

	// Load the trigger
	_, err := LoadTrigger(reader)
	if err == nil {
		t.Fatal("Expected error for invalid YAML, got nil")
	}
}
