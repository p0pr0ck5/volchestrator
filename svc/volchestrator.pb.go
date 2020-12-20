// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.12.1
// source: svc/volchestrator.proto

package volchestrator

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ClientStatus int32

const (
	ClientStatus_UNKNOWN ClientStatus = 0
	ClientStatus_ALIVE   ClientStatus = 1
	ClientStatus_DEAD    ClientStatus = 2
	ClientStatus_LEFT    ClientStatus = 3
)

// Enum value maps for ClientStatus.
var (
	ClientStatus_name = map[int32]string{
		0: "UNKNOWN",
		1: "ALIVE",
		2: "DEAD",
		3: "LEFT",
	}
	ClientStatus_value = map[string]int32{
		"UNKNOWN": 0,
		"ALIVE":   1,
		"DEAD":    2,
		"LEFT":    3,
	}
)

func (x ClientStatus) Enum() *ClientStatus {
	p := new(ClientStatus)
	*p = x
	return p
}

func (x ClientStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ClientStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_svc_volchestrator_proto_enumTypes[0].Descriptor()
}

func (ClientStatus) Type() protoreflect.EnumType {
	return &file_svc_volchestrator_proto_enumTypes[0]
}

func (x ClientStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ClientStatus.Descriptor instead.
func (ClientStatus) EnumDescriptor() ([]byte, []int) {
	return file_svc_volchestrator_proto_rawDescGZIP(), []int{0}
}

type HeartbeatMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *HeartbeatMessage) Reset() {
	*x = HeartbeatMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_svc_volchestrator_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HeartbeatMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HeartbeatMessage) ProtoMessage() {}

func (x *HeartbeatMessage) ProtoReflect() protoreflect.Message {
	mi := &file_svc_volchestrator_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HeartbeatMessage.ProtoReflect.Descriptor instead.
func (*HeartbeatMessage) Descriptor() ([]byte, []int) {
	return file_svc_volchestrator_proto_rawDescGZIP(), []int{0}
}

func (x *HeartbeatMessage) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type HeartbeatResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *HeartbeatResponse) Reset() {
	*x = HeartbeatResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_svc_volchestrator_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HeartbeatResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HeartbeatResponse) ProtoMessage() {}

func (x *HeartbeatResponse) ProtoReflect() protoreflect.Message {
	mi := &file_svc_volchestrator_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HeartbeatResponse.ProtoReflect.Descriptor instead.
func (*HeartbeatResponse) Descriptor() ([]byte, []int) {
	return file_svc_volchestrator_proto_rawDescGZIP(), []int{1}
}

func (x *HeartbeatResponse) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type ClientInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id           string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	ClientStatus ClientStatus           `protobuf:"varint,2,opt,name=clientStatus,proto3,enum=volchestrator.ClientStatus" json:"clientStatus,omitempty"`
	FirstSeen    *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=firstSeen,proto3" json:"firstSeen,omitempty"`
	LastSeen     *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=lastSeen,proto3" json:"lastSeen,omitempty"`
}

func (x *ClientInfo) Reset() {
	*x = ClientInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_svc_volchestrator_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClientInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClientInfo) ProtoMessage() {}

func (x *ClientInfo) ProtoReflect() protoreflect.Message {
	mi := &file_svc_volchestrator_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClientInfo.ProtoReflect.Descriptor instead.
func (*ClientInfo) Descriptor() ([]byte, []int) {
	return file_svc_volchestrator_proto_rawDescGZIP(), []int{2}
}

func (x *ClientInfo) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *ClientInfo) GetClientStatus() ClientStatus {
	if x != nil {
		return x.ClientStatus
	}
	return ClientStatus_UNKNOWN
}

func (x *ClientInfo) GetFirstSeen() *timestamppb.Timestamp {
	if x != nil {
		return x.FirstSeen
	}
	return nil
}

func (x *ClientInfo) GetLastSeen() *timestamppb.Timestamp {
	if x != nil {
		return x.LastSeen
	}
	return nil
}

type ClientList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Info []*ClientInfo `protobuf:"bytes,1,rep,name=info,proto3" json:"info,omitempty"`
}

