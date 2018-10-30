// Code generated by protoc-gen-go. DO NOT EDIT.
// source: persistent.proto

/*
Package ccprotos is a generated protocol buffer package.

It is generated from these files:
	persistent.proto
	subscriptiontx.proto

It has these top-level messages:
	Contract
	RegContract
	QueryContract
	RedeemContract
*/
package ccprotos

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import protos "hyperledger.abchain.org/protos"
import protos1 "hyperledger.abchain.org/protos"
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

type Contract struct {
	DelegatorPk *protos.PublicKey          `protobuf:"bytes,1,opt,name=delegatorPk" json:"delegatorPk,omitempty"`
	TotalRedeem []byte                     `protobuf:"bytes,3,opt,name=totalRedeem,proto3" json:"totalRedeem,omitempty"`
	Status      []*Contract_MemberStatus   `protobuf:"bytes,5,rep,name=status" json:"status,omitempty"`
	ContractTs  *google_protobuf.Timestamp `protobuf:"bytes,6,opt,name=contractTs" json:"contractTs,omitempty"`
	FrozenTo    *google_protobuf.Timestamp `protobuf:"bytes,8,opt,name=frozenTo" json:"frozenTo,omitempty"`
	IsFrozen    bool                       `protobuf:"varint,10,opt,name=isFrozen" json:"isFrozen,omitempty"`
	NextAddr    *protos1.TxAddr            `protobuf:"bytes,12,opt,name=nextAddr" json:"nextAddr,omitempty"`
}

func (m *Contract) Reset()                    { *m = Contract{} }
func (m *Contract) String() string            { return proto.CompactTextString(m) }
func (*Contract) ProtoMessage()               {}
func (*Contract) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Contract) GetDelegatorPk() *protos.PublicKey {
	if m != nil {
		return m.DelegatorPk
	}
	return nil
}

func (m *Contract) GetTotalRedeem() []byte {
	if m != nil {
		return m.TotalRedeem
	}
	return nil
}

func (m *Contract) GetStatus() []*Contract_MemberStatus {
	if m != nil {
		return m.Status
	}
	return nil
}

func (m *Contract) GetContractTs() *google_protobuf.Timestamp {
	if m != nil {
		return m.ContractTs
	}
	return nil
}

func (m *Contract) GetFrozenTo() *google_protobuf.Timestamp {
	if m != nil {
		return m.FrozenTo
	}
	return nil
}

func (m *Contract) GetIsFrozen() bool {
	if m != nil {
		return m.IsFrozen
	}
	return false
}

func (m *Contract) GetNextAddr() *protos1.TxAddr {
	if m != nil {
		return m.NextAddr
	}
	return nil
}

type Contract_MemberStatus struct {
	Weight      uint32 `protobuf:"varint,2,opt,name=weight" json:"weight,omitempty"`
	TotalRedeem []byte `protobuf:"bytes,3,opt,name=totalRedeem,proto3" json:"totalRedeem,omitempty"`
	MemberID    string `protobuf:"bytes,5,opt,name=memberID" json:"memberID,omitempty"`
}

func (m *Contract_MemberStatus) Reset()                    { *m = Contract_MemberStatus{} }
func (m *Contract_MemberStatus) String() string            { return proto.CompactTextString(m) }
func (*Contract_MemberStatus) ProtoMessage()               {}
func (*Contract_MemberStatus) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 0} }

func (m *Contract_MemberStatus) GetWeight() uint32 {
	if m != nil {
		return m.Weight
	}
	return 0
}

func (m *Contract_MemberStatus) GetTotalRedeem() []byte {
	if m != nil {
		return m.TotalRedeem
	}
	return nil
}

func (m *Contract_MemberStatus) GetMemberID() string {
	if m != nil {
		return m.MemberID
	}
	return ""
}

func init() {
	proto.RegisterType((*Contract)(nil), "ccprotos.Contract")
	proto.RegisterType((*Contract_MemberStatus)(nil), "ccprotos.Contract.MemberStatus")
}

