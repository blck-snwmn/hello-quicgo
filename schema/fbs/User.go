// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package fbs

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type UserT struct {
	Name string `json:"name"`
	Pos *PositionT `json:"pos"`
	Color Color `json:"color"`
	Inventory []*ItemT `json:"inventory"`
}

func (t *UserT) Pack(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	if t == nil { return 0 }
	nameOffset := flatbuffers.UOffsetT(0)
	if t.Name != "" {
		nameOffset = builder.CreateString(t.Name)
	}
	inventoryOffset := flatbuffers.UOffsetT(0)
	if t.Inventory != nil {
		inventoryLength := len(t.Inventory)
		inventoryOffsets := make([]flatbuffers.UOffsetT, inventoryLength)
		for j := 0; j < inventoryLength; j++ {
			inventoryOffsets[j] = t.Inventory[j].Pack(builder)
		}
		UserStartInventoryVector(builder, inventoryLength)
		for j := inventoryLength - 1; j >= 0; j-- {
			builder.PrependUOffsetT(inventoryOffsets[j])
		}
		inventoryOffset = builder.EndVector(inventoryLength)
	}
	UserStart(builder)
	UserAddName(builder, nameOffset)
	posOffset := t.Pos.Pack(builder)
	UserAddPos(builder, posOffset)
	UserAddColor(builder, t.Color)
	UserAddInventory(builder, inventoryOffset)
	return UserEnd(builder)
}

func (rcv *User) UnPackTo(t *UserT) {
	t.Name = string(rcv.Name())
	t.Pos = rcv.Pos(nil).UnPack()
	t.Color = rcv.Color()
	inventoryLength := rcv.InventoryLength()
	t.Inventory = make([]*ItemT, inventoryLength)
	for j := 0; j < inventoryLength; j++ {
		x := Item{}
		rcv.Inventory(&x, j)
		t.Inventory[j] = x.UnPack()
	}
}

func (rcv *User) UnPack() *UserT {
	if rcv == nil { return nil }
	t := &UserT{}
	rcv.UnPackTo(t)
	return t
}

type User struct {
	_tab flatbuffers.Table
}

func GetRootAsUser(buf []byte, offset flatbuffers.UOffsetT) *User {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &User{}
	x.Init(buf, n+offset)
	return x
}

func GetSizePrefixedRootAsUser(buf []byte, offset flatbuffers.UOffsetT) *User {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &User{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func (rcv *User) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *User) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *User) Name() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *User) Pos(obj *Position) *Position {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		x := o + rcv._tab.Pos
		if obj == nil {
			obj = new(Position)
		}
		obj.Init(rcv._tab.Bytes, x)
		return obj
	}
	return nil
}

func (rcv *User) Color() Color {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return Color(rcv._tab.GetInt8(o + rcv._tab.Pos))
	}
	return 2
}

func (rcv *User) MutateColor(n Color) bool {
	return rcv._tab.MutateInt8Slot(8, int8(n))
}

func (rcv *User) Inventory(obj *Item, j int) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		x := rcv._tab.Vector(o)
		x += flatbuffers.UOffsetT(j) * 4
		x = rcv._tab.Indirect(x)
		obj.Init(rcv._tab.Bytes, x)
		return true
	}
	return false
}

func (rcv *User) InventoryLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func UserStart(builder *flatbuffers.Builder) {
	builder.StartObject(4)
}
func UserAddName(builder *flatbuffers.Builder, name flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(name), 0)
}
func UserAddPos(builder *flatbuffers.Builder, pos flatbuffers.UOffsetT) {
	builder.PrependStructSlot(1, flatbuffers.UOffsetT(pos), 0)
}
func UserAddColor(builder *flatbuffers.Builder, color Color) {
	builder.PrependInt8Slot(2, int8(color), 2)
}
func UserAddInventory(builder *flatbuffers.Builder, inventory flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(3, flatbuffers.UOffsetT(inventory), 0)
}
func UserStartInventoryVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func UserEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}