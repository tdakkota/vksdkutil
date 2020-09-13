package botpoll

import (
	"strconv"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
)

func NewLongPoll(vk *api.VK) (*longpoll.LongPoll, TestLongPoll) {
	server := NewTestLongPoll()

	options := server.Subscribe()
	vk.Client = server.Client()
	vk.MethodURL = server.URL()

	lp := &longpoll.LongPoll{
		Server: options.Server,
		Key:    options.Key,
		Ts:     strconv.Itoa(options.Ts),
		Wait:   25,
		VK:     vk,
		Client: vk.Client,
	}

	return lp, server
}
