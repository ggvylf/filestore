// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.6.1
// source: proxy.proto

package proto

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

// DBProxyServiceClient is the client API for DBProxyService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DBProxyServiceClient interface {
	// 请求执行sql动作
	ExecuteAction(ctx context.Context, in *ReqExec, opts ...grpc.CallOption) (*RespExec, error)
}

type dBProxyServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDBProxyServiceClient(cc grpc.ClientConnInterface) DBProxyServiceClient {
	return &dBProxyServiceClient{cc}
}

func (c *dBProxyServiceClient) ExecuteAction(ctx context.Context, in *ReqExec, opts ...grpc.CallOption) (*RespExec, error) {
	out := new(RespExec)
	err := c.cc.Invoke(ctx, "/proto.DBProxyService/ExecuteAction", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DBProxyServiceServer is the server API for DBProxyService service.
// All implementations must embed UnimplementedDBProxyServiceServer
// for forward compatibility
type DBProxyServiceServer interface {
	// 请求执行sql动作
	ExecuteAction(context.Context, *ReqExec) (*RespExec, error)
	mustEmbedUnimplementedDBProxyServiceServer()
}

// UnimplementedDBProxyServiceServer must be embedded to have forward compatible implementations.
type UnimplementedDBProxyServiceServer struct {
}

func (UnimplementedDBProxyServiceServer) ExecuteAction(context.Context, *ReqExec) (*RespExec, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExecuteAction not implemented")
}
func (UnimplementedDBProxyServiceServer) mustEmbedUnimplementedDBProxyServiceServer() {}

// UnsafeDBProxyServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DBProxyServiceServer will
// result in compilation errors.
type UnsafeDBProxyServiceServer interface {
	mustEmbedUnimplementedDBProxyServiceServer()
}

func RegisterDBProxyServiceServer(s grpc.ServiceRegistrar, srv DBProxyServiceServer) {
	s.RegisterService(&DBProxyService_ServiceDesc, srv)
}

func _DBProxyService_ExecuteAction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReqExec)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DBProxyServiceServer).ExecuteAction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.DBProxyService/ExecuteAction",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DBProxyServiceServer).ExecuteAction(ctx, req.(*ReqExec))
	}
	return interceptor(ctx, in, info, handler)
}

// DBProxyService_ServiceDesc is the grpc.ServiceDesc for DBProxyService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DBProxyService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.DBProxyService",
	HandlerType: (*DBProxyServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ExecuteAction",
			Handler:    _DBProxyService_ExecuteAction_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proxy.proto",
}
