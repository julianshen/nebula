package server

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/julianshen/nebula/api/gen/api/proto"
	"github.com/julianshen/nebula/data"
	"github.com/julianshen/nebula/handlers/triggers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TriggerServer implements the TriggerService gRPC server
type TriggerServer struct {
	pb.UnimplementedTriggerServiceServer
	store triggers.TriggerStore
}

// NewTriggerServer creates a new TriggerServer
func NewTriggerServer(store triggers.TriggerStore) *TriggerServer {
	return &TriggerServer{
		store: store,
	}
}

// Start starts the gRPC server
func (s *TriggerServer) Start(address string) error {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterTriggerServiceServer(grpcServer, s)

	log.Printf("Starting gRPC server on %s", address)
	return grpcServer.Serve(lis)
}

// ListTriggers lists all triggers under a specified namespace
func (s *TriggerServer) ListTriggers(ctx context.Context, req *pb.ListTriggersRequest) (*pb.ListTriggersResponse, error) {
	triggers := s.store.GetTriggers(req.Namespace)

	pbTriggers := make([]*pb.Trigger, 0, len(triggers))
	for _, t := range triggers {
		pbTriggers = append(pbTriggers, convertToPbTrigger(t))
	}

	return &pb.ListTriggersResponse{
		Triggers: pbTriggers,
	}, nil
}

// AddTrigger adds a new trigger to a namespace
func (s *TriggerServer) AddTrigger(ctx context.Context, req *pb.AddTriggerRequest) (*pb.AddTriggerResponse, error) {
	if req.Trigger == nil {
		return nil, status.Error(codes.InvalidArgument, "trigger is required")
	}

	trigger := convertToDataTrigger(req.Trigger)
	err := s.store.SaveTrigger(ctx, trigger.Namespace, trigger.ID, trigger)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to save trigger: %v", err)
	}

	return &pb.AddTriggerResponse{
		Trigger: req.Trigger,
	}, nil
}

// UpdateTrigger updates an existing trigger
func (s *TriggerServer) UpdateTrigger(ctx context.Context, req *pb.UpdateTriggerRequest) (*pb.UpdateTriggerResponse, error) {
	if req.Trigger == nil {
		return nil, status.Error(codes.InvalidArgument, "trigger is required")
	}

	trigger := convertToDataTrigger(req.Trigger)
	err := s.store.SaveTrigger(ctx, trigger.Namespace, trigger.ID, trigger)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update trigger: %v", err)
	}

	return &pb.UpdateTriggerResponse{
		Trigger: req.Trigger,
	}, nil
}

// RemoveTrigger removes a trigger
func (s *TriggerServer) RemoveTrigger(ctx context.Context, req *pb.RemoveTriggerRequest) (*pb.RemoveTriggerResponse, error) {
	err := s.store.DeleteTrigger(ctx, req.Namespace, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete trigger: %v", err)
	}

	return &pb.RemoveTriggerResponse{
		Success: true,
	}, nil
}

// Helper functions to convert between protobuf and data types

func convertToPbTrigger(t *data.Trigger) *pb.Trigger {
	return &pb.Trigger{
		Id:         t.ID,
		Name:       t.Name,
		Namespace:  t.Namespace,
		ObjectType: t.ObjectType,
		EventType:  t.EventType,
		Enabled:    t.Enabled,
		RootGroup:  convertToPbConditionGroup(&t.RootGroup),
	}
}

func convertToPbConditionGroup(g *data.ConditionGroup) *pb.ConditionGroup {
	pbGroup := &pb.ConditionGroup{
		Operator:   g.Operator,
		Conditions: make([]*pb.Condition, 0, len(g.Conditions)),
		Groups:     make([]*pb.ConditionGroup, 0, len(g.Groups)),
	}

	for _, c := range g.Conditions {
		pbGroup.Conditions = append(pbGroup.Conditions, &pb.Condition{
			Field:    c.Field,
			Operator: c.Operator,
			Value:    c.Value,
		})
	}

	for _, subGroup := range g.Groups {
		pbGroup.Groups = append(pbGroup.Groups, convertToPbConditionGroup(&subGroup))
	}

	return pbGroup
}

func convertToDataTrigger(t *pb.Trigger) *data.Trigger {
	return &data.Trigger{
		ID:         t.Id,
		Name:       t.Name,
		Namespace:  t.Namespace,
		ObjectType: t.ObjectType,
		EventType:  t.EventType,
		Enabled:    t.Enabled,
		RootGroup:  *convertToDataConditionGroup(t.RootGroup),
	}
}

func convertToDataConditionGroup(g *pb.ConditionGroup) *data.ConditionGroup {
	dataGroup := &data.ConditionGroup{
		Operator:   g.Operator,
		Conditions: make([]data.Condition, 0, len(g.Conditions)),
		Groups:     make([]data.ConditionGroup, 0, len(g.Groups)),
	}

	for _, c := range g.Conditions {
		dataGroup.Conditions = append(dataGroup.Conditions, data.Condition{
			Field:    c.Field,
			Operator: c.Operator,
			Value:    c.Value,
		})
	}

	for _, subGroup := range g.Groups {
		dataGroup.Groups = append(dataGroup.Groups, *convertToDataConditionGroup(subGroup))
	}

	return dataGroup
}
