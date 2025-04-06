package triggers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/julianshen/nebula/data"
)

// MatchTrigger returns true if the event satisfies the trigger's root condition group
func MatchTrigger(trigger *data.Trigger, event *data.Event) bool {
	if trigger == nil || !trigger.Enabled {
		return false
	}
	return evalGroup(trigger.RootGroup, event)
}

func evalGroup(group data.ConditionGroup, event *data.Event) bool {
	fmt.Printf("Evaluating group: %+v\n", group)
	switch strings.ToUpper(group.Operator) {
	case "AND":
		for _, cond := range group.Conditions {
			res := evalCondition(cond, event)
			fmt.Printf("AND condition %+v result: %v\n", cond, res)
			if !res {
				return false
			}
		}
		for _, sub := range group.Groups {
			res := evalGroup(sub, event)
			fmt.Printf("AND subgroup %+v result: %v\n", sub, res)
			if !res {
				return false
			}
		}
		return true
	case "OR":
		for _, cond := range group.Conditions {
			res := evalCondition(cond, event)
			fmt.Printf("OR condition %+v result: %v\n", cond, res)
			if res {
				return true
			}
		}
		for _, sub := range group.Groups {
			res := evalGroup(sub, event)
			fmt.Printf("OR subgroup %+v result: %v\n", sub, res)
			if res {
				return true
			}
		}
		return false
	case "NOT":
		for _, cond := range group.Conditions {
			res := evalCondition(cond, event)
			fmt.Printf("NOT condition %+v result: %v\n", cond, res)
			if res {
				return false
			}
		}
		for _, sub := range group.Groups {
			res := evalGroup(sub, event)
			fmt.Printf("NOT subgroup %+v result: %v\n", sub, res)
			if res {
				return false
			}
		}
		return true
	default:
		return false
	}
}

func evalCondition(cond data.Condition, event *data.Event) bool {
	fieldVal := getFieldValue(event, cond.Field)
	fmt.Printf("Condition %+v resolved fieldVal: %#v\n", cond, fieldVal)
	if fieldVal == nil {
		return false
	}

	valStr := toString(fieldVal)
	fmt.Printf("Condition %+v resolved valStr: %q\n", cond, valStr)
	switch cond.Operator {
	case "eq", "==":
		return valStr == cond.Value
	case "neq", "!=":
		return valStr != cond.Value
	case "gt":
		return compare(valStr, cond.Value) > 0
	case "lt":
		return compare(valStr, cond.Value) < 0
	case "gte":
		return compare(valStr, cond.Value) >= 0
	case "lte":
		return compare(valStr, cond.Value) <= 0
	case "contains":
		return strings.Contains(valStr, cond.Value)
	case "regex":
		// TODO: implement regex matching
		return false
	default:
		return false
	}
}

func getFieldValue(event *data.Event, path string) interface{} {
	parts := strings.Split(path, ".")
	var current interface{} = event

	fmt.Printf("Resolving path %q with %d parts\n", path, len(parts))
	for i, part := range parts {
		fmt.Printf("Resolving part %d: %q, current type: %T\n", i, part, current)

		// Handle nil values
		if current == nil {
			return nil
		}

		switch obj := current.(type) {
		case *data.Event:
			switch part {
			case "event_id":
				current = obj.ID
			case "event_type":
				current = obj.EventType
			case "event_version":
				current = obj.EventVersion
			case "namespace":
				current = obj.Namespace
			case "object_type":
				current = obj.ObjectType
			case "object_id":
				current = obj.ObjectID
			case "timestamp":
				current = obj.Timestamp
			case "actor":
				current = obj.Actor
			case "context":
				current = obj.Context
			case "payload":
				current = obj.Payload
			case "nats_meta":
				current = obj.NatsMeta
			default:
				return nil
			}
		case map[string]interface{}:
			val, ok := obj[part]
			if !ok {
				return nil
			}
			current = val
		case struct {
			Type string
			ID   string
		}:
			if part == "type" {
				current = obj.Type
			} else if part == "id" {
				current = obj.ID
			} else {
				return nil
			}
		case struct {
			RequestID string
			TraceID   string
		}:
			if part == "request_id" {
				current = obj.RequestID
			} else if part == "trace_id" {
				current = obj.TraceID
			} else {
				return nil
			}
		case struct {
			Stream     string
			Sequence   uint64
			ReceivedAt time.Time
		}:
			if part == "stream" {
				current = obj.Stream
			} else if part == "sequence" {
				current = obj.Sequence
			} else if part == "received_at" {
				current = obj.ReceivedAt
			} else {
				return nil
			}
		default:
			// Special case for Event.Payload
			// Try to access fields by name using reflection
			if part == "before" || part == "after" {
				// Try to access Payload.Before or Payload.After
				fmt.Printf("Trying to access %s field\n", part)

				// Direct field access for Event.Payload
				if i == 1 && parts[0] == "payload" {
					if event.Payload.Before != nil && part == "before" {
						current = event.Payload.Before
						continue
					}
					if event.Payload.After != nil && part == "after" {
						current = event.Payload.After
						continue
					}
				}
			} else if i == 2 && parts[0] == "payload" && (parts[1] == "before" || parts[1] == "after") {
				// Handle payload.before.X or payload.after.X
				if parts[1] == "before" && event.Payload.Before != nil {
					if val, ok := event.Payload.Before[part]; ok {
						current = val
						continue
					}
				} else if parts[1] == "after" && event.Payload.After != nil {
					if val, ok := event.Payload.After[part]; ok {
						current = val
						continue
					}
				}
			}

			// Try reflection for unknown types
			fmt.Printf("Unknown type: %T for part %q\n", current, part)
			return nil
		}
	}
	return current
}

func toString(val interface{}) string {
	switch v := val.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 64)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		if v {
			return "true"
		}
		return "false"
	case time.Time:
		return v.Format(time.RFC3339)
	default:
		return ""
	}
}

func compare(a, b string) int {
	// try numeric comparison
	af, aerr := strconv.ParseFloat(a, 64)
	bf, berr := strconv.ParseFloat(b, 64)
	if aerr == nil && berr == nil {
		if af < bf {
			return -1
		} else if af > bf {
			return 1
		}
		return 0
	}
	// fallback to string comparison
	if a < b {
		return -1
	} else if a > b {
		return 1
	}
	return 0
}
