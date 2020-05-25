package middleware

import (
	"fmt"
	"testing"
	"time"

	"github.com/SevereCloud/vksdk/api"
	"github.com/stretchr/testify/assert"
)

func TestRetry(t *testing.T) {
	attempts, timeout := int64(2), 50*time.Millisecond
	m := Retry(int(attempts), timeout)

	counter := attempts
	handler := m(func(method string, params api.Params) (api.Response, error) {
		counter--
		return api.Response{}, fmt.Errorf("test error %d", counter)
	})

	start := time.Now()
	_, err := handler("", api.Params{})

	assert.Error(t, err)
	assert.Zero(t, counter)
	assert.LessOrEqual(t, attempts*timeout.Nanoseconds(), time.Since(start).Nanoseconds())
}
