// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package stats

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

// StatsClient is the client API for Stats service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type StatsClient interface {
	Get(ctx context.Context, in *Statistics, opts ...grpc.CallOption) (*Statistics, error)
}

type statsClient struct {
	cc grpc.ClientConnInterface
}

func NewStatsClient(cc grpc.ClientConnInterface) StatsClient {
	return &statsClient{cc}
}

func (c *statsClient) Get(ctx context.Context, in *Statistics, opts ...grpc.CallOption) (*Statistics, error) {
	out := new(Statistics)
	err := c.cc.Invoke(ctx, "/Stats/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// StatsServer is the server API for Stats service.
// All implementations must embed UnimplementedStatsServer
// for forward compatibility
type StatsServer interface {
	Get(context.Context, *Statistics) (*Statistics, error)
	mustEmbedUnimplementedStatsServer()
}

// UnimplementedStatsServer must be embedded to have forward compatible implementations.
type UnimplementedStatsServer struct {
}

func (UnimplementedStatsServer) Get(context.Context, *Statistics) (*Statistics, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedStatsServer) mustEmbedUnimplementedStatsServer() {}

// UnsafeStatsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StatsServer will
// result in compilation errors.
type UnsafeStatsServer interface {
	mustEmbedUnimplementedStatsServer()
}

func RegisterStatsServer(s grpc.ServiceRegistrar, srv StatsServer) {
	s.RegisterService(&Stats_ServiceDesc, srv)
}

func _Stats_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Statistics)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StatsServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Stats/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StatsServer).Get(ctx, req.(*Statistics))
	}
	return interceptor(ctx, in, info, handler)
}

// Stats_ServiceDesc is the grpc.ServiceDesc for Stats service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Stats_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Stats",
	HandlerType: (*StatsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Get",
			Handler:    _Stats_Get_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "wg.proto",
}
