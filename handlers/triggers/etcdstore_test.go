package triggers

import (
	"context"
	"testing"

	"github.com/julianshen/nebula/data"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// TestEtcdStore_Interface ensures EtcdStore implements TriggerStore
func TestEtcdStore_Interface(t *testing.T) {
	var _ TriggerStore = (*EtcdStore)(nil)
}

// MockEtcdClient is a mock implementation of the etcd client for testing
type MockEtcdClient struct {
	GetFunc    func(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.GetResponse, error)
	PutFunc    func(ctx context.Context, key, val string, opts ...clientv3.OpOption) (*clientv3.PutResponse, error)
	DeleteFunc func(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.DeleteResponse, error)
	WatchFunc  func(ctx context.Context, key string, opts ...clientv3.OpOption) clientv3.WatchChan
}

func (m *MockEtcdClient) Get(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.GetResponse, error) {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, key, opts...)
	}
	return &clientv3.GetResponse{}, nil
}

func (m *MockEtcdClient) Put(ctx context.Context, key, val string, opts ...clientv3.OpOption) (*clientv3.PutResponse, error) {
	if m.PutFunc != nil {
		return m.PutFunc(ctx, key, val, opts...)
	}
	return &clientv3.PutResponse{}, nil
}

func (m *MockEtcdClient) Delete(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.DeleteResponse, error) {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, key, opts...)
	}
	return &clientv3.DeleteResponse{}, nil
}

func (m *MockEtcdClient) Watch(ctx context.Context, key string, opts ...clientv3.OpOption) clientv3.WatchChan {
	if m.WatchFunc != nil {
		return m.WatchFunc(ctx, key, opts...)
	}
	ch := make(chan clientv3.WatchResponse)
	close(ch)
	return ch
}

func (m *MockEtcdClient) Close() error {
	return nil
}

// TestEtcdStore_ParseKey tests the parseKey function
func TestEtcdStore_ParseKey(t *testing.T) {
	store := &EtcdStore{
		prefix: "/triggers/",
	}

	tests := []struct {
		name            string
		key             string
		wantNamespace   string
		wantTriggerName string
		wantErr         bool
	}{
		{
			name:            "valid key",
			key:             "/triggers/namespace1/trigger1.yaml",
			wantNamespace:   "namespace1",
			wantTriggerName: "trigger1",
			wantErr:         false,
		},
		{
			name:            "valid key with nested path",
			key:             "/triggers/namespace1/subdir/trigger2.yaml",
			wantNamespace:   "namespace1",
			wantTriggerName: "trigger2",
			wantErr:         false,
		},
		{
			name:            "invalid key - wrong prefix",
			key:             "/wrong/namespace1/trigger1.yaml",
			wantNamespace:   "",
			wantTriggerName: "",
			wantErr:         true,
		},
		{
			name:            "invalid key - no namespace",
			key:             "/triggers/trigger1.yaml",
			wantNamespace:   "",
			wantTriggerName: "",
			wantErr:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			namespace, triggerName, err := store.parseKey([]byte(tt.key))
			if (err != nil) != tt.wantErr {
				t.Errorf("parseKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if namespace != tt.wantNamespace {
				t.Errorf("parseKey() namespace = %v, want %v", namespace, tt.wantNamespace)
			}
			if triggerName != tt.wantTriggerName {
				t.Errorf("parseKey() triggerName = %v, want %v", triggerName, tt.wantTriggerName)
			}
		})
	}
}

// TestEtcdStore_GetTriggers tests the GetTriggers function
func TestEtcdStore_GetTriggers(t *testing.T) {
	store := &EtcdStore{
		triggers: map[string]map[string]*data.Trigger{
			"namespace1": {
				"trigger1": &data.Trigger{ID: "1", Name: "Trigger 1"},
				"trigger2": &data.Trigger{ID: "2", Name: "Trigger 2"},
			},
			"namespace2": {
				"trigger3": &data.Trigger{ID: "3", Name: "Trigger 3"},
			},
		},
	}

	tests := []struct {
		name      string
		namespace string
		want      int
	}{
		{
			name:      "namespace1",
			namespace: "namespace1",
			want:      2,
		},
		{
			name:      "namespace2",
			namespace: "namespace2",
			want:      1,
		},
		{
			name:      "nonexistent namespace",
			namespace: "nonexistent",
			want:      0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			triggers := store.GetTriggers(tt.namespace)
			if len(triggers) != tt.want {
				t.Errorf("GetTriggers() returned %d triggers, want %d", len(triggers), tt.want)
			}
		})
	}
}

// TestEtcdStore_GetAllTriggers tests the GetAllTriggers function
func TestEtcdStore_GetAllTriggers(t *testing.T) {
	store := &EtcdStore{
		triggers: map[string]map[string]*data.Trigger{
			"namespace1": {
				"trigger1": &data.Trigger{ID: "1", Name: "Trigger 1"},
				"trigger2": &data.Trigger{ID: "2", Name: "Trigger 2"},
			},
			"namespace2": {
				"trigger3": &data.Trigger{ID: "3", Name: "Trigger 3"},
			},
		},
	}

	triggers := store.GetAllTriggers()
	if len(triggers) != 3 {
		t.Errorf("GetAllTriggers() returned %d triggers, want 3", len(triggers))
	}
}
