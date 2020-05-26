package botpoll

import (
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func testParseArgs(t *testing.T, query url.Values) (time.Duration, string, error) {
	req, err := http.NewRequest(http.MethodGet, "/testurl", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.URL.RawQuery = query.Encode()

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
