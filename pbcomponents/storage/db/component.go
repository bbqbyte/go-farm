package db

import (
	"keywea.com/cloud/pblib/pb/component"
)

type Component struct {
	component.DefaultComponent
}

func (dbc *Component) Create(instConfig *component.ComponentInstConfig) (interface{}, error) {
	return NewDB(instConfig.Name, *instConfig.Config)
}

func (dbc *Component) Update(inst interface{}, instConfig *component.ComponentInstConfig) error {
	r := inst.(*DS)
	_, err := r.UpdatePool(instConfig.Name, r.ParseConfig(instConfig.Config))
	return err
}
