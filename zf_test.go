package zeroformatter

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/shamaton/zeroformatter/char"
	"github.com/shamaton/zeroformatter/datetimeoffset"
)

func TestPrimitiveInt(t *testing.T) {
	var rInt8 int8
	vInt8 := int8(-8)
	if err := checkRoutine(t, vInt8, &rInt8, false); err != nil {
		t.Error(err)
	}

	var rInt16 int16
	vInt16 := int16(-16)
	if err := checkRoutine(t, vInt16, &rInt16, false); err != nil {
		t.Error(err)
	}

	var rInt int
	vInt := -65535
	if err := checkRoutine(t, vInt, &rInt, false); err != nil {
		t.Error(err)
	}

	var rInt32 int32
	vInt32 := int32(-32)
	if err := checkRoutine(t, vInt32, &rInt32, false); err != nil {
		t.Error(err)
	}

	var rInt64 int64
	vInt64 := int64(-64)
	if err := checkRoutine(t, vInt64, &rInt64, false); err != nil {
		t.Error(err)
	}

	// pointer
	pvInt8 := new(int8)
	prInt8 := new(int8)
	*pvInt8 = math.MaxInt8
	if err := checkRoutine(t, pvInt8, prInt8, false); err != nil {
		t.Error(err)
	}
	if err := checkRoutine(t, &pvInt8, &prInt8, false); err != nil {
		t.Error(err)
	}

	pvInt16 := new(int16)
	prInt16 := new(int16)
	*pvInt16 = math.MaxInt16
	if err := checkRoutine(t, pvInt16, prInt16, false); err != nil {
		t.Error(err)
	}
	if err := checkRoutine(t, &pvInt16, &prInt16, false); err != nil {
		t.Error(err)
	}

	pvInt := new(int)
	prInt := new(int)
	*pvInt = math.MaxInt32
	if err := checkRoutine(t, pvInt, prInt, false); err != nil {
		t.Error(err)
	}
	if err := checkRoutine(t, &pvInt, &prInt, false); err != nil {
		t.Error(err)
	}

	pvInt32 := new(int32)
	prInt32 := new(int32)
	*pvInt32 = math.MaxInt32
	if err := checkRoutine(t, pvInt32, prInt32, false); err != nil {
		t.Error(err)
	}
	if err := checkRoutine(t, &pvInt32, &prInt32, false); err != nil {
		t.Error(err)
	}

	pvInt64 := new(int64)
	prInt64 := new(int64)
	*pvInt64 = math.MaxInt64
	if err := checkRoutine(t, pvInt64, prInt64, false); err != nil {
		t.Error(err)
	}
	if err := checkRoutine(t, &pvInt64, &prInt64, false); err != nil {
		t.Error(err)
	}

	// error
	var rError int32
	vError := int(-1)
	if err := checkRoutine(t, vError, &rError, false); err != nil {
		if strings.Contains(err.Error(), "value diffrent") {
			t.Error(err)
		}
	}
}

