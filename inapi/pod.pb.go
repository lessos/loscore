// Code generated by protoc-gen-go.
// source: inapi/pod.proto
// DO NOT EDIT!

package inapi

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type PbPodSpecBoxImageDriver int32

const (
	PbPodSpecBoxImageDriver_Unknown PbPodSpecBoxImageDriver = 0
	PbPodSpecBoxImageDriver_Docker  PbPodSpecBoxImageDriver = 1
	PbPodSpecBoxImageDriver_Rkt     PbPodSpecBoxImageDriver = 2
)

var PbPodSpecBoxImageDriver_name = map[int32]string{
	0: "Unknown",
	1: "Docker",
	2: "Rkt",
}
var PbPodSpecBoxImageDriver_value = map[string]int32{
	"Unknown": 0,
	"Docker":  1,
	"Rkt":     2,
}

func (x PbPodSpecBoxImageDriver) String() string {
	return proto.EnumName(PbPodSpecBoxImageDriver_name, int32(x))
}
func (PbPodSpecBoxImageDriver) EnumDescriptor() ([]byte, []int) { return fileDescriptor3, []int{0} }

type PbPodRepStatus struct {
	Id      string             `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Rep     uint32             `protobuf:"varint,2,opt,name=rep" json:"rep,omitempty"`
	Action  uint32             `protobuf:"varint,3,opt,name=action" json:"action,omitempty"`
	Node    string             `protobuf:"bytes,4,opt,name=node" json:"node,omitempty"`
	Boxes   []*PbPodBoxStatus  `protobuf:"bytes,5,rep,name=boxes" json:"boxes,omitempty"`
	OpLog   *PbOpLogSets       `protobuf:"bytes,6,opt,name=op_log,json=opLog" json:"op_log,omitempty"`
	Stats   *PbStatsSampleFeed `protobuf:"bytes,7,opt,name=stats" json:"stats,omitempty"`
	Updated uint32             `protobuf:"varint,8,opt,name=updated" json:"updated,omitempty"`
}

func (m *PbPodRepStatus) Reset()                    { *m = PbPodRepStatus{} }
func (m *PbPodRepStatus) String() string            { return proto.CompactTextString(m) }
func (*PbPodRepStatus) ProtoMessage()               {}
func (*PbPodRepStatus) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{0} }

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

func (m *PbPodRepStatus) GetBoxes() []*PbPodBoxStatus {
	if m != nil {
		return m.Boxes
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

type PbVolumeMount struct {
	Name      string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	ReadOnly  bool   `protobuf:"varint,2,opt,name=read_only,json=readOnly" json:"read_only,omitempty"`
	MountPath string `protobuf:"bytes,3,opt,name=mount_path,json=mountPath" json:"mount_path,omitempty"`
	HostDir   string `protobuf:"bytes,4,opt,name=host_dir,json=hostDir" json:"host_dir,omitempty"`
}

func (m *PbVolumeMount) Reset()                    { *m = PbVolumeMount{} }
func (m *PbVolumeMount) String() string            { return proto.CompactTextString(m) }
func (*PbVolumeMount) ProtoMessage()               {}
func (*PbVolumeMount) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{1} }

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

type PbServicePort struct {
	Name     string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	BoxPort  uint32 `protobuf:"varint,2,opt,name=box_port,json=boxPort" json:"box_port,omitempty"`
	HostPort uint32 `protobuf:"varint,3,opt,name=host_port,json=hostPort" json:"host_port,omitempty"`
}

func (m *PbServicePort) Reset()                    { *m = PbServicePort{} }
func (m *PbServicePort) String() string            { return proto.CompactTextString(m) }
func (*PbServicePort) ProtoMessage()               {}
func (*PbServicePort) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{2} }

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
	Name         string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Phase        string `protobuf:"bytes,2,opt,name=phase" json:"phase,omitempty"`
	Retry        uint32 `protobuf:"varint,3,opt,name=retry" json:"retry,omitempty"`
	ErrorCode    uint32 `protobuf:"varint,4,opt,name=error_code,json=errorCode" json:"error_code,omitempty"`
	ErrorMessage string `protobuf:"bytes,5,opt,name=error_message,json=errorMessage" json:"error_message,omitempty"`
}

func (m *PbPodBoxStatusExecutor) Reset()                    { *m = PbPodBoxStatusExecutor{} }
func (m *PbPodBoxStatusExecutor) String() string            { return proto.CompactTextString(m) }
func (*PbPodBoxStatusExecutor) ProtoMessage()               {}
func (*PbPodBoxStatusExecutor) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{3} }

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
	Name         string                    `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	ImageDriver  PbPodSpecBoxImageDriver   `protobuf:"varint,2,opt,name=image_driver,json=imageDriver,enum=inapi.PbPodSpecBoxImageDriver" json:"image_driver,omitempty"`
	ImageOptions []*Label                  `protobuf:"bytes,3,rep,name=image_options,json=imageOptions" json:"image_options,omitempty"`
	ResCpuLimit  int64                     `protobuf:"varint,4,opt,name=res_cpu_limit,json=resCpuLimit" json:"res_cpu_limit,omitempty"`
	ResMemLimit  int64                     `protobuf:"varint,5,opt,name=res_mem_limit,json=resMemLimit" json:"res_mem_limit,omitempty"`
	Mounts       []*PbVolumeMount          `protobuf:"bytes,6,rep,name=mounts" json:"mounts,omitempty"`
	Ports        []*PbServicePort          `protobuf:"bytes,7,rep,name=ports" json:"ports,omitempty"`
	Command      []string                  `protobuf:"bytes,8,rep,name=command" json:"command,omitempty"`
	Executors    []*PbPodBoxStatusExecutor `protobuf:"bytes,9,rep,name=executors" json:"executors,omitempty"`
	Action       uint32                    `protobuf:"varint,10,opt,name=action" json:"action,omitempty"`
	Started      uint32                    `protobuf:"varint,11,opt,name=started" json:"started,omitempty"`
	Updated      uint32                    `protobuf:"varint,12,opt,name=updated" json:"updated,omitempty"`
}

