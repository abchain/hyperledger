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
	Signature *Signature `protobuf:"bytes,3,opt,name=signature" json:"signature,omitempty"`
}

func (m *TxCredential_UserCredential) Reset()                    { *m = TxCredential_UserCredential{} }
func (m *TxCredential_UserCredential) String() string            { return proto.CompactTextString(m) }
func (*TxCredential_UserCredential) ProtoMessage()               {}
func (*TxCredential_UserCredential) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{2, 0} }

func (m *TxCredential_UserCredential) GetSignature() *Signature {
	if m != nil {
		return m.Signature
	}
	return nil
}

type TxCredential_DataCredential struct {
	Key string `protobuf:"bytes,1,opt,name=key" json:"key,omitempty"`
	// Types that are valid to be assigned to Data:
	//	*TxCredential_DataCredential_Bts
	//	*TxCredential_DataCredential_Int
	//	*TxCredential_DataCredential_Str
	Data isTxCredential_DataCredential_Data `protobuf_oneof:"data"`
}

func (m *TxCredential_DataCredential) Reset()                    { *m = TxCredential_DataCredential{} }
func (m *TxCredential_DataCredential) String() string            { return proto.CompactTextString(m) }
func (*TxCredential_DataCredential) ProtoMessage()               {}
func (*TxCredential_DataCredential) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{2, 1} }

type isTxCredential_DataCredential_Data interface {
	isTxCredential_DataCredential_Data()
}

type TxCredential_DataCredential_Bts struct {
	Bts []byte `protobuf:"bytes,2,opt,name=bts,proto3,oneof"`
}
type TxCredential_DataCredential_Int struct {
	Int int32 `protobuf:"varint,3,opt,name=int,oneof"`
}
type TxCredential_DataCredential_Str struct {
	Str string `protobuf:"bytes,4,opt,name=str,oneof"`
}

func (*TxCredential_DataCredential_Bts) isTxCredential_DataCredential_Data() {}
func (*TxCredential_DataCredential_Int) isTxCredential_DataCredential_Data() {}
func (*TxCredential_DataCredential_Str) isTxCredential_DataCredential_Data() {}

func (m *TxCredential_DataCredential) GetData() isTxCredential_DataCredential_Data {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *TxCredential_DataCredential) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *TxCredential_DataCredential) GetBts() []byte {
	if x, ok := m.GetData().(*TxCredential_DataCredential_Bts); ok {
		return x.Bts
	}
	return nil
}

func (m *TxCredential_DataCredential) GetInt() int32 {
	if x, ok := m.GetData().(*TxCredential_DataCredential_Int); ok {
		return x.Int
	}
	return 0
}

func (m *TxCredential_DataCredential) GetStr() string {
	if x, ok := m.GetData().(*TxCredential_DataCredential_Str); ok {
		return x.Str
	}
	return ""
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*TxCredential_DataCredential) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _TxCredential_DataCredential_OneofMarshaler, _TxCredential_DataCredential_OneofUnmarshaler, _TxCredential_DataCredential_OneofSizer, []interface{}{
		(*TxCredential_DataCredential_Bts)(nil),
		(*TxCredential_DataCredential_Int)(nil),
		(*TxCredential_DataCredential_Str)(nil),
	}
}

func _TxCredential_DataCredential_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*TxCredential_DataCredential)
	// data
	switch x := m.Data.(type) {
	case *TxCredential_DataCredential_Bts:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		b.EncodeRawBytes(x.Bts)
	case *TxCredential_DataCredential_Int:
		b.EncodeVarint(3<<3 | proto.WireVarint)
		b.EncodeVarint(uint64(x.Int))
	case *TxCredential_DataCredential_Str:
		b.EncodeVarint(4<<3 | proto.WireBytes)
		b.EncodeStringBytes(x.Str)
	case nil:
	default:
		return fmt.Errorf("TxCredential_DataCredential.Data has unexpected type %T", x)
	}
	return nil
}

