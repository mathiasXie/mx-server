// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v6.30.2
// source: proto/tts.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Provider int32

const (
	Provider_MICROSOFT  Provider = 0
	Provider_VOLCENGINE Provider = 1
)

// Enum value maps for Provider.
var (
	Provider_name = map[int32]string{
		0: "MICROSOFT",
		1: "VOLCENGINE",
	}
	Provider_value = map[string]int32{
		"MICROSOFT":  0,
		"VOLCENGINE": 1,
	}
)

func (x Provider) Enum() *Provider {
	p := new(Provider)
	*p = x
	return p
}

func (x Provider) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Provider) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_tts_proto_enumTypes[0].Descriptor()
}

func (Provider) Type() protoreflect.EnumType {
	return &file_proto_tts_proto_enumTypes[0]
}

func (x Provider) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Provider.Descriptor instead.
func (Provider) EnumDescriptor() ([]byte, []int) {
	return file_proto_tts_proto_rawDescGZIP(), []int{0}
}

// 文本转语音请求
type TextToSpeechRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Text          string                 `protobuf:"bytes,1,opt,name=text,proto3" json:"text,omitempty"`                            // 要转换的文本
	Language      string                 `protobuf:"bytes,2,opt,name=language,proto3" json:"language,omitempty"`                    // 语言代码，如 zh-CN, en-US
	VoiceId       string                 `protobuf:"bytes,3,opt,name=voice_id,json=voiceId,proto3" json:"voice_id,omitempty"`       // 声音ID
	Speed         float32                `protobuf:"fixed32,4,opt,name=speed,proto3" json:"speed,omitempty"`                        // 语速，范围 0.5-2.0
	Pitch         float32                `protobuf:"fixed32,5,opt,name=pitch,proto3" json:"pitch,omitempty"`                        // 音调，范围 0.5-2.0
	Provider      Provider               `protobuf:"varint,6,opt,name=provider,proto3,enum=tts.Provider" json:"provider,omitempty"` // 平台，范围
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TextToSpeechRequest) Reset() {
	*x = TextToSpeechRequest{}
	mi := &file_proto_tts_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TextToSpeechRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TextToSpeechRequest) ProtoMessage() {}

func (x *TextToSpeechRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_tts_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TextToSpeechRequest.ProtoReflect.Descriptor instead.
func (*TextToSpeechRequest) Descriptor() ([]byte, []int) {
	return file_proto_tts_proto_rawDescGZIP(), []int{0}
}

func (x *TextToSpeechRequest) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

func (x *TextToSpeechRequest) GetLanguage() string {
	if x != nil {
		return x.Language
	}
	return ""
}

func (x *TextToSpeechRequest) GetVoiceId() string {
	if x != nil {
		return x.VoiceId
	}
	return ""
}

func (x *TextToSpeechRequest) GetSpeed() float32 {
	if x != nil {
		return x.Speed
	}
	return 0
}

func (x *TextToSpeechRequest) GetPitch() float32 {
	if x != nil {
		return x.Pitch
	}
	return 0
}

func (x *TextToSpeechRequest) GetProvider() Provider {
	if x != nil {
		return x.Provider
	}
	return Provider_MICROSOFT
}

// 文本转语音响应
type TextToSpeechResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	AudioData     []byte                 `protobuf:"bytes,1,opt,name=audio_data,json=audioData,proto3" json:"audio_data,omitempty"`     // 音频数据
	Format        string                 `protobuf:"bytes,2,opt,name=format,proto3" json:"format,omitempty"`                            // 音频格式，如 mp3, wav
	SampleRate    int32                  `protobuf:"varint,3,opt,name=sample_rate,json=sampleRate,proto3" json:"sample_rate,omitempty"` // 采样率
	Channels      int32                  `protobuf:"varint,4,opt,name=channels,proto3" json:"channels,omitempty"`                       // 声道数
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TextToSpeechResponse) Reset() {
	*x = TextToSpeechResponse{}
	mi := &file_proto_tts_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TextToSpeechResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TextToSpeechResponse) ProtoMessage() {}

