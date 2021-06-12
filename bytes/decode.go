package bytes

import (
	"io"
)

// DecodeUint8 consumes an encoded unsigned 8-bit integer from the buffer.
func (b *Buffer) DecodeUint8() (uint8, error) {
	v := b.buf[b.idx:]
	if len(v) < 1 {
		return 0, io.ErrUnexpectedEOF
	}
	b.idx += 1
	return v[0], nil
}

// DecodeInt8 consumes an encoded signed 8-bit integer from the buffer.
func (b *Buffer) DecodeInt8() (int8, error) {
	v, err := b.DecodeUint8()
	return int8(v), err
}

// DecodeByte consumes an encoded raw byte from the buffer.
func (b *Buffer) DecodeByte() (byte, error) {
	return b.DecodeUint8()
}

// DecodeBool consumes an encoded raw boolean from the buffer.
func (b *Buffer) DecodeBool() (bool, error) {
	v, err := b.DecodeUint8()
	if err != nil {
		return false, err
	}
	if v == 0x00 {
		return false, nil
	} else {
		return true, nil
	}
}

// DecodeUint16 consumes an encoded unsigned 16-bit big-endian integer from the buffer.
func (b *Buffer) DecodeUint16() (uint16, error) {
	v := b.buf[b.idx:]
	if len(v) < 2 {
		return 0, io.ErrUnexpectedEOF
	}
	b.idx += 2
	return uint16(v[0])<<8 | uint16(v[1])<<0, nil
}

// DecodeUint16 consumes an encoded signed 16-bit big-endian integer from the buffer.
func (b *Buffer) DecodeInt16() (int16, error) {
	v, err := b.DecodeUint16()
	return int16(v), err
}

// DecodeUint32 consumes an encoded unsigned 32-bit big-endian integer from the buffer.
func (b *Buffer) DecodeUint32() (uint32, error) {
	v := b.buf[b.idx:]
	if len(v) < 4 {
		return 0, io.ErrUnexpectedEOF
	}
	b.idx += 4
	return uint32(v[0])<<24 | uint32(v[1])<<16 | uint32(v[2])<<8 | uint32(v[3])<<0, nil
}

// DecodeInt32 consumes an encoded signed 32-bit big-endian integer from the buffer.
func (b *Buffer) DecodeInt32() (int32, error) {
	v, err := b.DecodeUint32()
	return int32(v), err
}

// DecodeRune consumes an encoded rune from the buffer.
func (b *Buffer) DecodeRune() (rune, error) {
	v, err := b.DecodeUint32()
	return rune(v), err
}

// DecodeUint64 consumes an encoded unsigned 64-bit big-endian integer from the buffer.
func (b *Buffer) DecodeUint64() (uint64, error) {
	v := b.buf[b.idx:]
	if len(v) < 8 {
		return 0, io.ErrUnexpectedEOF
	}
	b.idx += 8
	return uint64(v[0])<<56 | uint64(v[1])<<48 | uint64(v[2])<<40 | uint64(v[3])<<32 | uint64(v[4])<<24 | uint64(v[5])<<16 | uint64(v[6])<<8 | uint64(v[7])<<0, nil
}

// DecodeInt64 consumes an encoded signed 64-bit big-endian integer from the buffer.
func (b *Buffer) DecodeInt64() (int64, error) {
	v, err := b.DecodeUint64()
	return int64(v), err
}

// DecodeBytes consumes an encoded length-prefixed raw bytes from the buffer.
func (b *Buffer) DecodeBytes() ([]byte, error) {
	n, err := b.DecodeUint16()
	if err != nil {
		return nil, err
	}
	v := b.buf[b.idx:]
	if len(v) < int(n) {
		return nil, io.ErrUnexpectedEOF
	}
	b.idx += int(n)
	return v[:n], nil
}

// DecodeString consumes an encoded length-prefixed raw string from the buffer.
func (b *Buffer) DecodeString() (string, error) {
	v, err := b.DecodeBytes()
	return string(v), err
}

// DecodeBytesN consumes an n-limited encoded raw bytes to the buffer.
func (b *Buffer) DecodeBytesN(n uint16) ([]byte, error) {
	v := b.buf[b.idx:]
	if len(v) < int(n) {
		return nil, io.ErrUnexpectedEOF
	}
	b.idx += int(n)
	return v[:n], nil
}

// DecodeStringN consumes an n-limited encoded raw string to the buffer.
func (b *Buffer) DecodeStringN(n uint16) (string, error) {
	v, err := b.DecodeBytesN(n)
	return string(v), err
}
