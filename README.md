# VK SDK Utilities 
[![CI](https://github.com/tdakkota/vksdkutil/workflows/master/badge.svg)](https://github.com/tdakkota/vksdkutil/actions)
[![Documentation](https://godoc.org/github.com/tdakkota/vksdkutil?status.svg)](https://pkg.go.dev/github.com/tdakkota/vksdkutil?tab=subdirectories)
[![codecov](https://codecov.io/gh/tdakkota/vksdkutil/branch/master/graph/badge.svg)](https://codecov.io/gh/tdakkota/vksdkutil)
[![license](https://img.shields.io/github/license/tdakkota/vksdkutil.svg?maxAge=2592000)](https://github.com/tdakkota/vksdkutil/blob/master/LICENSE)

Some useful things for [vksdk](https://github.com/SevereCloud/vksdk)

## Features

- Handler middlewares
    - [Logging](https://github.com/tdakkota/vksdkutil/tree/master/middleware/log)
    - [Retrying](https://github.com/tdakkota/vksdkutil/blob/master/middleware/README.md)
    - [Caching](https://github.com/tdakkota/vksdkutil/blob/master/middleware/cache/README.md)
- `testutil` package for `api.VK` mocking

## Middleware example

```go
package main

import (
  "github.com/SevereCloud/vksdk/v2/api"
  sdkutil "github.com/tdakkota/vksdkutil/v3"
  "github.com/tdakkota/vksdkutil/v3/middleware/zapvk"
  "go.uber.org/zap"
  "go.uber.org/zap/zapcore"
)

func main() {
  vk := sdkutil.BuildSDK("token").WithMiddleware(
    zapvk.Log(zap.L(), zapcore.DebugLevel, true),
  ).Complete()
  // ...
}
```

## Testing example
You have a file

```go
package mypackage

import (
    "github.com/SevereCloud/vksdk/api"
)

func MarkAsRead(sdk *api.VK, peerID int) (int, error) {
    builder := params.NewMessagesMarkAsReadBuilder()
    builder.PeerID(peerID)
    
    return sdk.MessagesMarkAsRead(builder.Params)
}
```

So, with `testutil` you can test it
```go
package mypackage

import (
    "testing"

    "github.com/tdakkota/vksdkutil/testutil"
)

func TestMarkAsRead(t *testing.T) {
	sdk, expect := testutil.CreateSDK(t)

	peerID, count := 10, 2
	expect.ExpectCall("messages.markAsRead").WithParams(api.Params{
		"peer_id": peerID,
	}).ReturnsJSON(count)

	read, err := MarkAsRead(sdk, peerID)
	if err != nil {
		t.Fatal(err)
	}

	if count != read {
		t.Errorf("expected %d, got %d", count, read)
	}
}
```