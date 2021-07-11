// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.17.3
// source: message_hummer.proto

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

type MessageElementInfoServiceType33 struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Index  uint32 `protobuf:"varint,1,opt,name=index,proto3" json:"index,omitempty"`
	Text   []byte `protobuf:"bytes,2,opt,name=text,proto3" json:"text,omitempty"`
	Compat []byte `protobuf:"bytes,3,opt,name=compat,proto3" json:"compat,omitempty"`
	Buffer []byte `protobuf:"bytes,4,opt,name=buffer,proto3" json:"buffer,omitempty"`
}

func (x *MessageElementInfoServiceType33) Reset() {
	*x = MessageElementInfoServiceType33{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_hummer_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MessageElementInfoServiceType33) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageElementInfoServiceType33) ProtoMessage() {}

func (x *MessageElementInfoServiceType33) ProtoReflect() protoreflect.Message {
	mi := &file_message_hummer_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MessageElementInfoServiceType33.ProtoReflect.Descriptor instead.
func (*MessageElementInfoServiceType33) Descriptor() ([]byte, []int) {
	return file_message_hummer_proto_rawDescGZIP(), []int{0}
}

func (x *MessageElementInfoServiceType33) GetIndex() uint32 {
	if x != nil {
		return x.Index
	}
	return 0
}

func (x *MessageElementInfoServiceType33) GetText() []byte {
	if x != nil {
		return x.Text
	}
	return nil
}

func (x *MessageElementInfoServiceType33) GetCompat() []byte {
	if x != nil {
		return x.Compat
	}
	return nil
}

func (x *MessageElementInfoServiceType33) GetBuffer() []byte {
	if x != nil {
		return x.Buffer
	}
	return nil
}

type MessageElementInfoServiceType37 struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PackId      []byte `protobuf:"bytes,1,opt,name=packId,proto3" json:"packId,omitempty"`
	StickerId   []byte `protobuf:"bytes,2,opt,name=stickerId,proto3" json:"stickerId,omitempty"`
	QsId        uint32 `protobuf:"varint,3,opt,name=qsId,proto3" json:"qsId,omitempty"`
	SourceType  uint32 `protobuf:"varint,4,opt,name=sourceType,proto3" json:"sourceType,omitempty"`
	StickerType uint32 `protobuf:"varint,5,opt,name=stickerType,proto3" json:"stickerType,omitempty"`
	ResultId    []byte `protobuf:"bytes,6,opt,name=resultId,proto3" json:"resultId,omitempty"`
	Text        []byte `protobuf:"bytes,7,opt,name=text,proto3" json:"text,omitempty"`
	SurpriseId  []byte `protobuf:"bytes,8,opt,name=surpriseId,proto3" json:"surpriseId,omitempty"`
	RandomType  uint32 `protobuf:"varint,9,opt,name=randomType,proto3" json:"randomType,omitempty"`
}

func (x *MessageElementInfoServiceType37) Reset() {
	*x = MessageElementInfoServiceType37{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_hummer_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MessageElementInfoServiceType37) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageElementInfoServiceType37) ProtoMessage() {}

