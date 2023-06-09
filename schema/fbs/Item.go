// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package fbs

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type ItemT struct {
	Name string `json:"name"`
}

func (t *ItemT) Pack(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	if t == nil { return 0 }
	nameOffset := flatbuffers.UOffsetT(0)
	if t.Name != "" {
		nameOffset = builder.CreateString(t.Name)
	}
	ItemStart(builder)
	ItemAddName(builder, nameOffset)
	return ItemEnd(builder)
}

func (rcv *Item) UnPackTo(t *ItemT) {
	t.Name = string(rcv.Name())
}

func (rcv *Item) UnPack() *ItemT {
	if rcv == nil { return nil }
	t := &ItemT{}
	rcv.UnPackTo(t)
	return t
}

type Item struct {
	_tab flatbuffers.Table
}

func GetRootAsItem(buf []byte, offset flatbuffers.UOffsetT) *Item {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Item{}
	x.Init(buf, n+offset)
	return x
}

func GetSizePrefixedRootAsItem(buf []byte, offset flatbuffers.UOffsetT) *Item {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &Item{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func (rcv *Item) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Item) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *Item) Name() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func ItemStart(builder *flatbuffers.Builder) {
	builder.StartObject(1)
}
func ItemAddName(builder *flatbuffers.Builder, name flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(name), 0)
}
func ItemEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
