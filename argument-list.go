package args

import multierror "github.com/hashicorp/go-multierror"

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
