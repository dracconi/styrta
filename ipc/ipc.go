package ipc

import (
	"sync"
	"sync/atomic"
)

// Erlang-like IPC
// Each node has an outbox and an inbox.
// They are buffered channels.
// Internal goroutine "ipc" manages the messages and passed them through.
// Nodes are identified through Pid. Can be registered at registry.go

type Pid uint32

type MsgTyp uint16

type MsgBody interface {}

type MsgDie struct {}

type Message struct {
	Id   uint32
	From Pid
	To   Pid
	Body MsgBody
}

type POBox chan Message

type node struct {
	pid           Pid
	inbox, outbox chan Message
}

var mnodes sync.Mutex
var nodes map[Pid]node
var last_pid Pid = 1

func Initialize() {
	nodes = make(map[Pid]node)
	RegistryInitialize()
}

func Make() (Pid, chan Message, chan Message) {
	atomic.AddUint32((*uint32)(&last_pid), 1)
	return last_pid, make(chan Message, 16), make(chan Message, 16)
}

// Adds node to the search list and starts the searching for it
func Add(pid Pid, inbox, outbox chan Message) bool {
	mnodes.Lock()
	defer mnodes.Unlock()

	// Append the node to the list for further book-keeping
	nodes[pid] = node{pid, inbox, outbox}

	// Start the daemon for the node
	v, _ := nodes[pid]
	go postman(&(v))

	return false
}

func Start(f func(Pid, POBox, POBox, []interface{}), w ...interface{}) Pid {
	p, i, o := Make()
	Add(p, i, o)
	go f(p, i, o, w)
	return p
}
