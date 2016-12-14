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

	var b []byte
	var err error
	if t.Kind() == reflect.Struct && !isDateTime(t) && !isDateTimeOffset(t) {
		startOffset := (2 + t.NumField()) * intByte4

		// NOTE : memory allocation is not just size.
		b = make([]byte, startOffset, int(t.Type().Size())+startOffset)
		err = d.serializeStruct(t, &b, startOffset)

		hoge, _ := d.calcSize(t)
		fmt.Println(len(b))
		fmt.Println(hoge + uint32(startOffset))
	} else {
		b, err = d.serialize(t)
	}

	return b, err
}

func (d *serializer) serializeStruct(rv reflect.Value, b *[]byte, offset int) error {
	nf := rv.NumField()
	index := 2 * intByte4
	for i := 0; i < nf; i++ {
		// todo : want to receive binary size
		ab, err := d.serialize(rv.Field(i))
		if err != nil {
			return err
		}
		*b = append(*b, ab...)

		(*b)[index], (*b)[index+1], (*b)[index+2], (*b)[index+3] = byte(offset), byte(offset>>8), byte(offset>>16), byte(offset>>24)
		index += intByte4
		offset += len(ab)
	}
	// size
	si := len(*b)
	(*b)[0], (*b)[1], (*b)[2], (*b)[3] = byte(si), byte(si>>8), byte(si>>16), byte(si>>24)
	// index max
	im := nf - 1
	(*b)[4], (*b)[5], (*b)[6], (*b)[7] = byte(im), byte(im>>8), byte(im>>16), byte(im>>24)
	return nil
}

