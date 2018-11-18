// Code generated by protoc-gen-go. DO NOT EDIT.
// source: inapi/pod.proto

package inapi

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

type PbPodSpecBoxImageDriver int32

const (
	PbPodSpecBoxImageDriver_Unknown PbPodSpecBoxImageDriver = 0
	PbPodSpecBoxImageDriver_Docker  PbPodSpecBoxImageDriver = 1
	PbPodSpecBoxImageDriver_Rkt     PbPodSpecBoxImageDriver = 2
	PbPodSpecBoxImageDriver_Pouch   PbPodSpecBoxImageDriver = 3
)

var PbPodSpecBoxImageDriver_name = map[int32]string{
	0: "Unknown",
	1: "Docker",
	2: "Rkt",
	3: "Pouch",
}
var PbPodSpecBoxImageDriver_value = map[string]int32{
	"Unknown": 0,
	"Docker":  1,
	"Rkt":     2,
	"Pouch":   3,
}

func (x PbPodSpecBoxImageDriver) String() string {
	return proto.EnumName(PbPodSpecBoxImageDriver_name, int32(x))
}
func (PbPodSpecBoxImageDriver) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_pod_2de871fe32e65100, []int{0}
}

type PbPodRepStatus struct {
	Id                   string             `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Rep                  uint32             `protobuf:"varint,2,opt,name=rep" json:"rep,omitempty"`
	Action               uint32             `protobuf:"varint,3,opt,name=action" json:"action,omitempty"`
	Node                 string             `protobuf:"bytes,4,opt,name=node" json:"node,omitempty"`
	Box                  *PbPodBoxStatus    `protobuf:"bytes,10,opt,name=box" json:"box,omitempty"`
	OpLog                *PbOpLogSets       `protobuf:"bytes,6,opt,name=op_log,json=opLog" json:"op_log,omitempty"`
	Stats                *PbStatsSampleFeed `protobuf:"bytes,7,opt,name=stats" json:"stats,omitempty"`
	Updated              uint32             `protobuf:"varint,8,opt,name=updated" json:"updated,omitempty"`
	Volumes              []*PbVolumeStatus  `protobuf:"bytes,9,rep,name=volumes" json:"volumes,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *PbPodRepStatus) Reset()         { *m = PbPodRepStatus{} }
func (m *PbPodRepStatus) String() string { return proto.CompactTextString(m) }
func (*PbPodRepStatus) ProtoMessage()    {}
func (*PbPodRepStatus) Descriptor() ([]byte, []int) {
	return fileDescriptor_pod_2de871fe32e65100, []int{0}
}
func (m *PbPodRepStatus) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PbPodRepStatus.Unmarshal(m, b)
}
func (m *PbPodRepStatus) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PbPodRepStatus.Marshal(b, m, deterministic)
}
func (dst *PbPodRepStatus) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PbPodRepStatus.Merge(dst, src)
}
func (m *PbPodRepStatus) XXX_Size() int {
	return xxx_messageInfo_PbPodRepStatus.Size(m)
}
func (m *PbPodRepStatus) XXX_DiscardUnknown() {
	xxx_messageInfo_PbPodRepStatus.DiscardUnknown(m)
}

var xxx_messageInfo_PbPodRepStatus proto.InternalMessageInfo

func (m *PbPodRepStatus) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *PbPodRepStatus) GetRep() uint32 {
	if m != nil {
		return m.Rep
	}
	return 0
}

func (m *PbPodRepStatus) GetAction() uint32 {
	if m != nil {
		return m.Action
	}
	return 0
}

func (m *PbPodRepStatus) GetNode() string {
	if m != nil {
		return m.Node
	}
	return ""
}

func (m *PbPodRepStatus) GetBox() *PbPodBoxStatus {
	if m != nil {
		return m.Box
	}
	return nil
}

func (m *PbPodRepStatus) GetOpLog() *PbOpLogSets {
	if m != nil {
		return m.OpLog
	}
	return nil
}

