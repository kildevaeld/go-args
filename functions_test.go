package args

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFunction(t *testing.T) {

	fn := func(name string, m ArgumentMap) error {
		return errors.New("error")
	}

	if call, err := FunctionToCall(fn); err != nil {
		t.Fatal(err)
	} else {
		_, err := call.Call(ArgumentList{Must("World")})
		assert.EqualError(t, err, "invalid call count")
		a, aerr := call.Call(MustList("World", MustMap(Map{"test": "mig"})))
		assert.Nil(t, aerr)

		assert.Equal(t, a.Type(), ArgumentSliceType)
		first := a.Value().([]Argument)[0]
		assert.Equal(t, first.Type(), ErrorType)
		assert.EqualError(t, first.Value().(error), "error")

		bs, _ := ToJSONIndented(Must(call))
		fmt.Printf("%s\n", bs)
	}

	if _, err := FunctionToCall(func(m Map) {}); err != nil {
		t.Fatal(err)
	}

}
