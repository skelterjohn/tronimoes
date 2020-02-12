// Code generated by protoc-gen-go. DO NOT EDIT.
// source: tronimoes.proto

package proto

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
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type HelloRequest struct {
	Message              string   `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HelloRequest) Reset()         { *m = HelloRequest{} }
func (m *HelloRequest) String() string { return proto.CompactTextString(m) }
func (*HelloRequest) ProtoMessage()    {}
func (*HelloRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_786f716abeafd7fa, []int{0}
}

func (m *HelloRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HelloRequest.Unmarshal(m, b)
}
func (m *HelloRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HelloRequest.Marshal(b, m, deterministic)
}
func (m *HelloRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HelloRequest.Merge(m, src)
}
func (m *HelloRequest) XXX_Size() int {
	return xxx_messageInfo_HelloRequest.Size(m)
}
func (m *HelloRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_HelloRequest.DiscardUnknown(m)
}

var xxx_messageInfo_HelloRequest proto.InternalMessageInfo

func (m *HelloRequest) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type HelloResponse struct {
	Message              string   `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HelloResponse) Reset()         { *m = HelloResponse{} }
func (m *HelloResponse) String() string { return proto.CompactTextString(m) }
func (*HelloResponse) ProtoMessage()    {}
func (*HelloResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_786f716abeafd7fa, []int{1}
}

func (m *HelloResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HelloResponse.Unmarshal(m, b)
}
func (m *HelloResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HelloResponse.Marshal(b, m, deterministic)
}
func (m *HelloResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HelloResponse.Merge(m, src)
}
func (m *HelloResponse) XXX_Size() int {
	return xxx_messageInfo_HelloResponse.Size(m)
}
func (m *HelloResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_HelloResponse.DiscardUnknown(m)
}

var xxx_messageInfo_HelloResponse proto.InternalMessageInfo

func (m *HelloResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func init() {
	proto.RegisterType((*HelloRequest)(nil), "skelterjohn.tronimoes.HelloRequest")
	proto.RegisterType((*HelloResponse)(nil), "skelterjohn.tronimoes.HelloResponse")
}

func init() { proto.RegisterFile("tronimoes.proto", fileDescriptor_786f716abeafd7fa) }

var fileDescriptor_786f716abeafd7fa = []byte{
	// 147 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2f, 0x29, 0xca, 0xcf,
	0xcb, 0xcc, 0xcd, 0x4f, 0x2d, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x12, 0x2d, 0xce, 0x4e,
	0xcd, 0x29, 0x49, 0x2d, 0xca, 0xca, 0xcf, 0xc8, 0xd3, 0x83, 0x4b, 0x2a, 0x69, 0x70, 0xf1, 0x78,
	0xa4, 0xe6, 0xe4, 0xe4, 0x07, 0xa5, 0x16, 0x96, 0xa6, 0x16, 0x97, 0x08, 0x49, 0x70, 0xb1, 0xe7,
	0xa6, 0x16, 0x17, 0x27, 0xa6, 0xa7, 0x4a, 0x30, 0x2a, 0x30, 0x6a, 0x70, 0x06, 0xc1, 0xb8, 0x4a,
	0x9a, 0x5c, 0xbc, 0x50, 0x95, 0xc5, 0x05, 0xf9, 0x79, 0xc5, 0xa9, 0xb8, 0x95, 0x1a, 0x45, 0x71,
	0xb1, 0xb8, 0x27, 0xe6, 0xa6, 0x0a, 0x05, 0x71, 0xb1, 0x82, 0xb5, 0x08, 0x29, 0xeb, 0x61, 0xb5,
	0x5d, 0x0f, 0xd9, 0x6a, 0x29, 0x15, 0xfc, 0x8a, 0x20, 0xb6, 0x3a, 0xb1, 0x47, 0xb1, 0x82, 0x3d,
	0x94, 0xc4, 0x06, 0xa6, 0x8c, 0x01, 0x01, 0x00, 0x00, 0xff, 0xff, 0x73, 0x8a, 0x3f, 0x13, 0xea,
	0x00, 0x00, 0x00,
}
