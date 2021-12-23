package rawpkt

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
	"time"
)

// Marshal content to raw byte data
func (p *Packet) Marshal(structPtr interface{}) error {
	vp := reflect.ValueOf(structPtr)
	return p.deserialize(vp)
}

func (p *Packet) deserialize(fv reflect.Value) error {
	switch fv.Kind() {
	case reflect.Int:
		b := make([]byte, 4)
		binary.LittleEndian.PutUint32(b, uint32(fv.Int()))
		p.data = append(p.data, b...)
		p.addSize(4)

	case reflect.Uint:
		b := make([]byte, 4)
		binary.LittleEndian.PutUint32(b, uint32(fv.Uint()))
		p.data = append(p.data, b...)
		p.addSize(4)

	case reflect.Int8:
		p.data = append(p.data, byte(fv.Int()))
		p.addSize(1)

	case reflect.Int16:
		b := make([]byte, 2)
		binary.LittleEndian.PutUint16(b, uint16(fv.Int()))
		p.data = append(p.data, b...)
		p.addSize(2)

	case reflect.Int32:
		b := make([]byte, 4)
		binary.LittleEndian.PutUint32(b, uint32(fv.Int()))
		p.data = append(p.data, b...)
		p.addSize(4)

	case reflect.Int64:
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, uint64(fv.Int()))
		p.data = append(p.data, b...)
		p.addSize(8)

	case reflect.Uint8:
		p.data = append(p.data, byte(fv.Uint()))
		p.addSize(1)

	case reflect.Uint16:
		b := make([]byte, 2)
		binary.LittleEndian.PutUint16(b, uint16(fv.Uint()))
		p.data = append(p.data, b...)
		p.addSize(2)

	case reflect.Uint32:
		b := make([]byte, 4)
		binary.LittleEndian.PutUint32(b, uint32(fv.Uint()))
		p.data = append(p.data, b...)
		p.addSize(4)

	case reflect.Uint64:
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, uint64(fv.Uint()))
		p.data = append(p.data, b...)
		p.addSize(8)

	case reflect.Bool:
		var bitSet = byte(0)
		if fv.Bool() {
			bitSet = 1
		}
		p.data = append(p.data, bitSet)
		p.addSize(1)

	case reflect.Float32:
		var b bytes.Buffer
		err := binary.Write(&b, binary.LittleEndian, float32(fv.Float()))
		if err == nil {
			p.data = append(p.data, b.Bytes()...)
			p.addSize(4)
		}

	case reflect.Float64:
		var b bytes.Buffer
		err := binary.Write(&b, binary.LittleEndian, fv.Float())
		if err == nil {
			p.data = append(p.data, b.Bytes()...)
			p.addSize(8)
		}

	case reflect.Array:
		arrLen := uint16(fv.Len())
		for i := 0; i < int(arrLen); i++ {
			p.deserialize(fv.Index(i))
		}

	case reflect.String:
		b := []byte(fv.String())
		lenBinary := make([]byte, 4)
		binary.LittleEndian.PutUint32(lenBinary, uint32(len(b)))
		p.data = append(p.data, lenBinary...)
		p.data = append(p.data, b...)
		p.addSize(uint16(fv.Len() + 4))

	case reflect.Struct:
		if !isTime(fv) {
			for kk := 0; kk < fv.NumField(); kk++ {
				p.deserialize(fv.Field(kk))
			}
		} else {
			b := make([]byte, 8)
			t := fv.Interface().(time.Time)
			binary.LittleEndian.PutUint64(b, uint64(t.Unix()))
			p.data = append(p.data, b...)
			p.addSize(8)
		}

	case reflect.Ptr:
		if !fv.IsValid() {
			return errors.New("invalid pointer")
		}
		if fv.IsNil() || fv.IsZero() {
			x := reflect.New(fv.Type().Elem())
			if fv.CanSet() {
				fv.Set(x)
			} else {
				fmt.Println("Can not set value for ", fv.Type())
			}
		}
		p.deserialize(fv.Elem())

	case reflect.Slice:
		// write 2 byte for slice length
		b := make([]byte, 2)
		binary.LittleEndian.PutUint16(b, uint16(fv.Len()))
		p.data = append(p.data, b...)
		p.addSize(2)
		// write slice data
		for i := 0; i < fv.Len(); i++ {
			p.deserialize(fv.Index(i))
		}

	default:
		fmt.Println("Encode unsuported reflect type", fv.Kind())
	}
	return nil
}

func isTime(obj reflect.Value) bool {
	_, ok := obj.Interface().(time.Time)
	return ok
}
