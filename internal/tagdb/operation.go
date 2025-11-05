package tagdb

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type operator interface {
	serialize() []byte
	getTransactionId() string
}

type operationCode int

const (
	opCodeSet operationCode = iota
	opCodeDelete
	opCodeTag
	opCodeUntag
	opCodeCommit
)

func (op operationCode) String() string {
	switch op {
	case opCodeSet:
		return "SET"
	case opCodeDelete:
		return "DELETE"
	case opCodeTag:
		return "TAG"
	case opCodeUntag:
		return "UNTAG"
	case opCodeCommit:
		return "COMMIT"
	default:
		panic(fmt.Sprintf("unsupported operation code %d", op))
	}
}

const (
	opFieldSeparator           = "\x1F" // Unit Separator ASCII character.
	opRecordSeparator          = "\x1E" // Record Separator ASCII character.
	opRecordSeparatorByte byte = 30
)

type setOperation struct {
	transactionId string
	key           string
	value         string
}

func (op setOperation) serialize() []byte {
	fields := []string{op.transactionId, opCodeSet.String(), op.key, op.value}
	record := strings.Join(fields, opFieldSeparator) + opRecordSeparator
	return []byte(record)
}

func (op setOperation) getTransactionId() string {
	return op.transactionId
}

type deleteOperation struct {
	transactionId string
	key           string
}

func (op deleteOperation) serialize() []byte {
	fields := []string{op.transactionId, opCodeDelete.String(), op.key}
	record := strings.Join(fields, opFieldSeparator) + opRecordSeparator
	return []byte(record)
}

func (op deleteOperation) getTransactionId() string {
	return op.transactionId
}

type tagOperation struct {
	transactionId string
	key           string
	tag           string
}

func (op tagOperation) serialize() []byte {
	fields := []string{op.transactionId, opCodeTag.String(), op.key, op.tag}
	record := strings.Join(fields, opFieldSeparator) + opRecordSeparator
	return []byte(record)
}

func (op tagOperation) getTransactionId() string {
	return op.transactionId
}

type untagOperation struct {
	transactionId string
	key           string
	tag           string
}

func (op untagOperation) serialize() []byte {
	fields := []string{op.transactionId, opCodeUntag.String(), op.key, op.tag}
	record := strings.Join(fields, opFieldSeparator) + opRecordSeparator
	return []byte(record)
}

func (op untagOperation) getTransactionId() string {
	return op.transactionId
}

type commitOperation struct {
	transactionId string
}

func (op commitOperation) serialize() []byte {
	fields := []string{op.transactionId, opCodeCommit.String()}
	record := strings.Join(fields, opFieldSeparator) + opRecordSeparator
	return []byte(record)
}

func (op commitOperation) getTransactionId() string {
	return op.transactionId
}

func deserialize(data []byte) (operator, error) {
	// Remove optional trailing record separator.
	if len(data) > 0 && data[len(data)-1] == opRecordSeparatorByte {
		data = data[0 : len(data)-1]
	}

	// Convert to text early, for clearer log messages later.
	record := string(data)
	fields := strings.Split(record, opFieldSeparator)

	const txField = 0     // Required for all op codes.
	const opCodeField = 1 // Required for all op codes.
	const keyField = 2
	const valueField = 3 // Mutually exclusive with tagField.
	const tagField = 3   // Mutually exclusive with valueField.

	// Validation.
	if len(fields) < 2 {
		return nil, fmt.Errorf("cannot deserialize corrupted operation record: %s", record)
	}

	if err := uuid.Validate(fields[txField]); err != nil {
		return nil, fmt.Errorf("cannot deserialize unsupported transaction id: %s", fields[txField])
	}

	var opCode operationCode
	var expectedFieldCount int
	switch fields[opCodeField] {
	case opCodeSet.String():
		opCode = opCodeSet
		expectedFieldCount = 4
	case opCodeDelete.String():
		opCode = opCodeDelete
		expectedFieldCount = 3
	case opCodeTag.String():
		opCode = opCodeTag
		expectedFieldCount = 4
	case opCodeUntag.String():
		opCode = opCodeUntag
		expectedFieldCount = 4
	case opCodeCommit.String():
		opCode = opCodeCommit
		expectedFieldCount = 2
	default:
		return nil, fmt.Errorf("cannot deserialize unsupported operation code: %s", fields[opCodeField])
	}

	if len(fields) != expectedFieldCount {
		return nil, fmt.Errorf(
			"unexpected number of fields for %s operation, expected %d found %d in record: %s",
			opCode.String(),
			expectedFieldCount,
			len(fields),
			record)
	}

	// Deserialize.
	switch opCode {
	case opCodeSet:
		return &setOperation{
			transactionId: fields[txField],
			key:           fields[keyField],
			value:         fields[valueField],
		}, nil
	case opCodeDelete:
		return &deleteOperation{
			transactionId: fields[txField],
			key:           fields[keyField],
		}, nil
	case opCodeTag:
		return &tagOperation{
			transactionId: fields[txField],
			key:           fields[keyField],
			tag:           fields[tagField],
		}, nil
	case opCodeUntag:
		return &untagOperation{
			transactionId: fields[txField],
			key:           fields[keyField],
			tag:           fields[tagField],
		}, nil
	case opCodeCommit:
		return &commitOperation{
			transactionId: fields[txField],
		}, nil
	default:
		return nil, fmt.Errorf("cannot deserialize due to unsupported op code %d", opCode)
	}
}

func opSplit(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF {
		return 0, nil, nil
	}

	if i := bytes.IndexByte(data, opRecordSeparatorByte); i >= 0 {
		return i + 1, data[0 : i+1], nil
	}

	return 0, nil, nil
}
