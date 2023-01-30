package k6utils

import (
	"bytes"
	"encoding/binary"
)

type FidoUafTlvObject struct {
	backBuffer *bytes.Buffer
}

func NewFidoUafTlvObject(tag int16, value []byte) *FidoUafTlvObject {
	t := &FidoUafTlvObject{}
	t.initBuffer(tag, int16(len(value)))
	t.PutInt8Array(value)
	return t
}

func NewFidoUafTlvObjectFromArray(tag int16, values ...*FidoUafTlvObject) *FidoUafTlvObject {
	t := &FidoUafTlvObject{}
	var totalSize int16
	for _, value := range values {
		totalSize += int16(len(value.GetByteArray()))
	}

	t.initBuffer(tag, totalSize)

	for _, value := range values {
		t.PutInt8Array(value.GetByteArray())
	}
	return t
}

func NewFidoUafTlvObjectWithSize(tag int16, contentSize int16) *FidoUafTlvObject {
	t := &FidoUafTlvObject{}
	t.initBuffer(tag, contentSize)
	return t
}

func (t *FidoUafTlvObject) PutInt32(value int32) *FidoUafTlvObject {
	binary.Write(t.backBuffer, binary.LittleEndian, value)
	return t
}

func (t *FidoUafTlvObject) PutInt16(value int16) *FidoUafTlvObject {
	binary.Write(t.backBuffer, binary.LittleEndian, value)
	return t
}

func (t *FidoUafTlvObject) PutInt8(value int8) *FidoUafTlvObject {
	binary.Write(t.backBuffer, binary.LittleEndian, value)
	return t
}

func (t *FidoUafTlvObject) PutInt8Array(value []byte) *FidoUafTlvObject {
	t.backBuffer.Write(value)
	return t
}

func (t *FidoUafTlvObject) GetByteArray() []byte {
	return t.backBuffer.Bytes()
}

func (t *FidoUafTlvObject) initBuffer(tag int16, contentSize int16) *FidoUafTlvObject {
	t.backBuffer = new(bytes.Buffer)
	t.backBuffer.Write(make([]byte, 0, 2+2+int(contentSize)))
	binary.Write(t.backBuffer, binary.LittleEndian, tag)
	binary.Write(t.backBuffer, binary.LittleEndian, contentSize)
	return t
}
