// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.12.1
// source: svc/volchestrator_admin.proto

package svc

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

type Volume_Status int32

const (
	Volume_Available   Volume_Status = 0
	Volume_Unavailable Volume_Status = 1
	Volume_Attaching   Volume_Status = 2
	Volume_Attached    Volume_Status = 3
	Volume_Detaching   Volume_Status = 4
)

// Enum value maps for Volume_Status.
var (
	Volume_Status_name = map[int32]string{
		0: "Available",
		1: "Unavailable",
		2: "Attaching",
		3: "Attached",
		4: "Detaching",
	}
	Volume_Status_value = map[string]int32{
		"Available":   0,
		"Unavailable": 1,
		"Attaching":   2,
		"Attached":    3,
		"Detaching":   4,
	}
)

func (x Volume_Status) Enum() *Volume_Status {
	p := new(Volume_Status)
	*p = x
	return p
}

func (x Volume_Status) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Volume_Status) Descriptor() protoreflect.EnumDescriptor {
	return file_svc_volchestrator_admin_proto_enumTypes[0].Descriptor()
}

func (Volume_Status) Type() protoreflect.EnumType {
	return &file_svc_volchestrator_admin_proto_enumTypes[0]
}

func (x Volume_Status) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Volume_Status.Descriptor instead.
func (Volume_Status) EnumDescriptor() ([]byte, []int) {
	return file_svc_volchestrator_admin_proto_rawDescGZIP(), []int{5, 0}
}

type Client struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ClientId   string                 `protobuf:"bytes,1,opt,name=client_id,json=clientId,proto3" json:"client_id,omitempty"`
	Registered *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=registered,proto3" json:"registered,omitempty"`
	LastSeen   *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=last_seen,json=lastSeen,proto3" json:"last_seen,omitempty"`
}

func (x *Client) Reset() {
	*x = Client{}
	if protoimpl.UnsafeEnabled {
		mi := &file_svc_volchestrator_admin_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Client) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Client) ProtoMessage() {}

func (x *Client) ProtoReflect() protoreflect.Message {
	mi := &file_svc_volchestrator_admin_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Client.ProtoReflect.Descriptor instead.
func (*Client) Descriptor() ([]byte, []int) {
	return file_svc_volchestrator_admin_proto_rawDescGZIP(), []int{0}
}

func (x *Client) GetClientId() string {
	if x != nil {
		return x.ClientId
	}
	return ""
}

func (x *Client) GetRegistered() *timestamppb.Timestamp {
	if x != nil {
		return x.Registered
	}
	return nil
}

func (x *Client) GetLastSeen() *timestamppb.Timestamp {
	if x != nil {
		return x.LastSeen
	}
	return nil
}

type GetClientRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ClientId string `protobuf:"bytes,1,opt,name=client_id,json=clientId,proto3" json:"client_id,omitempty"`
}

func (x *GetClientRequest) Reset() {
	*x = GetClientRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_svc_volchestrator_admin_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetClientRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetClientRequest) ProtoMessage() {}

func (x *GetClientRequest) ProtoReflect() protoreflect.Message {
	mi := &file_svc_volchestrator_admin_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetClientRequest.ProtoReflect.Descriptor instead.
func (*GetClientRequest) Descriptor() ([]byte, []int) {
	return file_svc_volchestrator_admin_proto_rawDescGZIP(), []int{1}
}

func (x *GetClientRequest) GetClientId() string {
	if x != nil {
		return x.ClientId
	}
	return ""
}

type GetClientResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Client *Client `protobuf:"bytes,1,opt,name=client,proto3" json:"client,omitempty"`
}

func (x *GetClientResponse) Reset() {
	*x = GetClientResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_svc_volchestrator_admin_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetClientResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetClientResponse) ProtoMessage() {}

func (x *GetClientResponse) ProtoReflect() protoreflect.Message {
	mi := &file_svc_volchestrator_admin_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetClientResponse.ProtoReflect.Descriptor instead.
func (*GetClientResponse) Descriptor() ([]byte, []int) {
	return file_svc_volchestrator_admin_proto_rawDescGZIP(), []int{2}
}

func (x *GetClientResponse) GetClient() *Client {
	if x != nil {
		return x.Client
	}
	return nil
}

type ListClientsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ListClientsRequest) Reset() {
	*x = ListClientsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_svc_volchestrator_admin_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListClientsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListClientsRequest) ProtoMessage() {}

func (x *ListClientsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_svc_volchestrator_admin_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListClientsRequest.ProtoReflect.Descriptor instead.
func (*ListClientsRequest) Descriptor() ([]byte, []int) {
	return file_svc_volchestrator_admin_proto_rawDescGZIP(), []int{3}
}

type ListClientsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Clients []*Client `protobuf:"bytes,1,rep,name=clients,proto3" json:"clients,omitempty"`
}

