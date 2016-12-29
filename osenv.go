package osenv

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

	return eachStructFields(v, func(rsf reflect.StructField, rv reflect.Value, tags []string) error {

		n := len(tags)
		if n == 0 {
			return fmt.Errorf("env: %s too less args", rsf.Name)
		} else if n > 2 {
			return fmt.Errorf("env: %s too many args", rsf.Name)
		}

		arg := os.Getenv(strings.TrimSpace(tags[0]))
		if arg == "" && n == 2 {
			arg = strings.TrimSpace(tags[1])
		}

		if err := setField(rv, arg); err != nil {
			return fmt.Errorf("env: set field(%s, %s) %v", rsf.Name, arg, err)
		}

		return nil
	})

}

func eachStructFields(v interface{}, fn func(reflect.StructField, reflect.Value, []string) error) error {

	rv := reflect.ValueOf(v)

	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &invalidValueError{reflect.TypeOf(v)}
	}

	stv := rv.Type().Elem()
	rvv := rv.Elem()

NEXT:
	for i := 0; i < stv.NumField(); i++ {
		stf := stv.Field(i)
		tags := stf.Tag.Get(tagName)
		if tags == "-" || tags == "" {
			continue
		}
		if err := fn(stf, rvv.Field(i), strings.Split(tags, ",")); err != nil {
			return err
		}
		continue NEXT
	}

	return nil
}

func setField(v reflect.Value, envArg string) error {

	if !v.CanSet() {
		return nil
	}
	switch v.Kind() {
	case reflect.String:
		v.SetString(envArg)
	case reflect.Bool:
		n, err := strconv.ParseBool(envArg)
		if err != nil {
			return err
		}
		v.SetBool(n)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v.Type().Name() == "Duration" {
			d, err := time.ParseDuration(envArg)
			if err != nil {
				return err
			}
			v.SetInt(int64(d))
			return nil
		}

		n, err := strconv.ParseInt(envArg, 10, 64)
		if err != nil {
			return err
		}
		v.SetInt(n)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n, err := strconv.ParseUint(envArg, 10, 64)
		if err != nil {
			return err
		}
		v.SetUint(n)
	case reflect.Float32, reflect.Float64:
		n, err := strconv.ParseFloat(envArg, 64)
		if err != nil {
			return err
		}
		v.SetFloat(n)
	default:
		return fmt.Errorf("Invalid type %s", v.Kind().String())
	}

	return nil
}
