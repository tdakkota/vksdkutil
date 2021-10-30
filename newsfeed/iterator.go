package newsfeed

import (
	"context"
	"fmt"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
)

// Elem is a newsfeed item.
type Elem struct {
	object.NewsfeedNewsfeedItem
	object.ExtendedResponse
}

// Iterator is a newsfeed iterator.
type Iterator struct {
	client *api.VK

	// Current state.
	lastErr error
	// Buffer state.
	buf    []Elem
	bufCur int
	// Request state.
	nextFrom string
	// Request params.
	params func(b *params.NewsfeedGetBuilder)
}

// NewIterator creates new Iterator.
func NewIterator(client *api.VK, params func(b *params.NewsfeedGetBuilder)) *Iterator {
	return &Iterator{client: client, params: params}
}

// Next prepares the next message for reading with the Value method.
// It returns true on success, or false if there is no news anymore or an error happened while preparing it.
// Err should be consulted to distinguish between the two cases.
func (m *Iterator) Next(ctx context.Context) bool {
	if m.lastErr != nil {
		return false
	}

	if !m.bufNext() {
		// If buffer is empty, we should fetch next batch.
		if err := m.requestNext(ctx); err != nil {
			m.lastErr = err
			return false
		}
		// Try again with new buffer.
		return m.bufNext()
	}

	return true
}

func (m *Iterator) bufNext() bool {
	if len(m.buf)-1 <= m.bufCur {
		return false
	}

	m.bufCur++
	return true
}

type itemWrap object.NewsfeedNewsfeedItem

func (itemWrap) UnmarshalJSON() {}

func (m *Iterator) requestNext(ctx context.Context) error {
	b := params.NewNewsfeedGetBuilder()
	b.WithContext(ctx)
	if m.nextFrom != "" {
		b.StartFrom(m.nextFrom)
	}
	if m.params != nil {
		m.params(b)
	}

	var resp struct {
		Items []itemWrap `json:"items"`
		object.ExtendedResponse
		NextFrom string `json:"next_from"`
	}
	if err := m.client.RequestUnmarshal("newsfeed.get", &resp, b.Params); err != nil {
		return fmt.Errorf("request feed: %w", err)
	}
	m.nextFrom = resp.NextFrom

	m.bufCur = -1
	m.buf = m.buf[:0]
	for _, item := range resp.Items {
		m.buf = append(m.buf, Elem{
			NewsfeedNewsfeedItem: object.NewsfeedNewsfeedItem(item),
			ExtendedResponse:     resp.ExtendedResponse,
		})
	}

	return nil
}

// Value returns current item.
func (m *Iterator) Value() Elem {
	return m.buf[m.bufCur]
}

// Err returns the error, if any, that was encountered during iteration.
func (m *Iterator) Err() error {
	return m.lastErr
}
