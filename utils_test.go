package args

import (
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
	if err := Register(&T{}, Test); err != nil {
		t.Fatal(err)
	}

	assert.Nil(t, Register(&R{}, Test2))

	assert.NotNil(t, Register(&T{}, Test))

	arg, err := New(&T{
		Rapper: "Test",
	})

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, arg.Type(), Test)

}
