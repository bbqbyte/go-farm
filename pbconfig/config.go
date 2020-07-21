package pbconfig

import (
	"fmt"
)

type Configor interface {
	SetString(key, val string) error

	GetString(key string, defaultVal ...string) string
	GetInt(key string, defaultVal ...int) (int, error)
	GetInt64(key string, defaultVal ...int64) (int64, error)
	GetBool(key string) (bool, error)
	GetFloat(key string, defaultVal ...float64) (float64, error)

	GetRawValue(key string) (interface{}, error)
	SaveFile(file string) error
}

type Config interface {
	Load(file string) (Configor, error)
	ParseData(data []byte) (Configor, error)
}

var adapters = make(map[string]Config)

// Register
func Register(name string, adapter Config) {
	if adapter == nil {
		panic("config: Register adapter is nil")
	}
	if _, ok := adapters[name]; ok {
		panic("config: Register adapter duplicate")
	}
	adapters[name] = adapter
}

// NewConfig adapterName is json/yaml.
// filename is the config file path.
func NewConfig(adapterName, filename string) (Configor, error) {
	adapter, ok := adapters[adapterName]
	if !ok {
		return nil, fmt.Errorf("config: unknown adaptername %q", adapterName)
	}
	return adapter.Load(filename)
}

// NewConfigData adapterName is json/yaml.
// data is the config data.
func NewConfigData(adapterName string, data []byte) (Configor, error) {
	adapter, ok := adapters[adapterName]
	if !ok {
		return nil, fmt.Errorf("config: unknown adaptername %q", adapterName)
	}
	return adapter.ParseData(data)
}

// ExpandValueEnvForMap convert all string value with environment variable.
func ExpandValueEnvForMap(m map[string]interface{}) map[string]interface{} {
	for k, v := range m {
		switch value := v.(type) {
		case string:
			m[k] = ExpandValueEnv(value)
		case map[string]interface{}:
			m[k] = ExpandValueEnvForMap(value)
		case map[string]string:
			for k2, v2 := range value {
				value[k2] = ExpandValueEnv(v2)
			}
			m[k] = value
		}
	}
	return m
}
