// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.9
// source: ddl.proto

package model

import (
	grpc "google.golang.org/grpc"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// DDLClient is the client API for DDL service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DDLClient interface {
}

type dDLClient struct {
	cc grpc.ClientConnInterface
}

func NewDDLClient(cc grpc.ClientConnInterface) DDLClient {
	return &dDLClient{cc}
}

// DDLServer is the server API for DDL service.
// All implementations must embed UnimplementedDDLServer
// for forward compatibility
type DDLServer interface {
	mustEmbedUnimplementedDDLServer()
}

// UnimplementedDDLServer must be embedded to have forward compatible implementations.
type UnimplementedDDLServer struct {
}

func (UnimplementedDDLServer) mustEmbedUnimplementedDDLServer() {}

// UnsafeDDLServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DDLServer will
// result in compilation errors.
type UnsafeDDLServer interface {
	mustEmbedUnimplementedDDLServer()
}

func RegisterDDLServer(s grpc.ServiceRegistrar, srv DDLServer) {
	s.RegisterService(&DDL_ServiceDesc, srv)
}

// DDL_ServiceDesc is the grpc.ServiceDesc for DDL service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DDL_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "cql.grpc.v1.DDL",
	HandlerType: (*DDLServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams:     []grpc.StreamDesc{},
	Metadata:    "ddl.proto",
}