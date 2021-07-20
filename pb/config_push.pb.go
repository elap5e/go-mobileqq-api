// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.17.3
// source: config_push.proto

package pb

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

// DomainIp
type DomainIP struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DomainIP) Reset() {
	*x = DomainIP{}
	if protoimpl.UnsafeEnabled {
		mi := &file_config_push_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DomainIP) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DomainIP) ProtoMessage() {}

func (x *DomainIP) ProtoReflect() protoreflect.Message {
	mi := &file_config_push_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DomainIP.ProtoReflect.Descriptor instead.
func (*DomainIP) Descriptor() ([]byte, []int) {
	return file_config_push_proto_rawDescGZIP(), []int{0}
}

// NameReqBody
type DomainIP_Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Body *DomainIP_RequestBody `protobuf:"bytes,1,opt,name=body,proto3" json:"body,omitempty"` // "SubCmd_name_Req"
}

func (x *DomainIP_Request) Reset() {
	*x = DomainIP_Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_config_push_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DomainIP_Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DomainIP_Request) ProtoMessage() {}

func (x *DomainIP_Request) ProtoReflect() protoreflect.Message {
	mi := &file_config_push_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DomainIP_Request.ProtoReflect.Descriptor instead.
func (*DomainIP_Request) Descriptor() ([]byte, []int) {
	return file_config_push_proto_rawDescGZIP(), []int{0, 0}
}

func (x *DomainIP_Request) GetBody() *DomainIP_RequestBody {
	if x != nil {
		return x.Body
	}
	return nil
}

// NameRspBody
type DomainIP_Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Body *DomainIP_ResponseBody `protobuf:"bytes,1,opt,name=body,proto3" json:"body,omitempty"` // "SubCmd_name_Rsp"
}

func (x *DomainIP_Response) Reset() {
	*x = DomainIP_Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_config_push_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DomainIP_Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DomainIP_Response) ProtoMessage() {}

func (x *DomainIP_Response) ProtoReflect() protoreflect.Message {
	mi := &file_config_push_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DomainIP_Response.ProtoReflect.Descriptor instead.
func (*DomainIP_Response) Descriptor() ([]byte, []int) {
	return file_config_push_proto_rawDescGZIP(), []int{0, 1}
}

func (x *DomainIP_Response) GetBody() *DomainIP_ResponseBody {
	if x != nil {
		return x.Body
	}
	return nil
}

// ServerList
type DomainIP_ServerList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ip   string `protobuf:"bytes,1,opt,name=ip,proto3" json:"ip,omitempty"`      // "string_IP"
	Port int32  `protobuf:"varint,2,opt,name=port,proto3" json:"port,omitempty"` // "uint32_Port"
}

func (x *DomainIP_ServerList) Reset() {
	*x = DomainIP_ServerList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_config_push_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DomainIP_ServerList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DomainIP_ServerList) ProtoMessage() {}

func (x *DomainIP_ServerList) ProtoReflect() protoreflect.Message {
	mi := &file_config_push_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DomainIP_ServerList.ProtoReflect.Descriptor instead.
func (*DomainIP_ServerList) Descriptor() ([]byte, []int) {
	return file_config_push_proto_rawDescGZIP(), []int{0, 2}
}

func (x *DomainIP_ServerList) GetIp() string {
	if x != nil {
		return x.Ip
	}
	return ""
}

func (x *DomainIP_ServerList) GetPort() int32 {
	if x != nil {
		return x.Port
	}
	return 0
}

// SubCmd_name_Req
type DomainIP_RequestBody struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uip        uint32                 `protobuf:"varint,1,opt,name=uip,proto3" json:"uip,omitempty"`                                // "uint32_uip"
	Uin        int64                  `protobuf:"varint,2,opt,name=uin,proto3" json:"uin,omitempty"`                                // "uint32_qq"
	DomainList []*DomainIP_DomainInfo `protobuf:"bytes,3,rep,name=domain_list,json=domainList,proto3" json:"domain_list,omitempty"` // "dnameInfo"
}

func (x *DomainIP_RequestBody) Reset() {
	*x = DomainIP_RequestBody{}
	if protoimpl.UnsafeEnabled {
		mi := &file_config_push_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DomainIP_RequestBody) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DomainIP_RequestBody) ProtoMessage() {}

