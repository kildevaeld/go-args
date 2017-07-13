package args

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
