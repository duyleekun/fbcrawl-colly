// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// GrpcClient is the client API for Grpc service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GrpcClient interface {
	// Sends a greeting
	Init(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Pointer, error)
	FreeColly(ctx context.Context, in *Pointer, opts ...grpc.CallOption) (*Empty, error)
	Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error)
	LoginWithCookies(ctx context.Context, in *LoginWithCookiesRequest, opts ...grpc.CallOption) (*Empty, error)
	FetchGroupInfo(ctx context.Context, in *FetchGroupInfoRequest, opts ...grpc.CallOption) (*FacebookGroup, error)
	FetchUserInfo(ctx context.Context, in *FetchUserInfoRequest, opts ...grpc.CallOption) (*FacebookUser, error)
	FetchGroupFeed(ctx context.Context, in *FetchGroupFeedRequest, opts ...grpc.CallOption) (*FacebookPostList, error)
	FetchPost(ctx context.Context, in *FetchPostRequest, opts ...grpc.CallOption) (*FacebookPost, error)
	FetchContentImages(ctx context.Context, in *FetchContentImagesRequest, opts ...grpc.CallOption) (*FacebookImageList, error)
	FetchImageUrl(ctx context.Context, in *FetchImageUrlRequest, opts ...grpc.CallOption) (*FacebookImage, error)
}

type grpcClient struct {
	cc grpc.ClientConnInterface
}

func NewGrpcClient(cc grpc.ClientConnInterface) GrpcClient {
	return &grpcClient{cc}
}

