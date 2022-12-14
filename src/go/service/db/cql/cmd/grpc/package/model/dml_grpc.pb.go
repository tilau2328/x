// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.9
// source: dml.proto

package model

import (
	grpc "google.golang.org/grpc"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// DMLClient is the client API for DML service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DMLClient interface {
}

type dMLClient struct {
	cc grpc.ClientConnInterface
}

func NewDMLClient(cc grpc.ClientConnInterface) DMLClient {
	return &dMLClient{cc}
}

// DMLServer is the server API for DML service.
// All implementations must embed UnimplementedDMLServer
// for forward compatibility
type DMLServer interface {
	mustEmbedUnimplementedDMLServer()
}

// UnimplementedDMLServer must be embedded to have forward compatible implementations.
type UnimplementedDMLServer struct {
}

func (UnimplementedDMLServer) mustEmbedUnimplementedDMLServer() {}

// UnsafeDMLServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DMLServer will
// result in compilation errors.
type UnsafeDMLServer interface {
	mustEmbedUnimplementedDMLServer()
}

func RegisterDMLServer(s grpc.ServiceRegistrar, srv DMLServer) {
	s.RegisterService(&DML_ServiceDesc, srv)
}

// DML_ServiceDesc is the grpc.ServiceDesc for DML service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DML_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "cql.grpc.v1.DML",
	HandlerType: (*DMLServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams:     []grpc.StreamDesc{},
	Metadata:    "dml.proto",
}
