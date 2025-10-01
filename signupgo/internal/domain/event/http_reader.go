package event

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

var ErrNotFound = errors.New("event not found")

type httpReader struct {
	base string
	c    *http.Client
}

func NewHTTPReader(base string, timeoutMs int) Reader {
	if timeoutMs <= 0 {
		timeoutMs = 2000
	}
	return &httpReader{
		base: base,
		c:    &http.Client{Timeout: time.Duration(timeoutMs) * time.Millisecond},
	}
}

func (r *httpReader) FindActiveByID(ctx context.Context, id string) (*Event, error) {
	req, _ := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/event/%s", r.base, id), nil)
	resp, err := r.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("event service status %d", resp.StatusCode)
	}

	var e Event
	if err := json.NewDecoder(resp.Body).Decode(&e); err != nil {
		return nil, err
	}
	if e.Canceled != 0 {
		// existe pero cancelado
		return &e, nil
	}
	return &e, nil
}
