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

// ToString from interface
func ToString(v interface{}) string {
	ret := []string{}
	rv := reflect.ValueOf(v)
	eachStructFields(rv, func(field reflect.StructField, rv2 reflect.Value, tags []string) error {
		if tags[0] != "" {

			var value string

			switch rv2.Type().Name() {
			case "Time":
				// 有點多餘,但先保留著
				value = rv2.Interface().(time.Time).String()
			default:
				value = fmt.Sprintf("%v", rv2.Interface())
			}
			ret = append(ret, fmt.Sprintf("%s=%s", tags[0], value))
		}
		return nil
	})
	return strings.Join(ret, "\n")
}

// LoadTo 將 環境變數載入 struct 內
func LoadTo(v interface{}) error {
	rv := reflect.ValueOf(v)
	return eachStructFields(rv, func(field reflect.StructField, rv2 reflect.Value, tags []string) error {
		arg := strings.Split(os.Getenv(strings.TrimSpace(tags[0])), ",")
		if arg[0] == "" && len(tags) > 1 {
			arg = tags[1:]
		}
		if err := setField(rv2, arg); err != nil {
			return fmt.Errorf("env: set field(%s, %s) %v", field.Name, arg, err)
		}
		return nil
	})
}

func eachStructFields(rv reflect.Value, fn func(reflect.StructField, reflect.Value, []string) error) error {
	var rvt reflect.Type
	var rv2 reflect.Value
	switch rv.Kind() {
	case reflect.Ptr:
		if rv.IsNil() {
			return &invalidValueError{reflect.TypeOf(rv)}
		}
		rvt = rv.Type().Elem()
		rv2 = rv.Elem()
	case reflect.Struct:
		rvt = rv.Type()
		rv2 = rv
	default:
		return &invalidValueError{reflect.TypeOf(rv)}
	}

	for i := 0; i < rvt.NumField(); i++ {
		fieldType := rvt.Field(i)
		fieldValue := rv2.Field(i)
		if fieldValue.Kind() == reflect.Struct {
			err := eachStructFields(fieldValue, fn)
			if err != nil {
				return &invalidValueError{reflect.TypeOf(rv)}
			}
		}

		tagVal := fieldType.Tag.Get(tagName)

		if tagVal == "-" || tagVal == "" {
			continue
		}
		tags := strings.Split(tagVal, ",")
		n := len(tags)
		if n == 0 {
			return fmt.Errorf("env: %s too less args", fieldType.Name)
		} else if err := fn(fieldType, fieldValue, tags); err != nil {
			return err
		}
	}

	return nil
}

func setField(rv2 reflect.Value, envArg []string) error {
	if !rv2.CanSet() {
		return nil
	}
	switch rv2.Kind() {
	case reflect.String:
		rv2.SetString(envArg[0])
	case reflect.Bool:
		n, err := strconv.ParseBool(envArg[0])
		if err != nil {
			return err
		}
		rv2.SetBool(n)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if rv2.Type().Name() == "Duration" {
			d, err := time.ParseDuration(envArg[0])
			if err != nil {
				return err
			}
			rv2.SetInt(int64(d))
			return nil
		}
		n, err := strconv.ParseInt(envArg[0], 10, 64)
		if err != nil {
			return err
		}
		rv2.SetInt(n)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n, err := strconv.ParseUint(envArg[0], 10, 64)
		if err != nil {
			return err
		}
		rv2.SetUint(n)
	case reflect.Float32, reflect.Float64:
		n, err := strconv.ParseFloat(envArg[0], 64)
		if err != nil {
			return err
		}
		rv2.SetFloat(n)
	case reflect.Slice:
		switch rv2.Interface().(type) {
		case []int:
			t := reflect.TypeOf([]int{})
			s := reflect.MakeSlice(t, 0, 0)
			p, err := sliceAtoi(envArg)
			if err != nil {
				return err
			}
			a := append(s.Interface().([]int), p...)
			rv2.Set(reflect.ValueOf(a))
		case []string:
			t := reflect.TypeOf([]string{})
			s := reflect.MakeSlice(t, 0, 0)
			a := append(s.Interface().([]string), envArg...)
			rv2.Set(reflect.ValueOf(a))
		default:
			return fmt.Errorf("Invalid slice type %s", rv2.Type().String())
		}
	case reflect.Struct:
		if rv2.Type().Name() == "Time" {
			t, err := time.Parse(time.RFC3339, envArg[0]) // "2012-11-01T22:08:41+00:00"
			if err != nil {
				return err
			}
			rv2.Set(reflect.ValueOf(t)) //time.Now()
		}
		return nil
	default:
		return fmt.Errorf("Invalid type %s", rv2.Kind().String())
	}

	return nil
}

func sliceAtoi(as []string) ([]int, error) {
	is := make([]int, 0, len(as))
	for _, a := range as {
		i, err := strconv.Atoi(a)
		if err != nil {
			return is, err
		}
		is = append(is, i)
	}
	return is, nil
}
