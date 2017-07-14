package args

import (
	"reflect"

	multierror "github.com/hashicorp/go-multierror"
)

// TypeFromInterface returns a Type from a an interface
// or UndefinedType if it's unsupported
func TypeFromInterface(i interface{}) Type {
	switch i.(type) {
	case bool:
		return BoolType
	case int:
		return IntType
	case int64:
		return Int64Type
	case int32:
		return Int32Type
	case int16:
		return Int16Type
	case int8:
		return Int8Type
	case uint:
		return UintType
	case uint64:
		return Uint64Type
	case uint32:
		return Uint32Type
	case uint16:
		return Uint16Type
	case uint8:
		return Uint8Type
	case float32:
		return Float32Type
	case float64:
		return Float64Type
	case string:
		return StringType
	case Map:
		return MapType
	case []string:
		return StringSliceType
	case []byte:
		return ByteSliceType
	case []Argument:
		return ArgumentSliceType
	case ArgumentList:
		return ArgumentListType
	case ArgumentMap:
		return ArgumentMapType
	case func(ArgumentList) (Argument, error):
		return CallType
	default:
		if _, ok := i.(Call); ok {
			return CallType
		}
		if tt := findType(i); tt != nil && tt.e.Type != UndefinedType {
			return tt.e.Type
		}
		if _, ok := i.(error); ok {
			return ErrorType
		}
		return UndefinedType
	}

}

var reflectMap = reflect.TypeOf(Map{})
var reflectAMap = reflect.TypeOf(ArgumentMap{})
var reflectAList = reflect.TypeOf(ArgumentList{})
var reflectSSlice = reflect.TypeOf([]string{})
var reflectASlice = reflect.TypeOf([]Argument{})
var reflectError = reflect.TypeOf((*error)(nil)).Elem()

func TypeFromReflectType(t reflect.Type) Type {
	switch t.Kind() {
	case reflect.String:
		return StringType
	case reflect.Invalid:
		return UndefinedType
	case reflect.Bool:
		return BoolType
	case reflect.Int:
		return IntType
	case reflect.Int8:
		return Int8Type
	case reflect.Int16:
		return Int16Type
	case reflect.Int32:
		return Int32Type
	case reflect.Int64:
		return Int64Type
	case reflect.Uint:
		return UintType
	case reflect.Uint8:
		return Uint8Type
	case reflect.Uint16:
		return Uint16Type
	case reflect.Uint32:
		return Uint32Type
	case reflect.Uint64:
		return Uint64Type
	case reflect.Float32:
		return Float32Type
	case reflect.Float64:
		return Float64Type
	case reflect.Slice:
		if t == reflectAList {
			return ArgumentListType
		} else if t == reflectSSlice {
			return StringSliceType
		} else if t == reflectASlice {
			return ArgumentSliceType
		}
		return UndefinedType
	case reflect.Map:
		if t == reflectMap {
			return MapType
		}
		if (t.Key().Kind() == reflect.String && t.Elem().Kind() == reflect.Interface) || t == reflectAMap {
			return ArgumentMapType
		}
		return UndefinedType

	default:
		if tt := findType(t); tt != nil && tt.e.Type != UndefinedType {
			return tt.e.Type
		} else if t.Implements(reflectError) {
			return ErrorType
		}

		return UndefinedType
	}
}

func interfaceMapToArgumentMap(m map[string]interface{}) (ArgumentMap, error) {
	out := ArgumentMap{}
	var result error
	for k, v := range m {
		if a, e := New(v); e != nil {
			result = multierror.Append(result, e)
		} else {
			out[k] = a
		}
	}

	return out, result

}

// New makes a new argument from an interface or returns a error
func New(in ...interface{}) (Argument, error) {
	l := len(in)
	out := make([]Argument, l)
	for index, i := range in {
		if a, ok := i.(Argument); ok {
			return a, nil
		}

		a := argumentPool.Get().(*argument)
		a.v = i

		if i == nil {
			a.t = NilType
			return a, nil
		}

		if a.t = TypeFromInterface(i); a.t == UndefinedType || a.t == CallType {
			switch t := i.(type) {
			case map[string]interface{}:
				v, e := interfaceMapToArgumentMap(t)
				if e != nil {
					return nil, e
				}
				a.v = v
			case map[string]Argument:
				a.v = ArgumentMap(t)
			case func(ArgumentList) (Argument, error):
				a.v = CallFunc(t)
			}

		}

		out[index] = a
	}
	if l == 1 {
		return out[0], nil
	}
	return Must(out), nil
}

func Must(i ...interface{}) Argument {
	a, e := New(i...)
	if e != nil {
		panic(e)
	}
	return a
}

// NewOrNil create a new argument or returns a nil Arguments
// if an argument could not be created
func NewOrNil(i ...interface{}) Argument {
	if a, err := New(i...); err == nil {
		return a
	}
	return NilArgument()
}

// NilArgument create a new nil argument
func NilArgument() Argument {
	args := argumentPool.Get().(*argument)
	args.t = NilType
	args.v = nil
	return args
}

func Undefined() Argument {
	return argumentPool.Get().(*argument)
}