func (m *PbPodBoxStatus) Reset()                    { *m = PbPodBoxStatus{} }
func (m *PbPodBoxStatus) String() string            { return proto.CompactTextString(m) }
func (*PbPodBoxStatus) ProtoMessage()               {}
func (*PbPodBoxStatus) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{4} }

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
	proto.RegisterType((*PbServicePort)(nil), "inapi.PbServicePort")
	proto.RegisterType((*PbPodBoxStatusExecutor)(nil), "inapi.PbPodBoxStatusExecutor")
	proto.RegisterType((*PbPodBoxStatus)(nil), "inapi.PbPodBoxStatus")
	proto.RegisterEnum("inapi.PbPodSpecBoxImageDriver", PbPodSpecBoxImageDriver_name, PbPodSpecBoxImageDriver_value)
}

func init() { proto.RegisterFile("inapi/pod.proto", fileDescriptor3) }

var fileDescriptor3 = []byte{
	// 673 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x54, 0x41, 0x6f, 0x13, 0x3d,
	0x10, 0xfd, 0x92, 0xed, 0x26, 0xd9, 0x49, 0xd2, 0x2f, 0x58, 0xa5, 0xb8, 0xad, 0x8a, 0xa2, 0x70,
	0x09, 0x05, 0x05, 0x51, 0x8e, 0x3d, 0xd1, 0x16, 0x24, 0xa4, 0x56, 0x8d, 0x1c, 0xc1, 0x85, 0xc3,
	0xca, 0xbb, 0x3b, 0x4a, 0x56, 0xcd, 0xae, 0x2d, 0xdb, 0x29, 0x29, 0x27, 0x7e, 0x08, 0xff, 0x8d,
	0xbf, 0x82, 0x6c, 0x6f, 0xd2, 0x44, 0x0a, 0x37, 0xcf, 0xf3, 0xdb, 0x79, 0xe3, 0x99, 0x37, 0x0b,
	0xff, 0xe7, 0x25, 0x97, 0xf9, 0x3b, 0x29, 0xb2, 0x91, 0x54, 0xc2, 0x08, 0x12, 0x3a, 0xe0, 0xb8,
	0xe7, 0xf1, 0x84, 0x6b, 0xf4, 0x17, 0xc7, 0x07, 0x1e, 0x11, 0x12, 0x15, 0x37, 0x42, 0x55, 0xe8,
	0x33, 0x8f, 0x6a, 0xc3, 0x8d, 0xf6, 0xd0, 0xe0, 0x57, 0x1d, 0xf6, 0xc7, 0xc9, 0x58, 0x64, 0x0c,
	0xe5, 0xc4, 0x70, 0xb3, 0xd0, 0x64, 0x1f, 0xea, 0x79, 0x46, 0x6b, 0xfd, 0xda, 0x30, 0x62, 0xf5,
	0x3c, 0x23, 0x3d, 0x08, 0x14, 0x4a, 0x5a, 0xef, 0xd7, 0x86, 0x5d, 0x66, 0x8f, 0xe4, 0x10, 0x1a,
	0x3c, 0x35, 0xb9, 0x28, 0x69, 0xe0, 0xc0, 0x2a, 0x22, 0x04, 0xf6, 0x4a, 0x91, 0x21, 0xdd, 0x73,
	0xdf, 0xba, 0x33, 0x79, 0x03, 0x61, 0x22, 0x96, 0xa8, 0x69, 0xd8, 0x0f, 0x86, 0xed, 0xf3, 0xe7,
	0x23, 0x57, 0xc3, 0xc8, 0x69, 0x5e, 0x8a, 0xa5, 0xd7, 0x64, 0x9e, 0x43, 0x5e, 0x43, 0x43, 0xc8,
	0x78, 0x2e, 0xa6, 0xb4, 0xd1, 0xaf, 0x0d, 0xdb, 0xe7, 0x64, 0xcd, 0xbe, 0x93, 0x37, 0x62, 0x3a,
	0x41, 0xa3, 0x59, 0x28, 0xec, 0x91, 0x8c, 0x20, 0x74, 0xef, 0xa0, 0x4d, 0xc7, 0xa4, 0x6b, 0xa6,
	0xcd, 0xa8, 0x27, 0xbc, 0x90, 0x73, 0xfc, 0x8c, 0x98, 0x31, 0x4f, 0x23, 0x14, 0x9a, 0x0b, 0x99,
	0x71, 0x83, 0x19, 0x6d, 0xb9, 0xa2, 0x57, 0xe1, 0xe0, 0x27, 0x74, 0xc7, 0xc9, 0x37, 0x31, 0x5f,
	0x14, 0x78, 0x2b, 0x16, 0xa5, 0x71, 0xcf, 0xe0, 0x05, 0x56, 0x2d, 0x70, 0x67, 0x72, 0x02, 0x91,
	0x42, 0x9e, 0xc5, 0xa2, 0x9c, 0x3f, 0xba, 0x56, 0xb4, 0x58, 0xcb, 0x02, 0x77, 0xe5, 0xfc, 0x91,
	0x9c, 0x02, 0x14, 0xf6, 0xcb, 0x58, 0x72, 0x33, 0x73, 0x3d, 0x89, 0x58, 0xe4, 0x90, 0x31, 0x37,
	0x33, 0x72, 0x04, 0xad, 0x99, 0xd0, 0x26, 0xce, 0x72, 0x55, 0xb5, 0xa6, 0x69, 0xe3, 0xeb, 0x5c,
	0x0d, 0xbe, 0x5b, 0xed, 0x09, 0xaa, 0x87, 0x3c, 0xc5, 0xb1, 0x50, 0xbb, 0xb5, 0x8f, 0xa0, 0x95,
	0x88, 0x65, 0x2c, 0x85, 0x32, 0xd5, 0x14, 0x9a, 0x89, 0x58, 0x3a, 0xfa, 0x09, 0x44, 0x2e, 0xb5,
	0xbb, 0xf3, 0xc3, 0x70, 0x5a, 0xf6, 0x72, 0xf0, 0xbb, 0x06, 0x87, 0xdb, 0x7d, 0xfe, 0xb4, 0xc4,
	0x74, 0x61, 0x84, 0xda, 0x29, 0x73, 0x00, 0xa1, 0x9c, 0x71, 0x8d, 0x4e, 0x23, 0x62, 0x3e, 0xb0,
	0xa8, 0x42, 0xa3, 0x1e, 0xab, 0xec, 0x3e, 0xb0, 0x2f, 0x46, 0xa5, 0x84, 0x8a, 0xd3, 0xd5, 0xbc,
	0xbb, 0x2c, 0x72, 0xc8, 0x95, 0x1d, 0xfa, 0x2b, 0xe8, 0xfa, 0xeb, 0x02, 0xb5, 0xe6, 0x53, 0xa4,
	0xa1, 0x4b, 0xd9, 0x71, 0xe0, 0xad, 0xc7, 0x06, 0x7f, 0x82, 0xca, 0x7a, 0xeb, 0xf2, 0x76, 0x96,
	0xf5, 0x11, 0x3a, 0x79, 0xc1, 0xa7, 0x18, 0x67, 0x2a, 0x7f, 0x40, 0xe5, 0xaa, 0xdb, 0x3f, 0x7f,
	0xb9, 0xe9, 0xa3, 0x89, 0xc4, 0xf4, 0x52, 0x2c, 0xbf, 0x58, 0xda, 0xb5, 0x63, 0xb1, 0x76, 0xfe,
	0x14, 0x90, 0xf7, 0xd0, 0xf5, 0x29, 0x84, 0xb4, 0x3e, 0xd5, 0x34, 0x70, 0x5e, 0xec, 0x54, 0x39,
	0x6e, 0x78, 0x82, 0x73, 0xe6, 0x55, 0xee, 0x3c, 0x83, 0x0c, 0xa0, 0xab, 0x50, 0xc7, 0xa9, 0x5c,
	0xc4, 0xf3, 0xbc, 0xc8, 0x8d, 0x7b, 0x63, 0xc0, 0xda, 0x0a, 0xf5, 0x95, 0x5c, 0xdc, 0x58, 0x68,
	0xc5, 0x29, 0xb0, 0xa8, 0x38, 0xe1, 0x9a, 0x73, 0x8b, 0x85, 0xe7, 0xbc, 0x85, 0x86, 0x33, 0x82,
	0xa6, 0x0d, 0xa7, 0x79, 0xb0, 0xae, 0x7b, 0xc3, 0x71, 0xac, 0xe2, 0x90, 0x33, 0x08, 0xed, 0x24,
	0xad, 0xa9, 0xb7, 0xc9, 0x1b, 0x16, 0x61, 0x9e, 0x62, 0x0d, 0x9d, 0x8a, 0xa2, 0xe0, 0xa5, 0x35,
	0x74, 0x60, 0x4d, 0x55, 0x85, 0xe4, 0x02, 0x22, 0xac, 0x06, 0xad, 0x69, 0xe4, 0x32, 0x9d, 0xee,
	0x5c, 0xbb, 0x95, 0x1d, 0xd8, 0x13, 0x7f, 0x63, 0xb7, 0x61, 0x6b, 0xb7, 0x29, 0x34, 0xb5, 0xe1,
	0xca, 0xee, 0x4f, 0xdb, 0x7b, 0xb0, 0x0a, 0x37, 0x37, 0xab, 0xb3, 0xb5, 0x59, 0x67, 0x17, 0xf0,
	0xe2, 0x1f, 0xf3, 0x21, 0x6d, 0x68, 0x7e, 0x2d, 0xef, 0x4b, 0xf1, 0xa3, 0xec, 0xfd, 0x47, 0x00,
	0x1a, 0xd7, 0x22, 0xbd, 0x47, 0xd5, 0xab, 0x91, 0x26, 0x04, 0xec, 0xde, 0xf4, 0xea, 0x49, 0xc3,
	0xfd, 0xa0, 0x3e, 0xfc, 0x0d, 0x00, 0x00, 0xff, 0xff, 0xb3, 0xa3, 0x1f, 0x3e, 0xf5, 0x04, 0x00,
	0x00,
}
