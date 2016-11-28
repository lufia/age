package age

import (
	"time"
)

type item struct {
	key      interface{}
	v        interface{}
	deadline *time.Time
}

type Map struct {
	TimeToLive time.Duration
	m          map[interface{}]*item
	a          []*item
}

func (m *Map) Len() int {
	return len(m.m)
}

func (m *Map) Get(key interface{}) (interface{}, bool) {
	item, ok := m.m[key]
	if !ok {
		return nil, false
	}
	return item.v, true
}

func (m *Map) Put(key, v interface{}) {
	p := &item{key: key, v: v}
	if m.TimeToLive > 0 {
		t := time.Now().Add(m.TimeToLive)
		p.deadline = &t
	}
	if m.m == nil {
		m.m = make(map[interface{}]*item)
	}
	m.m[key] = p
	m.a = append(m.a, p)
}

func (m *Map) Delete(key interface{}) {
	p, ok := m.m[key]
	if !ok {
		return
	}
	delete(m.m, key)
	i := m.indexItem(p)
	if i < 0 {
		panic("can't happen")
	}
	m.a = append(m.a[:i], m.a[i+1:]...)
}

func (m *Map) indexItem(p *item) int {
	for i, v := range m.a {
		if v == p {
			return i
		}
	}
	return -1
}

func (m *Map) Outdated() []interface{} {
	now := time.Now()
	n := m.countOutdated(now)
	a := m.a[0:n]
	rv := make([]interface{}, len(a))
	for i, p := range a {
		delete(m.m, p.key)
		rv[i] = p.v
	}
	m.a = m.a[n:]
	return rv
}

func (m *Map) countOutdated(now time.Time) int {
	for i, v := range m.a {
		if v.deadline.After(now) {
			return i
		}
	}
	return len(m.a)
}
