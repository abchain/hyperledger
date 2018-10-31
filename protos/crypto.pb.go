// Code generated by protoc-gen-go. DO NOT EDIT.
// source: crypto.proto

/*
Package protos is a generated protocol buffer package.

It is generated from these files:
	crypto.proto
	deploytx.proto
	tx.proto
	txaddr.proto

It has these top-level messages:
	PrivateKey
	PublicKey
	Signature
	ECPoint
	ECSignature
	DeployTx
	TxBase
	TxHeader
	TxCredential
	TxAddr
	TxMsgExample
*/
package protos

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type PrivateKey struct {
	Version         int32  `protobuf:"varint,1,opt,name=version" json:"version,omitempty"`
	Curvetype       int32  `protobuf:"varint,2,opt,name=curvetype" json:"curvetype,omitempty"`
	RootFingerprint []byte `protobuf:"bytes,3,opt,name=rootFingerprint,proto3" json:"rootFingerprint,omitempty"`
	Index           []byte `protobuf:"bytes,4,opt,name=index,proto3" json:"index,omitempty"`
	Chaincode       []byte `protobuf:"bytes,5,opt,name=chaincode,proto3" json:"chaincode,omitempty"`
	D               []byte `protobuf:"bytes,6,opt,name=d,proto3" json:"d,omitempty"`
}

func (m *PrivateKey) Reset()                    { *m = PrivateKey{} }
func (m *PrivateKey) String() string            { return proto.CompactTextString(m) }
func (*PrivateKey) ProtoMessage()               {}
func (*PrivateKey) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *PrivateKey) GetVersion() int32 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *PrivateKey) GetCurvetype() int32 {
	if m != nil {
		return m.Curvetype
	}
	return 0
}

func (m *PrivateKey) GetRootFingerprint() []byte {
	if m != nil {
		return m.RootFingerprint
	}
	return nil
}

func (m *PrivateKey) GetIndex() []byte {
	if m != nil {
		return m.Index
	}
	return nil
}

func (m *PrivateKey) GetChaincode() []byte {
	if m != nil {
		return m.Chaincode
	}
	return nil
}

func (m *PrivateKey) GetD() []byte {
	if m != nil {
		return m.D
	}
	return nil
}

type PublicKey struct {
	Version         int32    `protobuf:"varint,1,opt,name=version" json:"version,omitempty"`
	Curvetype       int32    `protobuf:"varint,2,opt,name=curvetype" json:"curvetype,omitempty"`
	RootFingerprint []byte   `protobuf:"bytes,3,opt,name=rootFingerprint,proto3" json:"rootFingerprint,omitempty"`
	Index           []byte   `protobuf:"bytes,4,opt,name=index,proto3" json:"index,omitempty"`
	Chaincode       []byte   `protobuf:"bytes,5,opt,name=chaincode,proto3" json:"chaincode,omitempty"`
	P               *ECPoint `protobuf:"bytes,6,opt,name=p" json:"p,omitempty"`
}

func (m *PublicKey) Reset()                    { *m = PublicKey{} }
func (m *PublicKey) String() string            { return proto.CompactTextString(m) }
func (*PublicKey) ProtoMessage()               {}
func (*PublicKey) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *PublicKey) GetVersion() int32 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *PublicKey) GetCurvetype() int32 {
	if m != nil {
		return m.Curvetype
	}
	return 0
}

func (m *PublicKey) GetRootFingerprint() []byte {
	if m != nil {
		return m.RootFingerprint
	}
	return nil
}

func (m *PublicKey) GetIndex() []byte {
	if m != nil {
		return m.Index
	}
	return nil
}

func (m *PublicKey) GetChaincode() []byte {
	if m != nil {
		return m.Chaincode
	}
	return nil
}

func (m *PublicKey) GetP() *ECPoint {
	if m != nil {
		return m.P
	}
	return nil
}

type Signature struct {
	P *ECPoint `protobuf:"bytes,1,opt,name=p" json:"p,omitempty"`
}

func (m *Signature) Reset()                    { *m = Signature{} }
func (m *Signature) String() string            { return proto.CompactTextString(m) }
func (*Signature) ProtoMessage()               {}
func (*Signature) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *Signature) GetP() *ECPoint {
	if m != nil {
		return m.P
	}
	return nil
}

type ECPoint struct {
	X []byte `protobuf:"bytes,1,opt,name=x,proto3" json:"x,omitempty"`
	Y []byte `protobuf:"bytes,2,opt,name=y,proto3" json:"y,omitempty"`
}

func (m *ECPoint) Reset()                    { *m = ECPoint{} }
func (m *ECPoint) String() string            { return proto.CompactTextString(m) }
func (*ECPoint) ProtoMessage()               {}
func (*ECPoint) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *ECPoint) GetX() []byte {
	if m != nil {
		return m.X
	}
	return nil
}

func (m *ECPoint) GetY() []byte {
	if m != nil {
		return m.Y
	}
	return nil
}

