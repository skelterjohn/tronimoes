// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.15.0
// source: tiles/proto/tiles.proto

package proto

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

type Placement_Type int32

const (
	Placement_TYPE_UNKNOWN        Placement_Type = 0
	Placement_PLAYER_LEADER       Placement_Type = 1
	Placement_PLAYER_CONTINUATION Placement_Type = 2
	Placement_FREE_LEADER         Placement_Type = 3
	Placement_FREE_CONTINUATION   Placement_Type = 4
)

// Enum value maps for Placement_Type.
var (
	Placement_Type_name = map[int32]string{
		0: "TYPE_UNKNOWN",
		1: "PLAYER_LEADER",
		2: "PLAYER_CONTINUATION",
		3: "FREE_LEADER",
		4: "FREE_CONTINUATION",
	}
	Placement_Type_value = map[string]int32{
		"TYPE_UNKNOWN":        0,
		"PLAYER_LEADER":       1,
		"PLAYER_CONTINUATION": 2,
		"FREE_LEADER":         3,
		"FREE_CONTINUATION":   4,
	}
)

func (x Placement_Type) Enum() *Placement_Type {
	p := new(Placement_Type)
	*p = x
	return p
}

func (x Placement_Type) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Placement_Type) Descriptor() protoreflect.EnumDescriptor {
	return file_tiles_proto_tiles_proto_enumTypes[0].Descriptor()
}

func (Placement_Type) Type() protoreflect.EnumType {
	return &file_tiles_proto_tiles_proto_enumTypes[0]
}

func (x Placement_Type) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Placement_Type.Descriptor instead.
func (Placement_Type) EnumDescriptor() ([]byte, []int) {
	return file_tiles_proto_tiles_proto_rawDescGZIP(), []int{2, 0}
}

type Tile struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	A int32 `protobuf:"varint,1,opt,name=a,proto3" json:"a,omitempty"`
	B int32 `protobuf:"varint,2,opt,name=b,proto3" json:"b,omitempty"`
}

func (x *Tile) Reset() {
	*x = Tile{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tiles_proto_tiles_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Tile) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Tile) ProtoMessage() {}

func (x *Tile) ProtoReflect() protoreflect.Message {
	mi := &file_tiles_proto_tiles_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Tile.ProtoReflect.Descriptor instead.
func (*Tile) Descriptor() ([]byte, []int) {
	return file_tiles_proto_tiles_proto_rawDescGZIP(), []int{0}
}

func (x *Tile) GetA() int32 {
	if x != nil {
		return x.A
	}
	return 0
}

func (x *Tile) GetB() int32 {
	if x != nil {
		return x.B
	}
	return 0
}

type Coord struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	X int32 `protobuf:"varint,1,opt,name=x,proto3" json:"x,omitempty"`
	Y int32 `protobuf:"varint,2,opt,name=y,proto3" json:"y,omitempty"`
}

func (x *Coord) Reset() {
	*x = Coord{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tiles_proto_tiles_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Coord) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Coord) ProtoMessage() {}

func (x *Coord) ProtoReflect() protoreflect.Message {
	mi := &file_tiles_proto_tiles_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Coord.ProtoReflect.Descriptor instead.
func (*Coord) Descriptor() ([]byte, []int) {
	return file_tiles_proto_tiles_proto_rawDescGZIP(), []int{1}
}

func (x *Coord) GetX() int32 {
	if x != nil {
		return x.X
	}
	return 0
}

func (x *Coord) GetY() int32 {
	if x != nil {
		return x.Y
	}
	return 0
}

type Placement struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Tile *Tile          `protobuf:"bytes,1,opt,name=tile,proto3" json:"tile,omitempty"`
	A    *Coord         `protobuf:"bytes,2,opt,name=a,proto3" json:"a,omitempty"`
	B    *Coord         `protobuf:"bytes,3,opt,name=b,proto3" json:"b,omitempty"`
	Type Placement_Type `protobuf:"varint,4,opt,name=type,proto3,enum=skelterjohn.tronimoes.tiles.Placement_Type" json:"type,omitempty"`
}

func (x *Placement) Reset() {
	*x = Placement{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tiles_proto_tiles_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Placement) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Placement) ProtoMessage() {}