func (c *grpcClient) Init(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Pointer, error) {
	out := new(Pointer)
	err := c.cc.Invoke(ctx, "/fbcrawl_colly.Grpc/Init", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *grpcClient) FreeColly(ctx context.Context, in *Pointer, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/fbcrawl_colly.Grpc/FreeColly", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *grpcClient) Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error) {
	out := new(LoginResponse)
	err := c.cc.Invoke(ctx, "/fbcrawl_colly.Grpc/Login", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *grpcClient) LoginWithCookies(ctx context.Context, in *LoginWithCookiesRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/fbcrawl_colly.Grpc/LoginWithCookies", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *grpcClient) FetchGroupInfo(ctx context.Context, in *FetchGroupInfoRequest, opts ...grpc.CallOption) (*FacebookGroup, error) {
	out := new(FacebookGroup)
	err := c.cc.Invoke(ctx, "/fbcrawl_colly.Grpc/FetchGroupInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *grpcClient) FetchUserInfo(ctx context.Context, in *FetchUserInfoRequest, opts ...grpc.CallOption) (*FacebookUser, error) {
	out := new(FacebookUser)
	err := c.cc.Invoke(ctx, "/fbcrawl_colly.Grpc/FetchUserInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *grpcClient) FetchGroupFeed(ctx context.Context, in *FetchGroupFeedRequest, opts ...grpc.CallOption) (*FacebookPostList, error) {
	out := new(FacebookPostList)
	err := c.cc.Invoke(ctx, "/fbcrawl_colly.Grpc/FetchGroupFeed", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *grpcClient) FetchPost(ctx context.Context, in *FetchPostRequest, opts ...grpc.CallOption) (*FacebookPost, error) {
	out := new(FacebookPost)
	err := c.cc.Invoke(ctx, "/fbcrawl_colly.Grpc/FetchPost", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *grpcClient) FetchContentImages(ctx context.Context, in *FetchContentImagesRequest, opts ...grpc.CallOption) (*FacebookImageList, error) {
	out := new(FacebookImageList)
	err := c.cc.Invoke(ctx, "/fbcrawl_colly.Grpc/FetchContentImages", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *grpcClient) FetchImageUrl(ctx context.Context, in *FetchImageUrlRequest, opts ...grpc.CallOption) (*FacebookImage, error) {
	out := new(FacebookImage)
	err := c.cc.Invoke(ctx, "/fbcrawl_colly.Grpc/FetchImageUrl", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GrpcServer is the server API for Grpc service.
// All implementations must embed UnimplementedGrpcServer
// for forward compatibility
type GrpcServer interface {
	// Sends a greeting
	Init(context.Context, *Empty) (*Pointer, error)
	FreeColly(context.Context, *Pointer) (*Empty, error)
	Login(context.Context, *LoginRequest) (*LoginResponse, error)
	LoginWithCookies(context.Context, *LoginWithCookiesRequest) (*Empty, error)
	FetchGroupInfo(context.Context, *FetchGroupInfoRequest) (*FacebookGroup, error)
	FetchUserInfo(context.Context, *FetchUserInfoRequest) (*FacebookUser, error)
	FetchGroupFeed(context.Context, *FetchGroupFeedRequest) (*FacebookPostList, error)
	FetchPost(context.Context, *FetchPostRequest) (*FacebookPost, error)
	FetchContentImages(context.Context, *FetchContentImagesRequest) (*FacebookImageList, error)
	FetchImageUrl(context.Context, *FetchImageUrlRequest) (*FacebookImage, error)
	mustEmbedUnimplementedGrpcServer()
}

// UnimplementedGrpcServer must be embedded to have forward compatible implementations.
type UnimplementedGrpcServer struct {
}

func (*UnimplementedGrpcServer) Init(context.Context, *Empty) (*Pointer, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Init not implemented")
}
func (*UnimplementedGrpcServer) FreeColly(context.Context, *Pointer) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FreeColly not implemented")
}
func (*UnimplementedGrpcServer) Login(context.Context, *LoginRequest) (*LoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (*UnimplementedGrpcServer) LoginWithCookies(context.Context, *LoginWithCookiesRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoginWithCookies not implemented")
}
func (*UnimplementedGrpcServer) FetchGroupInfo(context.Context, *FetchGroupInfoRequest) (*FacebookGroup, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FetchGroupInfo not implemented")
}
func (*UnimplementedGrpcServer) FetchUserInfo(context.Context, *FetchUserInfoRequest) (*FacebookUser, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FetchUserInfo not implemented")
}
func (*UnimplementedGrpcServer) FetchGroupFeed(context.Context, *FetchGroupFeedRequest) (*FacebookPostList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FetchGroupFeed not implemented")
}
func (*UnimplementedGrpcServer) FetchPost(context.Context, *FetchPostRequest) (*FacebookPost, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FetchPost not implemented")
}
func (*UnimplementedGrpcServer) FetchContentImages(context.Context, *FetchContentImagesRequest) (*FacebookImageList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FetchContentImages not implemented")
}
func (*UnimplementedGrpcServer) FetchImageUrl(context.Context, *FetchImageUrlRequest) (*FacebookImage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FetchImageUrl not implemented")
}
func (*UnimplementedGrpcServer) mustEmbedUnimplementedGrpcServer() {}

func RegisterGrpcServer(s *grpc.Server, srv GrpcServer) {
	s.RegisterService(&_Grpc_serviceDesc, srv)
}

func _Grpc_Init_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GrpcServer).Init(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fbcrawl_colly.Grpc/Init",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GrpcServer).Init(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Grpc_FreeColly_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Pointer)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GrpcServer).FreeColly(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fbcrawl_colly.Grpc/FreeColly",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GrpcServer).FreeColly(ctx, req.(*Pointer))
	}
	return interceptor(ctx, in, info, handler)
}

func _Grpc_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GrpcServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fbcrawl_colly.Grpc/Login",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GrpcServer).Login(ctx, req.(*LoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Grpc_LoginWithCookies_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginWithCookiesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GrpcServer).LoginWithCookies(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fbcrawl_colly.Grpc/LoginWithCookies",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GrpcServer).LoginWithCookies(ctx, req.(*LoginWithCookiesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Grpc_FetchGroupInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FetchGroupInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GrpcServer).FetchGroupInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fbcrawl_colly.Grpc/FetchGroupInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GrpcServer).FetchGroupInfo(ctx, req.(*FetchGroupInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Grpc_FetchUserInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FetchUserInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GrpcServer).FetchUserInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fbcrawl_colly.Grpc/FetchUserInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GrpcServer).FetchUserInfo(ctx, req.(*FetchUserInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Grpc_FetchGroupFeed_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FetchGroupFeedRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GrpcServer).FetchGroupFeed(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fbcrawl_colly.Grpc/FetchGroupFeed",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GrpcServer).FetchGroupFeed(ctx, req.(*FetchGroupFeedRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Grpc_FetchPost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FetchPostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GrpcServer).FetchPost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fbcrawl_colly.Grpc/FetchPost",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GrpcServer).FetchPost(ctx, req.(*FetchPostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Grpc_FetchContentImages_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FetchContentImagesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GrpcServer).FetchContentImages(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fbcrawl_colly.Grpc/FetchContentImages",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GrpcServer).FetchContentImages(ctx, req.(*FetchContentImagesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Grpc_FetchImageUrl_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FetchImageUrlRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GrpcServer).FetchImageUrl(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fbcrawl_colly.Grpc/FetchImageUrl",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GrpcServer).FetchImageUrl(ctx, req.(*FetchImageUrlRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Grpc_serviceDesc = grpc.ServiceDesc{
	ServiceName: "fbcrawl_colly.Grpc",
	HandlerType: (*GrpcServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Init",
			Handler:    _Grpc_Init_Handler,
		},
		{
			MethodName: "FreeColly",
			Handler:    _Grpc_FreeColly_Handler,
		},
		{
			MethodName: "Login",
			Handler:    _Grpc_Login_Handler,
		},
		{
			MethodName: "LoginWithCookies",
			Handler:    _Grpc_LoginWithCookies_Handler,
		},
		{
			MethodName: "FetchGroupInfo",
			Handler:    _Grpc_FetchGroupInfo_Handler,
		},
		{
			MethodName: "FetchUserInfo",
			Handler:    _Grpc_FetchUserInfo_Handler,
		},
		{
			MethodName: "FetchGroupFeed",
			Handler:    _Grpc_FetchGroupFeed_Handler,
		},
		{
			MethodName: "FetchPost",
			Handler:    _Grpc_FetchPost_Handler,
		},
		{
			MethodName: "FetchContentImages",
			Handler:    _Grpc_FetchContentImages_Handler,
		},
		{
			MethodName: "FetchImageUrl",
			Handler:    _Grpc_FetchImageUrl_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "fbcrawl.proto",
}
