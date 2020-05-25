package main

import (
	"bufio"
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

const (
	defaultBufSize int64 = 1024
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

//nolint:funlen
func Copy(fromPath string, toPath string, offset, limit int64) error {
	// Read file
	rf, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer rf.Close()

	fi, err := rf.Stat()
	if err != nil {
		return err
	}

	if !fi.Mode().IsRegular() {
		return ErrUnsupportedFile
	}
	if fi.Size() < offset {
		return ErrOffsetExceedsFileSize
	}

	if _, err := rf.Seek(offset, io.SeekStart); err != nil {
		return err
	}

	// Write file
	wf, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer wf.Close()

	// Buffers
	bufSize := defaultBufSize
	if limit != 0 && limit < bufSize {
		bufSize = limit
	}

	br := bufio.NewReaderSize(rf, int(bufSize))
	bw := bufio.NewWriterSize(wf, int(bufSize))

	// Progress bar
	barTotal := fi.Size() - offset
	if limit != 0 && limit < barTotal {
		barTotal = limit
	}

	bar := pb.StartNew(int(barTotal))
	defer bar.Finish()

	// Copy files
	summaryWritten := int64(0)
	for {
		written, err := io.CopyN(bw, br, bufSize)

		bar.Add(int(written))

		summaryWritten += written
		if summaryWritten == limit {
			break
		}

		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}
	}

	bw.Flush()

	return nil
}
