package triggers

import (
	"context"

	"github.com/julianshen/nebula/data"
)

// TriggerStore defines the interface for a trigger store
type TriggerStore interface {
	// LoadAll loads all triggers from the store
	LoadAll(ctx context.Context) error

	// Watch starts watching for changes to triggers
	Watch(ctx context.Context)

	// GetTriggers returns all triggers for a namespace
	GetTriggers(namespace string) []*data.Trigger

	// GetAllTriggers returns all triggers from all namespaces
	GetAllTriggers() []*data.Trigger

	// SaveTrigger saves a trigger to the store
	SaveTrigger(ctx context.Context, namespace, name string, trigger *data.Trigger) error

	// DeleteTrigger deletes a trigger from the store
	DeleteTrigger(ctx context.Context, namespace, name string) error

	// Close closes the store
	Close() error
}
