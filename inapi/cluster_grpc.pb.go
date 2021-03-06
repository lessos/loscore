// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package inapi

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// ApiHostMemberClient is the client API for ApiHostMember service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ApiHostMemberClient interface {
	HostJoin(ctx context.Context, in *ResHostNew, opts ...grpc.CallOption) (*ResHost, error)
}

type apiHostMemberClient struct {
	cc grpc.ClientConnInterface
}

func NewApiHostMemberClient(cc grpc.ClientConnInterface) ApiHostMemberClient {
	return &apiHostMemberClient{cc}
}

func (c *apiHostMemberClient) HostJoin(ctx context.Context, in *ResHostNew, opts ...grpc.CallOption) (*ResHost, error) {
	out := new(ResHost)
	err := c.cc.Invoke(ctx, "/inapi.ApiHostMember/HostJoin", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ApiHostMemberServer is the server API for ApiHostMember service.
// All implementations must embed UnimplementedApiHostMemberServer
// for forward compatibility
type ApiHostMemberServer interface {
	HostJoin(context.Context, *ResHostNew) (*ResHost, error)
	mustEmbedUnimplementedApiHostMemberServer()
}

// UnimplementedApiHostMemberServer must be embedded to have forward compatible implementations.
type UnimplementedApiHostMemberServer struct {
}

func (*UnimplementedApiHostMemberServer) HostJoin(context.Context, *ResHostNew) (*ResHost, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HostJoin not implemented")
}
func (*UnimplementedApiHostMemberServer) mustEmbedUnimplementedApiHostMemberServer() {}

func RegisterApiHostMemberServer(s *grpc.Server, srv ApiHostMemberServer) {
	s.RegisterService(&_ApiHostMember_serviceDesc, srv)
}

func _ApiHostMember_HostJoin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResHostNew)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApiHostMemberServer).HostJoin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/inapi.ApiHostMember/HostJoin",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApiHostMemberServer).HostJoin(ctx, req.(*ResHostNew))
	}
	return interceptor(ctx, in, info, handler)
}

var _ApiHostMember_serviceDesc = grpc.ServiceDesc{
	ServiceName: "inapi.ApiHostMember",
	HandlerType: (*ApiHostMemberServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "HostJoin",
			Handler:    _ApiHostMember_HostJoin_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "cluster.proto",
}

// ApiZoneMasterClient is the client API for ApiZoneMaster service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ApiZoneMasterClient interface {
	HostConfig(ctx context.Context, in *ZoneHostConfigRequest, opts ...grpc.CallOption) (*ZoneHostConfigReply, error)
	HostStatusSync(ctx context.Context, in *ResHost, opts ...grpc.CallOption) (*ResHostBound, error)
}

type apiZoneMasterClient struct {
	cc grpc.ClientConnInterface
}

func NewApiZoneMasterClient(cc grpc.ClientConnInterface) ApiZoneMasterClient {
	return &apiZoneMasterClient{cc}
}

func (c *apiZoneMasterClient) HostConfig(ctx context.Context, in *ZoneHostConfigRequest, opts ...grpc.CallOption) (*ZoneHostConfigReply, error) {
	out := new(ZoneHostConfigReply)
	err := c.cc.Invoke(ctx, "/inapi.ApiZoneMaster/HostConfig", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *apiZoneMasterClient) HostStatusSync(ctx context.Context, in *ResHost, opts ...grpc.CallOption) (*ResHostBound, error) {
	out := new(ResHostBound)
	err := c.cc.Invoke(ctx, "/inapi.ApiZoneMaster/HostStatusSync", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ApiZoneMasterServer is the server API for ApiZoneMaster service.
// All implementations must embed UnimplementedApiZoneMasterServer
// for forward compatibility
type ApiZoneMasterServer interface {
	HostConfig(context.Context, *ZoneHostConfigRequest) (*ZoneHostConfigReply, error)
	HostStatusSync(context.Context, *ResHost) (*ResHostBound, error)
	mustEmbedUnimplementedApiZoneMasterServer()
}

// UnimplementedApiZoneMasterServer must be embedded to have forward compatible implementations.
type UnimplementedApiZoneMasterServer struct {
}

func (*UnimplementedApiZoneMasterServer) HostConfig(context.Context, *ZoneHostConfigRequest) (*ZoneHostConfigReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HostConfig not implemented")
}
func (*UnimplementedApiZoneMasterServer) HostStatusSync(context.Context, *ResHost) (*ResHostBound, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HostStatusSync not implemented")
}
func (*UnimplementedApiZoneMasterServer) mustEmbedUnimplementedApiZoneMasterServer() {}

func RegisterApiZoneMasterServer(s *grpc.Server, srv ApiZoneMasterServer) {
	s.RegisterService(&_ApiZoneMaster_serviceDesc, srv)
}

func _ApiZoneMaster_HostConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ZoneHostConfigRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApiZoneMasterServer).HostConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/inapi.ApiZoneMaster/HostConfig",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApiZoneMasterServer).HostConfig(ctx, req.(*ZoneHostConfigRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ApiZoneMaster_HostStatusSync_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResHost)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApiZoneMasterServer).HostStatusSync(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/inapi.ApiZoneMaster/HostStatusSync",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApiZoneMasterServer).HostStatusSync(ctx, req.(*ResHost))
	}
	return interceptor(ctx, in, info, handler)
}

var _ApiZoneMaster_serviceDesc = grpc.ServiceDesc{
	ServiceName: "inapi.ApiZoneMaster",
	HandlerType: (*ApiZoneMasterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "HostConfig",
			Handler:    _ApiZoneMaster_HostConfig_Handler,
		},
		{
			MethodName: "HostStatusSync",
			Handler:    _ApiZoneMaster_HostStatusSync_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "cluster.proto",
}
