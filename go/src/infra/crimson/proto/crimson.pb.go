// Code generated by protoc-gen-go.
// source: crimson.proto
// DO NOT EDIT!

/*
Package crimson is a generated protocol buffer package.

It is generated from these files:
	crimson.proto

It has these top-level messages:
	IPRanges
	IPRange
	IPRangeStatus
	IPRangeQuery
	IPRangeDelete
	IPRangeDeleteList
	Host
	HostList
	HostDelete
	HostDeleteList
	HostStatus
	HostQuery
*/
package crimson

import prpccommon "github.com/luci/luci-go/common/prpc"
import prpc "github.com/luci/luci-go/server/prpc"

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

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

type IPRanges struct {
	Ranges []*IPRange `protobuf:"bytes,1,rep,name=ranges" json:"ranges,omitempty"`
}

func (m *IPRanges) Reset()                    { *m = IPRanges{} }
func (m *IPRanges) String() string            { return proto.CompactTextString(m) }
func (*IPRanges) ProtoMessage()               {}
func (*IPRanges) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *IPRanges) GetRanges() []*IPRange {
	if m != nil {
		return m.Ranges
	}
	return nil
}

type IPRange struct {
	Site      string `protobuf:"bytes,1,opt,name=site" json:"site,omitempty"`
	VlanId    uint32 `protobuf:"varint,2,opt,name=vlan_id,json=vlanId" json:"vlan_id,omitempty"`
	StartIp   string `protobuf:"bytes,3,opt,name=start_ip,json=startIp" json:"start_ip,omitempty"`
	EndIp     string `protobuf:"bytes,4,opt,name=end_ip,json=endIp" json:"end_ip,omitempty"`
	VlanAlias string `protobuf:"bytes,5,opt,name=vlan_alias,json=vlanAlias" json:"vlan_alias,omitempty"`
}

func (m *IPRange) Reset()                    { *m = IPRange{} }
func (m *IPRange) String() string            { return proto.CompactTextString(m) }
func (*IPRange) ProtoMessage()               {}
func (*IPRange) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type IPRangeStatus struct {
	Error string `protobuf:"bytes,1,opt,name=error" json:"error,omitempty"`
}

func (m *IPRangeStatus) Reset()                    { *m = IPRangeStatus{} }
func (m *IPRangeStatus) String() string            { return proto.CompactTextString(m) }
func (*IPRangeStatus) ProtoMessage()               {}
func (*IPRangeStatus) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

type IPRangeQuery struct {
	Site      string `protobuf:"bytes,1,opt,name=site" json:"site,omitempty"`
	VlanId    uint32 `protobuf:"varint,2,opt,name=vlan_id,json=vlanId" json:"vlan_id,omitempty"`
	Limit     uint32 `protobuf:"varint,3,opt,name=limit" json:"limit,omitempty"`
	Ip        string `protobuf:"bytes,4,opt,name=ip" json:"ip,omitempty"`
	VlanAlias string `protobuf:"bytes,5,opt,name=vlan_alias,json=vlanAlias" json:"vlan_alias,omitempty"`
}

func (m *IPRangeQuery) Reset()                    { *m = IPRangeQuery{} }
func (m *IPRangeQuery) String() string            { return proto.CompactTextString(m) }
func (*IPRangeQuery) ProtoMessage()               {}
func (*IPRangeQuery) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

// Keep compatible with IPRangeQuery.
type IPRangeDelete struct {
	Site   string `protobuf:"bytes,1,opt,name=site" json:"site,omitempty"`
	VlanId uint32 `protobuf:"varint,2,opt,name=vlan_id,json=vlanId" json:"vlan_id,omitempty"`
}

func (m *IPRangeDelete) Reset()                    { *m = IPRangeDelete{} }
func (m *IPRangeDelete) String() string            { return proto.CompactTextString(m) }
func (*IPRangeDelete) ProtoMessage()               {}
func (*IPRangeDelete) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

type IPRangeDeleteList struct {
	Ranges []*IPRangeDelete `protobuf:"bytes,1,rep,name=ranges" json:"ranges,omitempty"`
}

