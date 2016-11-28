package age

import (
	"testing"
	"time"
)

func TestPut(t *testing.T) {
	var m Map
	tab := []interface{}{
		"test1",
		"test2",
		30,
		struct{ A string }{A: "aaa"},
	}
	for i, v := range tab {
		if _, ok := m.Get(v); ok {
			t.Errorf("Exist(%#v) = true; want false", v)
		}
		m.Put(v, v)
		if n := m.Len(); n != i+1 {
			t.Errorf("Len() = %d; want %d", n, i+1)
		}
		if _, ok := m.Get(v); !ok {
			t.Errorf("Exist(%#v) = false; want true", v)
		}
	}
}

func TestDelete(t *testing.T) {
	const key = "test"

	var m Map
	m.TimeToLive = time.Millisecond
	m.Delete(key)

	m.Put(key, 20)
	m.Delete(key)
	if _, ok := m.Get(key); ok {
		t.Errorf("Get(%q) = true; want false", key)
	}
	<-time.After(m.TimeToLive)
	a := m.Outdated()
	if len(a) != 0 {
		t.Errorf("Outdated() = %d items; want 0 itme", len(a))
	}
}

func TestOutdated(t *testing.T) {
	var m Map
	m.TimeToLive = 500 * time.Millisecond
	m.Put("test1", 1)
	m.Put("test2", 2)
	<-time.After(m.TimeToLive / 2)
	m.Put("test3", 3)
	<-time.After(m.TimeToLive / 2)
	m.Put("test4", 4)
	a := m.Outdated()
	if len(a) != 2 {
		t.Errorf("Outdated() = %d items; want %d items", len(a), 2)
	}
	if m.Len() != 2 {
		t.Errorf("Len() = %d; want %d", m.Len(), 2)
	}
}
