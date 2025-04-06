package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/julianshen/nebula/data"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// This utility checks MongoDB for stored events

func main() {
	// Parse command line flags
	var (
		mongoURI  = flag.String("mongo", "mongodb://localhost:27017", "MongoDB URI")
		database  = flag.String("db", "eventstore", "MongoDB database name")
		namespace = flag.String("namespace", "sales", "Namespace to filter events")
		limit     = flag.Int("limit", 10, "Maximum number of events to retrieve")
	)

	flag.Parse()

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(*mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	// Get the events collection
	collection := client.Database(*database).Collection("events")

	// Create a filter for the namespace
	filter := bson.M{"namespace": *namespace}

	// Find the events
	findOptions := options.Find().SetLimit(int64(*limit)).SetSort(bson.M{"timestamp": -1})
	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		log.Fatalf("Failed to query MongoDB: %v", err)
	}
	defer cursor.Close(ctx)

	// Iterate through the results
	fmt.Printf("Recent events in namespace '%s':\n\n", *namespace)
	count := 0
	for cursor.Next(ctx) {
		var event data.Event
		if err := cursor.Decode(&event); err != nil {
			log.Printf("Failed to decode event: %v", err)
			continue
		}

		// Print the event details
		count++
		fmt.Printf("Event %d:\n", count)
		fmt.Printf("  ID: %s\n", event.ID)
		fmt.Printf("  Type: %s\n", event.EventType)
		fmt.Printf("  Object: %s/%s\n", event.ObjectType, event.ObjectID)
		fmt.Printf("  Timestamp: %s\n", event.Timestamp.Format(time.RFC3339))

		// Print the payload
		if len(event.Payload.After) > 0 {
			fmt.Println("  Payload:")
			payloadJSON, err := json.MarshalIndent(event.Payload.After, "    ", "  ")
			if err != nil {
				fmt.Printf("    Error marshaling payload: %v\n", err)
			} else {
				fmt.Printf("    %s\n", string(payloadJSON))
			}
		}
		fmt.Println()
	}

	if count == 0 {
		fmt.Printf("No events found in namespace '%s'\n", *namespace)
	} else {
		fmt.Printf("Found %d events\n", count)
	}
}
