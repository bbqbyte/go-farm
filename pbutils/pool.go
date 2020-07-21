package pbutils

import (
	"sync"
	"errors"
)

var (
	ErrObjectNil     = errors.New("object is nil.")
	ErrFetchTimedOut = errors.New("fetch object timeout.")
)

type poolConfig struct {
	Factory  func() (interface{}, error)
	IsActive func(interface{}) bool
	Release  func(interface{})

	MaxCap  int
	MaxIdle int
	MinIdle int
}

type pool struct {
	objChan chan interface{}
	mu      *sync.Mutex
	config  poolConfig
}

type GenericPool interface {
	Get() (obj interface{}, err error)
	Put(obj interface{}) error
	Clear()
	Close()
	Len() (int)
}

func NewGenericPool(config poolConfig) (GenericPool, error) {
	p := &pool{
		objChan: make(chan interface{}, config.MaxCap),
		mu:      &sync.Mutex{},
		config:  config,
	}

	return p, nil
}

func (p *pool) Get() (obj interface{}, err error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	select {
	case obj = <-p.objChan:
		if p.config.IsActive(obj) {
			return
		}
		p.config.Release(obj)
	default:
		obj, err = p.config.Factory()
		if err != nil {
			return nil, err
		}
		return obj, nil
	}

	return
}

func (p *pool) Put(obj interface{}) error {
	if obj == nil {
		return ErrObjectNil
	}
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.config.IsActive(obj) {
		p.config.Release(obj)
	}

	select {
	case p.objChan <- obj:
	default:
		p.config.Release(obj)
	}

	return nil
}

func (p *pool) Clear() {

}

func (p *pool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	close(p.objChan)

	for c := range p.objChan {
		p.config.Release(c)
	}

	p.objChan = make(chan interface{}, p.config.MinIdle)
}

func (p *pool) Len() (int) {
	return len(p.objChan)
}
