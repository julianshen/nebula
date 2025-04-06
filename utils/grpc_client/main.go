package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	pb "github.com/julianshen/nebula/api/gen/api/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Parse command line flags
	var (
		serverAddr = flag.String("server", "localhost:50051", "The server address in the format host:port")
		command    = flag.String("cmd", "list", "Command to execute: list, add, update, remove")
		namespace  = flag.String("namespace", "sales", "Namespace for triggers")
		id         = flag.String("id", "", "Trigger ID (required for update and remove)")
		name       = flag.String("name", "", "Trigger name (required for add and update)")
		objectType = flag.String("object-type", "order", "Object type (for add and update)")
		eventType  = flag.String("event-type", "created", "Event type (for add and update)")
		field1     = flag.String("field1", "payload.after.amount", "First condition field")
		op1        = flag.String("op1", "gt", "First condition operator")
		value1     = flag.String("value1", "1000", "First condition value")
		field2     = flag.String("field2", "payload.after.region", "Second condition field")
		op2        = flag.String("op2", "eq", "Second condition operator")
		value2     = flag.String("value2", "US", "Second condition value")
	)

	flag.Parse()

	// Set up a connection to the server
	conn, err := grpc.Dial(*serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewTriggerServiceClient(conn)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Execute the requested command
	switch *command {
	case "list":
		listTriggers(ctx, client, *namespace)
	case "add":
		if *name == "" {
			log.Fatal("Trigger name is required for add command")
		}
		addTrigger(ctx, client, *namespace, *id, *name, *objectType, *eventType, *field1, *op1, *value1, *field2, *op2, *value2)
	case "update":
		if *id == "" || *name == "" {
			log.Fatal("Trigger ID and name are required for update command")
		}
		updateTrigger(ctx, client, *namespace, *id, *name, *objectType, *eventType, *field1, *op1, *value1, *field2, *op2, *value2)
	case "remove":
		if *id == "" {
			log.Fatal("Trigger ID is required for remove command")
		}
		removeTrigger(ctx, client, *namespace, *id)
	default:
		log.Fatalf("Unknown command: %s", *command)
	}
}

// listTriggers lists all triggers in a namespace
func listTriggers(ctx context.Context, client pb.TriggerServiceClient, namespace string) {
	resp, err := client.ListTriggers(ctx, &pb.ListTriggersRequest{
		Namespace: namespace,
	})
	if err != nil {
		log.Fatalf("Failed to list triggers: %v", err)
	}

	fmt.Printf("Triggers in namespace '%s':\n", namespace)
	if len(resp.Triggers) == 0 {
		fmt.Println("No triggers found")
		return
	}

	for i, trigger := range resp.Triggers {
		fmt.Printf("%d. %s - %s\n", i+1, trigger.Id, trigger.Name)
		fmt.Printf("   Namespace: %s\n", trigger.Namespace)
		fmt.Printf("   Object Type: %s\n", trigger.ObjectType)
		fmt.Printf("   Event Type: %s\n", trigger.EventType)
		fmt.Printf("   Enabled: %v\n", trigger.Enabled)
		fmt.Printf("   Conditions:\n")
		printConditionGroup(trigger.RootGroup, "     ")
		fmt.Println()
	}
}

// addTrigger adds a new trigger
func addTrigger(ctx context.Context, client pb.TriggerServiceClient, namespace, id, name, objectType, eventType, field1, op1, value1, field2, op2, value2 string) {
	// Generate ID if not provided
	if id == "" {
		id = fmt.Sprintf("%s-%d", name, time.Now().Unix())
	}

	trigger := createTrigger(namespace, id, name, objectType, eventType, field1, op1, value1, field2, op2, value2)

	resp, err := client.AddTrigger(ctx, &pb.AddTriggerRequest{
		Trigger: trigger,
	})
	if err != nil {
		log.Fatalf("Failed to add trigger: %v", err)
	}

	fmt.Printf("Successfully added trigger: %s - %s\n", resp.Trigger.Id, resp.Trigger.Name)
}

// updateTrigger updates an existing trigger
func updateTrigger(ctx context.Context, client pb.TriggerServiceClient, namespace, id, name, objectType, eventType, field1, op1, value1, field2, op2, value2 string) {
	trigger := createTrigger(namespace, id, name, objectType, eventType, field1, op1, value1, field2, op2, value2)

	resp, err := client.UpdateTrigger(ctx, &pb.UpdateTriggerRequest{
		Trigger: trigger,
	})
	if err != nil {
		log.Fatalf("Failed to update trigger: %v", err)
	}

	fmt.Printf("Successfully updated trigger: %s - %s\n", resp.Trigger.Id, resp.Trigger.Name)
}

// removeTrigger removes a trigger
func removeTrigger(ctx context.Context, client pb.TriggerServiceClient, namespace, id string) {
	resp, err := client.RemoveTrigger(ctx, &pb.RemoveTriggerRequest{
		Namespace: namespace,
		Id:        id,
	})
	if err != nil {
		log.Fatalf("Failed to remove trigger: %v", err)
	}

	if resp.Success {
		fmt.Printf("Successfully removed trigger: %s\n", id)
	} else {
		fmt.Printf("Failed to remove trigger: %s\n", id)
	}
}

// createTrigger creates a trigger with the specified parameters
func createTrigger(namespace, id, name, objectType, eventType, field1, op1, value1, field2, op2, value2 string) *pb.Trigger {
	return &pb.Trigger{
		Id:         id,
		Name:       name,
		Namespace:  namespace,
		ObjectType: objectType,
		EventType:  eventType,
		Enabled:    true,
		RootGroup: &pb.ConditionGroup{
			Operator: "AND",
			Conditions: []*pb.Condition{
				{
					Field:    field1,
					Operator: op1,
					Value:    value1,
				},
				{
					Field:    field2,
					Operator: op2,
					Value:    value2,
				},
			},
		},
	}
}

// printConditionGroup prints a condition group recursively
func printConditionGroup(group *pb.ConditionGroup, indent string) {
	fmt.Printf("%sOperator: %s\n", indent, group.Operator)

	for _, condition := range group.Conditions {
		fmt.Printf("%s- %s %s %s\n", indent, condition.Field, condition.Operator, condition.Value)
	}

	for i, subGroup := range group.Groups {
		fmt.Printf("%sSubgroup %d:\n", indent, i+1)
		printConditionGroup(subGroup, indent+"  ")
	}
}
