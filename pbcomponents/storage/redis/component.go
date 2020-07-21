package redis

import (
	"keywea.com/cloud/pblib/pb/component"
)

type Component struct {
	component.DefaultComponent
}

func (redisc *Component) Create(instConfig *component.ComponentInstConfig) (interface{}, error) {
	return NewPool(instConfig.Name, *instConfig.Config)
}

func (redisc *Component) Update(inst interface{}, instConfig *component.ComponentInstConfig) error {
	r := inst.(*RPool)
	config := r.ParseConfig(instConfig.Config)
	_, err := r.UpdatePool(config)
	return err
}

//func (redisc *Component) Destroy(inst interface{}, instConfig *pbcc.ComponentInstConfig) error {
//	r := inst.(*RPool)
//	return r.Destroy()
//}