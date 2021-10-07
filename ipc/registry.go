package ipc

import (
"sync" )

// Registry for the IPC.

var mnames sync.Mutex
var names map[string]Pid

func RegistryInitialize() {
	names = make(map[string]Pid)
}

func Register(pid Pid, name string) bool {
	mnames.Lock()
	defer mnames.Unlock()
	if _, ok := names[name]; ok {
	return true
	}

	names[name] = pid
	return false
}

func Resolve(name string) (Pid, bool) {
	mnames.Lock()
	defer mnames.Unlock()
	if _, ok := names[name]; !ok {
		return 0, false
	}
	return names[name], true
}