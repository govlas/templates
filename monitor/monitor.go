package mtutils

import (
	"runtime"
	"sync"
)

// template type Monitor(A)
type A int

type AccessFunc func(A)

type Monitor struct {
	sync.Mutex
	obj A
}

func NewMonitor(obj A) *Monitor {
	ret := new(Monitor)
	ret.obj = obj
	runtime.SetFinalizer(ret, func(m *Monitor) {
		m.Unlock()
	})
	return ret
}

func (m *Monitor) Capture() A {
	m.Lock()
	return m.obj
}

func (m *Monitor) Release() {
	m.Unlock()
}

func (m *Monitor) Access(f AccessFunc) {
	m.Lock()
	defer m.Unlock()
	f(m.obj)
}
