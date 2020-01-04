// Code generated by protoc-gen-go. DO NOT EDIT.
// source: facade.proto

package facade

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Mode int32

const (
	Mode_TERM  Mode = 0
	Mode_LINE  Mode = 1
	Mode_DRAFT Mode = 16
)

var Mode_name = map[int32]string{
	0:  "TERM",
	1:  "LINE",
	16: "DRAFT",
}

var Mode_value = map[string]int32{
	"TERM":  0,
	"LINE":  1,
	"DRAFT": 16,
}

func (x Mode) String() string {
	return proto.EnumName(Mode_name, int32(x))
}

func (Mode) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_5478f20a07eaa28e, []int{0}
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
	return fileDescriptor_5478f20a07eaa28e, []int{0}
}

func (m *Empty) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Empty.Unmarshal(m, b)
}
func (m *Empty) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Empty.Marshal(b, m, deterministic)
}
func (m *Empty) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Empty.Merge(m, src)
}
func (m *Empty) XXX_Size() int {
	return xxx_messageInfo_Empty.Size(m)
}
func (m *Empty) XXX_DiscardUnknown() {
	xxx_messageInfo_Empty.DiscardUnknown(m)
}

var xxx_messageInfo_Empty proto.InternalMessageInfo

type Status struct {
	Success              bool     `protobuf:"varint,1,opt,name=Success,proto3" json:"Success,omitempty"`
	Error                string   `protobuf:"bytes,2,opt,name=Error,proto3" json:"Error,omitempty"`
	Info                 string   `protobuf:"bytes,3,opt,name=Info,proto3" json:"Info,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Status) Reset()         { *m = Status{} }
func (m *Status) String() string { return proto.CompactTextString(m) }
func (*Status) ProtoMessage()    {}
func (*Status) Descriptor() ([]byte, []int) {
	return fileDescriptor_5478f20a07eaa28e, []int{1}
}

func (m *Status) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Status.Unmarshal(m, b)
}
func (m *Status) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Status.Marshal(b, m, deterministic)
}
func (m *Status) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Status.Merge(m, src)
}
func (m *Status) XXX_Size() int {
	return xxx_messageInfo_Status.Size(m)
}
func (m *Status) XXX_DiscardUnknown() {
	xxx_messageInfo_Status.DiscardUnknown(m)
}

var xxx_messageInfo_Status proto.InternalMessageInfo

func (m *Status) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func (m *Status) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

func (m *Status) GetInfo() string {
	if m != nil {
		return m.Info
	}
	return ""
}

type RawText struct {
	Raw                  []byte   `protobuf:"bytes,1,opt,name=Raw,proto3" json:"Raw,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RawText) Reset()         { *m = RawText{} }
func (m *RawText) String() string { return proto.CompactTextString(m) }
func (*RawText) ProtoMessage()    {}
func (*RawText) Descriptor() ([]byte, []int) {
	return fileDescriptor_5478f20a07eaa28e, []int{2}
}

func (m *RawText) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RawText.Unmarshal(m, b)
}
func (m *RawText) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RawText.Marshal(b, m, deterministic)
}
func (m *RawText) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RawText.Merge(m, src)
}
func (m *RawText) XXX_Size() int {
	return xxx_messageInfo_RawText.Size(m)
}
func (m *RawText) XXX_DiscardUnknown() {
	xxx_messageInfo_RawText.DiscardUnknown(m)
}

var xxx_messageInfo_RawText proto.InternalMessageInfo

