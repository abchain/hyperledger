// Code generated by protoc-gen-go. DO NOT EDIT.
// source: tx.proto

package protos

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/golang/protobuf/ptypes/timestamp"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type TxBase struct {
	Network string `protobuf:"bytes,1,opt,name=network" json:"network,omitempty"`
	Ccname  string `protobuf:"bytes,2,opt,name=ccname" json:"ccname,omitempty"`
	Method  string `protobuf:"bytes,3,opt,name=method" json:"method,omitempty"`
}

func (m *TxBase) Reset()                    { *m = TxBase{} }
func (m *TxBase) String() string            { return proto.CompactTextString(m) }
func (*TxBase) ProtoMessage()               {}
func (*TxBase) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{0} }

func (m *TxBase) GetNetwork() string {
	if m != nil {
		return m.Network
	}
	return ""
}

func (m *TxBase) GetCcname() string {
	if m != nil {
		return m.Ccname
	}
	return ""
}

func (m *TxBase) GetMethod() string {
	if m != nil {
		return m.Method
	}
	return ""
}

type TxHeader struct {
	Base      *TxBase                    `protobuf:"bytes,1,opt,name=base" json:"base,omitempty"`
	ExpiredTs *google_protobuf.Timestamp `protobuf:"bytes,2,opt,name=expiredTs" json:"expiredTs,omitempty"`
	Nonce     []byte                     `protobuf:"bytes,3,opt,name=nonce,proto3" json:"nonce,omitempty"`
}

func (m *TxHeader) Reset()                    { *m = TxHeader{} }
func (m *TxHeader) String() string            { return proto.CompactTextString(m) }
func (*TxHeader) ProtoMessage()               {}
func (*TxHeader) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{1} }

func (m *TxHeader) GetBase() *TxBase {
	if m != nil {
		return m.Base
	}
	return nil
}

func (m *TxHeader) GetExpiredTs() *google_protobuf.Timestamp {
	if m != nil {
		return m.ExpiredTs
	}
	return nil
}

func (m *TxHeader) GetNonce() []byte {
	if m != nil {
		return m.Nonce
	}
	return nil
}

type TxCredential struct {
	Addrc []*TxCredential_AddrCredentials `protobuf:"bytes,1,rep,name=addrc" json:"addrc,omitempty"`
}

func (m *TxCredential) Reset()                    { *m = TxCredential{} }
func (m *TxCredential) String() string            { return proto.CompactTextString(m) }
func (*TxCredential) ProtoMessage()               {}
func (*TxCredential) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{2} }

func (m *TxCredential) GetAddrc() []*TxCredential_AddrCredentials {
	if m != nil {
		return m.Addrc
	}
	return nil
}

type TxCredential_UserCredential struct {
	Signature *ECSignature `protobuf:"bytes,1,opt,name=signature" json:"signature,omitempty"`
	Pk        *PublicKey   `protobuf:"bytes,2,opt,name=pk" json:"pk,omitempty"`
}

func (m *TxCredential_UserCredential) Reset()                    { *m = TxCredential_UserCredential{} }
func (m *TxCredential_UserCredential) String() string            { return proto.CompactTextString(m) }
func (*TxCredential_UserCredential) ProtoMessage()               {}
func (*TxCredential_UserCredential) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{2, 0} }

func (m *TxCredential_UserCredential) GetSignature() *ECSignature {
	if m != nil {
		return m.Signature
	}
	return nil
}

func (m *TxCredential_UserCredential) GetPk() *PublicKey {
	if m != nil {
		return m.Pk
	}
	return nil
}

type TxCredential_InnerCredential struct {
	Ccname      string `protobuf:"bytes,1,opt,name=ccname" json:"ccname,omitempty"`
	Fingerprint []byte `protobuf:"bytes,2,opt,name=fingerprint,proto3" json:"fingerprint,omitempty"`
}

func (m *TxCredential_InnerCredential) Reset()                    { *m = TxCredential_InnerCredential{} }
func (m *TxCredential_InnerCredential) String() string            { return proto.CompactTextString(m) }
func (*TxCredential_InnerCredential) ProtoMessage()               {}
func (*TxCredential_InnerCredential) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{2, 1} }

func (m *TxCredential_InnerCredential) GetCcname() string {
	if m != nil {
		return m.Ccname
	}
	return ""
}

func (m *TxCredential_InnerCredential) GetFingerprint() []byte {
	if m != nil {
		return m.Fingerprint
	}
	return nil
}

type TxCredential_AddrCredentials struct {
	// Types that are valid to be assigned to Cred:
	//	*TxCredential_AddrCredentials_User
	//	*TxCredential_AddrCredentials_Cc
	Cred isTxCredential_AddrCredentials_Cred `protobuf_oneof:"cred"`
}

