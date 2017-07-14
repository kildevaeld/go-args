package args

import "encoding/json"

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
		vv[k] = ToDescription(v)
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
		vv[k] = ToDescription(v)
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

func ToDescription(a Argument) ArgumentDescription {
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

		return ArgumentDescription{
			Type:  a.Type(),
			Name:  a.Type().String(),
			Value: val,
		}
	}
}

func ToJSON(a Argument) ([]byte, error) {
	return json.Marshal(ToDescription(a))
}

func ToJSONIndented(a Argument) ([]byte, error) {
	return json.MarshalIndent(ToDescription(a), "", "  ")
}
