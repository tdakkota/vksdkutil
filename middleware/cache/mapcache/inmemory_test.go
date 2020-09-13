package mapcache

import (
	"testing"

	"github.com/tdakkota/vksdkutil/v2/middleware/cache/testutil"
)

func TestMap(t *testing.T) {
	t.Run("test-cache", testutil.TestCache(NewMap()))
}
