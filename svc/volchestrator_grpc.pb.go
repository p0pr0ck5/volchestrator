// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package volchestrator

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// VolchestratorClient is the client API for Volchestrator service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type VolchestratorClient interface {
	Heartbeat(ctx context.Context, in *HeartbeatMessage, opts ...grpc.CallOption) (*HeartbeatResponse, error)
}

type volchestratorClient struct {
	cc grpc.ClientConnInterface
}

func NewVolchestratorClient(cc grpc.ClientConnInterface) VolchestratorClient {
	return &volchestratorClient{cc}
}

func (c *volchestratorClient) Heartbeat(ctx context.Context, in *HeartbeatMessage, opts ...grpc.CallOption) (*HeartbeatResponse, error) {
	out := new(HeartbeatResponse)
	err := c.cc.Invoke(ctx, "/volchestrator.Volchestrator/Heartbeat", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// VolchestratorServer is the server API for Volchestrator service.
// All implementations must embed UnimplementedVolchestratorServer
// for forward compatibility
type VolchestratorServer interface {
	Heartbeat(context.Context, *HeartbeatMessage) (*HeartbeatResponse, error)
	mustEmbedUnimplementedVolchestratorServer()
}

// UnimplementedVolchestratorServer must be embedded to have forward compatible implementations.
type UnimplementedVolchestratorServer struct {
}

func (UnimplementedVolchestratorServer) Heartbeat(context.Context, *HeartbeatMessage) (*HeartbeatResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Heartbeat not implemented")
}
func (UnimplementedVolchestratorServer) mustEmbedUnimplementedVolchestratorServer() {}

// UnsafeVolchestratorServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to VolchestratorServer will
// result in compilation errors.
type UnsafeVolchestratorServer interface {
	mustEmbedUnimplementedVolchestratorServer()
}

func RegisterVolchestratorServer(s grpc.ServiceRegistrar, srv VolchestratorServer) {
	s.RegisterService(&Volchestrator_ServiceDesc, srv)
}

func _Volchestrator_Heartbeat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HeartbeatMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VolchestratorServer).Heartbeat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/volchestrator.Volchestrator/Heartbeat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VolchestratorServer).Heartbeat(ctx, req.(*HeartbeatMessage))
	}
	return interceptor(ctx, in, info, handler)
}

// Volchestrator_ServiceDesc is the grpc.ServiceDesc for Volchestrator service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Volchestrator_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "volchestrator.Volchestrator",
	HandlerType: (*VolchestratorServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Heartbeat",
			Handler:    _Volchestrator_Heartbeat_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "svc/volchestrator.proto",
}
