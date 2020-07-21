package pbapp

import (
	"keywea.com/cloud/pblib/pb/module"
	"keywea.com/cloud/pblib/pb/version"
	"keywea.com/cloud/pblib/pbconfig"
	"testing"
	"log"
)

func TestModule(t *testing.T) {
	mod := module.NewModule("ithink", "")
	module.RegisterModule(mod)
	s := "{\"major\": 1, \"minor\":0, \"revision\":1, \"buildver\":\"1\", \"preRelease\":\"Alpha\"}"
	c, err := pbconfig.NewConfigData("json", []byte(s))
	if err != nil {

	} else {
		kwversion, _ := version.NewPBVersion(c)
		mod.SetVersion(kwversion)
		log.Printf("%v", mod.Version())
	}
}

func TestVersion(t *testing.T) {
	s := "{\"major\": 1, \"minor\":0, \"revision\":1, \"buildver\":\"1\", \"preRelease\":\"Alpha\"}"
	c, err := pbconfig.NewConfigData("json", []byte(s))
	if err != nil {

	} else {
		kwversion, _ := version.NewPBVersion(c)
		log.Printf("%v", kwversion)
	}
}