func _TxCredential_DataCredential_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*TxCredential_DataCredential)
	switch tag {
	case 2: // data.bts
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeRawBytes(true)
		m.Data = &TxCredential_DataCredential_Bts{x}
		return true, err
	case 3: // data.int
		if wire != proto.WireVarint {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeVarint()
		m.Data = &TxCredential_DataCredential_Int{int32(x)}
		return true, err
	case 4: // data.str
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeStringBytes()
		m.Data = &TxCredential_DataCredential_Str{x}
		return true, err
	default:
		return false, nil
	}
}

func _TxCredential_DataCredential_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*TxCredential_DataCredential)
	// data
	switch x := m.Data.(type) {
	case *TxCredential_DataCredential_Bts:
		n += proto.SizeVarint(2<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(len(x.Bts)))
		n += len(x.Bts)
	case *TxCredential_DataCredential_Int:
		n += proto.SizeVarint(3<<3 | proto.WireVarint)
		n += proto.SizeVarint(uint64(x.Int))
	case *TxCredential_DataCredential_Str:
		n += proto.SizeVarint(4<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(len(x.Str)))
		n += len(x.Str)
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type TxCredential_AddrCredentials struct {
	// Types that are valid to be assigned to Cred:
	//	*TxCredential_AddrCredentials_User
	//	*TxCredential_AddrCredentials_Data
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
type TxCredential_AddrCredentials_Data struct {
	Data *TxCredential_DataCredential `protobuf:"bytes,3,opt,name=data,oneof"`
}

func (*TxCredential_AddrCredentials_User) isTxCredential_AddrCredentials_Cred() {}
func (*TxCredential_AddrCredentials_Data) isTxCredential_AddrCredentials_Cred() {}

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

func (m *TxCredential_AddrCredentials) GetData() *TxCredential_DataCredential {
	if x, ok := m.GetCred().(*TxCredential_AddrCredentials_Data); ok {
		return x.Data
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*TxCredential_AddrCredentials) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _TxCredential_AddrCredentials_OneofMarshaler, _TxCredential_AddrCredentials_OneofUnmarshaler, _TxCredential_AddrCredentials_OneofSizer, []interface{}{
		(*TxCredential_AddrCredentials_User)(nil),
		(*TxCredential_AddrCredentials_Data)(nil),
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
	case *TxCredential_AddrCredentials_Data:
		b.EncodeVarint(3<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Data); err != nil {
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
	case 3: // cred.data
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(TxCredential_DataCredential)
		err := b.DecodeMessage(msg)
		m.Cred = &TxCredential_AddrCredentials_Data{msg}
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
	case *TxCredential_AddrCredentials_Data:
		s := proto.Size(x.Data)
		n += proto.SizeVarint(3<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type TxBatch struct {
	Txs []*TxBatchSubTx `protobuf:"bytes,1,rep,name=txs" json:"txs,omitempty"`
}

func (m *TxBatch) Reset()                    { *m = TxBatch{} }
func (m *TxBatch) String() string            { return proto.CompactTextString(m) }
func (*TxBatch) ProtoMessage()               {}
func (*TxBatch) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{3} }

func (m *TxBatch) GetTxs() []*TxBatchSubTx {
	if m != nil {
		return m.Txs
	}
	return nil
}

type TxBatchSubTx struct {
	Method  string `protobuf:"bytes,1,opt,name=method" json:"method,omitempty"`
	Payload []byte `protobuf:"bytes,2,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (m *TxBatchSubTx) Reset()                    { *m = TxBatchSubTx{} }
func (m *TxBatchSubTx) String() string            { return proto.CompactTextString(m) }
func (*TxBatchSubTx) ProtoMessage()               {}
func (*TxBatchSubTx) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{3, 0} }

func (m *TxBatchSubTx) GetMethod() string {
	if m != nil {
		return m.Method
	}
	return ""
}

func (m *TxBatchSubTx) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

type TxBatchResp struct {
	Response [][]byte `protobuf:"bytes,1,rep,name=response,proto3" json:"response,omitempty"`
}

func (m *TxBatchResp) Reset()                    { *m = TxBatchResp{} }
func (m *TxBatchResp) String() string            { return proto.CompactTextString(m) }
func (*TxBatchResp) ProtoMessage()               {}
func (*TxBatchResp) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{4} }

func (m *TxBatchResp) GetResponse() [][]byte {
	if m != nil {
		return m.Response
	}
	return nil
}

func init() {
	proto.RegisterType((*TxBase)(nil), "protos.TxBase")
	proto.RegisterType((*TxHeader)(nil), "protos.TxHeader")
	proto.RegisterType((*TxCredential)(nil), "protos.TxCredential")
	proto.RegisterType((*TxCredential_UserCredential)(nil), "protos.TxCredential.UserCredential")
	proto.RegisterType((*TxCredential_DataCredential)(nil), "protos.TxCredential.DataCredential")
	proto.RegisterType((*TxCredential_AddrCredentials)(nil), "protos.TxCredential.AddrCredentials")
	proto.RegisterType((*TxBatch)(nil), "protos.TxBatch")
	proto.RegisterType((*TxBatchSubTx)(nil), "protos.TxBatch.subTx")
	proto.RegisterType((*TxBatchResp)(nil), "protos.TxBatchResp")
}

func init() { proto.RegisterFile("tx.proto", fileDescriptor1) }

var fileDescriptor1 = []byte{
	// 480 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x53, 0xc1, 0x8e, 0xd3, 0x30,
	0x10, 0xdd, 0x92, 0xb6, 0xdb, 0x4e, 0xab, 0x02, 0x16, 0xa0, 0x28, 0x07, 0xa8, 0x02, 0x12, 0xe5,
	0x92, 0x4a, 0xe5, 0xc2, 0x72, 0xdb, 0xc2, 0xa1, 0x67, 0x13, 0x2e, 0xdc, 0x9c, 0x78, 0xb6, 0x89,
	0xb6, 0xb1, 0x23, 0xdb, 0x15, 0xe9, 0x85, 0x2b, 0x47, 0x7e, 0x19, 0xd9, 0x4e, 0x48, 0x16, 0xed,
	0x29, 0x79, 0x6f, 0xe6, 0xcd, 0xcc, 0xb3, 0xc7, 0x30, 0x33, 0x4d, 0x52, 0x2b, 0x69, 0x24, 0x99,
	0xba, 0x8f, 0x8e, 0x96, 0xb9, 0xba, 0xd4, 0x46, 0x7a, 0x36, 0x7a, 0x73, 0x94, 0xf2, 0x78, 0xc2,
	0xad, 0x43, 0xd9, 0xf9, 0x6e, 0x6b, 0xca, 0x0a, 0xb5, 0x61, 0x55, 0xed, 0x13, 0x62, 0x0a, 0xd3,
	0xb4, 0xd9, 0x33, 0x8d, 0x24, 0x84, 0x6b, 0x81, 0xe6, 0xa7, 0x54, 0xf7, 0xe1, 0x68, 0x3d, 0xda,
	0xcc, 0x69, 0x07, 0xc9, 0x2b, 0x98, 0xe6, 0xb9, 0x60, 0x15, 0x86, 0x4f, 0x5c, 0xa0, 0x45, 0x96,
	0xaf, 0xd0, 0x14, 0x92, 0x87, 0x81, 0xe7, 0x3d, 0x8a, 0x7f, 0xc1, 0x2c, 0x6d, 0x0e, 0xc8, 0x38,
	0x2a, 0x12, 0xc3, 0x38, 0x63, 0x1a, 0x5d, 0xc9, 0xc5, 0x6e, 0xe5, 0xbb, 0xea, 0xc4, 0xf7, 0xa4,
	0x2e, 0x46, 0x3e, 0xc1, 0x1c, 0x9b, 0xba, 0x54, 0xc8, 0x53, 0xed, 0x5a, 0x2c, 0x76, 0x51, 0xe2,
	0x07, 0x4f, 0xba, 0xc1, 0x93, 0xb4, 0x1b, 0x9c, 0xf6, 0xc9, 0xe4, 0x05, 0x4c, 0x84, 0x14, 0x39,
	0xba, 0x01, 0x96, 0xd4, 0x83, 0xf8, 0x77, 0x00, 0xcb, 0xb4, 0xf9, 0xa2, 0x90, 0xa3, 0x30, 0x25,
	0x3b, 0x91, 0xcf, 0x30, 0x61, 0x9c, 0xab, 0x3c, 0x1c, 0xad, 0x83, 0xcd, 0x62, 0xf7, 0xae, 0x9f,
	0xa2, 0x4f, 0x4a, 0x6e, 0x39, 0x57, 0x3d, 0xd4, 0xd4, 0x4b, 0xa2, 0x5b, 0x58, 0x7d, 0xd7, 0x38,
	0x88, 0x90, 0x2d, 0xcc, 0x75, 0x79, 0x14, 0xcc, 0x9c, 0x95, 0x6f, 0xbc, 0xd8, 0x3d, 0xef, 0x2a,
	0x7e, 0xeb, 0x02, 0xb4, 0xcf, 0x89, 0xee, 0x60, 0xf5, 0x95, 0x19, 0x36, 0x28, 0xf1, 0x0c, 0x82,
	0x7b, 0xbc, 0xb4, 0xe7, 0x6c, 0x7f, 0x09, 0x81, 0x20, 0x33, 0xde, 0xfd, 0xf2, 0x70, 0x45, 0x2d,
	0xb0, 0x5c, 0x29, 0x8c, 0x6b, 0x31, 0xb1, 0x5c, 0x29, 0x8c, 0xe5, 0xb4, 0x51, 0xe1, 0xd8, 0x2a,
	0x2d, 0xa7, 0x8d, 0xda, 0x4f, 0x61, 0xcc, 0x99, 0x61, 0xd1, 0x9f, 0x11, 0x3c, 0xfd, 0xcf, 0x05,
	0xb9, 0x81, 0xf1, 0x59, 0xa3, 0x6a, 0xcf, 0xff, 0xed, 0xa3, 0xce, 0x1f, 0xfa, 0x3b, 0x5c, 0x51,
	0x27, 0xb1, 0x52, 0x5b, 0xb6, 0xb5, 0xf8, 0xb8, 0xf4, 0xa1, 0x2f, 0x2b, 0xb5, 0x12, 0x3b, 0x51,
	0xae, 0x90, 0xc7, 0x15, 0x5c, 0xdb, 0x9b, 0x36, 0x79, 0x41, 0xde, 0x43, 0x60, 0x1a, 0xdd, 0xde,
	0xc0, 0xcb, 0xe1, 0x1e, 0x98, 0xbc, 0x48, 0xf4, 0x39, 0x4b, 0x1b, 0x6a, 0x33, 0xa2, 0x1b, 0x98,
	0x38, 0x34, 0x58, 0xaf, 0xd1, 0x70, 0xbd, 0xec, 0xa2, 0xd6, 0xec, 0x72, 0x92, 0x8c, 0xfb, 0xe3,
	0xa2, 0x1d, 0x8c, 0x3f, 0xc0, 0xa2, 0x2d, 0x48, 0x51, 0xd7, 0x24, 0x82, 0x99, 0x42, 0x5d, 0x4b,
	0xe1, 0xf6, 0x2f, 0xd8, 0x2c, 0xe9, 0x3f, 0xbc, 0x5f, 0xff, 0x78, 0x5d, 0x5c, 0x6a, 0x54, 0x27,
	0xe4, 0x47, 0x54, 0x09, 0xcb, 0xf2, 0x82, 0x95, 0x22, 0x91, 0xea, 0xe8, 0xdf, 0x8a, 0xce, 0xfc,
	0x83, 0xfa, 0xf8, 0x37, 0x00, 0x00, 0xff, 0xff, 0xe6, 0xd0, 0xc0, 0x79, 0x63, 0x03, 0x00, 0x00,
}