func (m *RawText) GetRaw() []byte {
	if m != nil {
		return m.Raw
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
	Terminal             *TermConfig   `protobuf:"bytes,13,opt,name=Terminal,proto3" json:"Terminal,omitempty"`
	Lines                *LineConfig   `protobuf:"bytes,14,opt,name=Lines,proto3" json:"Lines,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *Config) Reset()         { *m = Config{} }
func (m *Config) String() string { return proto.CompactTextString(m) }
func (*Config) ProtoMessage()    {}
func (*Config) Descriptor() ([]byte, []int) {
	return fileDescriptor_5478f20a07eaa28e, []int{3}
}

func (m *Config) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Config.Unmarshal(m, b)
}
func (m *Config) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Config.Marshal(b, m, deterministic)
}
func (m *Config) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Config.Merge(m, src)
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
	return Mode_TERM
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

func (m *Config) GetTerminal() *TermConfig {
	if m != nil {
		return m.Terminal
	}
	return nil
}

func (m *Config) GetLines() *LineConfig {
	if m != nil {
		return m.Lines
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
	return fileDescriptor_5478f20a07eaa28e, []int{4}
}

func (m *FontConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FontConfig.Unmarshal(m, b)
}
func (m *FontConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FontConfig.Marshal(b, m, deterministic)
}
func (m *FontConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FontConfig.Merge(m, src)
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
	return fileDescriptor_5478f20a07eaa28e, []int{5}
}

func (m *CameraConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CameraConfig.Unmarshal(m, b)
}
func (m *CameraConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CameraConfig.Marshal(b, m, deterministic)
}
func (m *CameraConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CameraConfig.Merge(m, src)
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
	return fileDescriptor_5478f20a07eaa28e, []int{6}
}

func (m *MaskConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MaskConfig.Unmarshal(m, b)
}
func (m *MaskConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MaskConfig.Marshal(b, m, deterministic)
}
func (m *MaskConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MaskConfig.Merge(m, src)
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

type TermConfig struct {
	Grid                 *GridConfig `protobuf:"bytes,1,opt,name=Grid,proto3" json:"Grid,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *TermConfig) Reset()         { *m = TermConfig{} }
func (m *TermConfig) String() string { return proto.CompactTextString(m) }
func (*TermConfig) ProtoMessage()    {}
func (*TermConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_5478f20a07eaa28e, []int{7}
}

func (m *TermConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TermConfig.Unmarshal(m, b)
}
func (m *TermConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TermConfig.Marshal(b, m, deterministic)
}
func (m *TermConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TermConfig.Merge(m, src)
}
func (m *TermConfig) XXX_Size() int {
	return xxx_messageInfo_TermConfig.Size(m)
}
func (m *TermConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_TermConfig.DiscardUnknown(m)
}

var xxx_messageInfo_TermConfig proto.InternalMessageInfo

func (m *TermConfig) GetGrid() *GridConfig {
	if m != nil {
		return m.Grid
	}
	return nil
}

type LineConfig struct {
	Grid                 *GridConfig `protobuf:"bytes,1,opt,name=Grid,proto3" json:"Grid,omitempty"`
	SetDownward          bool        `protobuf:"varint,3,opt,name=SetDownward,proto3" json:"SetDownward,omitempty"`
	Downward             bool        `protobuf:"varint,4,opt,name=Downward,proto3" json:"Downward,omitempty"`
	SetSpeed             bool        `protobuf:"varint,5,opt,name=SetSpeed,proto3" json:"SetSpeed,omitempty"`
	Speed                float64     `protobuf:"fixed64,6,opt,name=Speed,proto3" json:"Speed,omitempty"`
	SetAdaptive          bool        `protobuf:"varint,7,opt,name=SetAdaptive,proto3" json:"SetAdaptive,omitempty"`
	Adaptive             bool        `protobuf:"varint,8,opt,name=Adaptive,proto3" json:"Adaptive,omitempty"`
	SetDrop              bool        `protobuf:"varint,9,opt,name=SetDrop,proto3" json:"SetDrop,omitempty"`
	Drop                 bool        `protobuf:"varint,10,opt,name=Drop,proto3" json:"Drop,omitempty"`
	SetSmooth            bool        `protobuf:"varint,11,opt,name=SetSmooth,proto3" json:"SetSmooth,omitempty"`
	Smooth               bool        `protobuf:"varint,12,opt,name=Smooth,proto3" json:"Smooth,omitempty"`
	SetBuffer            bool        `protobuf:"varint,13,opt,name=SetBuffer,proto3" json:"SetBuffer,omitempty"`
	Buffer               uint64      `protobuf:"varint,14,opt,name=Buffer,proto3" json:"Buffer,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *LineConfig) Reset()         { *m = LineConfig{} }
func (m *LineConfig) String() string { return proto.CompactTextString(m) }
func (*LineConfig) ProtoMessage()    {}
func (*LineConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_5478f20a07eaa28e, []int{8}
}

func (m *LineConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LineConfig.Unmarshal(m, b)
}
func (m *LineConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LineConfig.Marshal(b, m, deterministic)
}
func (m *LineConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LineConfig.Merge(m, src)
}
func (m *LineConfig) XXX_Size() int {
	return xxx_messageInfo_LineConfig.Size(m)
}
func (m *LineConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_LineConfig.DiscardUnknown(m)
}

var xxx_messageInfo_LineConfig proto.InternalMessageInfo

func (m *LineConfig) GetGrid() *GridConfig {
	if m != nil {
		return m.Grid
	}
	return nil
}

func (m *LineConfig) GetSetDownward() bool {
	if m != nil {
		return m.SetDownward
	}
	return false
}

func (m *LineConfig) GetDownward() bool {
	if m != nil {
		return m.Downward
	}
	return false
}

func (m *LineConfig) GetSetSpeed() bool {
	if m != nil {
		return m.SetSpeed
	}
	return false
}

func (m *LineConfig) GetSpeed() float64 {
	if m != nil {
		return m.Speed
	}
	return 0
}

func (m *LineConfig) GetSetAdaptive() bool {
	if m != nil {
		return m.SetAdaptive
	}
	return false
}

func (m *LineConfig) GetAdaptive() bool {
	if m != nil {
		return m.Adaptive
	}
	return false
}

func (m *LineConfig) GetSetDrop() bool {
	if m != nil {
		return m.SetDrop
	}
	return false
}

func (m *LineConfig) GetDrop() bool {
	if m != nil {
		return m.Drop
	}
	return false
}

func (m *LineConfig) GetSetSmooth() bool {
	if m != nil {
		return m.SetSmooth
	}
	return false
}

func (m *LineConfig) GetSmooth() bool {
	if m != nil {
		return m.Smooth
	}
	return false
}

func (m *LineConfig) GetSetBuffer() bool {
	if m != nil {
		return m.SetBuffer
	}
	return false
}

func (m *LineConfig) GetBuffer() uint64 {
	if m != nil {
		return m.Buffer
	}
	return 0
}

type GridConfig struct {
	SetWidth             bool     `protobuf:"varint,1,opt,name=SetWidth,proto3" json:"SetWidth,omitempty"`
	Width                uint64   `protobuf:"varint,2,opt,name=Width,proto3" json:"Width,omitempty"`
	SetHeight            bool     `protobuf:"varint,3,opt,name=SetHeight,proto3" json:"SetHeight,omitempty"`
	Height               uint64   `protobuf:"varint,4,opt,name=Height,proto3" json:"Height,omitempty"`
	SetFill              bool     `protobuf:"varint,7,opt,name=SetFill,proto3" json:"SetFill,omitempty"`
	Fill                 string   `protobuf:"bytes,8,opt,name=Fill,proto3" json:"Fill,omitempty"`
	SetVert              bool     `protobuf:"varint,9,opt,name=SetVert,proto3" json:"SetVert,omitempty"`
	Vert                 string   `protobuf:"bytes,10,opt,name=Vert,proto3" json:"Vert,omitempty"`
	SetFrag              bool     `protobuf:"varint,11,opt,name=SetFrag,proto3" json:"SetFrag,omitempty"`
	Frag                 string   `protobuf:"bytes,12,opt,name=Frag,proto3" json:"Frag,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GridConfig) Reset()         { *m = GridConfig{} }
func (m *GridConfig) String() string { return proto.CompactTextString(m) }
func (*GridConfig) ProtoMessage()    {}
func (*GridConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_5478f20a07eaa28e, []int{9}
}

func (m *GridConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GridConfig.Unmarshal(m, b)
}
func (m *GridConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GridConfig.Marshal(b, m, deterministic)
}
func (m *GridConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GridConfig.Merge(m, src)
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

func init() {
	proto.RegisterEnum("facade.Mode", Mode_name, Mode_value)
	proto.RegisterType((*Empty)(nil), "facade.Empty")
	proto.RegisterType((*Status)(nil), "facade.Status")
	proto.RegisterType((*RawText)(nil), "facade.RawText")
	proto.RegisterType((*Config)(nil), "facade.Config")
	proto.RegisterType((*FontConfig)(nil), "facade.FontConfig")
	proto.RegisterType((*CameraConfig)(nil), "facade.CameraConfig")
	proto.RegisterType((*MaskConfig)(nil), "facade.MaskConfig")
	proto.RegisterType((*TermConfig)(nil), "facade.TermConfig")
	proto.RegisterType((*LineConfig)(nil), "facade.LineConfig")
	proto.RegisterType((*GridConfig)(nil), "facade.GridConfig")
}

func init() { proto.RegisterFile("facade.proto", fileDescriptor_5478f20a07eaa28e) }

var fileDescriptor_5478f20a07eaa28e = []byte{
	// 705 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x55, 0xcf, 0x6f, 0xd3, 0x4a,
	0x10, 0xae, 0x53, 0xc7, 0x49, 0xa6, 0x69, 0x5f, 0xb4, 0xaa, 0x9e, 0xac, 0xbe, 0x77, 0x88, 0x2c,
	0x81, 0x22, 0x04, 0x39, 0x14, 0x4e, 0xdc, 0x4a, 0x93, 0x40, 0xa5, 0xb6, 0x12, 0x4e, 0x04, 0x12,
	0xb7, 0x6d, 0xbc, 0x49, 0x2d, 0xe2, 0x6c, 0xb4, 0xde, 0x10, 0x7a, 0xe1, 0x7f, 0xe0, 0x1f, 0xe3,
	0x2f, 0x42, 0x08, 0xcd, 0xcc, 0xfa, 0x47, 0xd4, 0x0b, 0xdc, 0xf6, 0xfb, 0xe6, 0x9b, 0x99, 0xf5,
	0x37, 0xa3, 0x35, 0x74, 0x17, 0x72, 0x2e, 0x13, 0x35, 0xdc, 0x18, 0x6d, 0xb5, 0x08, 0x18, 0x45,
	0x2d, 0x68, 0x8e, 0xb3, 0x8d, 0x7d, 0x88, 0xae, 0x21, 0x98, 0x5a, 0x69, 0xb7, 0xb9, 0x08, 0xa1,
	0x35, 0xdd, 0xce, 0xe7, 0x2a, 0xcf, 0x43, 0xaf, 0xef, 0x0d, 0xda, 0x71, 0x01, 0xc5, 0x29, 0x34,
	0xc7, 0xc6, 0x68, 0x13, 0x36, 0xfa, 0xde, 0xa0, 0x13, 0x33, 0x10, 0x02, 0xfc, 0xab, 0xf5, 0x42,
	0x87, 0x87, 0x44, 0xd2, 0x39, 0xfa, 0x0f, 0x5a, 0xb1, 0xdc, 0xcd, 0xd4, 0x57, 0x2b, 0x7a, 0x70,
	0x18, 0xcb, 0x1d, 0x95, 0xea, 0xc6, 0x78, 0x8c, 0x7e, 0x34, 0x20, 0xb8, 0xd4, 0xeb, 0x45, 0xba,
	0x14, 0x67, 0xd0, 0x9e, 0x2a, 0x3b, 0x52, 0x77, 0xdb, 0xa5, 0x6b, 0x56, 0x62, 0xec, 0xc6, 0x81,
	0x06, 0x05, 0x18, 0xd0, 0xed, 0x94, 0xbd, 0xd1, 0x89, 0xa2, 0x86, 0x78, 0x3b, 0x86, 0xa2, 0x0f,
	0x3e, 0xd1, 0x7e, 0xdf, 0x1b, 0x9c, 0x9c, 0x77, 0x87, 0xee, 0x7b, 0x91, 0x8b, 0x29, 0x22, 0x9e,
	0x82, 0x3f, 0xd1, 0x6b, 0x1b, 0x36, 0xfb, 0xde, 0xe0, 0xe8, 0x5c, 0x14, 0x0a, 0xe4, 0xf8, 0x3e,
	0x31, 0xc5, 0xc5, 0x73, 0x08, 0x2e, 0x65, 0xa6, 0x8c, 0x0c, 0x03, 0x52, 0x9e, 0x16, 0x4a, 0x66,
	0x9d, 0xd6, 0x69, 0xb0, 0xea, 0x8d, 0xcc, 0x3f, 0x87, 0xad, 0xfd, 0xaa, 0xc8, 0x15, 0x55, 0xf1,
	0x2c, 0x86, 0xd0, 0x9e, 0x29, 0x93, 0xa5, 0x6b, 0xb9, 0x0a, 0x8f, 0xf7, 0xb5, 0xc8, 0x3b, 0x6d,
	0xa9, 0x11, 0x03, 0x68, 0x5e, 0xa7, 0x6b, 0x95, 0x87, 0x27, 0xfb, 0x62, 0x24, 0x9d, 0x98, 0x05,
	0xd1, 0x6b, 0x80, 0xea, 0x1b, 0x9c, 0x43, 0xb7, 0x32, 0x53, 0xe5, 0xfc, 0x18, 0xe2, 0xa4, 0x88,
	0xe6, 0xf1, 0xd1, 0x39, 0xfa, 0x06, 0xdd, 0xfa, 0x57, 0xb9, 0xec, 0x4f, 0x5a, 0x67, 0xb5, 0x6c,
	0x84, 0x98, 0x4d, 0x34, 0x66, 0x7b, 0x31, 0x9d, 0x45, 0x04, 0xdd, 0xa9, 0xb2, 0x57, 0xb9, 0xce,
	0x94, 0x35, 0xe9, 0xdc, 0x8d, 0x64, 0x8f, 0x13, 0xff, 0x43, 0xa7, 0x12, 0xf8, 0x24, 0xa8, 0x08,
	0xbc, 0x7b, 0xe5, 0xd4, 0x5f, 0xde, 0xfd, 0x15, 0x40, 0xe5, 0x1c, 0xce, 0xe1, 0xad, 0x49, 0x13,
	0x4a, 0xac, 0xd9, 0x85, 0x5c, 0x31, 0x07, 0x3c, 0x47, 0x3f, 0x1b, 0x00, 0x95, 0x87, 0x7f, 0x9a,
	0x26, 0xfa, 0x70, 0x84, 0xab, 0xa9, 0x77, 0xeb, 0x9d, 0x34, 0x89, 0xfb, 0xd2, 0x3a, 0x85, 0xcb,
	0x5c, 0x86, 0xf9, 0x3b, 0xdb, 0xf5, 0xd8, 0x54, 0xd9, 0xe9, 0x46, 0xa9, 0x84, 0xd6, 0x8f, 0x17,
	0x9d, 0x30, 0x2e, 0x3a, 0x07, 0x02, 0x72, 0x96, 0x81, 0xeb, 0x77, 0x91, 0xc8, 0x8d, 0x4d, 0xbf,
	0x28, 0xda, 0x2e, 0xee, 0x57, 0x50, 0x58, 0xb3, 0x0c, 0xb7, 0xb9, 0x66, 0x19, 0x63, 0x23, 0x47,
	0x46, 0x6f, 0xc2, 0x4e, 0x69, 0x24, 0x42, 0x34, 0x92, 0x68, 0x20, 0x9a, 0xce, 0x38, 0x22, 0xbc,
	0x4d, 0xa6, 0xb5, 0xbd, 0x0f, 0x8f, 0x78, 0x44, 0x25, 0x21, 0xfe, 0x85, 0xc0, 0x85, 0xba, 0x14,
	0x72, 0xc8, 0x65, 0xbd, 0xd9, 0x2e, 0x16, 0xca, 0xd0, 0x46, 0x73, 0x16, 0x13, 0x98, 0xe5, 0x42,
	0xb8, 0xbf, 0x7e, 0xec, 0x50, 0xf4, 0xcb, 0x03, 0xa8, 0xcc, 0x75, 0xc6, 0x7c, 0x4c, 0x13, 0x7b,
	0x5f, 0x7b, 0x01, 0x08, 0xa3, 0x31, 0x1c, 0x68, 0x50, 0x05, 0x06, 0xae, 0xed, 0x3b, 0x95, 0x2e,
	0xef, 0xad, 0x1b, 0x43, 0x45, 0x60, 0x5b, 0x17, 0xf2, 0xb9, 0xad, 0xe3, 0xd9, 0x90, 0x49, 0xba,
	0x5a, 0x39, 0x2b, 0x0b, 0x88, 0x86, 0x10, 0xdd, 0xe6, 0xcd, 0x22, 0x8e, 0xd5, 0x1f, 0x94, 0xb1,
	0x35, 0xfb, 0x10, 0xa2, 0x9a, 0x68, 0x60, 0x35, 0x71, 0xae, 0xb6, 0x91, 0x4b, 0x67, 0x5e, 0x01,
	0xa9, 0x36, 0xd2, 0x5d, 0x57, 0xdb, 0xc8, 0xe5, 0xb3, 0x27, 0xfc, 0x4e, 0x89, 0x36, 0xf8, 0xb3,
	0x71, 0x7c, 0xd3, 0x3b, 0xc0, 0xd3, 0xf5, 0xd5, 0xed, 0xb8, 0xe7, 0x89, 0x0e, 0x34, 0x47, 0xf1,
	0xc5, 0x64, 0xd6, 0xeb, 0x9d, 0x7f, 0xf7, 0x20, 0x98, 0xd0, 0x2e, 0x8a, 0x17, 0xd0, 0x61, 0xb7,
	0xb6, 0x46, 0x89, 0x93, 0xf2, 0x31, 0x22, 0xea, 0xac, 0xc4, 0xfc, 0x7c, 0x47, 0x07, 0x62, 0x08,
	0xad, 0x51, 0x9a, 0x6f, 0x56, 0xf2, 0x41, 0xfc, 0x53, 0x04, 0xdd, 0x6b, 0xfc, 0x58, 0x3d, 0xf0,
	0xf0, 0xa1, 0x79, 0xbf, 0x55, 0xe6, 0x41, 0x1c, 0x17, 0x41, 0xfa, 0x25, 0x3c, 0xd6, 0xde, 0x05,
	0xf4, 0xf3, 0x78, 0xf9, 0x3b, 0x00, 0x00, 0xff, 0xff, 0x03, 0x0e, 0x4e, 0x8d, 0x4c, 0x06, 0x00,
	0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// FacadeClient is the client API for Facade service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type FacadeClient interface {
	Configure(ctx context.Context, in *Config, opts ...grpc.CallOption) (*Status, error)
	Display(ctx context.Context, opts ...grpc.CallOption) (Facade_DisplayClient, error)
	Query(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Status, error)
}

type facadeClient struct {
	cc *grpc.ClientConn
}

func NewFacadeClient(cc *grpc.ClientConn) FacadeClient {
	return &facadeClient{cc}
}

func (c *facadeClient) Configure(ctx context.Context, in *Config, opts ...grpc.CallOption) (*Status, error) {
	out := new(Status)
	err := c.cc.Invoke(ctx, "/facade.Facade/Configure", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *facadeClient) Display(ctx context.Context, opts ...grpc.CallOption) (Facade_DisplayClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Facade_serviceDesc.Streams[0], "/facade.Facade/Display", opts...)
	if err != nil {
		return nil, err
	}
	x := &facadeDisplayClient{stream}
	return x, nil
}

type Facade_DisplayClient interface {
	Send(*RawText) error
	CloseAndRecv() (*Status, error)
	grpc.ClientStream
}

type facadeDisplayClient struct {
	grpc.ClientStream
}

func (x *facadeDisplayClient) Send(m *RawText) error {
	return x.ClientStream.SendMsg(m)
}

func (x *facadeDisplayClient) CloseAndRecv() (*Status, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(Status)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *facadeClient) Query(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Status, error) {
	out := new(Status)
	err := c.cc.Invoke(ctx, "/facade.Facade/Query", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FacadeServer is the server API for Facade service.
type FacadeServer interface {
	Configure(context.Context, *Config) (*Status, error)
	Display(Facade_DisplayServer) error
	Query(context.Context, *Empty) (*Status, error)
}

// UnimplementedFacadeServer can be embedded to have forward compatible implementations.
type UnimplementedFacadeServer struct {
}

func (*UnimplementedFacadeServer) Configure(ctx context.Context, req *Config) (*Status, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Configure not implemented")
}
func (*UnimplementedFacadeServer) Display(srv Facade_DisplayServer) error {
	return status.Errorf(codes.Unimplemented, "method Display not implemented")
}
func (*UnimplementedFacadeServer) Query(ctx context.Context, req *Empty) (*Status, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Query not implemented")
}

func RegisterFacadeServer(s *grpc.Server, srv FacadeServer) {
	s.RegisterService(&_Facade_serviceDesc, srv)
}

func _Facade_Configure_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Config)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FacadeServer).Configure(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/facade.Facade/Configure",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FacadeServer).Configure(ctx, req.(*Config))
	}
	return interceptor(ctx, in, info, handler)
}

func _Facade_Display_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(FacadeServer).Display(&facadeDisplayServer{stream})
}

type Facade_DisplayServer interface {
	SendAndClose(*Status) error
	Recv() (*RawText, error)
	grpc.ServerStream
}

type facadeDisplayServer struct {
	grpc.ServerStream
}

func (x *facadeDisplayServer) SendAndClose(m *Status) error {
	return x.ServerStream.SendMsg(m)
}

func (x *facadeDisplayServer) Recv() (*RawText, error) {
	m := new(RawText)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Facade_Query_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FacadeServer).Query(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/facade.Facade/Query",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FacadeServer).Query(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _Facade_serviceDesc = grpc.ServiceDesc{
	ServiceName: "facade.Facade",
	HandlerType: (*FacadeServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Configure",
			Handler:    _Facade_Configure_Handler,
		},
		{
			MethodName: "Query",
			Handler:    _Facade_Query_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Display",
			Handler:       _Facade_Display_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "facade.proto",
}
