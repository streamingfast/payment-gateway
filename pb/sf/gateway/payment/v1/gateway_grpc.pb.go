// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: sf/gateway/payment/v1/gateway.proto

package pbgateway

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

const (
	DiscoverService_Services_FullMethodName = "/sf.gateway.payment.v1.DiscoverService/Services"
)

// DiscoverServiceClient is the client API for DiscoverService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DiscoverServiceClient interface {
	Services(ctx context.Context, in *ServicesRequest, opts ...grpc.CallOption) (*ServicesResponse, error)
}

type discoverServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDiscoverServiceClient(cc grpc.ClientConnInterface) DiscoverServiceClient {
	return &discoverServiceClient{cc}
}

func (c *discoverServiceClient) Services(ctx context.Context, in *ServicesRequest, opts ...grpc.CallOption) (*ServicesResponse, error) {
	out := new(ServicesResponse)
	err := c.cc.Invoke(ctx, DiscoverService_Services_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DiscoverServiceServer is the server API for DiscoverService service.
// All implementations must embed UnimplementedDiscoverServiceServer
// for forward compatibility
type DiscoverServiceServer interface {
	Services(context.Context, *ServicesRequest) (*ServicesResponse, error)
	mustEmbedUnimplementedDiscoverServiceServer()
}

// UnimplementedDiscoverServiceServer must be embedded to have forward compatible implementations.
type UnimplementedDiscoverServiceServer struct {
}

func (UnimplementedDiscoverServiceServer) Services(context.Context, *ServicesRequest) (*ServicesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Services not implemented")
}
func (UnimplementedDiscoverServiceServer) mustEmbedUnimplementedDiscoverServiceServer() {}

// UnsafeDiscoverServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DiscoverServiceServer will
// result in compilation errors.
type UnsafeDiscoverServiceServer interface {
	mustEmbedUnimplementedDiscoverServiceServer()
}

func RegisterDiscoverServiceServer(s grpc.ServiceRegistrar, srv DiscoverServiceServer) {
	s.RegisterService(&DiscoverService_ServiceDesc, srv)
}

func _DiscoverService_Services_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ServicesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DiscoverServiceServer).Services(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DiscoverService_Services_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DiscoverServiceServer).Services(ctx, req.(*ServicesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// DiscoverService_ServiceDesc is the grpc.ServiceDesc for DiscoverService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DiscoverService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "sf.gateway.payment.v1.DiscoverService",
	HandlerType: (*DiscoverServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Services",
			Handler:    _DiscoverService_Services_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "sf/gateway/payment/v1/gateway.proto",
}

const (
	UsageService_Report_FullMethodName = "/sf.gateway.payment.v1.UsageService/Report"
)

// UsageServiceClient is the client API for UsageService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UsageServiceClient interface {
	Report(ctx context.Context, in *ReportRequest, opts ...grpc.CallOption) (*ReportResponse, error)
}

type usageServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewUsageServiceClient(cc grpc.ClientConnInterface) UsageServiceClient {
	return &usageServiceClient{cc}
}

func (c *usageServiceClient) Report(ctx context.Context, in *ReportRequest, opts ...grpc.CallOption) (*ReportResponse, error) {
	out := new(ReportResponse)
	err := c.cc.Invoke(ctx, UsageService_Report_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UsageServiceServer is the server API for UsageService service.
// All implementations must embed UnimplementedUsageServiceServer
// for forward compatibility
type UsageServiceServer interface {
	Report(context.Context, *ReportRequest) (*ReportResponse, error)
	mustEmbedUnimplementedUsageServiceServer()
}

// UnimplementedUsageServiceServer must be embedded to have forward compatible implementations.
type UnimplementedUsageServiceServer struct {
}

func (UnimplementedUsageServiceServer) Report(context.Context, *ReportRequest) (*ReportResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Report not implemented")
}
func (UnimplementedUsageServiceServer) mustEmbedUnimplementedUsageServiceServer() {}

// UnsafeUsageServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UsageServiceServer will
// result in compilation errors.
type UnsafeUsageServiceServer interface {
	mustEmbedUnimplementedUsageServiceServer()
}

func RegisterUsageServiceServer(s grpc.ServiceRegistrar, srv UsageServiceServer) {
	s.RegisterService(&UsageService_ServiceDesc, srv)
}

func _UsageService_Report_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReportRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsageServiceServer).Report(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UsageService_Report_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsageServiceServer).Report(ctx, req.(*ReportRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// UsageService_ServiceDesc is the grpc.ServiceDesc for UsageService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UsageService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "sf.gateway.payment.v1.UsageService",
	HandlerType: (*UsageServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Report",
			Handler:    _UsageService_Report_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "sf/gateway/payment/v1/gateway.proto",
}

const (
	AuthService_Authenticate_FullMethodName = "/sf.gateway.payment.v1.AuthService/Authenticate"
)

// AuthServiceClient is the client API for AuthService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AuthServiceClient interface {
	Authenticate(ctx context.Context, in *AuthenticateRequest, opts ...grpc.CallOption) (*AuthenticateResponse, error)
}

type authServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAuthServiceClient(cc grpc.ClientConnInterface) AuthServiceClient {
	return &authServiceClient{cc}
}

func (c *authServiceClient) Authenticate(ctx context.Context, in *AuthenticateRequest, opts ...grpc.CallOption) (*AuthenticateResponse, error) {
	out := new(AuthenticateResponse)
	err := c.cc.Invoke(ctx, AuthService_Authenticate_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuthServiceServer is the server API for AuthService service.
// All implementations must embed UnimplementedAuthServiceServer
// for forward compatibility
type AuthServiceServer interface {
	Authenticate(context.Context, *AuthenticateRequest) (*AuthenticateResponse, error)
	mustEmbedUnimplementedAuthServiceServer()
}

// UnimplementedAuthServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAuthServiceServer struct {
}

func (UnimplementedAuthServiceServer) Authenticate(context.Context, *AuthenticateRequest) (*AuthenticateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Authenticate not implemented")
}
func (UnimplementedAuthServiceServer) mustEmbedUnimplementedAuthServiceServer() {}

// UnsafeAuthServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AuthServiceServer will
// result in compilation errors.
type UnsafeAuthServiceServer interface {
	mustEmbedUnimplementedAuthServiceServer()
}

func RegisterAuthServiceServer(s grpc.ServiceRegistrar, srv AuthServiceServer) {
	s.RegisterService(&AuthService_ServiceDesc, srv)
}

func _AuthService_Authenticate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AuthenticateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).Authenticate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthService_Authenticate_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).Authenticate(ctx, req.(*AuthenticateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AuthService_ServiceDesc is the grpc.ServiceDesc for AuthService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AuthService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "sf.gateway.payment.v1.AuthService",
	HandlerType: (*AuthServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Authenticate",
			Handler:    _AuthService_Authenticate_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "sf/gateway/payment/v1/gateway.proto",
}
