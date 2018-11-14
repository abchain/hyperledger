// Code generated by protoc-gen-go. DO NOT EDIT.
// source: events.proto

package protos

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

type EventType int32

const (
	EventType_REGISTER  EventType = 0
	EventType_BLOCK     EventType = 1
	EventType_CHAINCODE EventType = 2
	EventType_REJECTION EventType = 3
)

var EventType_name = map[int32]string{
	0: "REGISTER",
	1: "BLOCK",
	2: "CHAINCODE",
	3: "REJECTION",
}
var EventType_value = map[string]int32{
	"REGISTER":  0,
	"BLOCK":     1,
	"CHAINCODE": 2,
	"REJECTION": 3,
}

func (x EventType) String() string {
	return proto.EnumName(EventType_name, int32(x))
}
func (EventType) EnumDescriptor() ([]byte, []int) { return fileDescriptor4, []int{0} }

// ChaincodeReg is used for registering chaincode Interests
// when EventType is CHAINCODE
type ChaincodeReg struct {
	ChaincodeID string `protobuf:"bytes,1,opt,name=chaincodeID" json:"chaincodeID,omitempty"`
	EventName   string `protobuf:"bytes,2,opt,name=eventName" json:"eventName,omitempty"`
}

func (m *ChaincodeReg) Reset()                    { *m = ChaincodeReg{} }
func (m *ChaincodeReg) String() string            { return proto.CompactTextString(m) }
func (*ChaincodeReg) ProtoMessage()               {}
func (*ChaincodeReg) Descriptor() ([]byte, []int) { return fileDescriptor4, []int{0} }

func (m *ChaincodeReg) GetChaincodeID() string {
	if m != nil {
		return m.ChaincodeID
	}
	return ""
}

func (m *ChaincodeReg) GetEventName() string {
	if m != nil {
		return m.EventName
	}
	return ""
}

type Interest struct {
	EventType EventType `protobuf:"varint,1,opt,name=eventType,enum=protos.EventType" json:"eventType,omitempty"`
	// Ideally we should just have the following oneof for different
	// Reg types and get rid of EventType. But this is an API change
	// Additional Reg types may add messages specific to their type
	// to the oneof.
	//
	// Types that are valid to be assigned to RegInfo:
	//	*Interest_ChaincodeRegInfo
	RegInfo isInterest_RegInfo `protobuf_oneof:"RegInfo"`
}

func (m *Interest) Reset()                    { *m = Interest{} }
func (m *Interest) String() string            { return proto.CompactTextString(m) }
func (*Interest) ProtoMessage()               {}
func (*Interest) Descriptor() ([]byte, []int) { return fileDescriptor4, []int{1} }

type isInterest_RegInfo interface {
	isInterest_RegInfo()
}

type Interest_ChaincodeRegInfo struct {
	ChaincodeRegInfo *ChaincodeReg `protobuf:"bytes,2,opt,name=chaincodeRegInfo,oneof"`
}

func (*Interest_ChaincodeRegInfo) isInterest_RegInfo() {}

func (m *Interest) GetRegInfo() isInterest_RegInfo {
	if m != nil {
		return m.RegInfo
	}
	return nil
}

func (m *Interest) GetEventType() EventType {
	if m != nil {
		return m.EventType
	}
	return EventType_REGISTER
}

func (m *Interest) GetChaincodeRegInfo() *ChaincodeReg {
	if x, ok := m.GetRegInfo().(*Interest_ChaincodeRegInfo); ok {
		return x.ChaincodeRegInfo
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*Interest) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _Interest_OneofMarshaler, _Interest_OneofUnmarshaler, _Interest_OneofSizer, []interface{}{
		(*Interest_ChaincodeRegInfo)(nil),
	}
}

func _Interest_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*Interest)
	// RegInfo
	switch x := m.RegInfo.(type) {
	case *Interest_ChaincodeRegInfo:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.ChaincodeRegInfo); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("Interest.RegInfo has unexpected type %T", x)
	}
	return nil
}

func _Interest_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*Interest)
	switch tag {
	case 2: // RegInfo.chaincodeRegInfo
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ChaincodeReg)
		err := b.DecodeMessage(msg)
		m.RegInfo = &Interest_ChaincodeRegInfo{msg}
		return true, err
	default:
		return false, nil
	}
}

