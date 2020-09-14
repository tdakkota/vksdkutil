package paramsutil

import (
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/stretchr/testify/require"
	"testing"
)

func runTest(cb func(a *require.Assertions)) func(*testing.T) {
	return func(t *testing.T) {
		cb(require.New(t))
	}
}

func TestAppend(t *testing.T) {
	t.Run("empty", runTest(func(a *require.Assertions) {
		var params []api.Params
		params = Append("k", "v", params...)

		a.Len(params, 1)
		a.Equal(params[0]["k"], "v")
	}))

	t.Run("not-empty", runTest(func(a *require.Assertions) {
		params := []api.Params{
			{"key": "value"},
		}
		params = Append("k", "v", params...)

		a.Len(params, 1)
		a.Equal(params[0]["k"], "v")
		a.Equal(params[0]["key"], "value")
	}))
}

func TestDelete(t *testing.T) {
	t.Run("empty", runTest(func(a *require.Assertions) {
		var params []api.Params
		a.False(Delete("key", params...))
	}))

	t.Run("not-empty", runTest(func(a *require.Assertions) {
		params := []api.Params{
			{"key": "value"},
		}
		a.True(Delete("key", params...))
	}))
}

func TestFind(t *testing.T) {
	t.Run("empty", runTest(func(a *require.Assertions) {
		var params []api.Params
		_, ok := Find("key", params...)

		a.False(ok)
	}))

	t.Run("not-empty", runTest(func(a *require.Assertions) {
		params := []api.Params{
			{"key": "value"},
		}
		v, ok := Find("key", params...)

		a.True(ok)
		a.Equal(v, params[0]["key"])
	}))
}
