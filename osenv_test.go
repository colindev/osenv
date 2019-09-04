package osenv

import (
	"fmt"
	"os"
	"testing"
	"time"
)

type X2 struct {
	Sub int `env:"sub,-2"`
}

type X struct {
	unexport string
	NoTag    string
	Int      int           `env:"int,-1"`
	Int8     int8          `env:"int8"`
	Int16    int16         `env:"int16"`
	Int32    int32         `env:"int32"`
	Int64    int64         `env:"int64"`
	Uint     uint          `env:"uint"`
	Uint8    uint8         `env:"uint8,8"`
	Uint16   uint16        `env:"uint16"`
	Uint32   uint32        `env:"uint32"`
	Uint64   uint64        `env:"uint64"`
	Float32  float32       `env:"float32"`
	Float64  float64       `env:"float64"`
	String   string        `env:"string"`
	Bool     bool          `env:"bool,true"`
	Duration time.Duration `env:"duration"`
	Time     time.Time     `env:"time,2012-11-01T22:08:41+00:00"`
	Omit     []string      `env:"-"`
	Slice1   []string      `env:"slice1,a,b,c"`
	Slice2   []int         `env:"slice2,1,2,3"`
	X        X2
}

func TestHelp(t *testing.T) {

	// TODO test output
	Help(X{}, os.Stdout)
}

func Test(t *testing.T) {

	var (
		x  *X
		xx X
	)

	os.Setenv("unexport", "unexport")
	os.Setenv("NoTag", "NoTag")
	os.Setenv("int", "1")
	os.Setenv("INT", "2")
	os.Setenv("int8", "3")
	os.Setenv("int16", "4")
	os.Setenv("int32", "5")
	os.Setenv("int64", "6")
	os.Setenv("uint", "7")
	// 不設定 uint8, 測試預設
	os.Setenv("uint16", "9")
	os.Setenv("uint32", "10")
	os.Setenv("uint64", "11")
	os.Setenv("float32", "12")
	os.Setenv("float64", "13")
	os.Setenv("string", "abc")
	os.Setenv("duration", "3h")
	os.Setenv("slice1", "a,b,c")
	os.Setenv("slice2", "1,2,3")
	os.Setenv("sub", "-3")

	if err := LoadTo(struct{}{}); err != nil {
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
	if err := LoadTo(&xx); err != nil {
		t.Error(err)
	}

	Debug = true
	if err := LoadTo(xx); err != nil {
		t.Error(err)
	}
	Debug = false

	t.Logf("%#v", xx)
	if xx.unexport != "" {
		t.Errorf("unexport can't load %#v", xx.unexport)
	}
	if xx.NoTag != "" {
		t.Errorf("NoTag can't load %#v", xx.NoTag)
	}
	if xx.Int != int(1) {
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
	if !xx.Bool {
		t.Errorf("bool load fail %#v", xx.Bool)
	}
	t.Log("duration is ", xx.Duration.String())
	if xx.Duration != time.Hour*3 {
		t.Errorf("duration load fail %#v", time.Duration(xx.Duration))
	}

	t.Log(xx.Time.String(), "// time.String()")
	t.Log(xx.Time.Format(time.RFC3339), "// time.Format(RFC3339)")
	t.Log(fmt.Sprintf("%v", xx.Time), `// fmt.Sprintf("%v", Time)`)

	raw := ToString(&xx)
	expectTime, _ := time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
	expectRaw := fmt.Sprintf(`int=1
int8=3
int16=4
int32=5
int64=6
uint=7
uint8=8
uint16=9
uint32=10
uint64=11
float32=12
float64=13
string=abc
bool=true
duration=3h0m0s
time=%s
slice1=[a b c]
slice2=[1 2 3]
sub=-3`, expectTime.String())
	t.Log(raw)
	if raw != expectRaw {
		t.Errorf("want:\n%s\ngot:\n%s", expectRaw, raw)
	}
}
