package zeroformatter

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"reflect"
	"time"
	"unicode/utf16"
	"unsafe"

	"github.com/shamaton/zeroformatter/char"
	"github.com/shamaton/zeroformatter/datetimeoffset"
)

type deserializer struct {
	data []byte
}

const minStructDataSize = 9

func createDeserializer(data []byte) *deserializer {
	return &deserializer{
		data: data,
	}
}

// Deserialize analyzes byte data and set into holder.
func Deserialize(holder interface{}, data []byte) error {
	ds := createDeserializer(data)

	t := reflect.ValueOf(holder)
	if t.Kind() != reflect.Ptr {
		return fmt.Errorf("holder must set pointer value. but got: %t", holder)
	}

	t = t.Elem()
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// byte to Struct
	if t.Kind() == reflect.Struct && !isDateTime(t) && !isDateTimeOffset(t) {
		return ds.deserializeStruct(t)
	}

	// byte to primitive
	_, err := ds.deserialize(t, 0)
	return err
}

func (d *deserializer) deserializeStruct(t reflect.Value) error {
	dataLen := len(d.data)
	if dataLen < minStructDataSize {
		return fmt.Errorf("data size is not enough: %d", dataLen)
	}

	// data lookup
	offset := uint32(0)

	// size
	b, offset := d.readSize4(offset)
	size := binary.LittleEndian.Uint32(b)
	if size != uint32(dataLen) {
		return fmt.Errorf("data size is wrong [ %d : %d ]", size, dataLen)
	}

	// index
	b, offset = d.readSize4(offset)
	dataIndex := binary.LittleEndian.Uint32(b)
	numField := t.NumField()
	if dataIndex != uint32(numField-1) {
		return fmt.Errorf("data index is diffrent [ %d : %d ]", dataIndex, numField-1)
	}

	for i := 0; i < numField; i++ {
		b, offset = d.readSize4(offset)
		dataOffset := binary.LittleEndian.Uint32(b)
		if _, err := d.deserialize(t.Field(i), dataOffset); err != nil {
			return err
		}
	}
	return nil
}

func isDateTime(value reflect.Value) bool {
	i := value.Interface()
	switch i.(type) {
	case time.Time:
		return true
	}
	return false
}

func isDateTimeOffset(value reflect.Value) bool {
	i := value.Interface()
	switch i.(type) {
	case datetimeoffset.DateTimeOffset:
		return true
	}
	return false
}

func isDuration(value reflect.Value) bool {
	// check type
	i := value.Interface()
	switch i.(type) {
	case time.Duration:
		return true
	}
	return false
}

func isChar(value reflect.Value) bool {
	i := value.Interface()
	switch i.(type) {
	case char.Char:
		return true
	}
	return false
}

