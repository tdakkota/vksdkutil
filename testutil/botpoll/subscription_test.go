package botpoll

import (
	"context"
	"testing"

	"github.com/SevereCloud/vksdk/object"
	"github.com/stretchr/testify/assert"
)

func TestSubscription(t *testing.T) {
	s := newSubscription()
	defer s.Close()

	s.Notify([]object.GroupEvent{
		{Type: object.EventMessageNew},
	})

	s.Notify([]object.GroupEvent{
		{Type: object.EventMessageNew},
	})

	ctxt, cancel := context.WithCancel(context.Background())
	defer cancel()

	r := s.Poll(ctxt)
	assert.Len(t, r.Updates, 2)
}
