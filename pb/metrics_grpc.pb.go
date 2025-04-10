// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: metrics.proto

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

// MetricsClient is the client API for Metrics service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MetricsClient interface {
	// The total number of HTTP and HTTPS requests for all HTTP methods across all environments.
	ClusterRequests(ctx context.Context, in *ClusterRequestsRequest, opts ...grpc.CallOption) (*ClusterRequestsResponse, error)
	// The number of HTTP codes that originate from all load balancers.
	ClusterResponseCodes(ctx context.Context, in *ClusterResponseCodesRequest, opts ...grpc.CallOption) (*ClusterResponseCodesResponse, error)
	// The total number of HTTP and HTTPS requests for all HTTP methods for an environment.
	Requests(ctx context.Context, in *RequestsRequest, opts ...grpc.CallOption) (*RequestsResponse, error)
	// The number of HTTP codes that originate from an environments load balancer.
	ResponseCodes(ctx context.Context, in *ResponseCodesRequest, opts ...grpc.CallOption) (*ResponseCodesResponse, error)
	// The time elapsed after the request leaves the load balancer until a response from the target is received.
	ResponseTimes(ctx context.Context, in *ResponseTimesRequest, opts ...grpc.CallOption) (*ResponseTimesResponse, error)
	// The percentage of all cacheable requests for which CDN served the content from its cache.
	CacheRatio(ctx context.Context, in *CacheRatioRequest, opts ...grpc.CallOption) (*CacheRatioResponse, error)
	// This error indicates a communication problem between CDN and the origin.
	OriginErrors(ctx context.Context, in *OriginErrorsRequest, opts ...grpc.CallOption) (*OriginErrorsResponse, error)
	// The number of requests which were sent to invalidate the CDN.
	InvalidationRequests(ctx context.Context, in *InvalidationRequestsRequest, opts ...grpc.CallOption) (*InvalidationRequestsResponse, error)
	// The number of paths which were invalidated by requests.
	InvalidationPaths(ctx context.Context, in *InvalidationPathsRequest, opts ...grpc.CallOption) (*InvalidationPathsResponse, error)
	// The amount of resources used by an environment.
	ResourceUsage(ctx context.Context, in *ResourceUsageRequest, opts ...grpc.CallOption) (*ResourceUsageResponse, error)
}

type metricsClient struct {
	cc grpc.ClientConnInterface
}

func NewMetricsClient(cc grpc.ClientConnInterface) MetricsClient {
	return &metricsClient{cc}
}

