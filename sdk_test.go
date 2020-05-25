package sdkutil

import (
	"net/http"
	"testing"

	"github.com/SevereCloud/vksdk/api"
	"github.com/stretchr/testify/assert"
)

func TestSDKBuilder(t *testing.T) {
	httpClient := &http.Client{}
	builder := BuildSDK("token").
		WithHTTPClient(httpClient).
		WithRequestLimit(10).
		WithVersion("test-version").
		WithUserAgent("test-user-agent").
		WithMethodURL("http://test.vk.url/").
		WithMiddleware(func(handler Handler) Handler {
			return func(method string, params api.Params) (api.Response, error) {
				return api.Response{
					Response: []byte("1"),
				}, nil
			}
		})
	sdk := builder.Complete()

	assert.Equal(t, sdk, builder.sdk)
	assert.Equal(t, httpClient, sdk.Client)
	assert.Equal(t, "test-version", sdk.Version)
	assert.Equal(t, "test-user-agent", sdk.UserAgent)
	assert.Equal(t, 10, sdk.Limit)
	assert.Equal(t, "token", sdk.AccessToken)
	assert.Equal(t, "http://test.vk.url/", sdk.MethodURL)

	var r int

	err := sdk.RequestUnmarshal("test.call", api.Params{}, &r)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 1, r)
}

func TestPoolCreated(t *testing.T) {
	builder := BuildSDK("token", "token2").
		WithRequestLimitPerToken(2)
	sdk := builder.Complete()

	assert.Equal(t, true, sdk.IsPoolClient)
	assert.Equal(t, 2*2, sdk.Limit)
}
