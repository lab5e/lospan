// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        (unknown)
// source: lospan/entities.proto

package lospan

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// State of device
type DeviceState int32

const (
	DeviceState_UNSPECIFIED DeviceState = 0
	DeviceState_OTAA        DeviceState = 1
	DeviceState_ABP         DeviceState = 2
	DeviceState_DISABLED    DeviceState = 3
)

// Enum value maps for DeviceState.
var (
	DeviceState_name = map[int32]string{
		0: "UNSPECIFIED",
		1: "OTAA",
		2: "ABP",
		3: "DISABLED",
	}
	DeviceState_value = map[string]int32{
		"UNSPECIFIED": 0,
		"OTAA":        1,
		"ABP":         2,
		"DISABLED":    3,
	}
)

func (x DeviceState) Enum() *DeviceState {
	p := new(DeviceState)
	*p = x
	return p
}

func (x DeviceState) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (DeviceState) Descriptor() protoreflect.EnumDescriptor {
	return file_lospan_entities_proto_enumTypes[0].Descriptor()
}

func (DeviceState) Type() protoreflect.EnumType {
	return &file_lospan_entities_proto_enumTypes[0]
}

func (x DeviceState) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use DeviceState.Descriptor instead.
func (DeviceState) EnumDescriptor() ([]byte, []int) {
	return file_lospan_entities_proto_rawDescGZIP(), []int{0}
}

// Application is a logical construct on top of devices. Devices in the same application share the same
// application key
type Application struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Eui string  `protobuf:"bytes,1,opt,name=eui,proto3" json:"eui,omitempty"`
	Tag *string `protobuf:"bytes,2,opt,name=tag,proto3,oneof" json:"tag,omitempty"`
}

func (x *Application) Reset() {
	*x = Application{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lospan_entities_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Application) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Application) ProtoMessage() {}

func (x *Application) ProtoReflect() protoreflect.Message {
	mi := &file_lospan_entities_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Application.ProtoReflect.Descriptor instead.
func (*Application) Descriptor() ([]byte, []int) {
	return file_lospan_entities_proto_rawDescGZIP(), []int{0}
}

func (x *Application) GetEui() string {
	if x != nil {
		return x.Eui
	}
	return ""
}

func (x *Application) GetTag() string {
	if x != nil && x.Tag != nil {
		return *x.Tag
	}
	return ""
}

// Device is the ... device that connects to the gateway. "Node" might be a better name since it's
// part of the LoRaWAN implementation nomenclature.
type Device struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Eui               *string      `protobuf:"bytes,1,opt,name=eui,proto3,oneof" json:"eui,omitempty"` // Tecnhically required but can be skipped when creating a new
	ApplicationEui    *string      `protobuf:"bytes,2,opt,name=application_eui,json=applicationEui,proto3,oneof" json:"application_eui,omitempty"`
	State             *DeviceState `protobuf:"varint,3,opt,name=state,proto3,enum=lospan.DeviceState,oneof" json:"state,omitempty"`
	DevAddr           *uint32      `protobuf:"varint,4,opt,name=dev_addr,json=devAddr,proto3,oneof" json:"dev_addr,omitempty"`                                // 7+25 bits
	AppKey            []byte       `protobuf:"bytes,5,opt,name=app_key,json=appKey,proto3,oneof" json:"app_key,omitempty"`                                    // 16 bytes/256 bits
	AppSessionKey     []byte       `protobuf:"bytes,6,opt,name=app_session_key,json=appSessionKey,proto3,oneof" json:"app_session_key,omitempty"`             // 16 bytes/256 bits
	NetworkSessionKey []byte       `protobuf:"bytes,7,opt,name=network_session_key,json=networkSessionKey,proto3,oneof" json:"network_session_key,omitempty"` // 16 bytes/256 bits
	FrameCountUp      *int32       `protobuf:"varint,8,opt,name=frame_count_up,json=frameCountUp,proto3,oneof" json:"frame_count_up,omitempty"`               // in reality uint16
	FrameCountDown    *int32       `protobuf:"varint,9,opt,name=frame_count_down,json=frameCountDown,proto3,oneof" json:"frame_count_down,omitempty"`         // in reality uint16
	RelaxedCounter    *bool        `protobuf:"varint,10,opt,name=relaxed_counter,json=relaxedCounter,proto3,oneof" json:"relaxed_counter,omitempty"`
	KeyWarning        *bool        `protobuf:"varint,11,opt,name=key_warning,json=keyWarning,proto3,oneof" json:"key_warning,omitempty"` // Ignored on updates; set by service
	Tag               *string      `protobuf:"bytes,12,opt,name=tag,proto3,oneof" json:"tag,omitempty"`
	DevNonces         []int32      `protobuf:"varint,13,rep,packed,name=dev_nonces,json=devNonces,proto3" json:"dev_nonces,omitempty"` // in reality uint16
}

