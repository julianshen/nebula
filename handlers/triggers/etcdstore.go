package triggers

import (
	"bytes"
	"context"
	"fmt"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/julianshen/nebula/data"
	clientv3 "go.etcd.io/etcd/client/v3"
)

const (
	// DefaultTriggerPrefix is the default prefix for trigger keys in etcd
	DefaultTriggerPrefix = "/triggers/"
	// DefaultWatchTimeout is the default timeout for watch operations
	DefaultWatchTimeout = 5 * time.Second
)

// EtcdStore represents a trigger store backed by etcd
type EtcdStore struct {
	TriggerStore
	client      *clientv3.Client
	prefix      string
	triggers    map[string]map[string]*data.Trigger // namespace -> triggerName -> Trigger
	mu          sync.RWMutex
	watchCancel context.CancelFunc
}

// NewEtcdStore creates a new etcd-backed trigger store
func NewEtcdStore(endpoints []string, prefix string) (*EtcdStore, error) {
	if prefix == "" {
		prefix = DefaultTriggerPrefix
	}

	// Ensure prefix ends with "/"
	if !strings.HasSuffix(prefix, "/") {
		prefix = prefix + "/"
	}

	// Create etcd client
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create etcd client: %w", err)
	}

	return &EtcdStore{
		client:   client,
		prefix:   prefix,
		triggers: make(map[string]map[string]*data.Trigger),
	}, nil
}

// Close closes the etcd client and stops watching for changes
func (s *EtcdStore) Close() error {
	if s.watchCancel != nil {
		s.watchCancel()
	}
	return s.client.Close()
}

// LoadAll loads all triggers from etcd
func (s *EtcdStore) LoadAll(ctx context.Context) error {
	// Get all keys under the prefix
	resp, err := s.client.Get(ctx, s.prefix, clientv3.WithPrefix())
	if err != nil {
		return fmt.Errorf("failed to get triggers from etcd: %w", err)
	}

	// Clear existing triggers
	s.mu.Lock()
	s.triggers = make(map[string]map[string]*data.Trigger)
	s.mu.Unlock()

	// Process each key-value pair
	for _, kv := range resp.Kvs {
		if err := s.processTrigger(kv.Key, kv.Value); err != nil {
			return err
		}
	}

	return nil
}

// Watch starts watching for changes to triggers in etcd
func (s *EtcdStore) Watch(ctx context.Context) {
	// Cancel any existing watch
	if s.watchCancel != nil {
		s.watchCancel()
	}

	// Create a new context with cancel function
	watchCtx, cancel := context.WithCancel(ctx)
	s.watchCancel = cancel

	// Start watching for changes
	watchChan := s.client.Watch(watchCtx, s.prefix, clientv3.WithPrefix())

	// Process watch events in a goroutine
	go func() {
		for watchResp := range watchChan {
			for _, event := range watchResp.Events {
				switch event.Type {
				case clientv3.EventTypePut:
					// Process updated or new trigger
					if err := s.processTrigger(event.Kv.Key, event.Kv.Value); err != nil {
						fmt.Printf("Error processing trigger update: %v\n", err)
					}
				case clientv3.EventTypeDelete:
					// Remove deleted trigger
					if err := s.removeTrigger(event.Kv.Key); err != nil {
						fmt.Printf("Error removing trigger: %v\n", err)
					}
				}
			}
		}
	}()
}

// GetTriggers returns all triggers for a namespace
func (s *EtcdStore) GetTriggers(namespace string) []*data.Trigger {
	s.mu.RLock()
	defer s.mu.RUnlock()

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

// GetAllTriggers returns all triggers from all namespaces
func (s *EtcdStore) GetAllTriggers() []*data.Trigger {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var allTriggers []*data.Trigger
	for _, namespaceTriggers := range s.triggers {
		for _, trigger := range namespaceTriggers {
			allTriggers = append(allTriggers, trigger)
		}
	}

	return allTriggers
}

// processTrigger processes a trigger key-value pair from etcd
func (s *EtcdStore) processTrigger(key, value []byte) error {
	// Extract namespace and trigger name from key
	namespace, triggerName, err := s.parseKey(key)
	if err != nil {
		return err
	}

	// Parse trigger YAML
	reader := bytes.NewReader(value)
	trigger, err := LoadTrigger(reader)
	if err != nil {
		return fmt.Errorf("failed to parse trigger %s/%s: %w", namespace, triggerName, err)
	}

	// Store trigger in memory
	s.mu.Lock()
	defer s.mu.Unlock()

	// Initialize namespace map if it doesn't exist
	if _, ok := s.triggers[namespace]; !ok {
		s.triggers[namespace] = make(map[string]*data.Trigger)
	}

	// Store trigger
	s.triggers[namespace][triggerName] = trigger

	return nil
}

// removeTrigger removes a trigger from memory
func (s *EtcdStore) removeTrigger(key []byte) error {
	// Extract namespace and trigger name from key
	namespace, triggerName, err := s.parseKey(key)
	if err != nil {
		return err
	}

	// Remove trigger from memory
	s.mu.Lock()
	defer s.mu.Unlock()

	if namespaceTriggers, ok := s.triggers[namespace]; ok {
		delete(namespaceTriggers, triggerName)
	}

	return nil
}

// parseKey parses a key into namespace and trigger name
func (s *EtcdStore) parseKey(key []byte) (namespace, triggerName string, err error) {
	// Convert key to string and remove prefix
	keyStr := string(key)
	if !strings.HasPrefix(keyStr, s.prefix) {
		return "", "", fmt.Errorf("key %s does not have expected prefix %s", keyStr, s.prefix)
	}

	// Remove prefix
	keyStr = strings.TrimPrefix(keyStr, s.prefix)

	// Split into namespace and trigger name
	parts := strings.Split(keyStr, "/")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("invalid key format: %s", keyStr)
	}

	namespace = parts[0]
	triggerName = path.Base(keyStr)

	// Remove file extension if present
	if ext := path.Ext(triggerName); ext != "" {
		triggerName = strings.TrimSuffix(triggerName, ext)
	}

	return namespace, triggerName, nil
}

// SaveTrigger saves a trigger to etcd
func (s *EtcdStore) SaveTrigger(ctx context.Context, namespace, name string, trigger *data.Trigger) error {
	// Marshal trigger to YAML
	yamlData, err := trigger.ToYAML()
	if err != nil {
		return fmt.Errorf("failed to marshal trigger to YAML: %w", err)
	}

	// Create key
	key := s.prefix + namespace + "/" + name + ".yaml"

	// Save to etcd
	_, err = s.client.Put(ctx, key, string(yamlData))
	if err != nil {
		return fmt.Errorf("failed to save trigger to etcd: %w", err)
	}

	return nil
}

// DeleteTrigger deletes a trigger from etcd
func (s *EtcdStore) DeleteTrigger(ctx context.Context, namespace, name string) error {
	// Create key
	key := s.prefix + namespace + "/" + name + ".yaml"

	// Delete from etcd
	_, err := s.client.Delete(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to delete trigger from etcd: %w", err)
	}

	return nil
}
