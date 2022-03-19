package tftp

import (
	"errors"
	"strconv"
)

type Opcode uint16

const (
	OpcodeReadRequest Opcode = iota + 1
	OpcodeWriteRequest
	OpcodeData
	OpcodeAcknowledgment
	OpcodeError
)

func NewInvalidOpcodeError(opcode Opcode) error {
	return errors.New("invalid opcode: " + strconv.Itoa(int(opcode)))
}
