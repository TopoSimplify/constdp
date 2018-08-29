package constdp

import "sync"

const cachKeySize = 6

type cacheMap struct {
	sync.RWMutex
	dict map[[cachKeySize]int]struct{}
}

func newCacheMap(size int) *cacheMap {
	return &cacheMap{dict: make(map[[cachKeySize]int]struct{}, size)}
}

func (m *cacheMap) HasKey(key *[cachKeySize]int) bool {
	m.RLock()
	var _, ok = m.dict[*key]
	m.RUnlock()
	return ok
}

func (m *cacheMap) Set(key *[cachKeySize]int) {
	m.Lock()
	m.dict[*key] = struct{}{}
	m.Unlock()
}

func (m *cacheMap) Delete(key *[cachKeySize]int) {
	m.Lock()
	delete(m.dict, *key)
	m.Unlock()
}

func (m *cacheMap) Size() int {
	m.RLock()
	var v = len(m.dict)
	m.RUnlock()
	return v
}

func (m *cacheMap) Keys() [][cachKeySize]int {
	m.RLock()
	var keys = make([][cachKeySize]int, m.Size())
	for k := range m.dict {
		keys = append(keys, k)
	}
	m.RUnlock()
	return keys
}