func TestPrimitiveUint(t *testing.T) {
	var rUint8 uint8
	vUint8 := uint8(math.MaxUint8)
	if err := checkRoutine(t, vUint8, &rUint8, false); err != nil {
		t.Error(err)
	}

	var rUint16 uint16
	vUint16 := uint16(math.MaxUint16)
	if err := checkRoutine(t, vUint16, &rUint16, false); err != nil {
		t.Error(err)
	}

	var rUint uint
	vUint := uint(math.MaxUint32 / 2)
	if err := checkRoutine(t, vUint, &rUint, false); err != nil {
		t.Error(err)
	}

	var rUint32 uint32
	vUint32 := uint32(math.MaxUint32)
	if err := checkRoutine(t, vUint32, &rUint32, false); err != nil {
		t.Error(err)
	}

	var rUint64 uint64
	vUint64 := uint64(math.MaxUint64)
	if err := checkRoutine(t, vUint64, &rUint64, false); err != nil {
		t.Error(err)
	}

	// pointer
	pvUint8 := new(uint8)
	prUint8 := new(uint8)
	*pvUint8 = math.MaxUint8
	if err := checkRoutine(t, pvUint8, prUint8, false); err != nil {
		t.Error(err)
	}
	if err := checkRoutine(t, &pvUint8, &prUint8, false); err != nil {
		t.Error(err)
	}

	pvUint16 := new(uint16)
	prUint16 := new(uint16)
	*pvUint16 = math.MaxUint16
	if err := checkRoutine(t, pvUint16, prUint16, false); err != nil {
		t.Error(err)
	}
	if err := checkRoutine(t, &pvUint16, &prUint16, false); err != nil {
		t.Error(err)
	}

	pvUint := new(uint)
	prUint := new(uint)
	*pvUint = math.MaxUint32
	if err := checkRoutine(t, pvUint, prUint, false); err != nil {
		t.Error(err)
	}
	if err := checkRoutine(t, &pvUint, &prUint, false); err != nil {
		t.Error(err)
	}

	pvUint32 := new(uint32)
	prUint32 := new(uint32)
	*pvUint32 = math.MaxUint32
	if err := checkRoutine(t, pvUint32, prUint32, false); err != nil {
		t.Error(err)
	}
	if err := checkRoutine(t, &pvUint32, &prUint32, false); err != nil {
		t.Error(err)
	}

	pvUint64 := new(uint64)
	prUint64 := new(uint64)
	*pvUint64 = math.MaxUint64
	if err := checkRoutine(t, pvUint64, prUint64, false); err != nil {
		t.Error(err)
	}
	if err := checkRoutine(t, &pvUint64, &prUint64, false); err != nil {
		t.Error(err)
	}

	var rError uint32
	vError := uint(1)
	if err := checkRoutine(t, vError, &rError, false); err != nil {
		if strings.Contains(err.Error(), "value diffrent") {
			t.Error(err)
		}
	}
}

func TestPrimitiveFloat(t *testing.T) {

	var rFloat32 float32
	vFloat32 := float32(math.MaxFloat32)
	if err := checkRoutine(t, vFloat32, &rFloat32, false); err != nil {
		t.Error(err)
	}

	var rFloat64 float64
	vFloat64 := math.MaxFloat64
	if err := checkRoutine(t, vFloat64, &rFloat64, false); err != nil {
		t.Error(err)
	}

	// pointer
	pvFloat32 := new(float32)
	prFloat32 := new(float32)
	*pvFloat32 = math.MaxFloat32 / 2
	if err := checkRoutine(t, pvFloat32, prFloat32, false); err != nil {
		t.Error(err)
	}
	if err := checkRoutine(t, &pvFloat32, &prFloat32, false); err != nil {
		t.Error(err)
	}

	pvFloat64 := new(float64)
	prFloat64 := new(float64)
	*pvFloat64 = math.MaxFloat64 / 2
	if err := checkRoutine(t, pvFloat64, prFloat64, false); err != nil {
		t.Error(err)
	}
	if err := checkRoutine(t, &pvFloat64, &prFloat64, false); err != nil {
		t.Error(err)
	}

	// error
	var rError float32
	vError := float64(1)
	if err := checkRoutine(t, vError, &rError, false); err != nil {
		if strings.Contains(err.Error(), "value diffrent") {
			t.Error(err)
		}
	}
}

func TestPrimitiveBool(t *testing.T) {
	var rBool bool
	vBool := true
	if err := checkRoutine(t, vBool, &rBool, false); err != nil {
		t.Error(err)
	}

	pvBool := new(bool)
	prBool := new(bool)
	*pvBool = true
	if err := checkRoutine(t, pvBool, prBool, false); err != nil {
		t.Error(err)
	}
	if err := checkRoutine(t, &pvBool, &prBool, false); err != nil {
		t.Error(err)
	}
}

func TestPrimitiveString(t *testing.T) {

	var rChar char.Char
	vChar := char.Char('Z')
	if err := checkRoutine(t, vChar, &rChar, false); err != nil {
		t.Error(err)
	}
	//t.Logf("%#U", rChar)

	var rString string
	vString := "this string serialize and deserialize."
	if err := checkRoutine(t, vString, &rString, false); err != nil {
		t.Error(err)
	}

	// pointer
	pvChar := new(char.Char)
	prChar := new(char.Char)
	*pvChar = char.Char('Y')
	if err := checkRoutine(t, pvChar, prChar, false); err != nil {
		t.Error(err)
	}
	if err := checkRoutine(t, &pvChar, &prChar, false); err != nil {
		t.Error(err)
	}
	//t.Logf("%#U", *prChar)

	pvString := new(string)
	prString := new(string)
	*pvString = "this string is pointer value"
	if err := checkRoutine(t, pvString, prString, false); err != nil {
		t.Error(err)
	}
	if err := checkRoutine(t, &pvString, &prString, false); err != nil {
		t.Error(err)
	}
}

