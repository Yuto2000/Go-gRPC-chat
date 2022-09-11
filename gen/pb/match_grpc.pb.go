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

// MatchServiceClient is the client API for MatchService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MatchServiceClient interface {
	JoinRoom(ctx context.Context, in *JoinRoomRequest, opts ...grpc.CallOption) (MatchService_JoinRoomClient, error)
}

type matchServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMatchServiceClient(cc grpc.ClientConnInterface) MatchServiceClient {
	return &matchServiceClient{cc}
}

func (c *matchServiceClient) JoinRoom(ctx context.Context, in *JoinRoomRequest, opts ...grpc.CallOption) (MatchService_JoinRoomClient, error) {
	stream, err := c.cc.NewStream(ctx, &MatchService_ServiceDesc.Streams[0], "/chat.MatchService/JoinRoom", opts...)
	if err != nil {
		return nil, err
	}
	x := &matchServiceJoinRoomClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type MatchService_JoinRoomClient interface {
	Recv() (*JoinRoomResponse, error)
	grpc.ClientStream
}

type matchServiceJoinRoomClient struct {
	grpc.ClientStream
}

func (x *matchServiceJoinRoomClient) Recv() (*JoinRoomResponse, error) {
	m := new(JoinRoomResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// MatchServiceServer is the server API for MatchService service.
// All implementations should embed UnimplementedMatchServiceServer
// for forward compatibility
type MatchServiceServer interface {
	JoinRoom(*JoinRoomRequest, MatchService_JoinRoomServer) error
}

// UnimplementedMatchServiceServer should be embedded to have forward compatible implementations.
type UnimplementedMatchServiceServer struct {
}

func (UnimplementedMatchServiceServer) JoinRoom(*JoinRoomRequest, MatchService_JoinRoomServer) error {
	return status.Errorf(codes.Unimplemented, "method JoinRoom not implemented")
}

// UnsafeMatchServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MatchServiceServer will
// result in compilation errors.
type UnsafeMatchServiceServer interface {
	mustEmbedUnimplementedMatchServiceServer()
}

func RegisterMatchServiceServer(s grpc.ServiceRegistrar, srv MatchServiceServer) {
	s.RegisterService(&MatchService_ServiceDesc, srv)
}

func _MatchService_JoinRoom_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(JoinRoomRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(MatchServiceServer).JoinRoom(m, &matchServiceJoinRoomServer{stream})
}

type MatchService_JoinRoomServer interface {
	Send(*JoinRoomResponse) error
	grpc.ServerStream
}

type matchServiceJoinRoomServer struct {
	grpc.ServerStream
}

func (x *matchServiceJoinRoomServer) Send(m *JoinRoomResponse) error {
	return x.ServerStream.SendMsg(m)
}

// MatchService_ServiceDesc is the grpc.ServiceDesc for MatchService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MatchService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "chat.MatchService",
	HandlerType: (*MatchServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "JoinRoom",
			Handler:       _MatchService_JoinRoom_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "proto/match.proto",
}
