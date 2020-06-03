// Code generated by protoc-gen-go. DO NOT EDIT.
// source: transaction.proto

package transaction

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
	Nonce                uint64   `protobuf:"varint,1,opt,name=nonce,proto3" json:"nonce,omitempty"`
	Amount               []byte   `protobuf:"bytes,2,opt,name=amount,proto3" json:"amount,omitempty"`
	Address              string   `protobuf:"bytes,3,opt,name=address,proto3" json:"address,omitempty"`
	Address2             string   `protobuf:"bytes,4,opt,name=address2,proto3" json:"address2,omitempty"`
	TransactionHash      []string `protobuf:"bytes,5,rep,name=transactionHash,proto3" json:"transactionHash,omitempty"`
	GasLimit             uint64   `protobuf:"varint,6,opt,name=gasLimit,proto3" json:"gasLimit,omitempty"`
	GasPrice             uint64   `protobuf:"varint,7,opt,name=gasPrice,proto3" json:"gasPrice,omitempty"`
	Payload              []byte   `protobuf:"bytes,8,opt,name=payload,proto3" json:"payload,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GeneralRequest) Reset()         { *m = GeneralRequest{} }
func (m *GeneralRequest) String() string { return proto.CompactTextString(m) }
func (*GeneralRequest) ProtoMessage()    {}
func (*GeneralRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_2cc4e03d2c28c490, []int{0}
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

func (m *GeneralRequest) GetNonce() uint64 {
	if m != nil {
		return m.Nonce
	}
	return 0
}

func (m *GeneralRequest) GetAmount() []byte {
	if m != nil {
		return m.Amount
	}
	return nil
}

func (m *GeneralRequest) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *GeneralRequest) GetAddress2() string {
	if m != nil {
		return m.Address2
	}
	return ""
}

func (m *GeneralRequest) GetTransactionHash() []string {
	if m != nil {
		return m.TransactionHash
	}
	return nil
}

func (m *GeneralRequest) GetGasLimit() uint64 {
	if m != nil {
		return m.GasLimit
	}
	return 0
}

func (m *GeneralRequest) GetGasPrice() uint64 {
	if m != nil {
		return m.GasPrice
	}
	return 0
}

func (m *GeneralRequest) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
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
	return fileDescriptor_2cc4e03d2c28c490, []int{1}
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
	proto.RegisterType((*GeneralRequest)(nil), "transaction.GeneralRequest")
	proto.RegisterType((*GeneralResponse)(nil), "transaction.GeneralResponse")
}

func init() { proto.RegisterFile("transaction.proto", fileDescriptor_2cc4e03d2c28c490) }

var fileDescriptor_2cc4e03d2c28c490 = []byte{
	// 325 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x93, 0xcd, 0x6a, 0xf2, 0x40,
	0x14, 0x86, 0xbf, 0x7c, 0x6a, 0xd4, 0x63, 0x51, 0x3a, 0x2d, 0x65, 0xb0, 0x5d, 0x04, 0x57, 0x81,
	0x82, 0x0b, 0x7b, 0x09, 0xa5, 0x3f, 0x94, 0x2a, 0x32, 0x8a, 0xfb, 0x63, 0x9c, 0xc6, 0x81, 0x71,
	0xc6, 0xce, 0x4c, 0x28, 0xde, 0x4e, 0xaf, 0xb0, 0x97, 0x50, 0x9c, 0x68, 0x9a, 0x96, 0xee, 0xd2,
	0x5d, 0x9e, 0xf3, 0x86, 0x37, 0xcf, 0x39, 0x10, 0x38, 0x75, 0x06, 0x95, 0xc5, 0xc4, 0x09, 0xad,
	0x86, 0x5b, 0xa3, 0x9d, 0x26, 0x9d, 0xd2, 0x68, 0xf0, 0x11, 0x40, 0xf7, 0x81, 0x2b, 0x6e, 0x50,
	0x32, 0xfe, 0x9a, 0x71, 0xeb, 0xc8, 0x39, 0x34, 0x94, 0x56, 0x09, 0xa7, 0x41, 0x14, 0xc4, 0x75,
	0x96, 0x03, 0xb9, 0x80, 0x10, 0x37, 0x3a, 0x53, 0x8e, 0xfe, 0x8f, 0x82, 0xf8, 0x84, 0x1d, 0x88,
	0x50, 0x68, 0xe2, 0x6a, 0x65, 0xb8, 0xb5, 0xb4, 0x16, 0x05, 0x71, 0x9b, 0x1d, 0x91, 0xf4, 0xa1,
	0x75, 0x78, 0x1c, 0xd1, 0xba, 0x8f, 0x0a, 0x26, 0x31, 0xf4, 0x4a, 0x16, 0x8f, 0x68, 0xd7, 0xb4,
	0x11, 0xd5, 0xe2, 0x36, 0xfb, 0x39, 0xde, 0xb7, 0xa4, 0x68, 0x9f, 0xc5, 0x46, 0x38, 0x1a, 0x7a,
	0xa1, 0x82, 0x0f, 0xd9, 0xd4, 0x88, 0x84, 0xd3, 0x66, 0x91, 0x79, 0xde, 0x7b, 0x6d, 0x71, 0x27,
	0x35, 0xae, 0x68, 0xcb, 0x0b, 0x1f, 0x71, 0x70, 0x0d, 0xbd, 0x62, 0x63, 0xbb, 0xd5, 0xca, 0xfa,
	0x97, 0x37, 0xdc, 0x5a, 0x4c, 0xf3, 0xa5, 0xdb, 0xec, 0x88, 0xa3, 0xf7, 0x3a, 0x74, 0xe6, 0x5f,
	0x4a, 0x64, 0x0c, 0xdd, 0x09, 0x7f, 0x2b, 0x4f, 0x2e, 0x87, 0xe5, 0x13, 0x7f, 0xbf, 0x65, 0xff,
	0xea, 0xf7, 0x30, 0xff, 0xec, 0xe0, 0x1f, 0x61, 0x70, 0x76, 0x8b, 0x32, 0xc9, 0x24, 0x3a, 0x3e,
	0xd7, 0x0e, 0xe5, 0x02, 0x65, 0xc6, 0xab, 0x75, 0x4e, 0xa0, 0x37, 0x13, 0xa9, 0xfa, 0x33, 0xc7,
	0x7b, 0x68, 0x4e, 0xb3, 0xa5, 0x14, 0x76, 0x5d, 0xad, 0xe7, 0x09, 0x3a, 0x7b, 0xaf, 0x71, 0x7e,
	0xd9, 0x6a, 0x5d, 0x77, 0x10, 0x2e, 0xb8, 0x11, 0x2f, 0xbb, 0xca, 0x35, 0x33, 0x67, 0x84, 0x4a,
	0x2b, 0xd5, 0x2c, 0x43, 0xff, 0x63, 0xdd, 0x7c, 0x06, 0x00, 0x00, 0xff, 0xff, 0xfb, 0x8a, 0x42,
	0x1d, 0x6d, 0x03, 0x00, 0x00,
}