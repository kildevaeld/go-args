package args

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const Test Type = NilType + 2
const Test2 Type = NilType + 4

type T struct {
	Rapper string
}

type R struct {
}

func TestRegister(t *testing.T) {
	if err := Register(&T{}, Extension{
		Name: "Test",
		Type: Test,
		Describe: func(a interface{}) *ArgumentDescription {
			v := a.(*T)
			return &ArgumentDescription{
				Name: "Test",
				Type: Test,
				Value: ArgumentDescription{
					Name: ArgumentMapType.String(),
					Type: ArgumentMapType,
					Value: Map{
						"Rapper": Must(v.Rapper).Describe(),
					},
				},
			}
		},
	}); err != nil {
		t.Fatal(err)
	}

	assert.Nil(t, Register(&R{}, Extension{
		Name: "Test2",
		Type: Test2,
	}))

	assert.NotNil(t, Register(&T{}, Extension{
		Name: "Test",
		Type: Test,
	}))

	arg, err := New(&T{
		Rapper: "Test",
	})

	if err != nil {
		t.Fatal(err)
	}

	b, _ := json.MarshalIndent(arg.Describe(), "", "  ")
	fmt.Printf("%s\n", b)
	assert.Equal(t, arg.Type(), Test)

}
