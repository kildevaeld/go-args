package args

import (
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
	ErrorType
	NilType
)

func (t Type) String() string {
	switch t {
	case UndefinedType:
		return "undefined"
	case BoolType:
		return "bool"
	case IntType:
		return "int"
	case Int64Type:
		return "int64"
	case Int32Type:
		return "int32"
	case Int16Type:
		return "int16"
	case Int8Type:
		return "int8"
	case UintType:
		return "uint"
	case Uint64Type:
		return "uint64"
	case Uint32Type:
		return "uint32"
	case Uint16Type:
		return "uint16"
	case Uint8Type:
		return "uint8"
	case Float32Type:
		return "float32"
	case Float64Type:
		return "float64"
	case StringType:
		return "string"
	case MapType:
		return "map"
	case StringSliceType:
		return "string_slice"
	case ByteSliceType:
		return "byte_slice"
	case ArgumentSliceType:
		return "argument_slice"
	case CallType:
		return "call"
	case ArgumentListType:
		return "argument_list"
	case ArgumentMapType:
		return "argument_map"
	case ArgumentType:
		return "argument"
	case ErrorType:
		return "error"
	case NilType:
		return "nil"
	default:
		if v := findTypeWithType(t); v != nil {
			return v.e.Name
		}
	}
	return ""
}

type Freeable interface {
	Free()
}

type Argument interface {
	Freeable
	Type() Type
	Value() interface{}
	Valid() bool
	Is(...Type) bool
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
			return &argument{t: UndefinedType, v: nil}
		},
	}
}