func (m *IPRangeDeleteList) Reset()                    { *m = IPRangeDeleteList{} }
func (m *IPRangeDeleteList) String() string            { return proto.CompactTextString(m) }
func (*IPRangeDeleteList) ProtoMessage()               {}
func (*IPRangeDeleteList) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *IPRangeDeleteList) GetRanges() []*IPRangeDelete {
	if m != nil {
		return m.Ranges
	}
	return nil
}

type Host struct {
	Site      string `protobuf:"bytes,1,opt,name=site" json:"site,omitempty"`
	Hostname  string `protobuf:"bytes,2,opt,name=hostname" json:"hostname,omitempty"`
	MacAddr   string `protobuf:"bytes,3,opt,name=mac_addr,json=macAddr" json:"mac_addr,omitempty"`
	Ip        string `protobuf:"bytes,4,opt,name=ip" json:"ip,omitempty"`
	BootClass string `protobuf:"bytes,5,opt,name=boot_class,json=bootClass" json:"boot_class,omitempty"`
}

func (m *Host) Reset()                    { *m = Host{} }
func (m *Host) String() string            { return proto.CompactTextString(m) }
func (*Host) ProtoMessage()               {}
func (*Host) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

type HostList struct {
	Hosts []*Host `protobuf:"bytes,1,rep,name=hosts" json:"hosts,omitempty"`
}

func (m *HostList) Reset()                    { *m = HostList{} }
func (m *HostList) String() string            { return proto.CompactTextString(m) }
func (*HostList) ProtoMessage()               {}
func (*HostList) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *HostList) GetHosts() []*Host {
	if m != nil {
		return m.Hosts
	}
	return nil
}

// Use the same numbers and fields here than as the 'Host' message, so as to
// make it a compatible subset of it, in case we want to merge them in the
// future.
type HostDelete struct {
	Hostname string `protobuf:"bytes,2,opt,name=hostname" json:"hostname,omitempty"`
	MacAddr  string `protobuf:"bytes,3,opt,name=mac_addr,json=macAddr" json:"mac_addr,omitempty"`
}

func (m *HostDelete) Reset()                    { *m = HostDelete{} }
func (m *HostDelete) String() string            { return proto.CompactTextString(m) }
func (*HostDelete) ProtoMessage()               {}
func (*HostDelete) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

// Use the same numbers and fields here than as the 'Host' message, so as to
// make it a compatible subset of it, in case we want to merge them in the
// future.
type HostDeleteList struct {
	Hosts []*HostDelete `protobuf:"bytes,1,rep,name=hosts" json:"hosts,omitempty"`
}

func (m *HostDeleteList) Reset()                    { *m = HostDeleteList{} }
func (m *HostDeleteList) String() string            { return proto.CompactTextString(m) }
func (*HostDeleteList) ProtoMessage()               {}
func (*HostDeleteList) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

func (m *HostDeleteList) GetHosts() []*HostDelete {
	if m != nil {
		return m.Hosts
	}
	return nil
}

type HostStatus struct {
}

func (m *HostStatus) Reset()                    { *m = HostStatus{} }
func (m *HostStatus) String() string            { return proto.CompactTextString(m) }
func (*HostStatus) ProtoMessage()               {}
func (*HostStatus) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

type HostQuery struct {
	Limit     uint32 `protobuf:"varint,1,opt,name=limit" json:"limit,omitempty"`
	Site      string `protobuf:"bytes,2,opt,name=site" json:"site,omitempty"`
	Hostname  string `protobuf:"bytes,3,opt,name=hostname" json:"hostname,omitempty"`
	MacAddr   string `protobuf:"bytes,4,opt,name=mac_addr,json=macAddr" json:"mac_addr,omitempty"`
	Ip        string `protobuf:"bytes,5,opt,name=ip" json:"ip,omitempty"`
	BootClass string `protobuf:"bytes,6,opt,name=boot_class,json=bootClass" json:"boot_class,omitempty"`
}

