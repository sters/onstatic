package onstatic

import (
	"os"
	"path/filepath"
	"plugin"
	"strings"
	"sync"
	"time"

	"github.com/morikuni/failure"
	oplugin "github.com/sters/onstatic/onstatic/plugin"
	"gopkg.in/src-d/go-billy.v4"
)

type loadedPluginsStruct struct {
	m            sync.RWMutex
	lastModified time.Time
	plugins      map[string]oplugin.Handlers
}

var loadedPlugins = &loadedPluginsStruct{
	m:            sync.RWMutex{},
	lastModified: time.Time{},
	plugins:      map[string]oplugin.Handlers{},
}

func loadPluginIfPossible(requestPath string) error {
	pathes := strings.Split(requestPath, "/")
	if len(pathes) < 2 {
		return failure.Unexpected("request path is too short")
	}

	repoName := pathes[1]

	pluginDir := filepath.Join(getRepositoriesDir(), repoName)
	pluginFs := fsNew(pluginDir)

	fsInfo, err := pluginFs.Stat(oplugin.PluginFilePath)
	if err != nil {
		return failure.Wrap(err)
	}

	if err := checkLastModTime(fsInfo); err != nil {
		return failure.Wrap(err)
	}

	api, err := loadPlugin(pluginFs)
	if err != nil {
		return failure.Wrap(err)
	}

	loadedPlugins.m.Lock()
	defer loadedPlugins.m.Unlock()

	loadedPlugins.lastModified = fsInfo.ModTime()
	loadedPlugins.plugins[repoName] = api.Register()

	return nil
}

func checkLastModTime(fsInfo os.FileInfo) error {
	loadedPlugins.m.RLock()
	defer loadedPlugins.m.RUnlock()

	if fsInfo.ModTime().Before(loadedPlugins.lastModified) {
		return failure.Unexpected("the plugin is not updated")
	}

	return nil
}

func loadPlugin(pluginFs billy.Filesystem) (oplugin.API, error) {
	plug, err := plugin.Open(pluginFs.Join(oplugin.PluginFilePath))
	if err != nil {
		return nil, failure.Wrap(err)
	}

	sym, err := plug.Lookup(oplugin.PluginExportVariableName)
	if err != nil {
		return nil, failure.Wrap(err)
	}

	api, ok := sym.(oplugin.API)
	if !ok {
		return nil, failure.Unexpected("failed to load plugin: missing entry point")
	}

	return api, nil
}
