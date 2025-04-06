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
root_group:
  operator: AND
  conditions:
    - field: payload.after.amount
      operator: gt
      value: "1000"
    - field: payload.after.status
      operator: eq
      value: "confirmed"
  groups:
    - operator: OR
      conditions:
        - field: payload.after.region
          operator: eq
          value: "US"
        - field: payload.after.region
          operator: eq
          value: "EU"
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

	// Verify root group
	if trigger.RootGroup.Operator != "AND" {
		t.Errorf("Expected RootGroup.Operator 'AND', got '%s'", trigger.RootGroup.Operator)
	}

	// Verify conditions
	if len(trigger.RootGroup.Conditions) != 2 {
		t.Fatalf("Expected 2 conditions, got %d", len(trigger.RootGroup.Conditions))
	}
	if trigger.RootGroup.Conditions[0].Field != "payload.after.amount" {
		t.Errorf("Expected condition 0 field 'payload.after.amount', got '%s'", trigger.RootGroup.Conditions[0].Field)
	}
	if trigger.RootGroup.Conditions[0].Operator != "gt" {
		t.Errorf("Expected condition 0 operator 'gt', got '%s'", trigger.RootGroup.Conditions[0].Operator)
	}
	if trigger.RootGroup.Conditions[0].Value != "1000" {
		t.Errorf("Expected condition 0 value '1000', got '%s'", trigger.RootGroup.Conditions[0].Value)
	}

	// Verify nested groups
	if len(trigger.RootGroup.Groups) != 1 {
		t.Fatalf("Expected 1 nested group, got %d", len(trigger.RootGroup.Groups))
	}
	if trigger.RootGroup.Groups[0].Operator != "OR" {
		t.Errorf("Expected nested group operator 'OR', got '%s'", trigger.RootGroup.Groups[0].Operator)
	}
	if len(trigger.RootGroup.Groups[0].Conditions) != 2 {
		t.Fatalf("Expected 2 conditions in nested group, got %d", len(trigger.RootGroup.Groups[0].Conditions))
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
