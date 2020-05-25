package testutil

type Expectations []*Expectation

func (stack *Expectations) Pop() (result *Expectation, ok bool) {
	if len(*stack) == 0 {
		return nil, false
	}
	result, *stack = (*stack)[0], (*stack)[1:]
	return result, true
}

func (stack *Expectations) Push(element *Expectation) {
	*stack = append(*stack, element)
}
