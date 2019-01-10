// Code generated by protoc-gen-go. DO NOT EDIT.
// source: regtx.proto

/*
Package ccprotos is a generated protocol buffer package.

It is generated from these files:
	regtx.proto

It has these top-level messages:
	RegPublicKey
	RevokePublicKey
	ActivePublicKey
	Settings
	RegGlobalData
	RegData
*/
package ccprotos

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import protos "hyperledger.abchain.org/protos"
import google_protobuf "github.com/golang/protobuf/ptypes/timestamp"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// user can register a public key only if it has own some pais
type RegPublicKey struct {
	Region  string `protobuf:"bytes,2,opt,name=region" json:"region,omitempty"`
	PkBytes []byte `protobuf:"bytes,3,opt,name=pkBytes,proto3" json:"pkBytes,omitempty"`
}

func (m *RegPublicKey) Reset()                    { *m = RegPublicKey{} }
func (m *RegPublicKey) String() string            { return proto.CompactTextString(m) }
func (*RegPublicKey) ProtoMessage()               {}
func (*RegPublicKey) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *RegPublicKey) GetRegion() string {
	if m != nil {
		return m.Region
	}
	return ""
}

func (m *RegPublicKey) GetPkBytes() []byte {
	if m != nil {
		return m.PkBytes
	}
	return nil
}

type RevokePublicKey struct {
	Pk *protos.PublicKey `protobuf:"bytes,1,opt,name=pk" json:"pk,omitempty"`
}

func (m *RevokePublicKey) Reset()                    { *m = RevokePublicKey{} }
func (m *RevokePublicKey) String() string            { return proto.CompactTextString(m) }
func (*RevokePublicKey) ProtoMessage()               {}
func (*RevokePublicKey) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *RevokePublicKey) GetPk() *protos.PublicKey {
	if m != nil {
		return m.Pk
	}
	return nil
}

type ActivePublicKey struct {
	Key []byte `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
}

func (m *ActivePublicKey) Reset()                    { *m = ActivePublicKey{} }
func (m *ActivePublicKey) String() string            { return proto.CompactTextString(m) }
func (*ActivePublicKey) ProtoMessage()               {}
func (*ActivePublicKey) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *ActivePublicKey) GetKey() []byte {
	if m != nil {
		return m.Key
	}
	return nil
}

type Settings struct {
	DebugMode      bool   `protobuf:"varint,1,opt,name=debugMode" json:"debugMode,omitempty"`
	RegPrivilege   string `protobuf:"bytes,2,opt,name=regPrivilege" json:"regPrivilege,omitempty"`
	AdminPrivilege string `protobuf:"bytes,3,opt,name=adminPrivilege" json:"adminPrivilege,omitempty"`
}

func (m *Settings) Reset()                    { *m = Settings{} }
func (m *Settings) String() string            { return proto.CompactTextString(m) }
func (*Settings) ProtoMessage()               {}
func (*Settings) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *Settings) GetDebugMode() bool {
	if m != nil {
		return m.DebugMode
	}
	return false
}

func (m *Settings) GetRegPrivilege() string {
	if m != nil {
		return m.RegPrivilege
	}
	return ""
}

func (m *Settings) GetAdminPrivilege() string {
	if m != nil {
		return m.AdminPrivilege
	}
	return ""
}

type RegGlobalData struct {
	RegPrivilege   string           `protobuf:"bytes,1,opt,name=regPrivilege" json:"regPrivilege,omitempty"`
	AdminPrivilege string           `protobuf:"bytes,2,opt,name=adminPrivilege" json:"adminPrivilege,omitempty"`
	Chaincodes     map[int32]string `protobuf:"bytes,5,rep,name=chaincodes" json:"chaincodes,omitempty" protobuf_key:"varint,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	DeployFlag     []byte           `protobuf:"bytes,10,opt,name=deployFlag,proto3" json:"deployFlag,omitempty"`
}

func (m *RegGlobalData) Reset()                    { *m = RegGlobalData{} }
func (m *RegGlobalData) String() string            { return proto.CompactTextString(m) }
func (*RegGlobalData) ProtoMessage()               {}
func (*RegGlobalData) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *RegGlobalData) GetRegPrivilege() string {
	if m != nil {
		return m.RegPrivilege
	}
	return ""
}

func (m *RegGlobalData) GetAdminPrivilege() string {
	if m != nil {
		return m.AdminPrivilege
	}
	return ""
}

