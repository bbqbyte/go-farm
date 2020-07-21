package component

import (
	"errors"
	"keywea.com/cloud/pblib/pbconfig"
)

var ComponentNotImplemented = errors.New("Component Not Implemented")

type Component interface {
	Create(*ComponentInstConfig) (interface{}, error)
	Init(inst interface{}, c *ComponentInstConfig) error
	Start(inst interface{}, c *ComponentInstConfig) error
	Stop(inst interface{}, c *ComponentInstConfig) error
	Update(inst interface{}, c *ComponentInstConfig) error
	Destroy(inst interface{}, c *ComponentInstConfig) error
}

type ComponentInstConfig struct {
	CompID string // 组件ID
	Name string // 实例名称
	Config *pbconfig.Configor
}

// 默认组件实现
type DefaultComponent struct {}

func (dfc *DefaultComponent) Create(*ComponentInstConfig) (interface{}, error) {
	return nil, ComponentNotImplemented
}
func (dfc *DefaultComponent) Init(inst interface{}, c *ComponentInstConfig) error {
	return nil
}
func (dfc *DefaultComponent) Start(inst interface{}, c *ComponentInstConfig) error {
	return ComponentNotImplemented
}

func (dfc *DefaultComponent) Stop(inst interface{}, c *ComponentInstConfig) error {
	return nil
}

func (dfc *DefaultComponent) Update(inst interface{}, c *ComponentInstConfig) error {
	return nil
}

func (dfc *DefaultComponent) Destroy(inst interface{}, c *ComponentInstConfig) error {
	return nil
}
