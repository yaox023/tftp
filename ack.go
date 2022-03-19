package tftp

import (
	"bytes"
	"encoding/binary"
)

type AckPacket struct {
	Block uint16
}

func (a *AckPacket) Unmarshal(bs []byte) error {
	buffer := bytes.NewBuffer(bs)

	var opcode Opcode
	err := binary.Read(buffer, binary.BigEndian, &opcode)
	if err != nil {
		return err
	}

	if opcode != OpcodeAcknowledgment {
		return NewInvalidOpcodeError(opcode)
	}

	err = binary.Read(buffer, binary.BigEndian, &a.Block)
	if err != nil {
		return err
	}

	return nil
}

func (a *AckPacket) Marshal() ([]byte, error) {
	buffer := new(bytes.Buffer)

	err := binary.Write(buffer, binary.BigEndian, OpcodeAcknowledgment)
	if err != nil {
		return nil, err
	}

	err = binary.Write(buffer, binary.BigEndian, a.Block)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
