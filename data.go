package tftp

import (
	"bytes"
	"encoding/binary"
	"io"
)

type DataPacket struct {
	Payload io.ReadWriter
	Block   uint16
}

const DatagramSize = 516
const BlockSize = 512

func (d *DataPacket) Marshal() ([]byte, error) {
	buffer := new(bytes.Buffer)

	err := binary.Write(buffer, binary.BigEndian, OpcodeData)
	if err != nil {
		return nil, err
	}

	d.Block += 1
	err = binary.Write(buffer, binary.BigEndian, d.Block)
	if err != nil {
		return nil, err
	}

	_, err = io.CopyN(buffer, d.Payload, BlockSize)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (d *DataPacket) Unmarshal(bs []byte) (int64, error) {
	buffer := bytes.NewBuffer(bs)

	var opcode Opcode
	err := binary.Read(buffer, binary.BigEndian, &opcode)
	if err != nil {
		return 0, err
	}
	if opcode != OpcodeData {
		return 0, NewInvalidOpcodeError(opcode)
	}

	err = binary.Read(buffer, binary.BigEndian, &d.Block)
	if err != nil {
		return 0, err
	}

	n, err := io.CopyN(d.Payload, buffer, BlockSize)
	if err != nil && err != io.EOF {
		return 0, err
	}

	return n, nil
}