func (c *metricsClient) ClusterRequests(ctx context.Context, in *ClusterRequestsRequest, opts ...grpc.CallOption) (*ClusterRequestsResponse, error) {
	out := new(ClusterRequestsResponse)
	err := c.cc.Invoke(ctx, "/workflow.metrics/ClusterRequests", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metricsClient) ClusterResponseCodes(ctx context.Context, in *ClusterResponseCodesRequest, opts ...grpc.CallOption) (*ClusterResponseCodesResponse, error) {
	out := new(ClusterResponseCodesResponse)
	err := c.cc.Invoke(ctx, "/workflow.metrics/ClusterResponseCodes", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metricsClient) Requests(ctx context.Context, in *RequestsRequest, opts ...grpc.CallOption) (*RequestsResponse, error) {
	out := new(RequestsResponse)
	err := c.cc.Invoke(ctx, "/workflow.metrics/Requests", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metricsClient) ResponseCodes(ctx context.Context, in *ResponseCodesRequest, opts ...grpc.CallOption) (*ResponseCodesResponse, error) {
	out := new(ResponseCodesResponse)
	err := c.cc.Invoke(ctx, "/workflow.metrics/ResponseCodes", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metricsClient) ResponseTimes(ctx context.Context, in *ResponseTimesRequest, opts ...grpc.CallOption) (*ResponseTimesResponse, error) {
	out := new(ResponseTimesResponse)
	err := c.cc.Invoke(ctx, "/workflow.metrics/ResponseTimes", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metricsClient) CacheRatio(ctx context.Context, in *CacheRatioRequest, opts ...grpc.CallOption) (*CacheRatioResponse, error) {
	out := new(CacheRatioResponse)
	err := c.cc.Invoke(ctx, "/workflow.metrics/CacheRatio", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metricsClient) OriginErrors(ctx context.Context, in *OriginErrorsRequest, opts ...grpc.CallOption) (*OriginErrorsResponse, error) {
	out := new(OriginErrorsResponse)
	err := c.cc.Invoke(ctx, "/workflow.metrics/OriginErrors", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metricsClient) InvalidationRequests(ctx context.Context, in *InvalidationRequestsRequest, opts ...grpc.CallOption) (*InvalidationRequestsResponse, error) {
	out := new(InvalidationRequestsResponse)
	err := c.cc.Invoke(ctx, "/workflow.metrics/InvalidationRequests", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metricsClient) InvalidationPaths(ctx context.Context, in *InvalidationPathsRequest, opts ...grpc.CallOption) (*InvalidationPathsResponse, error) {
	out := new(InvalidationPathsResponse)
	err := c.cc.Invoke(ctx, "/workflow.metrics/InvalidationPaths", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metricsClient) ResourceUsage(ctx context.Context, in *ResourceUsageRequest, opts ...grpc.CallOption) (*ResourceUsageResponse, error) {
	out := new(ResourceUsageResponse)
	err := c.cc.Invoke(ctx, "/workflow.metrics/ResourceUsage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MetricsServer is the server API for Metrics service.
// All implementations must embed UnimplementedMetricsServer
// for forward compatibility
type MetricsServer interface {
	// The total number of HTTP and HTTPS requests for all HTTP methods across all environments.
	ClusterRequests(context.Context, *ClusterRequestsRequest) (*ClusterRequestsResponse, error)
	// The number of HTTP codes that originate from all load balancers.
	ClusterResponseCodes(context.Context, *ClusterResponseCodesRequest) (*ClusterResponseCodesResponse, error)
	// The total number of HTTP and HTTPS requests for all HTTP methods for an environment.
	Requests(context.Context, *RequestsRequest) (*RequestsResponse, error)
	// The number of HTTP codes that originate from an environments load balancer.
	ResponseCodes(context.Context, *ResponseCodesRequest) (*ResponseCodesResponse, error)
	// The time elapsed after the request leaves the load balancer until a response from the target is received.
	ResponseTimes(context.Context, *ResponseTimesRequest) (*ResponseTimesResponse, error)
	// The percentage of all cacheable requests for which CDN served the content from its cache.
	CacheRatio(context.Context, *CacheRatioRequest) (*CacheRatioResponse, error)
	// This error indicates a communication problem between CDN and the origin.
	OriginErrors(context.Context, *OriginErrorsRequest) (*OriginErrorsResponse, error)
	// The number of requests which were sent to invalidate the CDN.
	InvalidationRequests(context.Context, *InvalidationRequestsRequest) (*InvalidationRequestsResponse, error)
	// The number of paths which were invalidated by requests.
	InvalidationPaths(context.Context, *InvalidationPathsRequest) (*InvalidationPathsResponse, error)
	// The amount of resources used by an environment.
	ResourceUsage(context.Context, *ResourceUsageRequest) (*ResourceUsageResponse, error)
	mustEmbedUnimplementedMetricsServer()
}

// UnimplementedMetricsServer must be embedded to have forward compatible implementations.
type UnimplementedMetricsServer struct {
}

func (UnimplementedMetricsServer) ClusterRequests(context.Context, *ClusterRequestsRequest) (*ClusterRequestsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClusterRequests not implemented")
}
func (UnimplementedMetricsServer) ClusterResponseCodes(context.Context, *ClusterResponseCodesRequest) (*ClusterResponseCodesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClusterResponseCodes not implemented")
}
func (UnimplementedMetricsServer) Requests(context.Context, *RequestsRequest) (*RequestsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Requests not implemented")
}
func (UnimplementedMetricsServer) ResponseCodes(context.Context, *ResponseCodesRequest) (*ResponseCodesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResponseCodes not implemented")
}
func (UnimplementedMetricsServer) ResponseTimes(context.Context, *ResponseTimesRequest) (*ResponseTimesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResponseTimes not implemented")
}
func (UnimplementedMetricsServer) CacheRatio(context.Context, *CacheRatioRequest) (*CacheRatioResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CacheRatio not implemented")
}
func (UnimplementedMetricsServer) OriginErrors(context.Context, *OriginErrorsRequest) (*OriginErrorsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OriginErrors not implemented")
}
func (UnimplementedMetricsServer) InvalidationRequests(context.Context, *InvalidationRequestsRequest) (*InvalidationRequestsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InvalidationRequests not implemented")
}
func (UnimplementedMetricsServer) InvalidationPaths(context.Context, *InvalidationPathsRequest) (*InvalidationPathsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InvalidationPaths not implemented")
}
func (UnimplementedMetricsServer) ResourceUsage(context.Context, *ResourceUsageRequest) (*ResourceUsageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResourceUsage not implemented")
}
func (UnimplementedMetricsServer) mustEmbedUnimplementedMetricsServer() {}

// UnsafeMetricsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MetricsServer will
// result in compilation errors.
type UnsafeMetricsServer interface {
	mustEmbedUnimplementedMetricsServer()
}

func RegisterMetricsServer(s grpc.ServiceRegistrar, srv MetricsServer) {
	s.RegisterService(&Metrics_ServiceDesc, srv)
}

func _Metrics_ClusterRequests_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClusterRequestsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricsServer).ClusterRequests(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/workflow.metrics/ClusterRequests",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricsServer).ClusterRequests(ctx, req.(*ClusterRequestsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Metrics_ClusterResponseCodes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClusterResponseCodesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricsServer).ClusterResponseCodes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/workflow.metrics/ClusterResponseCodes",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricsServer).ClusterResponseCodes(ctx, req.(*ClusterResponseCodesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Metrics_Requests_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricsServer).Requests(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/workflow.metrics/Requests",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricsServer).Requests(ctx, req.(*RequestsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Metrics_ResponseCodes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResponseCodesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricsServer).ResponseCodes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/workflow.metrics/ResponseCodes",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricsServer).ResponseCodes(ctx, req.(*ResponseCodesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Metrics_ResponseTimes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResponseTimesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricsServer).ResponseTimes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/workflow.metrics/ResponseTimes",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricsServer).ResponseTimes(ctx, req.(*ResponseTimesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Metrics_CacheRatio_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CacheRatioRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricsServer).CacheRatio(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/workflow.metrics/CacheRatio",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricsServer).CacheRatio(ctx, req.(*CacheRatioRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Metrics_OriginErrors_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OriginErrorsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricsServer).OriginErrors(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/workflow.metrics/OriginErrors",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricsServer).OriginErrors(ctx, req.(*OriginErrorsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Metrics_InvalidationRequests_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InvalidationRequestsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricsServer).InvalidationRequests(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/workflow.metrics/InvalidationRequests",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricsServer).InvalidationRequests(ctx, req.(*InvalidationRequestsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Metrics_InvalidationPaths_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InvalidationPathsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricsServer).InvalidationPaths(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/workflow.metrics/InvalidationPaths",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricsServer).InvalidationPaths(ctx, req.(*InvalidationPathsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Metrics_ResourceUsage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResourceUsageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricsServer).ResourceUsage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/workflow.metrics/ResourceUsage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricsServer).ResourceUsage(ctx, req.(*ResourceUsageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Metrics_ServiceDesc is the grpc.ServiceDesc for Metrics service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Metrics_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "workflow.metrics",
	HandlerType: (*MetricsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ClusterRequests",
			Handler:    _Metrics_ClusterRequests_Handler,
		},
		{
			MethodName: "ClusterResponseCodes",
			Handler:    _Metrics_ClusterResponseCodes_Handler,
		},
		{
			MethodName: "Requests",
			Handler:    _Metrics_Requests_Handler,
		},
		{
			MethodName: "ResponseCodes",
			Handler:    _Metrics_ResponseCodes_Handler,
		},
		{
			MethodName: "ResponseTimes",
			Handler:    _Metrics_ResponseTimes_Handler,
		},
		{
			MethodName: "CacheRatio",
			Handler:    _Metrics_CacheRatio_Handler,
		},
		{
			MethodName: "OriginErrors",
			Handler:    _Metrics_OriginErrors_Handler,
		},
		{
			MethodName: "InvalidationRequests",
			Handler:    _Metrics_InvalidationRequests_Handler,
		},
		{
			MethodName: "InvalidationPaths",
			Handler:    _Metrics_InvalidationPaths_Handler,
		},
		{
			MethodName: "ResourceUsage",
			Handler:    _Metrics_ResourceUsage_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "metrics.proto",
}
