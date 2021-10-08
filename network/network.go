package network

import ("github.com/dracconi/styrta/ipc"
"reflect"
"time"
"net")

type MsgAdd struct {
	address string
}

var peers map[net.Addr]ipc.Pid
var Pid ipc.Pid

func connectTo(address string) {
	conn, err := net.DialTimeout("tcp", address, time.Second)
	if err != nil {
		return
	}
	peers[conn.RemoteAddr()] = ipc.Start(peerd)
}

// Daemon for a single peer
func peerd(pid ipc.Pid, inbox, outbox ipc.POBox) {}

// Daemon for peers comm.
// Handles incoming messages
func peersd(pid ipc.Pid, inbox, outbox ipc.POBox) {
	for {
	msg := <-inbox
	switch reflect.TypeOf(msg.Body).String() {
		case "network.MsgAdd":
			connectTo(msg.Body.(MsgAdd).address)
		default:
	}
	}
}

func Initialize() {
	peers = make(map[net.Addr]ipc.Pid)
	p, i, o := ipc.Make()
	ipc.Add(p, i, o)
	go peersd(p, i, o)
	ipc.Register(p, "peers")
	Pid = p
}