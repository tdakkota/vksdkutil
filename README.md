# VK SDK Utilities 
[![Build Status](https://travis-ci.com/tdakkota/vksdkutil.svg?branch=master)](https://travis-ci.com/tdakkota/vksdkutil)
[![Documentation](https://godoc.org/github.com/tdakkota/vksdkutil?status.svg)](https://pkg.go.dev/github.com/tdakkota/vksdkutil?tab=subdirectories)
[![codecov](https://codecov.io/gh/tdakkota/vksdkutil/branch/master/graph/badge.svg)](https://codecov.io/gh/tdakkota/vksdkutil)
[![license](https://img.shields.io/github/license/tdakkota/vksdkutil.svg?maxAge=2592000)](https://github.com/tdakkota/vksdkutil/blob/master/LICENSE)

Some useful things for [vksdk](https://github.com/SevereCloud/vksdk)

## Features

- Handler middlewares(logging, retrying, etc.)
- `testutil` package for `api.VK` mocking

### Middleware example

```go
package main

import (
    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
    sdkutil "github.com/tdakkota/vksdkutil"
    zlog "github.com/tdakkota/vksdkutil/middleware/zerolog"
)

func main() {
    vk := sdkutil.BuildSDK("token").WithMiddleware(zlog.LoggingMiddleware(
         log.With().Str("type", "vksdk").Logger().Level(zerolog.DebugLevel),
    ))
    // ...
}
```

### Testing example
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