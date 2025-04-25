package util

import "sync"

type SyncMap[K comparable, V any] struct {
	inner sync.Map
	empty V
}

func (m *SyncMap[K, V]) Load(key K) (value V, ok bool) {
	load, ok := m.inner.Load(key)
	if !ok {
		return m.empty, false
	}
	return load.(V), true
}

func (m *SyncMap[K, V]) Store(key K, value V) {
	m.inner.Store(key, value)
}

func (m *SyncMap[K, V]) Delete(key K) {
	m.inner.Delete(key)
}

func (m *SyncMap[K, V]) Range(f func(key K, value V) bool) {
	m.inner.Range(func(key, value any) bool {
		return f(key.(K), value.(V))
	})
}

func (m *SyncMap[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	v, vLoaded := m.inner.LoadOrStore(key, value)
	return v.(V), vLoaded
}

func (m *SyncMap[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	v, vLoaded := m.inner.LoadAndDelete(key)
	if !vLoaded {
		return m.empty, false
	}
	return v.(V), vLoaded
}
