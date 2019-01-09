// Code generated by protoc-gen-go. DO NOT EDIT.
// source: crypto.proto

/*
Package protos is a generated protocol buffer package.

It is generated from these files:
	crypto.proto
	tx.proto
	txaddr.proto

It has these top-level messages:
	KeyDerived
	PrivateKey
	PublicKey
	Signature
	ECPoint
	TxBase
	TxHeader
	TxCredential
	TxBatch
	TxBatchResp
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

type KeyDerived struct {
	RootFingerprint []byte `protobuf:"bytes,1,opt,name=rootFingerprint,proto3" json:"rootFingerprint,omitempty"`
	Index           []byte `protobuf:"bytes,2,opt,name=index,proto3" json:"index,omitempty"`
	Chaincode       []byte `protobuf:"bytes,3,opt,name=chaincode,proto3" json:"chaincode,omitempty"`
}

func (m *KeyDerived) Reset()                    { *m = KeyDerived{} }
func (m *KeyDerived) String() string            { return proto.CompactTextString(m) }
func (*KeyDerived) ProtoMessage()               {}
func (*KeyDerived) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *KeyDerived) GetRootFingerprint() []byte {
	if m != nil {
		return m.RootFingerprint
	}
	return nil
}

func (m *KeyDerived) GetIndex() []byte {
	if m != nil {
		return m.Index
	}
	return nil
}

func (m *KeyDerived) GetChaincode() []byte {
	if m != nil {
		return m.Chaincode
	}
	return nil
}

type PrivateKey struct {
	Version int32       `protobuf:"varint,1,opt,name=version" json:"version,omitempty"`
	Kd      *KeyDerived `protobuf:"bytes,7,opt,name=kd" json:"kd,omitempty"`
	// Types that are valid to be assigned to Priv:
	//	*PrivateKey_Ec
	Priv isPrivateKey_Priv `protobuf_oneof:"priv"`
}

func (m *PrivateKey) Reset()                    { *m = PrivateKey{} }
func (m *PrivateKey) String() string            { return proto.CompactTextString(m) }
func (*PrivateKey) ProtoMessage()               {}
func (*PrivateKey) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type isPrivateKey_Priv interface {
	isPrivateKey_Priv()
}

type PrivateKey_Ec struct {
	Ec *PrivateKey_ECDSA `protobuf:"bytes,8,opt,name=ec,oneof"`
}

func (*PrivateKey_Ec) isPrivateKey_Priv() {}

func (m *PrivateKey) GetPriv() isPrivateKey_Priv {
	if m != nil {
		return m.Priv
	}
	return nil
}

func (m *PrivateKey) GetVersion() int32 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *PrivateKey) GetKd() *KeyDerived {
	if m != nil {
		return m.Kd
	}
	return nil
}

func (m *PrivateKey) GetEc() *PrivateKey_ECDSA {
	if x, ok := m.GetPriv().(*PrivateKey_Ec); ok {
		return x.Ec
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*PrivateKey) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _PrivateKey_OneofMarshaler, _PrivateKey_OneofUnmarshaler, _PrivateKey_OneofSizer, []interface{}{
		(*PrivateKey_Ec)(nil),
	}
}

func _PrivateKey_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*PrivateKey)
	// priv
	switch x := m.Priv.(type) {
	case *PrivateKey_Ec:
		b.EncodeVarint(8<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Ec); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("PrivateKey.Priv has unexpected type %T", x)
	}
	return nil
}

func _PrivateKey_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*PrivateKey)
	switch tag {
	case 8: // priv.ec
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(PrivateKey_ECDSA)
		err := b.DecodeMessage(msg)
		m.Priv = &PrivateKey_Ec{msg}
		return true, err
	default:
		return false, nil
	}
}

func _PrivateKey_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*PrivateKey)
	// priv
	switch x := m.Priv.(type) {
	case *PrivateKey_Ec:
		s := proto.Size(x.Ec)
		n += proto.SizeVarint(8<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type PrivateKey_ECDSA struct {
	Curvetype int32  `protobuf:"varint,1,opt,name=curvetype" json:"curvetype,omitempty"`
	D         []byte `protobuf:"bytes,2,opt,name=d,proto3" json:"d,omitempty"`
}

func (m *PrivateKey_ECDSA) Reset()                    { *m = PrivateKey_ECDSA{} }
func (m *PrivateKey_ECDSA) String() string            { return proto.CompactTextString(m) }
func (*PrivateKey_ECDSA) ProtoMessage()               {}
func (*PrivateKey_ECDSA) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1, 0} }

func (m *PrivateKey_ECDSA) GetCurvetype() int32 {
	if m != nil {
		return m.Curvetype
	}
	return 0
}

func (m *PrivateKey_ECDSA) GetD() []byte {
	if m != nil {
		return m.D
	}
	return nil
}

type PublicKey struct {
	Version int32       `protobuf:"varint,1,opt,name=version" json:"version,omitempty"`
	Kd      *KeyDerived `protobuf:"bytes,7,opt,name=kd" json:"kd,omitempty"`
	// Types that are valid to be assigned to Pub:
	//	*PublicKey_Ec
	Pub isPublicKey_Pub `protobuf_oneof:"pub"`
}

func (m *PublicKey) Reset()                    { *m = PublicKey{} }
func (m *PublicKey) String() string            { return proto.CompactTextString(m) }
func (*PublicKey) ProtoMessage()               {}
func (*PublicKey) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

type isPublicKey_Pub interface {
	isPublicKey_Pub()
}

type PublicKey_Ec struct {
	Ec *PublicKey_ECDSA `protobuf:"bytes,8,opt,name=ec,oneof"`
}

func (*PublicKey_Ec) isPublicKey_Pub() {}

func (m *PublicKey) GetPub() isPublicKey_Pub {
	if m != nil {
		return m.Pub
	}
	return nil
}

func (m *PublicKey) GetVersion() int32 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *PublicKey) GetKd() *KeyDerived {
	if m != nil {
		return m.Kd
	}
	return nil
}

func (m *PublicKey) GetEc() *PublicKey_ECDSA {
	if x, ok := m.GetPub().(*PublicKey_Ec); ok {
		return x.Ec
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*PublicKey) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _PublicKey_OneofMarshaler, _PublicKey_OneofUnmarshaler, _PublicKey_OneofSizer, []interface{}{
		(*PublicKey_Ec)(nil),
	}
}

func _PublicKey_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*PublicKey)
	// pub
	switch x := m.Pub.(type) {
	case *PublicKey_Ec:
		b.EncodeVarint(8<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Ec); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("PublicKey.Pub has unexpected type %T", x)
	}
	return nil
}

func _PublicKey_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*PublicKey)
	switch tag {
	case 8: // pub.ec
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(PublicKey_ECDSA)
		err := b.DecodeMessage(msg)
		m.Pub = &PublicKey_Ec{msg}
		return true, err
	default:
		return false, nil
	}
}

func _PublicKey_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*PublicKey)
	// pub
	switch x := m.Pub.(type) {
	case *PublicKey_Ec:
		s := proto.Size(x.Ec)
		n += proto.SizeVarint(8<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type PublicKey_ECDSA struct {
	Curvetype int32    `protobuf:"varint,1,opt,name=curvetype" json:"curvetype,omitempty"`
	P         *ECPoint `protobuf:"bytes,2,opt,name=p" json:"p,omitempty"`
}

func (m *PublicKey_ECDSA) Reset()                    { *m = PublicKey_ECDSA{} }
func (m *PublicKey_ECDSA) String() string            { return proto.CompactTextString(m) }
func (*PublicKey_ECDSA) ProtoMessage()               {}
func (*PublicKey_ECDSA) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2, 0} }

func (m *PublicKey_ECDSA) GetCurvetype() int32 {
	if m != nil {
		return m.Curvetype
	}
	return 0
}

func (m *PublicKey_ECDSA) GetP() *ECPoint {
	if m != nil {
		return m.P
	}
	return nil
}

type Signature struct {
	// Types that are valid to be assigned to Data:
	//	*Signature_Ec
	Data isSignature_Data `protobuf_oneof:"data"`
}

func (m *Signature) Reset()                    { *m = Signature{} }
func (m *Signature) String() string            { return proto.CompactTextString(m) }
func (*Signature) ProtoMessage()               {}
func (*Signature) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

type isSignature_Data interface {
	isSignature_Data()
}

type Signature_Ec struct {
	Ec *Signature_ECDSA `protobuf:"bytes,4,opt,name=ec,oneof"`
}

func (*Signature_Ec) isSignature_Data() {}

func (m *Signature) GetData() isSignature_Data {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *Signature) GetEc() *Signature_ECDSA {
	if x, ok := m.GetData().(*Signature_Ec); ok {
		return x.Ec
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*Signature) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _Signature_OneofMarshaler, _Signature_OneofUnmarshaler, _Signature_OneofSizer, []interface{}{
		(*Signature_Ec)(nil),
	}
}

func _Signature_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*Signature)
	// data
	switch x := m.Data.(type) {
	case *Signature_Ec:
		b.EncodeVarint(4<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Ec); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("Signature.Data has unexpected type %T", x)
	}
	return nil
}

func _Signature_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*Signature)
	switch tag {
	case 4: // data.ec
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Signature_ECDSA)
		err := b.DecodeMessage(msg)
		m.Data = &Signature_Ec{msg}
		return true, err
	default:
		return false, nil
	}
}

func _Signature_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*Signature)
	// data
	switch x := m.Data.(type) {
	case *Signature_Ec:
		s := proto.Size(x.Ec)
		n += proto.SizeVarint(4<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type Signature_ECDSA struct {
	R []byte `protobuf:"bytes,1,opt,name=r,proto3" json:"r,omitempty"`
	S []byte `protobuf:"bytes,2,opt,name=s,proto3" json:"s,omitempty"`
	V int32  `protobuf:"varint,3,opt,name=v" json:"v,omitempty"`
}

func (m *Signature_ECDSA) Reset()                    { *m = Signature_ECDSA{} }
func (m *Signature_ECDSA) String() string            { return proto.CompactTextString(m) }
func (*Signature_ECDSA) ProtoMessage()               {}
func (*Signature_ECDSA) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3, 0} }

func (m *Signature_ECDSA) GetR() []byte {
	if m != nil {
		return m.R
	}
	return nil
}

func (m *Signature_ECDSA) GetS() []byte {
	if m != nil {
		return m.S
	}
	return nil
}

func (m *Signature_ECDSA) GetV() int32 {
	if m != nil {
		return m.V
	}
	return 0
}

type ECPoint struct {
	X []byte `protobuf:"bytes,1,opt,name=x,proto3" json:"x,omitempty"`
	Y []byte `protobuf:"bytes,2,opt,name=y,proto3" json:"y,omitempty"`
}

func (m *ECPoint) Reset()                    { *m = ECPoint{} }
func (m *ECPoint) String() string            { return proto.CompactTextString(m) }
func (*ECPoint) ProtoMessage()               {}
func (*ECPoint) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

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

func init() {
	proto.RegisterType((*KeyDerived)(nil), "protos.KeyDerived")
	proto.RegisterType((*PrivateKey)(nil), "protos.PrivateKey")
	proto.RegisterType((*PrivateKey_ECDSA)(nil), "protos.PrivateKey.ECDSA")
	proto.RegisterType((*PublicKey)(nil), "protos.PublicKey")
	proto.RegisterType((*PublicKey_ECDSA)(nil), "protos.PublicKey.ECDSA")
	proto.RegisterType((*Signature)(nil), "protos.Signature")
	proto.RegisterType((*Signature_ECDSA)(nil), "protos.Signature.ECDSA")
	proto.RegisterType((*ECPoint)(nil), "protos.ECPoint")
}

func init() { proto.RegisterFile("crypto.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 350 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x92, 0xdd, 0x4a, 0xf3, 0x40,
	0x10, 0x86, 0xbf, 0x4d, 0x9b, 0xf6, 0xeb, 0x34, 0x50, 0x58, 0x04, 0x83, 0x28, 0x94, 0x80, 0x50,
	0x3d, 0x28, 0xd8, 0x5e, 0x81, 0xfd, 0x11, 0xa1, 0x27, 0x25, 0xbd, 0x82, 0x34, 0x3b, 0xd4, 0xb5,
	0x92, 0x5d, 0x37, 0x9b, 0xa5, 0xb9, 0x2e, 0xaf, 0xc0, 0x3b, 0x93, 0x6c, 0xb3, 0x8d, 0xf5, 0x48,
	0xf0, 0x28, 0xbc, 0xf3, 0xf3, 0xce, 0x93, 0xd9, 0x81, 0x20, 0x55, 0xa5, 0xd4, 0x62, 0x2c, 0x95,
	0xd0, 0x82, 0x76, 0xec, 0x27, 0x8f, 0x5e, 0x01, 0x56, 0x58, 0x2e, 0x50, 0x71, 0x83, 0x8c, 0x8e,
	0x60, 0xa0, 0x84, 0xd0, 0x4f, 0x3c, 0xdb, 0xa1, 0x92, 0x8a, 0x67, 0x3a, 0x24, 0x43, 0x32, 0x0a,
	0xe2, 0x9f, 0x61, 0x7a, 0x01, 0x3e, 0xcf, 0x18, 0x1e, 0x42, 0xcf, 0xe6, 0x8f, 0x82, 0x5e, 0x43,
	0x2f, 0x7d, 0x49, 0x78, 0x96, 0x0a, 0x86, 0x61, 0xcb, 0x66, 0x9a, 0x40, 0xf4, 0x41, 0x00, 0xd6,
	0x8a, 0x9b, 0x44, 0xe3, 0x0a, 0x4b, 0x1a, 0x42, 0xd7, 0xa0, 0xca, 0xb9, 0xc8, 0xec, 0x10, 0x3f,
	0x76, 0x92, 0x46, 0xe0, 0xed, 0x59, 0xd8, 0x1d, 0x92, 0x51, 0x7f, 0x42, 0x8f, 0xc0, 0xf9, 0xb8,
	0xc1, 0x8c, 0xbd, 0x3d, 0xa3, 0xf7, 0xe0, 0x61, 0x1a, 0xfe, 0xb7, 0x35, 0xa1, 0xab, 0x69, 0xdc,
	0xc7, 0xcb, 0xf9, 0x62, 0xf3, 0xf8, 0xfc, 0x2f, 0xf6, 0x30, 0xbd, 0x9a, 0x82, 0x6f, 0xa5, 0xe5,
	0x2b, 0x94, 0x41, 0x5d, 0x4a, 0xac, 0x87, 0x36, 0x01, 0x1a, 0x00, 0x61, 0xf5, 0xff, 0x10, 0x36,
	0xeb, 0x40, 0x5b, 0x2a, 0x6e, 0xa2, 0x4f, 0x02, 0xbd, 0x75, 0xb1, 0x7d, 0xe3, 0xe9, 0xdf, 0xa1,
	0xef, 0xbe, 0x41, 0x5f, 0x9e, 0xa0, 0x9d, 0xf9, 0x19, 0xf3, 0xe2, 0x77, 0xcc, 0x37, 0x40, 0xa4,
	0x65, 0xee, 0x4f, 0x06, 0xce, 0x70, 0x39, 0x5f, 0x0b, 0x9e, 0xe9, 0x98, 0xc8, 0x99, 0x0f, 0x2d,
	0x59, 0x6c, 0xa3, 0x77, 0xe8, 0x6d, 0xf8, 0x2e, 0x4b, 0x74, 0xa1, 0xb0, 0x86, 0x68, 0x9f, 0x43,
	0x9c, 0xd2, 0x67, 0x10, 0x0f, 0x0e, 0x22, 0x00, 0xa2, 0xea, 0x53, 0x20, 0xaa, 0x52, 0xb9, 0x5b,
	0x54, 0x5e, 0x29, 0x63, 0x1f, 0xdb, 0x8f, 0x89, 0xa9, 0xd6, 0xc6, 0x12, 0x9d, 0x44, 0xb7, 0xd0,
	0xad, 0x39, 0xaa, 0x82, 0x83, 0x6b, 0x3e, 0x54, 0xaa, 0x74, 0xcd, 0xe5, 0xf6, 0x78, 0x87, 0xd3,
	0xaf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x46, 0x20, 0x91, 0x88, 0x9e, 0x02, 0x00, 0x00,
}