func (d *serializer) serialize(rv reflect.Value) ([]byte, error) {
	var ret []byte

	switch rv.Kind() {
	case reflect.Int8:
		b := make([]byte, uintByte1)
		v := rv.Int()
		b[0] = byte(v)
		ret = b

	case reflect.Int16:
		b := make([]byte, uintByte2)
		v := rv.Int()
		b[0] = byte(v)
		b[1] = byte(v >> 8)
		ret = b

	case reflect.Int32:
		if isChar(rv) {

			// rune [ushort(2)]
			enc := utf16.Encode([]rune{int32(rv.Int())})

			b := make([]byte, uintByte2)
			v := enc[0]
			b[0] = byte(v)
			b[1] = byte(v >> 8)
			ret = b
		} else {
			b := make([]byte, uintByte4)
			v := rv.Int()
			b[0] = byte(v)
			b[1] = byte(v >> 8)
			b[2] = byte(v >> 16)
			b[3] = byte(v >> 24)
			ret = b

		}

	case reflect.Int:
		b := make([]byte, uintByte4)
		v := rv.Int()
		b[0] = byte(v)
		b[1] = byte(v >> 8)
		b[2] = byte(v >> 16)
		b[3] = byte(v >> 24)
		ret = b

	case reflect.Int64:
		if isDuration(rv) {
			b := make([]byte, uintByte4+uintByte8, uintByte4+uintByte8)
			// seconds
			ns := rv.MethodByName("Nanoseconds").Call([]reflect.Value{})[0]
			nanoseconds := ns.Int()
			sec, nsec := nanoseconds/(1000*1000), int32(nanoseconds%(1000*1000))
			b[0] = byte(sec)
			b[1] = byte(sec >> 8)
			b[2] = byte(sec >> 16)
			b[3] = byte(sec >> 24)
			b[4] = byte(sec >> 32)
			b[5] = byte(sec >> 40)
			b[6] = byte(sec >> 48)
			b[7] = byte(sec >> 56)

			// nanos
			o := uintByte8
			b[o+0] = byte(nsec)
			b[o+1] = byte(nsec >> 8)
			b[o+2] = byte(nsec >> 16)
			b[o+3] = byte(nsec >> 24)
			ret = b
		} else {

			b := make([]byte, uintByte8)
			v := rv.Int()
			b[0] = byte(v)
			b[1] = byte(v >> 8)
			b[2] = byte(v >> 16)
			b[3] = byte(v >> 24)
			b[4] = byte(v >> 32)
			b[5] = byte(v >> 40)
			b[6] = byte(v >> 48)
			b[7] = byte(v >> 56)
			ret = b
		}

	case reflect.Uint8:
		b := make([]byte, uintByte1)
		v := rv.Uint()
		b[0] = byte(v)
		ret = b

	case reflect.Uint16:
		b := make([]byte, uintByte2)
		v := rv.Uint()
		b[0] = byte(v)
		b[1] = byte(v >> 8)
		ret = b

	case reflect.Uint32, reflect.Uint:
		b := make([]byte, uintByte4)
		v := rv.Uint()
		b[0] = byte(v)
		b[1] = byte(v >> 8)
		b[2] = byte(v >> 16)
		b[3] = byte(v >> 24)
		ret = b

	case reflect.Uint64:
		b := make([]byte, uintByte8)
		v := rv.Uint()
		b[0] = byte(v)
		b[1] = byte(v >> 8)
		b[2] = byte(v >> 16)
		b[3] = byte(v >> 24)
		b[4] = byte(v >> 32)
		b[5] = byte(v >> 40)
		b[6] = byte(v >> 48)
		b[7] = byte(v >> 56)
		ret = b

	case reflect.Float32:
		b := make([]byte, uintByte4)

		v := math.Float32bits(float32(rv.Float()))
		b[0] = byte(v)
		b[1] = byte(v >> 8)
		b[2] = byte(v >> 16)
		b[3] = byte(v >> 24)
		ret = b

	case reflect.Float64:
		b := make([]byte, uintByte8)

		v := math.Float64bits(rv.Float())
		b[0] = byte(v)
		b[1] = byte(v >> 8)
		b[2] = byte(v >> 16)
		b[3] = byte(v >> 24)
		b[4] = byte(v >> 32)
		b[5] = byte(v >> 40)
		b[6] = byte(v >> 48)
		b[7] = byte(v >> 56)
		ret = b

	case reflect.Bool:
		b := make([]byte, uintByte1)

		if rv.Bool() {
			b[0] = 0x01
		} else {
			b[0] = 0x00
		}
		ret = b

	case reflect.String:
		str := rv.String()
		l := uint32(len(str))
		b := make([]byte, 0, l+uintByte4)
		b = append(b, byte(l), byte(l>>8), byte(l>>16), byte(l>>24))

		// NOTE : unsafe
		strBytes := *(*[]byte)(unsafe.Pointer(&str))
		b = append(b, strBytes...)
		ret = b

	case reflect.Array, reflect.Slice:
		l := rv.Len()
		if l > 0 {
			// first : know element size
			fb, err := d.serialize(rv.Index(0))
			if err != nil {
				return []byte(""), err
			}

			// second : make byte array
			size := uint32(l*len(fb)) + uintByte4
			b := make([]byte, 0, size)

			// third : append data
			b = append(b, byte(l), byte(l>>8), byte(l>>16), byte(l>>24))
			b = append(b, fb...)

			for i := 1; i < l; i++ {
				ab, err := d.serialize(rv.Index(i))
				if err != nil {
					return []byte(""), err
				}
				b = append(b, ab...)
			}
			ret = b
		} else {
			// only make length info
			b := make([]byte, uintByte4)
			ret = b
		}

	case reflect.Struct:
		if isDateTimeOffset(rv) {

			b := make([]byte, uintByte4+uintByte8+uintByte2, uintByte4+uintByte8+uintByte2)

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
			b[0] = byte(seconds)
			b[1] = byte(seconds >> 8)
			b[2] = byte(seconds >> 16)
			b[3] = byte(seconds >> 24)
			b[4] = byte(seconds >> 32)
			b[5] = byte(seconds >> 40)
			b[6] = byte(seconds >> 48)
			b[7] = byte(seconds >> 56)

			// nanos to byte
			o := uintByte8
			b[o+0] = byte(nanos)
			b[o+1] = byte(nanos >> 8)
			b[o+2] = byte(nanos >> 16)
			b[o+3] = byte(nanos >> 24)

			// offset to byte
			o += uintByte4
			b[o+0] = byte(offMin)
			b[o+1] = byte(offMin >> 8)

			ret = b
		} else if isDateTime(rv) {
			b := make([]byte, uintByte4+uintByte8, uintByte4+uintByte8)
			// seconds
			unixTime := rv.MethodByName("Unix").Call([]reflect.Value{})
			sec := unixTime[0].Int()
			b[0] = byte(sec)
			b[1] = byte(sec >> 8)
			b[2] = byte(sec >> 16)
			b[3] = byte(sec >> 24)
			b[4] = byte(sec >> 32)
			b[5] = byte(sec >> 40)
			b[6] = byte(sec >> 48)
			b[7] = byte(sec >> 56)

			// nanos
			nsec := int32(rv.FieldByName("nsec").Int())
			o := uintByte8
			b[o+0] = byte(nsec)
			b[o+1] = byte(nsec >> 8)
			b[o+2] = byte(nsec >> 16)
			b[o+3] = byte(nsec >> 24)
			ret = b
		} else {
			b := make([]byte, 0, rv.Type().Size())
			for i := 0; i < rv.NumField(); i++ {
				ab, err := d.serialize(rv.Field(i))
				if err != nil {
					return []byte(""), err
				}
				b = append(b, ab...)
			}
			ret = b
		}

	case reflect.Ptr:
		if rv.IsNil() {
			return []byte(""), errors.New(fmt.Sprint("pointer is null : ", rv.Type()))
		}
		b, err := d.serialize(rv.Elem())
		if err != nil {
			return []byte(""), err
		}
		ret = b

	default:
		return []byte(""), errors.New(fmt.Sprint("this type is not supported : ", rv.Type()))
	}

	return ret, nil
}

