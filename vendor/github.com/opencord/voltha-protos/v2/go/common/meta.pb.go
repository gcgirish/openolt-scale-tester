// Code generated by protoc-gen-go. DO NOT EDIT.
// source: voltha_protos/meta.proto

package common

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
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

type Access int32

const (
	// read-write, stored attribute
	Access_CONFIG Access = 0
	// read-only field, stored with the model, covered by its hash
	Access_READ_ONLY Access = 1
	// A read-only attribute that is not stored in the model, not covered
	// by its hash, its value is filled real-time upon each request.
	Access_REAL_TIME Access = 2
)

var Access_name = map[int32]string{
	0: "CONFIG",
	1: "READ_ONLY",
	2: "REAL_TIME",
}

var Access_value = map[string]int32{
	"CONFIG":    0,
	"READ_ONLY": 1,
	"REAL_TIME": 2,
}

func (x Access) String() string {
	return proto.EnumName(Access_name, int32(x))
}

func (Access) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_96b320e8a67781f3, []int{0}
}

type ChildNode struct {
	Key                  string   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ChildNode) Reset()         { *m = ChildNode{} }
func (m *ChildNode) String() string { return proto.CompactTextString(m) }
func (*ChildNode) ProtoMessage()    {}
func (*ChildNode) Descriptor() ([]byte, []int) {
	return fileDescriptor_96b320e8a67781f3, []int{0}
}

func (m *ChildNode) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChildNode.Unmarshal(m, b)
}
func (m *ChildNode) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChildNode.Marshal(b, m, deterministic)
}
func (m *ChildNode) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChildNode.Merge(m, src)
}
func (m *ChildNode) XXX_Size() int {
	return xxx_messageInfo_ChildNode.Size(m)
}
func (m *ChildNode) XXX_DiscardUnknown() {
	xxx_messageInfo_ChildNode.DiscardUnknown(m)
}

var xxx_messageInfo_ChildNode proto.InternalMessageInfo

func (m *ChildNode) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

var E_ChildNode = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FieldOptions)(nil),
	ExtensionType: (*ChildNode)(nil),
	Field:         7761772,
	Name:          "voltha.child_node",
	Tag:           "bytes,7761772,opt,name=child_node",
	Filename:      "voltha_protos/meta.proto",
}

var E_Access = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FieldOptions)(nil),
	ExtensionType: (*Access)(nil),
	Field:         7761773,
	Name:          "voltha.access",
	Tag:           "varint,7761773,opt,name=access,enum=voltha.Access",
	Filename:      "voltha_protos/meta.proto",
}

func init() {
	proto.RegisterEnum("voltha.Access", Access_name, Access_value)
	proto.RegisterType((*ChildNode)(nil), "voltha.ChildNode")
	proto.RegisterExtension(E_ChildNode)
	proto.RegisterExtension(E_Access)
}

func init() { proto.RegisterFile("voltha_protos/meta.proto", fileDescriptor_96b320e8a67781f3) }

var fileDescriptor_96b320e8a67781f3 = []byte{
	// 273 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x90, 0x41, 0x4b, 0x84, 0x40,
	0x18, 0x86, 0xb3, 0x05, 0xc1, 0x2f, 0x5a, 0xcc, 0x93, 0x04, 0x0b, 0xd2, 0x69, 0x09, 0x9a, 0x09,
	0xbb, 0xed, 0x6d, 0xdb, 0x76, 0x6b, 0x61, 0x53, 0x90, 0x2e, 0x75, 0x11, 0x1d, 0x27, 0x1d, 0x52,
	0x3f, 0x71, 0x66, 0x17, 0xfa, 0xa9, 0x5d, 0xfa, 0x05, 0xf5, 0x1f, 0x42, 0x47, 0xbb, 0xee, 0xed,
	0x9d, 0x77, 0xde, 0x79, 0x78, 0x18, 0x70, 0x0f, 0x58, 0xaa, 0x22, 0x89, 0x9b, 0x16, 0x15, 0x4a,
	0x5a, 0x71, 0x95, 0x90, 0x3e, 0x3b, 0xa6, 0xbe, 0xb9, 0xf4, 0x72, 0xc4, 0xbc, 0xe4, 0xb4, 0x6f,
	0xd3, 0xfd, 0x3b, 0xcd, 0xb8, 0x64, 0xad, 0x68, 0x14, 0xb6, 0x7a, 0x79, 0x35, 0x03, 0x6b, 0x55,
	0x88, 0x32, 0x0b, 0x30, 0xe3, 0x8e, 0x0d, 0x93, 0x0f, 0xfe, 0xe9, 0x1a, 0x9e, 0x31, 0xb7, 0xa2,
	0x2e, 0x5e, 0xfb, 0x60, 0x2e, 0x19, 0xe3, 0x52, 0x3a, 0x00, 0xe6, 0x2a, 0x0c, 0x36, 0xdb, 0x47,
	0xfb, 0xc4, 0x39, 0x07, 0x2b, 0x5a, 0x2f, 0x1f, 0xe2, 0x30, 0xd8, 0xbd, 0xda, 0xc6, 0x70, 0xdc,
	0xc5, 0x2f, 0xdb, 0xe7, 0xb5, 0x7d, 0xba, 0x88, 0x00, 0x58, 0x87, 0x8c, 0xeb, 0x8e, 0x39, 0x23,
	0xda, 0x81, 0x8c, 0x0e, 0x64, 0x23, 0x78, 0x99, 0x85, 0x8d, 0x12, 0x58, 0x4b, 0xf7, 0xe7, 0xfb,
	0x6b, 0xe2, 0x19, 0xf3, 0x33, 0xff, 0x82, 0x68, 0x67, 0xf2, 0xaf, 0x13, 0x59, 0x6c, 0x8c, 0x8b,
	0x27, 0x30, 0x13, 0xed, 0x71, 0x84, 0xf7, 0xab, 0x79, 0x53, 0x7f, 0x3a, 0xf2, 0xb4, 0x7f, 0x34,
	0xbc, 0xbf, 0xbf, 0x7d, 0x23, 0xb9, 0x50, 0xc5, 0x3e, 0x25, 0x0c, 0x2b, 0x8a, 0x0d, 0xaf, 0x19,
	0xb6, 0x19, 0xd5, 0xe3, 0x9b, 0xe1, 0x2b, 0x0f, 0x3e, 0xcd, 0x91, 0x32, 0xac, 0x2a, 0xac, 0x53,
	0xb3, 0x2f, 0xef, 0xfe, 0x02, 0x00, 0x00, 0xff, 0xff, 0x1f, 0x75, 0x11, 0xd9, 0x6f, 0x01, 0x00,
	0x00,
}
