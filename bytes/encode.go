package bytes

import (
	"encoding/binary"
	"math"
)

func (b *Buffer) Write(p []byte) {
	b.buf = append(b.buf, p...)
}

func (b *Buffer) WriteUint8(u uint8) {
	b.buf = append(b.buf, u)
}

func (b *Buffer) WriteUint16(u uint16) {
	t := make([]byte, 2)
	binary.BigEndian.PutUint16(t, u)
	b.buf = append(b.buf, t...)
}

func (b *Buffer) WriteUint32(u uint32) {
	t := make([]byte, 4)
	binary.BigEndian.PutUint32(t, u)
	b.buf = append(b.buf, t...)
}

func (b *Buffer) WriteUint32At(u uint32, off int64) {
	if g := off + 4 - int64(len(b.buf)); g > 0 {
		b.buf = append(b.buf, make([]byte, g)...)
	}
	binary.BigEndian.PutUint32(b.buf[off:], u)
}

func (b *Buffer) WriteUint16Bytes(p []byte) {
	b.WriteUint16(uint16(len(p) + 2))
	b.Write(p)
}

func (b *Buffer) WriteUint16String(s string) {
	b.WriteUint16Bytes([]byte(s))
}

func (b *Buffer) WriteUint16LengthBytes(p []byte) {
	b.WriteUint16(uint16(len(p)))
	b.Write(p)
}

func (b *Buffer) WriteUint32Bytes(p []byte) {
	b.WriteUint32(uint32(len(p) + 4))
	b.Write(p)
}

func (b *Buffer) WriteUint32String(s string) {
	b.WriteUint32Bytes([]byte(s))
}

// EncodeUint8 appends an unsigned 8-bit integer to the buffer.
func (b *Buffer) EncodeUint8(v uint8) {
	b.buf = append(b.buf, v)
}

// EncodeInt8 appends a signed 8-bit integer to the buffer.
func (b *Buffer) EncodeInt8(v int8) {
	b.EncodeUint8(uint8(v))
}

// EncodeByte appends a raw byte to the buffer.
func (b *Buffer) EncodeByte(v byte) {
	b.EncodeUint8(v)
}

// EncodeBool appends a raw boolean to the buffer.
func (b *Buffer) EncodeBool(v bool) {
	if v {
		b.EncodeUint8(0x01)
	} else {
		b.EncodeUint8(0x00)
	}
}

// EncodeUint16 appends an unsigned 16-bit big-endian integer to the buffer.
func (b *Buffer) EncodeUint16(v uint16) {
	t := make([]byte, 2)
	binary.BigEndian.PutUint16(t, v)
	b.buf = append(b.buf, t...)
}

// EncodeUint16 appends a signed 16-bit big-endian integer to the buffer.
func (b *Buffer) EncodeInt16(v int16) {
	b.EncodeUint16(uint16(v))
}

// EncodeUint32 appends an unsigned 32-bit big-endian integer to the buffer.
func (b *Buffer) EncodeUint32(u uint32) {
	b.WriteUint32(u)
}

// EncodeInt32 appends a signed 32-bit big-endian integer to the buffer.
func (b *Buffer) EncodeInt32(v int32) {
	b.EncodeUint32(uint32(v))
}

// EncodeRune appends a rune to the buffer.
func (b *Buffer) EncodeRune(v rune) {
	b.EncodeUint32(uint32(v))
}

// EncodeUint64 appends an unsigned 64-bit big-endian integer to the buffer.
func (b *Buffer) EncodeUint64(v uint64) {
	t := make([]byte, 8)
	binary.BigEndian.PutUint64(t, v)
	b.buf = append(b.buf, t...)
}

// EncodeInt64 appends a signed 64-bit big-endian integer to the buffer.
func (b *Buffer) EncodeInt64(v int64) {
	b.EncodeUint64(uint64(v))
}

// EncodeFloat32 appends a float 32 big-endian integer to the buffer.
func (b *Buffer) EncodeFloat32(v float32) {
	b.EncodeUint32(math.Float32bits(v))
}

// EncodeFloat64 appends a float 64 big-endian integer to the buffer.
func (b *Buffer) EncodeFloat64(v float64) {
	b.EncodeUint64(math.Float64bits(v))
}

// EncodeBytes appends a length-prefixed raw bytes to the buffer.
func (b *Buffer) EncodeBytes(v []byte) {
	b.EncodeUint16(uint16(len(v)))
	b.EncodeRawBytes(v)
}

// EncodeString appends a length-prefixed raw string to the buffer.
func (b *Buffer) EncodeString(v string) {
	b.EncodeBytes([]byte(v))
}

// EncodeBytesN appends an n-limited length-prefixed raw bytes to the buffer.
func (b *Buffer) EncodeBytesN(v []byte, n uint16) {
	l := uint16(len(v))
	if l > n {
		l = n
	}
	b.EncodeUint16(l)
	b.EncodeRawBytes(v[:l])
}

// EncodeStringN appends an n-limited length-prefixed raw string to the buffer.
func (b *Buffer) EncodeStringN(v string, n uint16) {
	b.EncodeBytesN([]byte(v), n)
}

// EncodeRawBytes appends a raw bytes to the buffer.
func (b *Buffer) EncodeRawBytes(v []byte) {
	b.buf = append(b.buf, v...)
}

// EncodeRawString appends a raw string to the buffer.
func (b *Buffer) EncodeRawString(v string) {
	b.EncodeRawBytes([]byte(v))
}

func (b *Buffer) Reset() {
	b.buf = b.buf[:0]
	b.idx = 0
}
