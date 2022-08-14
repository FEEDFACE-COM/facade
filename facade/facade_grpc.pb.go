// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.17.1
// source: facade.proto

package facade

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

// FacadeClient is the client API for Facade service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FacadeClient interface {
	Conf(ctx context.Context, in *Config, opts ...grpc.CallOption) (*Status, error)
	Pipe(ctx context.Context, opts ...grpc.CallOption) (Facade_PipeClient, error)
}

type facadeClient struct {
	cc grpc.ClientConnInterface
}

func NewFacadeClient(cc grpc.ClientConnInterface) FacadeClient {
	return &facadeClient{cc}
}

func (c *facadeClient) Conf(ctx context.Context, in *Config, opts ...grpc.CallOption) (*Status, error) {
	out := new(Status)
	err := c.cc.Invoke(ctx, "/facade.Facade/Conf", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *facadeClient) Pipe(ctx context.Context, opts ...grpc.CallOption) (Facade_PipeClient, error) {
	stream, err := c.cc.NewStream(ctx, &Facade_ServiceDesc.Streams[0], "/facade.Facade/Pipe", opts...)
	if err != nil {
		return nil, err
	}
	x := &facadePipeClient{stream}
	return x, nil
}

type Facade_PipeClient interface {
	Send(*RawText) error
	CloseAndRecv() (*Status, error)
	grpc.ClientStream
}

type facadePipeClient struct {
	grpc.ClientStream
}

func (x *facadePipeClient) Send(m *RawText) error {
	return x.ClientStream.SendMsg(m)
}

func (x *facadePipeClient) CloseAndRecv() (*Status, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(Status)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// FacadeServer is the server API for Facade service.
// All implementations should embed UnimplementedFacadeServer
// for forward compatibility
type FacadeServer interface {
	Conf(context.Context, *Config) (*Status, error)
	Pipe(Facade_PipeServer) error
}

// UnimplementedFacadeServer should be embedded to have forward compatible implementations.
type UnimplementedFacadeServer struct {
}

func (UnimplementedFacadeServer) Conf(context.Context, *Config) (*Status, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Conf not implemented")
}
func (UnimplementedFacadeServer) Pipe(Facade_PipeServer) error {
	return status.Errorf(codes.Unimplemented, "method Pipe not implemented")
}

// UnsafeFacadeServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FacadeServer will
// result in compilation errors.
type UnsafeFacadeServer interface {
	mustEmbedUnimplementedFacadeServer()
}

func RegisterFacadeServer(s grpc.ServiceRegistrar, srv FacadeServer) {
	s.RegisterService(&Facade_ServiceDesc, srv)
}

func _Facade_Conf_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Config)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FacadeServer).Conf(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/facade.Facade/Conf",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FacadeServer).Conf(ctx, req.(*Config))
	}
	return interceptor(ctx, in, info, handler)
}

func _Facade_Pipe_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(FacadeServer).Pipe(&facadePipeServer{stream})
}

type Facade_PipeServer interface {
	SendAndClose(*Status) error
	Recv() (*RawText, error)
	grpc.ServerStream
}

type facadePipeServer struct {
	grpc.ServerStream
}

func (x *facadePipeServer) SendAndClose(m *Status) error {
	return x.ServerStream.SendMsg(m)
}

func (x *facadePipeServer) Recv() (*RawText, error) {
	m := new(RawText)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Facade_ServiceDesc is the grpc.ServiceDesc for Facade service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Facade_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "facade.Facade",
	HandlerType: (*FacadeServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Conf",
			Handler:    _Facade_Conf_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Pipe",
			Handler:       _Facade_Pipe_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "facade.proto",
}