func (m *HostQuery) Reset()                    { *m = HostQuery{} }
func (m *HostQuery) String() string            { return proto.CompactTextString(m) }
func (*HostQuery) ProtoMessage()               {}
func (*HostQuery) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{11} }

func init() {
	proto.RegisterType((*IPRanges)(nil), "crimson.IPRanges")
	proto.RegisterType((*IPRange)(nil), "crimson.IPRange")
	proto.RegisterType((*IPRangeStatus)(nil), "crimson.IPRangeStatus")
	proto.RegisterType((*IPRangeQuery)(nil), "crimson.IPRangeQuery")
	proto.RegisterType((*IPRangeDelete)(nil), "crimson.IPRangeDelete")
	proto.RegisterType((*IPRangeDeleteList)(nil), "crimson.IPRangeDeleteList")
	proto.RegisterType((*Host)(nil), "crimson.Host")
	proto.RegisterType((*HostList)(nil), "crimson.HostList")
	proto.RegisterType((*HostDelete)(nil), "crimson.HostDelete")
	proto.RegisterType((*HostDeleteList)(nil), "crimson.HostDeleteList")
	proto.RegisterType((*HostStatus)(nil), "crimson.HostStatus")
	proto.RegisterType((*HostQuery)(nil), "crimson.HostQuery")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion3

// Client API for Crimson service

type CrimsonClient interface {
	CreateIPRange(ctx context.Context, in *IPRange, opts ...grpc.CallOption) (*IPRangeStatus, error)
	ReadIPRange(ctx context.Context, in *IPRangeQuery, opts ...grpc.CallOption) (*IPRanges, error)
	DeleteIPRange(ctx context.Context, in *IPRangeDeleteList, opts ...grpc.CallOption) (*IPRangeStatus, error)
	CreateHost(ctx context.Context, in *HostList, opts ...grpc.CallOption) (*HostStatus, error)
	ReadHost(ctx context.Context, in *HostQuery, opts ...grpc.CallOption) (*HostList, error)
	DeleteHost(ctx context.Context, in *HostDeleteList, opts ...grpc.CallOption) (*HostStatus, error)
}
type crimsonPRPCClient struct {
	client *prpccommon.Client
}

func NewCrimsonPRPCClient(client *prpccommon.Client) CrimsonClient {
	return &crimsonPRPCClient{client}
}

func (c *crimsonPRPCClient) CreateIPRange(ctx context.Context, in *IPRange, opts ...grpc.CallOption) (*IPRangeStatus, error) {
	out := new(IPRangeStatus)
	err := c.client.Call(ctx, "crimson.Crimson", "CreateIPRange", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *crimsonPRPCClient) ReadIPRange(ctx context.Context, in *IPRangeQuery, opts ...grpc.CallOption) (*IPRanges, error) {
	out := new(IPRanges)
	err := c.client.Call(ctx, "crimson.Crimson", "ReadIPRange", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *crimsonPRPCClient) DeleteIPRange(ctx context.Context, in *IPRangeDeleteList, opts ...grpc.CallOption) (*IPRangeStatus, error) {
	out := new(IPRangeStatus)
	err := c.client.Call(ctx, "crimson.Crimson", "DeleteIPRange", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *crimsonPRPCClient) CreateHost(ctx context.Context, in *HostList, opts ...grpc.CallOption) (*HostStatus, error) {
	out := new(HostStatus)
	err := c.client.Call(ctx, "crimson.Crimson", "CreateHost", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *crimsonPRPCClient) ReadHost(ctx context.Context, in *HostQuery, opts ...grpc.CallOption) (*HostList, error) {
	out := new(HostList)
	err := c.client.Call(ctx, "crimson.Crimson", "ReadHost", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *crimsonPRPCClient) DeleteHost(ctx context.Context, in *HostDeleteList, opts ...grpc.CallOption) (*HostStatus, error) {
	out := new(HostStatus)
	err := c.client.Call(ctx, "crimson.Crimson", "DeleteHost", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type crimsonClient struct {
	cc *grpc.ClientConn
}

func NewCrimsonClient(cc *grpc.ClientConn) CrimsonClient {
	return &crimsonClient{cc}
}

func (c *crimsonClient) CreateIPRange(ctx context.Context, in *IPRange, opts ...grpc.CallOption) (*IPRangeStatus, error) {
	out := new(IPRangeStatus)
	err := grpc.Invoke(ctx, "/crimson.Crimson/CreateIPRange", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *crimsonClient) ReadIPRange(ctx context.Context, in *IPRangeQuery, opts ...grpc.CallOption) (*IPRanges, error) {
	out := new(IPRanges)
	err := grpc.Invoke(ctx, "/crimson.Crimson/ReadIPRange", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *crimsonClient) DeleteIPRange(ctx context.Context, in *IPRangeDeleteList, opts ...grpc.CallOption) (*IPRangeStatus, error) {
	out := new(IPRangeStatus)
	err := grpc.Invoke(ctx, "/crimson.Crimson/DeleteIPRange", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *crimsonClient) CreateHost(ctx context.Context, in *HostList, opts ...grpc.CallOption) (*HostStatus, error) {
	out := new(HostStatus)
	err := grpc.Invoke(ctx, "/crimson.Crimson/CreateHost", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *crimsonClient) ReadHost(ctx context.Context, in *HostQuery, opts ...grpc.CallOption) (*HostList, error) {
	out := new(HostList)
	err := grpc.Invoke(ctx, "/crimson.Crimson/ReadHost", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *crimsonClient) DeleteHost(ctx context.Context, in *HostDeleteList, opts ...grpc.CallOption) (*HostStatus, error) {
	out := new(HostStatus)
	err := grpc.Invoke(ctx, "/crimson.Crimson/DeleteHost", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Crimson service

type CrimsonServer interface {
	CreateIPRange(context.Context, *IPRange) (*IPRangeStatus, error)
	ReadIPRange(context.Context, *IPRangeQuery) (*IPRanges, error)
	DeleteIPRange(context.Context, *IPRangeDeleteList) (*IPRangeStatus, error)
	CreateHost(context.Context, *HostList) (*HostStatus, error)
	ReadHost(context.Context, *HostQuery) (*HostList, error)
	DeleteHost(context.Context, *HostDeleteList) (*HostStatus, error)
}

func RegisterCrimsonServer(s prpc.Registrar, srv CrimsonServer) {
	s.RegisterService(&_Crimson_serviceDesc, srv)
}

func _Crimson_CreateIPRange_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IPRange)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CrimsonServer).CreateIPRange(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/crimson.Crimson/CreateIPRange",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrimsonServer).CreateIPRange(ctx, req.(*IPRange))
	}
	return interceptor(ctx, in, info, handler)
}

func _Crimson_ReadIPRange_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IPRangeQuery)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CrimsonServer).ReadIPRange(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/crimson.Crimson/ReadIPRange",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrimsonServer).ReadIPRange(ctx, req.(*IPRangeQuery))
	}
	return interceptor(ctx, in, info, handler)
}

func _Crimson_DeleteIPRange_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IPRangeDeleteList)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CrimsonServer).DeleteIPRange(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/crimson.Crimson/DeleteIPRange",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrimsonServer).DeleteIPRange(ctx, req.(*IPRangeDeleteList))
	}
	return interceptor(ctx, in, info, handler)
}

func _Crimson_CreateHost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HostList)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CrimsonServer).CreateHost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/crimson.Crimson/CreateHost",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrimsonServer).CreateHost(ctx, req.(*HostList))
	}
	return interceptor(ctx, in, info, handler)
}

func _Crimson_ReadHost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HostQuery)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CrimsonServer).ReadHost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/crimson.Crimson/ReadHost",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrimsonServer).ReadHost(ctx, req.(*HostQuery))
	}
	return interceptor(ctx, in, info, handler)
}