func (x *Placement) ProtoReflect() protoreflect.Message {
	mi := &file_tiles_proto_tiles_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Placement.ProtoReflect.Descriptor instead.
func (*Placement) Descriptor() ([]byte, []int) {
	return file_tiles_proto_tiles_proto_rawDescGZIP(), []int{2}
}

func (x *Placement) GetTile() *Tile {
	if x != nil {
		return x.Tile
	}
	return nil
}

func (x *Placement) GetA() *Coord {
	if x != nil {
		return x.A
	}
	return nil
}

func (x *Placement) GetB() *Coord {
	if x != nil {
		return x.B
	}
	return nil
}

func (x *Placement) GetType() Placement_Type {
	if x != nil {
		return x.Type
	}
	return Placement_TYPE_UNKNOWN
}

type Line struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The first placement in a line is always the leader, including the round leader.
	Placements []*Placement `protobuf:"bytes,1,rep,name=placements,proto3" json:"placements,omitempty"`
	PlayerId   string       `protobuf:"bytes,2,opt,name=player_id,json=playerId,proto3" json:"player_id,omitempty"`
	Murderer   string       `protobuf:"bytes,3,opt,name=murderer,proto3" json:"murderer,omitempty"`
}

func (x *Line) Reset() {
	*x = Line{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tiles_proto_tiles_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Line) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Line) ProtoMessage() {}

