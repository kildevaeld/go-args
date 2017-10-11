package args

import (
	"fmt"

	multierror "github.com/hashicorp/go-multierror"
)

type ArgumentList []Argument

func (a ArgumentList) Len() int {
	return len(a)
}
func (a ArgumentList) First() Argument {
	if len(a) > 0 {
		return a[0]
	}
	return nil
}

func (a ArgumentList) Check(args ...Type) error {
	if len(a) != len(args) {
		return fmt.Errorf("invalid length. expected %d, got %d", len(args), len(a))
	}

	var result error

	for i, aa := range a {
		t := args[i]
		if t != aa.Type() {
			result = multierror.Append(result, fmt.Errorf("invalid type at index '%d'. expected %d, found %d", i, t, aa.Type()))
		}
	}

	return result

}

func (a *ArgumentList) Free() {
	for _, aa := range *a {
		aa.Free()
	}
	*a = ArgumentList{}
}

func (a ArgumentList) ToInterfaceSlice() []interface{} {
	var out []interface{}
	for _, i := range a {
		out = append(out, i.Value())
	}
	return out
}

func NewList(in ...interface{}) (ArgumentList, error) {
	out := make([]Argument, len(in))
	var result error
	for i, v := range in {
		a, e := New(v)
		if e != nil {
			result = multierror.Append(result, e)
			out[i] = Undefined()
		} else {
			out[i] = a
		}
	}

	return ArgumentList(out), result
}

func MustList(in ...interface{}) ArgumentList {
	v, e := NewList(in...)
	if e != nil {
		panic(e)
	}
	return v
}