func _Interest_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*Interest)
	// RegInfo
	switch x := m.RegInfo.(type) {
	case *Interest_ChaincodeRegInfo:
		s := proto.Size(x.ChaincodeRegInfo)
		n += proto.SizeVarint(2<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

// ---------- consumer events ---------
// Register is sent by consumers for registering events
// string type - "register"
type Register struct {
	Events []*Interest `protobuf:"bytes,1,rep,name=events" json:"events,omitempty"`
}

func (m *Register) Reset()                    { *m = Register{} }
func (m *Register) String() string            { return proto.CompactTextString(m) }
func (*Register) ProtoMessage()               {}
func (*Register) Descriptor() ([]byte, []int) { return fileDescriptor4, []int{2} }

func (m *Register) GetEvents() []*Interest {
	if m != nil {
		return m.Events
	}
	return nil
}

// Rejection is sent by consumers for erroneous transaction rejection events
// string type - "rejection"
type Rejection struct {
	Tx       *Transaction `protobuf:"bytes,1,opt,name=tx" json:"tx,omitempty"`
	ErrorMsg string       `protobuf:"bytes,2,opt,name=errorMsg" json:"errorMsg,omitempty"`
	Txid     string       `protobuf:"bytes,3,opt,name=txid" json:"txid,omitempty"`
}

func (m *Rejection) Reset()                    { *m = Rejection{} }
func (m *Rejection) String() string            { return proto.CompactTextString(m) }
func (*Rejection) ProtoMessage()               {}
func (*Rejection) Descriptor() ([]byte, []int) { return fileDescriptor4, []int{3} }

func (m *Rejection) GetTx() *Transaction {
	if m != nil {
		return m.Tx
	}
	return nil
}

func (m *Rejection) GetErrorMsg() string {
	if m != nil {
		return m.ErrorMsg
	}
	return ""
}

func (m *Rejection) GetTxid() string {
	if m != nil {
		return m.Txid
	}
	return ""
}

// ---------- producer events ---------
type Unregister struct {
	Events []*Interest `protobuf:"bytes,1,rep,name=events" json:"events,omitempty"`
}

func (m *Unregister) Reset()                    { *m = Unregister{} }
func (m *Unregister) String() string            { return proto.CompactTextString(m) }
func (*Unregister) ProtoMessage()               {}
func (*Unregister) Descriptor() ([]byte, []int) { return fileDescriptor4, []int{4} }

func (m *Unregister) GetEvents() []*Interest {
	if m != nil {
		return m.Events
	}
	return nil
}

// Event is used by
//  - consumers (adapters) to send Register
//  - producer to advertise supported types and events
type Event struct {
	// Types that are valid to be assigned to Event:
	//	*Event_Register
	//	*Event_Block
	//	*Event_ChaincodeEvent
	//	*Event_Rejection
	//	*Event_Unregister
	Event isEvent_Event `protobuf_oneof:"Event"`
}

func (m *Event) Reset()                    { *m = Event{} }
func (m *Event) String() string            { return proto.CompactTextString(m) }
func (*Event) ProtoMessage()               {}
func (*Event) Descriptor() ([]byte, []int) { return fileDescriptor4, []int{5} }

type isEvent_Event interface {
	isEvent_Event()
}

type Event_Register struct {
	Register *Register `protobuf:"bytes,1,opt,name=register,oneof"`
}
type Event_Block struct {
	Block *Block `protobuf:"bytes,2,opt,name=block,oneof"`
}
type Event_ChaincodeEvent struct {
	ChaincodeEvent *ChaincodeEvent `protobuf:"bytes,3,opt,name=chaincodeEvent,oneof"`
}
type Event_Rejection struct {
	Rejection *Rejection `protobuf:"bytes,4,opt,name=rejection,oneof"`
}
type Event_Unregister struct {
	Unregister *Unregister `protobuf:"bytes,5,opt,name=unregister,oneof"`
}

func (*Event_Register) isEvent_Event()       {}
func (*Event_Block) isEvent_Event()          {}
func (*Event_ChaincodeEvent) isEvent_Event() {}
func (*Event_Rejection) isEvent_Event()      {}
func (*Event_Unregister) isEvent_Event()     {}

func (m *Event) GetEvent() isEvent_Event {
	if m != nil {
		return m.Event
	}
	return nil
}

func (m *Event) GetRegister() *Register {
	if x, ok := m.GetEvent().(*Event_Register); ok {
		return x.Register
	}
	return nil
}

func (m *Event) GetBlock() *Block {
	if x, ok := m.GetEvent().(*Event_Block); ok {
		return x.Block
	}
	return nil
}

func (m *Event) GetChaincodeEvent() *ChaincodeEvent {
	if x, ok := m.GetEvent().(*Event_ChaincodeEvent); ok {
		return x.ChaincodeEvent
	}
	return nil
}

func (m *Event) GetRejection() *Rejection {
	if x, ok := m.GetEvent().(*Event_Rejection); ok {
		return x.Rejection
	}
	return nil
}

func (m *Event) GetUnregister() *Unregister {
	if x, ok := m.GetEvent().(*Event_Unregister); ok {
		return x.Unregister
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*Event) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _Event_OneofMarshaler, _Event_OneofUnmarshaler, _Event_OneofSizer, []interface{}{
		(*Event_Register)(nil),
		(*Event_Block)(nil),
		(*Event_ChaincodeEvent)(nil),
		(*Event_Rejection)(nil),
		(*Event_Unregister)(nil),
	}
}

func _Event_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*Event)
	// Event
	switch x := m.Event.(type) {
	case *Event_Register:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Register); err != nil {
			return err
		}
	case *Event_Block:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Block); err != nil {
			return err
		}
	case *Event_ChaincodeEvent:
		b.EncodeVarint(3<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.ChaincodeEvent); err != nil {
			return err
		}
	case *Event_Rejection:
		b.EncodeVarint(4<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Rejection); err != nil {
			return err
		}
	case *Event_Unregister:
		b.EncodeVarint(5<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Unregister); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("Event.Event has unexpected type %T", x)
	}
	return nil
}

