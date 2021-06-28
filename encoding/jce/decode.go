package jce

import (
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
	case reflect.Ptr:
	case reflect.Slice:
		return d.decodeSlice(v)
	case reflect.Struct:
		return d.decodeStruct(v)
	case reflect.Array:
	case reflect.Map:
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

func (d *decoder) decodeBytes(v reflect.Value) error {
	ttyp, _ := d.decodeHead()
	n := int(d.decodeUint(ttyp))
	v.SetBytes(d.data[d.off : d.off+n])
	d.off += n
	return nil
}

func (d *decoder) decodeSlice(v reflect.Value) error {
	if v.Elem().Kind() == reflect.Uint8 {
		return d.decodeBytes(v)
	}
	return nil
}

func (d *decoder) decodeStruct(v reflect.Value) error {
	tv := v.Type()

	var fields structFields

	switch v.Kind() {
	case reflect.Struct:
		fields = typeFields(tv)
	}

	var subv reflect.Value

	typ, tag := d.decodeHead()
	for {
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
		if d.off == d.length {
			break
		} else {
			typ, tag = d.decodeHead()
			if tag == 0x0b {
				break
			}
		}
	}
	return nil
}

// func (d *decoder) arrayInterface() []interface{} {
// 	v := make([]interface{}, 0)
// 	ttyp, _ := d.decodeHead()
// 	l := int(reflect.ValueOf(d.uintInterface(ttyp)).Uint())
// 	for i := 0; i < l; i++ {
// 		ttyp, _ := d.decodeHead()
// 		v = append(v, d.valueInterface(ttyp))
// 	}
// 	return v
// }

func (d *decoder) decodeMap(typ uint8) map[string]interface{} {
	if typ != 0x08 {
		log.Panicf("unexpected type 0x%02x (decode map)", typ)
	}
	m := make(map[string]interface{})
	ttyp, _ := d.decodeHead()
	l := int(d.decodeUint(ttyp))
	for i := 0; i < l; i++ {
		ttyp, _ = d.decodeHead()
		key := d.decodeString(ttyp)
		ttyp, _ = d.decodeHead()
		m[key] = d.decodeString(ttyp)
	}
	return m
}

func (d *decoder) decodeString(typ uint8) string {
	var l int
	switch typ {
	default:
		log.Panicf("unexpected type 0x%02x (decode string)", typ)
	case 0x07:
		ttyp, _ := d.decodeHead()
		l = int(d.decodeUint(ttyp))
		d.off++
	case 0x06:
		l = int(d.data[d.off])
		d.off++
	}
	val := string(d.data[d.off : d.off+l])
	d.off += l
	return val
}

func (d *decoder) decodeFloat(typ uint8) float64 {
	var val uint64
	switch typ {
	default:
		log.Panicf("unexpected type 0x%02x (decode float)", typ)
	case 0x05:
		val = val<<8 + uint64(d.data[d.off])
		d.off++
		val = val<<8 + uint64(d.data[d.off])
		d.off++
		val = val<<8 + uint64(d.data[d.off])
		d.off++
		val = val<<8 + uint64(d.data[d.off])
		d.off++
		fallthrough
	case 0x04:
		val = val<<8 + uint64(d.data[d.off])
		d.off++
		val = val<<8 + uint64(d.data[d.off])
		d.off++
		val = val<<8 + uint64(d.data[d.off])
		d.off++
		val = val<<8 + uint64(d.data[d.off])
		d.off++
	}
	return math.Float64frombits(val)
}

func (d *decoder) decodeUint(typ uint8) uint64 {
	var val uint64
	switch typ {
	default:
		log.Panicf("unexpected type 0x%02x (decode uint)", typ)
	case 0x03:
		val = val<<8 + uint64(d.data[d.off])
		d.off++
		val = val<<8 + uint64(d.data[d.off])
		d.off++
		val = val<<8 + uint64(d.data[d.off])
		d.off++
		val = val<<8 + uint64(d.data[d.off])
		d.off++
		fallthrough
	case 0x02:
		val = val<<8 + uint64(d.data[d.off])
		d.off++
		val = val<<8 + uint64(d.data[d.off])
		d.off++
		fallthrough
	case 0x01:
		val = val<<8 + uint64(d.data[d.off])
		d.off++
		fallthrough
	case 0x00:
		val = val<<8 + uint64(d.data[d.off])
		d.off++
		fallthrough
	case 0x0c:
	}
	return val
}
