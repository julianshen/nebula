package triggers

import (
	"fmt"
	"strings"

	"github.com/expr-lang/expr"
	"github.com/julianshen/nebula/data"
)

// MatchTrigger returns true if the event satisfies the trigger's criteria.
// It supports:
// Expression-based matching using the expr library (preferred)
//
// Expression-based matching evaluates a string expression against the event object.
// The expression must evaluate to a boolean value.
// Example: event.event_type == "user.created" && event.payload.after.role == "admin"
//
// See the event system specification for more details on the expression language.
func MatchTrigger(trigger *data.Trigger, event *data.Event) (bool, error) {
	if trigger == nil || !trigger.Enabled {
		return false, nil
	}

	// If criteria is empty, match based on event type and namespace
	if trigger.Criteria == "" {
		return (trigger.EventType == "" || trigger.EventType == event.EventType) &&
			(trigger.Namespace == "" || trigger.Namespace == event.Namespace) &&
			(trigger.ObjectType == "" || trigger.ObjectType == event.ObjectType), nil
	}

	// If the trigger has a criteria expression, evaluate it
	return evaluateTriggerCriteria(event, trigger.Criteria)
}

// has(obj, "a.b.c") returns true if all keys exist down the path
func has(args ...any) (any, error) {
	if len(args) != 2 {
		return false, fmt.Errorf("has() expects 2 arguments")
	}
	root, ok := args[0].(map[string]interface{})
	if !ok {
		return false, nil
	}
	path, ok := args[1].(string)
	if !ok {
		return false, nil
	}

	parts := strings.Split(path, ".")
	current := root
	for i, part := range parts {
		val, exists := current[part]
		if !exists {
			return false, nil
		}
		if i == len(parts)-1 {
			return true, nil
		}
		// If not final, it must be a nested map
		next, ok := val.(map[string]interface{})
		if !ok {
			return false, nil
		}
		current = next
	}
	return true, nil
}

// EvaluateTriggerCriteria safely evaluates a criteria string against the given event
func evaluateTriggerCriteria(event *data.Event, criteria string) (bool, error) {
	// If criteria is empty, match based on event type and namespace
	if criteria == "" {
		// For empty criteria, we'll just return true since we don't have trigger information here
		// The actual matching based on event type and namespace is done in the MatchTrigger function
		return true, nil
	}

	// Create a map representation of the event that matches JSON field names
	eventMap := map[string]interface{}{
		"event_id":      event.ID,
		"event_type":    event.EventType,
		"event_version": event.EventVersion,
		"namespace":     event.Namespace,
		"object_type":   event.ObjectType,
		"object_id":     event.ObjectID,
		"timestamp":     event.Timestamp,
		"actor": map[string]interface{}{
			"type": event.Actor.Type,
			"id":   event.Actor.ID,
		},
		"context": map[string]interface{}{
			"request_id": event.Context.RequestID,
			"trace_id":   event.Context.TraceID,
		},
		"payload": map[string]interface{}{
			"before": event.Payload.Before,
			"after":  event.Payload.After,
		},
		"nats_meta": map[string]interface{}{
			"stream":      event.NatsMeta.Stream,
			"sequence":    event.NatsMeta.Sequence,
			"received_at": event.NatsMeta.ReceivedAt,
		},
	}

	// Create environment with event as the root variable
	env := map[string]interface{}{
		"event": eventMap,
	}

	// Compile the expression with custom functions
	options := []expr.Option{
		expr.Env(env),
		expr.Function("has", has),
	}

	program, err := expr.Compile(criteria, options...)
	if err != nil {
		return false, fmt.Errorf("failed to compile criteria: %w", err)
	}

	// Run the compiled expression
	output, err := expr.Run(program, env)
	if err != nil {
		return false, fmt.Errorf("failed to evaluate criteria: %w", err)
	}

	// Must return boolean
	result, ok := output.(bool)
	if !ok {
		return false, fmt.Errorf("expression did not return a boolean")
	}

	return result, nil
}
