// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v6.30.2
// source: proto/asr.proto

package proto

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
	ASRService_SpeechToText_FullMethodName = "/asr.ASRService/SpeechToText"
)

// ASRServiceClient is the client API for ASRService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ASRServiceClient interface {
	// 语音转文本
	SpeechToText(ctx context.Context, in *SpeechToTextRequest, opts ...grpc.CallOption) (*SpeechToTextResponse, error)
}

type aSRServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewASRServiceClient(cc grpc.ClientConnInterface) ASRServiceClient {
	return &aSRServiceClient{cc}
}

func (c *aSRServiceClient) SpeechToText(ctx context.Context, in *SpeechToTextRequest, opts ...grpc.CallOption) (*SpeechToTextResponse, error) {
	out := new(SpeechToTextResponse)
	err := c.cc.Invoke(ctx, ASRService_SpeechToText_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ASRServiceServer is the server API for ASRService service.
// All implementations must embed UnimplementedASRServiceServer
// for forward compatibility
type ASRServiceServer interface {
	// 语音转文本
	SpeechToText(context.Context, *SpeechToTextRequest) (*SpeechToTextResponse, error)
	mustEmbedUnimplementedASRServiceServer()
}

// UnimplementedASRServiceServer must be embedded to have forward compatible implementations.
type UnimplementedASRServiceServer struct {
}

func (UnimplementedASRServiceServer) SpeechToText(context.Context, *SpeechToTextRequest) (*SpeechToTextResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SpeechToText not implemented")
}
func (UnimplementedASRServiceServer) mustEmbedUnimplementedASRServiceServer() {}

// UnsafeASRServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ASRServiceServer will
// result in compilation errors.
type UnsafeASRServiceServer interface {
	mustEmbedUnimplementedASRServiceServer()
}

func RegisterASRServiceServer(s grpc.ServiceRegistrar, srv ASRServiceServer) {
	s.RegisterService(&ASRService_ServiceDesc, srv)
}

func _ASRService_SpeechToText_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SpeechToTextRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ASRServiceServer).SpeechToText(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ASRService_SpeechToText_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ASRServiceServer).SpeechToText(ctx, req.(*SpeechToTextRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ASRService_ServiceDesc is the grpc.ServiceDesc for ASRService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ASRService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "asr.ASRService",
	HandlerType: (*ASRServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SpeechToText",
			Handler:    _ASRService_SpeechToText_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/asr.proto",
}
