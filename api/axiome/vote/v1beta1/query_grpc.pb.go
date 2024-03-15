// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: axiome/vote/v1beta1/query.proto

package votev1beta1

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
	Query_History_FullMethodName     = "/axiome.vote.v1beta1.Query/History"
	Query_Government_FullMethodName  = "/axiome.vote.v1beta1.Query/Government"
	Query_Current_FullMethodName     = "/axiome.vote.v1beta1.Query/Current"
	Query_Params_FullMethodName      = "/axiome.vote.v1beta1.Query/Params"
	Query_Poll_FullMethodName        = "/axiome.vote.v1beta1.Query/Poll"
	Query_PollHistory_FullMethodName = "/axiome.vote.v1beta1.Query/PollHistory"
)

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type QueryClient interface {
	History(ctx context.Context, in *HistoryRequest, opts ...grpc.CallOption) (*HistoryResponse, error)
	Government(ctx context.Context, in *GovernmentRequest, opts ...grpc.CallOption) (*GovernmentResponse, error)
	Current(ctx context.Context, in *CurrentRequest, opts ...grpc.CallOption) (*CurrentResponse, error)
	Params(ctx context.Context, in *ParamsRequest, opts ...grpc.CallOption) (*ParamsResponse, error)
	Poll(ctx context.Context, in *PollRequest, opts ...grpc.CallOption) (*PollResponse, error)
	PollHistory(ctx context.Context, in *PollHistoryRequest, opts ...grpc.CallOption) (*PollHistoryResponse, error)
}

type queryClient struct {
	cc grpc.ClientConnInterface
}

func NewQueryClient(cc grpc.ClientConnInterface) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) History(ctx context.Context, in *HistoryRequest, opts ...grpc.CallOption) (*HistoryResponse, error) {
	out := new(HistoryResponse)
	err := c.cc.Invoke(ctx, Query_History_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) Government(ctx context.Context, in *GovernmentRequest, opts ...grpc.CallOption) (*GovernmentResponse, error) {
	out := new(GovernmentResponse)
	err := c.cc.Invoke(ctx, Query_Government_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) Current(ctx context.Context, in *CurrentRequest, opts ...grpc.CallOption) (*CurrentResponse, error) {
	out := new(CurrentResponse)
	err := c.cc.Invoke(ctx, Query_Current_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) Params(ctx context.Context, in *ParamsRequest, opts ...grpc.CallOption) (*ParamsResponse, error) {
	out := new(ParamsResponse)
	err := c.cc.Invoke(ctx, Query_Params_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) Poll(ctx context.Context, in *PollRequest, opts ...grpc.CallOption) (*PollResponse, error) {
	out := new(PollResponse)
	err := c.cc.Invoke(ctx, Query_Poll_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) PollHistory(ctx context.Context, in *PollHistoryRequest, opts ...grpc.CallOption) (*PollHistoryResponse, error) {
	out := new(PollHistoryResponse)
	err := c.cc.Invoke(ctx, Query_PollHistory_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
// All implementations must embed UnimplementedQueryServer
// for forward compatibility
type QueryServer interface {
	History(context.Context, *HistoryRequest) (*HistoryResponse, error)
	Government(context.Context, *GovernmentRequest) (*GovernmentResponse, error)
	Current(context.Context, *CurrentRequest) (*CurrentResponse, error)
	Params(context.Context, *ParamsRequest) (*ParamsResponse, error)
	Poll(context.Context, *PollRequest) (*PollResponse, error)
	PollHistory(context.Context, *PollHistoryRequest) (*PollHistoryResponse, error)
	mustEmbedUnimplementedQueryServer()
}

// UnimplementedQueryServer must be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (UnimplementedQueryServer) History(context.Context, *HistoryRequest) (*HistoryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method History not implemented")
}
func (UnimplementedQueryServer) Government(context.Context, *GovernmentRequest) (*GovernmentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Government not implemented")
}
func (UnimplementedQueryServer) Current(context.Context, *CurrentRequest) (*CurrentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Current not implemented")
}
func (UnimplementedQueryServer) Params(context.Context, *ParamsRequest) (*ParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Params not implemented")
}
func (UnimplementedQueryServer) Poll(context.Context, *PollRequest) (*PollResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Poll not implemented")
}
func (UnimplementedQueryServer) PollHistory(context.Context, *PollHistoryRequest) (*PollHistoryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PollHistory not implemented")
}
func (UnimplementedQueryServer) mustEmbedUnimplementedQueryServer() {}

// UnsafeQueryServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to QueryServer will
// result in compilation errors.
type UnsafeQueryServer interface {
	mustEmbedUnimplementedQueryServer()
}

func RegisterQueryServer(s grpc.ServiceRegistrar, srv QueryServer) {
	s.RegisterService(&Query_ServiceDesc, srv)
}

func _Query_History_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HistoryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).History(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_History_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).History(ctx, req.(*HistoryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Government_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GovernmentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Government(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_Government_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Government(ctx, req.(*GovernmentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Current_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CurrentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Current(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_Current_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Current(ctx, req.(*CurrentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Params_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ParamsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Params(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_Params_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Params(ctx, req.(*ParamsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Poll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PollRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Poll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_Poll_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Poll(ctx, req.(*PollRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_PollHistory_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PollHistoryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).PollHistory(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_PollHistory_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).PollHistory(ctx, req.(*PollHistoryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Query_ServiceDesc is the grpc.ServiceDesc for Query service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Query_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "axiome.vote.v1beta1.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "History",
			Handler:    _Query_History_Handler,
		},
		{
			MethodName: "Government",
			Handler:    _Query_Government_Handler,
		},
		{
			MethodName: "Current",
			Handler:    _Query_Current_Handler,
		},
		{
			MethodName: "Params",
			Handler:    _Query_Params_Handler,
		},
		{
			MethodName: "Poll",
			Handler:    _Query_Poll_Handler,
		},
		{
			MethodName: "PollHistory",
			Handler:    _Query_PollHistory_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "axiome/vote/v1beta1/query.proto",
}