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
	return fileDescriptor_facade_5db679803abb8e98, []int{0}
}

type Empty struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Empty) Reset()         { *m = Empty{} }
func (m *Empty) String() string { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()    {}
func (*Empty) Descriptor() ([]byte, []int) {
	return fileDescriptor_facade_5db679803abb8e98, []int{0}
}
func (m *Empty) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Empty.Unmarshal(m, b)
}
func (m *Empty) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Empty.Marshal(b, m, deterministic)
}
func (dst *Empty) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Empty.Merge(dst, src)
}
func (m *Empty) XXX_Size() int {
	return xxx_messageInfo_Empty.Size(m)
}
func (m *Empty) XXX_DiscardUnknown() {
	xxx_messageInfo_Empty.DiscardUnknown(m)
}

var xxx_messageInfo_Empty proto.InternalMessageInfo

type Text struct {
	Text                 []byte   `protobuf:"bytes,1,opt,name=Text,proto3" json:"Text,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Text) Reset()         { *m = Text{} }
func (m *Text) String() string { return proto.CompactTextString(m) }
func (*Text) ProtoMessage()    {}
func (*Text) Descriptor() ([]byte, []int) {
	return fileDescriptor_facade_5db679803abb8e98, []int{1}
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
	SetDebug             bool          `protobuf:"varint,1,opt,name=SetDebug,proto3" json:"SetDebug,omitempty"`
	Debug                bool          `protobuf:"varint,2,opt,name=Debug,proto3" json:"Debug,omitempty"`
	SetMode              bool          `protobuf:"varint,3,opt,name=SetMode,proto3" json:"SetMode,omitempty"`
	Mode                 Mode          `protobuf:"varint,4,opt,name=Mode,proto3,enum=facade.Mode" json:"Mode,omitempty"`
	Font                 *FontConfig   `protobuf:"bytes,5,opt,name=Font,proto3" json:"Font,omitempty"`
	Camera               *CameraConfig `protobuf:"bytes,6,opt,name=Camera,proto3" json:"Camera,omitempty"`
	Mask                 *MaskConfig   `protobuf:"bytes,7,opt,name=Mask,proto3" json:"Mask,omitempty"`
	Grid                 *GridConfig   `protobuf:"bytes,8,opt,name=Grid,proto3" json:"Grid,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *Config) Reset()         { *m = Config{} }
func (m *Config) String() string { return proto.CompactTextString(m) }
func (*Config) ProtoMessage()    {}
func (*Config) Descriptor() ([]byte, []int) {
	return fileDescriptor_facade_5db679803abb8e98, []int{2}
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

func (m *Config) GetSetDebug() bool {
	if m != nil {
		return m.SetDebug
	}
	return false
}

func (m *Config) GetDebug() bool {
	if m != nil {
		return m.Debug
	}
	return false
}

func (m *Config) GetSetMode() bool {
	if m != nil {
		return m.SetMode
	}
	return false
}

func (m *Config) GetMode() Mode {
	if m != nil {
		return m.Mode
	}
	return Mode_GRID
}

func (m *Config) GetFont() *FontConfig {
	if m != nil {
		return m.Font
	}
	return nil
}

func (m *Config) GetCamera() *CameraConfig {
	if m != nil {
		return m.Camera
	}
	return nil
}

func (m *Config) GetMask() *MaskConfig {
	if m != nil {
		return m.Mask
	}
	return nil
}

func (m *Config) GetGrid() *GridConfig {
	if m != nil {
		return m.Grid
	}
	return nil
}

type FontConfig struct {
	SetName              bool     `protobuf:"varint,1,opt,name=SetName,proto3" json:"SetName,omitempty"`
	Name                 string   `protobuf:"bytes,2,opt,name=Name,proto3" json:"Name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FontConfig) Reset()         { *m = FontConfig{} }
func (m *FontConfig) String() string { return proto.CompactTextString(m) }
func (*FontConfig) ProtoMessage()    {}
func (*FontConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_facade_5db679803abb8e98, []int{3}
}
func (m *FontConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FontConfig.Unmarshal(m, b)
}
func (m *FontConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FontConfig.Marshal(b, m, deterministic)
}
func (dst *FontConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FontConfig.Merge(dst, src)
}
func (m *FontConfig) XXX_Size() int {
	return xxx_messageInfo_FontConfig.Size(m)
}
func (m *FontConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_FontConfig.DiscardUnknown(m)
}

var xxx_messageInfo_FontConfig proto.InternalMessageInfo

func (m *FontConfig) GetSetName() bool {
	if m != nil {
		return m.SetName
	}
	return false
}

func (m *FontConfig) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type CameraConfig struct {
	SetZoom              bool     `protobuf:"varint,1,opt,name=SetZoom,proto3" json:"SetZoom,omitempty"`
	Zoom                 float64  `protobuf:"fixed64,2,opt,name=Zoom,proto3" json:"Zoom,omitempty"`
	SetIsometric         bool     `protobuf:"varint,3,opt,name=SetIsometric,proto3" json:"SetIsometric,omitempty"`
	Isometric            bool     `protobuf:"varint,4,opt,name=Isometric,proto3" json:"Isometric,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CameraConfig) Reset()         { *m = CameraConfig{} }
func (m *CameraConfig) String() string { return proto.CompactTextString(m) }
func (*CameraConfig) ProtoMessage()    {}
func (*CameraConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_facade_5db679803abb8e98, []int{4}
}
func (m *CameraConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CameraConfig.Unmarshal(m, b)
}
func (m *CameraConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CameraConfig.Marshal(b, m, deterministic)
}
func (dst *CameraConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CameraConfig.Merge(dst, src)
}
func (m *CameraConfig) XXX_Size() int {
	return xxx_messageInfo_CameraConfig.Size(m)
}
func (m *CameraConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_CameraConfig.DiscardUnknown(m)
}

var xxx_messageInfo_CameraConfig proto.InternalMessageInfo

func (m *CameraConfig) GetSetZoom() bool {
	if m != nil {
		return m.SetZoom
	}
	return false
}

func (m *CameraConfig) GetZoom() float64 {
	if m != nil {
		return m.Zoom
	}
	return 0
}

func (m *CameraConfig) GetSetIsometric() bool {
	if m != nil {
		return m.SetIsometric
	}
	return false
}

func (m *CameraConfig) GetIsometric() bool {
	if m != nil {
		return m.Isometric
	}
	return false
}

type MaskConfig struct {
	SetName              bool     `protobuf:"varint,1,opt,name=SetName,proto3" json:"SetName,omitempty"`
	Name                 string   `protobuf:"bytes,2,opt,name=Name,proto3" json:"Name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MaskConfig) Reset()         { *m = MaskConfig{} }
func (m *MaskConfig) String() string { return proto.CompactTextString(m) }
func (*MaskConfig) ProtoMessage()    {}
func (*MaskConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_facade_5db679803abb8e98, []int{5}
}
func (m *MaskConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MaskConfig.Unmarshal(m, b)
}
func (m *MaskConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MaskConfig.Marshal(b, m, deterministic)
}
func (dst *MaskConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MaskConfig.Merge(dst, src)
}
func (m *MaskConfig) XXX_Size() int {
	return xxx_messageInfo_MaskConfig.Size(m)
}
func (m *MaskConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_MaskConfig.DiscardUnknown(m)
}

var xxx_messageInfo_MaskConfig proto.InternalMessageInfo

func (m *MaskConfig) GetSetName() bool {
	if m != nil {
		return m.SetName
	}
	return false
}

func (m *MaskConfig) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type GridConfig struct {
	SetWidth             bool     `protobuf:"varint,1,opt,name=SetWidth,proto3" json:"SetWidth,omitempty"`
	Width                uint64   `protobuf:"varint,2,opt,name=Width,proto3" json:"Width,omitempty"`
	SetHeight            bool     `protobuf:"varint,3,opt,name=SetHeight,proto3" json:"SetHeight,omitempty"`
	Height               uint64   `protobuf:"varint,4,opt,name=Height,proto3" json:"Height,omitempty"`
	SetDownward          bool     `protobuf:"varint,5,opt,name=SetDownward,proto3" json:"SetDownward,omitempty"`
	Downward             bool     `protobuf:"varint,6,opt,name=Downward,proto3" json:"Downward,omitempty"`
	SetSpeed             bool     `protobuf:"varint,7,opt,name=SetSpeed,proto3" json:"SetSpeed,omitempty"`
	Speed                float64  `protobuf:"fixed64,8,opt,name=Speed,proto3" json:"Speed,omitempty"`
	SetBuffer            bool     `protobuf:"varint,9,opt,name=SetBuffer,proto3" json:"SetBuffer,omitempty"`
	Buffer               uint64   `protobuf:"varint,10,opt,name=Buffer,proto3" json:"Buffer,omitempty"`
	SetTerminal          bool     `protobuf:"varint,11,opt,name=SetTerminal,proto3" json:"SetTerminal,omitempty"`
	Terminal             bool     `protobuf:"varint,12,opt,name=Terminal,proto3" json:"Terminal,omitempty"`
	SetVert              bool     `protobuf:"varint,13,opt,name=SetVert,proto3" json:"SetVert,omitempty"`
	Vert                 string   `protobuf:"bytes,14,opt,name=Vert,proto3" json:"Vert,omitempty"`
	SetFrag              bool     `protobuf:"varint,15,opt,name=SetFrag,proto3" json:"SetFrag,omitempty"`
	Frag                 string   `protobuf:"bytes,16,opt,name=Frag,proto3" json:"Frag,omitempty"`
	SetFill              bool     `protobuf:"varint,17,opt,name=SetFill,proto3" json:"SetFill,omitempty"`
	Fill                 string   `protobuf:"bytes,18,opt,name=Fill,proto3" json:"Fill,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GridConfig) Reset()         { *m = GridConfig{} }
func (m *GridConfig) String() string { return proto.CompactTextString(m) }
func (*GridConfig) ProtoMessage()    {}
func (*GridConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_facade_5db679803abb8e98, []int{6}
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

func (m *GridConfig) GetSetWidth() bool {
	if m != nil {
		return m.SetWidth
	}
	return false
}

func (m *GridConfig) GetWidth() uint64 {
	if m != nil {
		return m.Width
	}
	return 0
}

func (m *GridConfig) GetSetHeight() bool {
	if m != nil {
		return m.SetHeight
	}
	return false
}

func (m *GridConfig) GetHeight() uint64 {
	if m != nil {
		return m.Height
	}
	return 0
}

func (m *GridConfig) GetSetDownward() bool {
	if m != nil {
		return m.SetDownward
	}
	return false
}

func (m *GridConfig) GetDownward() bool {
	if m != nil {
		return m.Downward
	}
	return false
}

func (m *GridConfig) GetSetSpeed() bool {
	if m != nil {
		return m.SetSpeed
	}
	return false
}

func (m *GridConfig) GetSpeed() float64 {
	if m != nil {
		return m.Speed
	}
	return 0
}

func (m *GridConfig) GetSetBuffer() bool {
	if m != nil {
		return m.SetBuffer
	}
	return false
}

func (m *GridConfig) GetBuffer() uint64 {
	if m != nil {
		return m.Buffer
	}
	return 0
}

func (m *GridConfig) GetSetTerminal() bool {
	if m != nil {
		return m.SetTerminal
	}
	return false
}

func (m *GridConfig) GetTerminal() bool {
	if m != nil {
		return m.Terminal
	}
	return false
}

func (m *GridConfig) GetSetVert() bool {
	if m != nil {
		return m.SetVert
	}
	return false
}

func (m *GridConfig) GetVert() string {
	if m != nil {
		return m.Vert
	}
	return ""
}

func (m *GridConfig) GetSetFrag() bool {
	if m != nil {
		return m.SetFrag
	}
	return false
}

func (m *GridConfig) GetFrag() string {
	if m != nil {
		return m.Frag
	}
	return ""
}

func (m *GridConfig) GetSetFill() bool {
	if m != nil {
		return m.SetFill
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
	proto.RegisterType((*Empty)(nil), "facade.Empty")
	proto.RegisterType((*Text)(nil), "facade.Text")
	proto.RegisterType((*Config)(nil), "facade.Config")
	proto.RegisterType((*FontConfig)(nil), "facade.FontConfig")
	proto.RegisterType((*CameraConfig)(nil), "facade.CameraConfig")
	proto.RegisterType((*MaskConfig)(nil), "facade.MaskConfig")
	proto.RegisterType((*GridConfig)(nil), "facade.GridConfig")
	proto.RegisterEnum("facade.Mode", Mode_name, Mode_value)
}

func init() { proto.RegisterFile("facade.proto", fileDescriptor_facade_5db679803abb8e98) }

var fileDescriptor_facade_5db679803abb8e98 = []byte{
	// 552 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x54, 0xcd, 0x6e, 0xd3, 0x40,
	0x10, 0xc6, 0xad, 0xe3, 0x38, 0x13, 0x37, 0x84, 0x51, 0x85, 0x56, 0x81, 0x43, 0x64, 0x09, 0x14,
	0x55, 0xa8, 0x87, 0x70, 0xe3, 0x06, 0x0d, 0x29, 0x3d, 0x14, 0x21, 0x3b, 0x02, 0x09, 0x4e, 0x6e,
	0xbd, 0x49, 0x2d, 0xec, 0x38, 0x72, 0x37, 0x2a, 0x5c, 0x78, 0x12, 0x9e, 0x8e, 0x27, 0x41, 0x33,
	0xb3, 0xfe, 0x09, 0x37, 0x4e, 0xde, 0xef, 0x67, 0x56, 0x9f, 0x67, 0x67, 0x17, 0x82, 0x75, 0x72,
	0x9b, 0xa4, 0xfa, 0x7c, 0x57, 0x95, 0xa6, 0x44, 0x4f, 0x50, 0xd8, 0x87, 0xde, 0xfb, 0x62, 0x67,
	0x7e, 0x86, 0x13, 0x70, 0x57, 0xfa, 0x87, 0x41, 0x94, 0xaf, 0x72, 0xa6, 0xce, 0x2c, 0x88, 0x78,
	0x1d, 0xfe, 0x3e, 0x02, 0xef, 0xa2, 0xdc, 0xae, 0xb3, 0x0d, 0x4e, 0xc0, 0x8f, 0xb5, 0x59, 0xe8,
	0x9b, 0xfd, 0x86, 0x2d, 0x7e, 0xd4, 0x60, 0x3c, 0x85, 0x9e, 0x08, 0x47, 0x2c, 0x08, 0x40, 0x05,
	0xfd, 0x58, 0x9b, 0xeb, 0x32, 0xd5, 0xea, 0x98, 0xf9, 0x1a, 0xe2, 0x14, 0x5c, 0xa6, 0xdd, 0xa9,
	0x33, 0x1b, 0xcd, 0x83, 0x73, 0x1b, 0x90, 0xb8, 0x88, 0x15, 0x7c, 0x09, 0xee, 0xb2, 0xdc, 0x1a,
	0xd5, 0x9b, 0x3a, 0xb3, 0xe1, 0x1c, 0x6b, 0x07, 0x71, 0x92, 0x27, 0x62, 0x1d, 0x5f, 0x81, 0x77,
	0x91, 0x14, 0xba, 0x4a, 0x94, 0xc7, 0xce, 0xd3, 0xda, 0x29, 0xac, 0xf5, 0x5a, 0x0f, 0xed, 0x7a,
	0x9d, 0xdc, 0x7f, 0x57, 0xfd, 0xc3, 0x5d, 0x89, 0xab, 0x77, 0xa5, 0x35, 0xf9, 0x2e, 0xab, 0x2c,
	0x55, 0xfe, 0xa1, 0x8f, 0xb8, 0xda, 0x47, 0xeb, 0xf0, 0x0d, 0x40, 0x9b, 0xc8, 0xfe, 0xef, 0xc7,
	0xa4, 0xd0, 0xb6, 0x41, 0x35, 0xa4, 0xd6, 0x32, 0x4d, 0xed, 0x19, 0x44, 0xbc, 0x0e, 0x7f, 0x41,
	0xd0, 0xcd, 0x68, 0xab, 0xbf, 0x96, 0x65, 0xd1, 0xa9, 0x26, 0x48, 0xd5, 0x4c, 0x53, 0xb5, 0x13,
	0xf1, 0x1a, 0x43, 0x08, 0x62, 0x6d, 0xae, 0xee, 0xcb, 0x42, 0x9b, 0x2a, 0xbb, 0xb5, 0x0d, 0x3e,
	0xe0, 0xf0, 0x39, 0x0c, 0x5a, 0x83, 0xcb, 0x86, 0x96, 0xa0, 0xec, 0xed, 0x7f, 0xff, 0x67, 0xf6,
	0x3f, 0xc7, 0x00, 0x6d, 0x33, 0xec, 0x68, 0x7c, 0xc9, 0x52, 0x73, 0xd7, 0x19, 0x0d, 0xc6, 0x34,
	0x1a, 0x22, 0x50, 0xbd, 0x1b, 0x09, 0xa0, 0x68, 0xb1, 0x36, 0x1f, 0x74, 0xb6, 0xb9, 0x33, 0x36,
	0x7b, 0x4b, 0xe0, 0x53, 0xf0, 0xac, 0xe4, 0x72, 0x91, 0x45, 0x38, 0x85, 0x21, 0x8d, 0x5c, 0xf9,
	0xb0, 0x7d, 0x48, 0xaa, 0x94, 0x67, 0xc3, 0x8f, 0xba, 0x14, 0x25, 0x69, 0x64, 0x4f, 0x92, 0x74,
	0xb5, 0x58, 0x9b, 0x78, 0xa7, 0x75, 0xca, 0x03, 0x20, 0x29, 0x19, 0x53, 0x4a, 0x11, 0x7c, 0xee,
	0xb1, 0x00, 0x9b, 0xf2, 0xdd, 0x7e, 0xbd, 0xd6, 0x95, 0x1a, 0x34, 0x29, 0x85, 0xa0, 0x94, 0x56,
	0x02, 0x49, 0x69, 0x79, 0x49, 0xb9, 0xd2, 0x55, 0x91, 0x6d, 0x93, 0x5c, 0x0d, 0x9b, 0x94, 0x35,
	0x45, 0x49, 0x1a, 0x39, 0x90, 0x24, 0x8d, 0x26, 0x07, 0xf1, 0x59, 0x57, 0x46, 0x9d, 0x34, 0x07,
	0x41, 0x90, 0x0e, 0x82, 0xe9, 0x91, 0x1c, 0x04, 0x73, 0xe2, 0x5e, 0x56, 0xc9, 0x46, 0x3d, 0x6e,
	0xdc, 0x04, 0xc9, 0xcd, 0xf4, 0x58, 0xdc, 0xcc, 0x59, 0x77, 0x96, 0xe7, 0xea, 0x49, 0xeb, 0xce,
	0xf2, 0x9c, 0xdd, 0x44, 0xa3, 0x75, 0x67, 0x79, 0x7e, 0xf6, 0x4c, 0x2e, 0x29, 0xfa, 0xe0, 0x5e,
	0x46, 0x57, 0x8b, 0xf1, 0x23, 0x1c, 0x40, 0x6f, 0x11, 0xbd, 0x5d, 0xae, 0xc6, 0xe3, 0xf9, 0x37,
	0xf0, 0x96, 0x7c, 0x29, 0xf0, 0x0c, 0x06, 0x32, 0x06, 0xfb, 0x4a, 0xe3, 0xa8, 0xb9, 0x7e, 0x4c,
	0x4d, 0x4e, 0x6a, 0xcc, 0x4f, 0x0d, 0xbe, 0x00, 0xf7, 0x53, 0xb6, 0xd3, 0xd8, 0xdc, 0x78, 0x7a,
	0x64, 0xfe, 0x31, 0xcd, 0x9c, 0x1b, 0x8f, 0x5f, 0xaa, 0xd7, 0x7f, 0x03, 0x00, 0x00, 0xff, 0xff,
	0x54, 0xd5, 0x82, 0x18, 0xb9, 0x04, 0x00, 0x00,
}
