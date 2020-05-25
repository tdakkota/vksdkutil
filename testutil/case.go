package testutil

import (
	"fmt"
	"github.com/SevereCloud/vksdk/api"
	sdkutil "github.com/tdakkota/vksdkutil"
	"testing"
)

type TestCase struct {
	Expectations Expectations
}

func NewTestCase() *TestCase {
	return &TestCase{Expectations: Expectations{}}
}

func (test *TestCase) ExpectCall(method string) *Expectation {
	p := NewExpectation(method)
	test.Expectations.Push(p)
	return p
}

var ErrNotExpected = fmt.Errorf("call is not expected")
var ErrNotMatched = fmt.Errorf("call is not matched")

func (test *TestCase) Handler() sdkutil.Handler {
	return func(method string, params api.Params) (api.Response, error) {
		r, ok := test.Expectations.Pop()
		if !ok {
			return api.Response{}, ErrNotExpected
		}

		matches, response, err := r.Match(method, params)
		if !matches {
			return api.Response{}, ErrNotMatched
		}

		return response, err
	}
}

func CreateSDK(t *testing.T) (*api.VK, *TestCase) {
	testCase := NewTestCase()
	sdk := api.NewVK("")
	sdk.Handler = testCase.Handler()
	return sdk, testCase
}
