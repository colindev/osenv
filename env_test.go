package env

import (
	"os"
	"testing"
)

type X struct {
	Int     int     `env:"int,INT"`
	Int8    int8    `env:"int8"`
	Int16   int16   `env:"int16"`
	Int32   int32   `env:"int32"`
	Int64   int64   `env:"int64"`
	Uint    uint    `env:"uint"`
	Uint8   uint8   `env:"uint8"`
	Uint16  uint16  `env:"uint16"`
	Uint32  uint32  `env:"uint32"`
	Uint64  uint64  `env:"uint64"`
	Float32 float32 `env:"float32"`
	Float64 float64 `env:"float64"`
	String  string  `env:"string"`

	Omit []string `env:"-"`
}

type Error struct {
	Slice []string `env:"xxx"`
}

func Test(t *testing.T) {

	var (
		x  *X
		xx X
	)
	os.Setenv("int", "1")
	os.Setenv("INT", "2")
	os.Setenv("int8", "3")
	os.Setenv("int16", "4")
	os.Setenv("int32", "5")
	os.Setenv("int64", "6")
	os.Setenv("uint", "7")
	os.Setenv("uint8", "8")
	os.Setenv("uint16", "9")
	os.Setenv("uint32", "10")
	os.Setenv("uint64", "11")
	os.Setenv("float32", "12")
	os.Setenv("float64", "13")
	os.Setenv("string", "abc")

	// omit y

	if err := LoadTo(struct{}{}); err == nil {
		t.Error("判斷錯誤 struct value")
	}
	if err := LoadTo(123); err == nil {
		t.Error("判斷錯誤 int")
	}
	if err := LoadTo(nil); err == nil {
		t.Error("判斷錯誤 nil")
	}
	if err := LoadTo(x); err == nil {
		t.Error("判斷錯誤 (*X)(nil)")
	}
	if err := LoadTo(&Error{}); err == nil {
		t.Error("結構內有不支援型別slice 不應該沒檢查到")
	}

	//fmt.Println(LoadTo(&[]string{}))
	//fmt.Println(LoadTo(&map[string]interface{}{}))

	if err := LoadTo(&xx); err != nil {
		t.Error(err)
	}

	t.Logf("%#v", xx)
	if xx.Int != int(2) {
		t.Errorf("int load fail %#v", xx.Int)
	}
	if xx.Int8 != int8(3) {
		t.Errorf("int8 load fail %#v", xx.Int8)
	}
	if xx.Int16 != int16(4) {
		t.Errorf("int14 load fail %#v", xx.Int16)
	}
	if xx.Int32 != int32(5) {
		t.Errorf("int32 load fail %#v", xx.Int32)
	}
	if xx.Int64 != int64(6) {
		t.Errorf("int64 load fail %#v", xx.Int64)
	}
	if xx.Uint != uint(7) {
		t.Errorf("uint load fail %#v", xx.Uint)
	}
	if xx.Uint8 != uint8(8) {
		t.Errorf("uint8 load fail %#v", xx.Uint8)
	}
	if xx.Uint16 != uint16(9) {
		t.Errorf("uint16 load fail %#v", xx.Uint16)
	}
	if xx.Uint32 != uint32(10) {
		t.Errorf("uint32 load fail %#v", xx.Uint32)
	}
	if xx.Uint64 != uint64(11) {
		t.Errorf("uint64 load fail %#v", xx.Uint64)
	}
	if xx.Float32 != float32(12) {
		t.Errorf("float32 load fail %#v", xx.Float32)
	}
	if xx.Float64 != float64(13) {
		t.Errorf("float64 load fail %#v", xx.Float64)
	}
	if xx.String != "abc" {
		t.Errorf("string load fail %#v", xx.String)
	}
}
