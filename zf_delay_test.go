package zeroformatter

import (
	"testing"
	"time"

	"reflect"

	"github.com/shamaton/zeroformatter/char"
	"github.com/shamaton/zeroformatter/datetimeoffset"
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
		Time:       time.Now(),
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
			TimeArray:       []time.Time{time.Now(), time.Now(), time.Now()},
			TimeOffsetArray: []datetimeoffset.DateTimeOffset{datetimeoffset.Now(), datetimeoffset.Now(), datetimeoffset.Now()},
			DurationArray:   []time.Duration{time.Duration(1 * time.Nanosecond), time.Duration(2 * time.Nanosecond)},

			// childchild
			Child: child2{
				Int2Uint:        map[int]uint{-1: 2, -3: 4},
				Float2Bool:      map[float32]bool{-1.1: true, -2.2: false},
				Char2String:     map[char.Char]string{char.Char('A'): "AA", char.Char('B'): "BB"},
				Time2TimeOffset: map[time.Time]datetimeoffset.DateTimeOffset{time.Now(): datetimeoffset.Now(), time.Now(): datetimeoffset.Now()},
				Duration2Struct: map[time.Duration]child3{time.Duration(1 * time.Hour): child3{Int: 1}, time.Duration(2 * time.Hour): child3{Int: 2}},
			},
		},
	}

	b, err := Serialize(vSt)
	if err != nil {
		t.Error(err)
	}

	holder := &st{}
	dds, err := DelayDeserialize(holder, b)
	if err != nil {
		t.Error(err)
	}

	// simple
	if err := dds.DeserializeByElement(&holder.Int8); err != nil {
		t.Error(err)
	}
	if err := dds.DeserializeByElement(&holder.Uint8); err != nil {
		t.Error(err)
	}
	if err := dds.DeserializeByElement(&holder.Float); err != nil {
		t.Error(err)
	}
	if err := dds.DeserializeByElement(&holder.Char); err != nil {
		t.Error(err)
	}
	if err := dds.DeserializeByElement(&holder.Time); err != nil {
		t.Error(err)
	}

	// multiple
	if err := dds.DeserializeByElement(&holder.Int16, &holder.Int32, &holder.Int64); err != nil {
		t.Error(err)
	}
	if err := dds.DeserializeByElement(&holder.Uint16, &holder.Uint32, &holder.Uint64); err != nil {
		t.Error(err)
	}
	if err := dds.DeserializeByElement(&holder.Double, &holder.String, &holder.Bool); err != nil {
		t.Error(err)
	}
	if err := dds.DeserializeByElement(&holder.Duration, &holder.TimeOffset, &holder.Child); err != nil {
		t.Error(err)
	}

	// value equal ?
	if !reflect.DeepEqual(vSt, holder) {
		t.Error("value not equal!!")
	}
}
