# Cache middleware
`Storage` interface represents a cache storage
```go
type Storage interface {
	Save(key Key, value api.Response) error
	Load(key Key) (api.Response, error)
}
```

### Usage:
```go
package main
import (

"github.com/SevereCloud/vksdk/api"
sdkutil "github.com/tdakkota/vksdkutil"
"github.com/tdakkota/vksdkutil/middleware/cache"
"github.com/tdakkota/vksdkutil/middleware/cache/mapcache"
"strings"
)

// cache all requests with method which starts with "get"
func cacheable(method string, param api.Params) bool {
    return strings.Contains(method, "get")  
}

func main() {
    storage := mapcache.NewMap() // create in-memory cache based on Go map

    vk := sdkutil.BuildSDK("token").
        WithMiddleware(cache.Create(storage, cacheable)). 
        Complete()
    // ...
}
```