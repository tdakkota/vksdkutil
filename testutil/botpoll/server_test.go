package botpoll

import (
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/SevereCloud/vksdk/v2/object"

	"github.com/stretchr/testify/assert"
)

func createRequest(query url.Values) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, "/testurl", nil)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = query.Encode()
	return req, nil
}

func testParseArgs(t *testing.T, query url.Values) (time.Duration, string, error) {
	req, err := createRequest(query)
	if err != nil {
		t.Fatal(err)
	}

	return TestLongPoll{}.parseArgs(req)
}

func TestTestLongPoll_parseArgs(t *testing.T) {
	t.Run("invalid-wait", func(t *testing.T) {
		_, _, err := testParseArgs(t, url.Values{
			"wait": []string{"ainvalidwait"},
			"key":  []string{"testkey"},
		})

		assert.Error(t, err)
	})

	t.Run("zero-wait", func(t *testing.T) {
		_, _, err := testParseArgs(t, url.Values{
			"wait": []string{"0"},
			"key":  []string{"testkey"},
		})

		assert.Error(t, err)
	})

	t.Run("invalid-key", func(t *testing.T) {
		_, _, err := testParseArgs(t, url.Values{
			"wait": []string{"25"},
		})

		assert.Error(t, err)
	})

	t.Run("max-90-wait", func(t *testing.T) {
		wait, _, err := testParseArgs(t, url.Values{
			"wait": []string{"10000"},
			"key":  []string{"testkey"},
		})
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, wait.Nanoseconds(), (90 * time.Second).Nanoseconds())
	})
}

func testHandler(url string, args url.Values) func(t *testing.T) {
	return func(t *testing.T) {
		server := NewTestLongPoll()
		defer server.Close()

		assert.HTTPError(t, server.ServeHTTP, "GET", url, args)
	}
}

func TestTestLongPoll_ServeHTTP(t *testing.T) {
	t.Run("fails-parse-args", testHandler("/test", url.Values{
		"wait": []string{"what"},
		"key":  []string{"abc"},
	}))

	t.Run("fails-get-subscriptions", testHandler("/test", url.Values{
		"wait": []string{"25"},
		"key":  []string{"abc"},
	}))
}

func TestTestLongPoll_Unsubscribe(t *testing.T) {
	t.Run("exists", func(t *testing.T) {
		server := NewTestLongPoll()
		defer server.Close()

		params := server.Subscribe()
		assert.Len(t, server.subscriptions.subs, 1)

		server.Unsubscribe(params)
		assert.Empty(t, server.subscriptions.subs)
	})

	t.Run("not-exists", func(t *testing.T) {
		server := NewTestLongPoll()
		defer server.Close()

		server.Unsubscribe(object.MessagesLongPollParams{})
		assert.Empty(t, server.subscriptions.subs)
	})
}
