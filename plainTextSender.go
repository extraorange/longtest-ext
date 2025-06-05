package main

import (
	"math/rand"
	"sync"
	"time"
)

type PlainTextRequest struct {
	Line string
}

func (r *PlainTextRequest) Serialize() ([]byte, error) {
	return []byte(r.Line), nil
}

func (r *PlainTextRequest) Size() int {
	return len(r.Line)
}

func NewPlainTextSender(opts LogSenderOpts) *GenericSender {
	headers := make(map[string]string, len(opts.Headers)+1)
	if opts.Headers != nil {
		for k, v := range opts.Headers {
			headers[k] = v
		}
	}
	headers["Content-Type"] = "text/plain"
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	sender := &GenericSender{
		LogSenderOpts: LogSenderOpts{
			ID:         opts.ID,
			Containers: opts.Containers,
			Lines:      opts.Lines,
			LinesPS:    opts.LinesPS,
			URL:        opts.URL,
			ReaderURL:  opts.ReaderURL,
			Headers:    headers,
		},
		path:       "/test-lines",
		timeout:    time.Second,
		numOfSends: opts.LinesPS,
		rnd:        rnd,
		mtx:        sync.Mutex{},
	}
	sender.generate = func() IRequest {
		line := sender.pickRandom(sender.Lines)
		return &PlainTextRequest{Line: line}
	}

	return sender
}
