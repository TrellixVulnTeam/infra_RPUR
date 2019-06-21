// Code generated by protoc-gen-go. DO NOT EDIT.
// source: infra/appengine/drone-queen/internal/config/config.proto

package config

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	duration "github.com/golang/protobuf/ptypes/duration"
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

// Config is the configuration data served by luci-config for this app.
type Config struct {
	// access_groups are the luci-auth groups controlling access to RPC endpoints.
	AccessGroups *AccessGroups `protobuf:"bytes,1,opt,name=access_groups,json=accessGroups,proto3" json:"access_groups,omitempty"`
	// assignment_duration is the duration before expiration for drone
	// assignments.
	AssignmentDuration *duration.Duration `protobuf:"bytes,2,opt,name=assignment_duration,json=assignmentDuration,proto3" json:"assignment_duration,omitempty"`
	// instance identifies which instance of the service this is.  For
	// example, this could be prod for the prod instance.
	Instance             string   `protobuf:"bytes,3,opt,name=instance,proto3" json:"instance,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Config) Reset()         { *m = Config{} }
func (m *Config) String() string { return proto.CompactTextString(m) }
func (*Config) ProtoMessage()    {}
func (*Config) Descriptor() ([]byte, []int) {
	return fileDescriptor_9bcfef40975c8024, []int{0}
}

func (m *Config) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Config.Unmarshal(m, b)
}
func (m *Config) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Config.Marshal(b, m, deterministic)
}
func (m *Config) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Config.Merge(m, src)
}
func (m *Config) XXX_Size() int {
	return xxx_messageInfo_Config.Size(m)
}
func (m *Config) XXX_DiscardUnknown() {
	xxx_messageInfo_Config.DiscardUnknown(m)
}

var xxx_messageInfo_Config proto.InternalMessageInfo

func (m *Config) GetAccessGroups() *AccessGroups {
	if m != nil {
		return m.AccessGroups
	}
	return nil
}

func (m *Config) GetAssignmentDuration() *duration.Duration {
	if m != nil {
		return m.AssignmentDuration
	}
	return nil
}

func (m *Config) GetInstance() string {
	if m != nil {
		return m.Instance
	}
	return ""
}

// AccessGroups holds access group configuration
type AccessGroups struct {
	// drones is the group for calling drone RPCs.
	Drones string `protobuf:"bytes,1,opt,name=drones,proto3" json:"drones,omitempty"`
	// inventory_providers is the group for calling inventory RPCs.
	InventoryProviders   string   `protobuf:"bytes,2,opt,name=inventory_providers,json=inventoryProviders,proto3" json:"inventory_providers,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AccessGroups) Reset()         { *m = AccessGroups{} }
func (m *AccessGroups) String() string { return proto.CompactTextString(m) }
func (*AccessGroups) ProtoMessage()    {}
func (*AccessGroups) Descriptor() ([]byte, []int) {
	return fileDescriptor_9bcfef40975c8024, []int{1}
}

func (m *AccessGroups) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AccessGroups.Unmarshal(m, b)
}
func (m *AccessGroups) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AccessGroups.Marshal(b, m, deterministic)
}
func (m *AccessGroups) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AccessGroups.Merge(m, src)
}
func (m *AccessGroups) XXX_Size() int {
	return xxx_messageInfo_AccessGroups.Size(m)
}
func (m *AccessGroups) XXX_DiscardUnknown() {
	xxx_messageInfo_AccessGroups.DiscardUnknown(m)
}

var xxx_messageInfo_AccessGroups proto.InternalMessageInfo

func (m *AccessGroups) GetDrones() string {
	if m != nil {
		return m.Drones
	}
	return ""
}

func (m *AccessGroups) GetInventoryProviders() string {
	if m != nil {
		return m.InventoryProviders
	}
	return ""
}

func init() {
	proto.RegisterType((*Config)(nil), "drone_queen.config.Config")
	proto.RegisterType((*AccessGroups)(nil), "drone_queen.config.AccessGroups")
}

