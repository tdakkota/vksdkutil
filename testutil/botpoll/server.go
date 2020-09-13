package botpoll

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/events"
	"net/http"
	"net/http/httptest"
	"strconv"
	"time"

	"github.com/SevereCloud/vksdk/v2/object"
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

func (l TestLongPoll) URL() string {
	return l.server.URL
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
	message, err := json.Marshal(events.MessageNewObject{
		Message: msg,
	})
	if err != nil {
		return err
	}

	l.NotifyOne(events.GroupEvent{
		Type:   events.EventMessageNew,
		Object: message,
	})

	return nil
}

func (l TestLongPoll) NotifyOne(event events.GroupEvent) {
	l.Notify([]events.GroupEvent{event})
}

func (l TestLongPoll) Notify(events []events.GroupEvent) {
	l.subscriptions.Notify(events)
}

func (l TestLongPoll) Subscribe() object.MessagesLongPollParams {
	return object.MessagesLongPollParams{
		Key:    l.subscriptions.Create(),
		Server: l.server.URL,
		Ts:     int(time.Now().Unix()),
	}
}

func (l TestLongPoll) Unsubscribe(o object.MessagesLongPollParams) {
	l.subscriptions.Delete(o.Key)
}

func (l TestLongPoll) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	wait, key, err := l.parseArgs(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(api.Error{
			Message: fmt.Errorf("invalid test request: %w", err).Error(),
		})

		return
	}

	subs, ok := l.subscriptions.Get(key)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(api.Error{
			Message: fmt.Errorf("key doesn't exists").Error(),
		})

		return
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
