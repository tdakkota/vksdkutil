package botpoll

import (
	"context"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/object"
	"github.com/stretchr/testify/assert"
	"github.com/tdakkota/vksdkutil/v2/testutil"
	"sync/atomic"
	"testing"
)

func TestTestLongPoll(t *testing.T) {
	vk, e := testutil.CreateSDK(t)
	e.AllowNotExpected()
	e.DefaultResponse = api.Response{
		Response: []byte(`1`),
	}

	lp, server := NewLongPoll(vk)
	defer server.Close()

	messages := []string{"test message", "second test message"}
	messageCounter := int64(len(messages))

	wait := make(chan struct{}, 1)
	lp.MessageNew(func(ctx context.Context, msg events.MessageNewObject) {
		n := len(messages) - int(atomic.LoadInt64(&messageCounter))

		assert.Equal(t, 1, msg.Message.PeerID)
		assert.Equal(t, messages[n], msg.Message.Text)

		atomic.AddInt64(&messageCounter, -1)

		counter := atomic.LoadInt64(&messageCounter)
		if counter == 0 {
			wait <- struct{}{}
		}
	})

	go func() {
		_ = lp.Run()
	}()

	for _, message := range messages {
		err := server.SendMessage(object.MessagesMessage{
			PeerID: 1,
			Text:   message,
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	<-wait
	lp.Shutdown()

	assert.Zero(t, atomic.LoadInt64(&messageCounter))
}
