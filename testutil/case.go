package testutil

import (
	"fmt"
	"testing"

	"github.com/SevereCloud/vksdk/v2/api"
	sdkutil "github.com/tdakkota/vksdkutil/v2"
)

type TestCase struct {
	Expectations     Expectations
	T                *testing.T
	DefaultResponse  api.Response
	allowNotExpected bool
}

func NewTestCase(t *testing.T) *TestCase {
	return &TestCase{Expectations: Expectations{}, T: t}
}

func (test *TestCase) AllowNotExpected() {
	test.allowNotExpected = true
}

func (test *TestCase) ExpectCall(method string) *Expectation {
	p := NewExpectation(method)
	test.Expectations.Push(p)

	return p
}

// ExpectationsError returns error if not all call expectation were met.
func (test *TestCase) ExpectationsError() error {
	if len(test.Expectations) != 0 {
		return fmt.Errorf("expected %d calls yet", len(test.Expectations))
	}

	return nil
}

// ExpectationsWereMet checks that all call expectation were met.
// Calls T.Error(err).
func (test *TestCase) ExpectationsWereMet() error {
	err := test.ExpectationsError()
	if err != nil {
		test.T.Error(err)
		return err
	}

	return nil
}

var ErrNotExpected = fmt.Errorf("call is not expected")

func (test *TestCase) Handler() sdkutil.Handler {
	return func(method string, params ...api.Params) (api.Response, error) {
		r, ok := test.Expectations.Pop()
		if !ok {
			if test.allowNotExpected {
				return test.DefaultResponse, nil
			}
			return api.Response{}, ErrNotExpected
		}

		matches, response, err := r.Match(method, params...)
		if matches && err != nil {
			return api.Response{}, err
		}

		return response, err
	}
}

func WithSDK(t *testing.T, f func(*testing.T, *api.VK, *TestCase)) {
	t.Helper()

	sdk, testCase := CreateSDK(t)
	defer func() {
		if err := testCase.ExpectationsWereMet(); err != nil {
			t.Error(err)
		}
	}()

	f(t, sdk, testCase)
}

func CreateSDK(t *testing.T) (*api.VK, *TestCase) {
	t.Helper()

	sdk, testCase := api.NewVK(""), NewTestCase(t)
	sdk.Handler = testCase.Handler()

	return sdk, testCase
}
