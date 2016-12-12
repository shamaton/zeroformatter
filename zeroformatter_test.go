package zeroformatter

import (
	"fmt"
	"math"
	"reflect"
	"testing"
	"time"
	"unsafe"

	"github.com/shamaton/zeroformatter/char"
)

func TestZeroformatter(t *testing.T) {
	f := func(in interface{}, out interface{}, isDispByte bool) error {
		d, err := Serialize(in)
		if err != nil {
			return err
		}
		if isDispByte {
			t.Log(in, " -- to byte --> ", d)
		}
		if err := Deserialize(out, d); err != nil {
			return err
		}
		return nil
	}
	_p := func(in interface{}, out interface{}) string {
		return fmt.Sprint("value different [in]:", in, " [out]:", out)
	}

	var rInt8 int8
	vInt8 := int8(-8)
	if err := f(vInt8, &rInt8, false); err != nil {
		t.Error(err)
	}
	if vInt8 != rInt8 {
		t.Error(_p(vInt8, rInt8))
	}
	t.Log(rInt8)

	var rInt16 int16
	vInt16 := int16(-16)
	if err := f(vInt16, &rInt16, false); err != nil {
		t.Error(err)
	}
	if vInt16 != rInt16 {
		t.Error(_p(vInt16, rInt16))
	}

	var rInt int
	vInt := -65535
	if err := f(vInt, &rInt, false); err != nil {
		t.Error(err)
	}
	if vInt != rInt {
		t.Error(_p(vInt, rInt))
	}

	var rInt32 int32
	vInt32 := int32(-32)
	if err := f(vInt32, &rInt32, false); err != nil {
		t.Error(err)
	}
	if vInt32 != rInt32 {
		t.Error(_p(vInt32, rInt32))
	}

	var rInt64 int64
	vInt64 := int64(-64)
	if err := f(vInt64, &rInt64, false); err != nil {
		t.Error(err)
	}
	if vInt64 != rInt64 {
		t.Error(_p(vInt64, rInt64))
	}
	t.Log(rInt64)

	var rUint8 uint8
	vUint8 := uint8(math.MaxUint8)
	if err := f(vUint8, &rUint8, false); err != nil {
		t.Error(err)
	}
	if vUint8 != rUint8 {
		t.Error(_p(vUint8, rUint8))
	}
	t.Log(rUint8)

	var rUint16 uint16
	vUint16 := uint16(math.MaxUint16)
	if err := f(vUint16, &rUint16, false); err != nil {
		t.Error(err)
	}
	if vUint16 != rUint16 {
		t.Error(_p(vUint16, rUint16))
	}
	t.Log(rUint16)

	var rUint uint
	vUint := uint(math.MaxUint32 / 2)
	if err := f(vUint, &rUint, false); err != nil {
		t.Error(err)
	}
	if vUint != rUint {
		t.Error(_p(vUint, rUint))
	}
	t.Log(rUint)

	var rUint32 uint32
	vUint32 := uint32(math.MaxUint32)
	if err := f(vUint32, &rUint32, false); err != nil {
		t.Error(err)
	}
	if vUint32 != rUint32 {
		t.Error(_p(vUint32, rUint32))
	}
	t.Log(rUint32)

	var rUint64 uint64
	vUint64 := uint64(math.MaxUint64)
	if err := f(vUint64, &rUint64, false); err != nil {
		t.Error(err)
	}
	if vUint64 != rUint64 {
		t.Error(_p(vUint64, rUint64))
	}
	t.Log(rUint64)

	var rFloat32 float32
	vFloat32 := float32(math.MaxFloat32)
	if err := f(vFloat32, &rFloat32, false); err != nil {
		t.Error(err)
	}
	if vFloat32 != rFloat32 {
		t.Error(_p(vFloat32, rFloat32))
	}
	t.Log(rFloat32)

	var rFloat64 float64
	vFloat64 := float64(math.MaxFloat64)
	if err := f(vFloat64, &rFloat64, false); err != nil {
		t.Error(err)
	}
	if vFloat64 != rFloat64 {
		t.Error(_p(vFloat64, rFloat64))
	}
	t.Log(rFloat64)

	var rBool bool
	vBool := true
	if err := f(vBool, &rBool, false); err != nil {
		t.Error(err)
	}
	if vBool != rBool {
		t.Error(_p(vBool, rBool))
	}
	t.Log(rBool)

	var rChar char.Char
	vChar := char.Char('Z')
	if err := f(vChar, &rChar, false); err != nil {
		t.Error(err)
	}
	if vChar != rChar {
		t.Error(_p(vChar, rChar))
	}
	t.Logf("%#U", rChar)

	var rString string
	vString := "this string serialize and deserialize."
	if err := f(vString, &rString, false); err != nil {
		t.Error(err)
	}
	if vString != rString {
		t.Error(_p(vString, rString))
	}
	t.Log(rString)

	var rTime time.Time
	vTime := time.Now()
	if err := f(vTime, &rTime, false); err != nil {
		t.Error(err)
	}
	if vTime != rTime {
		t.Error(_p(vTime, rTime))
	}
	t.Log(rTime)

	var rDuration time.Duration
	vDuration := time.Duration(12*time.Hour + 34*time.Minute + 56*time.Second + 78*time.Nanosecond)
	if err := f(vDuration, &rDuration, false); err != nil {
		t.Error(err)
	}
	if vDuration != rDuration {
		t.Error(_p(vDuration, rDuration))
	}
	t.Log(rDuration)

	// todo : more array/slice test cases
	var rIntArr []int
	vIntArr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, math.MinInt32}
	if err := f(vIntArr, &rIntArr, false); err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(vIntArr, rIntArr) {
		t.Error(_p(vIntArr, rIntArr))
	}
	t.Log(rIntArr)

	var rStrArr []string
	vStrArr := []string{"this", "is", "string", "array", ".", "can", "you", "see", "?"}
	if err := f(vStrArr, &rStrArr, false); err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(vStrArr, rStrArr) {
		t.Error(_p(vStrArr, rStrArr))
	}
	t.Log(rStrArr)

	var rArrEmpty []string
	vArrEmpty := []string{}
	if err := f(vArrEmpty, &rArrEmpty, false); err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(vArrEmpty, rArrEmpty) {
		t.Error(_p(vArrEmpty, rArrEmpty))
	}
	t.Log(rArrEmpty)

	/*
		var _rUint8 int8
		_vUint8 := int8(-8)
		if err := f(_vUint8, &_rUint8, false); err != nil {
			t.Error(err)
		}
		if _vUint8 != _rUint8 {
			t.Error(_p(_vUint8, _rUint8))
		}
		t.Log(_rUint8)
	*/
	type childchild struct {
		String string
		Floats []float32
	}
	type child struct {
		Int   int
		Time  time.Time
		Child childchild
	}
	type st struct {
		Int16  int16
		Int    int
		Int64  int64
		Uint16 uint16
		Uint   uint
		Uint64 uint64
		Float  float32
		Double float64
		Bool   bool
		Uint8  byte
		Int8   int8
		String string
		Time   time.Time
		Child  child
	}
	vSt := &st{
		Int:    -32,
		Int8:   -8,
		Int16:  -16,
		Int64:  -64,
		Uint:   32,
		Uint8:  8,
		Uint16: 16,
		Uint64: 64,
		Float:  1.23,
		Double: 2.3456,
		Bool:   true,
		String: "hello",
		Time:   time.Now(),
		Child: child{
			Int:   1234567,
			Time:  time.Now(),
			Child: childchild{String: "this is child in child", Floats: []float32{1.2, 3.4, 5.6}},
		},
	}
	rSt := st{}
	if err := f(vSt, &rSt, false); err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(*vSt, rSt) {
		t.Error(_p(*vSt, rSt))
	}

	t.Log(rSt)
	t.Log("stst ", unsafe.Sizeof(*vSt), " : ", unsafe.Sizeof(rSt))

	// pointer test mmmm...
	hoge := new(int)
	*hoge = 123
	fuga := new(int)
	rrrr := reflect.ValueOf(&fuga)
	t.Log(rrrr.Type().Elem())
	if err := f(&hoge, &fuga, false); err != nil {
		t.Error(err)
	}
	t.Log(hoge, *hoge, fuga, *fuga)

}
