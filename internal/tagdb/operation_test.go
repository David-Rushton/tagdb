package tagdb

import (
	"reflect"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func Test_operator_RoundTrips(t *testing.T) {
	txId := uuid.NewString()
	testCases := []operator{
		&setOperation{txId, "key1", "value1"},
		&deleteOperation{txId, "key2"},
		&tagOperation{txId, "key3", "tag1"},
		&untagOperation{txId, "key4", "tag2"},
		&commitOperation{txId},
	}

	for _, expected := range testCases {
		b := expected.serialize()
		actual, err := deserialize(b)
		if err != nil {
			t.Errorf("unexpected error during deserialization: %v", err)
			continue
		}

		if !reflect.DeepEqual(expected, actual) {
			t.Errorf(
				"operation failed to round trip:\n\texpected: %+v\n\tactual:   %+v",
				expected,
				actual)
		}
	}
}

func Test_operator_ShouldError_OnInvalidOperator(t *testing.T) {
	nonsenseRecord := []byte("gibberish")
	_, err := deserialize(nonsenseRecord)
	if err == nil {
		t.Errorf("expected error on invalid operator, but got none")
	}
}

func Test_operator_ShouldError_OnInvalidTransactionId(t *testing.T) {
	txId := "invalid-uuid"
	op := &setOperation{txId, "key1", "value1"}
	data := op.serialize()

	_, err := deserialize(data)
	if err == nil {
		t.Errorf("expected error on invalid operator, but got none")
	}
}

func Test_operator_ShouldError_OnInvalidOpCode(t *testing.T) {
	txId := uuid.NewString()
	fields := []string{txId, "INVALID_OP_CODE", "key1", "value1"}
	invalidRecord := strings.Join(fields, opFieldSeparator) + opRecordSeparator
	data := []byte(invalidRecord)

	_, err := deserialize(data)
	if err == nil {
		t.Errorf("expected error on invalid operator, but got none")
	}
}

func Test_operator_ShouldError_OnInvalidRecord(t *testing.T) {
	txId := uuid.NewString()
	fields := []string{txId, opCodeSet.String(), "key1", "value1", "extra_field"}
	invalidRecord := strings.Join(fields, opFieldSeparator) + opRecordSeparator
	data := []byte(invalidRecord)

	_, err := deserialize(data)
	if err == nil {
		t.Errorf("expected error on invalid operator, but got none")
	}
}
