package main 

import (
    "sync"
)

type SyncMap struct {
    m sync.Map
}

func (sm *SyncMap) Set (k,v interface{}){
    sm.m.Store(k,v)
}

func (sm *SyncMap) Get (k interface{}) interface{}{
    v ,exit := sm.m.Load(k)
    if exit {
        return v
    }
    return nil
}

func (sm *SyncMap) Del(key interface{}){
    sm.m.Delete(key)
}

func (sm *SyncMap) Range (funcs func(key, value interface{}) bool) {
    sm.m.Range(funcs)
}