func (x *ListClientsResponse) Reset() {
	*x = ListClientsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_svc_volchestrator_admin_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListClientsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListClientsResponse) ProtoMessage() {}

func (x *ListClientsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_svc_volchestrator_admin_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListClientsResponse.ProtoReflect.Descriptor instead.
func (*ListClientsResponse) Descriptor() ([]byte, []int) {
	return file_svc_volchestrator_admin_proto_rawDescGZIP(), []int{4}
}

func (x *ListClientsResponse) GetClients() []*Client {
	if x != nil {
		return x.Clients
	}
	return nil
}

type Volume struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	VolumeId string        `protobuf:"bytes,1,opt,name=volume_id,json=volumeId,proto3" json:"volume_id,omitempty"`
	Region   string        `protobuf:"bytes,2,opt,name=region,proto3" json:"region,omitempty"`
	Tag      string        `protobuf:"bytes,3,opt,name=tag,proto3" json:"tag,omitempty"`
	Status   Volume_Status `protobuf:"varint,4,opt,name=status,proto3,enum=volchestrator.Volume_Status" json:"status,omitempty"`
}

func (x *Volume) Reset() {
	*x = Volume{}
	if protoimpl.UnsafeEnabled {
		mi := &file_svc_volchestrator_admin_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Volume) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Volume) ProtoMessage() {}

func (x *Volume) ProtoReflect() protoreflect.Message {
	mi := &file_svc_volchestrator_admin_proto_msgTypes[5]
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
	return file_svc_volchestrator_admin_proto_rawDescGZIP(), []int{5}
}

func (x *Volume) GetVolumeId() string {
	if x != nil {
		return x.VolumeId
	}
	return ""
}

func (x *Volume) GetRegion() string {
	if x != nil {
		return x.Region
	}
	return ""
}

func (x *Volume) GetTag() string {
	if x != nil {
		return x.Tag
	}
	return ""
}

func (x *Volume) GetStatus() Volume_Status {
	if x != nil {
		return x.Status
	}
	return Volume_Available
}

type GetVolumeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	VolumeId string `protobuf:"bytes,1,opt,name=volume_id,json=volumeId,proto3" json:"volume_id,omitempty"`
}

func (x *GetVolumeRequest) Reset() {
	*x = GetVolumeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_svc_volchestrator_admin_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetVolumeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetVolumeRequest) ProtoMessage() {}

func (x *GetVolumeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_svc_volchestrator_admin_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetVolumeRequest.ProtoReflect.Descriptor instead.
func (*GetVolumeRequest) Descriptor() ([]byte, []int) {
	return file_svc_volchestrator_admin_proto_rawDescGZIP(), []int{6}
}

func (x *GetVolumeRequest) GetVolumeId() string {
	if x != nil {
		return x.VolumeId
	}
	return ""
}

type GetVolumeResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Volume *Volume `protobuf:"bytes,1,opt,name=volume,proto3" json:"volume,omitempty"`
}

func (x *GetVolumeResponse) Reset() {
	*x = GetVolumeResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_svc_volchestrator_admin_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetVolumeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetVolumeResponse) ProtoMessage() {}

