// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.17.3
// source: config_domain.proto

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

type ConfigDomain struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DomainList *DomainList `protobuf:"bytes,1,opt,name=domainList,proto3" json:"domainList,omitempty"`
}

func (x *ConfigDomain) Reset() {
	*x = ConfigDomain{}
	if protoimpl.UnsafeEnabled {
		mi := &file_config_domain_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConfigDomain) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConfigDomain) ProtoMessage() {}

func (x *ConfigDomain) ProtoReflect() protoreflect.Message {
	mi := &file_config_domain_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConfigDomain.ProtoReflect.Descriptor instead.
func (*ConfigDomain) Descriptor() ([]byte, []int) {
	return file_config_domain_proto_rawDescGZIP(), []int{0}
}

func (x *ConfigDomain) GetDomainList() *DomainList {
	if x != nil {
		return x.DomainList
	}
	return nil
}

type DomainList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uip            uint32          `protobuf:"varint,1,opt,name=uip,proto3" json:"uip,omitempty"`
	Uin            uint32          `protobuf:"varint,2,opt,name=uin,proto3" json:"uin,omitempty"`
	DomainIPList   []*DomainIPList `protobuf:"bytes,3,rep,name=domainIPList,proto3" json:"domainIPList,omitempty"`
	DomainIPv6List []*DomainIPList `protobuf:"bytes,4,rep,name=domainIPv6List,proto3" json:"domainIPv6List,omitempty"`
}

func (x *DomainList) Reset() {
	*x = DomainList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_config_domain_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DomainList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DomainList) ProtoMessage() {}

func (x *DomainList) ProtoReflect() protoreflect.Message {
	mi := &file_config_domain_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DomainList.ProtoReflect.Descriptor instead.
func (*DomainList) Descriptor() ([]byte, []int) {
	return file_config_domain_proto_rawDescGZIP(), []int{1}
}

func (x *DomainList) GetUip() uint32 {
	if x != nil {
		return x.Uip
	}
	return 0
}

func (x *DomainList) GetUin() uint32 {
	if x != nil {
		return x.Uin
	}
	return 0
}

func (x *DomainList) GetDomainIPList() []*DomainIPList {
	if x != nil {
		return x.DomainIPList
	}
	return nil
}

func (x *DomainList) GetDomainIPv6List() []*DomainIPList {
	if x != nil {
		return x.DomainIPv6List
	}
	return nil
}

type DomainIPList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Result uint32      `protobuf:"varint,1,opt,name=result,proto3" json:"result,omitempty"`
	Domain string      `protobuf:"bytes,2,opt,name=domain,proto3" json:"domain,omitempty"`
	Type   uint32      `protobuf:"varint,3,opt,name=type,proto3" json:"type,omitempty"`
	Ttl    uint32      `protobuf:"varint,4,opt,name=ttl,proto3" json:"ttl,omitempty"`
	Ip     []*DomainIP `protobuf:"bytes,5,rep,name=ip,proto3" json:"ip,omitempty"`
}

func (x *DomainIPList) Reset() {
	*x = DomainIPList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_config_domain_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DomainIPList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DomainIPList) ProtoMessage() {}

func (x *DomainIPList) ProtoReflect() protoreflect.Message {
	mi := &file_config_domain_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DomainIPList.ProtoReflect.Descriptor instead.
func (*DomainIPList) Descriptor() ([]byte, []int) {
	return file_config_domain_proto_rawDescGZIP(), []int{2}
}

func (x *DomainIPList) GetResult() uint32 {
	if x != nil {
		return x.Result
	}
	return 0
}

func (x *DomainIPList) GetDomain() string {
	if x != nil {
		return x.Domain
	}
	return ""
}

func (x *DomainIPList) GetType() uint32 {
	if x != nil {
		return x.Type
	}
	return 0
}

func (x *DomainIPList) GetTtl() uint32 {
	if x != nil {
		return x.Ttl
	}
	return 0
}

func (x *DomainIPList) GetIp() []*DomainIP {
	if x != nil {
		return x.Ip
	}
	return nil
}

type DomainIP struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ip   string `protobuf:"bytes,1,opt,name=ip,proto3" json:"ip,omitempty"`
	Port uint32 `protobuf:"varint,2,opt,name=port,proto3" json:"port,omitempty"`
}

func (x *DomainIP) Reset() {
	*x = DomainIP{}
	if protoimpl.UnsafeEnabled {
		mi := &file_config_domain_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DomainIP) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DomainIP) ProtoMessage() {}

func (x *DomainIP) ProtoReflect() protoreflect.Message {
	mi := &file_config_domain_proto_msgTypes[3]
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
	return file_config_domain_proto_rawDescGZIP(), []int{3}
}

