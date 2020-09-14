package token

import (
	"errors"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/stretchr/testify/require"
	"testing"
)

func getter(token string, err error) Getter {
	return GetterFunc(func(string, ...api.Params) (string, error) {
		return token, err
	})
}

func testCreate(force bool, expectedToken string, params api.Params) func(t *testing.T) {
	m := Create(force, getter("token", nil))

	return func(t *testing.T) {
		handler := m(func(method string, params ...api.Params) (api.Response, error) {
			require.Equal(t, expectedToken, params[0]["access_token"])
			return api.Response{}, nil
		})

		_, err := handler("", params)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestCreate(t *testing.T) {
	t.Run("get", testCreate(false, "token", api.Params{}))

	t.Run("existing", testCreate(false, "notoken", api.Params{
		"access_token": "notoken",
	}))

	t.Run("force", testCreate(true, "token", api.Params{
		"access_token": "notoken",
	}))

	t.Run("err", func(t *testing.T) {
		e := errors.New("test-error")
		m := Create(false, getter("", e))
		handler := m(func(method string, params ...api.Params) (api.Response, error) {
			return api.Response{}, nil
		})

		_, err := handler("", api.Params{})
		require.Error(t, err)
	})
}
