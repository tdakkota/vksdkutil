package middleware

import (
	"github.com/tdakkota/vksdkutil/v2/middleware/paramsutil"
	"testing"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/tdakkota/vksdkutil/v2"
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
		func(method string, params ...api.Params) (api.Response, error) {
			if _, ok := paramsutil.Find("access_token", params...); !ok {
				t.Fatal("expected token")
			}

			if _, ok := paramsutil.Find("v", params...); !ok {
				t.Fatal("expected version")
			}

			return api.Response{}, nil
		}),
	)

	t.Run("one", testDeleteToken(
		[]string{"access_token"},
		func(method string, params ...api.Params) (api.Response, error) {
			if _, ok := paramsutil.Find("access_token", params...); ok {
				t.Fatal("expected no token")
			}

			if _, ok := paramsutil.Find("v", params...); !ok {
				t.Fatal("expected version")
			}

			return api.Response{}, nil
		}),
	)

	t.Run("multiple", testDeleteToken(
		[]string{"access_token", "v"},
		func(method string, params ...api.Params) (api.Response, error) {
			if _, ok := paramsutil.Find("access_token", params...); ok {
				t.Fatal("expected no token")
			}

			if _, ok := paramsutil.Find("v", params...); ok {
				t.Fatal("expected no version")
			}

			return api.Response{}, nil
		}),
	)
}
