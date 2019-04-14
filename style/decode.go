package style

import (
	"io"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func trimWS(s string) string {
	return strings.Trim(s, " \r\n\t")
}

func parseRaw(s string) (map[string]string, error) {
	raw := map[string]string{}
	styleStrs := strings.Split(s, ";")
	for _, styleStr := range styleStrs {
		styleStr = trimWS(styleStr)
		if len(styleStr) == 0 {
			continue
		}
		styleKV := strings.Split(styleStr, ":")
		if len(styleKV) != 2 {
			return nil, errors.Errorf("invalid style syntax (%s) must be of (key:val)", styleStr)
		}
		raw[trimWS(styleKV[0])] = trimWS(styleKV[1])
	}
	return raw, nil
}

type Unmarshaler interface {
	UnmarshalStyle(v string) error
}

type Decoder struct {
	reader io.Reader
	raw    map[string]string
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		reader: r,
	}
}

func (d *Decoder) Decode(styles *Styles) error {
	b, err := ioutil.ReadAll(d.reader)
	if err != nil {
		return err
	}
	d.raw, err = parseRaw(string(b))
	if err != nil {
		return errors.Wrap(err, "parse raw")
	}
	err = d.decode(styles, "")
	if err != nil {
		return err
	}
	return nil
}

func (d *Decoder) decode(v interface{}, styleName string) error {
	styleValue, containsStyle := d.raw[styleName]
	if containsStyle {
		if um, impls := v.(Unmarshaler); impls {
			return um.UnmarshalStyle(styleValue)
		}
	}
	k := reflect.ValueOf(v).Kind()
	if k != reflect.Ptr {
		return errors.Errorf("passed interface is not pointer but (%s)", k)
	}
	rVal := reflect.ValueOf(v).Elem()
	if rVal.Kind() == reflect.Struct {
		return d.decodePtrToStruct(v)
	}
	if !containsStyle {
		return nil
	}
	switch rVal.Kind() {
	case reflect.String:
		rVal.SetString(styleValue)
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(styleValue, 64)
		if err != nil {
			return errors.Wrapf(err, "parse-float (%s)", styleValue)
		}
		rVal.SetFloat(f)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := strconv.ParseInt(styleValue, 10, 64)
		if err != nil {
			return errors.Wrapf(err, "parse-int (%s)", styleValue)
		}
		rVal.SetInt(n)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n, err := strconv.ParseUint(styleValue, 10, 64)
		if err != nil {
			return errors.Wrapf(err, "parse-uint (%s)", styleValue)
		}
		rVal.SetUint(n)
	default:
		return errors.Errorf("cannot decode %s", rVal.Kind())
	}
	return nil
}

func (d *Decoder) decodePtrToStruct(v interface{}) error {
	rVal := reflect.ValueOf(v).Elem()
	rType := rVal.Type()
	for i := 0; i < rVal.NumField(); i++ {
		field := rVal.Field(i)
		if !field.CanAddr() {
			continue
		}
		if !field.CanInterface() {
			continue
		}
		err := d.decode(field.Addr().Interface(), rType.Field(i).Tag.Get("style"))
		if err != nil {
			return err
		}
	}
	return nil
}
