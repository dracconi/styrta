package network

import ( "testing"
"github.com/dracconi/styrta/ipc"
)

func TestParsing(t *testing.T) {
	mb := MsgAdd{"zombo.com"}
	if parseBody(serializeBody(mb)) != mb {
		t.Errorf("ident false")
	}

	m := ipc.Message{0, 0, 0, mb}

	if parseMsg(serializeMsg(m)) != m {
		t.Errorf("ident false")
	}
}