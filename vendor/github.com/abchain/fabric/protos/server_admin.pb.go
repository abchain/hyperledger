// Code generated by protoc-gen-go. DO NOT EDIT.
// source: server_admin.proto

package protos

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf1 "github.com/golang/protobuf/ptypes/empty"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type ServerStatus_StatusCode int32

const (
	ServerStatus_UNDEFINED ServerStatus_StatusCode = 0
	ServerStatus_STARTED   ServerStatus_StatusCode = 1
	ServerStatus_STOPPED   ServerStatus_StatusCode = 2
	ServerStatus_PAUSED    ServerStatus_StatusCode = 3
	ServerStatus_ERROR     ServerStatus_StatusCode = 4
	ServerStatus_UNKNOWN   ServerStatus_StatusCode = 5
)

var ServerStatus_StatusCode_name = map[int32]string{
	0: "UNDEFINED",
	1: "STARTED",
	2: "STOPPED",
	3: "PAUSED",
	4: "ERROR",
	5: "UNKNOWN",
}
var ServerStatus_StatusCode_value = map[string]int32{
	"UNDEFINED": 0,
	"STARTED":   1,
	"STOPPED":   2,
	"PAUSED":    3,
	"ERROR":     4,
	"UNKNOWN":   5,
}

func (x ServerStatus_StatusCode) String() string {
	return proto.EnumName(ServerStatus_StatusCode_name, int32(x))
}
func (ServerStatus_StatusCode) EnumDescriptor() ([]byte, []int) { return fileDescriptor7, []int{0, 0} }

type ServerStatus struct {
	Status ServerStatus_StatusCode `protobuf:"varint,1,opt,name=status,enum=protos.ServerStatus_StatusCode" json:"status,omitempty"`
}

func (m *ServerStatus) Reset()                    { *m = ServerStatus{} }
func (m *ServerStatus) String() string            { return proto.CompactTextString(m) }
func (*ServerStatus) ProtoMessage()               {}
func (*ServerStatus) Descriptor() ([]byte, []int) { return fileDescriptor7, []int{0} }

func (m *ServerStatus) GetStatus() ServerStatus_StatusCode {
	if m != nil {
		return m.Status
	}
	return ServerStatus_UNDEFINED
}

