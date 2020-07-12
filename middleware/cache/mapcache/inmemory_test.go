package mapcache

import (
	"testing"

	"github.com/tdakkota/vksdkutil/middleware/cache/testutil"
)

func TestMap(t *testing.T) {
	t.Run("test-cache", testutil.TestCache(NewMap()))
}
