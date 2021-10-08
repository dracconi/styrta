package network

import ("testing"
"net"
"time"
"github.com/dracconi/styrta/ipc")

func TestMsg(t *testing.T) {
	ipc.Initialize()
	Initialize()
	m := ipc.MsgBody(MsgAdd{"localhost:2270"})

	p, i, o := ipc.Make()
	ipc.Add(p, i, o)

	l, err := net.Listen("tcp", "localhost:2270")
	if err != nil {
		t.Errorf("failed to make a listener")
	}

	o<-ipc.Message{1, p, Pid, m}

	_, err = l.Accept()
	if err != nil {
		t.Errorf("failed to accept")
	}

	// A delay so that the other goroutine can finalize adding the thing
	time.Sleep(time.Millisecond*10)

	if len(peers) == 0 {
	t.Errorf("not dialed?")
	}
}