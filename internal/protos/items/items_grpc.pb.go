// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.3
// source: internal/protos/items/items.proto

package proto

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
	ItemsController_AddItem_FullMethodName      = "/items.ItemsController/AddItem"
	ItemsController_EditItem_FullMethodName     = "/items.ItemsController/EditItem"
	ItemsController_DeleteItem_FullMethodName   = "/items.ItemsController/DeleteItem"
	ItemsController_GetUserItems_FullMethodName = "/items.ItemsController/GetUserItems"
	ItemsController_TypesCounts_FullMethodName  = "/items.ItemsController/TypesCounts"
)

// ItemsControllerClient is the client API for ItemsController service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ItemsControllerClient interface {
	AddItem(ctx context.Context, in *AddItemRequest, opts ...grpc.CallOption) (*AddItemResponse, error)
	EditItem(ctx context.Context, in *EditItemRequest, opts ...grpc.CallOption) (*EditItemResponse, error)
	DeleteItem(ctx context.Context, in *DeleteItemRequest, opts ...grpc.CallOption) (*DeleteItemResponse, error)
	GetUserItems(ctx context.Context, in *GetUserItemsRequest, opts ...grpc.CallOption) (*GetUserItemsResponse, error)
	TypesCounts(ctx context.Context, in *TypesCountsRequest, opts ...grpc.CallOption) (*TypesCountsResponse, error)
}

type itemsControllerClient struct {
	cc grpc.ClientConnInterface
}

func NewItemsControllerClient(cc grpc.ClientConnInterface) ItemsControllerClient {
	return &itemsControllerClient{cc}
}

func (c *itemsControllerClient) AddItem(ctx context.Context, in *AddItemRequest, opts ...grpc.CallOption) (*AddItemResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AddItemResponse)
	err := c.cc.Invoke(ctx, ItemsController_AddItem_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *itemsControllerClient) EditItem(ctx context.Context, in *EditItemRequest, opts ...grpc.CallOption) (*EditItemResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(EditItemResponse)
	err := c.cc.Invoke(ctx, ItemsController_EditItem_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *itemsControllerClient) DeleteItem(ctx context.Context, in *DeleteItemRequest, opts ...grpc.CallOption) (*DeleteItemResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteItemResponse)
	err := c.cc.Invoke(ctx, ItemsController_DeleteItem_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *itemsControllerClient) GetUserItems(ctx context.Context, in *GetUserItemsRequest, opts ...grpc.CallOption) (*GetUserItemsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetUserItemsResponse)
	err := c.cc.Invoke(ctx, ItemsController_GetUserItems_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *itemsControllerClient) TypesCounts(ctx context.Context, in *TypesCountsRequest, opts ...grpc.CallOption) (*TypesCountsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TypesCountsResponse)
	err := c.cc.Invoke(ctx, ItemsController_TypesCounts_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ItemsControllerServer is the server API for ItemsController service.
// All implementations must embed UnimplementedItemsControllerServer
// for forward compatibility.
type ItemsControllerServer interface {
	AddItem(context.Context, *AddItemRequest) (*AddItemResponse, error)
	EditItem(context.Context, *EditItemRequest) (*EditItemResponse, error)
	DeleteItem(context.Context, *DeleteItemRequest) (*DeleteItemResponse, error)
	GetUserItems(context.Context, *GetUserItemsRequest) (*GetUserItemsResponse, error)
	TypesCounts(context.Context, *TypesCountsRequest) (*TypesCountsResponse, error)
	mustEmbedUnimplementedItemsControllerServer()
}

// UnimplementedItemsControllerServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedItemsControllerServer struct{}

func (UnimplementedItemsControllerServer) AddItem(context.Context, *AddItemRequest) (*AddItemResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddItem not implemented")
}
func (UnimplementedItemsControllerServer) EditItem(context.Context, *EditItemRequest) (*EditItemResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EditItem not implemented")
}
func (UnimplementedItemsControllerServer) DeleteItem(context.Context, *DeleteItemRequest) (*DeleteItemResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteItem not implemented")
}
func (UnimplementedItemsControllerServer) GetUserItems(context.Context, *GetUserItemsRequest) (*GetUserItemsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserItems not implemented")
}
func (UnimplementedItemsControllerServer) TypesCounts(context.Context, *TypesCountsRequest) (*TypesCountsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TypesCounts not implemented")
}
func (UnimplementedItemsControllerServer) mustEmbedUnimplementedItemsControllerServer() {}
func (UnimplementedItemsControllerServer) testEmbeddedByValue()                         {}

// UnsafeItemsControllerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ItemsControllerServer will
// result in compilation errors.
type UnsafeItemsControllerServer interface {
	mustEmbedUnimplementedItemsControllerServer()
}

func RegisterItemsControllerServer(s grpc.ServiceRegistrar, srv ItemsControllerServer) {
	// If the following call pancis, it indicates UnimplementedItemsControllerServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&ItemsController_ServiceDesc, srv)
}

func _ItemsController_AddItem_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddItemRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ItemsControllerServer).AddItem(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ItemsController_AddItem_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ItemsControllerServer).AddItem(ctx, req.(*AddItemRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ItemsController_EditItem_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EditItemRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ItemsControllerServer).EditItem(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ItemsController_EditItem_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ItemsControllerServer).EditItem(ctx, req.(*EditItemRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ItemsController_DeleteItem_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteItemRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ItemsControllerServer).DeleteItem(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ItemsController_DeleteItem_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ItemsControllerServer).DeleteItem(ctx, req.(*DeleteItemRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ItemsController_GetUserItems_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserItemsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ItemsControllerServer).GetUserItems(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ItemsController_GetUserItems_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ItemsControllerServer).GetUserItems(ctx, req.(*GetUserItemsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ItemsController_TypesCounts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TypesCountsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ItemsControllerServer).TypesCounts(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ItemsController_TypesCounts_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ItemsControllerServer).TypesCounts(ctx, req.(*TypesCountsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ItemsController_ServiceDesc is the grpc.ServiceDesc for ItemsController service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ItemsController_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "items.ItemsController",
	HandlerType: (*ItemsControllerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddItem",
			Handler:    _ItemsController_AddItem_Handler,
		},
		{
			MethodName: "EditItem",
			Handler:    _ItemsController_EditItem_Handler,
		},
		{
			MethodName: "DeleteItem",
			Handler:    _ItemsController_DeleteItem_Handler,
		},
		{
			MethodName: "GetUserItems",
			Handler:    _ItemsController_GetUserItems_Handler,
		},
		{
			MethodName: "TypesCounts",
			Handler:    _ItemsController_TypesCounts_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/protos/items/items.proto",
}