func (m *PbPodRepStatus) GetStats() *PbStatsSampleFeed {
	if m != nil {
		return m.Stats
	}
	return nil
}

func (m *PbPodRepStatus) GetUpdated() uint32 {
	if m != nil {
		return m.Updated
	}
	return 0
}

func (m *PbPodRepStatus) GetVolumes() []*PbVolumeStatus {
	if m != nil {
		return m.Volumes
	}
	return nil
}

type PbVolumeMount struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	ReadOnly             bool     `protobuf:"varint,2,opt,name=read_only,json=readOnly" json:"read_only,omitempty"`
	MountPath            string   `protobuf:"bytes,3,opt,name=mount_path,json=mountPath" json:"mount_path,omitempty"`
	HostDir              string   `protobuf:"bytes,4,opt,name=host_dir,json=hostDir" json:"host_dir,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PbVolumeMount) Reset()         { *m = PbVolumeMount{} }
func (m *PbVolumeMount) String() string { return proto.CompactTextString(m) }
func (*PbVolumeMount) ProtoMessage()    {}
func (*PbVolumeMount) Descriptor() ([]byte, []int) {
	return fileDescriptor_pod_2de871fe32e65100, []int{1}
}
func (m *PbVolumeMount) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PbVolumeMount.Unmarshal(m, b)
}
func (m *PbVolumeMount) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PbVolumeMount.Marshal(b, m, deterministic)
}
func (dst *PbVolumeMount) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PbVolumeMount.Merge(dst, src)
}
func (m *PbVolumeMount) XXX_Size() int {
	return xxx_messageInfo_PbVolumeMount.Size(m)
}
func (m *PbVolumeMount) XXX_DiscardUnknown() {
	xxx_messageInfo_PbVolumeMount.DiscardUnknown(m)
}

var xxx_messageInfo_PbVolumeMount proto.InternalMessageInfo

func (m *PbVolumeMount) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *PbVolumeMount) GetReadOnly() bool {
	if m != nil {
		return m.ReadOnly
	}
	return false
}

func (m *PbVolumeMount) GetMountPath() string {
	if m != nil {
		return m.MountPath
	}
	return ""
}

func (m *PbVolumeMount) GetHostDir() string {
	if m != nil {
		return m.HostDir
	}
	return ""
}

type PbVolumeStatus struct {
	MountPath            string   `protobuf:"bytes,1,opt,name=mount_path,json=mountPath" json:"mount_path,omitempty"`
	Used                 int64    `protobuf:"varint,2,opt,name=used" json:"used,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PbVolumeStatus) Reset()         { *m = PbVolumeStatus{} }
func (m *PbVolumeStatus) String() string { return proto.CompactTextString(m) }
func (*PbVolumeStatus) ProtoMessage()    {}
func (*PbVolumeStatus) Descriptor() ([]byte, []int) {
	return fileDescriptor_pod_2de871fe32e65100, []int{2}
}
func (m *PbVolumeStatus) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PbVolumeStatus.Unmarshal(m, b)
}
func (m *PbVolumeStatus) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PbVolumeStatus.Marshal(b, m, deterministic)
}
func (dst *PbVolumeStatus) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PbVolumeStatus.Merge(dst, src)
}
func (m *PbVolumeStatus) XXX_Size() int {
	return xxx_messageInfo_PbVolumeStatus.Size(m)
}
func (m *PbVolumeStatus) XXX_DiscardUnknown() {
	xxx_messageInfo_PbVolumeStatus.DiscardUnknown(m)
}

var xxx_messageInfo_PbVolumeStatus proto.InternalMessageInfo

func (m *PbVolumeStatus) GetMountPath() string {
	if m != nil {
		return m.MountPath
	}
	return ""
}

func (m *PbVolumeStatus) GetUsed() int64 {
	if m != nil {
		return m.Used
	}
	return 0
}