func (x *GetVolumeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_svc_volchestrator_admin_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetVolumeResponse.ProtoReflect.Descriptor instead.
func (*GetVolumeResponse) Descriptor() ([]byte, []int) {
	return file_svc_volchestrator_admin_proto_rawDescGZIP(), []int{7}
}

func (x *GetVolumeResponse) GetVolume() *Volume {
	if x != nil {
		return x.Volume
	}
	return nil
}

var File_svc_volchestrator_admin_proto protoreflect.FileDescriptor

var file_svc_volchestrator_admin_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x73, 0x76, 0x63, 0x2f, 0x76, 0x6f, 0x6c, 0x63, 0x68, 0x65, 0x73, 0x74, 0x72, 0x61,
	0x74, 0x6f, 0x72, 0x5f, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x0d, 0x76, 0x6f, 0x6c, 0x63, 0x68, 0x65, 0x73, 0x74, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x1a, 0x1f,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x9a, 0x01, 0x0a, 0x06, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x63, 0x6c,
	0x69, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63,
	0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x3a, 0x0a, 0x0a, 0x72, 0x65, 0x67, 0x69, 0x73,
	0x74, 0x65, 0x72, 0x65, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0a, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65,
	0x72, 0x65, 0x64, 0x12, 0x37, 0x0a, 0x09, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x73, 0x65, 0x65, 0x6e,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x52, 0x08, 0x6c, 0x61, 0x73, 0x74, 0x53, 0x65, 0x65, 0x6e, 0x22, 0x2f, 0x0a, 0x10,
	0x47, 0x65, 0x74, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x1b, 0x0a, 0x09, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x22, 0x42, 0x0a,
	0x11, 0x47, 0x65, 0x74, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x2d, 0x0a, 0x06, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x15, 0x2e, 0x76, 0x6f, 0x6c, 0x63, 0x68, 0x65, 0x73, 0x74, 0x72, 0x61, 0x74,
	0x6f, 0x72, 0x2e, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x52, 0x06, 0x63, 0x6c, 0x69, 0x65, 0x6e,
	0x74, 0x22, 0x14, 0x0a, 0x12, 0x4c, 0x69, 0x73, 0x74, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x73,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x46, 0x0a, 0x13, 0x4c, 0x69, 0x73, 0x74, 0x43,
	0x6c, 0x69, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2f,
	0x0a, 0x07, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x15, 0x2e, 0x76, 0x6f, 0x6c, 0x63, 0x68, 0x65, 0x73, 0x74, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x2e,
	0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x52, 0x07, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x73, 0x22,
	0xdb, 0x01, 0x0a, 0x06, 0x56, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x76, 0x6f,
	0x6c, 0x75, 0x6d, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x76,
	0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x67, 0x69, 0x6f,
	0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x12,
	0x10, 0x0a, 0x03, 0x74, 0x61, 0x67, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x74, 0x61,
	0x67, 0x12, 0x34, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x1c, 0x2e, 0x76, 0x6f, 0x6c, 0x63, 0x68, 0x65, 0x73, 0x74, 0x72, 0x61, 0x74, 0x6f,
	0x72, 0x2e, 0x56, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52,
	0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x54, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x12, 0x0d, 0x0a, 0x09, 0x41, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65, 0x10, 0x00,
	0x12, 0x0f, 0x0a, 0x0b, 0x55, 0x6e, 0x61, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65, 0x10,
	0x01, 0x12, 0x0d, 0x0a, 0x09, 0x41, 0x74, 0x74, 0x61, 0x63, 0x68, 0x69, 0x6e, 0x67, 0x10, 0x02,
	0x12, 0x0c, 0x0a, 0x08, 0x41, 0x74, 0x74, 0x61, 0x63, 0x68, 0x65, 0x64, 0x10, 0x03, 0x12, 0x0d,
	0x0a, 0x09, 0x44, 0x65, 0x74, 0x61, 0x63, 0x68, 0x69, 0x6e, 0x67, 0x10, 0x04, 0x22, 0x2f, 0x0a,
	0x10, 0x47, 0x65, 0x74, 0x56, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x1b, 0x0a, 0x09, 0x76, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x76, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x49, 0x64, 0x22, 0x42,
	0x0a, 0x11, 0x47, 0x65, 0x74, 0x56, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x2d, 0x0a, 0x06, 0x76, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x76, 0x6f, 0x6c, 0x63, 0x68, 0x65, 0x73, 0x74, 0x72, 0x61,
	0x74, 0x6f, 0x72, 0x2e, 0x56, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x52, 0x06, 0x76, 0x6f, 0x6c, 0x75,
	0x6d, 0x65, 0x32, 0x90, 0x02, 0x0a, 0x12, 0x56, 0x6f, 0x6c, 0x63, 0x68, 0x65, 0x73, 0x74, 0x72,
	0x61, 0x74, 0x6f, 0x72, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x12, 0x50, 0x0a, 0x09, 0x47, 0x65, 0x74,
	0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x12, 0x1f, 0x2e, 0x76, 0x6f, 0x6c, 0x63, 0x68, 0x65, 0x73,
	0x74, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x20, 0x2e, 0x76, 0x6f, 0x6c, 0x63, 0x68, 0x65,
	0x73, 0x74, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x6c, 0x69, 0x65, 0x6e,
	0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x56, 0x0a, 0x0b, 0x4c,
	0x69, 0x73, 0x74, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x21, 0x2e, 0x76, 0x6f, 0x6c,
	0x63, 0x68, 0x65, 0x73, 0x74, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x43,
	0x6c, 0x69, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x22, 0x2e,
	0x76, 0x6f, 0x6c, 0x63, 0x68, 0x65, 0x73, 0x74, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x4c, 0x69,
	0x73, 0x74, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x00, 0x12, 0x50, 0x0a, 0x09, 0x47, 0x65, 0x74, 0x56, 0x6f, 0x6c, 0x75, 0x6d, 0x65,
	0x12, 0x1f, 0x2e, 0x76, 0x6f, 0x6c, 0x63, 0x68, 0x65, 0x73, 0x74, 0x72, 0x61, 0x74, 0x6f, 0x72,
	0x2e, 0x47, 0x65, 0x74, 0x56, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x20, 0x2e, 0x76, 0x6f, 0x6c, 0x63, 0x68, 0x65, 0x73, 0x74, 0x72, 0x61, 0x74, 0x6f,
	0x72, 0x2e, 0x47, 0x65, 0x74, 0x56, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x27, 0x5a, 0x25, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x70, 0x30, 0x70, 0x72, 0x30, 0x63, 0x6b, 0x35, 0x2f, 0x76, 0x6f, 0x6c,
	0x63, 0x68, 0x65, 0x73, 0x74, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x2f, 0x73, 0x76, 0x63, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_svc_volchestrator_admin_proto_rawDescOnce sync.Once
	file_svc_volchestrator_admin_proto_rawDescData = file_svc_volchestrator_admin_proto_rawDesc
)

