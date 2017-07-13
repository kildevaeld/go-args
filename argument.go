package args

type argument struct {
	t Type
	v interface{}
}

func (a *argument) Type() Type {
	return a.t
}

func (a *argument) Value() interface{} {
	return a.v
}

func (a *argument) Free() {
	if free, ok := a.v.(Freeable); ok {
		free.Free()
	}
	a.t = NilType
	a.v = nil
	argumentPool.Put(a)
}

func (a *argument) Call(arg Argument) (Argument, error) {
	if a.t == CallType {
		if v, ok := a.v.(func(Argument) (Argument, error)); ok {
			return v(arg)
		} else if v, ok := a.v.(Call); ok {
			return v.Call(arg)
		}

	}
	return nil, nil
}
