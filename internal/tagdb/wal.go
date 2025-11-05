package tagdb

import (
	"bufio"
	"errors"
	"io"
	"os"

	"dev.azure.com/trayport/Hackathon/_git/Q/internal/logger"
)

// Write-ahead log.
// TODO: Add mock fs support for testing.
type wal struct {
	id   int64
	file *os.File
	rw   *bufio.ReadWriter
}

func openWal(id int64, path string) (*wal, error) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		logger.Errorf("failed to open wal file `%s` because `%s`", path, err)
		return nil, err
	}
	logger.Infof("opened wal file `%s`", path)

	reader := bufio.NewReader(file)
	writer := bufio.NewWriter(file)
	readWriter := bufio.NewReadWriter(reader, writer)

	wal := &wal{
		id:   id,
		file: file,
		rw:   readWriter,
	}

	return wal, nil
}

func (w *wal) flush() {
	logger.Info("flushing wal")
	if err := w.rw.Flush(); err != nil {
		logger.Warnf("failed to flush wal: %s", err)
	}
}

func (w *wal) close() error {
	logger.Info("closing wal connection")
	flushErr := w.rw.Flush()
	closeErr := w.file.Close()

	err := errors.Join(flushErr, closeErr)
	if err != nil {
		logger.Errorf("could not cleanly close wal because %s", err)
	}

	return err
}

func (w *wal) read() ([]operator, error) {
	// Move to start of file.
	if _, err := w.file.Seek(0, io.SeekStart); err != nil {
		logger.Errorf("failed to seek to start of wal file: %s", err)
		return nil, err
	}
	defer func() {
		// Return to end of file.
		if _, err := w.file.Seek(0, io.SeekEnd); err != nil {
			logger.Panicf("failed to seek to end of wal file: %s", err)
		}
	}()

	// Read contents.
	committedTx := map[string]bool{}
	buf := []operator{}
	scanner := bufio.NewScanner(w.rw)
	scanner.Split(opSplit)
	for scanner.Scan() {
		data := scanner.Bytes()
		op, err := deserialize(data)
		if err != nil {
			logger.Errorf("failed to deserialize wal record: %s", err)
			return nil, err
		}

		if _, isCommit := op.(*commitOperation); isCommit {
			committedTx[op.(*commitOperation).transactionId] = true
		}

		buf = append(buf, op)
	}

	if scanner.Err() != nil {
		logger.Errorf("failed to scan wal file: %s", scanner.Err())
		return nil, scanner.Err()
	}

	// Filter for committed transactions only.
	var result []operator
	for _, op := range buf {
		if committedTx[op.getTransactionId()] {
			result = append(result, op)
		}
	}

	return result, nil
}

func (w *wal) write(ops []operator) error {
	defer w.flush()

	logger.Infof("writing %d operation(s) to wal", len(ops))
	for _, op := range ops {
		data := op.serialize()
		if _, err := w.rw.Write(data); err != nil {
			return err
		}
	}

	return nil
}