func (x *TextToSpeechResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_tts_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TextToSpeechResponse.ProtoReflect.Descriptor instead.
func (*TextToSpeechResponse) Descriptor() ([]byte, []int) {
	return file_proto_tts_proto_rawDescGZIP(), []int{1}
}

func (x *TextToSpeechResponse) GetAudioData() []byte {
	if x != nil {
		return x.AudioData
	}
	return nil
}

func (x *TextToSpeechResponse) GetFormat() string {
	if x != nil {
		return x.Format
	}
	return ""
}

func (x *TextToSpeechResponse) GetSampleRate() int32 {
	if x != nil {
		return x.SampleRate
	}
	return 0
}

func (x *TextToSpeechResponse) GetChannels() int32 {
	if x != nil {
		return x.Channels
	}
	return 0
}

// 获取语音列表请求
type VoicesListRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Provider      Provider               `protobuf:"varint,1,opt,name=provider,proto3,enum=tts.Provider" json:"provider,omitempty"` // 平台，范围
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *VoicesListRequest) Reset() {
	*x = VoicesListRequest{}
	mi := &file_proto_tts_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *VoicesListRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VoicesListRequest) ProtoMessage() {}

func (x *VoicesListRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_tts_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VoicesListRequest.ProtoReflect.Descriptor instead.
func (*VoicesListRequest) Descriptor() ([]byte, []int) {
	return file_proto_tts_proto_rawDescGZIP(), []int{2}
}

func (x *VoicesListRequest) GetProvider() Provider {
	if x != nil {
		return x.Provider
	}
	return Provider_MICROSOFT
}

// 获取语音列表响应
type VoicesListResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Voices        []*Voice               `protobuf:"bytes,1,rep,name=voices,proto3" json:"voices,omitempty"` // 语音列表
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *VoicesListResponse) Reset() {
	*x = VoicesListResponse{}
	mi := &file_proto_tts_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *VoicesListResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VoicesListResponse) ProtoMessage() {}

func (x *VoicesListResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_tts_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VoicesListResponse.ProtoReflect.Descriptor instead.
func (*VoicesListResponse) Descriptor() ([]byte, []int) {
	return file_proto_tts_proto_rawDescGZIP(), []int{3}
}

func (x *VoicesListResponse) GetVoices() []*Voice {
	if x != nil {
		return x.Voices
	}
	return nil
}

// 语音
type Voice struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	VoiceId       string                 `protobuf:"bytes,1,opt,name=voice_id,json=voiceId,proto3" json:"voice_id,omitempty"`       // 声音ID
	VoiceName     string                 `protobuf:"bytes,2,opt,name=voice_name,json=voiceName,proto3" json:"voice_name,omitempty"` // 声音名称
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Voice) Reset() {
	*x = Voice{}
	mi := &file_proto_tts_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Voice) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Voice) ProtoMessage() {}

