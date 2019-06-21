// Code generated by protoc-gen-go. DO NOT EDIT.
// source: facade.proto

package facade

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

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
	return fileDescriptor_facade_e4f1de882b897b17, []int{0}
}

type Status struct {
	Success              bool     `protobuf:"varint,1,opt,name=Success,proto3" json:"Success,omitempty"`
	Error                string   `protobuf:"bytes,2,opt,name=Error,proto3" json:"Error,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Status) Reset()         { *m = Status{} }
func (m *Status) String() string { return proto.CompactTextString(m) }
func (*Status) ProtoMessage()    {}
func (*Status) Descriptor() ([]byte, []int) {
	return fileDescriptor_facade_e4f1de882b897b17, []int{0}
}
func (m *Status) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Status.Unmarshal(m, b)
}
func (m *Status) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Status.Marshal(b, m, deterministic)
}
func (dst *Status) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Status.Merge(dst, src)
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
	return fileDescriptor_facade_e4f1de882b897b17, []int{1}
}
func (m *RawText) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RawText.Unmarshal(m, b)
}
func (m *RawText) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RawText.Marshal(b, m, deterministic)
}
func (dst *RawText) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RawText.Merge(dst, src)
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
	Grid                 *GridConfig   `protobuf:"bytes,8,opt,name=Grid,proto3" json:"Grid,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *Config) Reset()         { *m = Config{} }
func (m *Config) String() string { return proto.CompactTextString(m) }
func (*Config) ProtoMessage()    {}
func (*Config) Descriptor() ([]byte, []int) {
	return fileDescriptor_facade_e4f1de882b897b17, []int{2}
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
	return fileDescriptor_facade_e4f1de882b897b17, []int{3}
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
	return fileDescriptor_facade_e4f1de882b897b17, []int{4}
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
	return fileDescriptor_facade_e4f1de882b897b17, []int{5}
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
	return fileDescriptor_facade_e4f1de882b897b17, []int{6}
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
	proto.RegisterType((*Status)(nil), "facade.Status")
	proto.RegisterType((*RawText)(nil), "facade.RawText")
	proto.RegisterType((*Config)(nil), "facade.Config")
	proto.RegisterType((*FontConfig)(nil), "facade.FontConfig")
	proto.RegisterType((*CameraConfig)(nil), "facade.CameraConfig")
	proto.RegisterType((*MaskConfig)(nil), "facade.MaskConfig")
	proto.RegisterType((*GridConfig)(nil), "facade.GridConfig")
	proto.RegisterEnum("facade.Mode", Mode_name, Mode_value)
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

// FacadeServer is the server API for Facade service.
type FacadeServer interface {
	Configure(context.Context, *Config) (*Status, error)
	Display(Facade_DisplayServer) error
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

var _Facade_serviceDesc = grpc.ServiceDesc{
	ServiceName: "facade.Facade",
	HandlerType: (*FacadeServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Configure",
			Handler:    _Facade_Configure_Handler,
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

func init() { proto.RegisterFile("facade.proto", fileDescriptor_facade_e4f1de882b897b17) }

var fileDescriptor_facade_e4f1de882b897b17 = []byte{
	// 583 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x54, 0x4f, 0x6f, 0xd3, 0x30,
	0x14, 0x5f, 0xb6, 0x2c, 0x4d, 0xdf, 0xc2, 0x56, 0xac, 0x09, 0x59, 0x1b, 0x87, 0x28, 0x07, 0x54,
	0x21, 0xd8, 0x61, 0x5c, 0x10, 0x37, 0x58, 0xe9, 0xd8, 0x61, 0x1c, 0x9c, 0x09, 0x24, 0x6e, 0x5e,
	0xe3, 0x76, 0x11, 0x69, 0x33, 0x39, 0xae, 0x0a, 0x17, 0x3e, 0x09, 0x9f, 0x8e, 0x4f, 0x82, 0xde,
	0x7b, 0x4e, 0xd2, 0x8a, 0x13, 0x37, 0xff, 0xfe, 0x3c, 0xeb, 0xe7, 0xe7, 0x67, 0x43, 0x32, 0xd7,
	0x33, 0x5d, 0x98, 0x8b, 0x47, 0x5b, 0xbb, 0x5a, 0x44, 0x8c, 0xb2, 0xb7, 0x10, 0xe5, 0x4e, 0xbb,
	0x75, 0x23, 0x24, 0x0c, 0xf2, 0xf5, 0x6c, 0x66, 0x9a, 0x46, 0x06, 0x69, 0x30, 0x8e, 0x55, 0x0b,
	0xc5, 0x29, 0x1c, 0x7e, 0xb4, 0xb6, 0xb6, 0x72, 0x3f, 0x0d, 0xc6, 0x43, 0xc5, 0x20, 0x3b, 0x87,
	0x81, 0xd2, 0x9b, 0x3b, 0xf3, 0xc3, 0x89, 0x11, 0x1c, 0x28, 0xbd, 0xa1, 0xb2, 0x44, 0xe1, 0x32,
	0xfb, 0xbd, 0x0f, 0xd1, 0x55, 0xbd, 0x9a, 0x97, 0x0b, 0x71, 0x06, 0x71, 0x6e, 0xdc, 0xc4, 0xdc,
	0xaf, 0x17, 0x7e, 0xe3, 0x0e, 0xe3, 0xce, 0x2c, 0xec, 0x93, 0xc0, 0x80, 0x92, 0x18, 0x77, 0x5b,
	0x17, 0x46, 0x1e, 0xf8, 0x24, 0x0c, 0x45, 0x0a, 0x21, 0xd1, 0x61, 0x1a, 0x8c, 0x8f, 0x2f, 0x93,
	0x0b, 0x7f, 0x24, 0xe4, 0x14, 0x29, 0xe2, 0x05, 0x84, 0xd3, 0x7a, 0xe5, 0xe4, 0x61, 0x1a, 0x8c,
	0x8f, 0x2e, 0x45, 0xeb, 0x40, 0x8e, 0xf3, 0x28, 0xd2, 0xc5, 0x2b, 0x88, 0xae, 0xf4, 0xd2, 0x58,
	0x2d, 0x23, 0x72, 0x9e, 0xb6, 0x4e, 0x66, 0xbd, 0xd7, 0x7b, 0x70, 0xd7, 0x5b, 0xdd, 0x7c, 0x97,
	0x83, 0xdd, 0x5d, 0x91, 0x6b, 0x77, 0xc5, 0x35, 0xfa, 0xae, 0x6d, 0x59, 0xc8, 0x78, 0xd7, 0x87,
	0x5c, 0xeb, 0xc3, 0x75, 0xf6, 0x0e, 0xa0, 0x4f, 0xe4, 0xcf, 0xfb, 0x59, 0x2f, 0x4d, 0xd7, 0x79,
	0x86, 0x42, 0x40, 0x48, 0x34, 0x37, 0x9e, 0xd6, 0xd9, 0x2f, 0x48, 0xb6, 0x33, 0xfa, 0xea, 0x6f,
	0x75, 0xbd, 0xdc, 0xaa, 0x46, 0x88, 0xd5, 0x44, 0x63, 0x75, 0xa0, 0x68, 0x2d, 0x32, 0x48, 0x72,
	0xe3, 0x6e, 0x9a, 0x7a, 0x69, 0x9c, 0x2d, 0x67, 0xbe, 0xc1, 0x3b, 0x9c, 0x78, 0x0e, 0xc3, 0xde,
	0x10, 0x92, 0xa1, 0x27, 0x30, 0x7b, 0x7f, 0xee, 0xff, 0xcc, 0xfe, 0xe7, 0x00, 0xa0, 0x6f, 0x86,
	0x1f, 0x8d, 0xaf, 0x65, 0xe1, 0x1e, 0xb6, 0x46, 0x83, 0x30, 0x8e, 0x06, 0x0b, 0x58, 0x1f, 0x2a,
	0x06, 0x18, 0x2d, 0x37, 0xee, 0x93, 0x29, 0x17, 0x0f, 0xce, 0x67, 0xef, 0x09, 0xf1, 0x0c, 0x22,
	0x2f, 0x85, 0x54, 0xe4, 0x91, 0x48, 0xe1, 0x08, 0x47, 0xae, 0xde, 0xac, 0x36, 0xda, 0x16, 0x34,
	0x1b, 0xb1, 0xda, 0xa6, 0x30, 0x49, 0x27, 0x47, 0x9c, 0x64, 0x5b, 0xcb, 0x8d, 0xcb, 0x1f, 0x8d,
	0x29, 0x68, 0x00, 0x38, 0x25, 0x61, 0x4c, 0xc9, 0x42, 0x4c, 0x3d, 0x66, 0xe0, 0x53, 0x7e, 0x58,
	0xcf, 0xe7, 0xc6, 0xca, 0x61, 0x97, 0x92, 0x09, 0x4c, 0xe9, 0x25, 0xe0, 0x94, 0x9e, 0xe7, 0x94,
	0x77, 0xc6, 0x2e, 0xcb, 0x95, 0xae, 0xe4, 0x51, 0x97, 0xb2, 0xa5, 0x30, 0x49, 0x27, 0x27, 0x9c,
	0xa4, 0xd3, 0xf8, 0x22, 0xbe, 0x18, 0xeb, 0xe4, 0x93, 0xee, 0x22, 0x10, 0xe2, 0x45, 0x10, 0x7d,
	0xcc, 0x17, 0x41, 0x1c, 0xbb, 0xa7, 0x56, 0x2f, 0xe4, 0x49, 0xe7, 0x46, 0x88, 0x6e, 0xa2, 0x47,
	0xec, 0x26, 0xce, 0xbb, 0xcb, 0xaa, 0x92, 0x4f, 0x7b, 0x77, 0x59, 0x55, 0xe4, 0x46, 0x5a, 0x78,
	0x77, 0x59, 0x55, 0x2f, 0xcf, 0xf9, 0x91, 0x8a, 0x18, 0xc2, 0x6b, 0x75, 0x33, 0x19, 0xed, 0x89,
	0x21, 0x1c, 0x4e, 0xd4, 0xfb, 0xe9, 0xdd, 0x68, 0x74, 0xb9, 0x80, 0x68, 0x4a, 0x8f, 0x42, 0xbc,
	0x86, 0x21, 0x8f, 0xc1, 0xda, 0x1a, 0x71, 0xdc, 0x3d, 0x3f, 0xa2, 0xce, 0x3a, 0xcc, 0x9f, 0x53,
	0xb6, 0x27, 0x2e, 0x60, 0x30, 0x29, 0x9b, 0xc7, 0x4a, 0xff, 0x14, 0x27, 0xad, 0xe8, 0xff, 0x9f,
	0x7f, 0xdd, 0xe3, 0xe0, 0x3e, 0xa2, 0x7f, 0xee, 0xcd, 0xdf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x8c,
	0xd7, 0x0a, 0x25, 0xf7, 0x04, 0x00, 0x00,
}
