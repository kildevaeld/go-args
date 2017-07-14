package args

import (
	"errors"
	"reflect"
)

type Extension struct {
	Type        Type
	Name        string
	Describe    func(i interface{}) *ArgumentDescription
	MarshalJSON func(i interface{}) ([]byte, error)
	MarshalYAML func(i interface{}) (interface{}, error)
}

type typemap struct {
	//t Type
	v reflect.Type
	//s string
	e Extension
}

var _types []typemap

func findType(v interface{}) *typemap {
	var val reflect.Type
	if vv, ok := v.(reflect.Type); ok {
		val = vv
	} else {
		val = reflect.TypeOf(v)
	}

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	for _, m := range _types {
		if m.v == val {
			return &m
		}
	}

	return nil
}

func findTypeWithType(t Type) *typemap {
	for _, v := range _types {
		if t == v.e.Type {
			return &v
		}
	}
	return nil
}

func Register(v interface{}, e Extension) error {
	val := reflect.TypeOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	for _, r := range _types {
		if r.e.Type == e.Type || r.v == val {
			return errors.New("type already registered")
		}
	}

	if e.Type <= NilType {
		return errors.New("cannot register core type")
	}

	_types = append(_types, typemap{
		//t: t,
		v: val,
		//s: name,
		e: e,
	})

	return nil
}
