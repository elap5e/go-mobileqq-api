package jce

import (
	"encoding/hex"
	"log"
	"math"
	"reflect"
)

func Unmarshal(data []byte, v interface{}, opts ...bool) error {
	simple := false
	if len(opts) != 0 && opts[0] {
		simple = true
	}

	d := decoder{
		data:   data,
		off:    0,
		length: len(data),
	}

	return d.unmarshal(v, simple)
}

type decoder struct {
	data   []byte
	off    int
	length int
	typ    uint8
}

func (d *decoder) unmarshal(v interface{}, simple bool) error {
	if !simple {
		_, _ = d.decodeHead()
	}
	return d.decodeValue(reflect.ValueOf(v))
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

func (d *decoder) decodeValue(v reflect.Value) error {
	switch v.Kind() {
	case reflect.Interface:
		return d.decodeInterface(v)
	case reflect.Ptr:
		return d.decodePtr(v)
	case reflect.Slice:
		return d.decodeSlice(v)
	case reflect.Struct:
		return d.decodeStruct(v)
	case reflect.Array:
		return d.decodeArray(v)
	case reflect.Map:
		return d.decodeMap(v)
	case reflect.String:
		v.SetString(d.decodeString(d.typ))
	case reflect.Float64, reflect.Float32:
		v.SetFloat(d.decodeFloat(d.typ))
	case reflect.Uint64, reflect.Uint32, reflect.Uint, reflect.Uint16, reflect.Uint8:
		v.SetUint(d.decodeUint(d.typ))
	case reflect.Bool:
		v.SetBool(d.decodeUint(d.typ) != 0)
	}
	return nil
}

func (d *decoder) decodeInterface(v reflect.Value) error {
	return d.decodeValue(v.Elem())
}

func (d *decoder) decodePtr(v reflect.Value) error {
	if v.IsNil() {
		v.Set(reflect.New(v.Type().Elem()))
	}
	return d.decodeValue(v.Elem())
}

func (d *decoder) decodeBytes(v reflect.Value) error {
	_, _ = d.decodeHead()
	ttyp, _ := d.decodeHead()
	n := int(d.decodeUint(ttyp))
	v.SetBytes(d.data[d.off : d.off+n])
	d.off += n
	return nil
}

func (d *decoder) decodeSlice(v reflect.Value) error {
	if v.Type().Elem().Kind() == reflect.Uint8 {
		return d.decodeBytes(v)
	}
	return d.decodeArray(v)
}

func (d *decoder) decodeStruct(v reflect.Value) error {
	tv := v.Type()

	var fields structFields

	switch v.Kind() {
	case reflect.Ptr:
		return d.decodePtr(v)
	case reflect.Struct:
		fields = typeFields(tv)
	}

	var subv reflect.Value

	for d.off < d.length {
		typ, tag := d.decodeHead()
		if typ == 0x0b {
			break
		}
		var f *field
		if i, ok := fields.tagIndex[tag]; ok {
			f = &fields.list[i]
		}
		if f != nil {
			subv = v
			for _, i := range f.index {
				if subv.Kind() == reflect.Ptr {
					if subv.IsNil() {
						subv.Set(reflect.New(subv.Type().Elem()))
					}
					subv = subv.Elem()
				}
				subv = subv.Field(i)
			}
		}
		d.typ = typ
		if err := d.decodeValue(subv); err != nil {
			return err
		}
	}
	return nil
}

func (d *decoder) decodeArray(v reflect.Value) error {
	ttyp, _ := d.decodeHead()
	n := int(d.decodeUint(ttyp))
	if v.Kind() == reflect.Slice {
		if n >= v.Cap() {
			newv := reflect.MakeSlice(v.Type(), v.Len(), n)
			reflect.Copy(newv, v)
			v.Set(newv)
		}
		if n >= v.Len() {
			v.SetLen(n)
		}
	}
	for i := 0; i < n; i++ {
		if i < v.Len() {
			d.typ, _ = d.decodeHead()
			if err := d.decodeValue(v.Index(i)); err != nil {
				return err
			}
		}
	}
	return nil
}

func (d *decoder) decodeMap(v reflect.Value) error {
	typ, _ := d.decodeHead()
	n := int(d.decodeUint(typ))
	t := v.Type()
	if v.IsNil() {
		v.Set(reflect.MakeMap(t))
	}
	var subv reflect.Value
	for i := 0; i < n; i++ {
		ttyp, _ := d.decodeHead()
		key := d.decodeString(ttyp)
		subv = reflect.New(t.Elem()).Elem()
		d.typ, _ = d.decodeHead()
		if err := d.decodeValue(subv); err != nil {
			return err
		}
		v.SetMapIndex(reflect.ValueOf(key), subv)
	}
	return nil
}

func (d *decoder) decodeString(typ uint8) string {
	var l int
	switch typ {
	default:
		log.Fatalf(
			"(decode string) unexpected type 0x%02x 0x%08x dump\n%s",
			typ, d.off, hex.Dump(d.data),
		)
	case 0x07:
		ttyp, _ := d.decodeHead()
		l = int(d.decodeUint(ttyp))
	case 0x06:
		l = int(d.data[d.off])
		d.off++
	}
	val := string(d.data[d.off : d.off+l])
	d.off += l
	return val
}

func (d *decoder) decodeFloat(typ uint8) float64 {
	switch typ {
	default:
		log.Fatalf(
			"(decode float) unexpected type 0x%02x 0x%08x dump\n%s",
			typ, d.off, hex.Dump(d.data),
		)
	case 0x05:
		val := uint64(d.data[d.off])<<24 + uint64(d.data[d.off+1])<<16 + uint64(d.data[d.off+2])<<8 + uint64(d.data[d.off+3])
		d.off += 4
		val = val<<32 + uint64(d.data[d.off])<<24 + uint64(d.data[d.off+1])<<16 + uint64(d.data[d.off+2])<<8 + uint64(d.data[d.off+3])
		d.off += 4
		return math.Float64frombits(val)
	case 0x04:
		val := uint32(d.data[d.off])<<24 + uint32(d.data[d.off+1])<<16 + uint32(d.data[d.off+2])<<8 + uint32(d.data[d.off+3])
		d.off += 4
		return float64(math.Float32frombits(val))
	}
	return 0
}

func (d *decoder) decodeUint(typ uint8) uint64 {
	var val uint64
	switch typ {
	default:
		log.Fatalf(
			"(decode uint) unexpected type 0x%02x 0x%08x dump\n%s",
			typ, d.off, hex.Dump(d.data),
		)
	case 0x03:
		val = uint64(d.data[d.off])<<24 + uint64(d.data[d.off+1])<<16 + uint64(d.data[d.off+2])<<8 + uint64(d.data[d.off+3])
		d.off += 4
		fallthrough
	case 0x02:
		val = val<<16 + uint64(d.data[d.off])<<8 + uint64(d.data[d.off+1])
		d.off += 2
		fallthrough
	case 0x01:
		val = val<<8 + uint64(d.data[d.off])
		d.off++
		fallthrough
	case 0x00:
		val = val<<8 + uint64(d.data[d.off])
		d.off++
	case 0x0c:
	}
	return val
}
