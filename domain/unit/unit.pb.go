// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v4.25.2
// source: domain/unit/unit.proto

package unit

import (
	battlefield "github.com/corstijank/mekstrike/domain/battlefield"
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

type Stats struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name       string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Model      string   `protobuf:"bytes,2,opt,name=model,proto3" json:"model,omitempty"`
	Pointvalue int32    `protobuf:"varint,3,opt,name=pointvalue,proto3" json:"pointvalue,omitempty"`
	Type       string   `protobuf:"bytes,4,opt,name=type,proto3" json:"type,omitempty"`
	Size       int32    `protobuf:"varint,5,opt,name=size,proto3" json:"size,omitempty"`
	Movement   string   `protobuf:"bytes,6,opt,name=movement,proto3" json:"movement,omitempty"`
	Role       string   `protobuf:"bytes,7,opt,name=role,proto3" json:"role,omitempty"`
	Shortdmg   int32    `protobuf:"varint,8,opt,name=shortdmg,proto3" json:"shortdmg,omitempty"`
	Meddmg     int32    `protobuf:"varint,9,opt,name=meddmg,proto3" json:"meddmg,omitempty"`
	Longdmg    int32    `protobuf:"varint,10,opt,name=longdmg,proto3" json:"longdmg,omitempty"`
	Ovhdmg     int32    `protobuf:"varint,11,opt,name=ovhdmg,proto3" json:"ovhdmg,omitempty"`
	Armor      int32    `protobuf:"varint,12,opt,name=armor,proto3" json:"armor,omitempty"`
	Struct     int32    `protobuf:"varint,13,opt,name=struct,proto3" json:"struct,omitempty"`
	Specials   []string `protobuf:"bytes,14,rep,name=specials,proto3" json:"specials,omitempty"`
	Image      string   `protobuf:"bytes,15,opt,name=image,proto3" json:"image,omitempty"`
}

func (x *Stats) Reset() {
	*x = Stats{}
	if protoimpl.UnsafeEnabled {
		mi := &file_domain_unit_unit_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Stats) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Stats) ProtoMessage() {}

func (x *Stats) ProtoReflect() protoreflect.Message {
	mi := &file_domain_unit_unit_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Stats.ProtoReflect.Descriptor instead.
func (*Stats) Descriptor() ([]byte, []int) {
	return file_domain_unit_unit_proto_rawDescGZIP(), []int{0}
}

func (x *Stats) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Stats) GetModel() string {
	if x != nil {
		return x.Model
	}
	return ""
}

func (x *Stats) GetPointvalue() int32 {
	if x != nil {
		return x.Pointvalue
	}
	return 0
}

func (x *Stats) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *Stats) GetSize() int32 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *Stats) GetMovement() string {
	if x != nil {
		return x.Movement
	}
	return ""
}

func (x *Stats) GetRole() string {
	if x != nil {
		return x.Role
	}
	return ""
}

func (x *Stats) GetShortdmg() int32 {
	if x != nil {
		return x.Shortdmg
	}
	return 0
}

func (x *Stats) GetMeddmg() int32 {
	if x != nil {
		return x.Meddmg
	}
	return 0
}

func (x *Stats) GetLongdmg() int32 {
	if x != nil {
		return x.Longdmg
	}
	return 0
}

func (x *Stats) GetOvhdmg() int32 {
	if x != nil {
		return x.Ovhdmg
	}
	return 0
}

func (x *Stats) GetArmor() int32 {
	if x != nil {
		return x.Armor
	}
	return 0
}

func (x *Stats) GetStruct() int32 {
	if x != nil {
		return x.Struct
	}
	return 0
}

func (x *Stats) GetSpecials() []string {
	if x != nil {
		return x.Specials
	}
	return nil
}

func (x *Stats) GetImage() string {
	if x != nil {
		return x.Image
	}
	return ""
}

type Location struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BattlefieldId string                   `protobuf:"bytes,1,opt,name=battlefieldId,proto3" json:"battlefieldId,omitempty"`
	Position      *battlefield.Coordinates `protobuf:"bytes,2,opt,name=position,proto3" json:"position,omitempty"`
	Heading       int32                    `protobuf:"varint,3,opt,name=heading,proto3" json:"heading,omitempty"`
}

