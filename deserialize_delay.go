package zeroformatter

import (
	"encoding/binary"
	"fmt"
	"reflect"
)

type delayDeserializer struct {
	*deserializer
	holder       reflect.Value
	processedMap map[uintptr]int
	indexArray   []uintptr
}

func createDelayDeserialize(deserializer *deserializer, holder reflect.Value, num int) *delayDeserializer {
	return &delayDeserializer{
		deserializer: deserializer,
		holder:       holder,
		processedMap: map[uintptr]int{},
		indexArray:   make([]uintptr, num),
	}
}

func DelayDeserialize(holder interface{}, data []byte) (*delayDeserializer, error) {

	t := reflect.ValueOf(holder)
	if t.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("holder must set pointer value. but got: %t", holder)
	}

	t = t.Elem()
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// delaying enable is struct only
	if t.Kind() != reflect.Struct || isDateTime(t) || isDateTimeOffset(t) {
		return nil, fmt.Errorf("only defined struct can delay deserialize: %t", holder)
	}

	// check before detail checking
	dataLen := len(data)
	if dataLen < minStructDataSize {
		return nil, fmt.Errorf("data size is not enough: %d", dataLen)
	}

	// create deserializer
	ds := createDeserializer(data)

	// check size
	offset := uint32(0)
	b, offset := ds.read_s4(offset)

	size := binary.LittleEndian.Uint32(b)
	if size != uint32(dataLen) {
		return nil, fmt.Errorf("data size is wrong [ %d : %d ]", size, dataLen)
	}

	// check index
	b, offset = ds.read_s4(offset)
	dataIndex := binary.LittleEndian.Uint32(b)
	numField := t.NumField()
	if dataIndex != uint32(numField-1) {
		return nil, fmt.Errorf("data index is diffrent [ %d : %d ]", dataIndex, numField-1)
	}

	// create delay deserializer
	dds := createDelayDeserialize(ds, t, numField)

	// make access info
	for i := 0; i < numField; i++ {
		e := t.Field(i)
		p := e.Addr().Pointer()
		dds.processedMap[p] = i
		dds.indexArray[i] = p
	}

	return dds, nil
}

func (d *delayDeserializer) DeserializeByIndex(i int, indexes ...int) error {
	// index
	if err := d.deserializeByIndex(i); err != nil {
		return err
	}

	// indexes
	for _, idx := range indexes {
		if err := d.deserializeByIndex(idx); err != nil {
			return err
		}
	}
	return nil
}

func (d *delayDeserializer) deserializeByIndex(i int) error {
	if i >= len(d.indexArray) {
		return fmt.Errorf("this index is out of range : %d", i)
	}

	addr := d.indexArray[i]
	return d.deserializeByAddress(addr)
}

func (d *delayDeserializer) DeserializeByElement(element interface{}, elements ...interface{}) error {
	// element
	if err := d.deserializeByElement(element); err != nil {
		return err
	}

	// elements
	for _, e := range elements {
		if err := d.deserializeByElement(e); err != nil {
			return err
		}
	}
	return nil
}

func (d *delayDeserializer) deserializeByElement(element interface{}) error {

	t := reflect.ValueOf(element)
	if t.Kind() != reflect.Ptr {
		return fmt.Errorf("element must set pointer value. but got: %t", element)
	}

	t = t.Elem()
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// address
	address := t.Addr().Pointer()
	return d.deserializeByAddress(address)
}

func (d *delayDeserializer) deserializeByAddress(address uintptr) error {
	index, ok := d.processedMap[address]
	if !ok {
		return fmt.Errorf("not found address: %t", address)
	}

	// already deserialized
	if index < 0 {
		return nil
	}

	// value
	rv := d.holder.Field(index)
	// offset
	off := 8 + uint32(index)*byte4
	b, _ := d.read_s4(off)
	dataIndex := binary.LittleEndian.Uint32(b)

	// deserialize and update flag
	d.deserialize(rv, dataIndex)
	d.processedMap[address] = -1
	return nil
}

func (d *delayDeserializer) IsDeserialized(element interface{}) (bool, error) {

	t := reflect.ValueOf(element)
	if t.Kind() != reflect.Ptr {
		return false, fmt.Errorf("holder must set pointer value. but got: %t", element)
	}

	t = t.Elem()
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// address
	address := t.Addr().Pointer()

	index, ok := d.processedMap[address]
	if !ok {
		return false, fmt.Errorf("not found element: %t", element)
	}

	return (index < 0), nil
}
