package botpoll

import (
	"context"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type subscription struct {
	notify int64

	events []events.GroupEvent
	lock   sync.Mutex
}

func newSubscription() *subscription {
	return &subscription{notify: 0}
}

func (s *subscription) Notify(events []events.GroupEvent) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.events = append(s.events, events...)
	atomic.StoreInt64(&s.notify, 1)
}

func pollResponse(events []events.GroupEvent) longpoll.Response {
	return longpoll.Response{
		Ts:      strconv.FormatInt(time.Now().Unix(), 10),
		Updates: events,
		Failed:  0,
	}
}

func (s *subscription) returnEvents() (r longpoll.Response) {
	s.lock.Lock()
	defer s.lock.Unlock()

	r = pollResponse(s.events)
	s.events = nil
	return
}

func (s *subscription) Poll(ctxt context.Context) longpoll.Response {
	for {
		select {
		case <-ctxt.Done():
			return pollResponse(nil)
		default:
			if atomic.LoadInt64(&s.notify) == 1 {
				atomic.StoreInt64(&s.notify, 0)
				return s.returnEvents()
			}
			runtime.Gosched()
		}
	}
}

func (s *subscription) Close() error {
	s.lock.Lock()
	defer s.lock.Unlock()

	atomic.StoreInt64(&s.notify, 1)
	return nil
}
