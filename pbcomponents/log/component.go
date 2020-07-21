package log

import (
	"keywea.com/cloud/pblib/pb/component"
)

type Component struct {
	component.DefaultComponent
}

func (l *Component) Create(instConfig *component.ComponentInstConfig) (interface{}, error) {
	return NewLogWriter(instConfig.Name, *instConfig.Config)
}