func (x *Location) Reset() {
	*x = Location{}
	if protoimpl.UnsafeEnabled {
		mi := &file_domain_unit_unit_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Location) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Location) ProtoMessage() {}

func (x *Location) ProtoReflect() protoreflect.Message {
	mi := &file_domain_unit_unit_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Location.ProtoReflect.Descriptor instead.
func (*Location) Descriptor() ([]byte, []int) {
	return file_domain_unit_unit_proto_rawDescGZIP(), []int{1}
}

func (x *Location) GetBattlefieldId() string {
	if x != nil {
		return x.BattlefieldId
	}
	return ""
}

func (x *Location) GetPosition() *battlefield.Coordinates {
	if x != nil {
		return x.Position
	}
	return nil
}

func (x *Location) GetHeading() int32 {
	if x != nil {
		return x.Heading
	}
	return 0
}

type Data struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Stats    *Stats    `protobuf:"bytes,1,opt,name=stats,proto3" json:"stats,omitempty"`
	Location *Location `protobuf:"bytes,2,opt,name=location,proto3" json:"location,omitempty"`
	Owner    string    `protobuf:"bytes,3,opt,name=owner,proto3" json:"owner,omitempty"`
	Active   bool      `protobuf:"varint,4,opt,name=active,proto3" json:"active,omitempty"`
}

func (x *Data) Reset() {
	*x = Data{}
	if protoimpl.UnsafeEnabled {
		mi := &file_domain_unit_unit_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Data) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Data) ProtoMessage() {}

func (x *Data) ProtoReflect() protoreflect.Message {
	mi := &file_domain_unit_unit_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Data.ProtoReflect.Descriptor instead.
func (*Data) Descriptor() ([]byte, []int) {
	return file_domain_unit_unit_proto_rawDescGZIP(), []int{2}
}

func (x *Data) GetStats() *Stats {
	if x != nil {
		return x.Stats
	}
	return nil
}

func (x *Data) GetLocation() *Location {
	if x != nil {
		return x.Location
	}
	return nil
}

func (x *Data) GetOwner() string {
	if x != nil {
		return x.Owner
	}
	return ""
}

func (x *Data) GetActive() bool {
	if x != nil {
		return x.Active
	}
	return false
}

type DeployRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BattlefieldId string `protobuf:"bytes,1,opt,name=battlefieldId,proto3" json:"battlefieldId,omitempty"`
	Owner         string `protobuf:"bytes,2,opt,name=owner,proto3" json:"owner,omitempty"`
	Stats         *Stats `protobuf:"bytes,3,opt,name=stats,proto3" json:"stats,omitempty"`
	Corner        string `protobuf:"bytes,4,opt,name=corner,proto3" json:"corner,omitempty"`
}

func (x *DeployRequest) Reset() {
	*x = DeployRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_domain_unit_unit_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeployRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeployRequest) ProtoMessage() {}

