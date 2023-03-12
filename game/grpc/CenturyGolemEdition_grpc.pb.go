// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: grpc/CenturyGolemEdition.proto

package game_server

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

// GameClient is the client API for Game service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GameClient interface {
	// A simple RPC.
	//
	// Obtains the feature at a given position.
	//
	// A feature with an empty name is returned if there's no feature at the given
	// position.
	GetCard(ctx context.Context, in *GetCardMessage, opts ...grpc.CallOption) (*GemCard, error)
	NewGame(ctx context.Context, in *CreateGameMessage, opts ...grpc.CallOption) (*CreateGameResponse, error)
	PlayGemCard(ctx context.Context, in *PlayGemCardMessage, opts ...grpc.CallOption) (*PlayGemCardResponse, error)
}

type gameClient struct {
	cc grpc.ClientConnInterface
}

func NewGameClient(cc grpc.ClientConnInterface) GameClient {
	return &gameClient{cc}
}

func (c *gameClient) GetCard(ctx context.Context, in *GetCardMessage, opts ...grpc.CallOption) (*GemCard, error) {
	out := new(GemCard)
	err := c.cc.Invoke(ctx, "/game_server.Game/GetCard", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gameClient) NewGame(ctx context.Context, in *CreateGameMessage, opts ...grpc.CallOption) (*CreateGameResponse, error) {
	out := new(CreateGameResponse)
	err := c.cc.Invoke(ctx, "/game_server.Game/NewGame", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gameClient) PlayGemCard(ctx context.Context, in *PlayGemCardMessage, opts ...grpc.CallOption) (*PlayGemCardResponse, error) {
	out := new(PlayGemCardResponse)
	err := c.cc.Invoke(ctx, "/game_server.Game/PlayGemCard", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GameServer is the server API for Game service.
// All implementations must embed UnimplementedGameServer
// for forward compatibility
type GameServer interface {
	// A simple RPC.
	//
	// Obtains the feature at a given position.
	//
	// A feature with an empty name is returned if there's no feature at the given
	// position.
	GetCard(context.Context, *GetCardMessage) (*GemCard, error)
	NewGame(context.Context, *CreateGameMessage) (*CreateGameResponse, error)
	PlayGemCard(context.Context, *PlayGemCardMessage) (*PlayGemCardResponse, error)
	mustEmbedUnimplementedGameServer()
}

// UnimplementedGameServer must be embedded to have forward compatible implementations.
type UnimplementedGameServer struct {
}

func (UnimplementedGameServer) GetCard(context.Context, *GetCardMessage) (*GemCard, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCard not implemented")
}
func (UnimplementedGameServer) NewGame(context.Context, *CreateGameMessage) (*CreateGameResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NewGame not implemented")
}
func (UnimplementedGameServer) PlayGemCard(context.Context, *PlayGemCardMessage) (*PlayGemCardResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PlayGemCard not implemented")
}
func (UnimplementedGameServer) mustEmbedUnimplementedGameServer() {}

// UnsafeGameServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GameServer will
// result in compilation errors.
type UnsafeGameServer interface {
	mustEmbedUnimplementedGameServer()
}

func RegisterGameServer(s grpc.ServiceRegistrar, srv GameServer) {
	s.RegisterService(&Game_ServiceDesc, srv)
}

func _Game_GetCard_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCardMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GameServer).GetCard(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/game_server.Game/GetCard",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GameServer).GetCard(ctx, req.(*GetCardMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _Game_NewGame_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateGameMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GameServer).NewGame(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/game_server.Game/NewGame",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GameServer).NewGame(ctx, req.(*CreateGameMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _Game_PlayGemCard_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PlayGemCardMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GameServer).PlayGemCard(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/game_server.Game/PlayGemCard",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GameServer).PlayGemCard(ctx, req.(*PlayGemCardMessage))
	}
	return interceptor(ctx, in, info, handler)
}

// Game_ServiceDesc is the grpc.ServiceDesc for Game service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Game_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "game_server.Game",
	HandlerType: (*GameServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetCard",
			Handler:    _Game_GetCard_Handler,
		},
		{
			MethodName: "NewGame",
			Handler:    _Game_NewGame_Handler,
		},
		{
			MethodName: "PlayGemCard",
			Handler:    _Game_PlayGemCard_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "grpc/CenturyGolemEdition.proto",
}
