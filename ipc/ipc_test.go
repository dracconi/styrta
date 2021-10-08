package ipc

import (
	"math/rand"
	"testing"
)

type TMessage struct{}

func TestAdd(t *testing.T) {
	Initialize()
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
	Initialize()
	p, i, o := Make()
	Add(p, i, o)
	q, j, a := Make()
	Add(q, j, a)
	if q == p {
		t.Errorf("collision pid")
	}
	m := Message{1, p, q, TMessage{}}
	o <- m
	if m != <-j {
		t.Errorf("msg not passed")
	}

	o <- Message{0, 0, 0, MsgDie{}}
	if (Message{0, 0, 0, MsgDie{}}) != <-i {
		t.Errorf("death not caught")
	}
}
