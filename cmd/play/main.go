package main

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/pkg/errors"
)

type Bazer interface {
	Baz(v string) error
}

type StyleBazer struct {
	private string
	Public  string
}

func (sb *StyleBazer) Baz(v string) error {
	sb.private = "private-styler"
	sb.Public = "public-styler " + v
	return nil
}

type F float64

type Bar struct {
	S    string     `style:"s"`
	F    F          `style:"f"`
	SBaz StyleBazer `style:"sbaz"`
}

type Foo struct {
	B Bar
	S string
}

func (f Foo) Print() {
	b, _ := json.MarshalIndent(f, "", "  ")
	fmt.Printf("Foo:\n%s\n", b)
}

func main() {
	f := Foo{
		B: Bar{
			S: "john",
			F: 42.42,
		},
		S: "doe",
	}
	f.Print()
	err := forge(&f, "")
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	f.Print()
}

func forge(v interface{}, ident string) error {
	//must be pointer
	k := reflect.ValueOf(v).Kind()
	if k != reflect.Ptr {
		return errors.Errorf("passed interface is not pointer but (%s)", k)
	}
	rVal := reflect.ValueOf(v).Elem()
	fmt.Printf("%svisit (%s)\n", ident, rVal.Kind())

	//ptr to struct?
	if rVal.Kind() == reflect.Struct {
		rType := rVal.Type()
		fmt.Printf("%siterate struct (%s)\n", ident, rType.Name())
		for i := 0; i < rVal.NumField(); i++ {
			field := rVal.Field(i)
			if !field.CanAddr() {
				continue
			}
			if !field.CanInterface() {
				continue
			}

			fieldType := rType.Field(i)
			fmt.Printf("%sfield %d: (%s, %s, %s)\n", ident, i, fieldType.Name, fieldType.Type.Name(), fieldType.Tag.Get("style"))
			ptrToField := field.Addr()
			bazer, ok := ptrToField.Interface().(Bazer)
			if ok {
				fmt.Printf("%s - is Bazer - call\n", ident)
				return bazer.Baz("do-baz")
			}

			err := forge(field.Addr().Interface(), ident+"  ")
			if err != nil {
				return err
			}
		}
		return nil
	}

	if !rVal.CanSet() {
		return errors.Errorf("cannot set")
	}
	switch rVal.Kind() {
	case reflect.String:
		rVal.SetString("baz")
	case reflect.Float32, reflect.Float64:
		rVal.SetFloat(3.14)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		rVal.SetInt(-314)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		rVal.SetInt(314)
	default:
		fmt.Printf("%sdon't handle (%s)\n", ident, k)
	}

	return nil
}
