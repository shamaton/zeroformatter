package zeroformatter

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"reflect"
	"time"
	"unicode/utf16"

	"github.com/shamaton/zeroformatter/char"
	"github.com/shamaton/zeroformatter/datetimeoffset"
)

type deserializer struct {
	data []byte
}

func createDeserializer(data []byte) *deserializer {
	return &deserializer{
		data: data,
	}
}

const minStructDataSize = 9

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
	b, offset := d.read_s4(offset)
	size := binary.LittleEndian.Uint32(b)
	if size != uint32(dataLen) {
		return fmt.Errorf("data size is wrong [ %d : %d ]", size, dataLen)
	}

	// index
	b, offset = d.read_s4(offset)
	dataIndex := binary.LittleEndian.Uint32(b)
	if dataIndex != uint32(t.NumField()-1) {
		return fmt.Errorf("data index is diffrent [ %d : %d ]", dataIndex, t.NumField()-1)
	}

	for i := 0; i < t.NumField(); i++ {
		b, offset = d.read_s4(offset)
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

func (d *deserializer) read_s1(index uint32) (byte, uint32) {
	rb := uint32(1)
	return d.data[index], index + rb
}

func (d *deserializer) read_s2(index uint32) ([]byte, uint32) {
	rb := uint32(2)
	return d.data[index : index+rb], index + rb
}

func (d *deserializer) read_s4(index uint32) ([]byte, uint32) {
	rb := uint32(4)
	return d.data[index : index+rb], index + rb
}

func (d *deserializer) read_s8(index uint32) ([]byte, uint32) {
	rb := uint32(8)
	return d.data[index : index+rb], index + rb
}

func (d *deserializer) deserialize(rv reflect.Value, offset uint32) (uint32, error) {
	var err error

	switch rv.Kind() {
	case reflect.Int8:
		b, o := d.read_s1(offset)
		v := int8(b)
		rv.Set(reflect.ValueOf(v))
		// update
		offset = o

	case reflect.Int16:
		// Int16 [short(2)]
		b, o := d.read_s2(offset)
		_v := binary.LittleEndian.Uint16(b)
		v := int16(_v)
		rv.Set(reflect.ValueOf(v))
		// update
		offset = o

	case reflect.Int32:
		// char is used instead of rune
		if isChar(rv) {
			// rune [ushort(2)]
			b, o := d.read_s2(offset)
			u16s := []uint16{binary.LittleEndian.Uint16(b)}
			_v := utf16.Decode(u16s)
			v := char.Char(_v[0])
			rv.Set(reflect.ValueOf(v))

			// update
			offset = o
		} else {
			// Int32 [int(4)]
			b, o := d.read_s4(offset)
			_v := binary.LittleEndian.Uint32(b)
			v := int32(int32(_v))
			rv.Set(reflect.ValueOf(v))
			// update
			offset = o
		}

	case reflect.Int:
		// Int32 [int(4)]
		b, o := d.read_s4(offset)
		_v := binary.LittleEndian.Uint32(b)
		// NOTE : double cast
		v := int(int32(_v))
		rv.Set(reflect.ValueOf(v))
		// update
		offset = o

	case reflect.Int64:
		if isDuration(rv) {
			// todo : NOTE procedure is as same as datetime
			b, o1 := d.read_s8(offset)
			seconds := binary.LittleEndian.Uint64(b)
			b, o2 := d.read_s4(o1)
			nanos := binary.LittleEndian.Uint32(b)
			v := time.Duration(int64(seconds)*1000*1000 + int64(nanos))

			rv.Set(reflect.ValueOf(v))
			// update
			offset = o2
		} else {
			// Int64 [long(8)]
			b, o := d.read_s8(offset)
			_v := binary.LittleEndian.Uint64(b)
			v := int64(_v)
			rv.SetInt(v)
			// update
			offset = o
		}

	case reflect.Uint8:
		// byte in cSharp
		_v, o := d.read_s1(offset)
		v := uint8(_v)
		rv.Set(reflect.ValueOf(v))
		// update
		offset = o

	case reflect.Uint16:
		// Uint16 / Char
		b, o := d.read_s2(offset)
		v := binary.LittleEndian.Uint16(b)
		rv.Set(reflect.ValueOf(v))
		// update
		offset = o

	case reflect.Uint32:
		b, o := d.read_s4(offset)
		v := binary.LittleEndian.Uint32(b)
		rv.Set(reflect.ValueOf(v))
		// update
		offset = o

	case reflect.Uint:
		b, o := d.read_s4(offset)
		_v := binary.LittleEndian.Uint32(b)
		v := uint(_v)
		rv.Set(reflect.ValueOf(v))
		// update
		offset = o

	case reflect.Uint64:
		b, o := d.read_s8(offset)
		v := binary.LittleEndian.Uint64(b)
		rv.SetUint(v)
		// update
		offset = o

	case reflect.Float32:
		// Single
		b, o := d.read_s4(offset)
		_v := binary.LittleEndian.Uint32(b)
		v := math.Float32frombits(_v)
		rv.Set(reflect.ValueOf(v))
		// update
		offset = o

	case reflect.Float64:
		// Double
		b, o := d.read_s8(offset)
		_v := binary.LittleEndian.Uint64(b)
		v := math.Float64frombits(_v)
		rv.Set(reflect.ValueOf(v))
		// update
		offset = o

	case reflect.Bool:
		b, o := d.read_s1(offset)
		if b == 0x01 {
			rv.SetBool(true)
		} else if b == 0x00 {
			rv.SetBool(false)
		}
		// update
		offset = o

	case reflect.String:
		b, o := d.read_s4(offset)
		l := binary.LittleEndian.Uint32(b)
		v := string(d.data[o : o+l])
		rv.SetString(v)
		// update
		offset = o + l

	case reflect.Struct:
		if isDateTimeOffset(rv) {
			b, o1 := d.read_s8(offset)
			seconds := binary.LittleEndian.Uint64(b)
			b, o2 := d.read_s4(o1)
			nanos := binary.LittleEndian.Uint32(b)
			b, o3 := d.read_s2(o2)
			offMin := binary.LittleEndian.Uint16(b)

			v := datetimeoffset.Unix(int64(seconds)-int64(offMin*60), int64(nanos))
			rv.Set(reflect.ValueOf(v))
			// update
			offset = o3

		} else if isDateTime(rv) {
			b, o1 := d.read_s8(offset)
			seconds := binary.LittleEndian.Uint64(b)
			b, o2 := d.read_s4(o1)
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
		// element type
		e := rv.Type().Elem()

		// length
		b, offset := d.read_s4(offset)
		l := int(int32(binary.LittleEndian.Uint32(b)))

		// data is null
		if l < 0 {
			return offset, nil
		}

		o := offset
		tmpSlice := reflect.MakeSlice(rv.Type(), l, l)

		for i := 0; i < l; i++ {
			v := reflect.New(e).Elem()
			o, err = d.deserialize(v, o)
			if err != nil {
				return 0, err
			}

			tmpSlice.Index(i).Set(v)
		}
		rv.Set(tmpSlice)

		// update
		offset = o

	case reflect.Array:
		// element type
		e := rv.Type().Elem()

		// length
		b, offset := d.read_s4(offset)
		l := int(int32(binary.LittleEndian.Uint32(b)))

		// data is null
		if l < 0 {
			return offset, nil
		}
		if l != rv.Len() {
			return 0, fmt.Errorf("Array Length is different : data[%d] array[%d]", l, rv.Len())
		}

		o := offset
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
