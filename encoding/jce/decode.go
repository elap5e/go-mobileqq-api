package jce

import (
	"math"
	"reflect"
)

func Unmarshal(data []byte, v interface{}) error {
	d := decoder{}
	d.init(data)
	return d.unmarshal(v)
}

type decoder struct {
	data []byte
	off  int
}

func (d *decoder) init(data []byte) *decoder {
	d.data = data
	d.off = 0
	return d
}

func (d *decoder) unmarshal(v interface{}) error {
	rv := reflect.ValueOf(v)
	_ = rv

	return nil
}

func (d *decoder) decodeHead() (uint8, uint8) {
	b := d.data[d.off]
	d.off++
	typ := (byte)(b & 15)
	tag := (b & 240) >> 4
	if tag != 15 {
		return typ, tag
	}
	tag = d.data[d.off] & 255
	d.off++
	return typ, tag
}

func (d *decoder) valueInterface(typ uint8) (val interface{}) {
	switch typ {
	default:
		panic("not implement")
	case 0x0c, 0x00, 0x01, 0x02, 0x03:
		val = d.uintInterface(typ)
	case 0x04:
		val = d.float32Interface()
	case 0x05:
		val = d.float64Interface()
	case 0x06, 0x07:
		val = d.stringInterface(typ)
	case 0x08:
		val = d.mapInterface()
	case 0x09:
		val = d.arrayInterface()
	case 0x0a, 0x0b:
		val = d.structInterface()
	case 0x0d:
		val = d.bytesInterface()
		// 	val = d.arrayInterface()
		// 	d.scanNext()
		// case scanBeginObject:
		// 	val = d.objectInterface()
		// 	d.scanNext()
		// case scanBeginLiteral:
		// 	val = d.literalInterface()
	}
	return
}

func (d *decoder) uintInterface(typ uint8) interface{} {
	switch typ {
	case 0x0c:
		return uint8(0x00)
	case 0x00:
		val := d.data[d.off]
		d.off++
		return val
	case 0x01:
		val := uint16(d.data[d.off]) << 8
		d.off++
		val += uint16(d.data[d.off])
		d.off++
		return val
	case 0x02:
		val := uint32(d.data[d.off]) << 24
		d.off++
		val += uint32(d.data[d.off]) << 16
		d.off++
		val += uint32(d.data[d.off]) << 8
		d.off++
		val += uint32(d.data[d.off])
		d.off++
		return val
	case 0x03:
		val := uint64(d.data[d.off]) << 56
		d.off++
		val += uint64(d.data[d.off]) << 48
		d.off++
		val += uint64(d.data[d.off]) << 40
		d.off++
		val += uint64(d.data[d.off]) << 32
		d.off++
		val += uint64(d.data[d.off]) << 24
		d.off++
		val += uint64(d.data[d.off]) << 16
		d.off++
		val += uint64(d.data[d.off]) << 8
		d.off++
		val += uint64(d.data[d.off])
		d.off++
		return val
	}
	panic("not implement")
}

func (d *decoder) float32Interface() interface{} {
	val := uint32(d.data[d.off]) << 24
	d.off++
	val += uint32(d.data[d.off]) << 16
	d.off++
	val += uint32(d.data[d.off]) << 8
	d.off++
	val += uint32(d.data[d.off])
	d.off++
	return math.Float32frombits(val)
}

func (d *decoder) float64Interface() interface{} {
	val := uint64(d.data[d.off]) << 56
	d.off++
	val += uint64(d.data[d.off]) << 48
	d.off++
	val += uint64(d.data[d.off]) << 40
	d.off++
	val += uint64(d.data[d.off]) << 32
	d.off++
	val += uint64(d.data[d.off]) << 24
	d.off++
	val += uint64(d.data[d.off]) << 16
	d.off++
	val += uint64(d.data[d.off]) << 8
	d.off++
	val += uint64(d.data[d.off])
	d.off++
	return math.Float64frombits(val)
}

func (d *decoder) stringInterface(typ uint8) interface{} {
	switch typ {
	case 0x06:
		l := int(d.data[d.off])
		d.off++
		val := d.data[d.off : d.off+l]
		d.off += l
		return val
	case 0x07:
		ttyp, _ := d.decodeHead()
		l := int(reflect.ValueOf(d.uintInterface(ttyp)).Uint())
		d.off++
		val := d.data[d.off : d.off+l]
		d.off += l
		return val
	}
	panic("not implement")
}

func (d *decoder) mapInterface() map[string]interface{} {
	m := make(map[string]interface{})
	ttyp, _ := d.decodeHead()
	l := int(reflect.ValueOf(d.uintInterface(ttyp)).Uint())
	for i := 0; i < l; i++ {
		ttyp, _ := d.decodeHead()
		key := d.stringInterface(ttyp).(string)
		ttyp, _ = d.decodeHead()
		m[key] = d.valueInterface(ttyp)
	}
	return m
}

func (d *decoder) arrayInterface() []interface{} {
	v := make([]interface{}, 0)
	ttyp, _ := d.decodeHead()
	l := int(reflect.ValueOf(d.uintInterface(ttyp)).Uint())
	for i := 0; i < l; i++ {
		ttyp, _ := d.decodeHead()
		v = append(v, d.valueInterface(ttyp))
	}
	return v
}

func (d *decoder) structInterface() interface{} {
	typ, _ := d.decodeHead()
	for typ != 0x0b {
		val := d.valueInterface(typ)
		_ = val
		typ, _ = d.decodeHead()
	}
	return nil
}

func (d *decoder) bytesInterface() interface{} {
	_, _ = d.decodeHead()
	ttyp, _ := d.decodeHead()
	n := int(reflect.ValueOf(d.uintInterface(ttyp)).Uint())
	val := d.data[d.off : d.off+n]
	d.off += n
	return val
}
