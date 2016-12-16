package zeroformatter

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"unicode/utf16"
	"unsafe"
)

const (
	uintByte1 uint32 = 1 << iota
	uintByte2
	uintByte4
	uintByte8
)

const (
	intByte1 = 1 << iota
	intByte2
	intByte4
	intByte8
)

type serializer struct {
	data   []byte
	create []byte
}

func createSerializer() *serializer {
	return &serializer{}
}

func Serialize(holder interface{}) ([]byte, error) {
	d := createSerializer()

	t := reflect.ValueOf(holder)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
	}

	var err error
	if t.Kind() == reflect.Struct && !isDateTime(t) && !isDateTimeOffset(t) {
		startOffset := (2 + t.NumField()) * intByte4

		dataPartSize, _ := d.calcSize(t)
		size := uint32(startOffset) + dataPartSize
		d.create = make([]byte, size)

		err = d.serializeStruct(t, uint32(startOffset), size)
	} else {
		size, _ := d.calcSize(t)
		d.create = make([]byte, size)
		_, err = d.serialize(t, 0)
	}

	return d.create, err
}

func (d *serializer) serializeStruct(rv reflect.Value, offset uint32, size uint32) error {
	nf := rv.NumField()
	index := 2 * intByte4
	for i := 0; i < nf; i++ {
		s, err := d.serialize(rv.Field(i), offset)
		if err != nil {
			return err
		}

		d.create[index], d.create[index+1], d.create[index+2], d.create[index+3] = byte(offset), byte(offset>>8), byte(offset>>16), byte(offset>>24)
		index += intByte4
		offset += s
	}
	// size
	d.create[0], d.create[1], d.create[2], d.create[3] = byte(size), byte(size>>8), byte(size>>16), byte(size>>24)
	// last index
	li := nf - 1
	d.create[4], d.create[5], d.create[6], d.create[7] = byte(li), byte(li>>8), byte(li>>16), byte(li>>24)
	return nil
}

func (d *serializer) calcSize(rv reflect.Value) (uint32, error) {
	ret := uint32(0)

	switch rv.Kind() {
	case reflect.Int8:
		ret = uintByte1

	case reflect.Int16:
		ret = uintByte2

	case reflect.Int32:
		if isChar(rv) {
			ret = uintByte2
		} else {
			ret = uintByte4
		}

	case reflect.Int:
		ret = uintByte4

	case reflect.Int64:
		if isDuration(rv) {
			ret = uintByte4 + uintByte8
		} else {
			ret = uintByte8
		}

	case reflect.Uint8:
		ret = uintByte1

	case reflect.Uint16:
		ret = uintByte2

	case reflect.Uint32, reflect.Uint:
		ret = uintByte4

	case reflect.Uint64:
		ret = uintByte8

	case reflect.Float32:
		ret = uintByte4

	case reflect.Float64:
		ret = uintByte8

	case reflect.Bool:
		ret = uintByte1

	case reflect.String:
		//str := rv.String()
		//l := uint32(len(str))
		l := uint32(rv.Len())
		ret = l + uintByte4

	case reflect.Array, reflect.Slice:
		l := rv.Len()
		if l > 0 {
			ret += uintByte4
			// todo : fixed or variable
			for i := 0; i < l; i++ {
				s, err := d.calcSize(rv.Index(i))
				if err != nil {
					return 0, err
				}
				ret += s
			}
		} else {
			// only length info
			ret = uintByte4
		}

	case reflect.Struct:
		if isDateTimeOffset(rv) {
			ret = uintByte4 + uintByte8 + uintByte2
		} else if isDateTime(rv) {
			ret = uintByte4 + uintByte8
		} else {
			for i := 0; i < rv.NumField(); i++ {
				s, err := d.calcSize(rv.Field(i))
				if err != nil {
					return 0, err
				}
				ret += s
			}
		}

	case reflect.Ptr:
		if rv.IsNil() {
			return 0, errors.New(fmt.Sprint("pointer is null : ", rv.Type()))
		}
		s, err := d.calcSize(rv.Elem())
		if err != nil {
			return 0, err
		}
		ret = s

	default:
		return 0, errors.New(fmt.Sprint("this type is not supported : ", rv.Type()))
	}

	return ret, nil
}

