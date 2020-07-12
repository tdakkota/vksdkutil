# Retry example
A middleware to retry failed API requests 3 times with a 1 second timeout.
```go
package main

import (
    "time"

    sdkutil "github.com/tdakkota/vksdkutil"
    "github.com/tdakkota/vksdkutil/middleware"
)

func main() {
    n, timeout := 3, time.Second
    vk := sdkutil.BuildSDK("token").
        WithMiddleware(middleware.Retry(n, timeout)).
        Complete()
    // ...
}
```