func (d *deserializer) deserialize(rv reflect.Value, offset uint32) (uint32, error) {
	var err error

	switch rv.Kind() {
	case reflect.Int8:
		b, o := d.readSize1(offset)
		rv.SetInt(int64(b))
		// update
		offset = o

	case reflect.Int16:
		// Int16 [short(2)]
		b, o := d.readSize2(offset)
		_v := binary.LittleEndian.Uint16(b)
		rv.SetInt(int64(_v))
		// update
		offset = o

	case reflect.Int32:
		// char is used instead of rune
		if isChar(rv) {
			// rune [ushort(2)]
			b, o := d.readSize2(offset)
			u16s := []uint16{binary.LittleEndian.Uint16(b)}
			_v := utf16.Decode(u16s)
			v := char.Char(_v[0])
			rv.Set(reflect.ValueOf(v))

			// update
			offset = o
		} else {
			// Int32 [int(4)]
			b, o := d.readSize4(offset)
			_v := binary.LittleEndian.Uint32(b)
			// NOTE : double cast
			rv.SetInt(int64(int32(_v)))
			// update
			offset = o
		}

	case reflect.Int:
		// Int32 [int(4)]
		b, o := d.readSize4(offset)
		_v := binary.LittleEndian.Uint32(b)
		// NOTE : double cast
		rv.SetInt(int64(int32(_v)))
		// update
		offset = o

	case reflect.Int64:
		if isDuration(rv) {
			// todo : NOTE procedure is as same as datetime
			b, o1 := d.readSize8(offset)
			seconds := binary.LittleEndian.Uint64(b)
			b, o2 := d.readSize4(o1)
			nanos := binary.LittleEndian.Uint32(b)
			v := time.Duration(int64(seconds)*1000*1000 + int64(nanos))

			rv.Set(reflect.ValueOf(v))
			// update
			offset = o2
		} else {
			// Int64 [long(8)]
			b, o := d.readSize8(offset)
			v := binary.LittleEndian.Uint64(b)
			rv.SetInt(int64(v))
			// update
			offset = o
		}

	case reflect.Uint8:
		// byte in cSharp
		_v, o := d.readSize1(offset)
		rv.SetUint(uint64(_v))
		// update
		offset = o

	case reflect.Uint16:
		// Uint16 / Char
		b, o := d.readSize2(offset)
		v := binary.LittleEndian.Uint16(b)
		rv.SetUint(uint64(v))
		// update
		offset = o

	case reflect.Uint32:
		b, o := d.readSize4(offset)
		v := binary.LittleEndian.Uint32(b)
		rv.SetUint(uint64(v))
		// update
		offset = o

	case reflect.Uint:
		b, o := d.readSize4(offset)
		v := binary.LittleEndian.Uint32(b)
		rv.SetUint(uint64(v))
		// update
		offset = o

	case reflect.Uint64:
		b, o := d.readSize8(offset)
		v := binary.LittleEndian.Uint64(b)
		rv.SetUint(v)
		// update
		offset = o

	case reflect.Float32:
		// Single
		b, o := d.readSize4(offset)
		_v := binary.LittleEndian.Uint32(b)
		v := math.Float32frombits(_v)
		rv.SetFloat(float64(v))
		// update
		offset = o

	case reflect.Float64:
		// Double
		b, o := d.readSize8(offset)
		_v := binary.LittleEndian.Uint64(b)
		v := math.Float64frombits(_v)
		rv.SetFloat(v)
		// update
		offset = o

	case reflect.Bool:
		b, o := d.readSize1(offset)
		if b == 0x01 {
			rv.SetBool(true)
		} else if b == 0x00 {
			rv.SetBool(false)
		}
		// update
		offset = o

	case reflect.String:
		b, o := d.readSize4(offset)
		l := binary.LittleEndian.Uint32(b)
		dd := d.data[o : o+l]
		v := *(*string)(unsafe.Pointer(&dd))
		rv.SetString(v)
		// update
		offset = o + l

	case reflect.Struct:
		if isDateTimeOffset(rv) {
			b, o1 := d.readSize8(offset)
			seconds := binary.LittleEndian.Uint64(b)
			b, o2 := d.readSize4(o1)
			nanos := binary.LittleEndian.Uint32(b)
			b, o3 := d.readSize2(o2)
			offMin := binary.LittleEndian.Uint16(b)

			v := datetimeoffset.Unix(int64(seconds)-int64(offMin*60), int64(nanos))
			rv.Set(reflect.ValueOf(v))
			// update
			offset = o3

		} else if isDateTime(rv) {
			b, o1 := d.readSize8(offset)
			seconds := binary.LittleEndian.Uint64(b)
			b, o2 := d.readSize4(o1)
			nanos := binary.LittleEndian.Uint32(b)
			v := time.Unix(int64(seconds), int64(nanos))

			rv.Set(reflect.ValueOf(v))
			// update
			offset = o2
		} else {
			for i := 0; i < rv.NumField(); i++ {
				offset, err = d.deserialize(rv.Field(i), offset)
				if err != nil {
					return 0, err
				}
			}
		}

	case reflect.Slice:

		// length
		b, o := d.readSize4(offset)
		l := int(int32(binary.LittleEndian.Uint32(b)))

		// data is null
		if l < 0 {
			return o, nil
		}

		tmpSlice := reflect.MakeSlice(rv.Type(), l, l)

		for i := 0; i < l; i++ {
			v := tmpSlice.Index(i)
			o, err = d.deserialize(v, o)
			if err != nil {
				return 0, err
			}
		}
		rv.Set(tmpSlice)

		// update
		offset = o

	case reflect.Array:
		// element type
		e := rv.Type().Elem()

		// length
		b, o := d.readSize4(offset)
		l := int(int32(binary.LittleEndian.Uint32(b)))

		// data is null
		if l < 0 {
			return o, nil
		}
		if l != rv.Len() {
			return 0, fmt.Errorf("Array Length is different : data[%d] array[%d]", l, rv.Len())
		}

		for i := 0; i < l; i++ {
			v := reflect.New(e).Elem()
			o, err = d.deserialize(v, o)
			if err != nil {
				return 0, err
			}
			rv.Index(i).Set(v)
		}

		// update
		offset = o

	case reflect.Map:
		key := rv.Type().Key()
		value := rv.Type().Elem()

		// map length
		b, o := d.readSize4(offset)
		l := int(binary.LittleEndian.Uint32(b))

		if rv.IsNil() {
			rv.Set(reflect.MakeMapWithSize(rv.Type(), l))
		}

		for i := 0; i < l; i++ {
			k := reflect.New(key).Elem()
			v := reflect.New(value).Elem()
			o, err = d.deserialize(k, o)
			if err != nil {
				return 0, err
			}
			o, err = d.deserialize(v, o)
			if err != nil {
				return 0, err
			}

			rv.SetMapIndex(k, v)
			offset = o
		}

	case reflect.Ptr:
		e := rv.Type().Elem()
		v := reflect.New(e).Elem()
		offset, err = d.deserialize(v, offset)
		rv.Set(v.Addr())

	default:
		err = errors.New(fmt.Sprint("this type is not supported : ", rv.Type()))
	}

	return offset, err
}
