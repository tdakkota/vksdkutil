package botpoll

import (
	"context"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/SevereCloud/vksdk/object"
)

type subscription struct {
	notify *int64

	events []object.GroupEvent
	lock   sync.Mutex
}

func newSubscription() *subscription {
	return &subscription{notify: new(int64)}
}

func (s *subscription) Notify(events []object.GroupEvent) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.events = append(s.events, events...)
	atomic.StoreInt64(s.notify, 1)
}

func pollResponse(events []object.GroupEvent) object.LongpollBotResponse {
	return object.LongpollBotResponse{
		Ts:      strconv.FormatInt(time.Now().Unix(), 10),
		Updates: events,
		Failed:  0,
	}
}

func (s *subscription) returnEvents() (r object.LongpollBotResponse) {
	s.lock.Lock()
	defer s.lock.Unlock()

	r = pollResponse(s.events)
	s.events = nil
	return
}

func (s *subscription) Poll(ctxt context.Context) object.LongpollBotResponse {
	for {
		select {
		case <-ctxt.Done():
			return pollResponse(nil)
		default:
			if atomic.LoadInt64(s.notify) == 1 {
				atomic.StoreInt64(s.notify, 0)
				return s.returnEvents()
			}
			runtime.Gosched()
		}
	}
}

func (s *subscription) Close() error {
	s.lock.Lock()
	defer s.lock.Unlock()

	atomic.StoreInt64(s.notify, 1)
	return nil
}
