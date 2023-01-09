// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.9
// source: pluginapi/plugin.proto

package pluginapi

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

type EmptyMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *EmptyMessage) Reset() {
	*x = EmptyMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pluginapi_plugin_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EmptyMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EmptyMessage) ProtoMessage() {}

func (x *EmptyMessage) ProtoReflect() protoreflect.Message {
	mi := &file_pluginapi_plugin_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EmptyMessage.ProtoReflect.Descriptor instead.
func (*EmptyMessage) Descriptor() ([]byte, []int) {
	return file_pluginapi_plugin_proto_rawDescGZIP(), []int{0}
}

type NameResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *NameResponse) Reset() {
	*x = NameResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pluginapi_plugin_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NameResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NameResponse) ProtoMessage() {}

func (x *NameResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pluginapi_plugin_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NameResponse.ProtoReflect.Descriptor instead.
func (*NameResponse) Descriptor() ([]byte, []int) {
	return file_pluginapi_plugin_proto_rawDescGZIP(), []int{1}
}

func (x *NameResponse) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type Header struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key   string   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value []string `protobuf:"bytes,2,rep,name=value,proto3" json:"value,omitempty"`
}

func (x *Header) Reset() {
	*x = Header{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pluginapi_plugin_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Header) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Header) ProtoMessage() {}

func (x *Header) ProtoReflect() protoreflect.Message {
	mi := &file_pluginapi_plugin_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Header.ProtoReflect.Descriptor instead.
func (*Header) Descriptor() ([]byte, []int) {
	return file_pluginapi_plugin_proto_rawDescGZIP(), []int{2}
}

func (x *Header) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *Header) GetValue() []string {
	if x != nil {
		return x.Value
	}
	return nil
}

type HandleRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Path   string    `protobuf:"bytes,1,opt,name=path,proto3" json:"path,omitempty"`
	Header []*Header `protobuf:"bytes,2,rep,name=header,proto3" json:"header,omitempty"`
	Body   string    `protobuf:"bytes,3,opt,name=body,proto3" json:"body,omitempty"`
}

func (x *HandleRequest) Reset() {
	*x = HandleRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pluginapi_plugin_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HandleRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HandleRequest) ProtoMessage() {}

func (x *HandleRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pluginapi_plugin_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HandleRequest.ProtoReflect.Descriptor instead.
func (*HandleRequest) Descriptor() ([]byte, []int) {
	return file_pluginapi_plugin_proto_rawDescGZIP(), []int{3}
}

func (x *HandleRequest) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

func (x *HandleRequest) GetHeader() []*Header {
	if x != nil {
		return x.Header
	}
	return nil
}

func (x *HandleRequest) GetBody() string {
	if x != nil {
		return x.Body
	}
	return ""
}

type HandleResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Header []*Header `protobuf:"bytes,1,rep,name=header,proto3" json:"header,omitempty"`
	Body   string    `protobuf:"bytes,2,opt,name=body,proto3" json:"body,omitempty"`
}

func (x *HandleResponse) Reset() {
	*x = HandleResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pluginapi_plugin_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HandleResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HandleResponse) ProtoMessage() {}

