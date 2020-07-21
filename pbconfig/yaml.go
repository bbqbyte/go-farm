package pbconfig

import (
	"os"
	"sync"
	"io/ioutil"
	"github.com/go-yaml/yaml"
	"keywea.com/cloud/pblib/pbconverter"
	"fmt"
	"errors"
)

type YAMLConfig struct {
}

func (yc *YAMLConfig) Load(file string) (Configor, error) {
	// File exists?
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return nil, err
	}
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return yc.ParseData(data)
}

func (yc *YAMLConfig) ParseData(data []byte) (Configor, error) {
	o := &YAMLObject{
		data: make(map[string]interface{}),
	}
	err := yaml.Unmarshal(data, &o.data)
	if err != nil {
		return nil, err
	}

	o.data = ExpandValueEnvForMap(o.data)

	return o, nil
}

type YAMLObject struct {
	data map[string]interface{}
	sync.RWMutex
}

func (yo *YAMLObject) getData(key string) interface{} {
	yo.Lock()
	defer yo.Unlock()

	if v, ok := yo.data[key]; ok {
		return v
	}
	return nil
}

func (yo *YAMLObject) SetString(key, val string) error {
	yo.Lock()
	defer yo.Unlock()
	yo.data[key] = val
	return nil
}

func (yo *YAMLObject) GetString(key string, defaultVal ...string) string {
	val := yo.getData(key)
	if val != nil {
		if v, ok := val.(string); ok {
			return v
		}
	}
	if len(defaultVal) > 0 {
		return defaultVal[0]
	}
	return ""
}

func (yo *YAMLObject) GetInt(key string, defaultVal ...int) (int, error) {
	val := yo.getData(key)
	if val != nil {
		if v, ok := val.(int); ok {
			return v, nil
		} else if v, ok := val.(int64); ok {
			return int(v), nil
		}
	}
	if len(defaultVal) > 0 {
		return defaultVal[0], nil
	}
	return 0, errors.New("get Int Error on key:" + key)
}

func (yo *YAMLObject) GetInt64(key string, defaultVal ...int64) (int64, error) {
	val := yo.getData(key)
	if val != nil {
		if v, ok := val.(float64); ok {
			return int64(v), nil
		}
	}
	if len(defaultVal) > 0 {
		return defaultVal[0], nil
	}
	return 0, errors.New("get Int64 Error on key:" + key)
}

func (yo *YAMLObject) GetBool(key string) (bool, error) {
	val := yo.getData(key)
	if val != nil {
		return pbconverter.ToBoolean(val)
	}
	return false, fmt.Errorf("not exist key: %q", key)
}

func (yo *YAMLObject) GetFloat(key string, defaultVal ...float64) (float64, error) {
	val := yo.getData(key)
	if val != nil {
		if v, ok := val.(float64); ok {
			return v, nil
		} else if v, ok := val.(int); ok {
			return float64(v), nil
		} else if v, ok := val.(int64); ok {
			return float64(v), nil
		}
	}
	if len(defaultVal) > 0 {
		return defaultVal[0], nil
	}
	return 0.0, errors.New("get Float Error on key:" + key)
}

func (yo *YAMLObject) GetRawValue(key string) (interface{}, error) {
	val := yo.getData(key)
	if val != nil {
		return val, nil
	}
	return nil, errors.New("not exist key")
}

func (yo *YAMLObject) SaveFile(file string) error {
	// Write configuration file by filename.
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	b, err := yaml.Marshal(yo.data)
	if err != nil {
		return err
	}
	_, err = f.Write(b)
	return err
}

func init() {
	Register("yaml", &YAMLConfig{})
}
