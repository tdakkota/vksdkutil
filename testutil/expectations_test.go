package testutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpectations_Pop(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		stack := Expectations{}

		e, ok := stack.Pop()
		assert.False(t, ok)
		assert.Nil(t, e)
	})

	t.Run("non-empty", func(t *testing.T) {
		e := NewExpectation("test.method")
		stack := Expectations{e}

		e2, ok := stack.Pop()
		assert.True(t, ok)
		assert.Equal(t, e, e2)
	})
}

func TestExpectations_Push(t *testing.T) {
	e := NewExpectation("test.method")
	stack := Expectations{}

	stack.Push(e)
	assert.Equal(t, e, stack[0])
}