func (d *serializer) serialize(rv reflect.Value, offset uint32) (uint32, error) {
	size := uint32(0)

	switch rv.Kind() {
	case reflect.Int8:
		v := rv.Int()
		d.create[offset] = byte(v)
		size += uintByte1

	case reflect.Int16:
		v := rv.Int()
		d.create[offset+0] = byte(v)
		d.create[offset+1] = byte(v >> 8)
		size += uintByte2

	case reflect.Int32:
		if isChar(rv) {

			// rune [ushort(2)]
			enc := utf16.Encode([]rune{int32(rv.Int())})
			v := enc[0]
			d.create[offset+0] = byte(v)
			d.create[offset+1] = byte(v >> 8)
			size += uintByte2
		} else {
			v := rv.Int()
			d.create[offset+0] = byte(v)
			d.create[offset+1] = byte(v >> 8)
			d.create[offset+2] = byte(v >> 16)
			d.create[offset+3] = byte(v >> 24)
			size += uintByte4
		}

	case reflect.Int:
		v := rv.Int()
		d.create[offset+0] = byte(v)
		d.create[offset+1] = byte(v >> 8)
		d.create[offset+2] = byte(v >> 16)
		d.create[offset+3] = byte(v >> 24)
		size += uintByte4

	case reflect.Int64:
		if isDuration(rv) {
			// seconds
			ns := rv.MethodByName("Nanoseconds").Call([]reflect.Value{})[0]
			nanoseconds := ns.Int()
			sec, nsec := nanoseconds/(1000*1000), int32(nanoseconds%(1000*1000))
			d.create[offset+0] = byte(sec)
			d.create[offset+1] = byte(sec >> 8)
			d.create[offset+2] = byte(sec >> 16)
			d.create[offset+3] = byte(sec >> 24)
			d.create[offset+4] = byte(sec >> 32)
			d.create[offset+5] = byte(sec >> 40)
			d.create[offset+6] = byte(sec >> 48)
			d.create[offset+7] = byte(sec >> 56)
			size += uintByte8
			offset += uintByte8

			// nanos
			d.create[offset+0] = byte(nsec)
			d.create[offset+1] = byte(nsec >> 8)
			d.create[offset+2] = byte(nsec >> 16)
			d.create[offset+3] = byte(nsec >> 24)
			size += uintByte4
		} else {

			v := rv.Int()
			d.create[offset+0] = byte(v)
			d.create[offset+1] = byte(v >> 8)
			d.create[offset+2] = byte(v >> 16)
			d.create[offset+3] = byte(v >> 24)
			d.create[offset+4] = byte(v >> 32)
			d.create[offset+5] = byte(v >> 40)
			d.create[offset+6] = byte(v >> 48)
			d.create[offset+7] = byte(v >> 56)
			size += uintByte8
		}

	case reflect.Uint8:
		v := rv.Uint()
		d.create[offset+0] = byte(v)
		size += uintByte1

	case reflect.Uint16:
		v := rv.Uint()
		d.create[offset+0] = byte(v)
		d.create[offset+1] = byte(v >> 8)
		size += uintByte2

	case reflect.Uint32, reflect.Uint:
		v := rv.Uint()
		d.create[offset+0] = byte(v)
		d.create[offset+1] = byte(v >> 8)
		d.create[offset+2] = byte(v >> 16)
		d.create[offset+3] = byte(v >> 24)
		size += uintByte4

	case reflect.Uint64:
		v := rv.Uint()
		d.create[offset+0] = byte(v)
		d.create[offset+1] = byte(v >> 8)
		d.create[offset+2] = byte(v >> 16)
		d.create[offset+3] = byte(v >> 24)
		d.create[offset+4] = byte(v >> 32)
		d.create[offset+5] = byte(v >> 40)
		d.create[offset+6] = byte(v >> 48)
		d.create[offset+7] = byte(v >> 56)
		size += uintByte8

	case reflect.Float32:
		v := math.Float32bits(float32(rv.Float()))
		d.create[offset+0] = byte(v)
		d.create[offset+1] = byte(v >> 8)
		d.create[offset+2] = byte(v >> 16)
		d.create[offset+3] = byte(v >> 24)
		size += uintByte4

	case reflect.Float64:
		v := math.Float64bits(rv.Float())
		d.create[offset+0] = byte(v)
		d.create[offset+1] = byte(v >> 8)
		d.create[offset+2] = byte(v >> 16)
		d.create[offset+3] = byte(v >> 24)
		d.create[offset+4] = byte(v >> 32)
		d.create[offset+5] = byte(v >> 40)
		d.create[offset+6] = byte(v >> 48)
		d.create[offset+7] = byte(v >> 56)
		size += uintByte8

	case reflect.Bool:

		if rv.Bool() {
			d.create[offset+0] = 0x01
		} else {
			d.create[offset+0] = 0x00
		}
		size += uintByte1

	case reflect.String:
		str := rv.String()
		l := uint32(len(str))
		d.create[offset+0] = byte(l)
		d.create[offset+1] = byte(l >> 8)
		d.create[offset+2] = byte(l >> 16)
		d.create[offset+3] = byte(l >> 24)
		size += uintByte4
		offset += uintByte4

		// NOTE : unsafe
		strBytes := *(*[]byte)(unsafe.Pointer(&str))
		for i := uint32(0); i < l; i++ {
			d.create[offset+i] = strBytes[i]
		}
		size += l

	case reflect.Array, reflect.Slice:
		l := rv.Len()
		if l > 0 {
			d.create[offset+0] = byte(l)
			d.create[offset+1] = byte(l >> 8)
			d.create[offset+2] = byte(l >> 16)
			d.create[offset+3] = byte(l >> 24)
			size += uintByte4
			offset += uintByte4

			for i := 0; i < l; i++ {
				s, err := d.serialize(rv.Index(i), offset)
				if err != nil {
					return 0, err
				}
				offset += s
				size += s
			}
		} else {
			// only make length info
			d.create[offset+0] = 0
			d.create[offset+1] = 0
			d.create[offset+2] = 0
			d.create[offset+3] = 0
			size += uintByte4
		}

	case reflect.Struct:
		if isDateTimeOffset(rv) {

			// offset
			rets := rv.MethodByName("Zone").Call([]reflect.Value{})
			_, offSec := rets[0] /*name*/, rets[1].Int() /*offset*/
			offMin := uint16(offSec / 60)

			// seconds
			rets = rv.MethodByName("Unix").Call([]reflect.Value{})
			seconds := rets[0].Int() + offSec

			// nanos
			nanos := int32(rv.FieldByName("nsec").Int())

			// seconds to byte
			d.create[offset+0] = byte(seconds)
			d.create[offset+1] = byte(seconds >> 8)
			d.create[offset+2] = byte(seconds >> 16)
			d.create[offset+3] = byte(seconds >> 24)
			d.create[offset+4] = byte(seconds >> 32)
			d.create[offset+5] = byte(seconds >> 40)
			d.create[offset+6] = byte(seconds >> 48)
			d.create[offset+7] = byte(seconds >> 56)
			size += uintByte8
			offset += uintByte8

			// nanos to byte
			d.create[offset+0] = byte(nanos)
			d.create[offset+1] = byte(nanos >> 8)
			d.create[offset+2] = byte(nanos >> 16)
			d.create[offset+3] = byte(nanos >> 24)
			size += uintByte4
			offset += uintByte4

			// offset to byte
			d.create[offset+0] = byte(offMin)
			d.create[offset+1] = byte(offMin >> 8)
			size += uintByte2
		} else if isDateTime(rv) {
			// seconds
			unixTime := rv.MethodByName("Unix").Call([]reflect.Value{})
			sec := unixTime[0].Int()
			d.create[offset+0] = byte(sec)
			d.create[offset+1] = byte(sec >> 8)
			d.create[offset+2] = byte(sec >> 16)
			d.create[offset+3] = byte(sec >> 24)
			d.create[offset+4] = byte(sec >> 32)
			d.create[offset+5] = byte(sec >> 40)
			d.create[offset+6] = byte(sec >> 48)
			d.create[offset+7] = byte(sec >> 56)
			size += uintByte8
			offset += uintByte8

			// nanos
			nsec := int32(rv.FieldByName("nsec").Int())
			d.create[offset+0] = byte(nsec)
			d.create[offset+1] = byte(nsec >> 8)
			d.create[offset+2] = byte(nsec >> 16)
			d.create[offset+3] = byte(nsec >> 24)
			size += uintByte4
		} else {
			for i := 0; i < rv.NumField(); i++ {
				s, err := d.serialize(rv.Field(i), offset)
				if err != nil {
					return 0, err
				}
				offset += s
				size += s
			}
		}

	case reflect.Ptr:
		if rv.IsNil() {
			return 0, errors.New(fmt.Sprint("pointer is null : ", rv.Type()))
		}
		s, err := d.serialize(rv.Elem(), offset)
		if err != nil {
			return 0, err
		}
		size += s

	default:
		return 0, errors.New(fmt.Sprint("this type is not supported : ", rv.Type()))
	}

	return size, nil
}
