// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: restore.proto

package pb

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

// RestoreClient is the client API for Restore service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RestoreClient interface {
	Create(ctx context.Context, in *RestoreCreateRequest, opts ...grpc.CallOption) (*RestoreCreateResponse, error)
	Get(ctx context.Context, in *RestoreGetRequest, opts ...grpc.CallOption) (*RestoreGetResponse, error)
	List(ctx context.Context, in *RestoreListRequest, opts ...grpc.CallOption) (*RestoreListResponse, error)
}

type restoreClient struct {
	cc grpc.ClientConnInterface
}

func NewRestoreClient(cc grpc.ClientConnInterface) RestoreClient {
	return &restoreClient{cc}
}

func (c *restoreClient) Create(ctx context.Context, in *RestoreCreateRequest, opts ...grpc.CallOption) (*RestoreCreateResponse, error) {
	out := new(RestoreCreateResponse)
	err := c.cc.Invoke(ctx, "/workflow.restore/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *restoreClient) Get(ctx context.Context, in *RestoreGetRequest, opts ...grpc.CallOption) (*RestoreGetResponse, error) {
	out := new(RestoreGetResponse)
	err := c.cc.Invoke(ctx, "/workflow.restore/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *restoreClient) List(ctx context.Context, in *RestoreListRequest, opts ...grpc.CallOption) (*RestoreListResponse, error) {
	out := new(RestoreListResponse)
	err := c.cc.Invoke(ctx, "/workflow.restore/List", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RestoreServer is the server API for Restore service.
// All implementations must embed UnimplementedRestoreServer
// for forward compatibility
type RestoreServer interface {
	Create(context.Context, *RestoreCreateRequest) (*RestoreCreateResponse, error)
	Get(context.Context, *RestoreGetRequest) (*RestoreGetResponse, error)
	List(context.Context, *RestoreListRequest) (*RestoreListResponse, error)
	mustEmbedUnimplementedRestoreServer()
}

// UnimplementedRestoreServer must be embedded to have forward compatible implementations.
type UnimplementedRestoreServer struct {
}

func (UnimplementedRestoreServer) Create(context.Context, *RestoreCreateRequest) (*RestoreCreateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedRestoreServer) Get(context.Context, *RestoreGetRequest) (*RestoreGetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedRestoreServer) List(context.Context, *RestoreListRequest) (*RestoreListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (UnimplementedRestoreServer) mustEmbedUnimplementedRestoreServer() {}

// UnsafeRestoreServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RestoreServer will
// result in compilation errors.
type UnsafeRestoreServer interface {
	mustEmbedUnimplementedRestoreServer()
}

func RegisterRestoreServer(s grpc.ServiceRegistrar, srv RestoreServer) {
	s.RegisterService(&Restore_ServiceDesc, srv)
}

func _Restore_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RestoreCreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RestoreServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/workflow.restore/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RestoreServer).Create(ctx, req.(*RestoreCreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Restore_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RestoreGetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RestoreServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/workflow.restore/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RestoreServer).Get(ctx, req.(*RestoreGetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Restore_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RestoreListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RestoreServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/workflow.restore/List",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RestoreServer).List(ctx, req.(*RestoreListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Restore_ServiceDesc is the grpc.ServiceDesc for Restore service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Restore_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "workflow.restore",
	HandlerType: (*RestoreServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _Restore_Create_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _Restore_Get_Handler,
		},
		{
			MethodName: "List",
			Handler:    _Restore_List_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "restore.proto",
}
