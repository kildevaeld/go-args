package args

import (
	"fmt"
	"reflect"

	multierror "github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

type func_argument struct {
	args []Type
	out  []Type
	v    reflect.Value
}

func checkArguments(a []Type, v []Argument) error {
	if len(v) != len(a) {
		return fmt.Errorf("invalid counts %d ~= %d", len(v), len(a))
	}
	var err error
	for i, ma := range a {
		ha := v[i]
		if !ha.Is(ma) {
			err = multierror.Append(err, fmt.Errorf("parameter at #%d want '%s' got '%s'", i, ma, ha.Type()))
		}
	}
	return err
}

func (f *func_argument) Call(a ArgumentList) (Argument, error) {

	if a.Len() != len(f.args) {
		return nil, errors.New("invalid call count")
	}

	if err := checkArguments(f.args, a); err != nil {
		return nil, errors.Wrap(err, "in parameters")
	}

	ins := make([]reflect.Value, len(f.args))
	for i, v := range a {
		ins[i] = reflect.ValueOf(v.Value())
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()

	outs := f.v.Call(ins)
	result := make([]Argument, len(outs))
	for i, o := range outs {
		result[i] = Must(o.Interface())

	}
	/*fmt.Printf("%#v\n", result[0])
	if err := checkArguments(f.out, result); err != nil {
		return nil, errors.Wrap(err, "output parameters")
	}*/

	return Must(result), nil
}

func FunctionToCall(i interface{}) (Call, error) {

	v := reflect.TypeOf(i)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Func {
		return nil, errors.New("argument is not a function")
	}

	nins := v.NumIn()
	nouts := v.NumOut()

	var args []Type
	for i := 0; i < nins; i++ {
		m := v.In(i)
		if t := TypeFromReflectType(m); t != UndefinedType {
			args = append(args, t)
			continue
		}
		return nil, fmt.Errorf("invalid in-type: %s", m.Kind())
	}
	out := make([]Type, nouts)
	for i := 0; i < nouts; i++ {
		m := v.Out(i)
		if t := TypeFromReflectType(m); t != UndefinedType {
			out[i] = t
			continue
		}
		return nil, fmt.Errorf("invalid out-type: %s", m.Kind())
	}

	return &func_argument{args, out, reflect.ValueOf(i)}, nil
}
