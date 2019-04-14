package style

import (
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
	"strconv"

	"github.com/pkg/errors"
)

type ApplyFnc func(styles *Styles)

func ApplyNone(styles *Styles) {}

type Applier struct {
	fncs []ApplyFnc
}

func DecodeApplier(r io.Reader) (*Applier, error) {
	a := &Applier{
		fncs: []ApplyFnc{},
	}
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, errors.Wrap(err, "read-all")
	}
	raw, err := parseRaw(string(b))
	if err != nil {
		return nil, errors.Wrap(err, "parse-raw")
	}
	protoType := reflect.TypeOf(Styles{})
	for k, v := range raw {
		fnc, found, err := makeApplyFncV2(protoType, k, v, []int{})
		if err != nil {
			return nil, errors.Wrapf(err, "make-apply-fnc (%s, %s)", k, v)
		}
		if !found {
			continue
		}
		a.fncs = append(a.fncs, fnc)
	}
	return a, nil
}

func (a *Applier) Apply(styles *Styles) {
	for _, fnc := range a.fncs {
		fnc(styles)
	}
}

func makeApplyFncV1(key, val string) (ApplyFnc, error) {
	switch key {
	case "font-family":
		return func(styles *Styles) {
			styles.Font.Family = val
		}, nil
	case "border":
		b := Border{}
		err := b.UnmarshalStyle(val)
		if err != nil {
			return nil, err
		}
		return func(styles *Styles) {
			styles.Box.Border = b
		}, nil
	default:
		return ApplyNone, nil
	}
}

func makeApplyFncV2(rt reflect.Type, key, val string, indexPath []int) (ApplyFnc, bool, error) {
	//find field with tag style=key
	fmt.Printf("make-fnc (%s) -> (%s, %s)\n", rt.Name(), key, val)
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		//fmt.Printf("check field (%s) (style:%s)\n", field.Name, field.Tag.Get(""))
		currIndexPath := appendedCopy(indexPath, i)
		if t := field.Tag.Get("style"); t == key {
			fmt.Printf("make-fnc found field (%s) (%s)\n", field.Name, field.Type.Name())
			//check if type implements unmarshaler

			var setValue func(reflect.Value)
			um, impls := reflect.New(field.Type).Interface().(Unmarshaler)
			if impls {
				err := um.UnmarshalStyle(val)
				if err != nil {
					return nil, false, err
				}
				setValue = func(rv reflect.Value) {
					rv.Set(reflect.ValueOf(um).Elem())
				}
			} else {
				var err error
				setValue, err = makeSetValueFnc(field.Type.Kind(), val)
				if err != nil {
					return nil, false, errors.Wrapf(err, "make set value func (%s, %s)", key, val)
				}
			}

			return func(s *Styles) {
				fmt.Printf("apply (%v)\n", currIndexPath)
				rVal := reflect.ValueOf(s).Elem()
				for _, fIdx := range currIndexPath {
					field := rVal.Field(fIdx)
					fmt.Printf("    field %d (%s) \n", fIdx, field.Type().Name())
					rVal = rVal.Field(fIdx)
				}
				setValue(rVal)
			}, true, nil
		}
		if field.Type.Kind() == reflect.Struct {
			fnc, found, err := makeApplyFncV2(field.Type, key, val, currIndexPath)
			if err != nil {
				return nil, true, err
			} else if found {
				return fnc, true, nil
			}
		}
	}
	return nil, false, nil
}

func makeSetValueFnc(kind reflect.Kind, styleValue string) (func(v reflect.Value), error) {
	switch kind {
	case reflect.String:
		return func(v reflect.Value) {
			v.SetString(styleValue)
		}, nil
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(styleValue, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "parse-float (%s)", styleValue)
		}
		return func(v reflect.Value) {
			v.SetFloat(f)
		}, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := strconv.ParseInt(styleValue, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "parse-int (%s)", styleValue)
		}
		return func(v reflect.Value) {
			v.SetInt(n)
		}, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n, err := strconv.ParseUint(styleValue, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "parse-uint (%s)", styleValue)
		}
		return func(v reflect.Value) {
			v.SetUint(n)
		}, nil
	default:
		return nil, errors.Errorf("unsupported style kind (%s)", kind)
	}
}

func appendedCopy(sl []int, a int) []int {
	c := make([]int, len(sl)+1)
	for i, v := range sl {
		c[i] = v
	}
	c[len(c)-1] = a
	return c
}

/*
Family     string         `style:"font-family"`
	PointSize  float64        `style:"font-point-size"`
	Style      FontStyle      `style:"font-style"`
	Weight     FontWeight     `style:"font-weight"`
	Decoration FontDecoration `style:"font-decoration"`
*/
