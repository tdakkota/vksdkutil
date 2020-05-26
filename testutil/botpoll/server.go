package botpoll

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"time"

	"github.com/SevereCloud/vksdk/object"
)

type TestLongPoll struct {
	server *httptest.Server

	subscriptions *subscriptions
}

func NewTestLongPoll() TestLongPoll {
	lps := TestLongPoll{
		subscriptions: newSubscriptions(),
	}
	lps.server = httptest.NewServer(lps)

	return lps
}

func (l TestLongPoll) Client() *http.Client {
	return l.server.Client()
}

func (l TestLongPoll) parseArgs(r *http.Request) (time.Duration, string, error) {
	err := r.ParseForm()
	if err != nil {
		return 0, "", err
	}

	wait, err := strconv.ParseInt(r.Form.Get("wait"), 10, 64)
	if err != nil {
		return 0, "", fmt.Errorf("argument 'wait' is invalid")
	}

	if wait < 1 {
		return 0, "", fmt.Errorf("argument 'wait' should be greater than zero")
	}

	if wait > 90 {
		wait = 90
	}

	key := r.Form.Get("key")
	if key == "" {
		return 0, "", fmt.Errorf("argument 'key' is invalid")
	}

	return time.Duration(wait) * time.Second, key, nil
}

func (l TestLongPoll) SendMessage(msg object.MessagesMessage) error {
	message, err := json.Marshal(object.MessageNewObject{
		Message: msg,
	})
	if err != nil {
		return err
	}

	l.NotifyOne(object.GroupEvent{
		Type:   object.EventMessageNew,
		Object: message,
	})

	return nil
}

func (l TestLongPoll) NotifyOne(event object.GroupEvent) {
	l.Notify([]object.GroupEvent{event})
}

func (l TestLongPoll) Notify(events []object.GroupEvent) {
	l.subscriptions.Notify(events)
}

func (l TestLongPoll) Subscribe() object.MessagesLongpollParams {
	return object.MessagesLongpollParams{
		Key:    l.subscriptions.Create(),
		Server: l.server.URL,
		Ts:     int(time.Now().Unix()),
	}
}

func (l TestLongPoll) Unsubscribe(o object.MessagesLongpollParams) {
	l.subscriptions.Delete(o.Key)
}

func (l TestLongPoll) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	wait, key, err := l.parseArgs(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(object.Error{
			Message: fmt.Errorf("invalid test request: %w", err).Error(),
		})
	}

	subs, ok := l.subscriptions.Get(key)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(object.Error{
			Message: fmt.Errorf("key doesn't exists").Error(),
		})
	}

	ctxt, cancel := context.WithTimeout(r.Context(), wait*time.Second)
	defer cancel()

	_ = json.NewEncoder(w).Encode(subs.Poll(ctxt))
}

func (l TestLongPoll) Close() error {
	l.server.CloseClientConnections()
	l.server.Close()

	_ = l.subscriptions.Close()
	return nil
}