func (m *RegGlobalData) GetChaincodes() map[int32]string {
	if m != nil {
		return m.Chaincodes
	}
	return nil
}

func (m *RegGlobalData) GetDeployFlag() []byte {
	if m != nil {
		return m.DeployFlag
	}
	return nil
}

type RegData struct {
	Pk        *protos.PublicKey          `protobuf:"bytes,1,opt,name=pk" json:"pk,omitempty"`
	RegTxid   string                     `protobuf:"bytes,3,opt,name=regTxid" json:"regTxid,omitempty"`
	RegTs     *google_protobuf.Timestamp `protobuf:"bytes,4,opt,name=regTs" json:"regTs,omitempty"`
	Region    string                     `protobuf:"bytes,5,opt,name=region" json:"region,omitempty"`
	Enabled   bool                       `protobuf:"varint,6,opt,name=enabled" json:"enabled,omitempty"`
	Authcodes []int32                    `protobuf:"varint,10,rep,packed,name=authcodes" json:"authcodes,omitempty"`
}

func (m *RegData) Reset()                    { *m = RegData{} }
func (m *RegData) String() string            { return proto.CompactTextString(m) }
func (*RegData) ProtoMessage()               {}
func (*RegData) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *RegData) GetPk() *protos.PublicKey {
	if m != nil {
		return m.Pk
	}
	return nil
}

func (m *RegData) GetRegTxid() string {
	if m != nil {
		return m.RegTxid
	}
	return ""
}

func (m *RegData) GetRegTs() *google_protobuf.Timestamp {
	if m != nil {
		return m.RegTs
	}
	return nil
}

func (m *RegData) GetRegion() string {
	if m != nil {
		return m.Region
	}
	return ""
}

func (m *RegData) GetEnabled() bool {
	if m != nil {
		return m.Enabled
	}
	return false
}

func (m *RegData) GetAuthcodes() []int32 {
	if m != nil {
		return m.Authcodes
	}
	return nil
}

func init() {
	proto.RegisterType((*RegPublicKey)(nil), "ccprotos.RegPublicKey")
	proto.RegisterType((*RevokePublicKey)(nil), "ccprotos.RevokePublicKey")
	proto.RegisterType((*ActivePublicKey)(nil), "ccprotos.ActivePublicKey")
	proto.RegisterType((*Settings)(nil), "ccprotos.Settings")
	proto.RegisterType((*RegGlobalData)(nil), "ccprotos.RegGlobalData")
	proto.RegisterType((*RegData)(nil), "ccprotos.RegData")
}

