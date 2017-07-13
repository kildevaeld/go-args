package args

import (
	"errors"
	"reflect"
	"sync"
)

type Type int

const (
	UndefinedType Type = iota
	BoolType
	IntType
	Int64Type
	Int32Type
	Int16Type
	Int8Type
	UintType
	Uint64Type
	Uint32Type
	Uint16Type
	Uint8Type
	Float32Type
	Float64Type
	StringType
	MapType
	StringSliceType
	ByteSliceType
	ArgumentSliceType
	CallType
	ArgumentListType
	ArgumentMapType
	ArgumentType
	NilType
)

type Freeable interface {
	Free()
}

type Argument interface {
	Freeable
	Type() Type
	Value() interface{}
	Valid() bool
}

type Map map[string]interface{}

type Call interface {
	Call(args ArgumentList) (Argument, error)
}

type CallFunc func(args ArgumentList) (Argument, error)

func (c CallFunc) Call(args ArgumentList) (Argument, error) {
	return c(args)
}

var argumentPool sync.Pool

func init() {
	argumentPool = sync.Pool{
		New: func() interface{} {
			return &argument{}
		},
	}
}

type typemap struct {
	t Type
	v reflect.Type
}

var _types []typemap

func findType(v interface{}) *typemap {
	val := reflect.TypeOf(v)
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

func Register(v interface{}, t Type) error {
	val := reflect.TypeOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	for _, r := range _types {
		if r.t == t || r.v == val {
			return errors.New("type already registered")
		}
	}

	if t <= NilType {
		return errors.New("cannot register core type")
	}

	_types = append(_types, typemap{
		t: t,
		v: val,
	})

	return nil
}
