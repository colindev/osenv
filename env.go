package env

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const tagName = "env"

type invalidValueError struct {
	Type reflect.Type
}

func (e *invalidValueError) Error() string {

	if e.Type == nil {
		return "env: LoadTo(nil)"
	}

	if e.Type.Kind() != reflect.Ptr {
		return "env: LoadTo(non-pointer " + e.Type.String() + ")"
	}

	return "env: LoadTo(nil " + e.Type.String() + ")"
}

// LoadTo 將 環境變數載入 struct 內
func LoadTo(v interface{}) error {

	rv := reflect.ValueOf(v)

	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &invalidValueError{reflect.TypeOf(v)}
	}

	stv := rv.Type().Elem()
	rvv := rv.Elem()

	for i := 0; i < stv.NumField(); i++ {
		stf := stv.Field(i)
		tags := stf.Tag.Get(tagName)
		if tags == "-" {
			continue
		}
		for _, tag := range strings.Split(tags, ",") {
			if err := setField(rvv, i, strings.Trim(tag, " ")); err != nil {
				return err
			}
			continue
		}
	}

	return nil
}

func setField(v reflect.Value, i int, envName string) error {

	f := v.Field(i)
	s := os.Getenv(envName)
	switch f.Kind() {
	case reflect.String:
		f.SetString(s)
	case reflect.Bool:
		n, err := strconv.ParseBool(s)
		if err != nil {
			return err
		}
		f.SetBool(n)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if f.Type().Name() == "Duration" {
			d, err := time.ParseDuration(s)
			if err != nil {
				return err
			}
			f.SetInt(int64(d))
			return nil
		}

		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		f.SetInt(n)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return err
		}
		f.SetUint(n)
	case reflect.Float32, reflect.Float64:
		n, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return err
		}
		f.SetFloat(n)
	default:
		return fmt.Errorf("Invalid type %s", f.Kind().String())
	}

	return nil
}