func (x *ClientList) Reset() {
	*x = ClientList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_svc_volchestrator_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClientList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClientList) ProtoMessage() {}

func (x *ClientList) ProtoReflect() protoreflect.Message {
	mi := &file_svc_volchestrator_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClientList.ProtoReflect.Descriptor instead.
func (*ClientList) Descriptor() ([]byte, []int) {
	return file_svc_volchestrator_proto_rawDescGZIP(), []int{3}
}

func (x *ClientList) GetInfo() []*ClientInfo {
	if x != nil {
		return x.Info
	}
	return nil
}

type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_svc_volchestrator_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_svc_volchestrator_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_svc_volchestrator_proto_rawDescGZIP(), []int{4}
}

type VolumeID struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *VolumeID) Reset() {
	*x = VolumeID{}
	if protoimpl.UnsafeEnabled {
		mi := &file_svc_volchestrator_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VolumeID) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VolumeID) ProtoMessage() {}

func (x *VolumeID) ProtoReflect() protoreflect.Message {
	mi := &file_svc_volchestrator_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VolumeID.ProtoReflect.Descriptor instead.
func (*VolumeID) Descriptor() ([]byte, []int) {
	return file_svc_volchestrator_proto_rawDescGZIP(), []int{5}
}

func (x *VolumeID) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type Volume struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id               string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Tags             []string `protobuf:"bytes,2,rep,name=tags,proto3" json:"tags,omitempty"`
	AvailabilityZone string   `protobuf:"bytes,3,opt,name=availabilityZone,proto3" json:"availabilityZone,omitempty"`
}

func (x *Volume) Reset() {
	*x = Volume{}
	if protoimpl.UnsafeEnabled {
		mi := &file_svc_volchestrator_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Volume) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Volume) ProtoMessage() {}

func (x *Volume) ProtoReflect() protoreflect.Message {
	mi := &file_svc_volchestrator_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Volume.ProtoReflect.Descriptor instead.
func (*Volume) Descriptor() ([]byte, []int) {
	return file_svc_volchestrator_proto_rawDescGZIP(), []int{6}
}

func (x *Volume) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Volume) GetTags() []string {
	if x != nil {
		return x.Tags
	}
	return nil
}

func (x *Volume) GetAvailabilityZone() string {
	if x != nil {
		return x.AvailabilityZone
	}
	return ""
}

type VolumeList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Volumes []*Volume `protobuf:"bytes,1,rep,name=volumes,proto3" json:"volumes,omitempty"`
}

func (x *VolumeList) Reset() {
	*x = VolumeList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_svc_volchestrator_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VolumeList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VolumeList) ProtoMessage() {}

