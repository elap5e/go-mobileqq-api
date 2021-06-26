package jce

import (
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
	op   int
	scan scanner
}

func (d *decoder) init(data []byte) *decoder {
	d.data = data
	d.off = 0
	return d
}

func (d *decoder) value(v reflect.Value) error {
	switch d.op {
	default:
		panic("phasePanicMsg")
	case scanBeginArray:
		// if v.IsValid() {
		// 	if err := d.array(v); err != nil {
		// 		return err
		// 	}
		// } else {
		// 	d.skip()
		// }
		// d.scanNext()
	case scanBeginObject:
		// if v.IsValid() {
		// 	if err := d.object(v); err != nil {
		// 		return err
		// 	}
		// } else {
		// 	d.skip()
		// }
		// d.scanNext()
	case scanBeginLiteral:
		// All bytes inside literal return scanContinue op code.
		// start := d.readIndex()
		// d.rescanLiteral()

		// if v.IsValid() {
		// 	if err := d.literalStore(d.data[start:d.readIndex()], v, false); err != nil {
		// 		return err
		// 	}
		// }
	}
	return nil
}

func (d *decoder) unmarshal(v interface{}) error {
	rv := reflect.ValueOf(v)
	d.scan.reset()
	// d.scanWhile(scanSkipSpace)
	err := d.value(rv)
	if err != nil {
		return err
	}
	return nil
}
