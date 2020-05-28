package sdkutil

import (
	"net/http"

	"github.com/SevereCloud/vksdk/api"
)

// SDKBuilder represents *api.VK builder.
type SDKBuilder struct {
	tokens int
	sdk    *api.VK
}

// WithMiddleware adds Handler middleware.
func (builder SDKBuilder) WithMiddleware(f Middleware) SDKBuilder {
	PatchHandler(builder.sdk, f)

	return builder
}

// WithHTTPClient sets HTTP client.
func (builder SDKBuilder) WithHTTPClient(client *http.Client) SDKBuilder {
	builder.sdk.Client = client
	return builder
}

// WithRequestLimit sets request limit by second.
func (builder SDKBuilder) WithRequestLimit(limit int) SDKBuilder {
	builder.sdk.Limit = limit
	return builder
}

// WithRequestLimitPerToken sets request limit by second per token.
func (builder SDKBuilder) WithRequestLimitPerToken(limit int) SDKBuilder {
	builder.sdk.Limit = limit * builder.tokens
	return builder
}

// WithMethodURL sets API endpoint URL.
func (builder SDKBuilder) WithMethodURL(url string) SDKBuilder {
	builder.sdk.MethodURL = url
	return builder
}

// WithUserAgent sets User-Agent header.
func (builder SDKBuilder) WithUserAgent(agent string) SDKBuilder {
	builder.sdk.UserAgent = agent
	return builder
}

// WithVersion sets API version.
func (builder SDKBuilder) WithVersion(v string) SDKBuilder {
	builder.sdk.Version = v
	return builder
}

// Complete returns built API client.
func (builder SDKBuilder) Complete() *api.VK {
	return builder.sdk
}

// BuildSDK creates new SDKBuilder.
func BuildSDK(token string, tokens ...string) SDKBuilder {
	var vk *api.VK

	tokensLen := len(tokens)
	if tokensLen == 0 {
		vk = api.NewVK(token)
	} else {
		tokens = append([]string{token}, tokens...)
		vk = api.NewVKWithPool(tokens...)
	}

	return SDKBuilder{
		tokens: tokensLen + 1,
		sdk:    vk,
	}
}
