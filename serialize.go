package rawpkt

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"reflect"
	"time"
)

// Unmarshal packet back into struct / data.
func (p *Packet) Unmarshal(structPtr interface{}) error {
	vp := reflect.ValueOf(structPtr)
	return p.serialize(vp)
}

func (p *Packet) serialize(fv reflect.Value) error {
	if !fv.IsValid() {
		return errors.New("invalid value")
	}
	rt := fv.Type()
	stopIndex := uint16(rt.Size() + HEADER_SIZE)

	switch fv.Kind() {
	case reflect.Int:
		if p.Size() >= stopIndex && fv.CanSet() {
			fv.SetInt(int64(binary.LittleEndian.Uint32(p.data[HEADER_SIZE:stopIndex])))
			p.removeSize(4)
		}

	case reflect.Uint:
		if p.Size() >= stopIndex && fv.CanSet() {
			fv.SetUint(uint64(binary.LittleEndian.Uint32(p.data[HEADER_SIZE:stopIndex])))
			p.removeSize(4)
		}

	case reflect.Int8:
		if p.Size() >= stopIndex && fv.CanSet() {
			fv.SetInt(int64(p.data[HEADER_SIZE]))
			p.removeSize(1)
		}

	case reflect.Int16:
		if p.Size() >= stopIndex && fv.CanSet() {
			fv.SetInt(int64(binary.LittleEndian.Uint16(p.data[HEADER_SIZE:stopIndex])))
			p.removeSize(2)
		}

	case reflect.Int32:
		if p.Size() >= stopIndex && fv.CanSet() {
			fv.SetInt(int64(binary.LittleEndian.Uint32(p.data[HEADER_SIZE:stopIndex])))
			p.removeSize(4)
		}

	case reflect.Int64:
		if p.Size() >= stopIndex && fv.CanSet() {
			fv.SetInt(int64(binary.LittleEndian.Uint64(p.data[HEADER_SIZE:stopIndex])))
			p.removeSize(8)
		}

	case reflect.Uint8:
		if p.Size() >= stopIndex && fv.CanSet() {
			fv.SetUint(uint64(p.data[HEADER_SIZE]))
			p.removeSize(1)
		}

	case reflect.Uint16:

		if p.Size() >= stopIndex && fv.CanSet() {
			fv.SetUint(uint64(binary.LittleEndian.Uint16(p.data[HEADER_SIZE:stopIndex])))
			p.removeSize(2)
		}

	case reflect.Uint32:
		if p.Size() >= stopIndex && fv.CanSet() {
			fv.SetUint(uint64(binary.LittleEndian.Uint32(p.data[HEADER_SIZE:stopIndex])))
			p.removeSize(4)
		}

	case reflect.Uint64:
		if p.Size() >= stopIndex && fv.CanSet() {
			fv.SetUint(uint64(binary.LittleEndian.Uint64(p.data[HEADER_SIZE:stopIndex])))
			p.removeSize(8)
		}

	case reflect.Bool:
		if p.Size() >= stopIndex && fv.CanSet() {
			bitSet := int8(p.data[HEADER_SIZE])
			if bitSet == 1 {
				fv.SetBool(true)
			} else {
				fv.SetBool(false)
			}
			p.removeSize(1)
		}
	case reflect.Float32:
		if p.Size() >= stopIndex && fv.CanSet() {
			fv.SetFloat(float64(math.Float32frombits(binary.LittleEndian.Uint32(p.data[HEADER_SIZE:stopIndex]))))
			p.removeSize(4)
		}

	case reflect.Float64:
		if p.Size() >= stopIndex && fv.CanSet() {
			fv.SetFloat(math.Float64frombits(binary.LittleEndian.Uint64(p.data[HEADER_SIZE:stopIndex])))
			p.removeSize(8)
		}

	case reflect.Array:
		arrLen := uint16(fv.Len())
		for i := 0; i < int(arrLen); i++ {
			p.serialize(fv.Index(i))
		}

	case reflect.String:
		strLen := binary.LittleEndian.Uint32(p.data[HEADER_SIZE:stopIndex])
		p.removeSize(4)
		if strLen > 0 && strLen+HEADER_SIZE <= uint32(len(p.data)) && fv.CanSet() {
			fv.SetString(string(p.data[HEADER_SIZE : HEADER_SIZE+strLen]))
			p.removeSize(uint16(strLen))
		}

	case reflect.Struct:
		if !isTime(fv) {
			for kk := 0; kk < fv.NumField(); kk++ {
				p.serialize(fv.Field(kk))
			}
		} else {
			ti := time.Unix(int64(binary.LittleEndian.Uint64(p.data[HEADER_SIZE:stopIndex])), 0)
			fv.Set(reflect.ValueOf(ti))
			p.removeSize(8)
		}

	case reflect.Ptr:
		if !fv.IsValid() {
			return errors.New("invalid value pointer")
		}
		if fv.IsNil() || fv.IsZero() {
			x := reflect.New(fv.Type().Elem())
			if fv.CanSet() {
				fv.Set(x)
			}
		}
		p.serialize(fv.Elem())

	case reflect.Slice:
		// read 2 byte for slice length
		sliceLength := binary.LittleEndian.Uint16(p.data[HEADER_SIZE:stopIndex])
		p.removeSize(2)
		fv.Set(reflect.MakeSlice(fv.Type(), int(sliceLength), int(sliceLength)))
		for i := 0; i < int(sliceLength); i++ {
			p.serialize(fv.Index(i))
		}

	default:
		return fmt.Errorf("decode unsuported reflect type %s", fv.Kind())
	}
	return nil
}