func TestPrimitiveTime(t *testing.T) {

	var rTime time.Time
	vTime := time.Now()
	if err := checkRoutine(t, vTime, &rTime, false); err != nil {
		t.Error(err)
	}

	var rDuration time.Duration
	vDuration := time.Duration(12*time.Hour + 34*time.Minute + 56*time.Second + 78*time.Nanosecond)
	if err := checkRoutine(t, vDuration, &rDuration, false); err != nil {
		t.Error(err)
	}

	var rOffset datetimeoffset.DateTimeOffset
	vOffset := datetimeoffset.Now()
	if err := checkRoutine(t, vOffset, &rOffset, false); err != nil {
		t.Error(err)
	}

	// pointer
	pvTime := new(time.Time)
	prTime := new(time.Time)
	*pvTime = time.Now()
	if err := checkRoutine(t, pvTime, prTime, false); err != nil {
		t.Error(err)
	}
	if err := checkRoutine(t, &pvTime, &prTime, false); err != nil {
		t.Error(err)
	}

	pvDuration := new(time.Duration)
	prDuration := new(time.Duration)
	*pvDuration = time.Duration(987654321 * time.Nanosecond)
	if err := checkRoutine(t, pvDuration, prDuration, false); err != nil {
		t.Error(err)
	}
	if err := checkRoutine(t, &pvDuration, &prDuration, false); err != nil {
		t.Error(err)
	}

	pvOffset := new(datetimeoffset.DateTimeOffset)
	prOffset := new(datetimeoffset.DateTimeOffset)
	*pvOffset = datetimeoffset.Now()
	if err := checkRoutine(t, pvOffset, prOffset, false); err != nil {
		t.Error(err)
	}
	if err := checkRoutine(t, &pvOffset, &prOffset, false); err != nil {
		t.Error(err)
	}

	// error
	var rError time.Time
	vError := datetimeoffset.Now()
	if err := checkRoutine(t, vError, &rError, false); err != nil {
		if strings.Contains(err.Error(), "value diffrent") {
			t.Error(err)
		}
	}
}

