package tagdb

import (
	"errors"
	"os"
	"path"
	"strconv"
	"strings"

	"dev.azure.com/trayport/Hackathon/_git/Q/internal/logger"
)

const (
	walFileExtension = ".wal"
)

type walManager struct {
	walRoot   string
	currentId int64
	walFiles  map[int64]*wal
}

func newWalManager(walRoot string) (*walManager, error) {
	err := createDirIfNotExists(walRoot)
	if err != nil {
		innerErr := logger.Error("failed to get or create wal directory")
		return nil, errors.Join(err, innerErr)
	}

	wals, currentId, err := openWals(walRoot)
	if err != nil {
		innerErr := logger.Error("failed to list wal files")
		return nil, errors.Join(err, innerErr)
	}

	// Ensure there is at least one wal file.
	if len(wals) == 0 {
		walPath := path.Join(walRoot, "0"+walFileExtension)
		wal, err := openWal(0, walPath)
		if err != nil {
			innerErr := logger.Error("failed to create initial wal file")
			return nil, errors.Join(err, innerErr)
		}

		wals[0] = wal
		currentId = 0
	}

	return &walManager{
		walRoot:   walRoot,
		currentId: currentId,
		walFiles:  wals,
	}, nil
}

func (wm *walManager) close() error {
	var err error

	logger.Info("closing wal manager")
	for _, wal := range wm.walFiles {
		closeErr := wal.close()
		if closeErr != nil {
			logger.Errorf("failed to close wal file %d because %s", wal.id, closeErr)
			err = errors.Join(err, closeErr)
		}
	}

	return err
}

func (wm *walManager) current() *wal {
	current, found := wm.walFiles[wm.currentId]
	if !found {
		logger.Panicf("wal file %d not found", wm.currentId)
	}

	return current
}

func (wm *walManager) shouldRoll(rollWalAfterBytes int64) bool {
	info, err := wm.current().file.Stat()
	if err != nil {
		logger.Warn("unable to state wal file")
		return false
	}

	return info.Size() > rollWalAfterBytes
}

func (wm *walManager) roll() {
	nextId := wm.currentId + 1
	nextIdStr := strconv.FormatInt(nextId, 10)
	walPath := path.Join(wm.walRoot, nextIdStr+walFileExtension)
	wal, err := openWal(nextId, walPath)
	if err != nil {
		logger.Warnf("failed to create new wal file `%s` because `%s`", walPath, err)
		return
	}

	wm.walFiles[nextId] = wal
	wm.currentId++
	logger.Infof("rolled wal file to %d", wm.currentId)
}

func openWals(walRoot string) (files map[int64]*wal, currentId int64, err error) {
	result := map[int64]*wal{}
	maxId := int64(-1)

	dirEntries, err := os.ReadDir(walRoot)
	if err != nil {
		logger.Errorf("failed to read wal directory `%s` because `%s`", walRoot, err)
		return result, maxId, err
	}

	for _, entry := range dirEntries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), walFileExtension) {
			continue
		}

		strId := strings.TrimSuffix(entry.Name(), walFileExtension)
		id, err := strconv.ParseInt(strId, 10, 64)
		if err != nil {
			logger.Warnf("skipping wal file with invalid name `%s`", entry.Name())
			continue
		}

		walPath := path.Join(walRoot, entry.Name())
		wal, err := openWal(id, walPath)
		if err != nil {
			logger.Errorf("failed to open wal file `%s` because `%s`", walPath, err)
			return result, maxId, err
		}

		result[id] = wal
		if id > maxId {
			maxId = id
		}
	}

	return result, maxId, nil
}
