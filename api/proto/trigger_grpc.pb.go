// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.23.4
// source: api/proto/trigger.proto

package api

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// TriggerServiceClient is the client API for TriggerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TriggerServiceClient interface {
	// ListTriggers lists all triggers under a specified namespace
	ListTriggers(ctx context.Context, in *ListTriggersRequest, opts ...grpc.CallOption) (*ListTriggersResponse, error)
	// AddTrigger adds a new trigger to a namespace
	AddTrigger(ctx context.Context, in *AddTriggerRequest, opts ...grpc.CallOption) (*AddTriggerResponse, error)
	// UpdateTrigger updates an existing trigger
	UpdateTrigger(ctx context.Context, in *UpdateTriggerRequest, opts ...grpc.CallOption) (*UpdateTriggerResponse, error)
	// RemoveTrigger removes a trigger
	RemoveTrigger(ctx context.Context, in *RemoveTriggerRequest, opts ...grpc.CallOption) (*RemoveTriggerResponse, error)
}

type triggerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTriggerServiceClient(cc grpc.ClientConnInterface) TriggerServiceClient {
	return &triggerServiceClient{cc}
}

func (c *triggerServiceClient) ListTriggers(ctx context.Context, in *ListTriggersRequest, opts ...grpc.CallOption) (*ListTriggersResponse, error) {
	out := new(ListTriggersResponse)
	err := c.cc.Invoke(ctx, "/api.TriggerService/ListTriggers", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *triggerServiceClient) AddTrigger(ctx context.Context, in *AddTriggerRequest, opts ...grpc.CallOption) (*AddTriggerResponse, error) {
	out := new(AddTriggerResponse)
	err := c.cc.Invoke(ctx, "/api.TriggerService/AddTrigger", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *triggerServiceClient) UpdateTrigger(ctx context.Context, in *UpdateTriggerRequest, opts ...grpc.CallOption) (*UpdateTriggerResponse, error) {
	out := new(UpdateTriggerResponse)
	err := c.cc.Invoke(ctx, "/api.TriggerService/UpdateTrigger", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *triggerServiceClient) RemoveTrigger(ctx context.Context, in *RemoveTriggerRequest, opts ...grpc.CallOption) (*RemoveTriggerResponse, error) {
	out := new(RemoveTriggerResponse)
	err := c.cc.Invoke(ctx, "/api.TriggerService/RemoveTrigger", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TriggerServiceServer is the server API for TriggerService service.
// All implementations must embed UnimplementedTriggerServiceServer
// for forward compatibility
type TriggerServiceServer interface {
	// ListTriggers lists all triggers under a specified namespace
	ListTriggers(context.Context, *ListTriggersRequest) (*ListTriggersResponse, error)
	// AddTrigger adds a new trigger to a namespace
	AddTrigger(context.Context, *AddTriggerRequest) (*AddTriggerResponse, error)
	// UpdateTrigger updates an existing trigger
	UpdateTrigger(context.Context, *UpdateTriggerRequest) (*UpdateTriggerResponse, error)
	// RemoveTrigger removes a trigger
	RemoveTrigger(context.Context, *RemoveTriggerRequest) (*RemoveTriggerResponse, error)
	mustEmbedUnimplementedTriggerServiceServer()
}

// UnimplementedTriggerServiceServer must be embedded to have forward compatible implementations.
type UnimplementedTriggerServiceServer struct {
}

func (UnimplementedTriggerServiceServer) ListTriggers(context.Context, *ListTriggersRequest) (*ListTriggersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListTriggers not implemented")
}
func (UnimplementedTriggerServiceServer) AddTrigger(context.Context, *AddTriggerRequest) (*AddTriggerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddTrigger not implemented")
}
func (UnimplementedTriggerServiceServer) UpdateTrigger(context.Context, *UpdateTriggerRequest) (*UpdateTriggerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateTrigger not implemented")
}
func (UnimplementedTriggerServiceServer) RemoveTrigger(context.Context, *RemoveTriggerRequest) (*RemoveTriggerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveTrigger not implemented")
}
func (UnimplementedTriggerServiceServer) mustEmbedUnimplementedTriggerServiceServer() {}

// UnsafeTriggerServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TriggerServiceServer will
// result in compilation errors.
type UnsafeTriggerServiceServer interface {
	mustEmbedUnimplementedTriggerServiceServer()
}

func RegisterTriggerServiceServer(s grpc.ServiceRegistrar, srv TriggerServiceServer) {
	s.RegisterService(&TriggerService_ServiceDesc, srv)
}

func _TriggerService_ListTriggers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListTriggersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TriggerServiceServer).ListTriggers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.TriggerService/ListTriggers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TriggerServiceServer).ListTriggers(ctx, req.(*ListTriggersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TriggerService_AddTrigger_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddTriggerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TriggerServiceServer).AddTrigger(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.TriggerService/AddTrigger",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TriggerServiceServer).AddTrigger(ctx, req.(*AddTriggerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TriggerService_UpdateTrigger_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateTriggerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TriggerServiceServer).UpdateTrigger(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.TriggerService/UpdateTrigger",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TriggerServiceServer).UpdateTrigger(ctx, req.(*UpdateTriggerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TriggerService_RemoveTrigger_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveTriggerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TriggerServiceServer).RemoveTrigger(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.TriggerService/RemoveTrigger",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TriggerServiceServer).RemoveTrigger(ctx, req.(*RemoveTriggerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// TriggerService_ServiceDesc is the grpc.ServiceDesc for TriggerService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TriggerService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.TriggerService",
	HandlerType: (*TriggerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListTriggers",
			Handler:    _TriggerService_ListTriggers_Handler,
		},
		{
			MethodName: "AddTrigger",
			Handler:    _TriggerService_AddTrigger_Handler,
		},
		{
			MethodName: "UpdateTrigger",
			Handler:    _TriggerService_UpdateTrigger_Handler,
		},
		{
			MethodName: "RemoveTrigger",
			Handler:    _TriggerService_RemoveTrigger_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/proto/trigger.proto",
}
