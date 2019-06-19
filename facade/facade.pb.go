// Code generated by protoc-gen-go. DO NOT EDIT.
// source: facade.proto

package facade

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Mode int32

const (
	Mode_GRID  Mode = 0
	Mode_DRAFT Mode = 16
)

var Mode_name = map[int32]string{
	0:  "GRID",
	16: "DRAFT",
}
var Mode_value = map[string]int32{
	"GRID":  0,
	"DRAFT": 16,
}

func (x Mode) String() string {
	return proto.EnumName(Mode_name, int32(x))
}
func (Mode) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_facade_49e7cc99e9a6ca2d, []int{0}
}

type Response struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Response) Reset()         { *m = Response{} }
func (m *Response) String() string { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()    {}
func (*Response) Descriptor() ([]byte, []int) {
	return fileDescriptor_facade_49e7cc99e9a6ca2d, []int{0}
}
func (m *Response) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Response.Unmarshal(m, b)
}
func (m *Response) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Response.Marshal(b, m, deterministic)
}
func (dst *Response) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Response.Merge(dst, src)
}
func (m *Response) XXX_Size() int {
	return xxx_messageInfo_Response.Size(m)
}
func (m *Response) XXX_DiscardUnknown() {
	xxx_messageInfo_Response.DiscardUnknown(m)
}

var xxx_messageInfo_Response proto.InternalMessageInfo

type Text struct {
	Text                 []byte   `protobuf:"bytes,1,opt,name=text,proto3" json:"text,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Text) Reset()         { *m = Text{} }
func (m *Text) String() string { return proto.CompactTextString(m) }
func (*Text) ProtoMessage()    {}
func (*Text) Descriptor() ([]byte, []int) {
	return fileDescriptor_facade_49e7cc99e9a6ca2d, []int{1}
}
func (m *Text) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Text.Unmarshal(m, b)
}
func (m *Text) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Text.Marshal(b, m, deterministic)
}
func (dst *Text) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Text.Merge(dst, src)
}
func (m *Text) XXX_Size() int {
	return xxx_messageInfo_Text.Size(m)
}
func (m *Text) XXX_DiscardUnknown() {
	xxx_messageInfo_Text.DiscardUnknown(m)
}

var xxx_messageInfo_Text proto.InternalMessageInfo

func (m *Text) GetText() []byte {
	if m != nil {
		return m.Text
	}
	return nil
}

type Config struct {
	CheckDebug           bool        `protobuf:"varint,1,opt,name=checkDebug,proto3" json:"checkDebug,omitempty"`
	Debug                bool        `protobuf:"varint,2,opt,name=Debug,proto3" json:"Debug,omitempty"`
	CheckMode            bool        `protobuf:"varint,3,opt,name=checkMode,proto3" json:"checkMode,omitempty"`
	Mode                 Mode        `protobuf:"varint,4,opt,name=Mode,proto3,enum=facade.Mode" json:"Mode,omitempty"`
	Grid                 *GridConfig `protobuf:"bytes,5,opt,name=Grid,proto3" json:"Grid,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *Config) Reset()         { *m = Config{} }
func (m *Config) String() string { return proto.CompactTextString(m) }
func (*Config) ProtoMessage()    {}
func (*Config) Descriptor() ([]byte, []int) {
	return fileDescriptor_facade_49e7cc99e9a6ca2d, []int{2}
}
func (m *Config) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Config.Unmarshal(m, b)
}
func (m *Config) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Config.Marshal(b, m, deterministic)
}
func (dst *Config) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Config.Merge(dst, src)
}
func (m *Config) XXX_Size() int {
	return xxx_messageInfo_Config.Size(m)
}
func (m *Config) XXX_DiscardUnknown() {
	xxx_messageInfo_Config.DiscardUnknown(m)
}

var xxx_messageInfo_Config proto.InternalMessageInfo

func (m *Config) GetCheckDebug() bool {
	if m != nil {
		return m.CheckDebug
	}
	return false
}

func (m *Config) GetDebug() bool {
	if m != nil {
		return m.Debug
	}
	return false
}

