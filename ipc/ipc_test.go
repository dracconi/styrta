package ipc

import (
"testing"
"math/rand"
)

type TMessage struct {}
func (t TMessage) Typ() MsgTyp {
	return 99
}

func TestAdd(t *testing.T) {
	IPCInitialize()
	p, i, o := Make()
	if Add(p, i, o) == true {
		t.Errorf("wanted false")
	}
}

func TestRegistry(t *testing.T) {
	RegistryInitialize()
	r := (Pid)(rand.Uint32())
	if Register(r, "test") {
		t.Errorf("wanted false")
	}

	if pid, ok := Resolve("test"); pid != r || ok == false {
		t.Errorf("mismatch pid")
	}
}

func TestPostman(t *testing.T) {
	IPCInitialize()
	p, i, o := Make()
	Add(p, i, o)
	q, j, a := Make()
	Add(q, j, a)
	if (q == p) {
		t.Errorf("collision pid")
	}
	m := Message{1, p, q, TMessage{}}
	o<-m
	if m != <-j {
		t.Errorf("msg not passed")
	}
}