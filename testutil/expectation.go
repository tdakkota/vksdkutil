package testutil

import (
	"encoding/json"
	"fmt"
	"github.com/tdakkota/vksdkutil/v2/middleware/paramsutil"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/object"
)

// Expectation is VK API call expectation.
type Expectation struct {
	Method        string
	Params        api.Params
	Response      api.Response
	ErrorResponse bool
	ErrorMessage  string
}

// NewExpectation creates new Expectation.
func NewExpectation(method string) *Expectation {
	return &Expectation{Method: method}
}

// WithParams sets expected params.
// If API call contains unexpected params, it wouldn't fail.
func (e *Expectation) WithParams(params api.Params) *Expectation {
	e.Params = params
	return e
}

// WithParamsF sets expected params using function.
// If API call contains unexpected params, it wouldn't fail.
func (e *Expectation) WithParamsF(f func() api.Params) *Expectation {
	return e.WithParams(f())
}

// Returns sets API call result.
func (e *Expectation) Returns(response api.Response) *Expectation {
	e.Response = response
	return e
}

// ReturnsF sets API call result using function.
func (e *Expectation) ReturnsF(f func() api.Response) *Expectation {
	return e.Returns(f())
}

// ReturnsJSON sets API call result.
func (e *Expectation) ReturnsJSON(v interface{}) *Expectation {
	e.Response.Response, _ = json.Marshal(v)
	return e
}

// ReturnsJSONF sets API call result using function.
func (e *Expectation) ReturnsJSONF(f func() interface{}) *Expectation {
	return e.ReturnsJSON(f())
}

// ReturnsBytes sets API call result.
func (e *Expectation) ReturnsBytes(data []byte) *Expectation {
	e.Response.Response = data
	return e
}

// ReturnsBytesF sets API call result using function.
func (e *Expectation) ReturnsBytesF(f func() []byte) *Expectation {
	return e.ReturnsBytes(f())
}

// WithError sets API call error.
// Calls Fails(true).
func (e *Expectation) WithError(message string) *Expectation {
	e.Fails(true)
	e.ErrorMessage = message
	return e
}

// Fails marks API result as failed.
func (e *Expectation) Fails(fails bool) *Expectation {
	e.ErrorResponse = fails
	return e
}

func paramsToRequestParams(params ...api.Params) []object.BaseRequestParam {
	r := make([]object.BaseRequestParam, 0, len(params))

	for _, set := range params {
		for name, param := range set {
			r = append(r, object.BaseRequestParam{
				Key:   name,
				Value: api.FmtValue(param, 0),
			})
		}
	}

	return r
}

func (e Expectation) generateError(params ...api.Params) api.Error {
	text := e.ErrorMessage
	if text == "" {
		text = fmt.Sprintf("Generated test error for %s method", e.Method)
	}

	return api.Error{
		Code:          9999,
		Message:       text,
		Text:          text,
		RequestParams: paramsToRequestParams(params...),
	}
}

func matchParams(expected api.Params, got []api.Params) error {
	for name, param := range expected {
		v, ok := paramsutil.Find(name, got...)
		if !ok {
			return fmt.Errorf("expected param %s", name)
		}

		gotParam, expectedParam := api.FmtValue(v, 0), api.FmtValue(param, 0)
		if expectedParam != gotParam {
			return fmt.Errorf("expected param %s=%s, got %s=%s", name, expectedParam, name, gotParam)
		}
	}

	return nil
}

func (e Expectation) matchParams(params ...api.Params) error {
	return matchParams(e.Params, params)
}

// Match compares API call.
func (e Expectation) Match(method string, params ...api.Params) (bool, api.Response, error) {
	if e.Method != method {
		return false, api.Response{}, fmt.Errorf("expected method %s, got %s", e.Method, method)
	}

	if err := e.matchParams(params...); err != nil {
		return false, api.Response{}, err
	}

	if e.ErrorResponse {
		err := e.generateError(params...)
		return true, api.Response{
			Error: err,
		}, err
	}

	if e.Response.Response == nil {
		e.Response.Response = []byte(`1`)
	}

	return true, e.Response, nil
}