func (x *DomainIP) GetIp() string {
	if x != nil {
		return x.Ip
	}
	return ""
}

func (x *DomainIP) GetPort() uint32 {
	if x != nil {
		return x.Port
	}
	return 0
}

var File_config_domain_proto protoreflect.FileDescriptor

var file_config_domain_proto_rawDesc = []byte{
	0x0a, 0x13, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x5f, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x3b, 0x0a, 0x0c, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x44,
	0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x12, 0x2b, 0x0a, 0x0a, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x4c,
	0x69, 0x73, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x44, 0x6f, 0x6d, 0x61,
	0x69, 0x6e, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x0a, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x4c, 0x69,
	0x73, 0x74, 0x22, 0x9a, 0x01, 0x0a, 0x0a, 0x44, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x4c, 0x69, 0x73,
	0x74, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x69, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x03,
	0x75, 0x69, 0x70, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x69, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x03, 0x75, 0x69, 0x6e, 0x12, 0x31, 0x0a, 0x0c, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x49,
	0x50, 0x4c, 0x69, 0x73, 0x74, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x44, 0x6f,
	0x6d, 0x61, 0x69, 0x6e, 0x49, 0x50, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x0c, 0x64, 0x6f, 0x6d, 0x61,
	0x69, 0x6e, 0x49, 0x50, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x35, 0x0a, 0x0e, 0x64, 0x6f, 0x6d, 0x61,
	0x69, 0x6e, 0x49, 0x50, 0x76, 0x36, 0x4c, 0x69, 0x73, 0x74, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x0d, 0x2e, 0x44, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x49, 0x50, 0x4c, 0x69, 0x73, 0x74, 0x52,
	0x0e, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x49, 0x50, 0x76, 0x36, 0x4c, 0x69, 0x73, 0x74, 0x22,
	0x7f, 0x0a, 0x0c, 0x44, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x49, 0x50, 0x4c, 0x69, 0x73, 0x74, 0x12,
	0x16, 0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52,
	0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x64, 0x6f, 0x6d, 0x61, 0x69,
	0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x12,
	0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x04, 0x74,
	0x79, 0x70, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x74, 0x74, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x03, 0x74, 0x74, 0x6c, 0x12, 0x19, 0x0a, 0x02, 0x69, 0x70, 0x18, 0x05, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x09, 0x2e, 0x44, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x49, 0x50, 0x52, 0x02, 0x69, 0x70,
	0x22, 0x2e, 0x0a, 0x08, 0x44, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x49, 0x50, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x70, 0x12, 0x12, 0x0a, 0x04,
	0x70, 0x6f, 0x72, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x04, 0x70, 0x6f, 0x72, 0x74,
	0x42, 0x07, 0x5a, 0x05, 0x2e, 0x2f, 0x3b, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_config_domain_proto_rawDescOnce sync.Once
	file_config_domain_proto_rawDescData = file_config_domain_proto_rawDesc
)

func file_config_domain_proto_rawDescGZIP() []byte {
	file_config_domain_proto_rawDescOnce.Do(func() {
		file_config_domain_proto_rawDescData = protoimpl.X.CompressGZIP(file_config_domain_proto_rawDescData)
	})
	return file_config_domain_proto_rawDescData
}

var file_config_domain_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_config_domain_proto_goTypes = []interface{}{
	(*ConfigDomain)(nil), // 0: ConfigDomain
	(*DomainList)(nil),   // 1: DomainList
	(*DomainIPList)(nil), // 2: DomainIPList
	(*DomainIP)(nil),     // 3: DomainIP
}
var file_config_domain_proto_depIdxs = []int32{
	1, // 0: ConfigDomain.domainList:type_name -> DomainList
	2, // 1: DomainList.domainIPList:type_name -> DomainIPList
	2, // 2: DomainList.domainIPv6List:type_name -> DomainIPList
	3, // 3: DomainIPList.ip:type_name -> DomainIP
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_config_domain_proto_init() }
func file_config_domain_proto_init() {
	if File_config_domain_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_config_domain_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConfigDomain); i {
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
		file_config_domain_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DomainList); i {
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
		file_config_domain_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DomainIPList); i {
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
		file_config_domain_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
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
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_config_domain_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_config_domain_proto_goTypes,
		DependencyIndexes: file_config_domain_proto_depIdxs,
		MessageInfos:      file_config_domain_proto_msgTypes,
	}.Build()
	File_config_domain_proto = out.File
	file_config_domain_proto_rawDesc = nil
	file_config_domain_proto_goTypes = nil
	file_config_domain_proto_depIdxs = nil
}