func (x *VolumeList) ProtoReflect() protoreflect.Message {
	mi := &file_svc_volchestrator_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VolumeList.ProtoReflect.Descriptor instead.
func (*VolumeList) Descriptor() ([]byte, []int) {
	return file_svc_volchestrator_proto_rawDescGZIP(), []int{7}
}

func (x *VolumeList) GetVolumes() []*Volume {
	if x != nil {
		return x.Volumes
	}
	return nil
}

var File_svc_volchestrator_proto protoreflect.FileDescriptor

var file_svc_volchestrator_proto_rawDesc = []byte{
	0x0a, 0x17, 0x73, 0x76, 0x63, 0x2f, 0x76, 0x6f, 0x6c, 0x63, 0x68, 0x65, 0x73, 0x74, 0x72, 0x61,
	0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0d, 0x76, 0x6f, 0x6c, 0x63, 0x68,
	0x65, 0x73, 0x74, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x22, 0x0a, 0x10, 0x48, 0x65, 0x61,
	0x72, 0x74, 0x62, 0x65, 0x61, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x23, 0x0a,
	0x11, 0x48, 0x65, 0x61, 0x72, 0x74, 0x62, 0x65, 0x61, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02,
	0x69, 0x64, 0x22, 0xcf, 0x01, 0x0a, 0x0a, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x6e, 0x66,
	0x6f, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x3f, 0x0a, 0x0c, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x53, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1b, 0x2e, 0x76, 0x6f, 0x6c, 0x63, 0x68, 0x65,
	0x73, 0x74, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x52, 0x0c, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x12, 0x38, 0x0a, 0x09, 0x66, 0x69, 0x72, 0x73, 0x74, 0x53, 0x65, 0x65, 0x6e, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x52, 0x09, 0x66, 0x69, 0x72, 0x73, 0x74, 0x53, 0x65, 0x65, 0x6e, 0x12, 0x36, 0x0a, 0x08,
	0x6c, 0x61, 0x73, 0x74, 0x53, 0x65, 0x65, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x08, 0x6c, 0x61, 0x73, 0x74,
	0x53, 0x65, 0x65, 0x6e, 0x22, 0x3b, 0x0a, 0x0a, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x4c, 0x69,
	0x73, 0x74, 0x12, 0x2d, 0x0a, 0x04, 0x69, 0x6e, 0x66, 0x6f, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x19, 0x2e, 0x76, 0x6f, 0x6c, 0x63, 0x68, 0x65, 0x73, 0x74, 0x72, 0x61, 0x74, 0x6f, 0x72,
	0x2e, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x04, 0x69, 0x6e, 0x66,
	0x6f, 0x22, 0x07, 0x0a, 0x05, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x1a, 0x0a, 0x08, 0x56, 0x6f,
	0x6c, 0x75, 0x6d, 0x65, 0x49, 0x44, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x58, 0x0a, 0x06, 0x56, 0x6f, 0x6c, 0x75, 0x6d, 0x65,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64,
	0x12, 0x12, 0x0a, 0x04, 0x74, 0x61, 0x67, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04,
	0x74, 0x61, 0x67, 0x73, 0x12, 0x2a, 0x0a, 0x10, 0x61, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x69,
	0x6c, 0x69, 0x74, 0x79, 0x5a, 0x6f, 0x6e, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x10,
	0x61, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x5a, 0x6f, 0x6e, 0x65,
	0x22, 0x3d, 0x0a, 0x0a, 0x56, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x2f,
	0x0a, 0x07, 0x76, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x15, 0x2e, 0x76, 0x6f, 0x6c, 0x63, 0x68, 0x65, 0x73, 0x74, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x2e,
	0x56, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x52, 0x07, 0x76, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x73, 0x2a,
	0x3a, 0x0a, 0x0c, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12,
	0x0b, 0x0a, 0x07, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x09, 0x0a, 0x05,
	0x41, 0x4c, 0x49, 0x56, 0x45, 0x10, 0x01, 0x12, 0x08, 0x0a, 0x04, 0x44, 0x45, 0x41, 0x44, 0x10,
	0x02, 0x12, 0x08, 0x0a, 0x04, 0x4c, 0x45, 0x46, 0x54, 0x10, 0x03, 0x32, 0x61, 0x0a, 0x0d, 0x56,
	0x6f, 0x6c, 0x63, 0x68, 0x65, 0x73, 0x74, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x12, 0x50, 0x0a, 0x09,
	0x48, 0x65, 0x61, 0x72, 0x74, 0x62, 0x65, 0x61, 0x74, 0x12, 0x1f, 0x2e, 0x76, 0x6f, 0x6c, 0x63,
	0x68, 0x65, 0x73, 0x74, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x48, 0x65, 0x61, 0x72, 0x74, 0x62,
	0x65, 0x61, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x20, 0x2e, 0x76, 0x6f, 0x6c,
	0x63, 0x68, 0x65, 0x73, 0x74, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x48, 0x65, 0x61, 0x72, 0x74,
	0x62, 0x65, 0x61, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x32, 0x95,
	0x03, 0x0a, 0x12, 0x56, 0x6f, 0x6c, 0x63, 0x68, 0x65, 0x73, 0x74, 0x72, 0x61, 0x74, 0x6f, 0x72,
	0x41, 0x64, 0x6d, 0x69, 0x6e, 0x12, 0x40, 0x0a, 0x0b, 0x4c, 0x69, 0x73, 0x74, 0x43, 0x6c, 0x69,
	0x65, 0x6e, 0x74, 0x73, 0x12, 0x14, 0x2e, 0x76, 0x6f, 0x6c, 0x63, 0x68, 0x65, 0x73, 0x74, 0x72,
	0x61, 0x74, 0x6f, 0x72, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x19, 0x2e, 0x76, 0x6f, 0x6c,
	0x63, 0x68, 0x65, 0x73, 0x74, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x43, 0x6c, 0x69, 0x65, 0x6e,
	0x74, 0x4c, 0x69, 0x73, 0x74, 0x22, 0x00, 0x12, 0x3d, 0x0a, 0x09, 0x47, 0x65, 0x74, 0x56, 0x6f,
	0x6c, 0x75, 0x6d, 0x65, 0x12, 0x17, 0x2e, 0x76, 0x6f, 0x6c, 0x63, 0x68, 0x65, 0x73, 0x74, 0x72,
	0x61, 0x74, 0x6f, 0x72, 0x2e, 0x56, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x49, 0x44, 0x1a, 0x15, 0x2e,
	0x76, 0x6f, 0x6c, 0x63, 0x68, 0x65, 0x73, 0x74, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x56, 0x6f,
	0x6c, 0x75, 0x6d, 0x65, 0x22, 0x00, 0x12, 0x40, 0x0a, 0x0b, 0x4c, 0x69, 0x73, 0x74, 0x56, 0x6f,
	0x6c, 0x75, 0x6d, 0x65, 0x73, 0x12, 0x14, 0x2e, 0x76, 0x6f, 0x6c, 0x63, 0x68, 0x65, 0x73, 0x74,
	0x72, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x19, 0x2e, 0x76, 0x6f,
	0x6c, 0x63, 0x68, 0x65, 0x73, 0x74, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x56, 0x6f, 0x6c, 0x75,
	0x6d, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x22, 0x00, 0x12, 0x3b, 0x0a, 0x09, 0x41, 0x64, 0x64, 0x56,
	0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x12, 0x15, 0x2e, 0x76, 0x6f, 0x6c, 0x63, 0x68, 0x65, 0x73, 0x74,
	0x72, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x56, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x1a, 0x15, 0x2e, 0x76,
	0x6f, 0x6c, 0x63, 0x68, 0x65, 0x73, 0x74, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x56, 0x6f, 0x6c,
	0x75, 0x6d, 0x65, 0x22, 0x00, 0x12, 0x3e, 0x0a, 0x0c, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x56,
	0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x12, 0x15, 0x2e, 0x76, 0x6f, 0x6c, 0x63, 0x68, 0x65, 0x73, 0x74,
	0x72, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x56, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x1a, 0x15, 0x2e, 0x76,
	0x6f, 0x6c, 0x63, 0x68, 0x65, 0x73, 0x74, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x56, 0x6f, 0x6c,
	0x75, 0x6d, 0x65, 0x22, 0x00, 0x12, 0x3f, 0x0a, 0x0c, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x56,
	0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x12, 0x17, 0x2e, 0x76, 0x6f, 0x6c, 0x63, 0x68, 0x65, 0x73, 0x74,
	0x72, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x56, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x49, 0x44, 0x1a, 0x14,
	0x2e, 0x76, 0x6f, 0x6c, 0x63, 0x68, 0x65, 0x73, 0x74, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_svc_volchestrator_proto_rawDescOnce sync.Once
	file_svc_volchestrator_proto_rawDescData = file_svc_volchestrator_proto_rawDesc
)

func file_svc_volchestrator_proto_rawDescGZIP() []byte {
	file_svc_volchestrator_proto_rawDescOnce.Do(func() {
		file_svc_volchestrator_proto_rawDescData = protoimpl.X.CompressGZIP(file_svc_volchestrator_proto_rawDescData)
	})
	return file_svc_volchestrator_proto_rawDescData
}

var file_svc_volchestrator_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_svc_volchestrator_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_svc_volchestrator_proto_goTypes = []interface{}{
	(ClientStatus)(0),             // 0: volchestrator.ClientStatus
	(*HeartbeatMessage)(nil),      // 1: volchestrator.HeartbeatMessage
	(*HeartbeatResponse)(nil),     // 2: volchestrator.HeartbeatResponse
	(*ClientInfo)(nil),            // 3: volchestrator.ClientInfo
	(*ClientList)(nil),            // 4: volchestrator.ClientList
	(*Empty)(nil),                 // 5: volchestrator.Empty
	(*VolumeID)(nil),              // 6: volchestrator.VolumeID
	(*Volume)(nil),                // 7: volchestrator.Volume
	(*VolumeList)(nil),            // 8: volchestrator.VolumeList
	(*timestamppb.Timestamp)(nil), // 9: google.protobuf.Timestamp
}
var file_svc_volchestrator_proto_depIdxs = []int32{
	0,  // 0: volchestrator.ClientInfo.clientStatus:type_name -> volchestrator.ClientStatus
	9,  // 1: volchestrator.ClientInfo.firstSeen:type_name -> google.protobuf.Timestamp
	9,  // 2: volchestrator.ClientInfo.lastSeen:type_name -> google.protobuf.Timestamp
	3,  // 3: volchestrator.ClientList.info:type_name -> volchestrator.ClientInfo
	7,  // 4: volchestrator.VolumeList.volumes:type_name -> volchestrator.Volume
	1,  // 5: volchestrator.Volchestrator.Heartbeat:input_type -> volchestrator.HeartbeatMessage
	5,  // 6: volchestrator.VolchestratorAdmin.ListClients:input_type -> volchestrator.Empty
	6,  // 7: volchestrator.VolchestratorAdmin.GetVolume:input_type -> volchestrator.VolumeID
	5,  // 8: volchestrator.VolchestratorAdmin.ListVolumes:input_type -> volchestrator.Empty
	7,  // 9: volchestrator.VolchestratorAdmin.AddVolume:input_type -> volchestrator.Volume
	7,  // 10: volchestrator.VolchestratorAdmin.UpdateVolume:input_type -> volchestrator.Volume
	6,  // 11: volchestrator.VolchestratorAdmin.DeleteVolume:input_type -> volchestrator.VolumeID
	2,  // 12: volchestrator.Volchestrator.Heartbeat:output_type -> volchestrator.HeartbeatResponse
	4,  // 13: volchestrator.VolchestratorAdmin.ListClients:output_type -> volchestrator.ClientList
	7,  // 14: volchestrator.VolchestratorAdmin.GetVolume:output_type -> volchestrator.Volume
	8,  // 15: volchestrator.VolchestratorAdmin.ListVolumes:output_type -> volchestrator.VolumeList
	7,  // 16: volchestrator.VolchestratorAdmin.AddVolume:output_type -> volchestrator.Volume
	7,  // 17: volchestrator.VolchestratorAdmin.UpdateVolume:output_type -> volchestrator.Volume
	5,  // 18: volchestrator.VolchestratorAdmin.DeleteVolume:output_type -> volchestrator.Empty
	12, // [12:19] is the sub-list for method output_type
	5,  // [5:12] is the sub-list for method input_type
	5,  // [5:5] is the sub-list for extension type_name
	5,  // [5:5] is the sub-list for extension extendee
	0,  // [0:5] is the sub-list for field type_name
}

func init() { file_svc_volchestrator_proto_init() }
func file_svc_volchestrator_proto_init() {
	if File_svc_volchestrator_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_svc_volchestrator_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HeartbeatMessage); i {
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
		file_svc_volchestrator_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HeartbeatResponse); i {
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
		file_svc_volchestrator_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClientInfo); i {
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
		file_svc_volchestrator_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClientList); i {
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
		file_svc_volchestrator_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Empty); i {
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
		file_svc_volchestrator_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VolumeID); i {
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
		file_svc_volchestrator_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Volume); i {
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
		file_svc_volchestrator_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VolumeList); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_svc_volchestrator_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   2,
		},
		GoTypes:           file_svc_volchestrator_proto_goTypes,
		DependencyIndexes: file_svc_volchestrator_proto_depIdxs,
		EnumInfos:         file_svc_volchestrator_proto_enumTypes,
		MessageInfos:      file_svc_volchestrator_proto_msgTypes,
	}.Build()
	File_svc_volchestrator_proto = out.File
	file_svc_volchestrator_proto_rawDesc = nil
	file_svc_volchestrator_proto_goTypes = nil
	file_svc_volchestrator_proto_depIdxs = nil
}
