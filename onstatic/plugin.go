package onstatic

import (
	"context"
	"net/http"
)

// type loadedPluginsStruct struct {
// 	m       sync.RWMutex
// 	plugins map[string]repoPlugins // reponame => plugin name => handlers
// }

// type repoPlugins struct {
// 	plugins map[string]repoPlugin // plugin name => handlers
// }

// type repoPlugin struct {
// 	lastModified time.Time
// 	handlers     oplugin.Handlers
// 	api          oplugin.API
// }

// var loadedPlugins = &loadedPluginsStruct{
// 	m:       sync.RWMutex{},
// 	plugins: map[string]repoPlugins{},
// }

func handlePlugin(ctx context.Context, requestPath string) http.HandlerFunc {
	return nil

	// pathes := strings.Split(requestPath, "/")
	// if len(pathes) < 2 {
	// 	return nil
	// }

	// repoName := pathes[1]
	// pathUnderRepo := "/" + strings.Join(pathes[2:], "/")
	// repoFs := getRepoFs(repoName)

	// fsInfos, err := repoFs.ReadDir(oplugin.PluginDir)
	// if err != nil {
	// 	return nil
	// }

	// loadPluginIfPossible(ctx, repoName, repoFs, fsInfos)

	// loadedPlugins.m.RLock()
	// defer loadedPlugins.m.RUnlock()

	// repo, ok := loadedPlugins.plugins[repoName]
	// if !ok {
	// 	return nil
	// }

	// for _, handlers := range repo.plugins {
	// 	for p, h := range handlers.handlers {
	// 		if pathUnderRepo == string(p) {
	// 			return http.HandlerFunc(h)
	// 		}
	// 	}
	// }

	// return nil
}

// func loadPluginIfPossible(ctx context.Context, repoName string, repoFs billy.Filesystem, fsInfos []os.FileInfo) {
// 	for _, fsInfo := range fsInfos {
// 		if err := checkLastModTime(fsInfo, repoName); err != nil {
// 			zap.L().Warn("failed to load plugin, skip", zap.Error(err))
// 			continue
// 		}

// 		pluginName := fsInfo.Name()

// 		ep, err := loadPlugin(repoFs, pluginName)
// 		if err != nil {
// 			zap.L().Warn("failed to load plugin, skip", zap.Error(err))
// 			continue
// 		}

// 		api := ep(ctx, zap.L())
// 		api.Initialize(ctx)
// 		handlers := api.Handlers()

// 		loadedPlugins.m.Lock()
// 		defer loadedPlugins.m.Unlock()

// 		if _, ok := loadedPlugins.plugins[repoName]; !ok {
// 			loadedPlugins.plugins[repoName] = repoPlugins{
// 				plugins: map[string]repoPlugin{},
// 			}
// 		}

// 		loadedPlugins.plugins[repoName].plugins[pluginName] = repoPlugin{
// 			lastModified: fsInfo.ModTime(),
// 			handlers:     handlers,
// 			api:          api,
// 		}

// 		handlePaths := make([]string, len(handlers))
// 		for p := range handlers {
// 			handlePaths = append(handlePaths, string(p))
// 		}
// 		zap.L().Info(
// 			"plugin loaded",
// 			zap.String("repository", repoName),
// 			zap.String("name", pluginName),
// 			zap.Strings("handlePaths", handlePaths),
// 		)
// 	}
// }

// func checkLastModTime(pluginFsInfo os.FileInfo, repoName string) error {
// 	loadedPlugins.m.RLock()
// 	defer loadedPlugins.m.RUnlock()

// 	repo, ok := loadedPlugins.plugins[repoName]
// 	if !ok {
// 		return nil
// 	}

// 	p := repo.plugins[pluginFsInfo.Name()]
// 	if !ok {
// 		return nil
// 	}

// 	if pluginFsInfo.ModTime().Before(p.lastModified) {
// 		return failure.Unexpected("the plugin is not updated")
// 	}

// 	return nil
// }

// var loadPlugin = loadPluginActual // for testing

// func loadPluginActual(repoFs billy.Filesystem, filename string) (ep oplugin.EntryPoint, reterr error) {
// 	path := repoFs.Join(repoFs.Root(), oplugin.PluginDir, filename)

// 	defer func() {
// 		if err := recover(); err != nil {
// 			reterr = failure.Unexpected(
// 				fmt.Sprintf("%+v", err),
// 				failure.Messagef("failed to load plugin: cannot open plugin: %s", path),
// 			)
// 		}
// 	}()
// 	p, err := plugin.Open(path)
// 	if err != nil {
// 		return nil, failure.Wrap(err, failure.Messagef("failed to load plugin: cannot open plugin: %s", path))
// 	}

// 	sym, err := p.Lookup(oplugin.PluginExportVariableName)
// 	if err != nil {
// 		return nil, failure.Wrap(err, failure.Messagef("failed to load plugin: cannot find entry point: %s", oplugin.PluginExportVariableName))
// 	}

// 	ep, ok := sym.(oplugin.EntryPoint)
// 	if !ok {
// 		epp, ok := sym.(*oplugin.EntryPoint)
// 		if !ok {
// 			fmt.Printf("%#q\n", sym)
// 			return nil, failure.Unexpected("failed to load plugin: missing entry point")
// 		}
// 		return *epp, nil
// 	}

// 	return ep, nil
// }

func CleanupLoadedPlugins(ctx context.Context) {
	// loadedPlugins.m.Lock()
	// defer loadedPlugins.m.Unlock()

	// for _, repoPlugins := range loadedPlugins.plugins {
	// 	for _, p := range repoPlugins.plugins {
	// 		p.api.Stop(ctx)
	// 	}
	// }
	// loadedPlugins.plugins = map[string]repoPlugins{}
}
