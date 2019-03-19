// Code generated by protoc-gen-go. DO NOT EDIT.
// source: config.proto

package config

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type GeneralRequest struct {
	FilePath             string   `protobuf:"bytes,1,opt,name=filePath,proto3" json:"filePath,omitempty"`
	Network              string   `protobuf:"bytes,2,opt,name=network,proto3" json:"network,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GeneralRequest) Reset()         { *m = GeneralRequest{} }
func (m *GeneralRequest) String() string { return proto.CompactTextString(m) }
func (*GeneralRequest) ProtoMessage()    {}
func (*GeneralRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_3eaf2c85e69e9ea4, []int{0}
}

func (m *GeneralRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GeneralRequest.Unmarshal(m, b)
}
func (m *GeneralRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GeneralRequest.Marshal(b, m, deterministic)
}
func (m *GeneralRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GeneralRequest.Merge(m, src)
}
func (m *GeneralRequest) XXX_Size() int {
	return xxx_messageInfo_GeneralRequest.Size(m)
}
func (m *GeneralRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GeneralRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GeneralRequest proto.InternalMessageInfo

func (m *GeneralRequest) GetFilePath() string {
	if m != nil {
		return m.FilePath
	}
	return ""
}

func (m *GeneralRequest) GetNetwork() string {
	if m != nil {
		return m.Network
	}
	return ""
}

type GeneralResponse struct {
	Message              string   `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GeneralResponse) Reset()         { *m = GeneralResponse{} }
func (m *GeneralResponse) String() string { return proto.CompactTextString(m) }
func (*GeneralResponse) ProtoMessage()    {}
func (*GeneralResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_3eaf2c85e69e9ea4, []int{1}
}

