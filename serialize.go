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
	byte1 uint32 = 1 << iota
	byte2
	byte4
	byte8
)

type serializer struct {
	create []byte

	queueMapKey   [][]reflect.Value
	queueMapValue []reflect.Value
}

func createSerializer() *serializer {
	return &serializer{
		queueMapKey:   [][]reflect.Value{},
		queueMapValue: []reflect.Value{},
	}
}

// Serialize analyzes holder and converts to byte datas.
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
		startOffset := uint32(2+t.NumField()) * byte4

		dataPartSize, _ := d.calcSize(t)
		size := startOffset + dataPartSize
		d.create = make([]byte, size)

		err = d.serializeStruct(t, startOffset, size)
	} else {
		size, _ := d.calcSize(t)
		d.create = make([]byte, size)
		_, err = d.serialize(t, 0)
	}

	return d.create, err
}

func (d *serializer) serializeStruct(rv reflect.Value, offset uint32, size uint32) error {
	nf := rv.NumField()
	index := 2 * byte4
	for i := 0; i < nf; i++ {
		s, err := d.serialize(rv.Field(i), offset)
		if err != nil {
			return err
		}

		d.create[index], d.create[index+1], d.create[index+2], d.create[index+3] = byte(offset), byte(offset>>8), byte(offset>>16), byte(offset>>24)
		index += byte4
		offset += s
	}
	// size
	d.create[0], d.create[1], d.create[2], d.create[3] = byte(size), byte(size>>8), byte(size>>16), byte(size>>24)
	// last index
	li := nf - 1
	d.create[4], d.create[5], d.create[6], d.create[7] = byte(li), byte(li>>8), byte(li>>16), byte(li>>24)
	return nil
}

func (d *serializer) isFixedSize(rv reflect.Value) bool {
	ret := false
	switch rv.Kind() {
	case
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint,
		reflect.Float32, reflect.Float64,
		reflect.Bool:
		ret = true

	case
		reflect.Struct:
		if isDateTimeOffset(rv) || isDateTime(rv) {
			ret = true
		}

	default:
	}
	return ret
}

