package module

import (
	"github.com/bbqbyte/go-farm/pbconfig"
	"github.com/bbqbyte/go-farm/pb"
	"github.com/bbqbyte/go-farm/pb/version"
)

var (
	pbModuleTree map[string]*PBModule
)

func init() {
	pbModuleTree = make(map[string]*PBModule)

	self() // 自身信息
}

type PBModule struct {
	name string
	path string
	pbVersion *version.PBVersion
	builtOn string
}

// 注册使用的模块
func RegisterModule(mod *PBModule) {
	if mod == nil {
		panic("pbModule: Register provider is nil")
	}
	if _, dup := pbModuleTree[mod.Key()]; dup {
		panic("pbModule: Register duplicate for provider " + mod.Key())
	}
	pbModuleTree[mod.Key()] = mod
}

func NewModule(name, path string) *PBModule {
	return &PBModule{
		name: name,
		path: path,
	}
}

func (m *PBModule) SetName(name string) {
	m.name = name
}

func (m *PBModule) SetPath(path string) {
	m.path = path
}

func (m *PBModule) SetVersion(version *version.PBVersion) {
	m.pbVersion = version
}

func (m *PBModule) SetBuiltOn(buildtime string) {
	m.builtOn = buildtime
}

func (m *PBModule) Name() string {
	return m.name
}

func (m *PBModule) Path() string {
	return m.path
}

func (m *PBModule) Key() string {
	return m.name + pb.SYMBOL_AT + m.path
}

func (m *PBModule) Version() *version.PBVersion {
	return m.pbVersion
}

func (m *PBModule) BuiltOn() string {
	return m.builtOn
}

func self() {
	s := NewModule("pblib", "pblib.core")
	vjson := "{\"major\": 1, \"minor\":0, \"revision\":0, \"buildver\":\"0\", \"preRelease\":\"\"}"
	c, err := pbconfig.NewConfigData("json", []byte(vjson))
	if err != nil {
		panic("pbModule: init pblib module Error")
	} else {
		pbv, _ := version.NewPBVersion(c)
		s.SetVersion(pbv)
	}

	RegisterModule(s)
}
