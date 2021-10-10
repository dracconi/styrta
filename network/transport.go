package network

import (
"github.com/dracconi/styrta/ipc"
"bytes"
"time"
"reflect"
)

// Please look at format.txt for more info about the format of messages

const (
	tHello = iota + 1
	tAddPeer
	tTx
	tBlock
	tProof
	tInvite
)

func atoUint32(b []byte) uint32 {
	return uint32(b[0]) | (uint32(b[1])<<8) |  (uint32(b[2])<<16) | (uint32(b[3])<<24)
}

func toaUint32(u uint32) []byte {
	return []byte{byte(u), byte(u>>8), byte(u>>16), byte(u>>24)}
}

type MsgHello struct {
	version int8
	date uint32
}

func parseBody(b []byte) ipc.MsgBody {
	switch int(b[0]) {
		case tHello:
			return MsgHello{int8(b[1]), atoUint32(b[1:5])}
		case tAddPeer:
			n := atoUint32(b[1:5])
			return MsgAdd{string(b[5:5+n])}
		default:
	}
	return nil
}

func serializeBody(m ipc.MsgBody) []byte {
	switch reflect.TypeOf(m).String() {
		case "network.MsgHello":
			return append([]byte{tHello}, toaUint32(uint32(time.Now().UnixMilli()))...)
		case "network.MsgAdd":
			addr := m.(MsgAdd).address
			return append([]byte{tAddPeer}, append(toaUint32(uint32(len(addr))), []byte(addr)...)...)
		default:
	}
	return nil
}

func parseMsg(b []byte) ipc.Message {
	if len(b) < 12 {
		return ipc.Message{}
	}
	m := ipc.Message{}
	m.Id = atoUint32(b[0:4])
	m.From = ipc.Pid(atoUint32(b[4:8]))
	m.To = ipc.Pid(atoUint32(b[8:12]))
	m.Body = parseBody(b[12:])
	return m
}

func serializeMsg(m ipc.Message) []byte {
	b := []byte{}
	buf := bytes.NewBuffer(b)
	buf.Write(toaUint32(m.Id))
	buf.Write(toaUint32(uint32(m.From)))
	buf.Write(toaUint32(uint32(m.To)))
	buf.Write(serializeBody(m.Body))
	return buf.Bytes()
}