func (m *GeneralResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GeneralResponse.Unmarshal(m, b)
}
func (m *GeneralResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GeneralResponse.Marshal(b, m, deterministic)
}
func (m *GeneralResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GeneralResponse.Merge(m, src)
}
func (m *GeneralResponse) XXX_Size() int {
	return xxx_messageInfo_GeneralResponse.Size(m)
}
func (m *GeneralResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GeneralResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GeneralResponse proto.InternalMessageInfo

func (m *GeneralResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func init() {
	proto.RegisterType((*GeneralRequest)(nil), "config.GeneralRequest")
	proto.RegisterType((*GeneralResponse)(nil), "config.GeneralResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ConfigClient is the client API for Config service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ConfigClient interface {
	NewDagConfig(ctx context.Context, in *GeneralRequest, opts ...grpc.CallOption) (*GeneralResponse, error)
	GetAllConfigs(ctx context.Context, in *GeneralRequest, opts ...grpc.CallOption) (*GeneralResponse, error)
	GetConfig(ctx context.Context, in *GeneralRequest, opts ...grpc.CallOption) (*GeneralResponse, error)
}

type configClient struct {
	cc *grpc.ClientConn
}

func NewConfigClient(cc *grpc.ClientConn) ConfigClient {
	return &configClient{cc}
}

func (c *configClient) NewDagConfig(ctx context.Context, in *GeneralRequest, opts ...grpc.CallOption) (*GeneralResponse, error) {
	out := new(GeneralResponse)
	err := c.cc.Invoke(ctx, "/config.Config/NewDagConfig", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *configClient) GetAllConfigs(ctx context.Context, in *GeneralRequest, opts ...grpc.CallOption) (*GeneralResponse, error) {
	out := new(GeneralResponse)
	err := c.cc.Invoke(ctx, "/config.Config/GetAllConfigs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *configClient) GetConfig(ctx context.Context, in *GeneralRequest, opts ...grpc.CallOption) (*GeneralResponse, error) {
	out := new(GeneralResponse)
	err := c.cc.Invoke(ctx, "/config.Config/GetConfig", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ConfigServer is the server API for Config service.
type ConfigServer interface {
	NewDagConfig(context.Context, *GeneralRequest) (*GeneralResponse, error)
	GetAllConfigs(context.Context, *GeneralRequest) (*GeneralResponse, error)
	GetConfig(context.Context, *GeneralRequest) (*GeneralResponse, error)
}

func RegisterConfigServer(s *grpc.Server, srv ConfigServer) {
	s.RegisterService(&_Config_serviceDesc, srv)
}

func _Config_NewDagConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GeneralRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConfigServer).NewDagConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/config.Config/NewDagConfig",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConfigServer).NewDagConfig(ctx, req.(*GeneralRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Config_GetAllConfigs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GeneralRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConfigServer).GetAllConfigs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/config.Config/GetAllConfigs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConfigServer).GetAllConfigs(ctx, req.(*GeneralRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Config_GetConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GeneralRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConfigServer).GetConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/config.Config/GetConfig",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConfigServer).GetConfig(ctx, req.(*GeneralRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Config_serviceDesc = grpc.ServiceDesc{
	ServiceName: "config.Config",
	HandlerType: (*ConfigServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "NewDagConfig",
			Handler:    _Config_NewDagConfig_Handler,
		},
		{
			MethodName: "GetAllConfigs",
			Handler:    _Config_GetAllConfigs_Handler,
		},
		{
			MethodName: "GetConfig",
			Handler:    _Config_GetConfig_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "config.proto",
}

func init() { proto.RegisterFile("config.proto", fileDescriptor_3eaf2c85e69e9ea4) }

var fileDescriptor_3eaf2c85e69e9ea4 = []byte{
	// 187 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x49, 0xce, 0xcf, 0x4b,
	0xcb, 0x4c, 0xd7, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x83, 0xf0, 0x94, 0xdc, 0xb8, 0xf8,
	0xdc, 0x53, 0xf3, 0x52, 0x8b, 0x12, 0x73, 0x82, 0x52, 0x0b, 0x4b, 0x53, 0x8b, 0x4b, 0x84, 0xa4,
	0xb8, 0x38, 0xd2, 0x32, 0x73, 0x52, 0x03, 0x12, 0x4b, 0x32, 0x24, 0x18, 0x15, 0x18, 0x35, 0x38,
	0x83, 0xe0, 0x7c, 0x21, 0x09, 0x2e, 0xf6, 0xbc, 0xd4, 0x92, 0xf2, 0xfc, 0xa2, 0x6c, 0x09, 0x26,
	0xb0, 0x14, 0x8c, 0xab, 0xa4, 0xcd, 0xc5, 0x0f, 0x37, 0xa7, 0xb8, 0x20, 0x3f, 0xaf, 0x38, 0x15,
	0xa4, 0x38, 0x37, 0xb5, 0xb8, 0x38, 0x31, 0x3d, 0x15, 0x6a, 0x0e, 0x8c, 0x6b, 0x74, 0x9e, 0x91,
	0x8b, 0xcd, 0x19, 0x6c, 0xbf, 0x90, 0x23, 0x17, 0x8f, 0x5f, 0x6a, 0xb9, 0x4b, 0x62, 0x3a, 0x94,
	0x2f, 0xa6, 0x07, 0x75, 0x26, 0xaa, 0xab, 0xa4, 0xc4, 0x31, 0xc4, 0x21, 0xb6, 0x28, 0x31, 0x08,
	0x39, 0x71, 0xf1, 0xba, 0xa7, 0x96, 0x38, 0xe6, 0xe4, 0x40, 0x8c, 0x28, 0x26, 0xc7, 0x0c, 0x3b,
	0x2e, 0x4e, 0xf7, 0xd4, 0x12, 0xb2, 0xdd, 0x90, 0xc4, 0x06, 0x0e, 0x55, 0x63, 0x40, 0x00, 0x00,
	0x00, 0xff, 0xff, 0xf2, 0xa9, 0x90, 0x2a, 0x65, 0x01, 0x00, 0x00,
}