type PbServicePort struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	BoxPort              uint32   `protobuf:"varint,2,opt,name=box_port,json=boxPort" json:"box_port,omitempty"`
	HostPort             uint32   `protobuf:"varint,3,opt,name=host_port,json=hostPort" json:"host_port,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PbServicePort) Reset()         { *m = PbServicePort{} }
func (m *PbServicePort) String() string { return proto.CompactTextString(m) }
func (*PbServicePort) ProtoMessage()    {}
func (*PbServicePort) Descriptor() ([]byte, []int) {
	return fileDescriptor_pod_2de871fe32e65100, []int{3}
}
func (m *PbServicePort) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PbServicePort.Unmarshal(m, b)
}
func (m *PbServicePort) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PbServicePort.Marshal(b, m, deterministic)
}
func (dst *PbServicePort) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PbServicePort.Merge(dst, src)
}
func (m *PbServicePort) XXX_Size() int {
	return xxx_messageInfo_PbServicePort.Size(m)
}
func (m *PbServicePort) XXX_DiscardUnknown() {
	xxx_messageInfo_PbServicePort.DiscardUnknown(m)
}

var xxx_messageInfo_PbServicePort proto.InternalMessageInfo

func (m *PbServicePort) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *PbServicePort) GetBoxPort() uint32 {
	if m != nil {
		return m.BoxPort
	}
	return 0
}

func (m *PbServicePort) GetHostPort() uint32 {
	if m != nil {
		return m.HostPort
	}
	return 0
}

type PbPodBoxStatusExecutor struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Phase                string   `protobuf:"bytes,2,opt,name=phase" json:"phase,omitempty"`
	Retry                uint32   `protobuf:"varint,3,opt,name=retry" json:"retry,omitempty"`
	ErrorCode            uint32   `protobuf:"varint,4,opt,name=error_code,json=errorCode" json:"error_code,omitempty"`
	ErrorMessage         string   `protobuf:"bytes,5,opt,name=error_message,json=errorMessage" json:"error_message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PbPodBoxStatusExecutor) Reset()         { *m = PbPodBoxStatusExecutor{} }
func (m *PbPodBoxStatusExecutor) String() string { return proto.CompactTextString(m) }
func (*PbPodBoxStatusExecutor) ProtoMessage()    {}
func (*PbPodBoxStatusExecutor) Descriptor() ([]byte, []int) {
	return fileDescriptor_pod_2de871fe32e65100, []int{4}
}
func (m *PbPodBoxStatusExecutor) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PbPodBoxStatusExecutor.Unmarshal(m, b)
}
func (m *PbPodBoxStatusExecutor) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PbPodBoxStatusExecutor.Marshal(b, m, deterministic)
}
func (dst *PbPodBoxStatusExecutor) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PbPodBoxStatusExecutor.Merge(dst, src)
}
func (m *PbPodBoxStatusExecutor) XXX_Size() int {
	return xxx_messageInfo_PbPodBoxStatusExecutor.Size(m)
}
func (m *PbPodBoxStatusExecutor) XXX_DiscardUnknown() {
	xxx_messageInfo_PbPodBoxStatusExecutor.DiscardUnknown(m)
}

var xxx_messageInfo_PbPodBoxStatusExecutor proto.InternalMessageInfo

func (m *PbPodBoxStatusExecutor) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *PbPodBoxStatusExecutor) GetPhase() string {
	if m != nil {
		return m.Phase
	}
	return ""
}

func (m *PbPodBoxStatusExecutor) GetRetry() uint32 {
	if m != nil {
		return m.Retry
	}
	return 0
}

func (m *PbPodBoxStatusExecutor) GetErrorCode() uint32 {
	if m != nil {
		return m.ErrorCode
	}
	return 0
}

func (m *PbPodBoxStatusExecutor) GetErrorMessage() string {
	if m != nil {
		return m.ErrorMessage
	}
	return ""
}

