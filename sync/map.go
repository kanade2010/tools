package sync 

import (
    "sync"
)

type Map struct {
    m sync.Map
}

func (sm *Map) Set (k,v interface{}) {
    sm.m.Store(k,v)
}

func (sm *Map) Get (k interface{}) (interface{}, bool) {
    return sm.m.Load(k)
}

func (sm *Map) Del(key interface{}) {
    sm.m.Delete(key)
}

func (sm *Map) Range (funcs func(key, value interface{}) bool) {
    sm.m.Range(funcs)
}