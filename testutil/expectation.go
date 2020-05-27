package testutil

import (
	"encoding/json"
	"fmt"

	"github.com/SevereCloud/vksdk/api"
	"github.com/SevereCloud/vksdk/api/errors"
	"github.com/SevereCloud/vksdk/object"
)

type Expectation struct {
	Method        string
	Params        api.Params
	Response      api.Response
	ErrorResponse bool
	ErrorMessage  string
}

func NewExpectation(method string) *Expectation {
	return &Expectation{Method: method}
}

func (e *Expectation) WithParams(params api.Params) *Expectation {
	e.Params = params
	return e
}

func (e *Expectation) WithParamsF(f func() api.Params) *Expectation {
	return e.WithParams(f())
}

func (e *Expectation) Returns(response api.Response) *Expectation {
	e.Response = response
	return e
}

func (e *Expectation) ReturnsF(f func() api.Response) *Expectation {
	return e.Returns(f())
}

func (e *Expectation) ReturnsJSON(v interface{}) *Expectation {
	e.Response.Response, _ = json.Marshal(v)
	return e
}

func (e *Expectation) ReturnsBytes(data []byte) *Expectation {
	e.Response.Response = data
	return e
}

func (e *Expectation) WithError(message string) *Expectation {
	e.ErrorResponse = true
	e.ErrorMessage = message
	return e
}

func (e *Expectation) Fails(fails bool) *Expectation {
	e.ErrorResponse = fails
	return e
}

func paramsToRequestParams(params api.Params) []object.BaseRequestParam {
	r := make([]object.BaseRequestParam, 0, len(params))

	for name, param := range params {
		r = append(r, object.BaseRequestParam{
			Key:   name,
			Value: api.FmtValue(param, 0),
		})
	}

	return r
}

func (e Expectation) generateError(params api.Params) object.Error {
	text := e.ErrorMessage
	if text == "" {
		text = fmt.Sprintf("Generated test error for %s method", e.Method)
	}

	return object.Error{
		Code:          9999,
		Message:       text,
		Text:          text,
		RequestParams: paramsToRequestParams(params),
	}
}

func matchParams(got, expected api.Params) bool {
	for name, param := range expected {
		v, ok := got[name]
		if !ok || api.FmtValue(param, 0) != api.FmtValue(v, 0) {
			return false
		}
	}

	return true
}

func (e Expectation) matchParams(params api.Params) bool {
	return matchParams(params, e.Params)
}

func (e Expectation) Match(method string, params api.Params) (bool, api.Response, error) {
	if e.Method != method {
		return false, api.Response{}, nil
	}

	if !e.matchParams(params) {
		return false, api.Response{}, nil
	}

	if e.ErrorResponse {
		err := e.generateError(params)
		return true, api.Response{
			Error: err,
		}, errors.New(err)
	}

	if e.Response.Response == nil {
		e.Response.Response = []byte(`1`)
	}

	return true, e.Response, nil
}