func (x *DeployRequest) ProtoReflect() protoreflect.Message {
	mi := &file_domain_unit_unit_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeployRequest.ProtoReflect.Descriptor instead.
func (*DeployRequest) Descriptor() ([]byte, []int) {
	return file_domain_unit_unit_proto_rawDescGZIP(), []int{3}
}

func (x *DeployRequest) GetBattlefieldId() string {
	if x != nil {
		return x.BattlefieldId
	}
	return ""
}

func (x *DeployRequest) GetOwner() string {
	if x != nil {
		return x.Owner
	}
	return ""
}

func (x *DeployRequest) GetStats() *Stats {
	if x != nil {
		return x.Stats
	}
	return nil
}

func (x *DeployRequest) GetCorner() string {
	if x != nil {
		return x.Corner
	}
	return ""
}

var File_domain_unit_unit_proto protoreflect.FileDescriptor

var file_domain_unit_unit_proto_rawDesc = []byte{
	0x0a, 0x16, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2f, 0x75, 0x6e, 0x69, 0x74, 0x2f, 0x75, 0x6e,
	0x69, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x19, 0x6e, 0x65, 0x74, 0x2e, 0x6d, 0x65,
	0x6b, 0x73, 0x74, 0x72, 0x69, 0x6b, 0x65, 0x2e, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x75,
	0x6e, 0x69, 0x74, 0x1a, 0x24, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2f, 0x62, 0x61, 0x74, 0x74,
	0x6c, 0x65, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x2f, 0x62, 0x61, 0x74, 0x74, 0x6c, 0x65, 0x66, 0x69,
	0x65, 0x6c, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xef, 0x02, 0x0a, 0x05, 0x53, 0x74,
	0x61, 0x74, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x6d, 0x6f, 0x64, 0x65, 0x6c,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x12, 0x1e, 0x0a,
	0x0a, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x0a, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x12, 0x0a,
	0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70,
	0x65, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x04, 0x73, 0x69, 0x7a, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x6d, 0x6f, 0x76, 0x65, 0x6d, 0x65, 0x6e,
	0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6d, 0x6f, 0x76, 0x65, 0x6d, 0x65, 0x6e,
	0x74, 0x12, 0x12, 0x0a, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x72, 0x6f, 0x6c, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x64, 0x6d,
	0x67, 0x18, 0x08, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x64, 0x6d,
	0x67, 0x12, 0x16, 0x0a, 0x06, 0x6d, 0x65, 0x64, 0x64, 0x6d, 0x67, 0x18, 0x09, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x06, 0x6d, 0x65, 0x64, 0x64, 0x6d, 0x67, 0x12, 0x18, 0x0a, 0x07, 0x6c, 0x6f, 0x6e,
	0x67, 0x64, 0x6d, 0x67, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x6c, 0x6f, 0x6e, 0x67,
	0x64, 0x6d, 0x67, 0x12, 0x16, 0x0a, 0x06, 0x6f, 0x76, 0x68, 0x64, 0x6d, 0x67, 0x18, 0x0b, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x06, 0x6f, 0x76, 0x68, 0x64, 0x6d, 0x67, 0x12, 0x14, 0x0a, 0x05, 0x61,
	0x72, 0x6d, 0x6f, 0x72, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x61, 0x72, 0x6d, 0x6f,
	0x72, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x18, 0x0d, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x06, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x70, 0x65,
	0x63, 0x69, 0x61, 0x6c, 0x73, 0x18, 0x0e, 0x20, 0x03, 0x28, 0x09, 0x52, 0x08, 0x73, 0x70, 0x65,
	0x63, 0x69, 0x61, 0x6c, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x18, 0x0f,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x22, 0x95, 0x01, 0x0a, 0x08,
	0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x24, 0x0a, 0x0d, 0x62, 0x61, 0x74, 0x74,
	0x6c, 0x65, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0d, 0x62, 0x61, 0x74, 0x74, 0x6c, 0x65, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x49, 0x64, 0x12, 0x49,
	0x0a, 0x08, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x2d, 0x2e, 0x6e, 0x65, 0x74, 0x2e, 0x6d, 0x65, 0x6b, 0x73, 0x74, 0x72, 0x69, 0x6b, 0x65,
	0x2e, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x62, 0x61, 0x74, 0x74, 0x6c, 0x65, 0x66, 0x69,
	0x65, 0x6c, 0x64, 0x2e, 0x43, 0x6f, 0x6f, 0x72, 0x64, 0x69, 0x6e, 0x61, 0x74, 0x65, 0x73, 0x52,
	0x08, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x18, 0x0a, 0x07, 0x68, 0x65, 0x61,
	0x64, 0x69, 0x6e, 0x67, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x68, 0x65, 0x61, 0x64,
	0x69, 0x6e, 0x67, 0x22, 0xad, 0x01, 0x0a, 0x04, 0x44, 0x61, 0x74, 0x61, 0x12, 0x36, 0x0a, 0x05,
	0x73, 0x74, 0x61, 0x74, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x6e, 0x65,
	0x74, 0x2e, 0x6d, 0x65, 0x6b, 0x73, 0x74, 0x72, 0x69, 0x6b, 0x65, 0x2e, 0x64, 0x6f, 0x6d, 0x61,
	0x69, 0x6e, 0x2e, 0x75, 0x6e, 0x69, 0x74, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x73, 0x52, 0x05, 0x73,
	0x74, 0x61, 0x74, 0x73, 0x12, 0x3f, 0x0a, 0x08, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x23, 0x2e, 0x6e, 0x65, 0x74, 0x2e, 0x6d, 0x65, 0x6b,
	0x73, 0x74, 0x72, 0x69, 0x6b, 0x65, 0x2e, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x75, 0x6e,
	0x69, 0x74, 0x2e, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x08, 0x6c, 0x6f, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x61,
	0x63, 0x74, 0x69, 0x76, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x61, 0x63, 0x74,
	0x69, 0x76, 0x65, 0x22, 0x9b, 0x01, 0x0a, 0x0d, 0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x24, 0x0a, 0x0d, 0x62, 0x61, 0x74, 0x74, 0x6c, 0x65, 0x66,
	0x69, 0x65, 0x6c, 0x64, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x62, 0x61,
	0x74, 0x74, 0x6c, 0x65, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x6f,
	0x77, 0x6e, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6f, 0x77, 0x6e, 0x65,
	0x72, 0x12, 0x36, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x20, 0x2e, 0x6e, 0x65, 0x74, 0x2e, 0x6d, 0x65, 0x6b, 0x73, 0x74, 0x72, 0x69, 0x6b, 0x65,
	0x2e, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x75, 0x6e, 0x69, 0x74, 0x2e, 0x53, 0x74, 0x61,
	0x74, 0x73, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x63, 0x6f, 0x72,
	0x6e, 0x65, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x63, 0x6f, 0x72, 0x6e, 0x65,
	0x72, 0x42, 0x4e, 0x0a, 0x19, 0x6e, 0x65, 0x74, 0x2e, 0x6d, 0x65, 0x6b, 0x73, 0x74, 0x72, 0x69,
	0x6b, 0x65, 0x2e, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x75, 0x6e, 0x69, 0x74, 0x42, 0x04,
	0x55, 0x6e, 0x69, 0x74, 0x5a, 0x2b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x63, 0x6f, 0x72, 0x73, 0x74, 0x69, 0x6a, 0x61, 0x6e, 0x6b, 0x2f, 0x6d, 0x65, 0x6b, 0x73,
	0x74, 0x72, 0x69, 0x6b, 0x65, 0x2f, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2f, 0x75, 0x6e, 0x69,
	0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_domain_unit_unit_proto_rawDescOnce sync.Once
	file_domain_unit_unit_proto_rawDescData = file_domain_unit_unit_proto_rawDesc
)

func file_domain_unit_unit_proto_rawDescGZIP() []byte {
	file_domain_unit_unit_proto_rawDescOnce.Do(func() {
		file_domain_unit_unit_proto_rawDescData = protoimpl.X.CompressGZIP(file_domain_unit_unit_proto_rawDescData)
	})
	return file_domain_unit_unit_proto_rawDescData
}

var file_domain_unit_unit_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_domain_unit_unit_proto_goTypes = []interface{}{
	(*Stats)(nil),                   // 0: net.mekstrike.domain.unit.Stats
	(*Location)(nil),                // 1: net.mekstrike.domain.unit.Location
	(*Data)(nil),                    // 2: net.mekstrike.domain.unit.Data
	(*DeployRequest)(nil),           // 3: net.mekstrike.domain.unit.DeployRequest
	(*battlefield.Coordinates)(nil), // 4: net.mekstrike.domain.battlefield.Coordinates
}
var file_domain_unit_unit_proto_depIdxs = []int32{
	4, // 0: net.mekstrike.domain.unit.Location.position:type_name -> net.mekstrike.domain.battlefield.Coordinates
	0, // 1: net.mekstrike.domain.unit.Data.stats:type_name -> net.mekstrike.domain.unit.Stats
	1, // 2: net.mekstrike.domain.unit.Data.location:type_name -> net.mekstrike.domain.unit.Location
	0, // 3: net.mekstrike.domain.unit.DeployRequest.stats:type_name -> net.mekstrike.domain.unit.Stats
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_domain_unit_unit_proto_init() }
func file_domain_unit_unit_proto_init() {
	if File_domain_unit_unit_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_domain_unit_unit_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Stats); i {
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
		file_domain_unit_unit_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Location); i {
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
		file_domain_unit_unit_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Data); i {
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
		file_domain_unit_unit_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeployRequest); i {
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
			RawDescriptor: file_domain_unit_unit_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_domain_unit_unit_proto_goTypes,
		DependencyIndexes: file_domain_unit_unit_proto_depIdxs,
		MessageInfos:      file_domain_unit_unit_proto_msgTypes,
	}.Build()
	File_domain_unit_unit_proto = out.File
	file_domain_unit_unit_proto_rawDesc = nil
	file_domain_unit_unit_proto_goTypes = nil
	file_domain_unit_unit_proto_depIdxs = nil
}