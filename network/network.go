package network

import ("github.com/dracconi/styrta/ipc"
"reflect"
"time"
"math/rand"
"bufio"
"golang.org/x/crypto/blake2s"
"strconv"
"net")

type MsgAdd struct {
	address string
}

var peers map[net.Addr]ipc.Pid
var Pid ipc.Pid

func connectTo(address string) (ipc.Pid, net.Conn) {
	conn, err := net.DialTimeout("tcp", address, time.Second)
	if err != nil {
		return 0, conn
	}
	p := ipc.Start(peerd, conn)
	peers[conn.RemoteAddr()] = p
	return p, conn
}

func peerdRead(conn net.Conn, com chan []byte) {
	buf := bufio.NewReader(conn)
	for {
		b, _ := buf.ReadBytes('\n')
		com <- b[:len(b)-1]
	}
}

// Daemon for a single peer
func peerd(pid ipc.Pid, inbox, outbox ipc.POBox, w []interface{}) {
	conn := w[0].(net.Conn)
	com := make(chan []byte, 4)
	router, _ := ipc.Resolve("router")
	go peerdRead(conn, com)

	for {
		select {
			case r := <-com: // data read
				m := parseMsg(r)
				m.To = router
				m.From = pid
				outbox<-m
			case w := <-inbox: // data to be written
				sigFrom := blake2s.Sum256([]byte(conn.LocalAddr().String()))
				w.From = ipc.Pid(atoUint32(sigFrom[0:4]))
				sigTo := blake2s.Sum256([]byte(conn.RemoteAddr().String()))
				w.To = ipc.Pid(atoUint32(sigTo[0:4]))
				sigId := blake2s.Sum256([]byte(time.Now().String() + " " + strconv.Itoa(int(pid)) + " " + conn.LocalAddr().String() + " " + strconv.Itoa(rand.Int())))
				w.Id = atoUint32(sigId[0:4])
				conn.Write(serializeMsg(w))
		}
	}
}

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