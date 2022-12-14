// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.9
// source: ddl.proto

package model

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

// DDLClient is the client API for DDL service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DDLClient interface {
	CreateKeySpaces(ctx context.Context, in *CreateKeySpacesRequest, opts ...grpc.CallOption) (*Empty, error)
	AlterKeySpaces(ctx context.Context, in *AlterKeySpacesRequest, opts ...grpc.CallOption) (*Empty, error)
	DropKeySpaces(ctx context.Context, in *DropKeySpacesRequest, opts ...grpc.CallOption) (*Empty, error)
	ListKeySpaces(ctx context.Context, in *ListKeySpacesRequest, opts ...grpc.CallOption) (*Empty, error)
	GetKeySpace(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
	CreateTables(ctx context.Context, in *CreateTablesRequest, opts ...grpc.CallOption) (*Empty, error)
	AlterTables(ctx context.Context, in *AlterTablesRequest, opts ...grpc.CallOption) (*Empty, error)
	DropTables(ctx context.Context, in *DropTablesRequest, opts ...grpc.CallOption) (*Empty, error)
	ListTables(ctx context.Context, in *ListTablesRequest, opts ...grpc.CallOption) (*Empty, error)
	GetTable(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
}

type dDLClient struct {
	cc grpc.ClientConnInterface
}

func NewDDLClient(cc grpc.ClientConnInterface) DDLClient {
	return &dDLClient{cc}
}

func (c *dDLClient) CreateKeySpaces(ctx context.Context, in *CreateKeySpacesRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/cql.grpc.v1.DDL/CreateKeySpaces", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dDLClient) AlterKeySpaces(ctx context.Context, in *AlterKeySpacesRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/cql.grpc.v1.DDL/AlterKeySpaces", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dDLClient) DropKeySpaces(ctx context.Context, in *DropKeySpacesRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/cql.grpc.v1.DDL/DropKeySpaces", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dDLClient) ListKeySpaces(ctx context.Context, in *ListKeySpacesRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/cql.grpc.v1.DDL/ListKeySpaces", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dDLClient) GetKeySpace(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/cql.grpc.v1.DDL/GetKeySpace", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dDLClient) CreateTables(ctx context.Context, in *CreateTablesRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/cql.grpc.v1.DDL/CreateTables", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dDLClient) AlterTables(ctx context.Context, in *AlterTablesRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/cql.grpc.v1.DDL/AlterTables", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dDLClient) DropTables(ctx context.Context, in *DropTablesRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/cql.grpc.v1.DDL/DropTables", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dDLClient) ListTables(ctx context.Context, in *ListTablesRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/cql.grpc.v1.DDL/ListTables", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dDLClient) GetTable(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/cql.grpc.v1.DDL/GetTable", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DDLServer is the server API for DDL service.
// All implementations must embed UnimplementedDDLServer
// for forward compatibility
type DDLServer interface {
	CreateKeySpaces(context.Context, *CreateKeySpacesRequest) (*Empty, error)
	AlterKeySpaces(context.Context, *AlterKeySpacesRequest) (*Empty, error)
	DropKeySpaces(context.Context, *DropKeySpacesRequest) (*Empty, error)
	ListKeySpaces(context.Context, *ListKeySpacesRequest) (*Empty, error)
	GetKeySpace(context.Context, *Empty) (*Empty, error)
	CreateTables(context.Context, *CreateTablesRequest) (*Empty, error)
	AlterTables(context.Context, *AlterTablesRequest) (*Empty, error)
	DropTables(context.Context, *DropTablesRequest) (*Empty, error)
	ListTables(context.Context, *ListTablesRequest) (*Empty, error)
	GetTable(context.Context, *Empty) (*Empty, error)
	mustEmbedUnimplementedDDLServer()
}

// UnimplementedDDLServer must be embedded to have forward compatible implementations.
type UnimplementedDDLServer struct {
}

func (UnimplementedDDLServer) CreateKeySpaces(context.Context, *CreateKeySpacesRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateKeySpaces not implemented")
}
func (UnimplementedDDLServer) AlterKeySpaces(context.Context, *AlterKeySpacesRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AlterKeySpaces not implemented")
}
func (UnimplementedDDLServer) DropKeySpaces(context.Context, *DropKeySpacesRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DropKeySpaces not implemented")
}
func (UnimplementedDDLServer) ListKeySpaces(context.Context, *ListKeySpacesRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListKeySpaces not implemented")
}
func (UnimplementedDDLServer) GetKeySpace(context.Context, *Empty) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetKeySpace not implemented")
}
func (UnimplementedDDLServer) CreateTables(context.Context, *CreateTablesRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateTables not implemented")
}
func (UnimplementedDDLServer) AlterTables(context.Context, *AlterTablesRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AlterTables not implemented")
}
func (UnimplementedDDLServer) DropTables(context.Context, *DropTablesRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DropTables not implemented")
}
func (UnimplementedDDLServer) ListTables(context.Context, *ListTablesRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListTables not implemented")
}
func (UnimplementedDDLServer) GetTable(context.Context, *Empty) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTable not implemented")
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

func _DDL_CreateKeySpaces_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateKeySpacesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DDLServer).CreateKeySpaces(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cql.grpc.v1.DDL/CreateKeySpaces",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DDLServer).CreateKeySpaces(ctx, req.(*CreateKeySpacesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DDL_AlterKeySpaces_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AlterKeySpacesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DDLServer).AlterKeySpaces(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cql.grpc.v1.DDL/AlterKeySpaces",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DDLServer).AlterKeySpaces(ctx, req.(*AlterKeySpacesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DDL_DropKeySpaces_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DropKeySpacesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DDLServer).DropKeySpaces(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cql.grpc.v1.DDL/DropKeySpaces",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DDLServer).DropKeySpaces(ctx, req.(*DropKeySpacesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DDL_ListKeySpaces_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListKeySpacesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DDLServer).ListKeySpaces(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cql.grpc.v1.DDL/ListKeySpaces",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DDLServer).ListKeySpaces(ctx, req.(*ListKeySpacesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DDL_GetKeySpace_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DDLServer).GetKeySpace(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cql.grpc.v1.DDL/GetKeySpace",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DDLServer).GetKeySpace(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _DDL_CreateTables_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateTablesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DDLServer).CreateTables(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cql.grpc.v1.DDL/CreateTables",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DDLServer).CreateTables(ctx, req.(*CreateTablesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DDL_AlterTables_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AlterTablesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DDLServer).AlterTables(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cql.grpc.v1.DDL/AlterTables",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DDLServer).AlterTables(ctx, req.(*AlterTablesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DDL_DropTables_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DropTablesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DDLServer).DropTables(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cql.grpc.v1.DDL/DropTables",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DDLServer).DropTables(ctx, req.(*DropTablesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DDL_ListTables_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListTablesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DDLServer).ListTables(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cql.grpc.v1.DDL/ListTables",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DDLServer).ListTables(ctx, req.(*ListTablesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DDL_GetTable_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DDLServer).GetTable(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cql.grpc.v1.DDL/GetTable",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DDLServer).GetTable(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// DDL_ServiceDesc is the grpc.ServiceDesc for DDL service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DDL_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "cql.grpc.v1.DDL",
	HandlerType: (*DDLServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateKeySpaces",
			Handler:    _DDL_CreateKeySpaces_Handler,
		},
		{
			MethodName: "AlterKeySpaces",
			Handler:    _DDL_AlterKeySpaces_Handler,
		},
		{
			MethodName: "DropKeySpaces",
			Handler:    _DDL_DropKeySpaces_Handler,
		},
		{
			MethodName: "ListKeySpaces",
			Handler:    _DDL_ListKeySpaces_Handler,
		},
		{
			MethodName: "GetKeySpace",
			Handler:    _DDL_GetKeySpace_Handler,
		},
		{
			MethodName: "CreateTables",
			Handler:    _DDL_CreateTables_Handler,
		},
		{
			MethodName: "AlterTables",
			Handler:    _DDL_AlterTables_Handler,
		},
		{
			MethodName: "DropTables",
			Handler:    _DDL_DropTables_Handler,
		},
		{
			MethodName: "ListTables",
			Handler:    _DDL_ListTables_Handler,
		},
		{
			MethodName: "GetTable",
			Handler:    _DDL_GetTable_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "ddl.proto",
}