func init() { proto.RegisterFile("persistent.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 342 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x90, 0xcf, 0x6a, 0xe3, 0x30,
	0x10, 0xc6, 0xf1, 0x86, 0x18, 0xad, 0x9c, 0x5d, 0x76, 0x75, 0x58, 0x84, 0x2f, 0x31, 0x7b, 0x32,
	0x2d, 0x28, 0x90, 0x40, 0x0b, 0xbd, 0x95, 0x96, 0x42, 0x29, 0x85, 0xa0, 0xfa, 0x05, 0x64, 0x6b,
	0xe2, 0x98, 0xda, 0x96, 0x91, 0x26, 0x34, 0xe9, 0xb3, 0xf5, 0xe1, 0x4a, 0xfd, 0x0f, 0xdf, 0xd2,
	0x9b, 0xbe, 0xf9, 0x7e, 0x33, 0xa3, 0x6f, 0xe8, 0x9f, 0x06, 0xac, 0x2b, 0x1c, 0x42, 0x8d, 0xa2,
	0xb1, 0x06, 0x0d, 0x23, 0x59, 0xd6, 0x3e, 0x5c, 0x78, 0xb9, 0x3f, 0x35, 0x60, 0x4b, 0xd0, 0x39,
	0x58, 0xa1, 0xd2, 0x6c, 0xaf, 0x8a, 0x5a, 0x18, 0x9b, 0xaf, 0x3a, 0x7f, 0x95, 0xd9, 0x53, 0x83,
	0xa6, 0x6b, 0x3b, 0x0b, 0xe3, 0x51, 0x69, 0x6d, 0x7b, 0x78, 0x99, 0x1b, 0x93, 0x97, 0xd0, 0x79,
	0xe9, 0x61, 0xb7, 0xc2, 0xa2, 0x02, 0x87, 0xaa, 0x6a, 0x3a, 0xe0, 0xff, 0xc7, 0x8c, 0x92, 0x3b,
	0x53, 0xa3, 0x55, 0x19, 0xb2, 0x0d, 0x0d, 0x34, 0x94, 0x90, 0x2b, 0x34, 0x76, 0xfb, 0xca, 0xbd,
	0xc8, 0x8b, 0x83, 0xf5, 0xdf, 0x8e, 0x74, 0x62, 0x7b, 0x48, 0xcb, 0x22, 0x7b, 0x82, 0x93, 0x9c,
	0x52, 0x2c, 0xa2, 0x01, 0x1a, 0x54, 0xa5, 0x04, 0x0d, 0x50, 0xf1, 0x59, 0xe4, 0xc5, 0x0b, 0x39,
	0x2d, 0xb1, 0x6b, 0xea, 0x3b, 0x54, 0x78, 0x70, 0x7c, 0x1e, 0xcd, 0xe2, 0x60, 0xbd, 0x14, 0x43,
	0x72, 0x31, 0xac, 0x16, 0xcf, 0x50, 0xa5, 0x60, 0x5f, 0x5a, 0x4c, 0xf6, 0x38, 0xbb, 0xa1, 0x34,
	0xeb, 0x81, 0xc4, 0x71, 0xbf, 0xfd, 0x4e, 0x28, 0xba, 0x48, 0x62, 0x88, 0x24, 0x92, 0x21, 0x92,
	0x9c, 0xd0, 0xec, 0x8a, 0x92, 0x9d, 0x35, 0xef, 0x50, 0x27, 0x86, 0x93, 0xb3, 0x9d, 0x23, 0xcb,
	0x42, 0x4a, 0x0a, 0xf7, 0xd0, 0x2a, 0x4e, 0x23, 0x2f, 0x26, 0x72, 0xd4, 0xec, 0x82, 0x92, 0x1a,
	0x8e, 0x78, 0xab, 0xb5, 0xe5, 0x8b, 0x76, 0xe6, 0xef, 0xe1, 0x38, 0xc9, 0xf1, 0xab, 0x2a, 0x47,
	0x3f, 0xd4, 0x74, 0x31, 0xcd, 0xc4, 0xfe, 0x51, 0xff, 0x0d, 0x8a, 0x7c, 0x8f, 0xfc, 0x47, 0xe4,
	0xc5, 0xbf, 0x64, 0xaf, 0xbe, 0x71, 0xbe, 0x90, 0x92, 0xaa, 0x9d, 0xf4, 0x78, 0xcf, 0xe7, 0x91,
	0x17, 0xff, 0x94, 0xa3, 0x4e, 0xfd, 0x76, 0xfd, 0xe6, 0x33, 0x00, 0x00, 0xff, 0xff, 0xe7, 0x42,
	0x84, 0xd9, 0x5e, 0x02, 0x00, 0x00,
}