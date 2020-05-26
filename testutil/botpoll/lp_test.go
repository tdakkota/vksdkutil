package botpoll

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/SevereCloud/vksdk/object"
	"github.com/stretchr/testify/assert"
	"github.com/tdakkota/vksdkutil/testutil"
)

func TestTestLongPoll(t *testing.T) {
	vk, _ := testutil.CreateSDK(t)

	lp, server := NewLongPoll(vk)
	defer server.Close()

	messages := []string{"test message", "second test message"}
	messageCounter := int64(len(messages))

	lp.MessageNew(func(msg object.MessageNewObject, i int) {
		n := len(messages) - int(atomic.LoadInt64(&messageCounter))

		assert.Equal(t, 1, msg.Message.PeerID)
		assert.Equal(t, messages[n], msg.Message.Text)

		atomic.AddInt64(&messageCounter, -1)
	})

	go func() {
		err := lp.Run()
		if err != nil {
			t.Error(err)
		}
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

	time.Sleep(1 * time.Second)
	lp.Shutdown()

	assert.Zero(t, atomic.LoadInt64(&messageCounter))
}
