package middleware

import (
	"testing"

	"github.com/SevereCloud/vksdk/api"
	"github.com/tdakkota/vksdkutil"
)

func testDeleteToken(params []string, handler sdkutil.Handler) func(t *testing.T) {
	return func(t *testing.T) {
		m := DeleteParams(params)

		_, err := m(handler)("", api.Params{"access_token": "", "v": ""})
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestDeleteToken(t *testing.T) {
	t.Run("no-one", testDeleteToken(
		[]string{},
		func(method string, params api.Params) (api.Response, error) {
			if _, ok := params["access_token"]; !ok {
				t.Fatal("expected token")
			}

			if _, ok := params["v"]; !ok {
				t.Fatal("expected version")
			}

			return api.Response{}, nil
		}),
	)

	t.Run("one", testDeleteToken(
		[]string{"access_token"},
		func(method string, params api.Params) (api.Response, error) {
			if _, ok := params["access_token"]; ok {
				t.Fatal("expected no token")
			}

			if _, ok := params["v"]; !ok {
				t.Fatal("expected version")
			}

			return api.Response{}, nil
		}),
	)

	t.Run("multiple", testDeleteToken(
		[]string{"access_token", "v"},
		func(method string, params api.Params) (api.Response, error) {
			if _, ok := params["access_token"]; ok {
				t.Fatal("expected no token")
			}

			if _, ok := params["v"]; ok {
				t.Fatal("expected no version")
			}

			return api.Response{}, nil
		}),
	)
}