func file_svc_volchestrator_admin_proto_rawDescGZIP() []byte {
	file_svc_volchestrator_admin_proto_rawDescOnce.Do(func() {
		file_svc_volchestrator_admin_proto_rawDescData = protoimpl.X.CompressGZIP(file_svc_volchestrator_admin_proto_rawDescData)
	})
	return file_svc_volchestrator_admin_proto_rawDescData
}

var file_svc_volchestrator_admin_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_svc_volchestrator_admin_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_svc_volchestrator_admin_proto_goTypes = []interface{}{
	(Volume_Status)(0),            // 0: volchestrator.Volume.Status
	(*Client)(nil),                // 1: volchestrator.Client
	(*GetClientRequest)(nil),      // 2: volchestrator.GetClientRequest
	(*GetClientResponse)(nil),     // 3: volchestrator.GetClientResponse
	(*ListClientsRequest)(nil),    // 4: volchestrator.ListClientsRequest
	(*ListClientsResponse)(nil),   // 5: volchestrator.ListClientsResponse
	(*Volume)(nil),                // 6: volchestrator.Volume
	(*GetVolumeRequest)(nil),      // 7: volchestrator.GetVolumeRequest
	(*GetVolumeResponse)(nil),     // 8: volchestrator.GetVolumeResponse
	(*timestamppb.Timestamp)(nil), // 9: google.protobuf.Timestamp
}
var file_svc_volchestrator_admin_proto_depIdxs = []int32{
	9, // 0: volchestrator.Client.registered:type_name -> google.protobuf.Timestamp
	9, // 1: volchestrator.Client.last_seen:type_name -> google.protobuf.Timestamp
	1, // 2: volchestrator.GetClientResponse.client:type_name -> volchestrator.Client
	1, // 3: volchestrator.ListClientsResponse.clients:type_name -> volchestrator.Client
	0, // 4: volchestrator.Volume.status:type_name -> volchestrator.Volume.Status
	6, // 5: volchestrator.GetVolumeResponse.volume:type_name -> volchestrator.Volume
	2, // 6: volchestrator.VolchestratorAdmin.GetClient:input_type -> volchestrator.GetClientRequest
	4, // 7: volchestrator.VolchestratorAdmin.ListClients:input_type -> volchestrator.ListClientsRequest
	7, // 8: volchestrator.VolchestratorAdmin.GetVolume:input_type -> volchestrator.GetVolumeRequest
	3, // 9: volchestrator.VolchestratorAdmin.GetClient:output_type -> volchestrator.GetClientResponse
	5, // 10: volchestrator.VolchestratorAdmin.ListClients:output_type -> volchestrator.ListClientsResponse
	8, // 11: volchestrator.VolchestratorAdmin.GetVolume:output_type -> volchestrator.GetVolumeResponse
	9, // [9:12] is the sub-list for method output_type
	6, // [6:9] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_svc_volchestrator_admin_proto_init() }
func file_svc_volchestrator_admin_proto_init() {
	if File_svc_volchestrator_admin_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_svc_volchestrator_admin_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Client); i {
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
		file_svc_volchestrator_admin_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetClientRequest); i {
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
		file_svc_volchestrator_admin_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetClientResponse); i {
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
		file_svc_volchestrator_admin_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListClientsRequest); i {
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
		file_svc_volchestrator_admin_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListClientsResponse); i {
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
		file_svc_volchestrator_admin_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
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
		file_svc_volchestrator_admin_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetVolumeRequest); i {
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
		file_svc_volchestrator_admin_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetVolumeResponse); i {
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
			RawDescriptor: file_svc_volchestrator_admin_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_svc_volchestrator_admin_proto_goTypes,
		DependencyIndexes: file_svc_volchestrator_admin_proto_depIdxs,
		EnumInfos:         file_svc_volchestrator_admin_proto_enumTypes,
		MessageInfos:      file_svc_volchestrator_admin_proto_msgTypes,
	}.Build()
	File_svc_volchestrator_admin_proto = out.File
	file_svc_volchestrator_admin_proto_rawDesc = nil
	file_svc_volchestrator_admin_proto_goTypes = nil
	file_svc_volchestrator_admin_proto_depIdxs = nil
}
