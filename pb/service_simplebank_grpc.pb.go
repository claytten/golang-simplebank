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
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// SimplebankClient is the client API for Simplebank service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SimplebankClient interface {
	// users
	LoginUser(ctx context.Context, in *LoginUserRequest, opts ...grpc.CallOption) (*LoginUserResponse, error)
	CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*CreateUserResponse, error)
	GetUser(ctx context.Context, in *GetUserRequest, opts ...grpc.CallOption) (*GetUserResponse, error)
	UpdateProfile(ctx context.Context, in *UpdateProfileRequest, opts ...grpc.CallOption) (*UpdateProfileResponse, error)
	UpdatePassword(ctx context.Context, in *UpdatePasswordRequest, opts ...grpc.CallOption) (*UpdatePasswordResponse, error)
	RenewToken(ctx context.Context, in *RenewTokenRequest, opts ...grpc.CallOption) (*RenewTokenResponse, error)
	//Account
	CreateAccount(ctx context.Context, in *CreateAccountRequest, opts ...grpc.CallOption) (*CreateAccountResponse, error)
	GetAccount(ctx context.Context, in *GetAccountRequest, opts ...grpc.CallOption) (*GetAccountResponse, error)
	UpdateAccount(ctx context.Context, in *UpdateAccountRequest, opts ...grpc.CallOption) (*UpdateAccountResponse, error)
	DeleteAccount(ctx context.Context, in *DeleteAccountRequest, opts ...grpc.CallOption) (*DeleteAccountResponse, error)
	TransferTxAccount(ctx context.Context, in *TransferTxAccountRequest, opts ...grpc.CallOption) (*TransferTxAccountResponse, error)
}

type simplebankClient struct {
	cc grpc.ClientConnInterface
}

func NewSimplebankClient(cc grpc.ClientConnInterface) SimplebankClient {
	return &simplebankClient{cc}
}

