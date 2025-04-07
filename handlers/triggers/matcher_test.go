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
			name: "basic matching - event type only",
			trigger: data.Trigger{
				Enabled:   true,
				EventType: "user.created",
			},
			want: true,
		},
		{
			name: "basic matching - namespace only",
			trigger: data.Trigger{
				Enabled:   true,
				Namespace: "core",
			},
			want: true,
		},
		{
			name: "basic matching - object type only",
			trigger: data.Trigger{
				Enabled:    true,
				ObjectType: "user",
			},
			want: true,
		},
		{
			name: "basic matching - all fields match",
			trigger: data.Trigger{
				Enabled:    true,
				EventType:  "user.created",
				Namespace:  "core",
				ObjectType: "user",
			},
			want: true,
		},
		{
			name: "basic matching - no match",
			trigger: data.Trigger{
				Enabled:   true,
				EventType: "user.updated",
			},
			want: false,
		},
		{
			name: "expr simple match",
			trigger: data.Trigger{
				Enabled:  true,
				Criteria: `event.event_type == "user.created"`,
			},
			want: true,
		},
		{
			name: "expr payload match",
			trigger: data.Trigger{
				Enabled:  true,
				Criteria: `event.event_type == "user.created" && event.payload.after.role == "admin"`,
			},
			want: true,
		},
		{
			name: "expr numeric comparison",
			trigger: data.Trigger{
				Enabled:  true,
				Criteria: `event.payload.after.amount > 1000`,
			},
			want: true,
		},
		{
			name: "expr no match",
			trigger: data.Trigger{
				Enabled:  true,
				Criteria: `event.event_type == "user.updated"`,
			},
			want: false,
		},
		{
			name: "expr field not found",
			trigger: data.Trigger{
				Enabled:  true,
				Criteria: `event.payload.after.nonexistent == "x"`,
			},
			want: false,
		},
		{
			name: "disabled trigger",
			trigger: data.Trigger{
				Enabled:   false,
				EventType: "user.created",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.trigger.Enabled {
				matched, err := MatchTrigger(&tt.trigger, event)
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if matched {
					t.Errorf("expected disabled trigger to not match")
				}
				return
			}
			got, err := MatchTrigger(&tt.trigger, event)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("expected %v, got %v", tt.want, got)
			}
		})
	}
}
