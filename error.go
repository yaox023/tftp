package tftp

import (
	"bytes"
	"encoding/binary"
	"strings"
)

type ErrorCode uint16

const (
	ErrorCodeNotDefined ErrorCode = iota + 1
	ErrorCodeFileNotFound
	ErrorCodeAccessViolation
	ErrorCodeDiskFull
	ErrorCodeIllegalOperation
	ErrorCodeUnknowTransferID
	ErrorCodeFileExists
	ErrorCodeNoUser
	ErrorCodeInvalidMode
)

type ErrorPacket struct {
	Code ErrorCode
	Msg  string
}

func (e *ErrorPacket) Unmarshal(bs []byte) error {
	buffer := bytes.NewBuffer(bs)

	var opcode Opcode
	err := binary.Read(buffer, binary.BigEndian, &opcode)
	if err != nil {
		return err
	}
	if opcode != OpcodeError {
		return NewInvalidOpcodeError(opcode)
	}

	err = binary.Read(buffer, binary.BigEndian, &e.Code)
	if err != nil {
		return err
	}

	e.Msg, err = buffer.ReadString(0)
	if err != nil {
		return err
	}
	e.Msg = strings.TrimRight(e.Msg, "\x00")
	return nil
}

func (e *ErrorPacket) Marshal() ([]byte, error) {
	// operation code + error code + message + 0-byte
	cap := 2 + 2 + len(e.Msg) + 1
	buffer := new(bytes.Buffer)
	buffer.Grow(cap)

	err := binary.Write(buffer, binary.BigEndian, OpcodeError)
	if err != nil {
		return nil, err
	}

	err = binary.Write(buffer, binary.BigEndian, e.Code)
	if err != nil {
		return nil, err
	}

	_, err = buffer.WriteString(e.Msg)
	if err != nil {
		return nil, err
	}

	err = buffer.WriteByte(0)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
