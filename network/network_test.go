package network

import ("testing"
"net"
"time"
"github.com/dracconi/styrta/ipc")

func temporaryRouter(p ipc.Pid, i, o ipc.POBox, w []interface{}) {
	tcom := w[0].(ipc.POBox)

	tcom<- (<-i)
}

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

	ret := make(ipc.POBox)
	pr := ipc.Start(temporaryRouter, ret)
	ipc.Register(pr, "router")

	o<-ipc.Message{1, p, Pid, m}

	conn, err := l.Accept()
	if err != nil {
		t.Errorf("failed to accept")
	}

	// A delay so that the other goroutine can finalize adding the thing
	time.Sleep(time.Millisecond*10)

	if len(peers) == 0 {
	t.Errorf("not dialed?")
	}

	mhelo := ipc.Message{1, 546, 678, MsgHello{1, uint32(time.Now().UnixMilli())}}
	conn.Write(serializeMsg(mhelo))

	time.Sleep(time.Millisecond*50)

	mhelo.To = pr

	rmsg := <-ret
	mhelo.From = rmsg.From

	if mhelo != rmsg {
		t.Errorf("identity fail %+v %+v", mhelo, rmsg)
	}
}