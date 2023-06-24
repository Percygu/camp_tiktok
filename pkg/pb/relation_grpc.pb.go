// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.10.0
// source: relation.proto

package pb

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
	RelationService_RelationAction_FullMethodName          = "/RelationService/RelationAction"
	RelationService_GetRelationFollowList_FullMethodName   = "/RelationService/GetRelationFollowList"
	RelationService_GetRelationFollowerList_FullMethodName = "/RelationService/GetRelationFollowerList"
)

// RelationServiceClient is the client API for RelationService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RelationServiceClient interface {
	RelationAction(ctx context.Context, in *RelationActionReq, opts ...grpc.CallOption) (*RelationActionRsp, error)
	GetRelationFollowList(ctx context.Context, in *GetRelationFollowListReq, opts ...grpc.CallOption) (*GetRelationFollowListRsp, error)
	GetRelationFollowerList(ctx context.Context, in *GetRelationFollowerListReq, opts ...grpc.CallOption) (*GetRelationFollowerListRsp, error)
}

type relationServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRelationServiceClient(cc grpc.ClientConnInterface) RelationServiceClient {
	return &relationServiceClient{cc}
}

func (c *relationServiceClient) RelationAction(ctx context.Context, in *RelationActionReq, opts ...grpc.CallOption) (*RelationActionRsp, error) {
	out := new(RelationActionRsp)
	err := c.cc.Invoke(ctx, RelationService_RelationAction_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *relationServiceClient) GetRelationFollowList(ctx context.Context, in *GetRelationFollowListReq, opts ...grpc.CallOption) (*GetRelationFollowListRsp, error) {
	out := new(GetRelationFollowListRsp)
	err := c.cc.Invoke(ctx, RelationService_GetRelationFollowList_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *relationServiceClient) GetRelationFollowerList(ctx context.Context, in *GetRelationFollowerListReq, opts ...grpc.CallOption) (*GetRelationFollowerListRsp, error) {
	out := new(GetRelationFollowerListRsp)
	err := c.cc.Invoke(ctx, RelationService_GetRelationFollowerList_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RelationServiceServer is the server API for RelationService service.
// All implementations should embed UnimplementedRelationServiceServer
// for forward compatibility
type RelationServiceServer interface {
	RelationAction(context.Context, *RelationActionReq) (*RelationActionRsp, error)
	GetRelationFollowList(context.Context, *GetRelationFollowListReq) (*GetRelationFollowListRsp, error)
	GetRelationFollowerList(context.Context, *GetRelationFollowerListReq) (*GetRelationFollowerListRsp, error)
}

// UnimplementedRelationServiceServer should be embedded to have forward compatible implementations.
type UnimplementedRelationServiceServer struct {
}

func (UnimplementedRelationServiceServer) RelationAction(context.Context, *RelationActionReq) (*RelationActionRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RelationAction not implemented")
}
func (UnimplementedRelationServiceServer) GetRelationFollowList(context.Context, *GetRelationFollowListReq) (*GetRelationFollowListRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRelationFollowList not implemented")
}
func (UnimplementedRelationServiceServer) GetRelationFollowerList(context.Context, *GetRelationFollowerListReq) (*GetRelationFollowerListRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRelationFollowerList not implemented")
}

// UnsafeRelationServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RelationServiceServer will
// result in compilation errors.
type UnsafeRelationServiceServer interface {
	mustEmbedUnimplementedRelationServiceServer()
}

func RegisterRelationServiceServer(s grpc.ServiceRegistrar, srv RelationServiceServer) {
	s.RegisterService(&RelationService_ServiceDesc, srv)
}

func _RelationService_RelationAction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RelationActionReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RelationServiceServer).RelationAction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RelationService_RelationAction_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RelationServiceServer).RelationAction(ctx, req.(*RelationActionReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _RelationService_GetRelationFollowList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRelationFollowListReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RelationServiceServer).GetRelationFollowList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RelationService_GetRelationFollowList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RelationServiceServer).GetRelationFollowList(ctx, req.(*GetRelationFollowListReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _RelationService_GetRelationFollowerList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRelationFollowerListReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RelationServiceServer).GetRelationFollowerList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RelationService_GetRelationFollowerList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RelationServiceServer).GetRelationFollowerList(ctx, req.(*GetRelationFollowerListReq))
	}
	return interceptor(ctx, in, info, handler)
}

// RelationService_ServiceDesc is the grpc.ServiceDesc for RelationService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RelationService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "RelationService",
	HandlerType: (*RelationServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RelationAction",
			Handler:    _RelationService_RelationAction_Handler,
		},
		{
			MethodName: "GetRelationFollowList",
			Handler:    _RelationService_GetRelationFollowList_Handler,
		},
		{
			MethodName: "GetRelationFollowerList",
			Handler:    _RelationService_GetRelationFollowerList_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "relation.proto",
}