func (m *Config) GetCheckMode() bool {
	if m != nil {
		return m.CheckMode
	}
	return false
}

func (m *Config) GetMode() Mode {
	if m != nil {
		return m.Mode
	}
	return Mode_GRID
}

func (m *Config) GetGrid() *GridConfig {
	if m != nil {
		return m.Grid
	}
	return nil
}

type GridConfig struct {
	CheckWidth           bool     `protobuf:"varint,1,opt,name=checkWidth,proto3" json:"checkWidth,omitempty"`
	Width                uint64   `protobuf:"varint,2,opt,name=Width,proto3" json:"Width,omitempty"`
	CheckHeight          bool     `protobuf:"varint,3,opt,name=checkHeight,proto3" json:"checkHeight,omitempty"`
	Height               uint64   `protobuf:"varint,4,opt,name=Height,proto3" json:"Height,omitempty"`
	CheckDownward        bool     `protobuf:"varint,5,opt,name=checkDownward,proto3" json:"checkDownward,omitempty"`
	Downward             bool     `protobuf:"varint,6,opt,name=Downward,proto3" json:"Downward,omitempty"`
	CheckSpeed           bool     `protobuf:"varint,7,opt,name=checkSpeed,proto3" json:"checkSpeed,omitempty"`
	Speed                float64  `protobuf:"fixed64,8,opt,name=Speed,proto3" json:"Speed,omitempty"`
	CheckBuffer          bool     `protobuf:"varint,9,opt,name=checkBuffer,proto3" json:"checkBuffer,omitempty"`
	Buffer               uint64   `protobuf:"varint,10,opt,name=Buffer,proto3" json:"Buffer,omitempty"`
	CheckTerminal        bool     `protobuf:"varint,11,opt,name=checkTerminal,proto3" json:"checkTerminal,omitempty"`
	Terminal             bool     `protobuf:"varint,12,opt,name=Terminal,proto3" json:"Terminal,omitempty"`
	CheckVert            bool     `protobuf:"varint,13,opt,name=checkVert,proto3" json:"checkVert,omitempty"`
	Vert                 string   `protobuf:"bytes,14,opt,name=Vert,proto3" json:"Vert,omitempty"`
	CheckFrag            bool     `protobuf:"varint,15,opt,name=checkFrag,proto3" json:"checkFrag,omitempty"`
	Frag                 string   `protobuf:"bytes,16,opt,name=Frag,proto3" json:"Frag,omitempty"`
	CheckFill            bool     `protobuf:"varint,17,opt,name=checkFill,proto3" json:"checkFill,omitempty"`
	Fill                 string   `protobuf:"bytes,18,opt,name=Fill,proto3" json:"Fill,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GridConfig) Reset()         { *m = GridConfig{} }
func (m *GridConfig) String() string { return proto.CompactTextString(m) }
func (*GridConfig) ProtoMessage()    {}
func (*GridConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_facade_49e7cc99e9a6ca2d, []int{3}
}
func (m *GridConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GridConfig.Unmarshal(m, b)
}
func (m *GridConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GridConfig.Marshal(b, m, deterministic)
}
func (dst *GridConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GridConfig.Merge(dst, src)
}
func (m *GridConfig) XXX_Size() int {
	return xxx_messageInfo_GridConfig.Size(m)
}
func (m *GridConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_GridConfig.DiscardUnknown(m)
}

var xxx_messageInfo_GridConfig proto.InternalMessageInfo

func (m *GridConfig) GetCheckWidth() bool {
	if m != nil {
		return m.CheckWidth
	}
	return false
}

func (m *GridConfig) GetWidth() uint64 {
	if m != nil {
		return m.Width
	}
	return 0
}

func (m *GridConfig) GetCheckHeight() bool {
	if m != nil {
		return m.CheckHeight
	}
	return false
}

func (m *GridConfig) GetHeight() uint64 {
	if m != nil {
		return m.Height
	}
	return 0
}

func (m *GridConfig) GetCheckDownward() bool {
	if m != nil {
		return m.CheckDownward
	}
	return false
}

func (m *GridConfig) GetDownward() bool {
	if m != nil {
		return m.Downward
	}
	return false
}

func (m *GridConfig) GetCheckSpeed() bool {
	if m != nil {
		return m.CheckSpeed
	}
	return false
}

func (m *GridConfig) GetSpeed() float64 {
	if m != nil {
		return m.Speed
	}
	return 0
}

func (m *GridConfig) GetCheckBuffer() bool {
	if m != nil {
		return m.CheckBuffer
	}
	return false
}

func (m *GridConfig) GetBuffer() uint64 {
	if m != nil {
		return m.Buffer
	}
	return 0
}

func (m *GridConfig) GetCheckTerminal() bool {
	if m != nil {
		return m.CheckTerminal
	}
	return false
}

func (m *GridConfig) GetTerminal() bool {
	if m != nil {
		return m.Terminal
	}
	return false
}

func (m *GridConfig) GetCheckVert() bool {
	if m != nil {
		return m.CheckVert
	}
	return false
}

func (m *GridConfig) GetVert() string {
	if m != nil {
		return m.Vert
	}
	return ""
}

func (m *GridConfig) GetCheckFrag() bool {
	if m != nil {
		return m.CheckFrag
	}
	return false
}

func (m *GridConfig) GetFrag() string {
	if m != nil {
		return m.Frag
	}
	return ""
}

func (m *GridConfig) GetCheckFill() bool {
	if m != nil {
		return m.CheckFill
	}
	return false
}

func (m *GridConfig) GetFill() string {
	if m != nil {
		return m.Fill
	}
	return ""
}

func init() {
	proto.RegisterType((*Response)(nil), "facade.Response")
	proto.RegisterType((*Text)(nil), "facade.Text")
	proto.RegisterType((*Config)(nil), "facade.Config")
	proto.RegisterType((*GridConfig)(nil), "facade.GridConfig")
	proto.RegisterEnum("facade.Mode", Mode_name, Mode_value)
}

func init() { proto.RegisterFile("facade.proto", fileDescriptor_facade_49e7cc99e9a6ca2d) }

var fileDescriptor_facade_49e7cc99e9a6ca2d = []byte{
	// 437 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x93, 0xe1, 0x8a, 0xd3, 0x40,
	0x10, 0xc7, 0x5d, 0xdd, 0xc6, 0x64, 0xda, 0xab, 0x71, 0x10, 0x59, 0xaa, 0x48, 0x28, 0x22, 0x41,
	0xf0, 0x3e, 0xd4, 0x27, 0x50, 0x4b, 0x4f, 0x3f, 0x08, 0xb2, 0x16, 0xfd, 0x9c, 0x6b, 0xa6, 0xed,
	0x62, 0x6c, 0x4a, 0x2e, 0xe5, 0xee, 0x09, 0x7c, 0x14, 0x9f, 0x53, 0x76, 0x26, 0xd9, 0x4b, 0xf1,
	0x5b, 0xe7, 0xf7, 0x9f, 0x3f, 0xfc, 0x68, 0x66, 0x61, 0xb2, 0x2d, 0x36, 0x45, 0x49, 0x97, 0xc7,
	0xa6, 0x6e, 0x6b, 0x8c, 0x64, 0x9a, 0x03, 0xc4, 0x96, 0x6e, 0x8e, 0xf5, 0xe1, 0x86, 0xe6, 0x33,
	0xd0, 0x6b, 0xba, 0x6b, 0x11, 0x41, 0xb7, 0x74, 0xd7, 0x1a, 0x95, 0xa9, 0x7c, 0x62, 0xf9, 0xf7,
	0xfc, 0xaf, 0x82, 0xe8, 0x53, 0x7d, 0xd8, 0xba, 0x1d, 0xbe, 0x02, 0xd8, 0xec, 0x69, 0xf3, 0x6b,
	0x49, 0xd7, 0xa7, 0x1d, 0x2f, 0xc5, 0x76, 0x40, 0xf0, 0x19, 0x8c, 0x24, 0x7a, 0xc8, 0x91, 0x0c,
	0xf8, 0x12, 0x12, 0xde, 0xf9, 0x5a, 0x97, 0x64, 0x1e, 0x71, 0x72, 0x0f, 0x30, 0x03, 0xcd, 0x81,
	0xce, 0x54, 0x3e, 0x5d, 0x4c, 0x2e, 0x3b, 0x57, 0xcf, 0x2c, 0x27, 0xf8, 0x06, 0xf4, 0x55, 0xe3,
	0x4a, 0x33, 0xca, 0x54, 0x3e, 0x5e, 0x60, 0xbf, 0xe1, 0x99, 0x78, 0x59, 0xce, 0xe7, 0x7f, 0x34,
	0xc0, 0x3d, 0x0c, 0xb2, 0x3f, 0x5d, 0xd9, 0xee, 0xcf, 0x64, 0x99, 0x78, 0x59, 0x89, 0xbc, 0xac,
	0xb6, 0x32, 0x60, 0x06, 0x63, 0xde, 0xf9, 0x4c, 0x6e, 0xb7, 0x6f, 0x3b, 0xdd, 0x21, 0xc2, 0xe7,
	0x10, 0x75, 0xa1, 0xe6, 0x62, 0x37, 0xe1, 0x6b, 0xb8, 0x90, 0xbf, 0xa2, 0xbe, 0x3d, 0xdc, 0x16,
	0x8d, 0xf8, 0xc6, 0xf6, 0x1c, 0xe2, 0x0c, 0xe2, 0xb0, 0x10, 0xf1, 0x42, 0x98, 0x83, 0xf1, 0xf7,
	0x23, 0x51, 0x69, 0x1e, 0x0f, 0x8c, 0x99, 0x78, 0x63, 0x89, 0xe2, 0x4c, 0xe5, 0xca, 0xca, 0x10,
	0x8c, 0x3f, 0x9e, 0xb6, 0x5b, 0x6a, 0x4c, 0x32, 0x30, 0x16, 0xe4, 0x8d, 0xbb, 0x10, 0xc4, 0xb8,
	0xe3, 0xbd, 0xf1, 0x9a, 0x9a, 0xdf, 0xee, 0x50, 0x54, 0x66, 0x3c, 0x30, 0xee, 0xa1, 0x37, 0x0e,
	0x0b, 0x13, 0x31, 0x0e, 0x59, 0xff, 0x69, 0x7f, 0x50, 0xd3, 0x9a, 0x8b, 0xc1, 0xa7, 0xf5, 0xc0,
	0x5f, 0x13, 0x07, 0xd3, 0x4c, 0xe5, 0x89, 0xe5, 0xdf, 0xa1, 0xb1, 0x6a, 0x8a, 0x9d, 0x79, 0x32,
	0x68, 0x78, 0xe0, 0x1b, 0x1c, 0xa4, 0xd2, 0x60, 0x16, 0x1a, 0xae, 0xaa, 0xcc, 0xd3, 0x61, 0xc3,
	0x55, 0x15, 0x37, 0x7c, 0x80, 0x5d, 0xc3, 0x55, 0xd5, 0xdb, 0x17, 0x72, 0x52, 0x18, 0x83, 0xbe,
	0xb2, 0x5f, 0x96, 0xe9, 0x03, 0x4c, 0x60, 0xb4, 0xb4, 0x1f, 0x56, 0xeb, 0x34, 0x5d, 0x14, 0x10,
	0xad, 0xf8, 0x80, 0xf0, 0x1d, 0x24, 0x72, 0x2a, 0xa7, 0x86, 0x70, 0xda, 0x9f, 0x95, 0xa0, 0x59,
	0xda, 0xcf, 0xfd, 0x1b, 0xc1, 0x1c, 0xf4, 0x37, 0x77, 0x24, 0x0c, 0x27, 0xea, 0x5f, 0xcc, 0xff,
	0x7b, 0xb9, 0xba, 0x8e, 0xf8, 0xa1, 0xbd, 0xff, 0x17, 0x00, 0x00, 0xff, 0xff, 0x14, 0x99, 0xdc,
	0xc3, 0x78, 0x03, 0x00, 0x00,
}
