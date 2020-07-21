package plugin

type Plugin interface {
	Name() string
}

type Loader interface {
	Plugins() map[string]Plugin
	Load(Plugin) error
	Unload(Plugin) error
	Find(string) Plugin
}

type plugin struct {
}

func (p *plugin) Name() string {
	return ""
}

func Plugins() map[string]Plugin {
	return defaultLoader.Plugins()
}

func Load(plugin Plugin) error {
	return defaultLoader.Load(plugin)
}

func NewLoader() Loader {
	return newLoader()
}
