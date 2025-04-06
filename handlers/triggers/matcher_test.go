package triggers

import (
	"fmt"
	"testing"
	"time"

	"github.com/julianshen/nebula/data"
)

func TestMatchTrigger(t *testing.T) {
	// Create a test event
	event := &data.Event{
		ID:           "evt1",
		EventType:    "user.created",
		EventVersion: "1.3.0",
		Namespace:    "core",
		ObjectType:   "user",
		ObjectID:     "u123",
		Timestamp:    time.Now(),
	}

	// Print event structure
	fmt.Printf("Event before payload init: %+v\n", event)
	fmt.Printf("Event.Payload before init: %+v\n", event.Payload)

	// Initialize payload
	event.Payload.Before = make(map[string]interface{})
	event.Payload.After = make(map[string]interface{})

	// Print payload after initialization
	fmt.Printf("Event.Payload after init: %+v\n", event.Payload)
	fmt.Printf("Event.Payload.After type: %T\n", event.Payload.After)

	// Set test values
	event.Payload.After["role"] = "admin"
	event.Payload.After["amount"] = 1500

	// Print final payload
	fmt.Printf("Event.Payload.After final: %+v\n", event.Payload.After)

	tests := []struct {
		name    string
		trigger data.Trigger
		want    bool
	}{
		{
			name: "simple eq match",
			trigger: data.Trigger{
				Enabled: true,
				RootGroup: data.ConditionGroup{
					Operator: "AND",
					Conditions: []data.Condition{
						{Field: "event_type", Operator: "eq", Value: "user.created"},
					},
				},
			},
			want: true,
		},
		{
			name: "nested AND OR match",
			trigger: data.Trigger{
				Enabled: true,
				RootGroup: data.ConditionGroup{
					Operator: "AND",
					Conditions: []data.Condition{
						{Field: "event_type", Operator: "eq", Value: "user.created"},
					},
					Groups: []data.ConditionGroup{
						{
							Operator: "OR",
							Conditions: []data.Condition{
								{Field: "payload.after.role", Operator: "eq", Value: "admin"},
								{Field: "payload.after.role", Operator: "eq", Value: "superuser"},
							},
						},
					},
				},
			},
			want: true,
		},
		{
			name: "numeric gt match",
			trigger: data.Trigger{
				Enabled: true,
				RootGroup: data.ConditionGroup{
					Operator: "AND",
					Conditions: []data.Condition{
						{Field: "payload.after.amount", Operator: "gt", Value: "1000"},
					},
				},
			},
			want: true,
		},
		{
			name: "numeric lt no match",
			trigger: data.Trigger{
				Enabled: true,
				RootGroup: data.ConditionGroup{
					Operator: "AND",
					Conditions: []data.Condition{
						{Field: "payload.after.amount", Operator: "lt", Value: "1000"},
					},
				},
			},
			want: false,
		},
		{
			name: "field not found",
			trigger: data.Trigger{
				Enabled: true,
				RootGroup: data.ConditionGroup{
					Operator: "AND",
					Conditions: []data.Condition{
						{Field: "payload.after.nonexistent", Operator: "eq", Value: "x"},
					},
				},
			},
			want: false,
		},
		{
			name: "disabled trigger",
			trigger: data.Trigger{
				Enabled: false,
				RootGroup: data.ConditionGroup{
					Operator: "AND",
					Conditions: []data.Condition{
						{Field: "event_type", Operator: "eq", Value: "user.created"},
					},
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.trigger.Enabled {
				if MatchTrigger(&tt.trigger, event) {
					t.Errorf("expected disabled trigger to not match")
				}
				return
			}
			got := MatchTrigger(&tt.trigger, event)
			if got != tt.want {
				t.Errorf("expected %v, got %v", tt.want, got)
			}
		})
	}
}
