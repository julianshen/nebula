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

type Condition struct {
	Field    string `json:"field" yaml:"field"`
	Operator string `json:"operator" yaml:"operator"` // eq, neq, gt, lt, gte, lte, contains, regex, etc.
	Value    string `json:"value" yaml:"value"`
}

type ConditionGroup struct {
	Operator   string           `json:"operator" yaml:"operator"` // AND, OR, NOT
	Conditions []Condition      `json:"conditions,omitempty" yaml:"conditions,omitempty"`
	Groups     []ConditionGroup `json:"groups,omitempty" yaml:"groups,omitempty"`
}

type Trigger struct {
	ID          string         `json:"id" yaml:"id"`
	Name        string         `json:"name" yaml:"name"`
	Namespace   string         `json:"namespace" yaml:"namespace"`
	ObjectType  string         `json:"object_type" yaml:"object_type"`
	EventType   string         `json:"event_type" yaml:"event_type"`
	RootGroup   ConditionGroup `json:"root_group" yaml:"root_group"`
	Description string         `json:"description,omitempty" yaml:"description,omitempty"`
	Enabled     bool           `json:"enabled" yaml:"enabled"`
}

// ToYAML marshals the trigger to YAML
func (t *Trigger) ToYAML() ([]byte, error) {
	return yaml.Marshal(t)
}