func _Event_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*Event)
	switch tag {
	case 1: // Event.register
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Register)
		err := b.DecodeMessage(msg)
		m.Event = &Event_Register{msg}
		return true, err
	case 2: // Event.block
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Block)
		err := b.DecodeMessage(msg)
		m.Event = &Event_Block{msg}
		return true, err
	case 3: // Event.chaincodeEvent
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ChaincodeEvent)
		err := b.DecodeMessage(msg)
		m.Event = &Event_ChaincodeEvent{msg}
		return true, err
	case 4: // Event.rejection
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Rejection)
		err := b.DecodeMessage(msg)
		m.Event = &Event_Rejection{msg}
		return true, err
	case 5: // Event.unregister
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Unregister)
		err := b.DecodeMessage(msg)
		m.Event = &Event_Unregister{msg}
		return true, err
	default:
		return false, nil
	}
}

func _Event_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*Event)
	// Event
	switch x := m.Event.(type) {
	case *Event_Register:
		s := proto.Size(x.Register)
		n += proto.SizeVarint(1<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Event_Block:
		s := proto.Size(x.Block)
		n += proto.SizeVarint(2<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Event_ChaincodeEvent:
		s := proto.Size(x.ChaincodeEvent)
		n += proto.SizeVarint(3<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Event_Rejection:
		s := proto.Size(x.Rejection)
		n += proto.SizeVarint(4<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Event_Unregister:
		s := proto.Size(x.Unregister)
		n += proto.SizeVarint(5<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

func init() {
	proto.RegisterType((*ChaincodeReg)(nil), "protos.ChaincodeReg")
	proto.RegisterType((*Interest)(nil), "protos.Interest")
	proto.RegisterType((*Register)(nil), "protos.Register")
	proto.RegisterType((*Rejection)(nil), "protos.Rejection")
	proto.RegisterType((*Unregister)(nil), "protos.Unregister")
	proto.RegisterType((*Event)(nil), "protos.Event")
	proto.RegisterEnum("protos.EventType", EventType_name, EventType_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Events service

type EventsClient interface {
	// event chatting using Event
	Chat(ctx context.Context, opts ...grpc.CallOption) (Events_ChatClient, error)
}

type eventsClient struct {
	cc *grpc.ClientConn
}

func NewEventsClient(cc *grpc.ClientConn) EventsClient {
	return &eventsClient{cc}
}

func (c *eventsClient) Chat(ctx context.Context, opts ...grpc.CallOption) (Events_ChatClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Events_serviceDesc.Streams[0], c.cc, "/protos.Events/Chat", opts...)
	if err != nil {
		return nil, err
	}
	x := &eventsChatClient{stream}
	return x, nil
}

type Events_ChatClient interface {
	Send(*Event) error
	Recv() (*Event, error)
	grpc.ClientStream
}

type eventsChatClient struct {
	grpc.ClientStream
}

func (x *eventsChatClient) Send(m *Event) error {
	return x.ClientStream.SendMsg(m)
}

func (x *eventsChatClient) Recv() (*Event, error) {
	m := new(Event)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for Events service

type EventsServer interface {
	// event chatting using Event
	Chat(Events_ChatServer) error
}

func RegisterEventsServer(s *grpc.Server, srv EventsServer) {
	s.RegisterService(&_Events_serviceDesc, srv)
}

func _Events_Chat_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(EventsServer).Chat(&eventsChatServer{stream})
}

type Events_ChatServer interface {
	Send(*Event) error
	Recv() (*Event, error)
	grpc.ServerStream
}

type eventsChatServer struct {
	grpc.ServerStream
}

func (x *eventsChatServer) Send(m *Event) error {
	return x.ServerStream.SendMsg(m)
}

func (x *eventsChatServer) Recv() (*Event, error) {
	m := new(Event)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _Events_serviceDesc = grpc.ServiceDesc{
	ServiceName: "protos.Events",
	HandlerType: (*EventsServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Chat",
			Handler:       _Events_Chat_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "events.proto",
}

func init() { proto.RegisterFile("events.proto", fileDescriptor4) }

var fileDescriptor4 = []byte{
	// 459 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x53, 0x5d, 0x8f, 0x93, 0x40,
	0x14, 0x05, 0xfa, 0xb1, 0x70, 0x69, 0x37, 0xec, 0x75, 0x63, 0x48, 0xe3, 0x43, 0x83, 0x31, 0x69,
	0xf6, 0xa1, 0x2a, 0x36, 0x3e, 0x2b, 0x2c, 0x11, 0xfc, 0x68, 0x93, 0xb1, 0xbe, 0x4b, 0xd9, 0xd9,
	0x2e, 0x7e, 0xc0, 0x66, 0x18, 0x4d, 0xfd, 0x0b, 0x3e, 0xfa, 0x8b, 0xcd, 0x5e, 0x18, 0x60, 0xed,
	0x93, 0x4f, 0xcc, 0xbd, 0xe7, 0x9c, 0x7b, 0x0f, 0x27, 0x33, 0x30, 0xe1, 0x3f, 0x79, 0x21, 0xab,
	0xe5, 0xad, 0x28, 0x65, 0x89, 0x63, 0xfa, 0x54, 0xb3, 0xf3, 0xec, 0x26, 0xcd, 0x8b, 0xac, 0xbc,
	0xe2, 0x04, 0xd7, 0xe8, 0x6c, 0x72, 0x9d, 0xee, 0x44, 0x9e, 0xd5, 0x95, 0xb7, 0x86, 0x49, 0xa8,
	0x58, 0x8c, 0xef, 0x71, 0x0e, 0x76, 0xab, 0x4a, 0x2e, 0x5d, 0x7d, 0xae, 0x2f, 0x2c, 0xd6, 0x6f,
	0xe1, 0x23, 0xb0, 0x68, 0xdc, 0x3a, 0xfd, 0xce, 0x5d, 0x83, 0xf0, 0xae, 0xe1, 0xfd, 0xd6, 0xc1,
	0x4c, 0x0a, 0xc9, 0x05, 0xaf, 0x24, 0x3e, 0x6d, 0xa8, 0xdb, 0x5f, 0xb7, 0x9c, 0x46, 0x9d, 0xfa,
	0x67, 0xf5, 0xde, 0x6a, 0x19, 0x29, 0x80, 0x75, 0x1c, 0x0c, 0xc0, 0xc9, 0x7a, 0x6e, 0x92, 0xe2,
	0xba, 0xa4, 0x15, 0xb6, 0x7f, 0xae, 0x74, 0x7d, 0xb7, 0xb1, 0xc6, 0x8e, 0xf8, 0x81, 0x05, 0x27,
	0xcd, 0xd1, 0x5b, 0x81, 0xc9, 0xf8, 0x3e, 0xaf, 0x24, 0x17, 0xb8, 0x80, 0x71, 0x1d, 0x92, 0xab,
	0xcf, 0x07, 0x0b, 0xdb, 0x77, 0xd4, 0x40, 0xe5, 0x96, 0x35, 0xb8, 0xf7, 0x19, 0x2c, 0xc6, 0xbf,
	0xf0, 0x4c, 0xe6, 0x65, 0x81, 0x8f, 0xc1, 0x90, 0x07, 0xf2, 0x6e, 0xfb, 0x0f, 0x94, 0x64, 0x2b,
	0xd2, 0xa2, 0x4a, 0x89, 0xc0, 0x0c, 0x79, 0xc0, 0x19, 0x98, 0x5c, 0x88, 0x52, 0x7c, 0xa8, 0xf6,
	0x4d, 0x22, 0x6d, 0x8d, 0x08, 0x43, 0x79, 0xc8, 0xaf, 0xdc, 0x01, 0xf5, 0xe9, 0xec, 0xbd, 0x04,
	0xf8, 0x54, 0x88, 0xff, 0x77, 0xf6, 0xc7, 0x80, 0x11, 0xe5, 0x86, 0x4b, 0x30, 0x95, 0xbe, 0x31,
	0xd7, 0xaa, 0xd4, 0x1f, 0xc7, 0x1a, 0x6b, 0x39, 0xf8, 0x04, 0x46, 0xbb, 0x6f, 0x65, 0xf6, 0xb5,
	0x49, 0x73, 0xaa, 0xc8, 0xc1, 0x5d, 0x33, 0xd6, 0x58, 0x8d, 0xe2, 0x2b, 0x38, 0x6d, 0xf3, 0xa4,
	0x45, 0x64, 0xdb, 0xf6, 0x1f, 0x1e, 0xa5, 0x4f, 0x68, 0xac, 0xb1, 0x7f, 0xf8, 0xf8, 0x1c, 0x2c,
	0xa1, 0xc2, 0x73, 0x87, 0x24, 0x3e, 0xeb, 0x9c, 0x35, 0x40, 0xac, 0xb1, 0x8e, 0x85, 0x2b, 0x80,
	0x1f, 0x6d, 0x1a, 0xee, 0x88, 0x34, 0xa8, 0x34, 0x5d, 0x4e, 0xb1, 0xc6, 0x7a, 0xbc, 0xe0, 0xa4,
	0x89, 0xe2, 0x22, 0x00, 0xab, 0xbd, 0x4b, 0x38, 0x01, 0x93, 0x45, 0x6f, 0x92, 0x8f, 0xdb, 0x88,
	0x39, 0x1a, 0x5a, 0x30, 0x0a, 0xde, 0x6f, 0xc2, 0x77, 0x8e, 0x8e, 0x53, 0xb0, 0xc2, 0xf8, 0x75,
	0xb2, 0x0e, 0x37, 0x97, 0x91, 0x63, 0xdc, 0x95, 0x2c, 0x7a, 0x1b, 0x85, 0xdb, 0x64, 0xb3, 0x76,
	0x06, 0xfe, 0x0a, 0xc6, 0x34, 0xa3, 0xc2, 0x0b, 0x18, 0x86, 0x37, 0xa9, 0xc4, 0xe9, 0xbd, 0x7b,
	0x3a, 0xbb, 0x5f, 0x7a, 0xda, 0x42, 0x7f, 0xa6, 0xef, 0xea, 0x77, 0xf6, 0xe2, 0x6f, 0x00, 0x00,
	0x00, 0xff, 0xff, 0x85, 0x98, 0x5e, 0xe7, 0x7e, 0x03, 0x00, 0x00,
}
