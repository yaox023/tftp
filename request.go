package tftp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"strings"
)

type RequestPacket struct {
	Mode     string
	FileName string
	Opcode   Opcode
}

func (r *RequestPacket) Unmarshal(bs []byte) error {
	buffer := bytes.NewBuffer(bs)

	err := binary.Read(buffer, binary.BigEndian, &r.Opcode)
	if err != nil {
		return err
	}

	if r.Opcode != OpcodeReadRequest && r.Opcode != OpcodeWriteRequest {
		return NewInvalidOpcodeError(r.Opcode)
	}

	r.FileName, err = buffer.ReadString(0)
	if err != nil {
		return err
	}
	r.FileName = strings.TrimRight(r.FileName, "\x00")
	if len(r.FileName) == 0 {
		return errors.New("empty filename")
	}

	r.Mode, err = buffer.ReadString(0)
	if err != nil {
		return err
	}
	r.Mode = strings.TrimRight(r.Mode, "\x00")
	r.Mode = strings.ToLower(r.Mode)
	if r.Mode != ModeOctet && r.Mode != ModeNetascii && r.Mode != ModeMail {
		return errors.New("invalid Mode: " + r.Mode)
	}

	return nil
}

func (r *RequestPacket) Marshal() ([]byte, error) {
	buffer := new(bytes.Buffer)

	err := binary.Write(buffer, binary.BigEndian, r.Opcode)
	if err != nil {
		return nil, err
	}

	_, err = buffer.WriteString(r.FileName)
	if err != nil {
		return nil, err
	}

	err = buffer.WriteByte(0)
	if err != nil {
		return nil, err
	}

	_, err = buffer.WriteString(r.Mode)
	if err != nil {
		return nil, err
	}
	err = buffer.WriteByte(0)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