func (x *Line) ProtoReflect() protoreflect.Message {
	mi := &file_tiles_proto_tiles_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Line.ProtoReflect.Descriptor instead.
func (*Line) Descriptor() ([]byte, []int) {
	return file_tiles_proto_tiles_proto_rawDescGZIP(), []int{3}
}

func (x *Line) GetPlacements() []*Placement {
	if x != nil {
		return x.Placements
	}
	return nil
}

func (x *Line) GetPlayerId() string {
	if x != nil {
		return x.PlayerId
	}
	return ""
}

func (x *Line) GetMurderer() string {
	if x != nil {
		return x.Murderer
	}
	return ""
}

type Player struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name          string  `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	PlayerId      string  `protobuf:"bytes,2,opt,name=player_id,json=playerId,proto3" json:"player_id,omitempty"`
	ChickenFooted bool    `protobuf:"varint,3,opt,name=chicken_footed,json=chickenFooted,proto3" json:"chicken_footed,omitempty"`
	Hand          []*Tile `protobuf:"bytes,4,rep,name=hand,proto3" json:"hand,omitempty"`
}

func (x *Player) Reset() {
	*x = Player{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tiles_proto_tiles_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Player) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Player) ProtoMessage() {}

func (x *Player) ProtoReflect() protoreflect.Message {
	mi := &file_tiles_proto_tiles_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Player.ProtoReflect.Descriptor instead.
func (*Player) Descriptor() ([]byte, []int) {
	return file_tiles_proto_tiles_proto_rawDescGZIP(), []int{4}
}

func (x *Player) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Player) GetPlayerId() string {
	if x != nil {
		return x.PlayerId
	}
	return ""
}

func (x *Player) GetChickenFooted() bool {
	if x != nil {
		return x.ChickenFooted
	}
	return false
}

func (x *Player) GetHand() []*Tile {
	if x != nil {
		return x.Hand
	}
	return nil
}

type Board struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Players      []*Player `protobuf:"bytes,1,rep,name=players,proto3" json:"players,omitempty"`
	PlayerLines  []*Line   `protobuf:"bytes,2,rep,name=player_lines,json=playerLines,proto3" json:"player_lines,omitempty"`
	FreeLines    []*Line   `protobuf:"bytes,3,rep,name=free_lines,json=freeLines,proto3" json:"free_lines,omitempty"`
	NextPlayerId string    `protobuf:"bytes,4,opt,name=next_player_id,json=nextPlayerId,proto3" json:"next_player_id,omitempty"`
	Bag          []*Tile   `protobuf:"bytes,5,rep,name=bag,proto3" json:"bag,omitempty"`
	Width        int32     `protobuf:"varint,6,opt,name=width,proto3" json:"width,omitempty"`
	Height       int32     `protobuf:"varint,7,opt,name=height,proto3" json:"height,omitempty"`
	Done         bool      `protobuf:"varint,8,opt,name=done,proto3" json:"done,omitempty"`
}

func (x *Board) Reset() {
	*x = Board{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tiles_proto_tiles_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Board) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Board) ProtoMessage() {}

func (x *Board) ProtoReflect() protoreflect.Message {
	mi := &file_tiles_proto_tiles_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Board.ProtoReflect.Descriptor instead.
func (*Board) Descriptor() ([]byte, []int) {
	return file_tiles_proto_tiles_proto_rawDescGZIP(), []int{5}
}

func (x *Board) GetPlayers() []*Player {
	if x != nil {
		return x.Players
	}
	return nil
}

func (x *Board) GetPlayerLines() []*Line {
	if x != nil {
		return x.PlayerLines
	}
	return nil
}

func (x *Board) GetFreeLines() []*Line {
	if x != nil {
		return x.FreeLines
	}
	return nil
}

func (x *Board) GetNextPlayerId() string {
	if x != nil {
		return x.NextPlayerId
	}
	return ""
}

func (x *Board) GetBag() []*Tile {
	if x != nil {
		return x.Bag
	}
	return nil
}

func (x *Board) GetWidth() int32 {
	if x != nil {
		return x.Width
	}
	return 0
}

func (x *Board) GetHeight() int32 {
	if x != nil {
		return x.Height
	}
	return 0
}

func (x *Board) GetDone() bool {
	if x != nil {
		return x.Done
	}
	return false
}

var File_tiles_proto_tiles_proto protoreflect.FileDescriptor

var file_tiles_proto_tiles_proto_rawDesc = []byte{
	0x0a, 0x17, 0x74, 0x69, 0x6c, 0x65, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x74, 0x69,
	0x6c, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1b, 0x73, 0x6b, 0x65, 0x6c, 0x74,
	0x65, 0x72, 0x6a, 0x6f, 0x68, 0x6e, 0x2e, 0x74, 0x72, 0x6f, 0x6e, 0x69, 0x6d, 0x6f, 0x65, 0x73,
	0x2e, 0x74, 0x69, 0x6c, 0x65, 0x73, 0x22, 0x22, 0x0a, 0x04, 0x54, 0x69, 0x6c, 0x65, 0x12, 0x0c,
	0x0a, 0x01, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x01, 0x61, 0x12, 0x0c, 0x0a, 0x01,
	0x62, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x01, 0x62, 0x22, 0x23, 0x0a, 0x05, 0x43, 0x6f,
	0x6f, 0x72, 0x64, 0x12, 0x0c, 0x0a, 0x01, 0x78, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x01,
	0x78, 0x12, 0x0c, 0x0a, 0x01, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x01, 0x79, 0x22,
	0xd5, 0x02, 0x0a, 0x09, 0x50, 0x6c, 0x61, 0x63, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x35, 0x0a,
	0x04, 0x74, 0x69, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x73, 0x6b,
	0x65, 0x6c, 0x74, 0x65, 0x72, 0x6a, 0x6f, 0x68, 0x6e, 0x2e, 0x74, 0x72, 0x6f, 0x6e, 0x69, 0x6d,
	0x6f, 0x65, 0x73, 0x2e, 0x74, 0x69, 0x6c, 0x65, 0x73, 0x2e, 0x54, 0x69, 0x6c, 0x65, 0x52, 0x04,
	0x74, 0x69, 0x6c, 0x65, 0x12, 0x30, 0x0a, 0x01, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x22, 0x2e, 0x73, 0x6b, 0x65, 0x6c, 0x74, 0x65, 0x72, 0x6a, 0x6f, 0x68, 0x6e, 0x2e, 0x74, 0x72,
	0x6f, 0x6e, 0x69, 0x6d, 0x6f, 0x65, 0x73, 0x2e, 0x74, 0x69, 0x6c, 0x65, 0x73, 0x2e, 0x43, 0x6f,
	0x6f, 0x72, 0x64, 0x52, 0x01, 0x61, 0x12, 0x30, 0x0a, 0x01, 0x62, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x22, 0x2e, 0x73, 0x6b, 0x65, 0x6c, 0x74, 0x65, 0x72, 0x6a, 0x6f, 0x68, 0x6e, 0x2e,
	0x74, 0x72, 0x6f, 0x6e, 0x69, 0x6d, 0x6f, 0x65, 0x73, 0x2e, 0x74, 0x69, 0x6c, 0x65, 0x73, 0x2e,
	0x43, 0x6f, 0x6f, 0x72, 0x64, 0x52, 0x01, 0x62, 0x12, 0x3f, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x2b, 0x2e, 0x73, 0x6b, 0x65, 0x6c, 0x74, 0x65, 0x72,
	0x6a, 0x6f, 0x68, 0x6e, 0x2e, 0x74, 0x72, 0x6f, 0x6e, 0x69, 0x6d, 0x6f, 0x65, 0x73, 0x2e, 0x74,
	0x69, 0x6c, 0x65, 0x73, 0x2e, 0x50, 0x6c, 0x61, 0x63, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x54,
	0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x22, 0x6c, 0x0a, 0x04, 0x54, 0x79, 0x70,
	0x65, 0x12, 0x10, 0x0a, 0x0c, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57,
	0x4e, 0x10, 0x00, 0x12, 0x11, 0x0a, 0x0d, 0x50, 0x4c, 0x41, 0x59, 0x45, 0x52, 0x5f, 0x4c, 0x45,
	0x41, 0x44, 0x45, 0x52, 0x10, 0x01, 0x12, 0x17, 0x0a, 0x13, 0x50, 0x4c, 0x41, 0x59, 0x45, 0x52,
	0x5f, 0x43, 0x4f, 0x4e, 0x54, 0x49, 0x4e, 0x55, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x10, 0x02, 0x12,
	0x0f, 0x0a, 0x0b, 0x46, 0x52, 0x45, 0x45, 0x5f, 0x4c, 0x45, 0x41, 0x44, 0x45, 0x52, 0x10, 0x03,
	0x12, 0x15, 0x0a, 0x11, 0x46, 0x52, 0x45, 0x45, 0x5f, 0x43, 0x4f, 0x4e, 0x54, 0x49, 0x4e, 0x55,
	0x41, 0x54, 0x49, 0x4f, 0x4e, 0x10, 0x04, 0x22, 0x87, 0x01, 0x0a, 0x04, 0x4c, 0x69, 0x6e, 0x65,
	0x12, 0x46, 0x0a, 0x0a, 0x70, 0x6c, 0x61, 0x63, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x26, 0x2e, 0x73, 0x6b, 0x65, 0x6c, 0x74, 0x65, 0x72, 0x6a, 0x6f,
	0x68, 0x6e, 0x2e, 0x74, 0x72, 0x6f, 0x6e, 0x69, 0x6d, 0x6f, 0x65, 0x73, 0x2e, 0x74, 0x69, 0x6c,
	0x65, 0x73, 0x2e, 0x50, 0x6c, 0x61, 0x63, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x0a, 0x70, 0x6c,
	0x61, 0x63, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x1b, 0x0a, 0x09, 0x70, 0x6c, 0x61, 0x79,
	0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x6c, 0x61,
	0x79, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x6d, 0x75, 0x72, 0x64, 0x65, 0x72, 0x65,
	0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6d, 0x75, 0x72, 0x64, 0x65, 0x72, 0x65,
	0x72, 0x22, 0x97, 0x01, 0x0a, 0x06, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x12, 0x12, 0x0a, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x12, 0x1b, 0x0a, 0x09, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x49, 0x64, 0x12, 0x25, 0x0a,
	0x0e, 0x63, 0x68, 0x69, 0x63, 0x6b, 0x65, 0x6e, 0x5f, 0x66, 0x6f, 0x6f, 0x74, 0x65, 0x64, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0d, 0x63, 0x68, 0x69, 0x63, 0x6b, 0x65, 0x6e, 0x46, 0x6f,
	0x6f, 0x74, 0x65, 0x64, 0x12, 0x35, 0x0a, 0x04, 0x68, 0x61, 0x6e, 0x64, 0x18, 0x04, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x21, 0x2e, 0x73, 0x6b, 0x65, 0x6c, 0x74, 0x65, 0x72, 0x6a, 0x6f, 0x68, 0x6e,
	0x2e, 0x74, 0x72, 0x6f, 0x6e, 0x69, 0x6d, 0x6f, 0x65, 0x73, 0x2e, 0x74, 0x69, 0x6c, 0x65, 0x73,
	0x2e, 0x54, 0x69, 0x6c, 0x65, 0x52, 0x04, 0x68, 0x61, 0x6e, 0x64, 0x22, 0xeb, 0x02, 0x0a, 0x05,
	0x42, 0x6f, 0x61, 0x72, 0x64, 0x12, 0x3d, 0x0a, 0x07, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x23, 0x2e, 0x73, 0x6b, 0x65, 0x6c, 0x74, 0x65, 0x72,
	0x6a, 0x6f, 0x68, 0x6e, 0x2e, 0x74, 0x72, 0x6f, 0x6e, 0x69, 0x6d, 0x6f, 0x65, 0x73, 0x2e, 0x74,
	0x69, 0x6c, 0x65, 0x73, 0x2e, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x52, 0x07, 0x70, 0x6c, 0x61,
	0x79, 0x65, 0x72, 0x73, 0x12, 0x44, 0x0a, 0x0c, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x5f, 0x6c,
	0x69, 0x6e, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x73, 0x6b, 0x65,
	0x6c, 0x74, 0x65, 0x72, 0x6a, 0x6f, 0x68, 0x6e, 0x2e, 0x74, 0x72, 0x6f, 0x6e, 0x69, 0x6d, 0x6f,
	0x65, 0x73, 0x2e, 0x74, 0x69, 0x6c, 0x65, 0x73, 0x2e, 0x4c, 0x69, 0x6e, 0x65, 0x52, 0x0b, 0x70,
	0x6c, 0x61, 0x79, 0x65, 0x72, 0x4c, 0x69, 0x6e, 0x65, 0x73, 0x12, 0x40, 0x0a, 0x0a, 0x66, 0x72,
	0x65, 0x65, 0x5f, 0x6c, 0x69, 0x6e, 0x65, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x21,
	0x2e, 0x73, 0x6b, 0x65, 0x6c, 0x74, 0x65, 0x72, 0x6a, 0x6f, 0x68, 0x6e, 0x2e, 0x74, 0x72, 0x6f,
	0x6e, 0x69, 0x6d, 0x6f, 0x65, 0x73, 0x2e, 0x74, 0x69, 0x6c, 0x65, 0x73, 0x2e, 0x4c, 0x69, 0x6e,
	0x65, 0x52, 0x09, 0x66, 0x72, 0x65, 0x65, 0x4c, 0x69, 0x6e, 0x65, 0x73, 0x12, 0x24, 0x0a, 0x0e,
	0x6e, 0x65, 0x78, 0x74, 0x5f, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x6e, 0x65, 0x78, 0x74, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72,
	0x49, 0x64, 0x12, 0x33, 0x0a, 0x03, 0x62, 0x61, 0x67, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x21, 0x2e, 0x73, 0x6b, 0x65, 0x6c, 0x74, 0x65, 0x72, 0x6a, 0x6f, 0x68, 0x6e, 0x2e, 0x74, 0x72,
	0x6f, 0x6e, 0x69, 0x6d, 0x6f, 0x65, 0x73, 0x2e, 0x74, 0x69, 0x6c, 0x65, 0x73, 0x2e, 0x54, 0x69,
	0x6c, 0x65, 0x52, 0x03, 0x62, 0x61, 0x67, 0x12, 0x14, 0x0a, 0x05, 0x77, 0x69, 0x64, 0x74, 0x68,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x77, 0x69, 0x64, 0x74, 0x68, 0x12, 0x16, 0x0a,
	0x06, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x68,
	0x65, 0x69, 0x67, 0x68, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x6f, 0x6e, 0x65, 0x18, 0x08, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x04, 0x64, 0x6f, 0x6e, 0x65, 0x42, 0x3b, 0x5a, 0x39, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x6b, 0x65, 0x6c, 0x74, 0x65, 0x72, 0x6a,
	0x6f, 0x68, 0x6e, 0x2f, 0x74, 0x72, 0x6f, 0x6e, 0x69, 0x6d, 0x6f, 0x65, 0x73, 0x2f, 0x73, 0x65,
	0x72, 0x76, 0x65, 0x72, 0x2f, 0x74, 0x69, 0x6c, 0x65, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x3b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_tiles_proto_tiles_proto_rawDescOnce sync.Once
	file_tiles_proto_tiles_proto_rawDescData = file_tiles_proto_tiles_proto_rawDesc
)

func file_tiles_proto_tiles_proto_rawDescGZIP() []byte {
	file_tiles_proto_tiles_proto_rawDescOnce.Do(func() {
		file_tiles_proto_tiles_proto_rawDescData = protoimpl.X.CompressGZIP(file_tiles_proto_tiles_proto_rawDescData)
	})
	return file_tiles_proto_tiles_proto_rawDescData
}

var file_tiles_proto_tiles_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_tiles_proto_tiles_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_tiles_proto_tiles_proto_goTypes = []interface{}{
	(Placement_Type)(0), // 0: skelterjohn.tronimoes.tiles.Placement.Type
	(*Tile)(nil),        // 1: skelterjohn.tronimoes.tiles.Tile
	(*Coord)(nil),       // 2: skelterjohn.tronimoes.tiles.Coord
	(*Placement)(nil),   // 3: skelterjohn.tronimoes.tiles.Placement
	(*Line)(nil),        // 4: skelterjohn.tronimoes.tiles.Line
	(*Player)(nil),      // 5: skelterjohn.tronimoes.tiles.Player
	(*Board)(nil),       // 6: skelterjohn.tronimoes.tiles.Board
}
var file_tiles_proto_tiles_proto_depIdxs = []int32{
	1,  // 0: skelterjohn.tronimoes.tiles.Placement.tile:type_name -> skelterjohn.tronimoes.tiles.Tile
	2,  // 1: skelterjohn.tronimoes.tiles.Placement.a:type_name -> skelterjohn.tronimoes.tiles.Coord
	2,  // 2: skelterjohn.tronimoes.tiles.Placement.b:type_name -> skelterjohn.tronimoes.tiles.Coord
	0,  // 3: skelterjohn.tronimoes.tiles.Placement.type:type_name -> skelterjohn.tronimoes.tiles.Placement.Type
	3,  // 4: skelterjohn.tronimoes.tiles.Line.placements:type_name -> skelterjohn.tronimoes.tiles.Placement
	1,  // 5: skelterjohn.tronimoes.tiles.Player.hand:type_name -> skelterjohn.tronimoes.tiles.Tile
	5,  // 6: skelterjohn.tronimoes.tiles.Board.players:type_name -> skelterjohn.tronimoes.tiles.Player
	4,  // 7: skelterjohn.tronimoes.tiles.Board.player_lines:type_name -> skelterjohn.tronimoes.tiles.Line
	4,  // 8: skelterjohn.tronimoes.tiles.Board.free_lines:type_name -> skelterjohn.tronimoes.tiles.Line
	1,  // 9: skelterjohn.tronimoes.tiles.Board.bag:type_name -> skelterjohn.tronimoes.tiles.Tile
	10, // [10:10] is the sub-list for method output_type
	10, // [10:10] is the sub-list for method input_type
	10, // [10:10] is the sub-list for extension type_name
	10, // [10:10] is the sub-list for extension extendee
	0,  // [0:10] is the sub-list for field type_name
}

func init() { file_tiles_proto_tiles_proto_init() }
func file_tiles_proto_tiles_proto_init() {
	if File_tiles_proto_tiles_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_tiles_proto_tiles_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Tile); i {
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
		file_tiles_proto_tiles_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Coord); i {
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
		file_tiles_proto_tiles_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Placement); i {
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
		file_tiles_proto_tiles_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Line); i {
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
		file_tiles_proto_tiles_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Player); i {
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
		file_tiles_proto_tiles_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Board); i {
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
			RawDescriptor: file_tiles_proto_tiles_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_tiles_proto_tiles_proto_goTypes,
		DependencyIndexes: file_tiles_proto_tiles_proto_depIdxs,
		EnumInfos:         file_tiles_proto_tiles_proto_enumTypes,
		MessageInfos:      file_tiles_proto_tiles_proto_msgTypes,
	}.Build()
	File_tiles_proto_tiles_proto = out.File
	file_tiles_proto_tiles_proto_rawDesc = nil
	file_tiles_proto_tiles_proto_goTypes = nil
	file_tiles_proto_tiles_proto_depIdxs = nil
}
