package main

import (
	"context"
	"net/http"
	"sync"
	"time"

	plugin "github.com/sters/onstatic/pluginapi"
)

type timer struct {
	plugin.BasicServer

	latestUpdate    time.Time
	mux             sync.RWMutex
	ticker          *time.Ticker
	tickerCtx       context.Context
	tickerCtxCancel context.CancelFunc
}

func (*timer) Name(context.Context, *plugin.EmptyMessage) (*plugin.NameResponse, error) {
	return &plugin.NameResponse{
		Name: "timer",
	}, nil
}

func (t *timer) Start(context.Context, *plugin.EmptyMessage) (*plugin.EmptyMessage, error) {
	t.updateLatestUpdate(time.Now())
	t.ticker = time.NewTicker(5 * time.Second)
	t.tickerCtx, t.tickerCtxCancel = context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-t.ticker.C:
				plugin.Log("tick")
				t.updateLatestUpdate(time.Now())
			case <-t.tickerCtx.Done():
				t.ticker.Stop()
				return
			}
		}
	}()

	t.RegisterHandler(plugin.HTTPMethodGET, "/api/timer", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(t.readLatestUpdate()))
	})

	return &plugin.EmptyMessage{}, nil
}

func (t *timer) Stop(context.Context, *plugin.EmptyMessage) (*plugin.EmptyMessage, error) {
	t.tickerCtxCancel()

	return &plugin.EmptyMessage{}, nil
}

func (t *timer) updateLatestUpdate(tt time.Time) {
	t.mux.Lock()
	defer t.mux.Unlock()
	t.latestUpdate = tt
}

func (t *timer) readLatestUpdate() string {
	t.mux.RLock()
	defer t.mux.RUnlock()

	return t.latestUpdate.String()
}

func main() {
	plugin.Serve(&timer{})
}
