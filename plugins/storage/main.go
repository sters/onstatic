package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/alphadose/haxmap"
	plugin "github.com/sters/onstatic/pluginapi"
)

const storageFilename = "data.json"

type storage struct {
	plugin.OnstaticPluginServer

	store *haxmap.Map[string, int64]
}

func (*storage) Name(context.Context, *plugin.EmptyMessage) (*plugin.NameResponse, error) {
	return &plugin.NameResponse{
		Name: "storage",
	}, nil
}

func (app *storage) Start(context.Context, *plugin.EmptyMessage) (*plugin.EmptyMessage, error) {
	app.store = haxmap.New[string, int64]()
	app.load()

	return &plugin.EmptyMessage{}, nil
}

func (app *storage) Stop(context.Context, *plugin.EmptyMessage) (*plugin.EmptyMessage, error) {
	app.flush()
	return &plugin.EmptyMessage{}, nil
}

func (app *storage) Handle(ctx context.Context, req *plugin.HandleRequest) (*plugin.HandleResponse, error) {
	// plugin.Handle(ctx, "get", "/api/storage/read", func() {
	// })

	switch req.Path {
	case "/api/storage/read":
		d, _ := app.store.Get("data")

		return &plugin.HandleResponse{
			Body: fmt.Sprintf("Data = %d", d),
		}, nil

	case "/api/storage/write":
		if d, ok := app.store.Get("data"); ok {
			app.store.Set("data", d+1)
		} else {
			app.store.Set("data", 1)
		}

		return &plugin.HandleResponse{
			Body: "ok",
		}, nil
	}

	return nil, plugin.ErrPluginNotHandledPath
}

func (app *storage) load() {
	f, err := os.OpenFile(storageFilename, os.O_RDONLY|os.O_CREATE, 0o666)
	if err != nil {
		plugin.Log("failed to open storage file: %v", err)
		return // skip
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	dec.UseNumber()
	result := map[string]json.Number{}
	if err := dec.Decode(&result); err != nil {
		plugin.Log("failed to read file: %v", err)
		return
	}

	for k, v := range result {
		i, err := v.Int64()
		if err != nil {
			continue
		}

		app.store.Set(k, i)
	}
}

func (app *storage) flush() {
	result := map[string]int64{}
	app.store.ForEach(func(key string, value int64) bool {
		result[key] = value
		return true
	})

	f, err := os.Create(storageFilename)
	if err != nil {
		plugin.Log("failed to open storage file: %v", err)
		return // skip
	}
	defer f.Close()

	plugin.Log("%+v", app.store)
	plugin.Log("%+v", result)

	enc := json.NewEncoder(f)
	if err := enc.Encode(result); err != nil {
		plugin.Log("failed to write file: %v", err)
		return
	}

	plugin.Log("data flushed")
}

func main() {
	plugin.Serve(&storage{})
}
