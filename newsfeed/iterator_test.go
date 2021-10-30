package newsfeed

import (
	"context"
	"os"
	"testing"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/stretchr/testify/require"
)

func TestIterator_Next(t *testing.T) {
	token := os.Getenv("E2E_VK_TOKEN")
	if token == "" {
		t.Skip("Set E2E_VK_TOKEN environment variable to run test")
	}

	a := require.New(t)
	ctx := context.Background()

	vk := api.NewVK(token)
	iter := NewIterator(vk, func(b *params.NewsfeedGetBuilder) {
		b.TestMode(true)
		b.Count(3)
		b.Filters([]string{"post"})
	})

	for i := 0; i < 10; i++ {
		if !iter.Next(ctx) {
			break
		}
		v := iter.Value()
		t.Logf("%+v", v.NewsfeedNewsfeedItem)
	}
	a.NoError(iter.Err())
}
