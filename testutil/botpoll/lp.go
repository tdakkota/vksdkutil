package botpoll

import (
	"strconv"

	"github.com/SevereCloud/vksdk/api"
	"github.com/SevereCloud/vksdk/longpoll-bot"
)

func NewLongPoll(vk *api.VK) (*longpoll.Longpoll, TestLongPoll) {
	server := NewTestLongPoll()

	options := server.Subscribe()
	lp := &longpoll.Longpoll{
		Server: options.Server,
		Key:    options.Key,
		Ts:     strconv.Itoa(options.Ts),
		Wait:   25,
		VK:     vk,
		Client: server.Client(),
	}

	return lp, server
}
