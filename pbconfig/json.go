package pbconfig

import (
	"io/ioutil"
	"keywea.com/cloud/pblib/pbconverter"
	"sync"
	"os"
	"encoding/json"
	"errors"
	"fmt"
)

type JSONConfig struct {
}

func (jsc *JSONConfig) Load(file string) (Configor, error) {
	// File exists?
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return nil, err
	}
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return jsc.ParseData(data)
}

func (jsc *JSONConfig) ParseData(data []byte) (Configor, error) {
	o := &JSONObject{
		data: make(map[string]interface{}),
	}
	err := json.Unmarshal(data, &o.data)
	if err != nil {
		var arrData []interface{}
		err2 := json.Unmarshal(data, &arrData)
		if err2 != nil {
			return nil, err
		}
		o.data["data"] = arrData
	}

	o.data = ExpandValueEnvForMap(o.data)

	return o, nil
}

type JSONObject struct {
	data map[string]interface{}
	sync.RWMutex
}

func (jo *JSONObject) getData(key string) interface{} {
	jo.Lock()
	defer jo.Unlock()

	if v, ok := jo.data[key]; ok {
		return v
	}
	return nil
}

func (jo *JSONObject) SetString(key, val string) error {
	jo.Lock()
	defer jo.Unlock()
	jo.data[key] = val
	return nil
}

func (jo *JSONObject) GetString(key string, defaultVal ...string) string {
	val := jo.getData(key)
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

func (jo *JSONObject) GetInt(key string, defaultVal ...int) (int, error) {
	val := jo.getData(key)
	if val != nil {
		if v, ok := val.(float64); ok {
			return int(v), nil
		}
	}
	if len(defaultVal) > 0 {
		return defaultVal[0], nil
	}
	return 0, errors.New("get Int Error on key:" + key)
}

func (jo *JSONObject) GetInt64(key string, defaultVal ...int64) (int64, error) {
	val := jo.getData(key)
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

func (jo *JSONObject) GetBool(key string) (bool, error) {
	val := jo.getData(key)
	if val != nil {
		return pbconverter.ToBoolean(val)
	}
	return false, fmt.Errorf("not exist key: %q", key)
}

func (jo *JSONObject) GetFloat(key string, defaultVal ...float64) (float64, error) {
	val := jo.getData(key)
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

func (jo *JSONObject) GetRawValue(key string) (interface{}, error) {
	val := jo.getData(key)
	if val != nil {
		return val, nil
	}
	return nil, errors.New("not exist key")
}

func (jo *JSONObject) SaveFile(file string) error {
	// Write configuration file by filename.
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	b, err := json.MarshalIndent(jo.data, "", "  ")
	if err != nil {
		return err
	}
	_, err = f.Write(b)
	return err
}

func init() {
	Register("json", &JSONConfig{})
}
