package tftp

import (
	"testing"
)

func TestRequest_Marshal(t *testing.T) {
	req := RequestPacket{}
	bs := []byte{0, 1, 49, 50, 51, 0, 110, 101, 116, 97, 115, 99, 105, 105, 0}
	err := req.Unmarshal(bs)
	if err != nil {
		t.Fatal(err)
	}
	if req.Mode != ModeNetascii {
		t.Fatal(req.Mode)
	}
	if req.FileName != "123" {
		t.Fatal(req.FileName)
	}
	if req.Opcode != OpcodeReadRequest {
		t.Fatal(req.Opcode)
	}

}