func (x *HandleResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pluginapi_plugin_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HandleResponse.ProtoReflect.Descriptor instead.
func (*HandleResponse) Descriptor() ([]byte, []int) {
	return file_pluginapi_plugin_proto_rawDescGZIP(), []int{4}
}

func (x *HandleResponse) GetHeader() []*Header {
	if x != nil {
		return x.Header
	}
	return nil
}

func (x *HandleResponse) GetBody() string {
	if x != nil {
		return x.Body
	}
	return ""
}

var File_pluginapi_plugin_proto protoreflect.FileDescriptor

var file_pluginapi_plugin_proto_rawDesc = []byte{
	0x0a, 0x16, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x6c, 0x75, 0x67,
	0x69, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e,
	0x22, 0x0e, 0x0a, 0x0c, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x22, 0x22, 0x0a, 0x0c, 0x4e, 0x61, 0x6d, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x22, 0x30, 0x0a, 0x06, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12, 0x10,
	0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79,
	0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x5f, 0x0a, 0x0d, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x74, 0x68, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x70, 0x61, 0x74, 0x68, 0x12, 0x26, 0x0a, 0x06, 0x68,
	0x65, 0x61, 0x64, 0x65, 0x72, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x70, 0x6c,
	0x75, 0x67, 0x69, 0x6e, 0x2e, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x52, 0x06, 0x68, 0x65, 0x61,
	0x64, 0x65, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x22, 0x4c, 0x0a, 0x0e, 0x48, 0x61, 0x6e, 0x64, 0x6c,
	0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x26, 0x0a, 0x06, 0x68, 0x65, 0x61,
	0x64, 0x65, 0x72, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x70, 0x6c, 0x75, 0x67,
	0x69, 0x6e, 0x2e, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x52, 0x06, 0x68, 0x65, 0x61, 0x64, 0x65,
	0x72, 0x12, 0x12, 0x0a, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x62, 0x6f, 0x64, 0x79, 0x32, 0xee, 0x01, 0x0a, 0x0e, 0x4f, 0x6e, 0x73, 0x74, 0x61, 0x74,
	0x69, 0x63, 0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x12, 0x34, 0x0a, 0x04, 0x4e, 0x61, 0x6d, 0x65,
	0x12, 0x14, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x14, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e,
	0x4e, 0x61, 0x6d, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x35,
	0x0a, 0x05, 0x53, 0x74, 0x61, 0x72, 0x74, 0x12, 0x14, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e,
	0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x14, 0x2e,
	0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x22, 0x00, 0x12, 0x34, 0x0a, 0x04, 0x53, 0x74, 0x6f, 0x70, 0x12, 0x14, 0x2e,
	0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x1a, 0x14, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x45, 0x6d, 0x70,
	0x74, 0x79, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x00, 0x12, 0x39, 0x0a, 0x06, 0x48,
	0x61, 0x6e, 0x64, 0x6c, 0x65, 0x12, 0x15, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x48,
	0x61, 0x6e, 0x64, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x70,
	0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x25, 0x5a, 0x23, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x74, 0x65, 0x72, 0x73, 0x2f, 0x6f, 0x6e, 0x73, 0x74, 0x61,
	0x74, 0x69, 0x63, 0x2f, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x61, 0x70, 0x69, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pluginapi_plugin_proto_rawDescOnce sync.Once
	file_pluginapi_plugin_proto_rawDescData = file_pluginapi_plugin_proto_rawDesc
)

func file_pluginapi_plugin_proto_rawDescGZIP() []byte {
	file_pluginapi_plugin_proto_rawDescOnce.Do(func() {
		file_pluginapi_plugin_proto_rawDescData = protoimpl.X.CompressGZIP(file_pluginapi_plugin_proto_rawDescData)
	})
	return file_pluginapi_plugin_proto_rawDescData
}

var file_pluginapi_plugin_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_pluginapi_plugin_proto_goTypes = []interface{}{
	(*EmptyMessage)(nil),   // 0: plugin.EmptyMessage
	(*NameResponse)(nil),   // 1: plugin.NameResponse
	(*Header)(nil),         // 2: plugin.Header
	(*HandleRequest)(nil),  // 3: plugin.HandleRequest
	(*HandleResponse)(nil), // 4: plugin.HandleResponse
}
var file_pluginapi_plugin_proto_depIdxs = []int32{
	2, // 0: plugin.HandleRequest.header:type_name -> plugin.Header
	2, // 1: plugin.HandleResponse.header:type_name -> plugin.Header
	0, // 2: plugin.OnstaticPlugin.Name:input_type -> plugin.EmptyMessage
	0, // 3: plugin.OnstaticPlugin.Start:input_type -> plugin.EmptyMessage
	0, // 4: plugin.OnstaticPlugin.Stop:input_type -> plugin.EmptyMessage
	3, // 5: plugin.OnstaticPlugin.Handle:input_type -> plugin.HandleRequest
	1, // 6: plugin.OnstaticPlugin.Name:output_type -> plugin.NameResponse
	0, // 7: plugin.OnstaticPlugin.Start:output_type -> plugin.EmptyMessage
	0, // 8: plugin.OnstaticPlugin.Stop:output_type -> plugin.EmptyMessage
	4, // 9: plugin.OnstaticPlugin.Handle:output_type -> plugin.HandleResponse
	6, // [6:10] is the sub-list for method output_type
	2, // [2:6] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_pluginapi_plugin_proto_init() }
func file_pluginapi_plugin_proto_init() {
	if File_pluginapi_plugin_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pluginapi_plugin_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EmptyMessage); i {
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
		file_pluginapi_plugin_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NameResponse); i {
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
		file_pluginapi_plugin_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Header); i {
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
		file_pluginapi_plugin_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HandleRequest); i {
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
		file_pluginapi_plugin_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HandleResponse); i {
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
			RawDescriptor: file_pluginapi_plugin_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pluginapi_plugin_proto_goTypes,
		DependencyIndexes: file_pluginapi_plugin_proto_depIdxs,
		MessageInfos:      file_pluginapi_plugin_proto_msgTypes,
	}.Build()
	File_pluginapi_plugin_proto = out.File
	file_pluginapi_plugin_proto_rawDesc = nil
	file_pluginapi_plugin_proto_goTypes = nil
	file_pluginapi_plugin_proto_depIdxs = nil
}
