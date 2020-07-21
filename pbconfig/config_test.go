package pbconfig

import (
	"testing"
)

func TestParseData(t *testing.T) {
	s := "{\"name\": \"bbq\", \"tag\":{\"tname\": \"zz\"}}"
	c, err := NewConfigData("json", []byte(s))
	if err != nil {
		t.Fatal("parse json error")
	} else {
		t.Log(c.GetString("name", "haha"))
		t.Log(c.GetRawValue("tag"))
	}
}


func TestLoad(t *testing.T) {
	p, err := Expand("~/test/yaml.yaml")
	c, err := NewConfig("yaml", p)
	if err != nil {
		t.Fatal("parse yaml error", err.Error())
	} else {
		t.Log(c.GetString("languages"))
		t.Log(c.GetInt("level", 3))
		t.Log(c.GetBool("oo"))
	}
}