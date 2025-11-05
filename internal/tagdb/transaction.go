package tagdb

import (
	"fmt"
	"sync"

	"dev.azure.com/trayport/Hackathon/_git/Q/internal/logger"
	"github.com/google/uuid"
)

type transaction struct {
	transactionId string
}

type readOnlyTransaction struct {
	transaction
	isOpen bool
	store  *inMemStore
	mu     *sync.RWMutex
}

func newReadOnlyTransaction(store *inMemStore, mu *sync.RWMutex) *readOnlyTransaction {
	mu.RLock()

	id := uuid.NewString()
	logger.Infof("creating read-only transaction %s", id)

	return &readOnlyTransaction{
		transaction: transaction{transactionId: id},
		isOpen:      true,
		store:       store,
		mu:          mu,
	}
}

func (tx *readOnlyTransaction) list(tags []string) ([]TaggedKV, error) {
	if !tx.isOpen {
		err := fmt.Errorf("cannot read from closed transaction %s", tx.transactionId)
		return []TaggedKV{}, err
	}

	return tx.store.list(tags), nil
}

func (tx *readOnlyTransaction) get(key string) (taggedKV TaggedKV, found bool, err error) {
	if !tx.isOpen {
		err := fmt.Errorf("cannot read from closed transaction %s", tx.transactionId)
		return TaggedKV{}, false, err
	}

	taggedKV, found = tx.store.get(key)
	return taggedKV, found, nil
}

func (tx *readOnlyTransaction) close() error {
	// Validation.
	if !tx.isOpen {
		err := fmt.Errorf("cannot close transaction %s", tx.transactionId)
		return err
	}

	logger.Infof("transaction %s closed", tx.transactionId)
	defer tx.mu.RUnlock()
	tx.isOpen = false

	return nil
}

type readWriteTransaction struct {
	transaction
	isOpen     bool
	operations []operator
	store      *inMemStore
	wal        *wal
	mu         *sync.RWMutex
}

func newReadWriteTransaction(store *inMemStore, wal *wal, mu *sync.RWMutex) *readWriteTransaction {
	mu.Lock()

	id := uuid.NewString()
	logger.Infof("creating read-only transaction %s", id)

	return &readWriteTransaction{
		transaction: transaction{transactionId: id},
		isOpen:      true,
		operations:  []operator{},
		store:       store,
		wal:         wal,
		mu:          mu,
	}
}

func (tx *readWriteTransaction) get(key string) (taggedKV TaggedKV, found bool, err error) {
	if !tx.isOpen {
		err := fmt.Errorf("cannot read from closed transaction %s", tx.transactionId)
		return TaggedKV{}, false, err
	}

	taggedKV, found = tx.store.get(key)
	return taggedKV, found, nil
}

func (tx *readWriteTransaction) set(key, value string) {
	// Validation.
	if !tx.isOpen {
		logger.Errorf("cannot update closed transaction %s", tx.transactionId)
	}

	tx.operations = append(tx.operations, &setOperation{
		transactionId: tx.transactionId,
		key:           key,
		value:         value,
	})
}

func (tx *readWriteTransaction) delete(key string) {
	// Validation.
	if !tx.isOpen {
		logger.Errorf("cannot update closed transaction %s", tx.transactionId)
	}

	tx.operations = append(tx.operations, &deleteOperation{
		transactionId: tx.transactionId,
		key:           key,
	})
}

func (tx *readWriteTransaction) tag(key string, tag string) {
	// Validation.
	if !tx.isOpen {
		logger.Errorf("cannot update closed transaction %s", tx.transactionId)
	}

	tx.operations = append(tx.operations, &tagOperation{
		transactionId: tx.transactionId,
		key:           key,
		tag:           tag,
	})
}

func (tx *readWriteTransaction) untag(key string, tag string) {
	// Validation.
	if !tx.isOpen {
		logger.Errorf("cannot update closed transaction %s", tx.transactionId)
	}

	tx.operations = append(tx.operations, &untagOperation{
		transactionId: tx.transactionId,
		key:           key,
		tag:           tag,
	})
}

func (tx *readWriteTransaction) cancel() {
	// Validation.
	if !tx.isOpen {
		return
	}

	logger.Infof("transaction %s cancelled", tx.transactionId)
	defer tx.mu.Unlock()
	tx.isOpen = false
	tx.operations = []operator{}
}

func (tx *readWriteTransaction) commit() error {
	// Validation.
	if !tx.isOpen {
		err := fmt.Errorf("cannot commit closed transaction %s", tx.transactionId)
		return err
	}

	logger.Infof("committing transaction %s", tx.transactionId)
	defer tx.mu.Unlock()
	defer func() { tx.isOpen = false }()
	tx.operations = append(tx.operations, &commitOperation{
		transactionId: tx.transactionId,
	})

	// Write to wal.
	if err := tx.wal.write(tx.operations); err != nil {
		logger.Errorf("failed to write transaction %s to wal because %s", tx.transactionId, err)
		return err
	}

	// Update in-memory store.
	tx.store.apply(tx.operations)

	return nil
}
