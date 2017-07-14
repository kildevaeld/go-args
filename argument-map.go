package args

import (
	"fmt"

	multierror "github.com/hashicorp/go-multierror"
)

type ArgumentMap map[string]Argument

func (a ArgumentMap) Len() int {
	return len(a)
}

func (a *ArgumentMap) Free() {
	for _, aa := range *a {
		aa.Free()
	}
	*a = ArgumentMap{}
}

func (a ArgumentMap) ToInterfaceMap() map[string]interface{} {
	out := make(map[string]interface{})
	for key, i := range a {
		out[key] = i.Value()
	}
	return out
}

type CheckMap map[string]Type

func (a ArgumentMap) GetField(field string, t Type) Argument {
	f, ok := a[field]
	if !ok {
		return nil
	}
	if f.Type() == t {
		return f
	} else if f.Type() == ArgumentListType {
		if first := f.Value().(ArgumentList).First(); first != nil {
			if first.Type() == t {
				return first
			}
		}
	}
	return nil
}

func (a ArgumentMap) CheckField(field string, t Type) bool {
	f, ok := a[field]
	if !ok {
		return false
	}
	if f.Type() == t {
		return true
	} else if f.Type() == ArgumentListType {
		if first := f.Value().(ArgumentList).First(); first != nil {
			return first.Type() == t
		}
	} else if f.Type() == ArgumentSliceType {
		if first := f.Value().([]Argument)[0]; first != nil {
			return first.Type() == t
		}
	}
	return false
}

func (a ArgumentMap) Check(m CheckMap) error {
	var result error

	for k, v := range m {
		if !a.CheckField(k, v) {
			result = multierror.Append(result, fmt.Errorf("field not found: %s", k))
		}
	}
	return result
}

func NewMap(m map[string]interface{}) (ArgumentMap, error) {
	out := make(map[string]Argument)
	var result error
	for k, v := range m {
		a, e := New(v)
		if e != nil {
			result = multierror.Append(result, e)
			a = Undefined()
		}
		out[k] = a
	}

	return out, result
}

func MustMap(m map[string]interface{}) ArgumentMap {
	v, e := NewMap(m)
	if e != nil {
		panic(e)
	}
	return v
}
