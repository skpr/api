// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: cron.proto

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

// CronClient is the client API for Cron service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CronClient interface {
	Suspend(ctx context.Context, in *CronSuspendRequest, opts ...grpc.CallOption) (*CronSuspendResponse, error)
	Resume(ctx context.Context, in *CronResumeRequest, opts ...grpc.CallOption) (*CronResumeResponse, error)
	List(ctx context.Context, in *CronListRequest, opts ...grpc.CallOption) (*CronListResponse, error)
	JobList(ctx context.Context, in *CronJobListRequest, opts ...grpc.CallOption) (*CronJobListResponse, error)
}

type cronClient struct {
	cc grpc.ClientConnInterface
}

func NewCronClient(cc grpc.ClientConnInterface) CronClient {
	return &cronClient{cc}
}

func (c *cronClient) Suspend(ctx context.Context, in *CronSuspendRequest, opts ...grpc.CallOption) (*CronSuspendResponse, error) {
	out := new(CronSuspendResponse)
	err := c.cc.Invoke(ctx, "/workflow.cron/Suspend", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cronClient) Resume(ctx context.Context, in *CronResumeRequest, opts ...grpc.CallOption) (*CronResumeResponse, error) {
	out := new(CronResumeResponse)
	err := c.cc.Invoke(ctx, "/workflow.cron/Resume", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cronClient) List(ctx context.Context, in *CronListRequest, opts ...grpc.CallOption) (*CronListResponse, error) {
	out := new(CronListResponse)
	err := c.cc.Invoke(ctx, "/workflow.cron/List", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cronClient) JobList(ctx context.Context, in *CronJobListRequest, opts ...grpc.CallOption) (*CronJobListResponse, error) {
	out := new(CronJobListResponse)
	err := c.cc.Invoke(ctx, "/workflow.cron/JobList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CronServer is the server API for Cron service.
// All implementations must embed UnimplementedCronServer
// for forward compatibility
type CronServer interface {
	Suspend(context.Context, *CronSuspendRequest) (*CronSuspendResponse, error)
	Resume(context.Context, *CronResumeRequest) (*CronResumeResponse, error)
	List(context.Context, *CronListRequest) (*CronListResponse, error)
	JobList(context.Context, *CronJobListRequest) (*CronJobListResponse, error)
	mustEmbedUnimplementedCronServer()
}

// UnimplementedCronServer must be embedded to have forward compatible implementations.
type UnimplementedCronServer struct {
}

func (UnimplementedCronServer) Suspend(context.Context, *CronSuspendRequest) (*CronSuspendResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Suspend not implemented")
}
func (UnimplementedCronServer) Resume(context.Context, *CronResumeRequest) (*CronResumeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Resume not implemented")
}
func (UnimplementedCronServer) List(context.Context, *CronListRequest) (*CronListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (UnimplementedCronServer) JobList(context.Context, *CronJobListRequest) (*CronJobListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method JobList not implemented")
}
func (UnimplementedCronServer) mustEmbedUnimplementedCronServer() {}

// UnsafeCronServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CronServer will
// result in compilation errors.
type UnsafeCronServer interface {
	mustEmbedUnimplementedCronServer()
}

func RegisterCronServer(s grpc.ServiceRegistrar, srv CronServer) {
	s.RegisterService(&Cron_ServiceDesc, srv)
}

func _Cron_Suspend_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CronSuspendRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CronServer).Suspend(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/workflow.cron/Suspend",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CronServer).Suspend(ctx, req.(*CronSuspendRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Cron_Resume_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CronResumeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CronServer).Resume(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/workflow.cron/Resume",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CronServer).Resume(ctx, req.(*CronResumeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Cron_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CronListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CronServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/workflow.cron/List",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CronServer).List(ctx, req.(*CronListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Cron_JobList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CronJobListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CronServer).JobList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/workflow.cron/JobList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CronServer).JobList(ctx, req.(*CronJobListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Cron_ServiceDesc is the grpc.ServiceDesc for Cron service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Cron_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "workflow.cron",
	HandlerType: (*CronServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Suspend",
			Handler:    _Cron_Suspend_Handler,
		},
		{
			MethodName: "Resume",
			Handler:    _Cron_Resume_Handler,
		},
		{
			MethodName: "List",
			Handler:    _Cron_List_Handler,
		},
		{
			MethodName: "JobList",
			Handler:    _Cron_JobList_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "cron.proto",
}
