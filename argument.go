package args

import (
	"gopkg.in/yaml.v2"

	"encoding/json"
)

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

func (a *argument) Valid() bool {
	return a.v != UndefinedType
}

func (a *argument) Free() {
	a.t = UndefinedType
	if a.v != nil {
		if free, ok := a.v.(Freeable); ok {
			free.Free()
		}
	}
	a.v = nil
	argumentPool.Put(a)
}

func (a *argument) Is(t ...Type) bool {
	for _, tt := range t {
		if tt == a.t {
			return true
		}
	}
	return false
}

func (a *argument) Call(arg ArgumentList) (Argument, error) {
	if a.t == CallType {
		if v, ok := a.v.(func(ArgumentList) (Argument, error)); ok {
			return v(arg)
		} else if v, ok := a.v.(Call); ok {
			return v.Call(arg)
		}

	}
	return nil, nil
}

func (a *argument) MarshalJSON() ([]byte, error) {
	if t := findTypeWithType(a.t); t != nil && t.e.MarshalJSON != nil {
		return t.e.MarshalJSON(a.v)
	}
	return json.Marshal(a.v)
}

func (a *argument) MarshalYAML() (interface{}, error) {
	if t := findTypeWithType(a.t); t != nil && t.e.MarshalYAML != nil {
		return t.e.MarshalYAML(a.v)
	}

	return yaml.Marshal(a.v)
}