func (x *MessageElementInfoServiceType37) ProtoReflect() protoreflect.Message {
	mi := &file_message_hummer_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MessageElementInfoServiceType37.ProtoReflect.Descriptor instead.
func (*MessageElementInfoServiceType37) Descriptor() ([]byte, []int) {
	return file_message_hummer_proto_rawDescGZIP(), []int{1}
}

func (x *MessageElementInfoServiceType37) GetPackId() []byte {
	if x != nil {
		return x.PackId
	}
	return nil
}

func (x *MessageElementInfoServiceType37) GetStickerId() []byte {
	if x != nil {
		return x.StickerId
	}
	return nil
}

func (x *MessageElementInfoServiceType37) GetQsId() uint32 {
	if x != nil {
		return x.QsId
	}
	return 0
}

func (x *MessageElementInfoServiceType37) GetSourceType() uint32 {
	if x != nil {
		return x.SourceType
	}
	return 0
}

func (x *MessageElementInfoServiceType37) GetStickerType() uint32 {
	if x != nil {
		return x.StickerType
	}
	return 0
}

func (x *MessageElementInfoServiceType37) GetResultId() []byte {
	if x != nil {
		return x.ResultId
	}
	return nil
}

func (x *MessageElementInfoServiceType37) GetText() []byte {
	if x != nil {
		return x.Text
	}
	return nil
}

func (x *MessageElementInfoServiceType37) GetSurpriseId() []byte {
	if x != nil {
		return x.SurpriseId
	}
	return nil
}

func (x *MessageElementInfoServiceType37) GetRandomType() uint32 {
	if x != nil {
		return x.RandomType
	}
	return 0
}

var File_message_hummer_proto protoreflect.FileDescriptor

var file_message_hummer_proto_rawDesc = []byte{
	0x0a, 0x14, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x68, 0x75, 0x6d, 0x6d, 0x65, 0x72,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x7b, 0x0a, 0x1f, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x45, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x54, 0x79, 0x70, 0x65, 0x33, 0x33, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x6e, 0x64,
	0x65, 0x78, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x12,
	0x12, 0x0a, 0x04, 0x74, 0x65, 0x78, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x74,
	0x65, 0x78, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x63, 0x6f, 0x6d, 0x70, 0x61, 0x74, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x06, 0x63, 0x6f, 0x6d, 0x70, 0x61, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x62,
	0x75, 0x66, 0x66, 0x65, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x62, 0x75, 0x66,
	0x66, 0x65, 0x72, 0x22, 0x9d, 0x02, 0x0a, 0x1f, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x45,
	0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x54, 0x79, 0x70, 0x65, 0x33, 0x37, 0x12, 0x16, 0x0a, 0x06, 0x70, 0x61, 0x63, 0x6b, 0x49,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x70, 0x61, 0x63, 0x6b, 0x49, 0x64, 0x12,
	0x1c, 0x0a, 0x09, 0x73, 0x74, 0x69, 0x63, 0x6b, 0x65, 0x72, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x09, 0x73, 0x74, 0x69, 0x63, 0x6b, 0x65, 0x72, 0x49, 0x64, 0x12, 0x12, 0x0a,
	0x04, 0x71, 0x73, 0x49, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x04, 0x71, 0x73, 0x49,
	0x64, 0x12, 0x1e, 0x0a, 0x0a, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x54, 0x79, 0x70, 0x65, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0a, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x54, 0x79, 0x70,
	0x65, 0x12, 0x20, 0x0a, 0x0b, 0x73, 0x74, 0x69, 0x63, 0x6b, 0x65, 0x72, 0x54, 0x79, 0x70, 0x65,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0b, 0x73, 0x74, 0x69, 0x63, 0x6b, 0x65, 0x72, 0x54,
	0x79, 0x70, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x49, 0x64, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x49, 0x64, 0x12,
	0x12, 0x0a, 0x04, 0x74, 0x65, 0x78, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x74,
	0x65, 0x78, 0x74, 0x12, 0x1e, 0x0a, 0x0a, 0x73, 0x75, 0x72, 0x70, 0x72, 0x69, 0x73, 0x65, 0x49,
	0x64, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0a, 0x73, 0x75, 0x72, 0x70, 0x72, 0x69, 0x73,
	0x65, 0x49, 0x64, 0x12, 0x1e, 0x0a, 0x0a, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x54, 0x79, 0x70,
	0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0a, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x54,
	0x79, 0x70, 0x65, 0x42, 0x07, 0x5a, 0x05, 0x2e, 0x2f, 0x3b, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_message_hummer_proto_rawDescOnce sync.Once
	file_message_hummer_proto_rawDescData = file_message_hummer_proto_rawDesc
)

func file_message_hummer_proto_rawDescGZIP() []byte {
	file_message_hummer_proto_rawDescOnce.Do(func() {
		file_message_hummer_proto_rawDescData = protoimpl.X.CompressGZIP(file_message_hummer_proto_rawDescData)
	})
	return file_message_hummer_proto_rawDescData
}

var file_message_hummer_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_message_hummer_proto_goTypes = []interface{}{
	(*MessageElementInfoServiceType33)(nil), // 0: MessageElementInfoServiceType33
	(*MessageElementInfoServiceType37)(nil), // 1: MessageElementInfoServiceType37
}
var file_message_hummer_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_message_hummer_proto_init() }
func file_message_hummer_proto_init() {
	if File_message_hummer_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_message_hummer_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MessageElementInfoServiceType33); i {
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
		file_message_hummer_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MessageElementInfoServiceType37); i {
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
			RawDescriptor: file_message_hummer_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_message_hummer_proto_goTypes,
		DependencyIndexes: file_message_hummer_proto_depIdxs,
		MessageInfos:      file_message_hummer_proto_msgTypes,
	}.Build()
	File_message_hummer_proto = out.File
	file_message_hummer_proto_rawDesc = nil
	file_message_hummer_proto_goTypes = nil
	file_message_hummer_proto_depIdxs = nil
}