func init() {
	proto.RegisterType((*ServerStatus)(nil), "protos.ServerStatus")
	proto.RegisterEnum("protos.ServerStatus_StatusCode", ServerStatus_StatusCode_name, ServerStatus_StatusCode_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Admin service

type AdminClient interface {
	// Return the serve status.
	GetStatus(ctx context.Context, in *google_protobuf1.Empty, opts ...grpc.CallOption) (*ServerStatus, error)
	StartServer(ctx context.Context, in *google_protobuf1.Empty, opts ...grpc.CallOption) (*ServerStatus, error)
	StopServer(ctx context.Context, in *google_protobuf1.Empty, opts ...grpc.CallOption) (*ServerStatus, error)
}

type adminClient struct {
	cc *grpc.ClientConn
}

func NewAdminClient(cc *grpc.ClientConn) AdminClient {
	return &adminClient{cc}
}

func (c *adminClient) GetStatus(ctx context.Context, in *google_protobuf1.Empty, opts ...grpc.CallOption) (*ServerStatus, error) {
	out := new(ServerStatus)
	err := grpc.Invoke(ctx, "/protos.Admin/GetStatus", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminClient) StartServer(ctx context.Context, in *google_protobuf1.Empty, opts ...grpc.CallOption) (*ServerStatus, error) {
	out := new(ServerStatus)
	err := grpc.Invoke(ctx, "/protos.Admin/StartServer", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminClient) StopServer(ctx context.Context, in *google_protobuf1.Empty, opts ...grpc.CallOption) (*ServerStatus, error) {
	out := new(ServerStatus)
	err := grpc.Invoke(ctx, "/protos.Admin/StopServer", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Admin service

type AdminServer interface {
	// Return the serve status.
	GetStatus(context.Context, *google_protobuf1.Empty) (*ServerStatus, error)
	StartServer(context.Context, *google_protobuf1.Empty) (*ServerStatus, error)
	StopServer(context.Context, *google_protobuf1.Empty) (*ServerStatus, error)
}

func RegisterAdminServer(s *grpc.Server, srv AdminServer) {
	s.RegisterService(&_Admin_serviceDesc, srv)
}

func _Admin_GetStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(google_protobuf1.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServer).GetStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.Admin/GetStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServer).GetStatus(ctx, req.(*google_protobuf1.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Admin_StartServer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(google_protobuf1.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServer).StartServer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.Admin/StartServer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServer).StartServer(ctx, req.(*google_protobuf1.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Admin_StopServer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(google_protobuf1.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServer).StopServer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.Admin/StopServer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServer).StopServer(ctx, req.(*google_protobuf1.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _Admin_serviceDesc = grpc.ServiceDesc{
	ServiceName: "protos.Admin",
	HandlerType: (*AdminServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetStatus",
			Handler:    _Admin_GetStatus_Handler,
		},
		{
			MethodName: "StartServer",
			Handler:    _Admin_StartServer_Handler,
		},
		{
			MethodName: "StopServer",
			Handler:    _Admin_StopServer_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "server_admin.proto",
}

func init() { proto.RegisterFile("server_admin.proto", fileDescriptor7) }

var fileDescriptor7 = []byte{
	// 253 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x2a, 0x4e, 0x2d, 0x2a,
	0x4b, 0x2d, 0x8a, 0x4f, 0x4c, 0xc9, 0xcd, 0xcc, 0xd3, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62,
	0x03, 0x53, 0xc5, 0x52, 0xd2, 0xe9, 0xf9, 0xf9, 0xe9, 0x39, 0xa9, 0xfa, 0x60, 0x6e, 0x52, 0x69,
	0x9a, 0x7e, 0x6a, 0x6e, 0x41, 0x49, 0x25, 0x44, 0x91, 0xd2, 0x22, 0x46, 0x2e, 0x9e, 0x60, 0xb0,
	0xde, 0xe0, 0x92, 0xc4, 0x92, 0xd2, 0x62, 0x21, 0x73, 0x2e, 0xb6, 0x62, 0x30, 0x4b, 0x82, 0x51,
	0x81, 0x51, 0x83, 0xcf, 0x48, 0x1e, 0xa2, 0xb0, 0x58, 0x0f, 0x59, 0x95, 0x1e, 0x84, 0x72, 0xce,
	0x4f, 0x49, 0x0d, 0x82, 0x2a, 0x57, 0x8a, 0xe4, 0xe2, 0x42, 0x88, 0x0a, 0xf1, 0x72, 0x71, 0x86,
	0xfa, 0xb9, 0xb8, 0xba, 0x79, 0xfa, 0xb9, 0xba, 0x08, 0x30, 0x08, 0x71, 0x73, 0xb1, 0x07, 0x87,
	0x38, 0x06, 0x85, 0xb8, 0xba, 0x08, 0x30, 0x42, 0x38, 0xfe, 0x01, 0x01, 0xae, 0x2e, 0x02, 0x4c,
	0x42, 0x5c, 0x5c, 0x6c, 0x01, 0x8e, 0xa1, 0xc1, 0xae, 0x2e, 0x02, 0xcc, 0x42, 0x9c, 0x5c, 0xac,
	0xae, 0x41, 0x41, 0xfe, 0x41, 0x02, 0x2c, 0x20, 0x35, 0xa1, 0x7e, 0xde, 0x7e, 0xfe, 0xe1, 0x7e,
	0x02, 0xac, 0x46, 0x07, 0x19, 0xb9, 0x58, 0x1d, 0x41, 0x3e, 0x13, 0xb2, 0xe6, 0xe2, 0x74, 0x4f,
	0x2d, 0x81, 0x3a, 0x55, 0x4c, 0x0f, 0xe2, 0x33, 0x3d, 0x98, 0xcf, 0xf4, 0x5c, 0x41, 0x3e, 0x93,
	0x12, 0xc1, 0xe6, 0x64, 0x25, 0x06, 0x21, 0x5b, 0x2e, 0xee, 0xe0, 0x92, 0xc4, 0xa2, 0x12, 0x88,
	0x30, 0xc9, 0xda, 0x6d, 0x40, 0x1e, 0xcc, 0x2f, 0x20, 0x4f, 0x77, 0x12, 0x24, 0x36, 0x8c, 0x01,
	0x01, 0x00, 0x00, 0xff, 0xff, 0x5b, 0x2d, 0xff, 0xc2, 0xaa, 0x01, 0x00, 0x00,
}