func (m *TxCredential_AddrCredentials) Reset()                    { *m = TxCredential_AddrCredentials{} }
func (m *TxCredential_AddrCredentials) String() string            { return proto.CompactTextString(m) }
func (*TxCredential_AddrCredentials) ProtoMessage()               {}
func (*TxCredential_AddrCredentials) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{2, 2} }

type isTxCredential_AddrCredentials_Cred interface {
	isTxCredential_AddrCredentials_Cred()
}

type TxCredential_AddrCredentials_User struct {
	User *TxCredential_UserCredential `protobuf:"bytes,1,opt,name=user,oneof"`
}
type TxCredential_AddrCredentials_Cc struct {
	Cc *TxCredential_InnerCredential `protobuf:"bytes,2,opt,name=cc,oneof"`
}

func (*TxCredential_AddrCredentials_User) isTxCredential_AddrCredentials_Cred() {}
func (*TxCredential_AddrCredentials_Cc) isTxCredential_AddrCredentials_Cred()   {}

func (m *TxCredential_AddrCredentials) GetCred() isTxCredential_AddrCredentials_Cred {
	if m != nil {
		return m.Cred
	}
	return nil
}

func (m *TxCredential_AddrCredentials) GetUser() *TxCredential_UserCredential {
	if x, ok := m.GetCred().(*TxCredential_AddrCredentials_User); ok {
		return x.User
	}
	return nil
}

func (m *TxCredential_AddrCredentials) GetCc() *TxCredential_InnerCredential {
	if x, ok := m.GetCred().(*TxCredential_AddrCredentials_Cc); ok {
		return x.Cc
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*TxCredential_AddrCredentials) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _TxCredential_AddrCredentials_OneofMarshaler, _TxCredential_AddrCredentials_OneofUnmarshaler, _TxCredential_AddrCredentials_OneofSizer, []interface{}{
		(*TxCredential_AddrCredentials_User)(nil),
		(*TxCredential_AddrCredentials_Cc)(nil),
	}
}

func _TxCredential_AddrCredentials_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*TxCredential_AddrCredentials)
	// cred
	switch x := m.Cred.(type) {
	case *TxCredential_AddrCredentials_User:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.User); err != nil {
			return err
		}
	case *TxCredential_AddrCredentials_Cc:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Cc); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("TxCredential_AddrCredentials.Cred has unexpected type %T", x)
	}
	return nil
}

func _TxCredential_AddrCredentials_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*TxCredential_AddrCredentials)
	switch tag {
	case 1: // cred.user
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(TxCredential_UserCredential)
		err := b.DecodeMessage(msg)
		m.Cred = &TxCredential_AddrCredentials_User{msg}
		return true, err
	case 2: // cred.cc
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(TxCredential_InnerCredential)
		err := b.DecodeMessage(msg)
		m.Cred = &TxCredential_AddrCredentials_Cc{msg}
		return true, err
	default:
		return false, nil
	}
}

