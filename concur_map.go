package constdp

//@formatter:off

import (
    "sync"
)

type Map struct {
    sync.RWMutex
    dict map[string]struct{}
}

func NewMap() *Map {
    return &Map{dict:make(map[string]struct{})}
}

func (m *Map) HasKey(key string) bool {
    m.RLock(); defer m.RUnlock()
    var _, ok = m.dict[key]
    return ok
}

func (m *Map) Set(key string) {
    m.Lock(); defer m.Unlock()
    m.dict[key] = struct{}{}
}

func (m *Map) Size() int {
    m.RLock(); defer m.RUnlock()
    var v = len(m.dict)
    return v
}

func (m *Map) Keys () []string {
    m.Lock(); defer m.Unlock()
    var keys = make([]string , m.Size())
    for k := range m.dict {
        keys = append(keys, k)
    }
    return keys
}