// PodBoxStatus represents the current information about a box
type PbPodBoxStatus struct {
	Name                 string                    `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	ImageDriver          PbPodSpecBoxImageDriver   `protobuf:"varint,2,opt,name=image_driver,json=imageDriver,enum=inapi.PbPodSpecBoxImageDriver" json:"image_driver,omitempty"`
	ImageOptions         []*Label                  `protobuf:"bytes,3,rep,name=image_options,json=imageOptions" json:"image_options,omitempty"`
	ResCpuLimit          int64                     `protobuf:"varint,4,opt,name=res_cpu_limit,json=resCpuLimit" json:"res_cpu_limit,omitempty"`
	ResMemLimit          int64                     `protobuf:"varint,5,opt,name=res_mem_limit,json=resMemLimit" json:"res_mem_limit,omitempty"`
	Mounts               []*PbVolumeMount          `protobuf:"bytes,6,rep,name=mounts" json:"mounts,omitempty"`
	Ports                []*PbServicePort          `protobuf:"bytes,7,rep,name=ports" json:"ports,omitempty"`
	Command              []string                  `protobuf:"bytes,8,rep,name=command" json:"command,omitempty"`
	Executors            []*PbPodBoxStatusExecutor `protobuf:"bytes,9,rep,name=executors" json:"executors,omitempty"`
	Action               uint32                    `protobuf:"varint,10,opt,name=action" json:"action,omitempty"`
	Started              uint32                    `protobuf:"varint,11,opt,name=started" json:"started,omitempty"`
	Updated              uint32                    `protobuf:"varint,12,opt,name=updated" json:"updated,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                  `json:"-"`
	XXX_unrecognized     []byte                    `json:"-"`
	XXX_sizecache        int32                     `json:"-"`
}

func (m *PbPodBoxStatus) Reset()         { *m = PbPodBoxStatus{} }
func (m *PbPodBoxStatus) String() string { return proto.CompactTextString(m) }
func (*PbPodBoxStatus) ProtoMessage()    {}
func (*PbPodBoxStatus) Descriptor() ([]byte, []int) {
	return fileDescriptor_pod_2de871fe32e65100, []int{5}
}
func (m *PbPodBoxStatus) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PbPodBoxStatus.Unmarshal(m, b)
}
func (m *PbPodBoxStatus) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PbPodBoxStatus.Marshal(b, m, deterministic)
}
func (dst *PbPodBoxStatus) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PbPodBoxStatus.Merge(dst, src)
}
func (m *PbPodBoxStatus) XXX_Size() int {
	return xxx_messageInfo_PbPodBoxStatus.Size(m)
}
func (m *PbPodBoxStatus) XXX_DiscardUnknown() {
	xxx_messageInfo_PbPodBoxStatus.DiscardUnknown(m)
}

var xxx_messageInfo_PbPodBoxStatus proto.InternalMessageInfo

func (m *PbPodBoxStatus) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *PbPodBoxStatus) GetImageDriver() PbPodSpecBoxImageDriver {
	if m != nil {
		return m.ImageDriver
	}
	return PbPodSpecBoxImageDriver_Unknown
}

func (m *PbPodBoxStatus) GetImageOptions() []*Label {
	if m != nil {
		return m.ImageOptions
	}
	return nil
}

func (m *PbPodBoxStatus) GetResCpuLimit() int64 {
	if m != nil {
		return m.ResCpuLimit
	}
	return 0
}

func (m *PbPodBoxStatus) GetResMemLimit() int64 {
	if m != nil {
		return m.ResMemLimit
	}
	return 0
}

func (m *PbPodBoxStatus) GetMounts() []*PbVolumeMount {
	if m != nil {
		return m.Mounts
	}
	return nil
}

func (m *PbPodBoxStatus) GetPorts() []*PbServicePort {
	if m != nil {
		return m.Ports
	}
	return nil
}

