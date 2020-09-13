package botpoll

import (
	"context"
	"github.com/SevereCloud/vksdk/v2/events"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubscription(t *testing.T) {
	s := newSubscription()
	defer s.Close()

	s.Notify([]events.GroupEvent{
		{Type: events.EventMessageNew},
	})

	s.Notify([]events.GroupEvent{
		{Type: events.EventMessageNew},
	})

	ctxt, cancel := context.WithCancel(context.Background())
	defer cancel()

	r := s.Poll(ctxt)
	assert.Len(t, r.Updates, 2)
}
