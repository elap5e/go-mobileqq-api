// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.17.3
// source: online.proto

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

type OnlinePushMessageTransfer struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FromUin         uint64 `protobuf:"varint,1,opt,name=fromUin,proto3" json:"fromUin,omitempty"`
	ToUin           uint64 `protobuf:"varint,2,opt,name=toUin,proto3" json:"toUin,omitempty"`
	MessageType     uint32 `protobuf:"varint,3,opt,name=messageType,proto3" json:"messageType,omitempty"`
	MessageSubtype  uint32 `protobuf:"varint,4,opt,name=messageSubtype,proto3" json:"messageSubtype,omitempty"`
	MessageSeq      uint32 `protobuf:"varint,5,opt,name=messageSeq,proto3" json:"messageSeq,omitempty"`
	MessageUid      uint64 `protobuf:"varint,6,opt,name=messageUid,proto3" json:"messageUid,omitempty"`
	MessageTime     uint32 `protobuf:"varint,7,opt,name=messageTime,proto3" json:"messageTime,omitempty"`
	RealMessageTime uint32 `protobuf:"varint,8,opt,name=realMessageTime,proto3" json:"realMessageTime,omitempty"`
	NickName        string `protobuf:"bytes,9,opt,name=nickName,proto3" json:"nickName,omitempty"`
	MessageData     []byte `protobuf:"bytes,10,opt,name=messageData,proto3" json:"messageData,omitempty"`
	ServerIp        uint32 `protobuf:"varint,11,opt,name=serverIp,proto3" json:"serverIp,omitempty"`
}

func (x *OnlinePushMessageTransfer) Reset() {
	*x = OnlinePushMessageTransfer{}
	if protoimpl.UnsafeEnabled {
		mi := &file_online_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OnlinePushMessageTransfer) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OnlinePushMessageTransfer) ProtoMessage() {}

func (x *OnlinePushMessageTransfer) ProtoReflect() protoreflect.Message {
	mi := &file_online_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OnlinePushMessageTransfer.ProtoReflect.Descriptor instead.
func (*OnlinePushMessageTransfer) Descriptor() ([]byte, []int) {
	return file_online_proto_rawDescGZIP(), []int{0}
}

func (x *OnlinePushMessageTransfer) GetFromUin() uint64 {
	if x != nil {
		return x.FromUin
	}
	return 0
}

func (x *OnlinePushMessageTransfer) GetToUin() uint64 {
	if x != nil {
		return x.ToUin
	}
	return 0
}

func (x *OnlinePushMessageTransfer) GetMessageType() uint32 {
	if x != nil {
		return x.MessageType
	}
	return 0
}

func (x *OnlinePushMessageTransfer) GetMessageSubtype() uint32 {
	if x != nil {
		return x.MessageSubtype
	}
	return 0
}

func (x *OnlinePushMessageTransfer) GetMessageSeq() uint32 {
	if x != nil {
		return x.MessageSeq
	}
	return 0
}

func (x *OnlinePushMessageTransfer) GetMessageUid() uint64 {
	if x != nil {
		return x.MessageUid
	}
	return 0
}

func (x *OnlinePushMessageTransfer) GetMessageTime() uint32 {
	if x != nil {
		return x.MessageTime
	}
	return 0
}

func (x *OnlinePushMessageTransfer) GetRealMessageTime() uint32 {
	if x != nil {
		return x.RealMessageTime
	}
	return 0
}

func (x *OnlinePushMessageTransfer) GetNickName() string {
	if x != nil {
		return x.NickName
	}
	return ""
}

func (x *OnlinePushMessageTransfer) GetMessageData() []byte {
	if x != nil {
		return x.MessageData
	}
	return nil
}

func (x *OnlinePushMessageTransfer) GetServerIp() uint32 {
	if x != nil {
		return x.ServerIp
	}
	return 0
}

var File_online_proto protoreflect.FileDescriptor

var file_online_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x6f, 0x6e, 0x6c, 0x69, 0x6e, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xfb,
	0x02, 0x0a, 0x19, 0x4f, 0x6e, 0x6c, 0x69, 0x6e, 0x65, 0x50, 0x75, 0x73, 0x68, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x12, 0x18, 0x0a, 0x07,
	0x66, 0x72, 0x6f, 0x6d, 0x55, 0x69, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x07, 0x66,
	0x72, 0x6f, 0x6d, 0x55, 0x69, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x55, 0x69, 0x6e, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x74, 0x6f, 0x55, 0x69, 0x6e, 0x12, 0x20, 0x0a, 0x0b,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x54, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x0b, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x26,
	0x0a, 0x0e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x53, 0x75, 0x62, 0x74, 0x79, 0x70, 0x65,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x53,
	0x75, 0x62, 0x74, 0x79, 0x70, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x53, 0x65, 0x71, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0a, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x53, 0x65, 0x71, 0x12, 0x1e, 0x0a, 0x0a, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x55, 0x69, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0a, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x55, 0x69, 0x64, 0x12, 0x20, 0x0a, 0x0b, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0b, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x28, 0x0a, 0x0f, 0x72, 0x65, 0x61, 0x6c,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x0f, 0x72, 0x65, 0x61, 0x6c, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x54, 0x69,
	0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x6e, 0x69, 0x63, 0x6b, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x09,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6e, 0x69, 0x63, 0x6b, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x20,
	0x0a, 0x0b, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x44, 0x61, 0x74, 0x61, 0x18, 0x0a, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x0b, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x44, 0x61, 0x74, 0x61,
	0x12, 0x1a, 0x0a, 0x08, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x49, 0x70, 0x18, 0x0b, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x08, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x49, 0x70, 0x42, 0x07, 0x5a, 0x05,
	0x2e, 0x2f, 0x3b, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_online_proto_rawDescOnce sync.Once
	file_online_proto_rawDescData = file_online_proto_rawDesc
)

func file_online_proto_rawDescGZIP() []byte {
	file_online_proto_rawDescOnce.Do(func() {
		file_online_proto_rawDescData = protoimpl.X.CompressGZIP(file_online_proto_rawDescData)
	})
	return file_online_proto_rawDescData
}

var file_online_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_online_proto_goTypes = []interface{}{
	(*OnlinePushMessageTransfer)(nil), // 0: OnlinePushMessageTransfer
}
var file_online_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_online_proto_init() }
func file_online_proto_init() {
	if File_online_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_online_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OnlinePushMessageTransfer); i {
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
			RawDescriptor: file_online_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_online_proto_goTypes,
		DependencyIndexes: file_online_proto_depIdxs,
		MessageInfos:      file_online_proto_msgTypes,
	}.Build()
	File_online_proto = out.File
	file_online_proto_rawDesc = nil
	file_online_proto_goTypes = nil
	file_online_proto_depIdxs = nil
}
