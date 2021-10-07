package ipc

// Postman for the IPC.
// Daemon for delivering messages.

func postman(n *node) {
	for {
	msg := <-(n.outbox)
	to, ok := nodes[msg.To]
	if !ok { continue }
	to.inbox<-msg
	}
}