func (d *serializer) calcSize(rv reflect.Value) (uint32, error) {
	ret := uint32(0)

	switch rv.Kind() {
	case reflect.Int8:
		ret = byte1

	case reflect.Int16:
		ret = byte2

	case reflect.Int32:
		if isChar(rv) {
			ret = byte2
		} else {
			ret = byte4
		}

	case reflect.Int:
		ret = byte4

	case reflect.Int64:
		if isDuration(rv) {
			ret = byte4 + byte8
		} else {
			ret = byte8
		}

	case reflect.Uint8:
		ret = byte1

	case reflect.Uint16:
		ret = byte2

	case reflect.Uint32, reflect.Uint:
		ret = byte4

	case reflect.Uint64:
		ret = byte8

	case reflect.Float32:
		ret = byte4

	case reflect.Float64:
		ret = byte8

	case reflect.Bool:
		ret = byte1

	case reflect.String:
		l := uint32(rv.Len())
		ret = l + byte4

	case reflect.Array, reflect.Slice:
		l := rv.Len()
		if l > 0 {
			ret += byte4
			isTypeFixed := d.isFixedSize(rv.Index(0))
			if isTypeFixed {
				s, err := d.calcSize(rv.Index(0))
				if err != nil {
					return 0, err
				}
				ret += s * uint32(l)
			} else {
				for i := 0; i < l; i++ {
					s, err := d.calcSize(rv.Index(i))
					if err != nil {
						return 0, err
					}
					ret += s
				}
			}
		} else {
			// only length info
			ret = byte4
		}

	case reflect.Struct:
		if isDateTimeOffset(rv) {
			ret = byte4 + byte8 + byte2
		} else if isDateTime(rv) {
			ret = byte4 + byte8
		} else {
			for i := 0; i < rv.NumField(); i++ {
				s, err := d.calcSize(rv.Field(i))
				if err != nil {
					return 0, err
				}
				ret += s
			}
		}

	case reflect.Map:
		// length
		ret += byte4
		l := uint32(rv.Len())

		if l < 1 {
			return ret, nil
		}
		// check fixed type
		keys := rv.MapKeys()

		d.queueMapKey = append(d.queueMapKey, keys)

		isFixedKey := d.isFixedSize(keys[0])
		if isFixedKey {
			sizeK, err := d.calcSize(keys[0])
			if err != nil {
				return 0, err
			}
			startI := len(d.queueMapValue)
			add := make([]reflect.Value, l)
			d.queueMapValue = append(d.queueMapValue, add...)
			for i, k := range keys {
				value := rv.MapIndex(k)
				d.queueMapValue[startI+i] = value
				sizeV, err := d.calcSize(value)
				if err != nil {
					return 0, err
				}
				ret += sizeK + sizeV
			}

		} else {
			startI := len(d.queueMapValue)
			add := make([]reflect.Value, l)
			d.queueMapValue = append(d.queueMapValue, add...)
			for i, k := range keys {
				sizeK, err := d.calcSize(k)
				if err != nil {
					return 0, err
				}
				value := rv.MapIndex(k)
				d.queueMapValue[startI+i] = value
				sizeV, err := d.calcSize(value)
				if err != nil {
					return 0, err
				}
				ret += sizeK + sizeV
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
		d.writeSize1Int64(rv.Int(), offset)
		size += byte1

	case reflect.Int16:
		d.writeSize2Int64(rv.Int(), offset)
		size += byte2

	case reflect.Int32:
		if isChar(rv) {

			// rune [ushort(2)]
			enc := utf16.Encode([]rune{int32(rv.Int())})
			v := enc[0]
			d.create[offset+0] = byte(v)
			d.create[offset+1] = byte(v >> 8)
			size += byte2
		} else {
			d.writeSize4Int64(rv.Int(), offset)
			size += byte4
		}

	case reflect.Int:
		d.writeSize4Int64(rv.Int(), offset)
		size += byte4

	case reflect.Int64:
		if isDuration(rv) {
			// seconds
			ns := rv.MethodByName("Nanoseconds").Call([]reflect.Value{})[0]
			nanoseconds := ns.Int()
			sec, nsec := nanoseconds/(1000*1000), int64(nanoseconds%(1000*1000))
			d.writeSize8Int64(sec, offset)
			size += byte8
			offset += byte8

			// nanos
			d.writeSize4Int64(nsec, offset)
			size += byte4
		} else {
			d.writeSize8Int64(rv.Int(), offset)
			size += byte8
		}

	case reflect.Uint8:
		d.writeSize1Uint64(rv.Uint(), offset)
		size += byte1

	case reflect.Uint16:
		d.writeSize2Uint64(rv.Uint(), offset)
		size += byte2

	case reflect.Uint32, reflect.Uint:
		d.writeSize4Uint64(rv.Uint(), offset)
		size += byte4

	case reflect.Uint64:
		d.writeSize8Uint64(rv.Uint(), offset)
		size += byte8

	case reflect.Float32:
		v := math.Float32bits(float32(rv.Float()))
		d.writeSize4Uint32(v, offset)
		size += byte4

	case reflect.Float64:
		v := math.Float64bits(rv.Float())
		d.writeSize8Uint64(v, offset)
		size += byte8

	case reflect.Bool:

		if rv.Bool() {
			d.writeSize1Uint64(0x01, offset)
		} else {
			d.writeSize1Uint64(0x00, offset)
		}
		size += byte1

	case reflect.String:
		str := rv.String()
		l := uint32(len(str))
		d.writeSize4Uint32(l, offset)
		size += byte4
		offset += byte4

		// NOTE : unsafe
		strBytes := *(*[]byte)(unsafe.Pointer(&str))
		for i := uint32(0); i < l; i++ {
			d.create[offset+i] = strBytes[i]
		}
		size += l

	case reflect.Array, reflect.Slice:
		l := rv.Len()
		if l > 0 {
			d.writeSize4Int(l, offset)
			size += byte4
			offset += byte4

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
			d.writeSize4Int(0, offset)
			size += byte4
		}

	case reflect.Struct:
		if isDateTimeOffset(rv) {

			// offset
			rets := rv.MethodByName("Zone").Call([]reflect.Value{})
			_, offSec := rets[0] /*name*/, rets[1].Int() /*offset*/
			offMin := offSec / 60

			// seconds
			rets = rv.MethodByName("Unix").Call([]reflect.Value{})
			seconds := rets[0].Int() + offSec

			// nanos
			rets = rv.MethodByName("Nanosecond").Call([]reflect.Value{})
			nanos := rets[0].Int()

			// seconds to byte
			d.writeSize8Int64(seconds, offset)
			size += byte8
			offset += byte8

			// nanos to byte
			d.writeSize4Int64(nanos, offset)
			size += byte4
			offset += byte4

			// offset to byte
			d.writeSize2Int64(offMin, offset)
			size += byte2
		} else if isDateTime(rv) {
			// seconds
			unixTime := rv.MethodByName("Unix").Call([]reflect.Value{})
			sec := unixTime[0].Int()
			d.writeSize8Int64(sec, offset)
			size += byte8
			offset += byte8

			// nanos
			rets := rv.MethodByName("Nanosecond").Call([]reflect.Value{})
			nsec := rets[0].Int()
			d.writeSize4Int64(nsec, offset)
			size += byte4
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

	case reflect.Map:
		// length
		l := rv.Len()
		d.writeSize4Int(l, offset)
		size += byte4
		offset += byte4

		if l < 1 {
			return size, nil
		}

		keys := d.queueMapKey[0]
		keysLen := len(keys)
		values := d.queueMapValue[:keysLen]

		for i, k := range keys {
			addOffByK, err := d.serialize(k, offset)
			if err != nil {
				return 0, err
			}
			addOffByV, err := d.serialize(values[i], offset+addOffByK)
			if err != nil {
				return 0, err
			}
			offset += addOffByK + addOffByV
			size += addOffByK + addOffByV
		}

		// update queue
		if len(d.queueMapKey) > 0 {
			d.queueMapKey = d.queueMapKey[1:]
			d.queueMapValue = d.queueMapValue[keysLen:]
		} else {
			d.queueMapKey = d.queueMapKey[:0]
			d.queueMapValue = d.queueMapValue[:0]
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
