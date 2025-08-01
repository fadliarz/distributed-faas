// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.6.1
// source: registrar-service/v1/api.proto

package registrar_service_v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	RegistrarService_RegisterMachine_FullMethodName = "/fadliarz.distributed_faas.registrar_service.v1.RegistrarService/RegisterMachine"
)

// RegistrarServiceClient is the client API for RegistrarService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RegistrarServiceClient interface {
	RegisterMachine(ctx context.Context, in *RegisterMachineRequest, opts ...grpc.CallOption) (*RegisterMachineResponse, error)
}

type registrarServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRegistrarServiceClient(cc grpc.ClientConnInterface) RegistrarServiceClient {
	return &registrarServiceClient{cc}
}

func (c *registrarServiceClient) RegisterMachine(ctx context.Context, in *RegisterMachineRequest, opts ...grpc.CallOption) (*RegisterMachineResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RegisterMachineResponse)
	err := c.cc.Invoke(ctx, RegistrarService_RegisterMachine_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RegistrarServiceServer is the server API for RegistrarService service.
// All implementations must embed UnimplementedRegistrarServiceServer
// for forward compatibility.
type RegistrarServiceServer interface {
	RegisterMachine(context.Context, *RegisterMachineRequest) (*RegisterMachineResponse, error)
	mustEmbedUnimplementedRegistrarServiceServer()
}

// UnimplementedRegistrarServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedRegistrarServiceServer struct{}

func (UnimplementedRegistrarServiceServer) RegisterMachine(context.Context, *RegisterMachineRequest) (*RegisterMachineResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterMachine not implemented")
}
func (UnimplementedRegistrarServiceServer) mustEmbedUnimplementedRegistrarServiceServer() {}
func (UnimplementedRegistrarServiceServer) testEmbeddedByValue()                          {}

// UnsafeRegistrarServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RegistrarServiceServer will
// result in compilation errors.
type UnsafeRegistrarServiceServer interface {
	mustEmbedUnimplementedRegistrarServiceServer()
}

func RegisterRegistrarServiceServer(s grpc.ServiceRegistrar, srv RegistrarServiceServer) {
	// If the following call pancis, it indicates UnimplementedRegistrarServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&RegistrarService_ServiceDesc, srv)
}

func _RegistrarService_RegisterMachine_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterMachineRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegistrarServiceServer).RegisterMachine(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RegistrarService_RegisterMachine_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegistrarServiceServer).RegisterMachine(ctx, req.(*RegisterMachineRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RegistrarService_ServiceDesc is the grpc.ServiceDesc for RegistrarService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RegistrarService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "fadliarz.distributed_faas.registrar_service.v1.RegistrarService",
	HandlerType: (*RegistrarServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RegisterMachine",
			Handler:    _RegistrarService_RegisterMachine_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "registrar-service/v1/api.proto",
}