func (m *PbPodBoxStatus) GetCommand() []string {
	if m != nil {
		return m.Command
	}
	return nil
}

func (m *PbPodBoxStatus) GetExecutors() []*PbPodBoxStatusExecutor {
	if m != nil {
		return m.Executors
	}
	return nil
}

func (m *PbPodBoxStatus) GetAction() uint32 {
	if m != nil {
		return m.Action
	}
	return 0
}

func (m *PbPodBoxStatus) GetStarted() uint32 {
	if m != nil {
		return m.Started
	}
	return 0
}

func (m *PbPodBoxStatus) GetUpdated() uint32 {
	if m != nil {
		return m.Updated
	}
	return 0
}

func init() {
	proto.RegisterType((*PbPodRepStatus)(nil), "inapi.PbPodRepStatus")
	proto.RegisterType((*PbVolumeMount)(nil), "inapi.PbVolumeMount")
	proto.RegisterType((*PbVolumeStatus)(nil), "inapi.PbVolumeStatus")
	proto.RegisterType((*PbServicePort)(nil), "inapi.PbServicePort")
	proto.RegisterType((*PbPodBoxStatusExecutor)(nil), "inapi.PbPodBoxStatusExecutor")
	proto.RegisterType((*PbPodBoxStatus)(nil), "inapi.PbPodBoxStatus")
	proto.RegisterEnum("inapi.PbPodSpecBoxImageDriver", PbPodSpecBoxImageDriver_name, PbPodSpecBoxImageDriver_value)
}

func init() { proto.RegisterFile("inapi/pod.proto", fileDescriptor_pod_2de871fe32e65100) }

