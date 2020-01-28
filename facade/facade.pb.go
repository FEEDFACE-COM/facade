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
	Mode_TAGS  Mode = 2
	Mode_DRAFT Mode = 16
)

var Mode_name = map[int32]string{
	0:  "TERM",
	1:  "LINE",
	2:  "TAGS",
	16: "DRAFT",
}

var Mode_value = map[string]int32{
	"TERM":  0,
	"LINE":  1,
	"TAGS":  2,
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
	Tags                 *TagConfig    `protobuf:"bytes,15,opt,name=Tags,proto3" json:"Tags,omitempty"`
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

func (m *Config) GetTags() *TagConfig {
	if m != nil {
		return m.Tags
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

type ShaderConfig struct {
	SetVert              bool     `protobuf:"varint,1,opt,name=SetVert,proto3" json:"SetVert,omitempty"`
	Vert                 string   `protobuf:"bytes,2,opt,name=Vert,proto3" json:"Vert,omitempty"`
	SetFrag              bool     `protobuf:"varint,3,opt,name=SetFrag,proto3" json:"SetFrag,omitempty"`
	Frag                 string   `protobuf:"bytes,4,opt,name=Frag,proto3" json:"Frag,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ShaderConfig) Reset()         { *m = ShaderConfig{} }
func (m *ShaderConfig) String() string { return proto.CompactTextString(m) }
func (*ShaderConfig) ProtoMessage()    {}
func (*ShaderConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_5478f20a07eaa28e, []int{7}
}

func (m *ShaderConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ShaderConfig.Unmarshal(m, b)
}
func (m *ShaderConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ShaderConfig.Marshal(b, m, deterministic)
}
func (m *ShaderConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ShaderConfig.Merge(m, src)
}
func (m *ShaderConfig) XXX_Size() int {
	return xxx_messageInfo_ShaderConfig.Size(m)
}
func (m *ShaderConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_ShaderConfig.DiscardUnknown(m)
}

var xxx_messageInfo_ShaderConfig proto.InternalMessageInfo

func (m *ShaderConfig) GetSetVert() bool {
	if m != nil {
		return m.SetVert
	}
	return false
}

func (m *ShaderConfig) GetVert() string {
	if m != nil {
		return m.Vert
	}
	return ""
}

func (m *ShaderConfig) GetSetFrag() bool {
	if m != nil {
		return m.SetFrag
	}
	return false
}

func (m *ShaderConfig) GetFrag() string {
	if m != nil {
		return m.Frag
	}
	return ""
}

type TermConfig struct {
	Shader               *ShaderConfig `protobuf:"bytes,1,opt,name=Shader,proto3" json:"Shader,omitempty"`
	Grid                 *GridConfig   `protobuf:"bytes,2,opt,name=Grid,proto3" json:"Grid,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *TermConfig) Reset()         { *m = TermConfig{} }
func (m *TermConfig) String() string { return proto.CompactTextString(m) }
func (*TermConfig) ProtoMessage()    {}
func (*TermConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_5478f20a07eaa28e, []int{8}
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

func (m *TermConfig) GetShader() *ShaderConfig {
	if m != nil {
		return m.Shader
	}
	return nil
}

func (m *TermConfig) GetGrid() *GridConfig {
	if m != nil {
		return m.Grid
	}
	return nil
}

type LineConfig struct {
	Shader               *ShaderConfig `protobuf:"bytes,1,opt,name=Shader,proto3" json:"Shader,omitempty"`
	Grid                 *GridConfig   `protobuf:"bytes,2,opt,name=Grid,proto3" json:"Grid,omitempty"`
	SetDownward          bool          `protobuf:"varint,3,opt,name=SetDownward,proto3" json:"SetDownward,omitempty"`
	Downward             bool          `protobuf:"varint,4,opt,name=Downward,proto3" json:"Downward,omitempty"`
	SetSpeed             bool          `protobuf:"varint,5,opt,name=SetSpeed,proto3" json:"SetSpeed,omitempty"`
	Speed                float64       `protobuf:"fixed64,6,opt,name=Speed,proto3" json:"Speed,omitempty"`
	SetFixed             bool          `protobuf:"varint,7,opt,name=SetFixed,proto3" json:"SetFixed,omitempty"`
	Fixed                bool          `protobuf:"varint,8,opt,name=Fixed,proto3" json:"Fixed,omitempty"`
	SetDrop              bool          `protobuf:"varint,9,opt,name=SetDrop,proto3" json:"SetDrop,omitempty"`
	Drop                 bool          `protobuf:"varint,10,opt,name=Drop,proto3" json:"Drop,omitempty"`
	SetStop              bool          `protobuf:"varint,11,opt,name=SetStop,proto3" json:"SetStop,omitempty"`
	Stop                 bool          `protobuf:"varint,12,opt,name=Stop,proto3" json:"Stop,omitempty"`
	SetBuffer            bool          `protobuf:"varint,13,opt,name=SetBuffer,proto3" json:"SetBuffer,omitempty"`
	Buffer               uint64        `protobuf:"varint,14,opt,name=Buffer,proto3" json:"Buffer,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *LineConfig) Reset()         { *m = LineConfig{} }
func (m *LineConfig) String() string { return proto.CompactTextString(m) }
func (*LineConfig) ProtoMessage()    {}
func (*LineConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_5478f20a07eaa28e, []int{9}
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

func (m *LineConfig) GetShader() *ShaderConfig {
	if m != nil {
		return m.Shader
	}
	return nil
}

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

func (m *LineConfig) GetSetFixed() bool {
	if m != nil {
		return m.SetFixed
	}
	return false
}

func (m *LineConfig) GetFixed() bool {
	if m != nil {
		return m.Fixed
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

func (m *LineConfig) GetSetStop() bool {
	if m != nil {
		return m.SetStop
	}
	return false
}

func (m *LineConfig) GetStop() bool {
	if m != nil {
		return m.Stop
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
	SetFill              bool     `protobuf:"varint,5,opt,name=SetFill,proto3" json:"SetFill,omitempty"`
	Fill                 string   `protobuf:"bytes,6,opt,name=Fill,proto3" json:"Fill,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GridConfig) Reset()         { *m = GridConfig{} }
func (m *GridConfig) String() string { return proto.CompactTextString(m) }
func (*GridConfig) ProtoMessage()    {}
func (*GridConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_5478f20a07eaa28e, []int{10}
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

type TagConfig struct {
	Shader               *ShaderConfig `protobuf:"bytes,1,opt,name=Shader,proto3" json:"Shader,omitempty"`
	SetDuration          bool          `protobuf:"varint,3,opt,name=SetDuration,proto3" json:"SetDuration,omitempty"`
	Duration             float64       `protobuf:"fixed64,4,opt,name=Duration,proto3" json:"Duration,omitempty"`
	SetFill              bool          `protobuf:"varint,5,opt,name=SetFill,proto3" json:"SetFill,omitempty"`
	Fill                 string        `protobuf:"bytes,6,opt,name=Fill,proto3" json:"Fill,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *TagConfig) Reset()         { *m = TagConfig{} }
func (m *TagConfig) String() string { return proto.CompactTextString(m) }
func (*TagConfig) ProtoMessage()    {}
func (*TagConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_5478f20a07eaa28e, []int{11}
}

func (m *TagConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TagConfig.Unmarshal(m, b)
}
func (m *TagConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TagConfig.Marshal(b, m, deterministic)
}
func (m *TagConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TagConfig.Merge(m, src)
}
func (m *TagConfig) XXX_Size() int {
	return xxx_messageInfo_TagConfig.Size(m)
}
func (m *TagConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_TagConfig.DiscardUnknown(m)
}

var xxx_messageInfo_TagConfig proto.InternalMessageInfo

func (m *TagConfig) GetShader() *ShaderConfig {
	if m != nil {
		return m.Shader
	}
	return nil
}

func (m *TagConfig) GetSetDuration() bool {
	if m != nil {
		return m.SetDuration
	}
	return false
}

func (m *TagConfig) GetDuration() float64 {
	if m != nil {
		return m.Duration
	}
	return 0
}

func (m *TagConfig) GetSetFill() bool {
	if m != nil {
		return m.SetFill
	}
	return false
}

func (m *TagConfig) GetFill() string {
	if m != nil {
		return m.Fill
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
	proto.RegisterType((*ShaderConfig)(nil), "facade.ShaderConfig")
	proto.RegisterType((*TermConfig)(nil), "facade.TermConfig")
	proto.RegisterType((*LineConfig)(nil), "facade.LineConfig")
	proto.RegisterType((*GridConfig)(nil), "facade.GridConfig")
	proto.RegisterType((*TagConfig)(nil), "facade.TagConfig")
}

func init() { proto.RegisterFile("facade.proto", fileDescriptor_5478f20a07eaa28e) }

var fileDescriptor_5478f20a07eaa28e = []byte{
	// 780 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x55, 0xcf, 0x6e, 0xd3, 0x4e,
	0x10, 0xae, 0x53, 0xc7, 0x49, 0xa6, 0x69, 0x9b, 0xdf, 0xaa, 0xfa, 0xc9, 0x2a, 0x3d, 0x44, 0x96,
	0x40, 0x11, 0x82, 0x20, 0x95, 0x1b, 0xb7, 0xd2, 0x24, 0xa5, 0x52, 0x5b, 0x09, 0x3b, 0x02, 0x89,
	0xdb, 0x36, 0xde, 0xa4, 0x86, 0x24, 0x8e, 0xd6, 0x1b, 0xa5, 0xbd, 0xf0, 0x0e, 0x5c, 0x78, 0x00,
	0x24, 0x1e, 0x91, 0x3b, 0x9a, 0x99, 0xf5, 0x9f, 0xa8, 0xa7, 0x1e, 0xb8, 0xcd, 0xf7, 0xcd, 0x37,
	0x3b, 0xb3, 0x3b, 0xb3, 0xbb, 0xd0, 0x9e, 0xca, 0x89, 0x8c, 0x55, 0x7f, 0xa5, 0x53, 0x93, 0x0a,
	0x8f, 0x51, 0xd0, 0x80, 0xfa, 0x70, 0xb1, 0x32, 0x0f, 0xc1, 0x15, 0x78, 0x91, 0x91, 0x66, 0x9d,
	0x09, 0x1f, 0x1a, 0xd1, 0x7a, 0x32, 0x51, 0x59, 0xe6, 0x3b, 0x5d, 0xa7, 0xd7, 0x0c, 0x73, 0x28,
	0x8e, 0xa0, 0x3e, 0xd4, 0x3a, 0xd5, 0x7e, 0xad, 0xeb, 0xf4, 0x5a, 0x21, 0x03, 0x21, 0xc0, 0xbd,
	0x5c, 0x4e, 0x53, 0x7f, 0x97, 0x48, 0xb2, 0x83, 0x67, 0xd0, 0x08, 0xe5, 0x66, 0xac, 0xee, 0x8d,
	0xe8, 0xc0, 0x6e, 0x28, 0x37, 0xb4, 0x54, 0x3b, 0x44, 0x33, 0xf8, 0x53, 0x03, 0xef, 0x3c, 0x5d,
	0x4e, 0x93, 0x99, 0x38, 0x86, 0x66, 0xa4, 0xcc, 0x40, 0xdd, 0xae, 0x67, 0x36, 0x59, 0x81, 0x31,
	0x1b, 0x3b, 0x6a, 0xe4, 0x60, 0x40, 0xd5, 0x29, 0x73, 0x9d, 0xc6, 0x8a, 0x12, 0x62, 0x75, 0x0c,
	0x45, 0x17, 0x5c, 0xa2, 0xdd, 0xae, 0xd3, 0x3b, 0x38, 0x6d, 0xf7, 0xed, 0x7e, 0x91, 0x0b, 0xc9,
	0x23, 0x5e, 0x80, 0x3b, 0x4a, 0x97, 0xc6, 0xaf, 0x77, 0x9d, 0xde, 0xde, 0xa9, 0xc8, 0x15, 0xc8,
	0x71, 0x3d, 0x21, 0xf9, 0xc5, 0x2b, 0xf0, 0xce, 0xe5, 0x42, 0x69, 0xe9, 0x7b, 0xa4, 0x3c, 0xca,
	0x95, 0xcc, 0x5a, 0xad, 0xd5, 0xe0, 0xaa, 0xd7, 0x32, 0xfb, 0xe6, 0x37, 0xb6, 0x57, 0x45, 0x2e,
	0x5f, 0x15, 0x6d, 0xd1, 0x87, 0xe6, 0x58, 0xe9, 0x45, 0xb2, 0x94, 0x73, 0x7f, 0x7f, 0x5b, 0x8b,
	0xbc, 0xd5, 0x16, 0x1a, 0xd1, 0x83, 0xfa, 0x55, 0xb2, 0x54, 0x99, 0x7f, 0xb0, 0x2d, 0x46, 0xd2,
	0x8a, 0x59, 0x20, 0x9e, 0x83, 0x3b, 0x96, 0xb3, 0xcc, 0x3f, 0x24, 0xe1, 0x7f, 0xc5, 0xaa, 0x72,
	0x96, 0x17, 0x80, 0xee, 0xe0, 0x1d, 0x40, 0xb9, 0x55, 0x7b, 0x90, 0x37, 0x72, 0xa1, 0x8a, 0x36,
	0x33, 0xc4, 0x86, 0x12, 0xcd, 0x5d, 0x26, 0x3b, 0xf8, 0x0e, 0xed, 0xea, 0xe6, 0x6d, 0xf4, 0x97,
	0x34, 0x5d, 0x54, 0xa2, 0x11, 0x62, 0x34, 0xd1, 0x18, 0xed, 0x84, 0x64, 0x8b, 0x00, 0xda, 0x91,
	0x32, 0x97, 0x59, 0xba, 0x50, 0x46, 0x27, 0x13, 0xdb, 0xb9, 0x2d, 0x4e, 0x9c, 0x40, 0xab, 0x14,
	0xb8, 0x24, 0x28, 0x09, 0xac, 0xbd, 0x3c, 0xd0, 0x27, 0xd6, 0xfe, 0x15, 0xda, 0xd1, 0x9d, 0x8c,
	0x95, 0xde, 0x8a, 0xfe, 0xa4, 0xb4, 0xa9, 0x44, 0x23, 0xc4, 0x68, 0xa2, 0x6d, 0x34, 0x71, 0xac,
	0x1e, 0x69, 0x39, 0xab, 0x0c, 0x1c, 0x42, 0x54, 0x13, 0xed, 0xb2, 0x1a, 0xed, 0xe0, 0x16, 0xa0,
	0x6c, 0x26, 0x0e, 0x12, 0x67, 0xa6, 0x44, 0x95, 0x41, 0xaa, 0xd6, 0x13, 0x5a, 0x0d, 0x0e, 0xd2,
	0x85, 0x4e, 0x62, 0xca, 0x5e, 0xe9, 0x37, 0x72, 0x79, 0x1f, 0xd1, 0x0e, 0x7e, 0xee, 0x02, 0x94,
	0x43, 0xf0, 0x6f, 0x92, 0x88, 0x2e, 0xec, 0xe1, 0x4d, 0x4c, 0x37, 0xcb, 0x8d, 0xd4, 0xb1, 0xdd,
	0x7a, 0x95, 0xc2, 0xbb, 0x5b, 0xb8, 0xb9, 0x5f, 0xcd, 0xaa, 0x2f, 0x52, 0x26, 0x5a, 0x29, 0x15,
	0xd3, 0x6d, 0xe3, 0x7b, 0x4d, 0x18, 0xef, 0x35, 0x3b, 0x3c, 0x9a, 0x10, 0x06, 0x36, 0x62, 0x94,
	0xdc, 0xab, 0x98, 0x6e, 0x12, 0x47, 0x10, 0xc6, 0x08, 0x76, 0x34, 0xf9, 0x25, 0x60, 0x96, 0x1b,
	0x33, 0xd0, 0xe9, 0xca, 0x6f, 0x15, 0x8d, 0x41, 0x88, 0x8d, 0x21, 0x1a, 0x88, 0x26, 0xdb, 0xaa,
	0x23, 0x93, 0xae, 0xfc, 0xbd, 0x42, 0x8d, 0x10, 0xd5, 0x44, 0xb7, 0x59, 0x4d, 0xdc, 0x09, 0xb4,
	0x22, 0x65, 0xde, 0xaf, 0xa7, 0x53, 0xa5, 0xe9, 0xb2, 0x36, 0xc3, 0x92, 0x10, 0xff, 0x83, 0x67,
	0x5d, 0x78, 0x35, 0xdd, 0xd0, 0xa2, 0xe0, 0x97, 0x03, 0x50, 0x1e, 0xa4, 0xdd, 0xd2, 0xe7, 0x24,
	0x36, 0x77, 0x95, 0xc7, 0x8d, 0x30, 0x6e, 0x89, 0x1d, 0x35, 0x5a, 0x81, 0x81, 0x4d, 0xfb, 0x41,
	0x25, 0xb3, 0x3b, 0x63, 0x8f, 0xbc, 0x24, 0x30, 0xad, 0x75, 0xb9, 0x9c, 0xd6, 0xf2, 0x76, 0x42,
	0x93, 0xf9, 0xdc, 0x9e, 0x75, 0x0e, 0x69, 0x42, 0x91, 0xf6, 0xec, 0x84, 0x26, 0xf3, 0x79, 0xf0,
	0xdb, 0x81, 0x56, 0xf1, 0x32, 0x3c, 0x71, 0x78, 0xec, 0x50, 0xac, 0xb5, 0x34, 0x49, 0xba, 0xac,
	0x0e, 0x85, 0xa5, 0x68, 0x28, 0x72, 0xb7, 0x4b, 0xfd, 0x2d, 0xf0, 0xd3, 0xea, 0x7c, 0xf9, 0x86,
	0x9f, 0x73, 0xd1, 0x04, 0x77, 0x3c, 0x0c, 0xaf, 0x3b, 0x3b, 0x68, 0x5d, 0x5d, 0xde, 0x0c, 0x3b,
	0x0e, 0x71, 0x67, 0x17, 0x51, 0xa7, 0x26, 0x5a, 0x50, 0x1f, 0x84, 0x67, 0xa3, 0x71, 0xa7, 0x73,
	0xfa, 0xc3, 0x01, 0x6f, 0x44, 0xc5, 0x8b, 0xd7, 0xd0, 0xe2, 0xca, 0xd7, 0x5a, 0x89, 0x83, 0xe2,
	0xf5, 0x26, 0xea, 0xb8, 0xc0, 0xfc, 0xdf, 0x05, 0x3b, 0xa2, 0x0f, 0x8d, 0x41, 0x92, 0xad, 0xe6,
	0xf2, 0x41, 0x1c, 0xe6, 0x4e, 0xfb, 0x7d, 0x3d, 0x56, 0xf7, 0x1c, 0x7c, 0x99, 0x3f, 0xae, 0x95,
	0x7e, 0x10, 0xfb, 0xb9, 0x93, 0xfe, 0xd0, 0xc7, 0xda, 0x5b, 0x8f, 0x7e, 0xdb, 0xb7, 0x7f, 0x03,
	0x00, 0x00, 0xff, 0xff, 0x8d, 0x6b, 0xe2, 0x86, 0x7d, 0x07, 0x00, 0x00,
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