func TestArray(t *testing.T) {

	var rIntA [10]int
	vIntA := [10]int{1, 2, 3, 4, 5, 6, 7, 8, 9, math.MinInt32}
	if err := checkRoutine(t, vIntA, &rIntA, false); err != nil {
		t.Error(err)
	}

	var rIntS []int
	vIntS := []int{math.MinInt32, 9, 8, 7, 6, 5, 4, 3, 2, 1}
	if err := checkRoutine(t, vIntS, &rIntS, false); err != nil {
		t.Error(err)
	}

	var rUintA [10]uint
	vUintA := [10]uint{1, 2, 3, 4, 5, 6, 7, 8, 9, math.MaxUint32}
	if err := checkRoutine(t, vUintA, &rUintA, false); err != nil {
		t.Error(err)
	}

	var rUintS []uint
	vUintS := []uint{math.MaxUint32, 9, 8, 7, 6, 5, 4, 3, 2, 1}
	if err := checkRoutine(t, vUintS, &rUintS, false); err != nil {
		t.Error(err)
	}

	var rFloatA [5]float64
	vFloatA := [5]float64{1.2, 3.4, 5.6, 7.8, math.MaxFloat64}
	if err := checkRoutine(t, vFloatA, &rFloatA, false); err != nil {
		t.Error(err)
	}

	var rFloatS []float64
	vFloatS := []float64{math.MaxFloat64, 9.8, 7.6, 5.4, 3.2}
	if err := checkRoutine(t, vFloatS, &rFloatS, false); err != nil {
		t.Error(err)
	}

	var rBoolA [5]bool
	vBoolA := [5]bool{true, false, true, false, true}
	if err := checkRoutine(t, vBoolA, &rBoolA, false); err != nil {
		t.Error(err)
	}

	var rBoolS []bool
	vBoolS := []bool{false, true, false, true, false}
	if err := checkRoutine(t, vBoolS, &rBoolS, false); err != nil {
		t.Error(err)
	}

	var rStrA []string
	vStrA := []string{"this", "is", "string", "array", ".", "can", "you", "see", "?"}
	if err := checkRoutine(t, vStrA, &rStrA, false); err != nil {
		t.Error(err)
	}

	var rStrS []string
	vStrS := []string{"this", "is", "string", "slice", ".", "can", "you", "see", "?"}
	if err := checkRoutine(t, vStrS, &rStrS, false); err != nil {
		t.Error(err)
	}

	var rEmptyA [0]string
	vEmptyA := [0]string{}
	if err := checkRoutine(t, vEmptyA, &rEmptyA, false); err != nil {
		t.Error(err)
	}

	var rEmptyS []string
	vEmptyS := []string{}
	if err := checkRoutine(t, vEmptyS, &rEmptyS, false); err != nil {
		t.Error(err)
	}

	var rCharA [3]char.Char
	vCharA := [3]char.Char{'A', 'B', 'C'}
	if err := checkRoutine(t, vCharA, &rCharA, false); err != nil {
		t.Error(err)
	}

	var rCharS []char.Char
	vCharS := []char.Char{'C', 'B', 'A'}
	if err := checkRoutine(t, vCharS, &rCharS, false); err != nil {
		t.Error(err)
	}

	var rOffsetA [3]datetimeoffset.DateTimeOffset
	vOffsetA := [3]datetimeoffset.DateTimeOffset{datetimeoffset.Now(), datetimeoffset.Now(), datetimeoffset.Now()}
	if err := checkRoutine(t, vOffsetA, &rOffsetA, false); err != nil {
		t.Error(err)
	}

	var rOffsetS []datetimeoffset.DateTimeOffset
	vOffsetS := []datetimeoffset.DateTimeOffset{datetimeoffset.Now(), datetimeoffset.Now(), datetimeoffset.Now()}
	if err := checkRoutine(t, vOffsetS, &rOffsetS, false); err != nil {
		t.Error(err)
	}

	// pointer
	pvStrS := new([]string)
	prStrS := new([]string)
	*pvStrS = []string{"this", "is", "pointer", "string", "slice", ".", "can", "you", "see", "?"}
	if err := checkRoutine(t, pvStrS, prStrS, false); err != nil {
		t.Error(err)
	}
	if err := checkRoutine(t, &pvStrS, &prStrS, false); err != nil {
		t.Error(err)
	}
}

func TestStruct(t *testing.T) {
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
	rSt := st{}
	if err := checkRoutine(t, vSt, &rSt, false); err != nil {
		t.Error(err)
	}

	// pointer
	prSt := new(st)
	if err := checkRoutine(t, &vSt, &prSt, false); err != nil {
		t.Error(err)
	}
}

func TestMap(t *testing.T) {
	rMapInt := map[int]int{1: 2, 3: 4, math.MaxInt32: math.MinInt32}
	vMapInt := map[int]int{}
	if err := checkRoutine(t, rMapInt, &vMapInt, false); err != nil {
		t.Error(err)
	}

	rMapStr := map[string]float32{"this": 1.2, "is": 3.4, "float map": 56.789}
	vMapStr := map[string]float32{}
	if err := checkRoutine(t, rMapStr, &vMapStr, false); err != nil {
		t.Error(err)
	}
}

// for test
func checkRoutine(t *testing.T, in interface{}, out interface{}, isDebug bool) error {
	d, err := Serialize(in)
	if err != nil {
		return err
	}

	if isDebug {
		t.Log(in, " -- to byte --> ", d)
	}

	if err := Deserialize(out, d); err != nil {
		return err
	}

	i := getValue(in)
	o := getValue(out)
	if isDebug {
		t.Log("value [in]:", i, " [out]:", o)
	}

	if !reflect.DeepEqual(i, o) {
		return errors.New(fmt.Sprint("value different [in]:", in, " [out]:", out))
	}
	return nil
}

// for check value
func getValue(v interface{}) interface{} {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	return rv.Interface()
}