var fileDescriptor_pod_2de871fe32e65100 = []byte{
	// 714 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x54, 0xc1, 0x6e, 0xf3, 0x44,
	0x10, 0x26, 0x71, 0x1d, 0xc7, 0x93, 0xa4, 0x84, 0x55, 0xf9, 0x71, 0x5b, 0x15, 0x45, 0xe1, 0x40,
	0xa9, 0x50, 0x2a, 0xca, 0x91, 0x13, 0x6d, 0xa9, 0x84, 0xd4, 0xaa, 0xd1, 0x46, 0x70, 0xe1, 0x60,
	0x6d, 0xec, 0x51, 0x62, 0x35, 0xf6, 0xae, 0x76, 0xd7, 0x25, 0xe5, 0x59, 0x78, 0x05, 0x9e, 0x89,
	0x57, 0x41, 0x3b, 0xeb, 0xa4, 0x49, 0x15, 0x6e, 0x3b, 0xdf, 0x7e, 0x9e, 0xf9, 0x66, 0xe7, 0x1b,
	0xc3, 0xe7, 0x45, 0x25, 0x54, 0x71, 0xad, 0x64, 0x3e, 0x51, 0x5a, 0x5a, 0xc9, 0x42, 0x02, 0xce,
	0x86, 0x1e, 0x9f, 0x0b, 0x83, 0xfe, 0xe2, 0xec, 0xc4, 0x23, 0x52, 0xa1, 0x16, 0x56, 0xea, 0x06,
	0xfd, 0xc2, 0xa3, 0xc6, 0x0a, 0x6b, 0x3c, 0x34, 0xfe, 0xa7, 0x0d, 0xc7, 0xd3, 0xf9, 0x54, 0xe6,
	0x1c, 0xd5, 0xcc, 0x0a, 0x5b, 0x1b, 0x76, 0x0c, 0xed, 0x22, 0x4f, 0x5a, 0xa3, 0xd6, 0x65, 0xcc,
	0xdb, 0x45, 0xce, 0x86, 0x10, 0x68, 0x54, 0x49, 0x7b, 0xd4, 0xba, 0x1c, 0x70, 0x77, 0x64, 0x9f,
	0xa0, 0x23, 0x32, 0x5b, 0xc8, 0x2a, 0x09, 0x08, 0x6c, 0x22, 0xc6, 0xe0, 0xa8, 0x92, 0x39, 0x26,
	0x47, 0xf4, 0x2d, 0x9d, 0xd9, 0xb7, 0x10, 0xcc, 0xe5, 0x3a, 0x81, 0x51, 0xeb, 0xb2, 0x77, 0xf3,
	0xe5, 0x84, 0x14, 0x4c, 0xa8, 0xe2, 0xad, 0x5c, 0xfb, 0x8a, 0xdc, 0x31, 0xd8, 0x77, 0xd0, 0x91,
	0x2a, 0x5d, 0xc9, 0x45, 0xd2, 0x21, 0x2e, 0xdb, 0x72, 0x9f, 0xd5, 0xa3, 0x5c, 0xcc, 0xd0, 0x1a,
	0x1e, 0x4a, 0x77, 0x64, 0x13, 0x08, 0xa9, 0x87, 0x24, 0x22, 0x66, 0xb2, 0x65, 0xba, 0x7c, 0x66,
	0x26, 0x4a, 0xb5, 0xc2, 0x07, 0xc4, 0x9c, 0x7b, 0x1a, 0x4b, 0x20, 0xaa, 0x55, 0x2e, 0x2c, 0xe6,
	0x49, 0x97, 0x04, 0x6f, 0x42, 0x76, 0x0d, 0xd1, 0xab, 0x5c, 0xd5, 0x25, 0x9a, 0x24, 0x1e, 0x05,
	0x7b, 0x0a, 0x7f, 0x27, 0xbc, 0x51, 0xb8, 0x61, 0x8d, 0xff, 0x82, 0xc1, 0xe6, 0xea, 0x49, 0xd6,
	0x95, 0xa5, 0x9e, 0x45, 0x89, 0xcd, 0x7b, 0xd1, 0x99, 0x9d, 0x43, 0xac, 0x51, 0xe4, 0xa9, 0xac,
	0x56, 0x6f, 0xf4, 0x6e, 0x5d, 0xde, 0x75, 0xc0, 0x73, 0xb5, 0x7a, 0x63, 0x17, 0x00, 0xa5, 0xfb,
	0x32, 0x55, 0xc2, 0x2e, 0xe9, 0x01, 0x63, 0x1e, 0x13, 0x32, 0x15, 0x76, 0xc9, 0x4e, 0xa1, 0xbb,
	0x94, 0xc6, 0xa6, 0x79, 0xa1, 0x9b, 0x77, 0x8c, 0x5c, 0x7c, 0x5f, 0xe8, 0xf1, 0x9d, 0x1b, 0xd5,
	0xae, 0xac, 0x0f, 0xb9, 0x5a, 0x1f, 0x73, 0x31, 0x38, 0xaa, 0x0d, 0xe6, 0x24, 0x21, 0xe0, 0x74,
	0x1e, 0xff, 0xe1, 0x1a, 0x98, 0xa1, 0x7e, 0x2d, 0x32, 0x9c, 0x4a, 0x7d, 0xb8, 0x81, 0x53, 0xe8,
	0xce, 0xe5, 0x3a, 0x55, 0x52, 0xdb, 0x66, 0xee, 0xd1, 0x5c, 0xae, 0x89, 0x7e, 0x0e, 0x31, 0xe9,
	0xa3, 0x3b, 0x3f, 0x7e, 0x12, 0xec, 0x2e, 0xc7, 0x7f, 0xb7, 0xe0, 0xd3, 0xfe, 0x6c, 0x7f, 0x59,
	0x63, 0x56, 0x5b, 0xa9, 0x0f, 0x96, 0x39, 0x81, 0x50, 0x2d, 0x85, 0x41, 0xaa, 0x11, 0x73, 0x1f,
	0x38, 0x54, 0xa3, 0xd5, 0x6f, 0x4d, 0x76, 0x1f, 0xb8, 0x56, 0x51, 0x6b, 0xa9, 0xd3, 0x6c, 0xe3,
	0xb0, 0x01, 0x8f, 0x09, 0xb9, 0x73, 0x36, 0xfb, 0x06, 0x06, 0xfe, 0xba, 0x44, 0x63, 0xc4, 0x02,
	0x93, 0x90, 0x52, 0xf6, 0x09, 0x7c, 0xf2, 0xd8, 0xf8, 0xdf, 0xa0, 0x31, 0xfb, 0x56, 0xde, 0x41,
	0x59, 0x3f, 0x43, 0xbf, 0x28, 0xc5, 0x02, 0xd3, 0x5c, 0x17, 0xaf, 0xa8, 0x49, 0xdd, 0xf1, 0xcd,
	0xd7, 0xbb, 0xde, 0x9d, 0x29, 0xcc, 0x6e, 0xe5, 0xfa, 0x57, 0x47, 0xbb, 0x27, 0x16, 0xef, 0x15,
	0xef, 0x01, 0xfb, 0x01, 0x06, 0x3e, 0x85, 0x54, 0x6e, 0x33, 0x4c, 0x12, 0x90, 0xbb, 0xfa, 0x4d,
	0x8e, 0x47, 0x31, 0xc7, 0x15, 0xf7, 0x55, 0x9e, 0x3d, 0x83, 0x8d, 0x61, 0xa0, 0xd1, 0xa4, 0x99,
	0xaa, 0xd3, 0x55, 0x51, 0x16, 0x96, 0x7a, 0x0c, 0x78, 0x4f, 0xa3, 0xb9, 0x53, 0xf5, 0xa3, 0x83,
	0x36, 0x9c, 0x12, 0xcb, 0x86, 0x13, 0x6e, 0x39, 0x4f, 0x58, 0x7a, 0xce, 0xf7, 0xd0, 0x21, 0x07,
	0x98, 0xa4, 0x43, 0x35, 0x4f, 0x3e, 0x38, 0x9a, 0x6c, 0xcb, 0x1b, 0x0e, 0xbb, 0x82, 0xd0, 0x4d,
	0xd2, 0xad, 0xd2, 0x3e, 0x79, 0xc7, 0x22, 0xdc, 0x53, 0xdc, 0x1a, 0x65, 0xb2, 0x2c, 0x45, 0xe5,
	0xd6, 0x28, 0x70, 0xce, 0x6c, 0x42, 0xf6, 0x13, 0xc4, 0xd8, 0x0c, 0x7a, 0xb3, 0x48, 0x17, 0x07,
	0x57, 0x7d, 0x63, 0x07, 0xfe, 0xce, 0xdf, 0xf9, 0x9b, 0xc0, 0xde, 0xdf, 0x24, 0x81, 0xc8, 0x58,
	0xa1, 0xdd, 0xd6, 0xf6, 0xbc, 0x07, 0x9b, 0x70, 0x77, 0x9f, 0xfb, 0x7b, 0xfb, 0x7c, 0xf5, 0x00,
	0x5f, 0xfd, 0xcf, 0x7c, 0x58, 0x0f, 0xa2, 0xdf, 0xaa, 0x97, 0x4a, 0xfe, 0x59, 0x0d, 0x3f, 0x63,
	0x00, 0x9d, 0x7b, 0x99, 0xbd, 0xa0, 0x1e, 0xb6, 0x58, 0x04, 0x01, 0x7f, 0xb1, 0xc3, 0x36, 0x8b,
	0x21, 0x9c, 0xca, 0x3a, 0x5b, 0x0e, 0x83, 0x79, 0x87, 0xfe, 0x8e, 0x3f, 0xfe, 0x17, 0x00, 0x00,
	0xff, 0xff, 0x77, 0xbf, 0x19, 0x06, 0x72, 0x05, 0x00, 0x00,
}
