package data

import (
	"time"

	yaml "gopkg.in/yaml.v3"
)

// Event represents a state change in the system following v1.2 spec
type Event struct {
	ID           string    `json:"event_id"`
	EventType    string    `json:"event_type"`
	EventVersion string    `json:"event_version"`
	Namespace    string    `json:"namespace"`
	ObjectType   string    `json:"object_type"`
	ObjectID     string    `json:"object_id"`
	Timestamp    time.Time `json:"timestamp"`
	Actor        struct {
		Type string `json:"type"`
		ID   string `json:"id"`
	} `json:"actor"`
	Context struct {
		RequestID string `json:"request_id"`
		TraceID   string `json:"trace_id"`
	} `json:"context"`
	Payload struct {
		Before map[string]interface{} `json:"before,omitempty"`
		After  map[string]interface{} `json:"after,omitempty"`
	} `json:"payload"`
	NatsMeta struct {
		Stream     string    `json:"stream"`
		Sequence   uint64    `json:"sequence"`
		ReceivedAt time.Time `json:"received_at"`
	} `json:"nats_meta"`
}

type Trigger struct {
	ID         string `json:"id" yaml:"id"`
	Name       string `json:"name" yaml:"name"`
	Namespace  string `json:"namespace" yaml:"namespace"`
	ObjectType string `json:"object_type" yaml:"object_type"`
	EventType  string `json:"event_type" yaml:"event_type"`
	// Criteria is an expression that is evaluated against the event.
	// It uses the expr language (https://github.com/expr-lang/expr) and must evaluate to a boolean.
	// Example: event.event_type == "user.created" && event.payload.after.role == "admin"
	Criteria    string `json:"criteria" yaml:"criteria"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	Enabled     bool   `json:"enabled" yaml:"enabled"`
}

// ToYAML marshals the trigger to YAML
func (t *Trigger) ToYAML() ([]byte, error) {
	return yaml.Marshal(t)
}