func (x *Device) Reset() {
	*x = Device{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lospan_entities_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Device) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Device) ProtoMessage() {}

func (x *Device) ProtoReflect() protoreflect.Message {
	mi := &file_lospan_entities_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Device.ProtoReflect.Descriptor instead.
func (*Device) Descriptor() ([]byte, []int) {
	return file_lospan_entities_proto_rawDescGZIP(), []int{1}
}

func (x *Device) GetEui() string {
	if x != nil && x.Eui != nil {
		return *x.Eui
	}
	return ""
}

func (x *Device) GetApplicationEui() string {
	if x != nil && x.ApplicationEui != nil {
		return *x.ApplicationEui
	}
	return ""
}

func (x *Device) GetState() DeviceState {
	if x != nil && x.State != nil {
		return *x.State
	}
	return DeviceState_UNSPECIFIED
}

func (x *Device) GetDevAddr() uint32 {
	if x != nil && x.DevAddr != nil {
		return *x.DevAddr
	}
	return 0
}

func (x *Device) GetAppKey() []byte {
	if x != nil {
		return x.AppKey
	}
	return nil
}

func (x *Device) GetAppSessionKey() []byte {
	if x != nil {
		return x.AppSessionKey
	}
	return nil
}

func (x *Device) GetNetworkSessionKey() []byte {
	if x != nil {
		return x.NetworkSessionKey
	}
	return nil
}

func (x *Device) GetFrameCountUp() int32 {
	if x != nil && x.FrameCountUp != nil {
		return *x.FrameCountUp
	}
	return 0
}

func (x *Device) GetFrameCountDown() int32 {
	if x != nil && x.FrameCountDown != nil {
		return *x.FrameCountDown
	}
	return 0
}

func (x *Device) GetRelaxedCounter() bool {
	if x != nil && x.RelaxedCounter != nil {
		return *x.RelaxedCounter
	}
	return false
}

func (x *Device) GetKeyWarning() bool {
	if x != nil && x.KeyWarning != nil {
		return *x.KeyWarning
	}
	return false
}

func (x *Device) GetTag() string {
	if x != nil && x.Tag != nil {
		return *x.Tag
	}
	return ""
}

func (x *Device) GetDevNonces() []int32 {
	if x != nil {
		return x.DevNonces
	}
	return nil
}

// UpstreamMessage is a message from one of the devices
type UpstreamMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Eui        string  `protobuf:"bytes,1,opt,name=eui,proto3" json:"eui,omitempty"`
	Timestamp  int64   `protobuf:"varint,2,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Payload    []byte  `protobuf:"bytes,3,opt,name=payload,proto3" json:"payload,omitempty"`
	GatewayEui string  `protobuf:"bytes,4,opt,name=gateway_eui,json=gatewayEui,proto3" json:"gateway_eui,omitempty"`
	Rssi       int32   `protobuf:"varint,5,opt,name=rssi,proto3" json:"rssi,omitempty"`
	Snr        float32 `protobuf:"fixed32,6,opt,name=snr,proto3" json:"snr,omitempty"`
	Frequency  float32 `protobuf:"fixed32,7,opt,name=frequency,proto3" json:"frequency,omitempty"`
	DataRate   string  `protobuf:"bytes,8,opt,name=data_rate,json=dataRate,proto3" json:"data_rate,omitempty"`
	DevAddr    uint32  `protobuf:"varint,9,opt,name=dev_addr,json=devAddr,proto3" json:"dev_addr,omitempty"`
}

func (x *UpstreamMessage) Reset() {
	*x = UpstreamMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lospan_entities_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpstreamMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpstreamMessage) ProtoMessage() {}

func (x *UpstreamMessage) ProtoReflect() protoreflect.Message {
	mi := &file_lospan_entities_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpstreamMessage.ProtoReflect.Descriptor instead.
func (*UpstreamMessage) Descriptor() ([]byte, []int) {
	return file_lospan_entities_proto_rawDescGZIP(), []int{2}
}

func (x *UpstreamMessage) GetEui() string {
	if x != nil {
		return x.Eui
	}
	return ""
}

func (x *UpstreamMessage) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *UpstreamMessage) GetPayload() []byte {
	if x != nil {
		return x.Payload
	}
	return nil
}

func (x *UpstreamMessage) GetGatewayEui() string {
	if x != nil {
		return x.GatewayEui
	}
	return ""
}

func (x *UpstreamMessage) GetRssi() int32 {
	if x != nil {
		return x.Rssi
	}
	return 0
}

func (x *UpstreamMessage) GetSnr() float32 {
	if x != nil {
		return x.Snr
	}
	return 0
}

func (x *UpstreamMessage) GetFrequency() float32 {
	if x != nil {
		return x.Frequency
	}
	return 0
}

func (x *UpstreamMessage) GetDataRate() string {
	if x != nil {
		return x.DataRate
	}
	return ""
}

func (x *UpstreamMessage) GetDevAddr() uint32 {
	if x != nil {
		return x.DevAddr
	}
	return 0
}

// DownstreamMessage is a message that should be or is sent to one of the devices
type DownstreamMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Eui     string `protobuf:"bytes,1,opt,name=eui,proto3" json:"eui,omitempty"`
	Payload []byte `protobuf:"bytes,2,opt,name=payload,proto3" json:"payload,omitempty"`
	Port    int32  `protobuf:"varint,3,opt,name=port,proto3" json:"port,omitempty"`
	Ack     bool   `protobuf:"varint,4,opt,name=ack,proto3" json:"ack,omitempty"`
	Created *int64 `protobuf:"varint,5,opt,name=created,proto3,oneof" json:"created,omitempty"`
	Sent    *int64 `protobuf:"varint,6,opt,name=sent,proto3,oneof" json:"sent,omitempty"`
	AckTime *int64 `protobuf:"varint,7,opt,name=ack_time,json=ackTime,proto3,oneof" json:"ack_time,omitempty"`
}

func (x *DownstreamMessage) Reset() {
	*x = DownstreamMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lospan_entities_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DownstreamMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DownstreamMessage) ProtoMessage() {}

func (x *DownstreamMessage) ProtoReflect() protoreflect.Message {
	mi := &file_lospan_entities_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DownstreamMessage.ProtoReflect.Descriptor instead.
func (*DownstreamMessage) Descriptor() ([]byte, []int) {
	return file_lospan_entities_proto_rawDescGZIP(), []int{3}
}

func (x *DownstreamMessage) GetEui() string {
	if x != nil {
		return x.Eui
	}
	return ""
}

func (x *DownstreamMessage) GetPayload() []byte {
	if x != nil {
		return x.Payload
	}
	return nil
}

func (x *DownstreamMessage) GetPort() int32 {
	if x != nil {
		return x.Port
	}
	return 0
}

func (x *DownstreamMessage) GetAck() bool {
	if x != nil {
		return x.Ack
	}
	return false
}

func (x *DownstreamMessage) GetCreated() int64 {
	if x != nil && x.Created != nil {
		return *x.Created
	}
	return 0
}

func (x *DownstreamMessage) GetSent() int64 {
	if x != nil && x.Sent != nil {
		return *x.Sent
	}
	return 0
}

func (x *DownstreamMessage) GetAckTime() int64 {
	if x != nil && x.AckTime != nil {
		return *x.AckTime
	}
	return 0
}

// Gateway is a LoRaWAN gateway/concentrator.
type Gateway struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Eui       string   `protobuf:"bytes,1,opt,name=eui,proto3" json:"eui,omitempty"`
	Ip        *string  `protobuf:"bytes,2,opt,name=ip,proto3,oneof" json:"ip,omitempty"` // Strictly not optional but used when updating
	StrictIp  *bool    `protobuf:"varint,3,opt,name=strict_ip,json=strictIp,proto3,oneof" json:"strict_ip,omitempty"`
	Latitude  *float32 `protobuf:"fixed32,4,opt,name=latitude,proto3,oneof" json:"latitude,omitempty"`
	Longitude *float32 `protobuf:"fixed32,5,opt,name=longitude,proto3,oneof" json:"longitude,omitempty"`
	Altitude  *float32 `protobuf:"fixed32,6,opt,name=altitude,proto3,oneof" json:"altitude,omitempty"`
}

func (x *Gateway) Reset() {
	*x = Gateway{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lospan_entities_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Gateway) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Gateway) ProtoMessage() {}

func (x *Gateway) ProtoReflect() protoreflect.Message {
	mi := &file_lospan_entities_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Gateway.ProtoReflect.Descriptor instead.
func (*Gateway) Descriptor() ([]byte, []int) {
	return file_lospan_entities_proto_rawDescGZIP(), []int{4}
}

func (x *Gateway) GetEui() string {
	if x != nil {
		return x.Eui
	}
	return ""
}

func (x *Gateway) GetIp() string {
	if x != nil && x.Ip != nil {
		return *x.Ip
	}
	return ""
}

func (x *Gateway) GetStrictIp() bool {
	if x != nil && x.StrictIp != nil {
		return *x.StrictIp
	}
	return false
}

func (x *Gateway) GetLatitude() float32 {
	if x != nil && x.Latitude != nil {
		return *x.Latitude
	}
	return 0
}

func (x *Gateway) GetLongitude() float32 {
	if x != nil && x.Longitude != nil {
		return *x.Longitude
	}
	return 0
}

func (x *Gateway) GetAltitude() float32 {
	if x != nil && x.Altitude != nil {
		return *x.Altitude
	}
	return 0
}

// GatewayMessage is a monitoring message to and from the gateway. This reflects the LoRaWAN gateway UDP
// protocol which again is more or less a 1:1 representation of the radio traffic with acks on top.
type GatewayMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GatewayMessage) Reset() {
	*x = GatewayMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lospan_entities_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GatewayMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GatewayMessage) ProtoMessage() {}

func (x *GatewayMessage) ProtoReflect() protoreflect.Message {
	mi := &file_lospan_entities_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GatewayMessage.ProtoReflect.Descriptor instead.
func (*GatewayMessage) Descriptor() ([]byte, []int) {
	return file_lospan_entities_proto_rawDescGZIP(), []int{5}
}

var File_lospan_entities_proto protoreflect.FileDescriptor

var file_lospan_entities_proto_rawDesc = []byte{
	0x0a, 0x15, 0x6c, 0x6f, 0x73, 0x70, 0x61, 0x6e, 0x2f, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x69, 0x65,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x6c, 0x6f, 0x73, 0x70, 0x61, 0x6e, 0x22,
	0x3e, 0x0a, 0x0b, 0x41, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x10,
	0x0a, 0x03, 0x65, 0x75, 0x69, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x65, 0x75, 0x69,
	0x12, 0x15, 0x0a, 0x03, 0x74, 0x61, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52,
	0x03, 0x74, 0x61, 0x67, 0x88, 0x01, 0x01, 0x42, 0x06, 0x0a, 0x04, 0x5f, 0x74, 0x61, 0x67, 0x22,
	0xc0, 0x05, 0x0a, 0x06, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x12, 0x15, 0x0a, 0x03, 0x65, 0x75,
	0x69, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x03, 0x65, 0x75, 0x69, 0x88, 0x01,
	0x01, 0x12, 0x2c, 0x0a, 0x0f, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x5f, 0x65, 0x75, 0x69, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x48, 0x01, 0x52, 0x0e, 0x61, 0x70,
	0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x45, 0x75, 0x69, 0x88, 0x01, 0x01, 0x12,
	0x2e, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x13,
	0x2e, 0x6c, 0x6f, 0x73, 0x70, 0x61, 0x6e, 0x2e, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x53, 0x74,
	0x61, 0x74, 0x65, 0x48, 0x02, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x88, 0x01, 0x01, 0x12,
	0x1e, 0x0a, 0x08, 0x64, 0x65, 0x76, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x0d, 0x48, 0x03, 0x52, 0x07, 0x64, 0x65, 0x76, 0x41, 0x64, 0x64, 0x72, 0x88, 0x01, 0x01, 0x12,
	0x1c, 0x0a, 0x07, 0x61, 0x70, 0x70, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c,
	0x48, 0x04, 0x52, 0x06, 0x61, 0x70, 0x70, 0x4b, 0x65, 0x79, 0x88, 0x01, 0x01, 0x12, 0x2b, 0x0a,
	0x0f, 0x61, 0x70, 0x70, 0x5f, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x6b, 0x65, 0x79,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x0c, 0x48, 0x05, 0x52, 0x0d, 0x61, 0x70, 0x70, 0x53, 0x65, 0x73,
	0x73, 0x69, 0x6f, 0x6e, 0x4b, 0x65, 0x79, 0x88, 0x01, 0x01, 0x12, 0x33, 0x0a, 0x13, 0x6e, 0x65,
	0x74, 0x77, 0x6f, 0x72, 0x6b, 0x5f, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x6b, 0x65,
	0x79, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0c, 0x48, 0x06, 0x52, 0x11, 0x6e, 0x65, 0x74, 0x77, 0x6f,
	0x72, 0x6b, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x4b, 0x65, 0x79, 0x88, 0x01, 0x01, 0x12,
	0x29, 0x0a, 0x0e, 0x66, 0x72, 0x61, 0x6d, 0x65, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x75,
	0x70, 0x18, 0x08, 0x20, 0x01, 0x28, 0x05, 0x48, 0x07, 0x52, 0x0c, 0x66, 0x72, 0x61, 0x6d, 0x65,
	0x43, 0x6f, 0x75, 0x6e, 0x74, 0x55, 0x70, 0x88, 0x01, 0x01, 0x12, 0x2d, 0x0a, 0x10, 0x66, 0x72,
	0x61, 0x6d, 0x65, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x64, 0x6f, 0x77, 0x6e, 0x18, 0x09,
	0x20, 0x01, 0x28, 0x05, 0x48, 0x08, 0x52, 0x0e, 0x66, 0x72, 0x61, 0x6d, 0x65, 0x43, 0x6f, 0x75,
	0x6e, 0x74, 0x44, 0x6f, 0x77, 0x6e, 0x88, 0x01, 0x01, 0x12, 0x2c, 0x0a, 0x0f, 0x72, 0x65, 0x6c,
	0x61, 0x78, 0x65, 0x64, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x18, 0x0a, 0x20, 0x01,
	0x28, 0x08, 0x48, 0x09, 0x52, 0x0e, 0x72, 0x65, 0x6c, 0x61, 0x78, 0x65, 0x64, 0x43, 0x6f, 0x75,
	0x6e, 0x74, 0x65, 0x72, 0x88, 0x01, 0x01, 0x12, 0x24, 0x0a, 0x0b, 0x6b, 0x65, 0x79, 0x5f, 0x77,
	0x61, 0x72, 0x6e, 0x69, 0x6e, 0x67, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x08, 0x48, 0x0a, 0x52, 0x0a,
	0x6b, 0x65, 0x79, 0x57, 0x61, 0x72, 0x6e, 0x69, 0x6e, 0x67, 0x88, 0x01, 0x01, 0x12, 0x15, 0x0a,
	0x03, 0x74, 0x61, 0x67, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x09, 0x48, 0x0b, 0x52, 0x03, 0x74, 0x61,
	0x67, 0x88, 0x01, 0x01, 0x12, 0x1d, 0x0a, 0x0a, 0x64, 0x65, 0x76, 0x5f, 0x6e, 0x6f, 0x6e, 0x63,
	0x65, 0x73, 0x18, 0x0d, 0x20, 0x03, 0x28, 0x05, 0x52, 0x09, 0x64, 0x65, 0x76, 0x4e, 0x6f, 0x6e,
	0x63, 0x65, 0x73, 0x42, 0x06, 0x0a, 0x04, 0x5f, 0x65, 0x75, 0x69, 0x42, 0x12, 0x0a, 0x10, 0x5f,
	0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x65, 0x75, 0x69, 0x42,
	0x08, 0x0a, 0x06, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x42, 0x0b, 0x0a, 0x09, 0x5f, 0x64, 0x65,
	0x76, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x42, 0x0a, 0x0a, 0x08, 0x5f, 0x61, 0x70, 0x70, 0x5f, 0x6b,
	0x65, 0x79, 0x42, 0x12, 0x0a, 0x10, 0x5f, 0x61, 0x70, 0x70, 0x5f, 0x73, 0x65, 0x73, 0x73, 0x69,
	0x6f, 0x6e, 0x5f, 0x6b, 0x65, 0x79, 0x42, 0x16, 0x0a, 0x14, 0x5f, 0x6e, 0x65, 0x74, 0x77, 0x6f,
	0x72, 0x6b, 0x5f, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x6b, 0x65, 0x79, 0x42, 0x11,
	0x0a, 0x0f, 0x5f, 0x66, 0x72, 0x61, 0x6d, 0x65, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x75,
	0x70, 0x42, 0x13, 0x0a, 0x11, 0x5f, 0x66, 0x72, 0x61, 0x6d, 0x65, 0x5f, 0x63, 0x6f, 0x75, 0x6e,
	0x74, 0x5f, 0x64, 0x6f, 0x77, 0x6e, 0x42, 0x12, 0x0a, 0x10, 0x5f, 0x72, 0x65, 0x6c, 0x61, 0x78,
	0x65, 0x64, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x42, 0x0e, 0x0a, 0x0c, 0x5f, 0x6b,
	0x65, 0x79, 0x5f, 0x77, 0x61, 0x72, 0x6e, 0x69, 0x6e, 0x67, 0x42, 0x06, 0x0a, 0x04, 0x5f, 0x74,
	0x61, 0x67, 0x22, 0xf8, 0x01, 0x0a, 0x0f, 0x55, 0x70, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x65, 0x75, 0x69, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x65, 0x75, 0x69, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x74, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61,
	0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64,
	0x12, 0x1f, 0x0a, 0x0b, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x5f, 0x65, 0x75, 0x69, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x45, 0x75,
	0x69, 0x12, 0x12, 0x0a, 0x04, 0x72, 0x73, 0x73, 0x69, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x04, 0x72, 0x73, 0x73, 0x69, 0x12, 0x10, 0x0a, 0x03, 0x73, 0x6e, 0x72, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x02, 0x52, 0x03, 0x73, 0x6e, 0x72, 0x12, 0x1c, 0x0a, 0x09, 0x66, 0x72, 0x65, 0x71, 0x75,
	0x65, 0x6e, 0x63, 0x79, 0x18, 0x07, 0x20, 0x01, 0x28, 0x02, 0x52, 0x09, 0x66, 0x72, 0x65, 0x71,
	0x75, 0x65, 0x6e, 0x63, 0x79, 0x12, 0x1b, 0x0a, 0x09, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x72, 0x61,
	0x74, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x64, 0x61, 0x74, 0x61, 0x52, 0x61,
	0x74, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x64, 0x65, 0x76, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x18, 0x09,
	0x20, 0x01, 0x28, 0x0d, 0x52, 0x07, 0x64, 0x65, 0x76, 0x41, 0x64, 0x64, 0x72, 0x22, 0xdf, 0x01,
	0x0a, 0x11, 0x44, 0x6f, 0x77, 0x6e, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x65, 0x75, 0x69, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x65, 0x75, 0x69, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x12,
	0x12, 0x0a, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x70,
	0x6f, 0x72, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x61, 0x63, 0x6b, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x03, 0x61, 0x63, 0x6b, 0x12, 0x1d, 0x0a, 0x07, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x48, 0x00, 0x52, 0x07, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x64, 0x88, 0x01, 0x01, 0x12, 0x17, 0x0a, 0x04, 0x73, 0x65, 0x6e, 0x74, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x03, 0x48, 0x01, 0x52, 0x04, 0x73, 0x65, 0x6e, 0x74, 0x88, 0x01, 0x01, 0x12, 0x1e, 0x0a,
	0x08, 0x61, 0x63, 0x6b, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x03, 0x48,
	0x02, 0x52, 0x07, 0x61, 0x63, 0x6b, 0x54, 0x69, 0x6d, 0x65, 0x88, 0x01, 0x01, 0x42, 0x0a, 0x0a,
	0x08, 0x5f, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x73, 0x65,
	0x6e, 0x74, 0x42, 0x0b, 0x0a, 0x09, 0x5f, 0x61, 0x63, 0x6b, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x22,
	0xf4, 0x01, 0x0a, 0x07, 0x47, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x65,
	0x75, 0x69, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x65, 0x75, 0x69, 0x12, 0x13, 0x0a,
	0x02, 0x69, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x02, 0x69, 0x70, 0x88,
	0x01, 0x01, 0x12, 0x20, 0x0a, 0x09, 0x73, 0x74, 0x72, 0x69, 0x63, 0x74, 0x5f, 0x69, 0x70, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x08, 0x48, 0x01, 0x52, 0x08, 0x73, 0x74, 0x72, 0x69, 0x63, 0x74, 0x49,
	0x70, 0x88, 0x01, 0x01, 0x12, 0x1f, 0x0a, 0x08, 0x6c, 0x61, 0x74, 0x69, 0x74, 0x75, 0x64, 0x65,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x02, 0x48, 0x02, 0x52, 0x08, 0x6c, 0x61, 0x74, 0x69, 0x74, 0x75,
	0x64, 0x65, 0x88, 0x01, 0x01, 0x12, 0x21, 0x0a, 0x09, 0x6c, 0x6f, 0x6e, 0x67, 0x69, 0x74, 0x75,
	0x64, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x02, 0x48, 0x03, 0x52, 0x09, 0x6c, 0x6f, 0x6e, 0x67,
	0x69, 0x74, 0x75, 0x64, 0x65, 0x88, 0x01, 0x01, 0x12, 0x1f, 0x0a, 0x08, 0x61, 0x6c, 0x74, 0x69,
	0x74, 0x75, 0x64, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x02, 0x48, 0x04, 0x52, 0x08, 0x61, 0x6c,
	0x74, 0x69, 0x74, 0x75, 0x64, 0x65, 0x88, 0x01, 0x01, 0x42, 0x05, 0x0a, 0x03, 0x5f, 0x69, 0x70,
	0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x73, 0x74, 0x72, 0x69, 0x63, 0x74, 0x5f, 0x69, 0x70, 0x42, 0x0b,
	0x0a, 0x09, 0x5f, 0x6c, 0x61, 0x74, 0x69, 0x74, 0x75, 0x64, 0x65, 0x42, 0x0c, 0x0a, 0x0a, 0x5f,
	0x6c, 0x6f, 0x6e, 0x67, 0x69, 0x74, 0x75, 0x64, 0x65, 0x42, 0x0b, 0x0a, 0x09, 0x5f, 0x61, 0x6c,
	0x74, 0x69, 0x74, 0x75, 0x64, 0x65, 0x22, 0x10, 0x0a, 0x0e, 0x47, 0x61, 0x74, 0x65, 0x77, 0x61,
	0x79, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2a, 0x3f, 0x0a, 0x0b, 0x44, 0x65, 0x76, 0x69,
	0x63, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x0f, 0x0a, 0x0b, 0x55, 0x4e, 0x53, 0x50, 0x45,
	0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x08, 0x0a, 0x04, 0x4f, 0x54, 0x41, 0x41,
	0x10, 0x01, 0x12, 0x07, 0x0a, 0x03, 0x41, 0x42, 0x50, 0x10, 0x02, 0x12, 0x0c, 0x0a, 0x08, 0x44,
	0x49, 0x53, 0x41, 0x42, 0x4c, 0x45, 0x44, 0x10, 0x03, 0x42, 0x0a, 0x5a, 0x08, 0x2e, 0x2f, 0x6c,
	0x6f, 0x73, 0x70, 0x61, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_lospan_entities_proto_rawDescOnce sync.Once
	file_lospan_entities_proto_rawDescData = file_lospan_entities_proto_rawDesc
)

func file_lospan_entities_proto_rawDescGZIP() []byte {
	file_lospan_entities_proto_rawDescOnce.Do(func() {
		file_lospan_entities_proto_rawDescData = protoimpl.X.CompressGZIP(file_lospan_entities_proto_rawDescData)
	})
	return file_lospan_entities_proto_rawDescData
}

var file_lospan_entities_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_lospan_entities_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_lospan_entities_proto_goTypes = []interface{}{
	(DeviceState)(0),          // 0: lospan.DeviceState
	(*Application)(nil),       // 1: lospan.Application
	(*Device)(nil),            // 2: lospan.Device
	(*UpstreamMessage)(nil),   // 3: lospan.UpstreamMessage
	(*DownstreamMessage)(nil), // 4: lospan.DownstreamMessage
	(*Gateway)(nil),           // 5: lospan.Gateway
	(*GatewayMessage)(nil),    // 6: lospan.GatewayMessage
}
var file_lospan_entities_proto_depIdxs = []int32{
	0, // 0: lospan.Device.state:type_name -> lospan.DeviceState
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_lospan_entities_proto_init() }
func file_lospan_entities_proto_init() {
	if File_lospan_entities_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_lospan_entities_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Application); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_lospan_entities_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Device); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_lospan_entities_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpstreamMessage); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_lospan_entities_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DownstreamMessage); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_lospan_entities_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Gateway); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_lospan_entities_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GatewayMessage); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_lospan_entities_proto_msgTypes[0].OneofWrappers = []interface{}{}
	file_lospan_entities_proto_msgTypes[1].OneofWrappers = []interface{}{}
	file_lospan_entities_proto_msgTypes[3].OneofWrappers = []interface{}{}
	file_lospan_entities_proto_msgTypes[4].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_lospan_entities_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_lospan_entities_proto_goTypes,
		DependencyIndexes: file_lospan_entities_proto_depIdxs,
		EnumInfos:         file_lospan_entities_proto_enumTypes,
		MessageInfos:      file_lospan_entities_proto_msgTypes,
	}.Build()
	File_lospan_entities_proto = out.File
	file_lospan_entities_proto_rawDesc = nil
	file_lospan_entities_proto_goTypes = nil
	file_lospan_entities_proto_depIdxs = nil
}