type ECSignature struct {
	R []byte `protobuf:"bytes,1,opt,name=r,proto3" json:"r,omitempty"`
	S []byte `protobuf:"bytes,2,opt,name=s,proto3" json:"s,omitempty"`
	V int32  `protobuf:"varint,3,opt,name=v" json:"v,omitempty"`
}

func (m *ECSignature) Reset()                    { *m = ECSignature{} }
func (m *ECSignature) String() string            { return proto.CompactTextString(m) }
func (*ECSignature) ProtoMessage()               {}
func (*ECSignature) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *ECSignature) GetR() []byte {
	if m != nil {
		return m.R
	}
	return nil
}

func (m *ECSignature) GetS() []byte {
	if m != nil {
		return m.S
	}
	return nil
}

func (m *ECSignature) GetV() int32 {
	if m != nil {
		return m.V
	}
	return 0
}

func init() {
	proto.RegisterType((*PrivateKey)(nil), "protos.PrivateKey")
	proto.RegisterType((*PublicKey)(nil), "protos.PublicKey")
	proto.RegisterType((*Signature)(nil), "protos.Signature")
	proto.RegisterType((*ECPoint)(nil), "protos.ECPoint")
	proto.RegisterType((*ECSignature)(nil), "protos.ECSignature")
}

func init() { proto.RegisterFile("crypto.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 266 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xcc, 0x92, 0xc1, 0x4a, 0xc4, 0x30,
	0x10, 0x86, 0x19, 0xb5, 0xbb, 0x74, 0x36, 0xb0, 0x10, 0x3c, 0xe4, 0xa0, 0xb0, 0x14, 0x84, 0xe2,
	0x61, 0x0f, 0x7a, 0xf0, 0x01, 0x96, 0xf5, 0xe2, 0xa5, 0xc4, 0x27, 0xe8, 0xb6, 0x61, 0x0d, 0x48,
	0x12, 0xa6, 0x69, 0x69, 0xdf, 0xca, 0x47, 0xf0, 0xd1, 0xa4, 0xd3, 0x2d, 0x05, 0xc1, 0xbb, 0xa7,
	0xf0, 0xfd, 0xf3, 0x31, 0xcc, 0x0c, 0x41, 0x51, 0xd1, 0x10, 0xa2, 0xdf, 0x07, 0xf2, 0xd1, 0xcb,
	0x15, 0x3f, 0x4d, 0xf6, 0x05, 0x88, 0x05, 0xd9, 0xae, 0x8c, 0xe6, 0xcd, 0x0c, 0x52, 0xe1, 0xba,
	0x33, 0xd4, 0x58, 0xef, 0x14, 0xec, 0x20, 0x4f, 0xf4, 0x8c, 0xf2, 0x0e, 0xd3, 0xaa, 0xa5, 0xce,
	0xc4, 0x21, 0x18, 0x75, 0xc5, 0xb5, 0x25, 0x90, 0x39, 0x6e, 0xc9, 0xfb, 0xf8, 0x6a, 0xdd, 0xd9,
	0x50, 0x20, 0xeb, 0xa2, 0xba, 0xde, 0x41, 0x2e, 0xf4, 0xef, 0x58, 0xde, 0x62, 0x62, 0x5d, 0x6d,
	0x7a, 0x75, 0xc3, 0xf5, 0x09, 0xb8, 0xfb, 0x47, 0x69, 0x5d, 0xe5, 0x6b, 0xa3, 0x12, 0xae, 0x2c,
	0x81, 0x14, 0x08, 0xb5, 0x5a, 0x71, 0x0a, 0x75, 0xf6, 0x0d, 0x98, 0x16, 0xed, 0xe9, 0xd3, 0x56,
	0xff, 0x77, 0xe2, 0x7b, 0x84, 0xc0, 0x13, 0x6f, 0x9e, 0xb6, 0xd3, 0xc5, 0x9b, 0xfd, 0xf1, 0x50,
	0x78, 0xeb, 0xa2, 0x86, 0x90, 0x3d, 0x62, 0xfa, 0x6e, 0xcf, 0xae, 0x8c, 0x2d, 0x5d, 0x5c, 0xf8,
	0xd3, 0x7d, 0xc0, 0xf5, 0x85, 0xc6, 0x3b, 0xf4, 0x6c, 0x0a, 0x0d, 0xfd, 0x48, 0x03, 0xef, 0x25,
	0x34, 0x0c, 0xd9, 0x0b, 0x6e, 0x8e, 0x87, 0xa5, 0xa9, 0x40, 0xa0, 0x59, 0xa5, 0x91, 0x9a, 0x59,
	0x6d, 0x46, 0xea, 0x78, 0xd9, 0x44, 0x43, 0x77, 0x9a, 0x7e, 0xc2, 0xf3, 0x4f, 0x00, 0x00, 0x00,
	0xff, 0xff, 0xb3, 0xc9, 0x97, 0x72, 0x20, 0x02, 0x00, 0x00,
}