func (x *Voice) ProtoReflect() protoreflect.Message {
	mi := &file_proto_tts_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Voice.ProtoReflect.Descriptor instead.
func (*Voice) Descriptor() ([]byte, []int) {
	return file_proto_tts_proto_rawDescGZIP(), []int{4}
}

func (x *Voice) GetVoiceId() string {
	if x != nil {
		return x.VoiceId
	}
	return ""
}

func (x *Voice) GetVoiceName() string {
	if x != nil {
		return x.VoiceName
	}
	return ""
}

var File_proto_tts_proto protoreflect.FileDescriptor

const file_proto_tts_proto_rawDesc = "" +
	"\n" +
	"\x0fproto/tts.proto\x12\x03tts\"\xb7\x01\n" +
	"\x13TextToSpeechRequest\x12\x12\n" +
	"\x04text\x18\x01 \x01(\tR\x04text\x12\x1a\n" +
	"\blanguage\x18\x02 \x01(\tR\blanguage\x12\x19\n" +
	"\bvoice_id\x18\x03 \x01(\tR\avoiceId\x12\x14\n" +
	"\x05speed\x18\x04 \x01(\x02R\x05speed\x12\x14\n" +
	"\x05pitch\x18\x05 \x01(\x02R\x05pitch\x12)\n" +
	"\bprovider\x18\x06 \x01(\x0e2\r.tts.ProviderR\bprovider\"\x8a\x01\n" +
	"\x14TextToSpeechResponse\x12\x1d\n" +
	"\n" +
	"audio_data\x18\x01 \x01(\fR\taudioData\x12\x16\n" +
	"\x06format\x18\x02 \x01(\tR\x06format\x12\x1f\n" +
	"\vsample_rate\x18\x03 \x01(\x05R\n" +
	"sampleRate\x12\x1a\n" +
	"\bchannels\x18\x04 \x01(\x05R\bchannels\">\n" +
	"\x11VoicesListRequest\x12)\n" +
	"\bprovider\x18\x01 \x01(\x0e2\r.tts.ProviderR\bprovider\"8\n" +
	"\x12VoicesListResponse\x12\"\n" +
	"\x06voices\x18\x01 \x03(\v2\n" +
	".tts.VoiceR\x06voices\"A\n" +
	"\x05Voice\x12\x19\n" +
	"\bvoice_id\x18\x01 \x01(\tR\avoiceId\x12\x1d\n" +
	"\n" +
	"voice_name\x18\x02 \x01(\tR\tvoiceName*)\n" +
	"\bProvider\x12\r\n" +
	"\tMICROSOFT\x10\x00\x12\x0e\n" +
	"\n" +
	"VOLCENGINE\x10\x012\x94\x01\n" +
	"\n" +
	"TTSService\x12E\n" +
	"\fTextToSpeech\x12\x18.tts.TextToSpeechRequest\x1a\x19.tts.TextToSpeechResponse\"\x00\x12?\n" +
	"\n" +
	"VoicesList\x12\x16.tts.VoicesListRequest\x1a\x17.tts.VoicesListResponse\"\x00B:Z8github.com/mathiasXie/gin-web/applications/tts-rpc/protob\x06proto3"

var (
	file_proto_tts_proto_rawDescOnce sync.Once
	file_proto_tts_proto_rawDescData []byte
)

func file_proto_tts_proto_rawDescGZIP() []byte {
	file_proto_tts_proto_rawDescOnce.Do(func() {
		file_proto_tts_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_proto_tts_proto_rawDesc), len(file_proto_tts_proto_rawDesc)))
	})
	return file_proto_tts_proto_rawDescData
}

var file_proto_tts_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_proto_tts_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_proto_tts_proto_goTypes = []any{
	(Provider)(0),                // 0: tts.Provider
	(*TextToSpeechRequest)(nil),  // 1: tts.TextToSpeechRequest
	(*TextToSpeechResponse)(nil), // 2: tts.TextToSpeechResponse
	(*VoicesListRequest)(nil),    // 3: tts.VoicesListRequest
	(*VoicesListResponse)(nil),   // 4: tts.VoicesListResponse
	(*Voice)(nil),                // 5: tts.Voice
}
var file_proto_tts_proto_depIdxs = []int32{
	0, // 0: tts.TextToSpeechRequest.provider:type_name -> tts.Provider
	0, // 1: tts.VoicesListRequest.provider:type_name -> tts.Provider
	5, // 2: tts.VoicesListResponse.voices:type_name -> tts.Voice
	1, // 3: tts.TTSService.TextToSpeech:input_type -> tts.TextToSpeechRequest
	3, // 4: tts.TTSService.VoicesList:input_type -> tts.VoicesListRequest
	2, // 5: tts.TTSService.TextToSpeech:output_type -> tts.TextToSpeechResponse
	4, // 6: tts.TTSService.VoicesList:output_type -> tts.VoicesListResponse
	5, // [5:7] is the sub-list for method output_type
	3, // [3:5] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_proto_tts_proto_init() }
func file_proto_tts_proto_init() {
	if File_proto_tts_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_proto_tts_proto_rawDesc), len(file_proto_tts_proto_rawDesc)),
			NumEnums:      1,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_tts_proto_goTypes,
		DependencyIndexes: file_proto_tts_proto_depIdxs,
		EnumInfos:         file_proto_tts_proto_enumTypes,
		MessageInfos:      file_proto_tts_proto_msgTypes,
	}.Build()
	File_proto_tts_proto = out.File
	file_proto_tts_proto_goTypes = nil
	file_proto_tts_proto_depIdxs = nil
}
