package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/julianshen/nebula/api/server"
	"github.com/julianshen/nebula/data"
	"github.com/julianshen/nebula/handlers/triggers"
	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"
)

func main() {
	// Load config
	cfgFile := "config.yaml"
	viper.SetConfigFile(cfgFile)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize etcd store
	etcdEndpoints := viper.GetStringSlice("etcd.endpoints")
	if len(etcdEndpoints) == 0 {
		etcdEndpoints = []string{"localhost:2379"}
	}

	triggerPrefix := viper.GetString("etcd.trigger_prefix")
	store, err := triggers.NewEtcdStore(etcdEndpoints, triggerPrefix)
	if err != nil {
		log.Fatalf("Failed to create etcd store: %v", err)
	}
	defer store.Close()

	// Load all triggers from etcd
	if err := store.LoadAll(ctx); err != nil {
		log.Printf("Warning: Failed to load triggers from etcd: %v", err)
	} else {
		log.Println("Successfully loaded triggers from etcd")
	}

	// Start watching for changes
	store.Watch(ctx)
	log.Println("Watching for trigger changes in etcd")

	// Start gRPC server
	grpcAddress := viper.GetString("grpc.address")
	if grpcAddress == "" {
		grpcAddress = ":50051"
	}

	triggerServer := server.NewTriggerServer(store)
	go func() {
		if err := triggerServer.Start(grpcAddress); err != nil {
			log.Fatalf("Failed to start gRPC server: %v", err)
		}
	}()
	log.Printf("gRPC server started on %s", grpcAddress)

	// Connect to NATS
	natsURL := viper.GetString("nats.url")
	if natsURL == "" {
		natsURL = nats.DefaultURL
	}
	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	// Subscribe to subject
	subject := viper.GetString("nats.subject")
	if subject == "" {
		subject = "event.>"
	}

	queueGroup := viper.GetString("nats.queue_group")

	// Create subscription
	var sub *nats.Subscription
	if queueGroup != "" {
		sub, err = nc.QueueSubscribe(subject, queueGroup, func(msg *nats.Msg) {
			handleMessage(msg, store)
		})
	} else {
		sub, err = nc.Subscribe(subject, func(msg *nats.Msg) {
			handleMessage(msg, store)
		})
	}

	if err != nil {
		log.Fatal(err)
	}
	defer sub.Unsubscribe()

	log.Println("Listening for events on:", subject)

	// Wait for termination signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	log.Println("Shutting down...")
}

// Cache of triggers by namespace
var (
	triggerCache     = make(map[string][]*data.Trigger)
	triggerCacheMu   sync.RWMutex
	lastCacheRefresh time.Time
)

// refreshTriggerCache refreshes the trigger cache if needed
func refreshTriggerCache(store *triggers.EtcdStore) {
	triggerCacheMu.RLock()
	// Only refresh if it's been more than 5 seconds since the last refresh
	if time.Since(lastCacheRefresh) < 5*time.Second {
		triggerCacheMu.RUnlock()
		return
	}
	triggerCacheMu.RUnlock()

	// Acquire write lock to update cache
	triggerCacheMu.Lock()
	defer triggerCacheMu.Unlock()

	// Check again in case another goroutine refreshed while we were waiting
	if time.Since(lastCacheRefresh) < 5*time.Second {
		return
	}

	// Get all namespaces
	allTriggers := store.GetAllTriggers()

	// Group triggers by namespace
	newCache := make(map[string][]*data.Trigger)
	for _, trigger := range allTriggers {
		newCache[trigger.Namespace] = append(newCache[trigger.Namespace], trigger)
	}

	// Update cache
	triggerCache = newCache
	lastCacheRefresh = time.Now()

	log.Printf("Trigger cache refreshed with %d namespaces", len(triggerCache))
}

// getTriggersByNamespace gets triggers for a namespace from the cache
func getTriggersByNamespace(namespace string, store *triggers.EtcdStore) []*data.Trigger {
	// Refresh cache if needed
	refreshTriggerCache(store)

	// Get triggers from cache
	triggerCacheMu.RLock()
	defer triggerCacheMu.RUnlock()

	return triggerCache[namespace]
}

// handleMessage processes an incoming NATS message
func handleMessage(msg *nats.Msg, store *triggers.EtcdStore) {
	// Parse event
	var event data.Event
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		log.Printf("Error parsing event: %v", err)
		return
	}

	log.Printf("Received event: %s, type: %s, namespace: %s",
		event.ID, event.EventType, event.Namespace)

	// Get triggers for this namespace from cache
	namespaceTriggers := getTriggersByNamespace(event.Namespace, store)
	log.Printf("Evaluating %d triggers for namespace %s", len(namespaceTriggers), event.Namespace)

	// Evaluate each trigger
	for _, trigger := range namespaceTriggers {
		// Skip disabled triggers
		if !trigger.Enabled {
			continue
		}

		// Skip triggers that don't match the event type or object type
		if trigger.EventType != "" && trigger.EventType != event.EventType {
			continue
		}
		if trigger.ObjectType != "" && trigger.ObjectType != event.ObjectType {
			continue
		}

		// Evaluate trigger conditions
		if triggers.MatchTrigger(trigger, &event) {
			log.Printf("Trigger matched: %s - %s", trigger.ID, trigger.Name)
			// TODO: Execute trigger action
		}
	}
}
