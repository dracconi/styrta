package ipc

import "reflect"

// Postman for the IPC.
// Daemon for delivering messages.

func postman(n *node) {
	for {
		msg := <-(n.outbox)
		if msg.To == 0 && reflect.TypeOf(msg.Body).String() == "ipc.MsgDie" {
			n.inbox<-Message{0, 0, 0, MsgDie{}}
			return
		}
		to, ok := nodes[msg.To]
		if !ok {
			continue
		}
		to.inbox <- msg
	}
}
