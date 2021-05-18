package onstatic

import (
	"fmt"
	"net/http"
	"os"
	"plugin"
	"strings"
	"sync"
	"time"

	"github.com/morikuni/failure"
	oplugin "github.com/sters/onstatic/onstatic/plugin"
	"go.uber.org/zap"
	"gopkg.in/src-d/go-billy.v4"
)

type loadedPluginsStruct struct {
	m       sync.RWMutex
	plugins map[string]repoPlugins // reponame => plugin name => handlers
}

type repoPlugins struct {
	plugins map[string]repoPlugin // plugin name => handlers
}

type repoPlugin struct {
	lastModified time.Time
	handlers     oplugin.Handlers
}

var loadedPlugins = &loadedPluginsStruct{
	m:       sync.RWMutex{},
	plugins: map[string]repoPlugins{},
}

func handlePlugin(requestPath string) http.HandlerFunc {
	pathes := strings.Split(requestPath, "/")
	if len(pathes) < 2 {
		return nil
	}

	repoName := pathes[1]
	pathUnderRepo := "/" + strings.Join(pathes[2:], "/")
	repoFs := getRepoFs(repoName)

	fsInfos, err := repoFs.ReadDir(oplugin.PluginDir)
	if err != nil {
		return nil
	}

	loadPluginIfPossible(repoName, repoFs, fsInfos)

	loadedPlugins.m.RLock()
	defer loadedPlugins.m.RUnlock()

	repo, ok := loadedPlugins.plugins[repoName]
	if !ok {
		return nil
	}

	for _, handlers := range repo.plugins {
		for p, h := range handlers.handlers {
			if pathUnderRepo == string(p) {
				return http.HandlerFunc(h)
			}
		}
	}

	return nil
}

func loadPluginIfPossible(repoName string, repoFs billy.Filesystem, fsInfos []os.FileInfo) {
	for _, fsInfo := range fsInfos {
		if err := checkLastModTime(fsInfo, repoName); err != nil {
			zap.L().Warn("failed to load plugin, skip", zap.Error(err))
			continue
		}

		pluginName := fsInfo.Name()

		ep, err := loadPlugin(repoFs, pluginName)
		if err != nil {
			zap.L().Warn("failed to load plugin, skip", zap.Error(err))
			continue
		}

		handlers := ep(zap.L()).Register()

		loadedPlugins.m.Lock()
		defer loadedPlugins.m.Unlock()

		if _, ok := loadedPlugins.plugins[repoName]; !ok {
			loadedPlugins.plugins[repoName] = repoPlugins{
				plugins: map[string]repoPlugin{},
			}
		}

		loadedPlugins.plugins[repoName].plugins[pluginName] = repoPlugin{
			lastModified: fsInfo.ModTime(),
			handlers:     handlers,
		}

		handlePaths := make([]string, len(handlers))
		for p := range handlers {
			handlePaths = append(handlePaths, string(p))
		}
		zap.L().Info(
			"plugin loaded",
			zap.String("repository", repoName),
			zap.String("name", pluginName),
			zap.Strings("handlePaths", handlePaths),
		)
	}
}

func checkLastModTime(pluginFsInfo os.FileInfo, repoName string) error {
	loadedPlugins.m.RLock()
	defer loadedPlugins.m.RUnlock()

	repo, ok := loadedPlugins.plugins[repoName]
	if !ok {
		return nil
	}

	p := repo.plugins[pluginFsInfo.Name()]
	if !ok {
		return nil
	}

	if pluginFsInfo.ModTime().Before(p.lastModified) {
		return failure.Unexpected("the plugin is not updated")
	}

	return nil
}

var loadPlugin = loadPluginActual // for testing

func loadPluginActual(repoFs billy.Filesystem, filename string) (oplugin.EntryPoint, error) {
	path := repoFs.Join(repoFs.Root(), oplugin.PluginDir, filename)

	p, err := plugin.Open(path)
	if err != nil {
		return nil, failure.Wrap(err, failure.Messagef("failed to load plugin: cannot open plugin: %s", path))
	}

	sym, err := p.Lookup(oplugin.PluginExportVariableName)
	if err != nil {
		return nil, failure.Wrap(err, failure.Messagef("failed to load plugin: cannot find entry point: %s", oplugin.PluginExportVariableName))
	}

	ep, ok := sym.(oplugin.EntryPoint)
	if !ok {
		epp, ok := sym.(*oplugin.EntryPoint)
		if !ok {
			fmt.Printf("%#q\n", sym)
			return nil, failure.Unexpected("failed to load plugin: missing entry point")
		}
		return *epp, nil
	}

	return ep, nil
}
