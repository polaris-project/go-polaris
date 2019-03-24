// Code generated by protoc-gen-go. DO NOT EDIT.
// source: dag.proto

package dag

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
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
	TransactionHash      string   `protobuf:"bytes,3,opt,name=transactionHash,proto3" json:"transactionHash,omitempty"`
	Address              string   `protobuf:"bytes,4,opt,name=address,proto3" json:"address,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GeneralRequest) Reset()         { *m = GeneralRequest{} }
func (m *GeneralRequest) String() string { return proto.CompactTextString(m) }
func (*GeneralRequest) ProtoMessage()    {}
func (*GeneralRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_228b96b95413374c, []int{0}
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

func (m *GeneralRequest) GetTransactionHash() string {
	if m != nil {
		return m.TransactionHash
	}
	return ""
}

func (m *GeneralRequest) GetAddress() string {
	if m != nil {
		return m.Address
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
	return fileDescriptor_228b96b95413374c, []int{1}
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
	proto.RegisterType((*GeneralRequest)(nil), "dag.GeneralRequest")
	proto.RegisterType((*GeneralResponse)(nil), "dag.GeneralResponse")
}

func init() { proto.RegisterFile("dag.proto", fileDescriptor_228b96b95413374c) }

var fileDescriptor_228b96b95413374c = []byte{
	// 295 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x92, 0xcf, 0x4a, 0x3b, 0x31,
	0x10, 0xc7, 0x7f, 0xfd, 0xb5, 0x54, 0x3b, 0x82, 0x85, 0x58, 0x34, 0xf4, 0x24, 0x7b, 0x2a, 0x08,
	0x3d, 0x28, 0x5e, 0xbc, 0x75, 0x57, 0x5d, 0x2f, 0x8a, 0x54, 0x5f, 0x60, 0xdc, 0x8c, 0xbb, 0x4b,
	0x63, 0xb6, 0x66, 0x52, 0x4a, 0x5f, 0xc1, 0xa7, 0xf0, 0x51, 0x65, 0xff, 0xb4, 0xba, 0x7a, 0xca,
	0x2d, 0xdf, 0x99, 0xc9, 0x87, 0xf9, 0x84, 0xc0, 0x40, 0x61, 0x3a, 0x5d, 0xda, 0xc2, 0x15, 0xa2,
	0xab, 0x30, 0x0d, 0x3e, 0x3a, 0x70, 0x18, 0x93, 0x21, 0x8b, 0x7a, 0x4e, 0xef, 0x2b, 0x62, 0x27,
	0xc6, 0xb0, 0xff, 0x9a, 0x6b, 0x7a, 0x44, 0x97, 0xc9, 0xce, 0x69, 0x67, 0x32, 0x98, 0xef, 0xb2,
	0x90, 0xb0, 0x67, 0xc8, 0xad, 0x0b, 0xbb, 0x90, 0xff, 0xab, 0xd6, 0x36, 0x8a, 0x09, 0x0c, 0x9d,
	0x45, 0xc3, 0x98, 0xb8, 0xbc, 0x30, 0x77, 0xc8, 0x99, 0xec, 0x56, 0x13, 0xbf, 0xcb, 0x25, 0x03,
	0x95, 0xb2, 0xc4, 0x2c, 0x7b, 0x35, 0xa3, 0x89, 0xc1, 0x19, 0x0c, 0x77, 0xbb, 0xf0, 0xb2, 0x30,
	0x4c, 0xe5, 0xf0, 0x1b, 0x31, 0x63, 0x4a, 0xcd, 0x2e, 0xdb, 0x78, 0xfe, 0xd9, 0x83, 0xee, 0x35,
	0xa6, 0xe2, 0x12, 0xfa, 0x0f, 0xb4, 0x2e, 0x4f, 0x47, 0xd3, 0x52, 0xae, 0x6d, 0x33, 0x1e, 0xb5,
	0x8b, 0x35, 0x36, 0xf8, 0x27, 0xae, 0xe0, 0xe0, 0x1e, 0x17, 0x54, 0x36, 0x38, 0x67, 0xbf, 0xbb,
	0x11, 0x8c, 0x62, 0x72, 0xcf, 0xdf, 0x5e, 0xe1, 0xa6, 0x32, 0xf3, 0x82, 0xdc, 0xc0, 0x71, 0x1b,
	0x12, 0x65, 0xb9, 0x56, 0x96, 0x8c, 0x1f, 0x26, 0x06, 0xd9, 0xc6, 0x70, 0xb8, 0x99, 0xd5, 0xef,
	0xe9, 0x07, 0xba, 0x85, 0x93, 0x3f, 0xa0, 0x27, 0x32, 0x8a, 0xac, 0x1f, 0x67, 0x06, 0x22, 0x26,
	0x17, 0x12, 0xff, 0x64, 0x79, 0xaf, 0x12, 0xa1, 0x4e, 0x56, 0x1a, 0x1d, 0x35, 0x2e, 0x21, 0x6a,
	0x34, 0x09, 0x79, 0x71, 0x5e, 0xfa, 0xd5, 0x47, 0xbf, 0xf8, 0x0a, 0x00, 0x00, 0xff, 0xff, 0x21,
	0xa4, 0xa3, 0x66, 0xf5, 0x02, 0x00, 0x00,
}
