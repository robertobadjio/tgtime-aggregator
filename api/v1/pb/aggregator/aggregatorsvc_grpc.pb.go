// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.3
// source: aggregatorsvc.proto

package aggregatorsvc

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
	Aggregator_CreateTime_FullMethodName = "/Aggregator/CreateTime"
)

// AggregatorClient is the client API for Aggregator service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AggregatorClient interface {
	CreateTime(ctx context.Context, in *CreateTimeRequest, opts ...grpc.CallOption) (*CreateTimeResponse, error)
}

type aggregatorClient struct {
	cc grpc.ClientConnInterface
}

func NewAggregatorClient(cc grpc.ClientConnInterface) AggregatorClient {
	return &aggregatorClient{cc}
}

func (c *aggregatorClient) CreateTime(ctx context.Context, in *CreateTimeRequest, opts ...grpc.CallOption) (*CreateTimeResponse, error) {
	out := new(CreateTimeResponse)
	err := c.cc.Invoke(ctx, Aggregator_CreateTime_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AggregatorServer is the server API for Aggregator service.
// All implementations must embed UnimplementedAggregatorServer
// for forward compatibility
type AggregatorServer interface {
	CreateTime(context.Context, *CreateTimeRequest) (*CreateTimeResponse, error)
	mustEmbedUnimplementedAggregatorServer()
}

// UnimplementedAggregatorServer must be embedded to have forward compatible implementations.
type UnimplementedAggregatorServer struct {
}

func (UnimplementedAggregatorServer) CreateTime(context.Context, *CreateTimeRequest) (*CreateTimeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateTime not implemented")
}
func (UnimplementedAggregatorServer) mustEmbedUnimplementedAggregatorServer() {}

// UnsafeAggregatorServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AggregatorServer will
// result in compilation errors.
type UnsafeAggregatorServer interface {
	mustEmbedUnimplementedAggregatorServer()
}

func RegisterAggregatorServer(s grpc.ServiceRegistrar, srv AggregatorServer) {
	s.RegisterService(&Aggregator_ServiceDesc, srv)
}

func _Aggregator_CreateTime_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateTimeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AggregatorServer).CreateTime(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Aggregator_CreateTime_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AggregatorServer).CreateTime(ctx, req.(*CreateTimeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Aggregator_ServiceDesc is the grpc.ServiceDesc for Aggregator service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Aggregator_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Aggregator",
	HandlerType: (*AggregatorServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateTime",
			Handler:    _Aggregator_CreateTime_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "aggregatorsvc.proto",
}
