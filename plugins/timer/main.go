package main

import (
	"context"
	"sync"
	"time"

	plugin "github.com/sters/onstatic/onstatic/plugin"
)

type timer struct {
	plugin.OnstaticPluginServer

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

	return &plugin.EmptyMessage{}, nil
}

func (t *timer) Stop(context.Context, *plugin.EmptyMessage) (*plugin.EmptyMessage, error) {
	t.tickerCtxCancel()

	return &plugin.EmptyMessage{}, nil
}

func (t *timer) Handle(ctx context.Context, req *plugin.HandleRequest) (*plugin.HandleResponse, error) {
	switch req.Path {
	case "/api/timer":
		return &plugin.HandleResponse{
			Body: t.readLatestUpdate(),
		}, nil
	}

	return nil, plugin.ErrPluginNotHandledPath
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