func init() { proto.RegisterFile("regtx.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 451 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x52, 0x5f, 0x8b, 0xd3, 0x4e,
	0x14, 0x25, 0xc9, 0x2f, 0xdd, 0xee, 0x6d, 0x7f, 0x56, 0x07, 0x91, 0x50, 0x44, 0x63, 0x04, 0x0d,
	0x08, 0x53, 0x59, 0x7d, 0x10, 0x41, 0xf0, 0xff, 0x3e, 0x88, 0xb0, 0x8c, 0xfb, 0x05, 0x26, 0xc9,
	0x75, 0x3a, 0x64, 0x9a, 0x09, 0x93, 0x49, 0xd9, 0x3c, 0xfb, 0xe5, 0xfc, 0x58, 0xd2, 0x49, 0xb3,
	0xcd, 0x56, 0x41, 0xdf, 0xee, 0x9f, 0x73, 0xee, 0x9d, 0x39, 0xf7, 0xc0, 0xcc, 0xa0, 0xb0, 0x57,
	0xb4, 0x36, 0xda, 0x6a, 0x32, 0xcd, 0x73, 0x17, 0x34, 0xcb, 0x67, 0xeb, 0xae, 0x46, 0xa3, 0xb0,
	0x10, 0x68, 0x28, 0xcf, 0xf2, 0x35, 0x97, 0x15, 0xd5, 0x46, 0xac, 0xfa, 0xfe, 0x2a, 0x37, 0x5d,
	0x6d, 0x75, 0x4f, 0x5b, 0x3e, 0x14, 0x5a, 0x0b, 0x85, 0x7d, 0x2f, 0x6b, 0xbf, 0xaf, 0xac, 0xdc,
	0x60, 0x63, 0xf9, 0xa6, 0xee, 0x01, 0xc9, 0x5b, 0x98, 0x33, 0x14, 0x17, 0x6d, 0xa6, 0x64, 0xfe,
	0x05, 0x3b, 0x72, 0x0f, 0x26, 0x06, 0x85, 0xd4, 0x55, 0xe4, 0xc7, 0x5e, 0x7a, 0xca, 0xf6, 0x19,
	0x89, 0xe0, 0xa4, 0x2e, 0xdf, 0x77, 0x16, 0x9b, 0x28, 0x88, 0xbd, 0x74, 0xce, 0x86, 0x34, 0x79,
	0x09, 0x0b, 0x86, 0x5b, 0x5d, 0xe2, 0x61, 0xc8, 0x23, 0xf0, 0xeb, 0x32, 0xf2, 0x62, 0x2f, 0x9d,
	0x9d, 0xdd, 0xe9, 0x17, 0x35, 0xf4, 0xba, 0xcd, 0xfc, 0xba, 0x4c, 0x1e, 0xc3, 0xe2, 0x5d, 0x6e,
	0xe5, 0x76, 0xc4, 0xba, 0x0d, 0x41, 0x89, 0x9d, 0xa3, 0xcd, 0xd9, 0x2e, 0x4c, 0x2c, 0x4c, 0xbf,
	0xa1, 0xb5, 0xb2, 0x12, 0x0d, 0xb9, 0x0f, 0xa7, 0x05, 0x66, 0xad, 0xf8, 0xaa, 0x0b, 0x74, 0x98,
	0x29, 0x3b, 0x14, 0x48, 0x02, 0x73, 0x83, 0xe2, 0xc2, 0xc8, 0xad, 0x54, 0x28, 0x70, 0xff, 0xf8,
	0x1b, 0x35, 0xf2, 0x04, 0x6e, 0xf1, 0x62, 0x23, 0xab, 0x03, 0x2a, 0x70, 0xa8, 0xa3, 0x6a, 0xf2,
	0xc3, 0x87, 0xff, 0x19, 0x8a, 0x73, 0xa5, 0x33, 0xae, 0x3e, 0x72, 0xcb, 0x7f, 0x9b, 0xee, 0xfd,
	0xd3, 0x74, 0xff, 0x4f, 0xd3, 0xc9, 0x39, 0x80, 0x3b, 0x59, 0xae, 0x0b, 0x6c, 0xa2, 0x30, 0x0e,
	0xd2, 0xd9, 0xd9, 0x53, 0x3a, 0x5c, 0x97, 0xde, 0x58, 0x4c, 0x3f, 0x5c, 0x23, 0x3f, 0x55, 0xd6,
	0x74, 0x6c, 0x44, 0x25, 0x0f, 0x00, 0x0a, 0xac, 0x95, 0xee, 0x3e, 0x2b, 0x2e, 0x22, 0x70, 0xaa,
	0x8d, 0x2a, 0xcb, 0x37, 0xb0, 0x38, 0xa2, 0x8f, 0x15, 0x0e, 0x9d, 0xc2, 0xe4, 0x2e, 0x84, 0x5b,
	0xae, 0xda, 0xe1, 0xb1, 0x7d, 0xf2, 0xda, 0x7f, 0xe5, 0x25, 0x3f, 0x3d, 0x38, 0x61, 0x28, 0xdc,
	0xff, 0xff, 0x7e, 0xcf, 0x9d, 0x3f, 0x0c, 0x8a, 0xcb, 0x2b, 0x59, 0xec, 0x55, 0x1d, 0x52, 0xf2,
	0x1c, 0xc2, 0x5d, 0xd8, 0x44, 0xff, 0x39, 0xfe, 0x92, 0xf6, 0x96, 0xa4, 0x83, 0x25, 0xe9, 0xe5,
	0x60, 0x49, 0xd6, 0x03, 0x47, 0x1e, 0x0c, 0x8f, 0x3d, 0x88, 0x15, 0xcf, 0x14, 0x16, 0xd1, 0xc4,
	0x19, 0x60, 0x48, 0x77, 0xe6, 0xe0, 0xad, 0x5d, 0xf7, 0x9a, 0x42, 0x1c, 0xa4, 0x21, 0x3b, 0x14,
	0xb2, 0x89, 0x5b, 0xf5, 0xe2, 0x57, 0x00, 0x00, 0x00, 0xff, 0xff, 0xd9, 0x27, 0xfc, 0x6f, 0x51,
	0x03, 0x00, 0x00,
}
