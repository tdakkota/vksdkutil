package botpoll

import (
	"context"
	"testing"

	"github.com/SevereCloud/vksdk/object"
	"github.com/stretchr/testify/assert"
)

func TestSubscriptions(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		s := newSubscriptions()

		_, ok := s.Get("key")
		assert.False(t, ok)

		key := s.Create()
		_, ok = s.Get(key)
		assert.True(t, ok)

		s.Delete(key)
		_, ok = s.Get(key)
		assert.False(t, ok)
	})

	t.Run("notify", func(t *testing.T) {
		subs := newSubscriptions()

		s, ok := subs.Get(subs.Create())
		assert.True(t, ok)

		s2, ok := subs.Get(subs.Create())
		assert.True(t, ok)

		subs.Notify([]object.GroupEvent{
			{Type: object.EventMessageNew},
		})

		subs.Notify([]object.GroupEvent{
			{Type: object.EventMessageNew},
		})

		assert.Len(t, s.Poll(context.Background()).Updates, 2)
		assert.Len(t, s2.Poll(context.Background()).Updates, 2)
	})
}
