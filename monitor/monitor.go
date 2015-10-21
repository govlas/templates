package mtutils

import (
	"runtime"
	"sync"
)

// template type Monitor(A)
type A int

type MonitoredObj A
type AccessFunc func(MonitoredObj)

type Monitor struct {
	sync.Mutex
	obj MonitoredObj
}

func NewMonitor(obj MonitoredObj) *Monitor {
	ret := new(Monitor)
	ret.obj = obj
	runtime.SetFinalizer(ret, func(m *Monitor) {
		m.Unlock()
	})
	return ret
}

func (m *Monitor) Capture() MonitoredObj {
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