func _Crimson_DeleteHost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HostDeleteList)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CrimsonServer).DeleteHost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/crimson.Crimson/DeleteHost",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrimsonServer).DeleteHost(ctx, req.(*HostDeleteList))
	}
	return interceptor(ctx, in, info, handler)
}

var _Crimson_serviceDesc = grpc.ServiceDesc{
	ServiceName: "crimson.Crimson",
	HandlerType: (*CrimsonServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateIPRange",
			Handler:    _Crimson_CreateIPRange_Handler,
		},
		{
			MethodName: "ReadIPRange",
			Handler:    _Crimson_ReadIPRange_Handler,
		},
		{
			MethodName: "DeleteIPRange",
			Handler:    _Crimson_DeleteIPRange_Handler,
		},
		{
			MethodName: "CreateHost",
			Handler:    _Crimson_CreateHost_Handler,
		},
		{
			MethodName: "ReadHost",
			Handler:    _Crimson_ReadHost_Handler,
		},
		{
			MethodName: "DeleteHost",
			Handler:    _Crimson_DeleteHost_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: fileDescriptor0,
}

func init() { proto.RegisterFile("crimson.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 519 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x9c, 0x54, 0xdb, 0x6e, 0xd3, 0x40,
	0x10, 0xc5, 0x49, 0x7c, 0xc9, 0xb4, 0xae, 0xc8, 0x42, 0x69, 0x88, 0x84, 0x84, 0x8c, 0x90, 0xc2,
	0x4b, 0x90, 0x5a, 0x84, 0x84, 0xda, 0x97, 0x2a, 0x20, 0x11, 0x89, 0x07, 0x30, 0x1f, 0x10, 0x6d,
	0xe3, 0x15, 0x58, 0xf2, 0x25, 0x5a, 0x6f, 0x91, 0x78, 0xe3, 0x8d, 0x8f, 0xe0, 0x33, 0xf9, 0x01,
	0x76, 0x67, 0x2f, 0x60, 0xd7, 0x0d, 0x4a, 0xdf, 0xf6, 0xcc, 0xcc, 0xd9, 0x39, 0x7b, 0x66, 0x6c,
	0x88, 0x37, 0x3c, 0x2f, 0x9b, 0xba, 0x5a, 0x6c, 0x79, 0x2d, 0x6a, 0x12, 0x1a, 0x98, 0xbc, 0x82,
	0x68, 0xf5, 0x31, 0xa5, 0xd5, 0x17, 0xd6, 0x90, 0x39, 0x04, 0x1c, 0x4f, 0x53, 0xef, 0xe9, 0x70,
	0x7e, 0x70, 0x7a, 0x7f, 0x61, 0x49, 0xa6, 0x24, 0x35, 0xf9, 0xe4, 0xa7, 0x07, 0xa1, 0x89, 0x11,
	0x02, 0xa3, 0x26, 0x17, 0x4c, 0x72, 0xbc, 0xf9, 0x38, 0xc5, 0x33, 0x39, 0x81, 0xf0, 0x5b, 0x41,
	0xab, 0x75, 0x9e, 0x4d, 0x07, 0x32, 0x1c, 0xa7, 0x81, 0x82, 0xab, 0x8c, 0x3c, 0x86, 0xa8, 0x11,
	0x94, 0x8b, 0x75, 0xbe, 0x9d, 0x0e, 0x91, 0x10, 0x22, 0x5e, 0x6d, 0xc9, 0x31, 0x04, 0xac, 0xca,
	0x54, 0x62, 0x84, 0x09, 0x5f, 0x22, 0x19, 0x7e, 0x02, 0x80, 0x57, 0xd1, 0x22, 0xa7, 0xcd, 0xd4,
	0xc7, 0xd4, 0x58, 0x45, 0x2e, 0x55, 0x20, 0x79, 0x0e, 0xb1, 0x11, 0xf2, 0x59, 0x50, 0x71, 0xdd,
	0x90, 0x87, 0xe0, 0x33, 0xce, 0x6b, 0x6e, 0xf4, 0x68, 0x90, 0xfc, 0xf0, 0xe0, 0xd0, 0xd4, 0x7d,
	0xba, 0x66, 0xfc, 0xfb, 0x7e, 0xaa, 0xe5, 0x9d, 0x45, 0x5e, 0xe6, 0x02, 0x25, 0xc7, 0xa9, 0x06,
	0xe4, 0x08, 0x06, 0x4e, 0xac, 0x3c, 0xfd, 0x4f, 0xe9, 0x85, 0x53, 0xfa, 0x96, 0x15, 0x4c, 0xec,
	0x67, 0x5c, 0xb2, 0x84, 0x49, 0x8b, 0xfd, 0x21, 0x6f, 0x04, 0x59, 0x74, 0x06, 0xf6, 0xa8, 0x3b,
	0x30, 0x5d, 0xeb, 0xc6, 0x26, 0x5d, 0x18, 0xbd, 0xaf, 0x25, 0xb1, 0xaf, 0xf5, 0x0c, 0xa2, 0xaf,
	0x32, 0x57, 0xd1, 0x92, 0x61, 0xef, 0x71, 0xea, 0xb0, 0x1a, 0x5b, 0x49, 0x37, 0x6b, 0x9a, 0x65,
	0xdc, 0x8e, 0x4d, 0xe2, 0x4b, 0x09, 0xfb, 0x5c, 0xb8, 0xaa, 0x6b, 0xb1, 0xde, 0x14, 0xb4, 0x71,
	0x2e, 0xa8, 0xc8, 0x52, 0x05, 0x92, 0x97, 0x10, 0x29, 0x05, 0x28, 0xff, 0x19, 0xf8, 0xaa, 0x83,
	0x55, 0x1f, 0x3b, 0xf5, 0xaa, 0x22, 0xd5, 0x39, 0xf9, 0x70, 0x50, 0xd0, 0x78, 0x76, 0x37, 0x91,
	0xc9, 0x39, 0x1c, 0xfd, 0xbd, 0x04, 0x7b, 0xbf, 0x68, 0xf7, 0x7e, 0xd0, 0xea, 0x6d, 0x6c, 0x33,
	0x0a, 0x0e, 0xb5, 0x02, 0xbd, 0x5f, 0xc9, 0x2f, 0x0f, 0xc6, 0x0a, 0xea, 0x35, 0x72, 0x9b, 0xe1,
	0xfd, 0xbb, 0x19, 0xd6, 0xde, 0xc1, 0x2d, 0xf6, 0x0e, 0x77, 0x28, 0x1f, 0xf5, 0xd9, 0xeb, 0xdf,
	0x62, 0x6f, 0xd0, 0xb1, 0xf7, 0xf4, 0xf7, 0x00, 0xc2, 0xa5, 0x7e, 0x09, 0x39, 0x87, 0x78, 0xc9,
	0x19, 0x15, 0xcc, 0x7e, 0xa9, 0x37, 0xbe, 0xe7, 0xd9, 0x8d, 0x85, 0x31, 0x8f, 0xbc, 0x47, 0xde,
	0xc0, 0x41, 0xca, 0x68, 0x66, 0xa9, 0xc7, 0xdd, 0x42, 0x7c, 0xfe, 0x6c, 0xd2, 0x0d, 0x2b, 0xea,
	0x3b, 0x88, 0xb5, 0x81, 0x96, 0x3c, 0xeb, 0x5f, 0x4b, 0x35, 0x87, 0x1d, 0x0a, 0x5e, 0x03, 0x68,
	0xf9, 0xb8, 0xb1, 0x93, 0xd6, 0x80, 0x90, 0xda, 0x9e, 0x99, 0xe3, 0x9d, 0x41, 0xa4, 0x94, 0xeb,
	0x3d, 0x6f, 0x95, 0x74, 0x35, 0xdb, 0x9b, 0x24, 0xe9, 0x02, 0x40, 0x8b, 0x42, 0xda, 0x49, 0xcf,
	0x36, 0xec, 0x68, 0x79, 0x15, 0xe0, 0x4f, 0xf5, 0xec, 0x4f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x55,
	0x40, 0x17, 0x61, 0x65, 0x05, 0x00, 0x00,
}
