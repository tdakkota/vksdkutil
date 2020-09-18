package testutil

import (
	"fmt"
	"testing"

	"github.com/SevereCloud/vksdk/v2/api"
	sdkutil "github.com/tdakkota/vksdkutil/v2"
)

// TestCase is VK SDK testcase.
type TestCase struct {
	Expectations     Expectations
	T                *testing.T
	DefaultResponse  api.Response
	allowNotExpected bool
}

// NewTestCase creates new TestCase.
func NewTestCase(t *testing.T) *TestCase {
	return &TestCase{Expectations: Expectations{}, T: t}
}

// AllowNotExpected allows not expected calls.
func (test *TestCase) AllowNotExpected() {
	test.allowNotExpected = true
}

// ExpectCall adds call expectation.
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

// ErrNotExpected is expectation error.
var ErrNotExpected = fmt.Errorf("call is not expected")

// Handler is mocking VK SDK handler.
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

// WithSDK calls CreateSDK and passes result to callback.
// Fails if not ExpectationsWereMet after f parameter call.
func WithSDK(t *testing.T, f func(*testing.T, *api.VK, *TestCase)) {
	t.Helper()

	sdk, testCase := CreateSDK(t)
	defer testCase.ExpectationsWereMet()

	f(t, sdk, testCase)
}

// CreateSDK creates new testing VK SDK and TestCase.
func CreateSDK(t *testing.T) (*api.VK, *TestCase) {
	t.Helper()

	sdk, testCase := api.NewVK(""), NewTestCase(t)
	sdk.Handler = testCase.Handler()

	return sdk, testCase
}