func init() {
	proto.RegisterFile("infra/appengine/drone-queen/internal/config/config.proto", fileDescriptor_9bcfef40975c8024)
}

var fileDescriptor_9bcfef40975c8024 = []byte{
	// 266 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x8f, 0x4d, 0x4b, 0xc4, 0x30,
	0x10, 0x86, 0xa9, 0x42, 0xd9, 0x8d, 0xeb, 0x25, 0x0b, 0x52, 0xf7, 0x20, 0x65, 0x4f, 0x7b, 0x31,
	0x01, 0xbd, 0x78, 0xf5, 0x0b, 0xc1, 0x93, 0xf4, 0x22, 0x78, 0x29, 0xd9, 0x76, 0x1a, 0x02, 0xeb,
	0x4c, 0x4d, 0xd2, 0x05, 0x7f, 0x99, 0x7f, 0x4f, 0x98, 0xb4, 0xeb, 0x82, 0xa7, 0xf0, 0xe6, 0x7d,
	0x66, 0xf2, 0x44, 0xdc, 0x39, 0xec, 0xbc, 0xd1, 0xa6, 0xef, 0x01, 0xad, 0x43, 0xd0, 0xad, 0x27,
	0x84, 0xeb, 0xaf, 0x01, 0x00, 0xb5, 0xc3, 0x08, 0x1e, 0xcd, 0x4e, 0x37, 0x84, 0x9d, 0xb3, 0xe3,
	0xa1, 0x7a, 0x4f, 0x91, 0xa4, 0x64, 0xb2, 0x66, 0x52, 0xa5, 0x66, 0x75, 0x65, 0x89, 0xec, 0x0e,
	0x34, 0x13, 0xdb, 0xa1, 0xd3, 0xed, 0xe0, 0x4d, 0x74, 0x84, 0x69, 0x66, 0xfd, 0x93, 0x89, 0xfc,
	0x91, 0x51, 0xf9, 0x2c, 0xce, 0x4d, 0xd3, 0x40, 0x08, 0xb5, 0xf5, 0x34, 0xf4, 0xa1, 0xc8, 0xca,
	0x6c, 0x73, 0x76, 0x53, 0xaa, 0xff, 0x6b, 0xd5, 0x3d, 0x83, 0x2f, 0xcc, 0x55, 0x0b, 0x73, 0x94,
	0xe4, 0xab, 0x58, 0x9a, 0x10, 0x9c, 0xc5, 0x4f, 0xc0, 0x58, 0x4f, 0xcf, 0x15, 0x27, 0xbc, 0xec,
	0x52, 0x25, 0x1f, 0x35, 0xf9, 0xa8, 0xa7, 0x11, 0xa8, 0xe4, 0xdf, 0xd4, 0x74, 0x27, 0x57, 0x62,
	0xe6, 0x30, 0x44, 0x83, 0x0d, 0x14, 0xa7, 0x65, 0xb6, 0x99, 0x57, 0x87, 0xbc, 0x7e, 0x17, 0x8b,
	0x63, 0x0b, 0x79, 0x21, 0x72, 0x16, 0x4d, 0xde, 0xf3, 0x6a, 0x4c, 0x52, 0x8b, 0xa5, 0xc3, 0x3d,
	0x60, 0x24, 0xff, 0x5d, 0xf7, 0x9e, 0xf6, 0xae, 0x05, 0x1f, 0xd8, 0x67, 0x5e, 0xc9, 0x43, 0xf5,
	0x36, 0x35, 0x0f, 0xb3, 0x8f, 0x3c, 0xfd, 0x72, 0x9b, 0xb3, 0xe5, 0xed, 0x6f, 0x00, 0x00, 0x00,
	0xff, 0xff, 0xc4, 0x73, 0x85, 0x3d, 0x93, 0x01, 0x00, 0x00,
}