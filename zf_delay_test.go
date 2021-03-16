package zeroformatter_test

import (
	"testing"
	"time"

	"reflect"

	"github.com/shamaton/zeroformatter/v2"
	"github.com/shamaton/zeroformatter/v2/char"
	"github.com/shamaton/zeroformatter/v2/datetimeoffset"
)

func TestDelayDeserialize(t *testing.T) {
	type child3 struct {
		Int int
	}
	type child2 struct {
		Int2Uint        map[int]uint
		Float2Bool      map[float32]bool
		Char2String     map[char.Char]string
		Time2TimeOffset map[time.Time]datetimeoffset.DateTimeOffset
		Duration2Struct map[time.Duration]child3
	}
	type child struct {
		IntArray        []int
		UintArray       []uint
		FloatArray      []float32
		BoolArray       []bool
		CharArray       []char.Char
		StringArray     []string
		TimeArray       []time.Time
		TimeOffsetArray []datetimeoffset.DateTimeOffset
		DurationArray   []time.Duration
		Child           child2
	}
	type st struct {
		Int8       int8
		Int16      int16
		Int32      int32
		Int64      int64
		Uint8      byte
		Uint16     uint16
		Uint32     uint32
		Uint64     uint64
		Float      float32
		Double     float64
		Bool       bool
		Char       char.Char
		String     string
		Time       time.Time
		Duration   time.Duration
		TimeOffset datetimeoffset.DateTimeOffset
		Child      child
	}
	vSt := &st{
		Int32:      -32,
		Int8:       -8,
		Int16:      -16,
		Int64:      -64,
		Uint32:     32,
		Uint8:      8,
		Uint16:     16,
		Uint64:     64,
		Float:      1.23,
		Double:     2.3456,
		Bool:       true,
		Char:       char.Char('A'),
		String:     "Parent",
		Time:       now,
		Duration:   time.Duration(123 * time.Second),
		TimeOffset: datetimeoffset.Now(),

		// child
		Child: child{
			IntArray:        []int{-1, -2, -3, -4, -5},
			UintArray:       []uint{1, 2, 3, 4, 5},
			FloatArray:      []float32{-1.2, -3.4, -5.6, -7.8},
			BoolArray:       []bool{true, true, false, false, true},
			CharArray:       []char.Char{char.Char('X'), char.Char('Y'), char.Char('Z')},
			StringArray:     []string{"str", "ing", "arr", "ay"},
			TimeArray:       []time.Time{now, now, now},
			TimeOffsetArray: []datetimeoffset.DateTimeOffset{datetimeoffset.Now(), datetimeoffset.Now(), datetimeoffset.Now()},
			DurationArray:   []time.Duration{time.Duration(1 * time.Nanosecond), time.Duration(2 * time.Nanosecond)},

			// childchild
			Child: child2{
				Int2Uint:        map[int]uint{-1: 2, -3: 4},
				Float2Bool:      map[float32]bool{-1.1: true, -2.2: false},
				Char2String:     map[char.Char]string{char.Char('A'): "AA", char.Char('B'): "BB"},
				Time2TimeOffset: map[time.Time]datetimeoffset.DateTimeOffset{now: datetimeoffset.Now()},
				Duration2Struct: map[time.Duration]child3{time.Duration(1 * time.Hour): child3{Int: 1}, time.Duration(2 * time.Hour): child3{Int: 2}},
			},
		},
	}

	b, err := zeroformatter.Serialize(vSt)
	if err != nil {
		t.Error(err)
	}

	eHolder := &st{}
	dds, err := zeroformatter.DelayDeserialize(eHolder, b)
	if err != nil {
		t.Error(err)
	}

	// simple
	if err := dds.DeserializeByElement(&eHolder.Int8); err != nil {
		t.Error(err)
	}
	if err := dds.DeserializeByElement(&eHolder.Uint8); err != nil {
		t.Error(err)
	}
	if err := dds.DeserializeByElement(&eHolder.Float); err != nil {
		t.Error(err)
	}
	if err := dds.DeserializeByElement(&eHolder.Char); err != nil {
		t.Error(err)
	}
	if err := dds.DeserializeByElement(&eHolder.Time); err != nil {
		t.Error(err)
	}

	if _bool, err := dds.IsDeserialized(&eHolder.Int16); err != nil || _bool {
		t.Error("deserialized error")
	}

	// multiple
	if err := dds.DeserializeByElement(&eHolder.Int16, &eHolder.Int32, &eHolder.Int64); err != nil {
		t.Error(err)
	}
	if err := dds.DeserializeByElement(&eHolder.Uint16, &eHolder.Uint32, &eHolder.Uint64); err != nil {
		t.Error(err)
	}
	if err := dds.DeserializeByElement(&eHolder.Double, &eHolder.String, &eHolder.Bool); err != nil {
		t.Error(err)
	}
	if err := dds.DeserializeByElement(&eHolder.Duration, &eHolder.TimeOffset, &eHolder.Child); err != nil {
		t.Error(err)
	}

	if _bool, err := dds.IsDeserialized(&eHolder.Int16); err != nil || !_bool {
		t.Error("deserialized error")
	}

	// value equal ?
	if !reflect.DeepEqual(vSt, eHolder) {
		t.Error("value not equal!!")
	}

	iHolder := &st{}

	dds, err = zeroformatter.DelayDeserialize(iHolder, b)
	if err != nil {
		t.Error(err)
	}
	if err := dds.DeserializeByIndex(0); err != nil {
		t.Error(err)
	}
	if err := dds.DeserializeByIndex(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16); err != nil {
		t.Error(err)
	}

	if err := dds.DeserializeByIndex(17); err == nil {
		t.Error("index error")
	}

	// value equal ?
	if !reflect.DeepEqual(vSt, iHolder) {
		t.Error("value not equal!!")
	}
}