func (x *DomainIP_RequestBody) ProtoReflect() protoreflect.Message {
	mi := &file_config_push_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DomainIP_RequestBody.ProtoReflect.Descriptor instead.
func (*DomainIP_RequestBody) Descriptor() ([]byte, []int) {
	return file_config_push_proto_rawDescGZIP(), []int{0, 3}
}

func (x *DomainIP_RequestBody) GetUip() uint32 {
	if x != nil {
		return x.Uip
	}
	return 0
}

func (x *DomainIP_RequestBody) GetUin() int64 {
	if x != nil {
		return x.Uin
	}
	return 0
}

func (x *DomainIP_RequestBody) GetDomainList() []*DomainIP_DomainInfo {
	if x != nil {
		return x.DomainList
	}
	return nil
}

// SubCmd_name_Rsp
type DomainIP_ResponseBody struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uip      uint32             `protobuf:"varint,1,opt,name=uip,proto3" json:"uip,omitempty"`                          // "uint32_uip"
	Uin      int64              `protobuf:"varint,2,opt,name=uin,proto3" json:"uin,omitempty"`                          // "uint32_qq"
	IpList   []*DomainIP_IPInfo `protobuf:"bytes,3,rep,name=ip_list,json=ipList,proto3" json:"ip_list,omitempty"`       // "iplistInfo"
	Ipv6List []*DomainIP_IPInfo `protobuf:"bytes,4,rep,name=ipv6_list,json=ipv6List,proto3" json:"ipv6_list,omitempty"` // "iplistInfoV6"
}

func (x *DomainIP_ResponseBody) Reset() {
	*x = DomainIP_ResponseBody{}
	if protoimpl.UnsafeEnabled {
		mi := &file_config_push_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DomainIP_ResponseBody) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DomainIP_ResponseBody) ProtoMessage() {}

func (x *DomainIP_ResponseBody) ProtoReflect() protoreflect.Message {
	mi := &file_config_push_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DomainIP_ResponseBody.ProtoReflect.Descriptor instead.
func (*DomainIP_ResponseBody) Descriptor() ([]byte, []int) {
	return file_config_push_proto_rawDescGZIP(), []int{0, 4}
}

func (x *DomainIP_ResponseBody) GetUip() uint32 {
	if x != nil {
		return x.Uip
	}
	return 0
}

func (x *DomainIP_ResponseBody) GetUin() int64 {
	if x != nil {
		return x.Uin
	}
	return 0
}

func (x *DomainIP_ResponseBody) GetIpList() []*DomainIP_IPInfo {
	if x != nil {
		return x.IpList
	}
	return nil
}

func (x *DomainIP_ResponseBody) GetIpv6List() []*DomainIP_IPInfo {
	if x != nil {
		return x.Ipv6List
	}
	return nil
}

// dnameInfo
type DomainIP_DomainInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Domain string `protobuf:"bytes,2,opt,name=domain,proto3" json:"domain,omitempty"` // "string_dname"
	Type   uint32 `protobuf:"varint,3,opt,name=type,proto3" json:"type,omitempty"`    // "uint32_type"
}

func (x *DomainIP_DomainInfo) Reset() {
	*x = DomainIP_DomainInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_config_push_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DomainIP_DomainInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DomainIP_DomainInfo) ProtoMessage() {}

func (x *DomainIP_DomainInfo) ProtoReflect() protoreflect.Message {
	mi := &file_config_push_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DomainIP_DomainInfo.ProtoReflect.Descriptor instead.
func (*DomainIP_DomainInfo) Descriptor() ([]byte, []int) {
	return file_config_push_proto_rawDescGZIP(), []int{0, 5}
}

func (x *DomainIP_DomainInfo) GetDomain() string {
	if x != nil {
		return x.Domain
	}
	return ""
}

func (x *DomainIP_DomainInfo) GetType() uint32 {
	if x != nil {
		return x.Type
	}
	return 0
}

// iplistInfo
type DomainIP_IPInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Result int32                  `protobuf:"varint,1,opt,name=result,proto3" json:"result,omitempty"` // "int32_result"
	Domain string                 `protobuf:"bytes,2,opt,name=domain,proto3" json:"domain,omitempty"`  // "string_dname"
	Type   int32                  `protobuf:"varint,3,opt,name=type,proto3" json:"type,omitempty"`     // "uint32_type"
	Ttl    int32                  `protobuf:"varint,4,opt,name=ttl,proto3" json:"ttl,omitempty"`       // "uint32_ttl"
	Items  []*DomainIP_ServerList `protobuf:"bytes,5,rep,name=items,proto3" json:"items,omitempty"`    // "ServerList_data"
}

