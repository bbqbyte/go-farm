package component

import (
	"fmt"
	"sync"

	"keywea.com/cloud/pblib/pb/log"
)

// keywea component
type PBC struct {
	components map[string]Component // 注册的组件 compID:component

	instNames []string // 实例名称
	instConfigs map[string]*ComponentInstConfig // 实例配置
	instances map[string]interface{}

	mu sync.RWMutex
}

func NewPBC() *PBC {
	return &PBC{
		components: make(map[string]Component),
		instNames: make([]string, 0),
		instConfigs: make(map[string]*ComponentInstConfig),
		instances: make(map[string]interface{}),
	}
}

// 注册组件
func (pb *PBC) RegisterComponent(compID string, comp Component) error {
	pb.mu.Lock()
	defer pb.mu.Unlock()

	if _, ok := pb.components[compID]; ok {
		return fmt.Errorf("Component `%v` already registered", compID)
	}
	pb.components[compID] = comp
	return nil
}

func (pb *PBC) Init(instConfigs []*ComponentInstConfig) error {
	for _, instConfig := range instConfigs {
		_, err := pb.CreateInstance(instConfig)
		if err != nil {
			return err
		}
	}
	for _, instName := range pb.instNames {
		if e := pb.InitInstance(instName); e != nil {
			return e
		}
	}
	return nil
}

// 组件实例
func (pb *PBC) CreateInstance(instConfig *ComponentInstConfig) (interface{}, error) {
	pb.mu.Lock()
	defer pb.mu.Unlock()

	if comp, ok := pb.components[instConfig.CompID]; ok {
		if _, ok := pb.instConfigs[instConfig.Name]; ok {
			return pb.instances[instConfig.Name], fmt.Errorf("Component Inst `%v` already created", instConfig.Name)
		}
		instance, err := comp.Create(instConfig)
		if err == nil {
			pb.instNames = append(pb.instNames, instConfig.Name)
			pb.instConfigs[instConfig.Name] = instConfig
			pb.instances[instConfig.Name] = instance
		}
		return instance, err
	}
	return nil, fmt.Errorf("Component `%v` Not Registered", instConfig.Name)
}

// 获取实例
func (pb *PBC) Instance(name string) interface{} {
	if inst, ok := pb.instances[name]; ok {
		return inst
	}
	return nil
}

// 实例初始化
func (pb *PBC) InitInstance(name string) error {
	pb.mu.Lock()
	defer pb.mu.Unlock()

	if inst, ok := pb.instances[name]; ok {
		instConfig := pb.instConfigs[name]
		return pb.components[instConfig.CompID].Init(inst, instConfig)
	}
	return fmt.Errorf("Component Inst `%v` Not created", name)
}

// 实例启动
func (pb *PBC) StartInstance(name string) error {
	pb.mu.Lock()
	defer pb.mu.Unlock()

	if inst, ok := pb.instances[name]; ok {
		instConfig := pb.instConfigs[name]
		return pb.components[instConfig.CompID].Start(inst, instConfig)
	}
	return fmt.Errorf("Component Inst `%v` Not created", name)
}

// 实例停止
func (pb *PBC) StopInstance(name string) error {
	pb.mu.Lock()
	defer pb.mu.Unlock()

	if inst, ok := pb.instances[name]; ok {
		instConfig := pb.instConfigs[name]
		return pb.components[instConfig.CompID].Stop(inst, instConfig)
	}
	return fmt.Errorf("Component Inst `%v` Not created", name)
}

// 实例更新
func (pb *PBC) UpdateInstance(name string, instConfig *ComponentInstConfig) error {
	pb.mu.Lock()
	defer pb.mu.Unlock()

	if inst, ok := pb.instances[name]; ok {
		delete(pb.instConfigs, name)
		pb.instConfigs[name] = instConfig
		return pb.components[instConfig.CompID].Update(inst, instConfig)
	}
	return fmt.Errorf("Component Inst `%v` Not created", name)
}

// 实例销毁
func (pb *PBC) DestroyInstance(name string) error {
	pb.mu.Lock()
	defer pb.mu.Unlock()

	if inst, ok := pb.instances[name]; ok {
		instConfig := pb.instConfigs[name]
		err := pb.components[instConfig.CompID].Destroy(inst, instConfig)
		if err == nil {
			delete(pb.instConfigs, name)
			delete(pb.instances, name)
			for i, v := range pb.instNames {
				if v == name {
					pb.instNames = append(pb.instNames[:i], pb.instNames[i+1:]...)
					break
				}
			}
		}
	}
	return fmt.Errorf("Component Inst `%v` Not created", name)
}

func (pb *PBC) DestroyAll() error {
	for i := len(pb.instNames) - 1; i >= 0; i-- {
		err := pb.DestroyInstance(pb.instNames[i])
		if err != nil {
			plog.Error("Component instance destroy Error", log.Error(err))
		}
	}
	return nil
}