func (c *simplebankClient) LoginUser(ctx context.Context, in *LoginUserRequest, opts ...grpc.CallOption) (*LoginUserResponse, error) {
	out := new(LoginUserResponse)
	err := c.cc.Invoke(ctx, "/pb.Simplebank/LoginUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *simplebankClient) CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*CreateUserResponse, error) {
	out := new(CreateUserResponse)
	err := c.cc.Invoke(ctx, "/pb.Simplebank/CreateUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *simplebankClient) GetUser(ctx context.Context, in *GetUserRequest, opts ...grpc.CallOption) (*GetUserResponse, error) {
	out := new(GetUserResponse)
	err := c.cc.Invoke(ctx, "/pb.Simplebank/GetUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *simplebankClient) UpdateProfile(ctx context.Context, in *UpdateProfileRequest, opts ...grpc.CallOption) (*UpdateProfileResponse, error) {
	out := new(UpdateProfileResponse)
	err := c.cc.Invoke(ctx, "/pb.Simplebank/UpdateProfile", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *simplebankClient) UpdatePassword(ctx context.Context, in *UpdatePasswordRequest, opts ...grpc.CallOption) (*UpdatePasswordResponse, error) {
	out := new(UpdatePasswordResponse)
	err := c.cc.Invoke(ctx, "/pb.Simplebank/UpdatePassword", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *simplebankClient) RenewToken(ctx context.Context, in *RenewTokenRequest, opts ...grpc.CallOption) (*RenewTokenResponse, error) {
	out := new(RenewTokenResponse)
	err := c.cc.Invoke(ctx, "/pb.Simplebank/RenewToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *simplebankClient) CreateAccount(ctx context.Context, in *CreateAccountRequest, opts ...grpc.CallOption) (*CreateAccountResponse, error) {
	out := new(CreateAccountResponse)
	err := c.cc.Invoke(ctx, "/pb.Simplebank/CreateAccount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *simplebankClient) GetAccount(ctx context.Context, in *GetAccountRequest, opts ...grpc.CallOption) (*GetAccountResponse, error) {
	out := new(GetAccountResponse)
	err := c.cc.Invoke(ctx, "/pb.Simplebank/GetAccount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *simplebankClient) UpdateAccount(ctx context.Context, in *UpdateAccountRequest, opts ...grpc.CallOption) (*UpdateAccountResponse, error) {
	out := new(UpdateAccountResponse)
	err := c.cc.Invoke(ctx, "/pb.Simplebank/UpdateAccount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *simplebankClient) DeleteAccount(ctx context.Context, in *DeleteAccountRequest, opts ...grpc.CallOption) (*DeleteAccountResponse, error) {
	out := new(DeleteAccountResponse)
	err := c.cc.Invoke(ctx, "/pb.Simplebank/DeleteAccount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *simplebankClient) TransferTxAccount(ctx context.Context, in *TransferTxAccountRequest, opts ...grpc.CallOption) (*TransferTxAccountResponse, error) {
	out := new(TransferTxAccountResponse)
	err := c.cc.Invoke(ctx, "/pb.Simplebank/TransferTxAccount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SimplebankServer is the server API for Simplebank service.
// All implementations must embed UnimplementedSimplebankServer
// for forward compatibility
type SimplebankServer interface {
	// users
	LoginUser(context.Context, *LoginUserRequest) (*LoginUserResponse, error)
	CreateUser(context.Context, *CreateUserRequest) (*CreateUserResponse, error)
	GetUser(context.Context, *GetUserRequest) (*GetUserResponse, error)
	UpdateProfile(context.Context, *UpdateProfileRequest) (*UpdateProfileResponse, error)
	UpdatePassword(context.Context, *UpdatePasswordRequest) (*UpdatePasswordResponse, error)
	RenewToken(context.Context, *RenewTokenRequest) (*RenewTokenResponse, error)
	//Account
	CreateAccount(context.Context, *CreateAccountRequest) (*CreateAccountResponse, error)
	GetAccount(context.Context, *GetAccountRequest) (*GetAccountResponse, error)
	UpdateAccount(context.Context, *UpdateAccountRequest) (*UpdateAccountResponse, error)
	DeleteAccount(context.Context, *DeleteAccountRequest) (*DeleteAccountResponse, error)
	TransferTxAccount(context.Context, *TransferTxAccountRequest) (*TransferTxAccountResponse, error)
	mustEmbedUnimplementedSimplebankServer()
}

// UnimplementedSimplebankServer must be embedded to have forward compatible implementations.
type UnimplementedSimplebankServer struct {
}

func (UnimplementedSimplebankServer) LoginUser(context.Context, *LoginUserRequest) (*LoginUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoginUser not implemented")
}
func (UnimplementedSimplebankServer) CreateUser(context.Context, *CreateUserRequest) (*CreateUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateUser not implemented")
}
func (UnimplementedSimplebankServer) GetUser(context.Context, *GetUserRequest) (*GetUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUser not implemented")
}
func (UnimplementedSimplebankServer) UpdateProfile(context.Context, *UpdateProfileRequest) (*UpdateProfileResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateProfile not implemented")
}
func (UnimplementedSimplebankServer) UpdatePassword(context.Context, *UpdatePasswordRequest) (*UpdatePasswordResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdatePassword not implemented")
}
func (UnimplementedSimplebankServer) RenewToken(context.Context, *RenewTokenRequest) (*RenewTokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RenewToken not implemented")
}
func (UnimplementedSimplebankServer) CreateAccount(context.Context, *CreateAccountRequest) (*CreateAccountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateAccount not implemented")
}
func (UnimplementedSimplebankServer) GetAccount(context.Context, *GetAccountRequest) (*GetAccountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAccount not implemented")
}
func (UnimplementedSimplebankServer) UpdateAccount(context.Context, *UpdateAccountRequest) (*UpdateAccountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateAccount not implemented")
}
func (UnimplementedSimplebankServer) DeleteAccount(context.Context, *DeleteAccountRequest) (*DeleteAccountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteAccount not implemented")
}
func (UnimplementedSimplebankServer) TransferTxAccount(context.Context, *TransferTxAccountRequest) (*TransferTxAccountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TransferTxAccount not implemented")
}
func (UnimplementedSimplebankServer) mustEmbedUnimplementedSimplebankServer() {}

// UnsafeSimplebankServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SimplebankServer will
// result in compilation errors.
type UnsafeSimplebankServer interface {
	mustEmbedUnimplementedSimplebankServer()
}

func RegisterSimplebankServer(s grpc.ServiceRegistrar, srv SimplebankServer) {
	s.RegisterService(&Simplebank_ServiceDesc, srv)
}

func _Simplebank_LoginUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SimplebankServer).LoginUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Simplebank/LoginUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SimplebankServer).LoginUser(ctx, req.(*LoginUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Simplebank_CreateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SimplebankServer).CreateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Simplebank/CreateUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SimplebankServer).CreateUser(ctx, req.(*CreateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Simplebank_GetUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SimplebankServer).GetUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Simplebank/GetUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SimplebankServer).GetUser(ctx, req.(*GetUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Simplebank_UpdateProfile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateProfileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SimplebankServer).UpdateProfile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Simplebank/UpdateProfile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SimplebankServer).UpdateProfile(ctx, req.(*UpdateProfileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Simplebank_UpdatePassword_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdatePasswordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SimplebankServer).UpdatePassword(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Simplebank/UpdatePassword",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SimplebankServer).UpdatePassword(ctx, req.(*UpdatePasswordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Simplebank_RenewToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RenewTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SimplebankServer).RenewToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Simplebank/RenewToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SimplebankServer).RenewToken(ctx, req.(*RenewTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Simplebank_CreateAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateAccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SimplebankServer).CreateAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Simplebank/CreateAccount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SimplebankServer).CreateAccount(ctx, req.(*CreateAccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Simplebank_GetAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SimplebankServer).GetAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Simplebank/GetAccount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SimplebankServer).GetAccount(ctx, req.(*GetAccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Simplebank_UpdateAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateAccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SimplebankServer).UpdateAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Simplebank/UpdateAccount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SimplebankServer).UpdateAccount(ctx, req.(*UpdateAccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Simplebank_DeleteAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteAccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SimplebankServer).DeleteAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Simplebank/DeleteAccount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SimplebankServer).DeleteAccount(ctx, req.(*DeleteAccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Simplebank_TransferTxAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TransferTxAccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SimplebankServer).TransferTxAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Simplebank/TransferTxAccount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SimplebankServer).TransferTxAccount(ctx, req.(*TransferTxAccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Simplebank_ServiceDesc is the grpc.ServiceDesc for Simplebank service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Simplebank_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.Simplebank",
	HandlerType: (*SimplebankServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "LoginUser",
			Handler:    _Simplebank_LoginUser_Handler,
		},
		{
			MethodName: "CreateUser",
			Handler:    _Simplebank_CreateUser_Handler,
		},
		{
			MethodName: "GetUser",
			Handler:    _Simplebank_GetUser_Handler,
		},
		{
			MethodName: "UpdateProfile",
			Handler:    _Simplebank_UpdateProfile_Handler,
		},
		{
			MethodName: "UpdatePassword",
			Handler:    _Simplebank_UpdatePassword_Handler,
		},
		{
			MethodName: "RenewToken",
			Handler:    _Simplebank_RenewToken_Handler,
		},
		{
			MethodName: "CreateAccount",
			Handler:    _Simplebank_CreateAccount_Handler,
		},
		{
			MethodName: "GetAccount",
			Handler:    _Simplebank_GetAccount_Handler,
		},
		{
			MethodName: "UpdateAccount",
			Handler:    _Simplebank_UpdateAccount_Handler,
		},
		{
			MethodName: "DeleteAccount",
			Handler:    _Simplebank_DeleteAccount_Handler,
		},
		{
			MethodName: "TransferTxAccount",
			Handler:    _Simplebank_TransferTxAccount_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service_simplebank.proto",
}
