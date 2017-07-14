package args

type ArgumentDescription struct {
	Type  Type        `json:"type"`
	Name  string      `json:"name"`
	Value interface{} `json:"value,omitempty"`
}

type FunctionDescription struct {
	In  []ArgumentDescription `json:"in"`
	Out []ArgumentDescription `json:"out"`
}

func mapToDescription(a ArgumentMap) ArgumentDescription {
	o := ArgumentDescription{
		Type: ArgumentType,
		Name: ArgumentType.String(),
	}
	vv := make(map[string]ArgumentDescription)
	for k, v := range a {
		vv[k] = Describe(v)
	}
	o.Value = vv
	return o
}

func listToDescription(a ArgumentList) ArgumentDescription {
	o := ArgumentDescription{
		Type: ArgumentSliceType,
		Name: ArgumentSliceType.String(),
	}
	vv := make([]ArgumentDescription, a.Len())
	for k, v := range a {
		vv[k] = Describe(v)
	}
	o.Value = vv
	return o
}

func funcToDescription(f *func_argument) FunctionDescription {
	fn := func(t []Type) []ArgumentDescription {
		var o []ArgumentDescription
		for _, tt := range t {
			o = append(o, ArgumentDescription{
				Type: tt,
				Name: tt.String(),
			})
		}
		return o
	}

	return FunctionDescription{
		In:  fn(f.args),
		Out: fn(f.out),
	}

}

func Describe(a Argument) ArgumentDescription {
	switch a.Type() {
	case ArgumentListType:
		return listToDescription(a.Value().(ArgumentList))
	case ArgumentSliceType:
		return listToDescription(ArgumentList(a.Value().([]Argument)))
	case ArgumentMapType:
		return mapToDescription(a.Value().(ArgumentMap))
	default:
		val := a.Value()
		if a.Type() == ErrorType {
			val = a.Value().(error).Error()
		} else if a.Type() == CallType {
			if f, ok := a.Value().(*func_argument); ok {
				val = funcToDescription(f)
			}
		}

		if t := findTypeWithType(a.Type()); t != nil && t.e.Describe != nil {
			if d := t.e.Describe(a.Value()); d != nil {
				return *d
			}
		}

		return ArgumentDescription{
			Type:  a.Type(),
			Name:  a.Type().String(),
			Value: val,
		}
	}
}

func DescriptionToArgument(a ArgumentDescription) (Argument, error) {
	return nil, nil
}
