package testutil

import (
	"github.com/SevereCloud/vksdk/api"
	"github.com/stretchr/testify/assert"
	"testing"
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
	})

	t.Run("not-expected", func(t *testing.T) {
		vksdk, _ := CreateSDK(t)

		_, err := vksdk.Request("notexpect.method", api.Params{})
		assert.Equal(t, ErrNotExpected, err)
	})

	t.Run("not-matched", func(t *testing.T) {
		vksdk, testCase := CreateSDK(t)

		params := api.Params{
			"param": 1,
		}
		testCase.ExpectCall("expect.method").WithParams(params)

		_, err := vksdk.Request("expect.method", api.Params{})
		assert.Equal(t, ErrNotMatched, err)
	})
}
