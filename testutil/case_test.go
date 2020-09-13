package testutil

import (
	"testing"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/stretchr/testify/assert"
)

func TestTestCase(t *testing.T) {
	t.Run("all-expected", func(t *testing.T) {
		vksdk, testCase := CreateSDK(t)

		params := api.Params{
			"param": 1,
		}
		testCase.ExpectCall("expect.method").WithParams(params)
		testCase.ExpectCall("expect.method2").Fails(true)

		_, err := vksdk.Request("expect.method", params)
		assert.NoError(t, err)

		_, err = vksdk.Request("expect.method2", params)
		assert.Error(t, err)

		assert.NoError(t, testCase.ExpectationsWereMet())
	})

	t.Run("not-expected", func(t *testing.T) {
		vksdk, testCase := CreateSDK(t)

		_, err := vksdk.Request("notexpect.method", api.Params{})
		assert.Equal(t, ErrNotExpected, err)

		assert.NoError(t, testCase.ExpectationsWereMet())
	})

	t.Run("not-matched", func(t *testing.T) {
		vksdk, testCase := CreateSDK(t)

		params := api.Params{
			"param": 1,
		}
		testCase.ExpectCall("expect.method").WithParams(params)

		_, err := vksdk.Request("expect.method", api.Params{})
		assert.Error(t, err)

		assert.NoError(t, testCase.ExpectationsWereMet())
	})

	t.Run("not-all-called", func(t *testing.T) {
		_, testCase := CreateSDK(t)

		params := api.Params{
			"param": 1,
		}
		testCase.ExpectCall("expect.method").WithParams(params)

		assert.Error(t, testCase.ExpectationsWereMet())
	})

	t.Run("with-sdk", func(t *testing.T) {
		WithSDK(t, func(t *testing.T, vk *api.VK, testCase *TestCase) {})

		assert.False(t, t.Failed())
	})
}