func _TxCredential_AddrCredentials_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*TxCredential_AddrCredentials)
	// cred
	switch x := m.Cred.(type) {
	case *TxCredential_AddrCredentials_User:
		s := proto.Size(x.User)
		n += proto.SizeVarint(1<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *TxCredential_AddrCredentials_Cc:
		s := proto.Size(x.Cc)
		n += proto.SizeVarint(2<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

func init() {
	proto.RegisterType((*TxBase)(nil), "protos.TxBase")
	proto.RegisterType((*TxHeader)(nil), "protos.TxHeader")
	proto.RegisterType((*TxCredential)(nil), "protos.TxCredential")
	proto.RegisterType((*TxCredential_UserCredential)(nil), "protos.TxCredential.UserCredential")
	proto.RegisterType((*TxCredential_InnerCredential)(nil), "protos.TxCredential.InnerCredential")
	proto.RegisterType((*TxCredential_AddrCredentials)(nil), "protos.TxCredential.AddrCredentials")
}

func init() { proto.RegisterFile("tx.proto", fileDescriptor1) }

var fileDescriptor1 = []byte{
	// 389 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x92, 0xc1, 0x6e, 0xd4, 0x30,
	0x10, 0x86, 0x9b, 0xec, 0x36, 0x74, 0x27, 0x51, 0x2b, 0x0c, 0x42, 0x51, 0x2e, 0x2c, 0x81, 0x43,
	0x4f, 0xa9, 0x58, 0x24, 0x04, 0xdc, 0x68, 0x85, 0xb4, 0xa8, 0x17, 0x64, 0xc2, 0x03, 0x38, 0xf6,
	0x6c, 0x88, 0x76, 0x63, 0x47, 0xb6, 0x23, 0xd2, 0x0b, 0x4f, 0x80, 0x78, 0x66, 0x84, 0x9d, 0x28,
	0xdd, 0xaa, 0xa7, 0xe8, 0x9f, 0xf9, 0x3d, 0xf3, 0xcd, 0xaf, 0xc0, 0x99, 0x1d, 0x8a, 0x4e, 0x2b,
	0xab, 0x48, 0xe4, 0x3e, 0x26, 0x4b, 0xb8, 0xbe, 0xeb, 0xac, 0xf2, 0xd5, 0xec, 0x65, 0xad, 0x54,
	0x7d, 0xc0, 0x2b, 0xa7, 0xaa, 0x7e, 0x77, 0x65, 0x9b, 0x16, 0x8d, 0x65, 0x6d, 0xe7, 0x0d, 0x39,
	0x85, 0xa8, 0x1c, 0xae, 0x99, 0x41, 0x92, 0xc2, 0x13, 0x89, 0xf6, 0x97, 0xd2, 0xfb, 0x34, 0x58,
	0x07, 0x97, 0x2b, 0x3a, 0x49, 0xf2, 0x02, 0x22, 0xce, 0x25, 0x6b, 0x31, 0x0d, 0x5d, 0x63, 0x54,
	0xff, 0xeb, 0x2d, 0xda, 0x9f, 0x4a, 0xa4, 0x0b, 0x5f, 0xf7, 0x2a, 0xff, 0x0d, 0x67, 0xe5, 0xb0,
	0x45, 0x26, 0x50, 0x93, 0x1c, 0x96, 0x15, 0x33, 0xe8, 0x46, 0xc6, 0x9b, 0x73, 0xbf, 0xd5, 0x14,
	0x7e, 0x27, 0x75, 0x3d, 0xf2, 0x01, 0x56, 0x38, 0x74, 0x8d, 0x46, 0x51, 0x1a, 0xb7, 0x22, 0xde,
	0x64, 0x85, 0x07, 0x2f, 0x26, 0xf0, 0xa2, 0x9c, 0xc0, 0xe9, 0x6c, 0x26, 0xcf, 0xe1, 0x54, 0x2a,
	0xc9, 0xd1, 0x01, 0x24, 0xd4, 0x8b, 0xfc, 0xef, 0x02, 0x92, 0x72, 0xb8, 0xd1, 0x28, 0x50, 0xda,
	0x86, 0x1d, 0xc8, 0x27, 0x38, 0x65, 0x42, 0x68, 0x9e, 0x06, 0xeb, 0xc5, 0x65, 0xbc, 0x79, 0x33,
	0x53, 0xcc, 0xa6, 0xe2, 0xb3, 0x10, 0x7a, 0x96, 0x86, 0xfa, 0x27, 0xd9, 0x0e, 0xce, 0x7f, 0x18,
	0xbc, 0xd7, 0x21, 0x6f, 0x61, 0x65, 0x9a, 0x5a, 0x32, 0xdb, 0xeb, 0xe9, 0xae, 0x67, 0xd3, 0xc4,
	0x2f, 0x37, 0xdf, 0xa7, 0x16, 0x9d, 0x5d, 0xe4, 0x15, 0x84, 0xdd, 0x7e, 0x3c, 0xed, 0xe9, 0xe4,
	0xfd, 0xd6, 0x57, 0x87, 0x86, 0xdf, 0xe2, 0x1d, 0x0d, 0xbb, 0x7d, 0x76, 0x0b, 0x17, 0x5f, 0xa5,
	0x3c, 0x5a, 0x34, 0xe7, 0x1e, 0x1c, 0xe5, 0xbe, 0x86, 0x78, 0xd7, 0xc8, 0x1a, 0x75, 0xa7, 0x1b,
	0x69, 0xdd, 0xd8, 0x84, 0xde, 0x2f, 0x65, 0x7f, 0x02, 0xb8, 0x78, 0x70, 0x0f, 0xf9, 0x08, 0xcb,
	0xde, 0xa0, 0x1e, 0x89, 0x5f, 0x3f, 0x9a, 0xc1, 0xf1, 0xa5, 0xdb, 0x13, 0xea, 0x9e, 0x90, 0xf7,
	0x10, 0x72, 0x3e, 0xe2, 0x3f, 0x1e, 0xde, 0x03, 0xf4, 0xed, 0x09, 0x0d, 0x39, 0xbf, 0x8e, 0x60,
	0xc9, 0x35, 0x8a, 0xca, 0xff, 0x9b, 0xef, 0xfe, 0x05, 0x00, 0x00, 0xff, 0xff, 0x27, 0x29, 0x67,
	0x4f, 0xae, 0x02, 0x00, 0x00,
}
