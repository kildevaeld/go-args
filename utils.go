package args

import multierror "github.com/hashicorp/go-multierror"

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
		if tt := findType(i); tt != nil {
			return tt.t
		}

		return UndefinedType
	}

}

func interfaceMapToArgumentMap(m map[string]interface{}) (ArgumentMap, error) {
	out := ArgumentMap{}
	var result error
	for k, v := range m {
		if a, e := NewArgument(v); e != nil {
			result = multierror.Append(result, e)
		} else {
			out[k] = a
		}
	}

	return out, result

}

func NewArgument(i interface{}) (Argument, error) {

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

	return a, nil
}

func NewArgumentOrNil(i interface{}) Argument {
	if a, err := NewArgument(i); err == nil {
		return a
	}
	return NilArgument()
}

func NilArgument() Argument {
	args := argumentPool.Get().(*argument)
	args.t = NilType
	args.v = nil
	return args
}
