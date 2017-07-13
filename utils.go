package args

import (
	"errors"
)

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

	switch t := i.(type) {
	case string:
		a.t = StringType
	case int:
		a.t = NumberType
	case bool:
		a.t = BoolType
	case map[string]interface{}:
		a.v = Map(t)
	case Map:
		a.t = MapType
	case []Argument:
		a.v = ArgumentList(t)
		a.t = ArgumentListType
	case ArgumentList:
		a.t = ArgumentListType
	case map[string]Argument:
		a.v = ArgumentMap(t)
		a.t = ArgumentMapType
	case ArgumentMap:
		a.t = ArgumentMapType
	case func(Argument) (Argument, error):
		a.t = CallType
		a.v = CallFunc(t)
	default:
		if _, ok := i.(Call); ok {
			a.t = CallType
		} else {

			if tt := findType(i); tt != nil {
				a.t = tt.t
			} else {
				a.Free()
				return nil, errors.New("invalid type")
			}

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
