package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSyncMap(t *testing.T) {
	var m SyncMap[string, *int]
	v := 1
	m.Store("1", &v)
	value, ok := m.Load("1")
	assert.True(t, ok)
	assert.Equal(t, *value, v)

	m.Delete("1")
	value2, ok := m.Load("1")
	assert.False(t, ok)
	assert.Nil(t, value2)
}

func TestSyncMap2(t *testing.T) {
	var m SyncMap[string, *int]
	v1 := 1
	v2 := 2
	m.Store("1", &v1)
	actual, loaded := m.LoadOrStore("1", &v2)
	assert.Equal(t, actual, &v1)
	assert.True(t, loaded)

	actual, loaded = m.LoadOrStore("2", &v2)
	assert.Equal(t, actual, &v2)
	assert.False(t, loaded)
}

func TestSyncMap3(t *testing.T) {
	var m SyncMap[string, *int]
	v1 := 1
	m.Store("1", &v1)
	actual, loaded := m.LoadAndDelete("1")
	assert.Equal(t, actual, &v1)
	assert.True(t, loaded)

	actual, loaded = m.LoadAndDelete("1")
	assert.Nil(t, actual)
	assert.False(t, loaded)
}
