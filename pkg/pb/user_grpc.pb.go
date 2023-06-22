// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.19.4
// source: user.proto

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
	User_GetUserInfo_FullMethodName          = "/User/GetUserInfo"
	User_CheckPassWord_FullMethodName        = "/User/CheckPassWord"
	User_Register_FullMethodName             = "/User/Register"
	User_GetUserInfoList_FullMethodName      = "/User/GetUserInfoList"
	User_CacheChangeUserCount_FullMethodName = "/User/CacheChangeUserCount"
	User_CacheGetAuthor_FullMethodName       = "/User/CacheGetAuthor"
)

// UserClient is the client API for User service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserClient interface {
	GetUserInfo(ctx context.Context, in *GetUserInfoRequest, opts ...grpc.CallOption) (*GetUserInfoResponse, error)
	CheckPassWord(ctx context.Context, in *CheckPassWordRequest, opts ...grpc.CallOption) (*CheckPassWordResponse, error)
	Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error)
	GetUserInfoList(ctx context.Context, in *GetUserInfoListRequest, opts ...grpc.CallOption) (*GetUserInfoListResponse, error)
	CacheChangeUserCount(ctx context.Context, in *CacheChangeUserCountReq, opts ...grpc.CallOption) (*CacheChangeUserCountRsp, error)
	CacheGetAuthor(ctx context.Context, in *CacheGetAuthorReq, opts ...grpc.CallOption) (*CacheGetAuthorRsp, error)
}

type userClient struct {
	cc grpc.ClientConnInterface
}

func NewUserClient(cc grpc.ClientConnInterface) UserClient {
	return &userClient{cc}
}

func (c *userClient) GetUserInfo(ctx context.Context, in *GetUserInfoRequest, opts ...grpc.CallOption) (*GetUserInfoResponse, error) {
	out := new(GetUserInfoResponse)
	err := c.cc.Invoke(ctx, User_GetUserInfo_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) CheckPassWord(ctx context.Context, in *CheckPassWordRequest, opts ...grpc.CallOption) (*CheckPassWordResponse, error) {
	out := new(CheckPassWordResponse)
	err := c.cc.Invoke(ctx, User_CheckPassWord_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error) {
	out := new(RegisterResponse)
	err := c.cc.Invoke(ctx, User_Register_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) GetUserInfoList(ctx context.Context, in *GetUserInfoListRequest, opts ...grpc.CallOption) (*GetUserInfoListResponse, error) {
	out := new(GetUserInfoListResponse)
	err := c.cc.Invoke(ctx, User_GetUserInfoList_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) CacheChangeUserCount(ctx context.Context, in *CacheChangeUserCountReq, opts ...grpc.CallOption) (*CacheChangeUserCountRsp, error) {
	out := new(CacheChangeUserCountRsp)
	err := c.cc.Invoke(ctx, User_CacheChangeUserCount_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) CacheGetAuthor(ctx context.Context, in *CacheGetAuthorReq, opts ...grpc.CallOption) (*CacheGetAuthorRsp, error) {
	out := new(CacheGetAuthorRsp)
	err := c.cc.Invoke(ctx, User_CacheGetAuthor_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserServer is the server API for User service.
// All implementations should embed UnimplementedUserServer
// for forward compatibility
type UserServer interface {
	GetUserInfo(context.Context, *GetUserInfoRequest) (*GetUserInfoResponse, error)
	CheckPassWord(context.Context, *CheckPassWordRequest) (*CheckPassWordResponse, error)
	Register(context.Context, *RegisterRequest) (*RegisterResponse, error)
	GetUserInfoList(context.Context, *GetUserInfoListRequest) (*GetUserInfoListResponse, error)
	CacheChangeUserCount(context.Context, *CacheChangeUserCountReq) (*CacheChangeUserCountRsp, error)
	CacheGetAuthor(context.Context, *CacheGetAuthorReq) (*CacheGetAuthorRsp, error)
}

// UnimplementedUserServer should be embedded to have forward compatible implementations.
type UnimplementedUserServer struct {
}

func (UnimplementedUserServer) GetUserInfo(context.Context, *GetUserInfoRequest) (*GetUserInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserInfo not implemented")
}
func (UnimplementedUserServer) CheckPassWord(context.Context, *CheckPassWordRequest) (*CheckPassWordResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckPassWord not implemented")
}
func (UnimplementedUserServer) Register(context.Context, *RegisterRequest) (*RegisterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (UnimplementedUserServer) GetUserInfoList(context.Context, *GetUserInfoListRequest) (*GetUserInfoListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserInfoList not implemented")
}
func (UnimplementedUserServer) CacheChangeUserCount(context.Context, *CacheChangeUserCountReq) (*CacheChangeUserCountRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CacheChangeUserCount not implemented")
}
func (UnimplementedUserServer) CacheGetAuthor(context.Context, *CacheGetAuthorReq) (*CacheGetAuthorRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CacheGetAuthor not implemented")
}

// UnsafeUserServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserServer will
// result in compilation errors.
type UnsafeUserServer interface {
	mustEmbedUnimplementedUserServer()
}

func RegisterUserServer(s grpc.ServiceRegistrar, srv UserServer) {
	s.RegisterService(&User_ServiceDesc, srv)
}

func _User_GetUserInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).GetUserInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: User_GetUserInfo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).GetUserInfo(ctx, req.(*GetUserInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_CheckPassWord_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckPassWordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).CheckPassWord(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: User_CheckPassWord_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).CheckPassWord(ctx, req.(*CheckPassWordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: User_Register_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).Register(ctx, req.(*RegisterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_GetUserInfoList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserInfoListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).GetUserInfoList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: User_GetUserInfoList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).GetUserInfoList(ctx, req.(*GetUserInfoListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_CacheChangeUserCount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CacheChangeUserCountReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).CacheChangeUserCount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: User_CacheChangeUserCount_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).CacheChangeUserCount(ctx, req.(*CacheChangeUserCountReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_CacheGetAuthor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CacheGetAuthorReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).CacheGetAuthor(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: User_CacheGetAuthor_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).CacheGetAuthor(ctx, req.(*CacheGetAuthorReq))
	}
	return interceptor(ctx, in, info, handler)
}

// User_ServiceDesc is the grpc.ServiceDesc for User service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var User_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "User",
	HandlerType: (*UserServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetUserInfo",
			Handler:    _User_GetUserInfo_Handler,
		},
		{
			MethodName: "CheckPassWord",
			Handler:    _User_CheckPassWord_Handler,
		},
		{
			MethodName: "Register",
			Handler:    _User_Register_Handler,
		},
		{
			MethodName: "GetUserInfoList",
			Handler:    _User_GetUserInfoList_Handler,
		},
		{
			MethodName: "CacheChangeUserCount",
			Handler:    _User_CacheChangeUserCount_Handler,
		},
		{
			MethodName: "CacheGetAuthor",
			Handler:    _User_CacheGetAuthor_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "user.proto",
}
