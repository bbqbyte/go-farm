package plugin

import "sync"

type loader struct {
	sync.RWMutex
	plugins map[string]Plugin
}

func newLoader() *loader {
	return &loader{
		plugins: make(map[string]Plugin),
	}
}

var (
	defaultLoader = newLoader()
)

func (l *loader) Plugins() map[string]Plugin {
	l.RLock()
	defer l.RUnlock()

	return l.plugins
}

func (l *loader) Load(plugin Plugin) error {
	l.Lock()
	defer l.Unlock()

	_, ok := l.plugins[plugin.Name()]
	if !ok {
		l.plugins[plugin.Name()] = plugin
	}

	return nil
}

func (l *loader) Unload(plugin Plugin) error {
	l.Lock()
	defer l.Unlock()

	_, ok := l.plugins[plugin.Name()]
	if ok {
		delete(l.plugins, plugin.Name())
	}
	return nil
}

func (l *loader) Find(name string) Plugin {
	l.RLock()
	defer l.RUnlock()

	return l.plugins[name]
}