func (x *DomainIP_IPInfo) Reset() {
	*x = DomainIP_IPInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_config_push_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DomainIP_IPInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DomainIP_IPInfo) ProtoMessage() {}

func (x *DomainIP_IPInfo) ProtoReflect() protoreflect.Message {
	mi := &file_config_push_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DomainIP_IPInfo.ProtoReflect.Descriptor instead.
func (*DomainIP_IPInfo) Descriptor() ([]byte, []int) {
	return file_config_push_proto_rawDescGZIP(), []int{0, 6}
}

func (x *DomainIP_IPInfo) GetResult() int32 {
	if x != nil {
		return x.Result
	}
	return 0
}

func (x *DomainIP_IPInfo) GetDomain() string {
	if x != nil {
		return x.Domain
	}
	return ""
}

func (x *DomainIP_IPInfo) GetType() int32 {
	if x != nil {
		return x.Type
	}
	return 0
}

func (x *DomainIP_IPInfo) GetTtl() int32 {
	if x != nil {
		return x.Ttl
	}
	return 0
}

func (x *DomainIP_IPInfo) GetItems() []*DomainIP_ServerList {
	if x != nil {
		return x.Items
	}
	return nil
}

var File_config_push_proto protoreflect.FileDescriptor

var file_config_push_proto_rawDesc = []byte{
	0x0a, 0x11, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x5f, 0x70, 0x75, 0x73, 0x68, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0xea, 0x04, 0x0a, 0x08, 0x44, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x49, 0x50,
	0x1a, 0x34, 0x0a, 0x07, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x29, 0x0a, 0x04, 0x62,
	0x6f, 0x64, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x44, 0x6f, 0x6d, 0x61,
	0x69, 0x6e, 0x49, 0x50, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x42, 0x6f, 0x64, 0x79,
	0x52, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x1a, 0x36, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x2a, 0x0a, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x16, 0x2e, 0x44, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x49, 0x50, 0x2e, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x42, 0x6f, 0x64, 0x79, 0x52, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x1a, 0x30,
	0x0a, 0x0a, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x70, 0x12, 0x12, 0x0a, 0x04,
	0x70, 0x6f, 0x72, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x70, 0x6f, 0x72, 0x74,
	0x1a, 0x68, 0x0a, 0x0b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x42, 0x6f, 0x64, 0x79, 0x12,
	0x10, 0x0a, 0x03, 0x75, 0x69, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x03, 0x75, 0x69,
	0x70, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x69, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03,
	0x75, 0x69, 0x6e, 0x12, 0x35, 0x0a, 0x0b, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x5f, 0x6c, 0x69,
	0x73, 0x74, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x44, 0x6f, 0x6d, 0x61, 0x69,
	0x6e, 0x49, 0x50, 0x2e, 0x44, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x0a,
	0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x4c, 0x69, 0x73, 0x74, 0x1a, 0x8c, 0x01, 0x0a, 0x0c, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x6f, 0x64, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x75,
	0x69, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x03, 0x75, 0x69, 0x70, 0x12, 0x10, 0x0a,
	0x03, 0x75, 0x69, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x75, 0x69, 0x6e, 0x12,
	0x29, 0x0a, 0x07, 0x69, 0x70, 0x5f, 0x6c, 0x69, 0x73, 0x74, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x10, 0x2e, 0x44, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x49, 0x50, 0x2e, 0x49, 0x50, 0x49, 0x6e,
	0x66, 0x6f, 0x52, 0x06, 0x69, 0x70, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x2d, 0x0a, 0x09, 0x69, 0x70,
	0x76, 0x36, 0x5f, 0x6c, 0x69, 0x73, 0x74, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x10, 0x2e,
	0x44, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x49, 0x50, 0x2e, 0x49, 0x50, 0x49, 0x6e, 0x66, 0x6f, 0x52,
	0x08, 0x69, 0x70, 0x76, 0x36, 0x4c, 0x69, 0x73, 0x74, 0x1a, 0x38, 0x0a, 0x0a, 0x44, 0x6f, 0x6d,
	0x61, 0x69, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x16, 0x0a, 0x06, 0x64, 0x6f, 0x6d, 0x61, 0x69,
	0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x12,
	0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x04, 0x74,
	0x79, 0x70, 0x65, 0x1a, 0x8a, 0x01, 0x0a, 0x06, 0x49, 0x50, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x16,
	0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06,
	0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x12, 0x12,
	0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x74, 0x79,
	0x70, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x74, 0x74, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x03, 0x74, 0x74, 0x6c, 0x12, 0x2a, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x05, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x44, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x49, 0x50, 0x2e, 0x53,
	0x65, 0x72, 0x76, 0x65, 0x72, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73,
	0x42, 0x07, 0x5a, 0x05, 0x2e, 0x2f, 0x3b, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_config_push_proto_rawDescOnce sync.Once
	file_config_push_proto_rawDescData = file_config_push_proto_rawDesc
)

func file_config_push_proto_rawDescGZIP() []byte {
	file_config_push_proto_rawDescOnce.Do(func() {
		file_config_push_proto_rawDescData = protoimpl.X.CompressGZIP(file_config_push_proto_rawDescData)
	})
	return file_config_push_proto_rawDescData
}

var file_config_push_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_config_push_proto_goTypes = []interface{}{
	(*DomainIP)(nil),              // 0: DomainIP
	(*DomainIP_Request)(nil),      // 1: DomainIP.Request
	(*DomainIP_Response)(nil),     // 2: DomainIP.Response
	(*DomainIP_ServerList)(nil),   // 3: DomainIP.ServerList
	(*DomainIP_RequestBody)(nil),  // 4: DomainIP.RequestBody
	(*DomainIP_ResponseBody)(nil), // 5: DomainIP.ResponseBody
	(*DomainIP_DomainInfo)(nil),   // 6: DomainIP.DomainInfo
	(*DomainIP_IPInfo)(nil),       // 7: DomainIP.IPInfo
}
var file_config_push_proto_depIdxs = []int32{
	4, // 0: DomainIP.Request.body:type_name -> DomainIP.RequestBody
	5, // 1: DomainIP.Response.body:type_name -> DomainIP.ResponseBody
	6, // 2: DomainIP.RequestBody.domain_list:type_name -> DomainIP.DomainInfo
	7, // 3: DomainIP.ResponseBody.ip_list:type_name -> DomainIP.IPInfo
	7, // 4: DomainIP.ResponseBody.ipv6_list:type_name -> DomainIP.IPInfo
	3, // 5: DomainIP.IPInfo.items:type_name -> DomainIP.ServerList
	6, // [6:6] is the sub-list for method output_type
	6, // [6:6] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_config_push_proto_init() }
func file_config_push_proto_init() {
	if File_config_push_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_config_push_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DomainIP); i {
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
		file_config_push_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DomainIP_Request); i {
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
		file_config_push_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DomainIP_Response); i {
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
		file_config_push_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DomainIP_ServerList); i {
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
		file_config_push_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DomainIP_RequestBody); i {
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
		file_config_push_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DomainIP_ResponseBody); i {
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
		file_config_push_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DomainIP_DomainInfo); i {
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
		file_config_push_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DomainIP_IPInfo); i {
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
			RawDescriptor: file_config_push_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_config_push_proto_goTypes,
		DependencyIndexes: file_config_push_proto_depIdxs,
		MessageInfos:      file_config_push_proto_msgTypes,
	}.Build()
	File_config_push_proto = out.File
	file_config_push_proto_rawDesc = nil
	file_config_push_proto_goTypes = nil
	file_config_push_proto_depIdxs = nil
}
