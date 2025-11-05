package tagdb

import (
	"path"
	"reflect"
	"testing"

	"github.com/google/uuid"
)

const (
	testWalId = 0
)

func Test_RoundTrips(t *testing.T) {
	// Arrange
	txId := uuid.NewString()
	expected := []operator{
		&setOperation{transactionId: txId, key: "key-1", value: "value-1"},
		&setOperation{transactionId: txId, key: "key-2", value: "value-2"},
		&setOperation{transactionId: txId, key: "key-3", value: "value-3"},
		&deleteOperation{transactionId: txId, key: "key-3"},
		&tagOperation{transactionId: txId, key: "key-2", tag: "tag-1"},
		&tagOperation{transactionId: txId, key: "key-2", tag: "tag-2"},
		&tagOperation{transactionId: txId, key: "key-2", tag: "tag-3"},
		&untagOperation{transactionId: txId, key: "key-2", tag: "tag-3"},
		&commitOperation{transactionId: txId},
	}

	path := path.Join(t.TempDir(), "test.wal")
	wal, err := openWal(testWalId, path)
	if err != nil {
		t.Fatalf("unexpected error connecting to wal: %v", err)
	}
	defer wal.close()

	// Act
	wal.write(expected)
	actual, err := wal.read()
	if err != nil {
		t.Fatalf("unexpected error reading wal: %v", err)
	}

	// Assert
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf(
			"wal operation failed to round trip:\n\texpected: %+v\n\tactual:   %+v",
			expected,
			actual)
	}
}

// TODO: Fix test.
// OPTION 1: In-mem apply only updates on commit.
// OPTION 2: Wal only yields committed transactions.
func Test_read_ReturnsCommittedTransactionsOnly(t *testing.T) {
	// Arrange
	tx1 := uuid.NewString()
	tx2 := uuid.NewString()
	committedTx := []operator{
		&setOperation{transactionId: tx1, key: "key-1", value: "value-1"},
		&setOperation{transactionId: tx1, key: "key-2", value: "value-2"},
		&commitOperation{transactionId: tx1},
	}
	uncommittedTx := []operator{
		&setOperation{transactionId: tx2, key: "key-3", value: "value-3"},
		&setOperation{transactionId: tx2, key: "key-4", value: "value-4"},
	}
	operations := append(committedTx, uncommittedTx...)
	expected := committedTx

	path := path.Join(t.TempDir(), "test.wal")
	wal, err := openWal(testWalId, path)
	if err != nil {
		t.Fatalf("unexpected error connecting to wal: %v", err)
	}
	defer wal.close()

	// Act
	wal.write(operations)
	actual, err := wal.read()
	if err != nil {
		t.Fatalf("unexpected error reading wal: %v", err)
	}

	// Assert
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf(
			"wal operation failed to round trip:\n\texpected: %+v\n\tactual:   %+v",
			expected,
			actual)
	}
}
