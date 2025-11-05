package tagdb

import (
	"errors"
	"fmt"
	"path"
	"slices"
	"sync"

	"dev.azure.com/trayport/Hackathon/_git/Q/internal/logger"
)

// Write-ahead log.
type storage struct {
	root       string
	walDir     string
	inMemStore *inMemStore
	walManager *walManager
	mu         sync.RWMutex
}

func openStorage(root string) (*storage, error) {
	// Ensure wal dir exists.
	walDir := path.Join(root, "wal")
	if err := createDirIfNotExists(walDir); err != nil {
		innerErr := logger.Error("cannot open wal directory")
		return nil, errors.Join(innerErr, err)
	}

	// Get current wal file name.
	walManager, err := newWalManager(walDir)
	if err != nil {
		innerErr := logger.Error("cannot create wal manager")
		return nil, errors.Join(innerErr, err)
	}

	// Rehydrate in-mem store from wals.
	var operations []operator
	for i := range walManager.currentId + 1 {
		wal := walManager.walFiles[int64(i)]
		walOps, err := wal.read()
		if err != nil {
			innerErr := logger.Error("cannot read wal operations")
			return nil, errors.Join(innerErr, err)
		}

		operations = append(operations, walOps...)
	}

	inMemStore := newInMemStore()
	inMemStore.apply(operations)

	// Create and return storage connection.
	storageConnection := &storage{
		root:       root,
		walDir:     walDir,
		inMemStore: inMemStore,
		walManager: walManager,
		mu:         sync.RWMutex{},
	}

	return storageConnection, nil
}

func (w *storage) close() error {
	logger.Info("closing storage connection")

	return w.walManager.close()
}

func (s *storage) list(tags []string) ([]TaggedKV, error) {
	tx := newReadOnlyTransaction(s.inMemStore, &s.mu)
	defer tx.close()

	return tx.list(tags)
}

func (s *storage) get(key string) (taggedKV TaggedKV, found bool, err error) {
	tx := newReadOnlyTransaction(s.inMemStore, &s.mu)
	defer tx.close()

	return tx.get(key)
}

func (s *storage) set(key, value string) error {
	tx := newReadWriteTransaction(s.inMemStore, s.walManager.current(), &s.mu)
	tx.set(key, value)

	return tx.commit()
}

func (s *storage) delete(key string) error {
	tx := newReadWriteTransaction(s.inMemStore, s.walManager.current(), &s.mu)

	old, found, err := tx.get(key)
	if err != nil {
		tx.cancel()
		return err
	}

	if !found {
		tx.cancel()
		return fmt.Errorf("key not found `%s` ", key)
	}

	for _, tag := range old.Tags {
		tx.untag(key, tag)
	}

	tx.delete(key)

	return tx.commit()
}

func (s *storage) tag(key, tag string) error {
	tx := newReadWriteTransaction(s.inMemStore, s.walManager.current(), &s.mu)
	defer tx.cancel()

	taggedKV, found, err := tx.get(key)
	if err != nil {
		return err
	}

	if !found {
		return fmt.Errorf("key not found `%s` ", key)
	}
	if slices.Contains(taggedKV.Tags, tag) {
		return fmt.Errorf("Tag `%s` already exists on key `%s`", tag, key)
	}

	tx.tag(key, tag)

	return tx.commit()
}

func (s *storage) untag(key, tag string) error {
	tx := newReadWriteTransaction(s.inMemStore, s.walManager.current(), &s.mu)
	defer tx.cancel()

	taggedKV, found, err := tx.get(key)
	if err != nil {
		return err
	}

	if !found {
		return fmt.Errorf("key not found `%s` ", key)
	}

	if !slices.Contains(taggedKV.Tags, tag) {
		return fmt.Errorf("Tag `%s` not found on key `%s`", tag, key)
	}

	tx.untag(key, tag)

	return tx.commit()
}

func (s *storage) maybeRoll(rollWalAfterBytes int64) {
	if s.walManager.shouldRoll(rollWalAfterBytes) {
		tx := newReadWriteTransaction(s.inMemStore, s.walManager.current(), &s.mu)
		defer tx.commit()

		logger.Info("rolling wal")
		s.walManager.roll()
	}
}

func (s *storage) maybeCompact() {
	panic("not implemented")
}
