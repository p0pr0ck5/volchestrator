// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package svc

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// VolchestratorAdminClient is the client API for VolchestratorAdmin service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type VolchestratorAdminClient interface {
	GetClient(ctx context.Context, in *GetClientRequest, opts ...grpc.CallOption) (*GetClientResponse, error)
	ListClients(ctx context.Context, in *ListClientsRequest, opts ...grpc.CallOption) (*ListClientsResponse, error)
	GetVolume(ctx context.Context, in *GetVolumeRequest, opts ...grpc.CallOption) (*GetVolumeResponse, error)
	AddVolume(ctx context.Context, in *AddVolumeRequest, opts ...grpc.CallOption) (*AddVolumeResponse, error)
	UpdateVolume(ctx context.Context, in *UpdateVolumeRequest, opts ...grpc.CallOption) (*UpdateVolumeResponse, error)
	DeleteVolume(ctx context.Context, in *DeleteVolumeRequest, opts ...grpc.CallOption) (*DeleteVolumeResponse, error)
}

type volchestratorAdminClient struct {
	cc grpc.ClientConnInterface
}

func NewVolchestratorAdminClient(cc grpc.ClientConnInterface) VolchestratorAdminClient {
	return &volchestratorAdminClient{cc}
}

func (c *volchestratorAdminClient) GetClient(ctx context.Context, in *GetClientRequest, opts ...grpc.CallOption) (*GetClientResponse, error) {
	out := new(GetClientResponse)
	err := c.cc.Invoke(ctx, "/volchestrator.VolchestratorAdmin/GetClient", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *volchestratorAdminClient) ListClients(ctx context.Context, in *ListClientsRequest, opts ...grpc.CallOption) (*ListClientsResponse, error) {
	out := new(ListClientsResponse)
	err := c.cc.Invoke(ctx, "/volchestrator.VolchestratorAdmin/ListClients", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *volchestratorAdminClient) GetVolume(ctx context.Context, in *GetVolumeRequest, opts ...grpc.CallOption) (*GetVolumeResponse, error) {
	out := new(GetVolumeResponse)
	err := c.cc.Invoke(ctx, "/volchestrator.VolchestratorAdmin/GetVolume", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *volchestratorAdminClient) AddVolume(ctx context.Context, in *AddVolumeRequest, opts ...grpc.CallOption) (*AddVolumeResponse, error) {
	out := new(AddVolumeResponse)
	err := c.cc.Invoke(ctx, "/volchestrator.VolchestratorAdmin/AddVolume", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *volchestratorAdminClient) UpdateVolume(ctx context.Context, in *UpdateVolumeRequest, opts ...grpc.CallOption) (*UpdateVolumeResponse, error) {
	out := new(UpdateVolumeResponse)
	err := c.cc.Invoke(ctx, "/volchestrator.VolchestratorAdmin/UpdateVolume", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *volchestratorAdminClient) DeleteVolume(ctx context.Context, in *DeleteVolumeRequest, opts ...grpc.CallOption) (*DeleteVolumeResponse, error) {
	out := new(DeleteVolumeResponse)
	err := c.cc.Invoke(ctx, "/volchestrator.VolchestratorAdmin/DeleteVolume", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// VolchestratorAdminServer is the server API for VolchestratorAdmin service.
// All implementations must embed UnimplementedVolchestratorAdminServer
// for forward compatibility
type VolchestratorAdminServer interface {
	GetClient(context.Context, *GetClientRequest) (*GetClientResponse, error)
	ListClients(context.Context, *ListClientsRequest) (*ListClientsResponse, error)
	GetVolume(context.Context, *GetVolumeRequest) (*GetVolumeResponse, error)
	AddVolume(context.Context, *AddVolumeRequest) (*AddVolumeResponse, error)
	UpdateVolume(context.Context, *UpdateVolumeRequest) (*UpdateVolumeResponse, error)
	DeleteVolume(context.Context, *DeleteVolumeRequest) (*DeleteVolumeResponse, error)
	mustEmbedUnimplementedVolchestratorAdminServer()
}

// UnimplementedVolchestratorAdminServer must be embedded to have forward compatible implementations.
type UnimplementedVolchestratorAdminServer struct {
}

func (UnimplementedVolchestratorAdminServer) GetClient(context.Context, *GetClientRequest) (*GetClientResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetClient not implemented")
}
func (UnimplementedVolchestratorAdminServer) ListClients(context.Context, *ListClientsRequest) (*ListClientsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListClients not implemented")
}
func (UnimplementedVolchestratorAdminServer) GetVolume(context.Context, *GetVolumeRequest) (*GetVolumeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetVolume not implemented")
}
func (UnimplementedVolchestratorAdminServer) AddVolume(context.Context, *AddVolumeRequest) (*AddVolumeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddVolume not implemented")
}
func (UnimplementedVolchestratorAdminServer) UpdateVolume(context.Context, *UpdateVolumeRequest) (*UpdateVolumeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateVolume not implemented")
}
func (UnimplementedVolchestratorAdminServer) DeleteVolume(context.Context, *DeleteVolumeRequest) (*DeleteVolumeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteVolume not implemented")
}
func (UnimplementedVolchestratorAdminServer) mustEmbedUnimplementedVolchestratorAdminServer() {}

// UnsafeVolchestratorAdminServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to VolchestratorAdminServer will
// result in compilation errors.
type UnsafeVolchestratorAdminServer interface {
	mustEmbedUnimplementedVolchestratorAdminServer()
}

func RegisterVolchestratorAdminServer(s grpc.ServiceRegistrar, srv VolchestratorAdminServer) {
	s.RegisterService(&VolchestratorAdmin_ServiceDesc, srv)
}

func _VolchestratorAdmin_GetClient_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetClientRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VolchestratorAdminServer).GetClient(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/volchestrator.VolchestratorAdmin/GetClient",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VolchestratorAdminServer).GetClient(ctx, req.(*GetClientRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VolchestratorAdmin_ListClients_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListClientsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VolchestratorAdminServer).ListClients(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/volchestrator.VolchestratorAdmin/ListClients",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VolchestratorAdminServer).ListClients(ctx, req.(*ListClientsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VolchestratorAdmin_GetVolume_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetVolumeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VolchestratorAdminServer).GetVolume(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/volchestrator.VolchestratorAdmin/GetVolume",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VolchestratorAdminServer).GetVolume(ctx, req.(*GetVolumeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VolchestratorAdmin_AddVolume_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddVolumeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VolchestratorAdminServer).AddVolume(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/volchestrator.VolchestratorAdmin/AddVolume",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VolchestratorAdminServer).AddVolume(ctx, req.(*AddVolumeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VolchestratorAdmin_UpdateVolume_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateVolumeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VolchestratorAdminServer).UpdateVolume(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/volchestrator.VolchestratorAdmin/UpdateVolume",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VolchestratorAdminServer).UpdateVolume(ctx, req.(*UpdateVolumeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VolchestratorAdmin_DeleteVolume_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteVolumeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VolchestratorAdminServer).DeleteVolume(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/volchestrator.VolchestratorAdmin/DeleteVolume",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VolchestratorAdminServer).DeleteVolume(ctx, req.(*DeleteVolumeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// VolchestratorAdmin_ServiceDesc is the grpc.ServiceDesc for VolchestratorAdmin service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var VolchestratorAdmin_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "volchestrator.VolchestratorAdmin",
	HandlerType: (*VolchestratorAdminServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetClient",
			Handler:    _VolchestratorAdmin_GetClient_Handler,
		},
		{
			MethodName: "ListClients",
			Handler:    _VolchestratorAdmin_ListClients_Handler,
		},
		{
			MethodName: "GetVolume",
			Handler:    _VolchestratorAdmin_GetVolume_Handler,
		},
		{
			MethodName: "AddVolume",
			Handler:    _VolchestratorAdmin_AddVolume_Handler,
		},
		{
			MethodName: "UpdateVolume",
			Handler:    _VolchestratorAdmin_UpdateVolume_Handler,
		},
		{
			MethodName: "DeleteVolume",
			Handler:    _VolchestratorAdmin_DeleteVolume_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "svc/volchestrator_admin.proto",
}