func Serialize2(holder interface{}) ([]byte, error) {
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

		hoge, _ := d.calcSize(t)
		d.create = make([]byte, hoge+uint32(startOffset))

		err = d.serializeStruct2(t, uint32(startOffset))
	} else {
		hoge, _ := d.calcSize(t)
		d.create = make([]byte, hoge)
		_, err = d.serialize2(t, 0)
	}

	return d.create, err
}

func (d *serializer) serializeStruct2(rv reflect.Value, offset uint32) error {
	nf := rv.NumField()
	index := 2 * intByte4
	for i := 0; i < nf; i++ {
		// todo : want to receive binary size
		o, err := d.serialize2(rv.Field(i), offset)
		if err != nil {
			return err
		}

		d.create[index], d.create[index+1], d.create[index+2], d.create[index+3] = byte(offset), byte(offset>>8), byte(offset>>16), byte(offset>>24)
		index += intByte4
		offset = o
	}
	// size
	si := len(d.create)
	d.create[0], d.create[1], d.create[2], d.create[3] = byte(si), byte(si>>8), byte(si>>16), byte(si>>24)
	// index max
	im := nf - 1
	d.create[4], d.create[5], d.create[6], d.create[7] = byte(im), byte(im>>8), byte(im>>16), byte(im>>24)
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

func (d *serializer) serialize2(rv reflect.Value, offset uint32) (uint32, error) {

	switch rv.Kind() {
	case reflect.Int8:
		v := rv.Int()
		d.create[offset] = byte(v)
		offset += uintByte1

	case reflect.Int16:
		v := rv.Int()
		d.create[offset+0] = byte(v)
		d.create[offset+1] = byte(v >> 8)
		offset += uintByte2

	case reflect.Int32:
		if isChar(rv) {

			// rune [ushort(2)]
			enc := utf16.Encode([]rune{int32(rv.Int())})
			v := enc[0]
			d.create[offset+0] = byte(v)
			d.create[offset+1] = byte(v >> 8)
			offset += uintByte2
		} else {
			v := rv.Int()
			d.create[offset+0] = byte(v)
			d.create[offset+1] = byte(v >> 8)
			d.create[offset+2] = byte(v >> 16)
			d.create[offset+3] = byte(v >> 24)
			offset += uintByte4
		}

	case reflect.Int:
		v := rv.Int()
		d.create[offset+0] = byte(v)
		d.create[offset+1] = byte(v >> 8)
		d.create[offset+2] = byte(v >> 16)
		d.create[offset+3] = byte(v >> 24)
		offset += uintByte4

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

			// nanos
			o := uintByte8
			d.create[offset+o+0] = byte(nsec)
			d.create[offset+o+1] = byte(nsec >> 8)
			d.create[offset+o+2] = byte(nsec >> 16)
			d.create[offset+o+3] = byte(nsec >> 24)
			offset += uintByte4 + uintByte8
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
			offset += uintByte8
		}

	case reflect.Uint8:
		v := rv.Uint()
		d.create[offset+0] = byte(v)
		offset += uintByte1

	case reflect.Uint16:
		v := rv.Uint()
		d.create[offset+0] = byte(v)
		d.create[offset+1] = byte(v >> 8)
		offset += uintByte2

	case reflect.Uint32, reflect.Uint:
		v := rv.Uint()
		d.create[offset+0] = byte(v)
		d.create[offset+1] = byte(v >> 8)
		d.create[offset+2] = byte(v >> 16)
		d.create[offset+3] = byte(v >> 24)
		offset += uintByte4

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
		offset += uintByte8

	case reflect.Float32:
		v := math.Float32bits(float32(rv.Float()))
		d.create[offset+0] = byte(v)
		d.create[offset+1] = byte(v >> 8)
		d.create[offset+2] = byte(v >> 16)
		d.create[offset+3] = byte(v >> 24)
		offset += uintByte4

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
		offset += uintByte8

	case reflect.Bool:

		if rv.Bool() {
			d.create[offset+0] = 0x01
		} else {
			d.create[offset+0] = 0x00
		}
		offset += uintByte1

	case reflect.String:
		str := rv.String()
		l := uint32(len(str))
		d.create[offset+0] = byte(l)
		d.create[offset+1] = byte(l >> 8)
		d.create[offset+2] = byte(l >> 16)
		d.create[offset+3] = byte(l >> 24)
		offset += uintByte4

		// NOTE : unsafe
		strBytes := *(*[]byte)(unsafe.Pointer(&str))
		for i := uint32(0); i < l; i++ {
			d.create[offset+i] = strBytes[i]
		}
		offset += l

	case reflect.Array, reflect.Slice:
		l := rv.Len()
		if l > 0 {
			d.create[offset+0] = byte(l)
			d.create[offset+1] = byte(l >> 8)
			d.create[offset+2] = byte(l >> 16)
			d.create[offset+3] = byte(l >> 24)

			for i := 0; i < l; i++ {
				o, err := d.serialize2(rv.Index(i), offset)
				if err != nil {
					return 0, err
				}
				offset += o
			}
		} else {
			// only make length info
			d.create[offset+0] = 0
			d.create[offset+1] = 0
			d.create[offset+2] = 0
			d.create[offset+3] = 0
			offset += uintByte4
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

			// nanos to byte
			offset += uintByte8
			d.create[offset+0] = byte(nanos)
			d.create[offset+1] = byte(nanos >> 8)
			d.create[offset+2] = byte(nanos >> 16)
			d.create[offset+3] = byte(nanos >> 24)

			// offset to byte
			offset += uintByte4
			d.create[offset+0] = byte(offMin)
			d.create[offset+1] = byte(offMin >> 8)

			offset += uintByte2
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

			// nanos
			nsec := int32(rv.FieldByName("nsec").Int())
			offset += uintByte8
			d.create[offset+0] = byte(nsec)
			d.create[offset+1] = byte(nsec >> 8)
			d.create[offset+2] = byte(nsec >> 16)
			d.create[offset+3] = byte(nsec >> 24)

			offset += uintByte4
		} else {
			for i := 0; i < rv.NumField(); i++ {
				o, err := d.serialize2(rv.Field(i), offset)
				if err != nil {
					return 0, err
				}
				offset += o
			}
		}

	case reflect.Ptr:
		if rv.IsNil() {
			return 0, errors.New(fmt.Sprint("pointer is null : ", rv.Type()))
		}
		o, err := d.serialize2(rv.Elem(), offset)
		if err != nil {
			return 0, err
		}
		offset += o

	default:
		return 0, errors.New(fmt.Sprint("this type is not supported : ", rv.Type()))
	}

	return offset, nil
}
