package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/alphadose/haxmap"
	plugin "github.com/sters/onstatic/pluginapi"
)

const storageFilename = "data.json"

type storage struct {
	plugin.BasicServer

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

	app.RegisterHandler(plugin.HTTPMethodGET, "/api/storage/read", func(w http.ResponseWriter, r *http.Request) {
		d, _ := app.store.Get("data")

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(fmt.Sprintf("Data = %d", d)))
	})

	app.RegisterHandler(plugin.HTTPMethodGET, "/api/storage/write", func(w http.ResponseWriter, r *http.Request) {
		if d, ok := app.store.Get("data"); ok {
			app.store.Set("data", d+1)
		} else {
			app.store.Set("data", 1)
		}
	})

	return &plugin.EmptyMessage{}, nil
}

func (app *storage) Stop(context.Context, *plugin.EmptyMessage) (*plugin.EmptyMessage, error) {
	app.flush()
	return &plugin.EmptyMessage{}, nil
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
