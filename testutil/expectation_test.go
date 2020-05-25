package testutil

import (
	"encoding/json"
	"testing"

	"github.com/SevereCloud/vksdk/api"
	"github.com/SevereCloud/vksdk/object"
	"github.com/stretchr/testify/assert"
)

func createExpectation() (string, api.Params, *Expectation) {
	method := "test.method"
	params := api.Params{
		"key":  1,
		"key2": "value",
	}
	e := NewExpectation(method).
		WithParams(params)
	return method, params, e
}

func TestExpectation(t *testing.T) {
	t.Run("match", func(t *testing.T) {
		method, params, e := createExpectation()

		match, response, err := e.Match(method, params)
		if err != nil {
			t.Fatal(err)
		}

		assert.True(t, match)
		assert.Equal(t, api.Response{Response: []byte(`1`)}, response)
	})

	t.Run("method-not-match", func(t *testing.T) {
		_, params, e := createExpectation()

		match, _, err := e.Match("not.expected", params)
		if err != nil {
			t.Fatal(err)
		}

		assert.False(t, match)
	})

	t.Run("params-not-match", func(t *testing.T) {
		method, _, e := createExpectation()

		match, _, err := e.Match(method, api.Params{
			"key": 1,
		})
		if err != nil {
			t.Fatal(err)
		}

		assert.False(t, match)
	})
}

func TestWithResponse(t *testing.T) {
	t.Run("with-response", func(t *testing.T) {
		method, params, e := createExpectation()

		returns := api.Response{
			Response: []byte(`2`),
		}
		e.Returns(returns)

		_, response, err := e.Match(method, params)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, returns, response)
	})

	t.Run("with-json", func(t *testing.T) {
		method, params, e := createExpectation()

		returns := object.MessagesMessage{
			PeerID: 10,
		}
		e.ReturnsJSON(returns)

		_, response, err := e.Match(method, params)
		if err != nil {
			t.Fatal(err)
		}

		marshaled, err := json.Marshal(returns)
		if err != nil {
			t.Fatal(err)
		}
		assert.JSONEq(t, string(marshaled), string(response.Response))
	})

	t.Run("with-bytes", func(t *testing.T) {
		method, params, e := createExpectation()

		returns := []byte(`10`)
		e.ReturnsBytes(returns)

		_, response, err := e.Match(method, params)
		if err != nil {
			t.Fatal(err)
		}

		assert.JSONEq(t, string(returns), string(response.Response))
	})
}

func TestWithError(t *testing.T) {
	t.Run("with-error", func(t *testing.T) {
		method, params, e := createExpectation()
		e.WithError(`test error`)

		_, response, err := e.Match(method, params)

		assert.Error(t, err)
		assert.Equal(t, `test error`, response.Error.Text)
		assert.Equal(t, `test error`, response.Error.Message)
	})

	t.Run("fails", func(t *testing.T) {
		method, params, e := createExpectation()
		e.Fails(true)

		_, response, err := e.Match(method, params)

		assert.Error(t, err)
		assert.NotEmpty(t, response.Error.Text)
	})
}
