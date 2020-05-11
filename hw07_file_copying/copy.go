package main

import (
	"bytes"
	"errors"
	"io"
	"log"
	"math"
	"os"

	"github.com/PavelVershinin/home_work/hw07_file_copying/progressbar"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath string, toPath string, offset, limit int64) error {
	var sourceFile, destinationFile, err = open(fromPath, toPath)
	if err != nil {
		return err
	}
	defer func() {
		if err := sourceFile.Close(); err != nil {
			log.Println(err)
		}
		if err := destinationFile.Close(); err != nil {
			log.Println(err)
		}
	}()

	sourceSize, err := size(sourceFile)
	if err != nil {
		return err
	}

	if sourceSize == 0 {
		return ErrUnsupportedFile
	}

	if offset > sourceSize {
		return ErrOffsetExceedsFileSize
	}

	if limit <= 0 {
		limit = sourceSize - offset
	} else {
		limit = int64(math.Min(float64(limit), float64(sourceSize-offset)))
	}

	return cp(sourceFile, destinationFile, limit)
}

func open(sourcePath, destinationPath string) (source, destination *os.File, err error) {
	source, err = os.Open(sourcePath)
	if err != nil {
		return nil, nil, err
	}

	if _, err := source.Seek(offset, 0); err != nil {
		return nil, nil, err
	}

	destination, err = os.OpenFile(destinationPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, nil, err
	}

	if err := destination.Truncate(0); err != nil {
		return nil, nil, err
	}

	return
}

func size(file *os.File) (int64, error) {
	info, err := file.Stat()
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

func cp(source io.Reader, destination io.Writer, limit int64) error {
	pBar := progressbar.New().
		Min(float64(offset)).
		Max(float64(offset + limit)).
		Val(float64(offset))

	defer func() {
		if err := pBar.Close(); err != nil {
			log.Println(err)
		}
	}()

	var buffCap int64 = bytes.MinRead
	if limit < buffCap {
		buffCap = limit
	}
	buff := make([]byte, buffCap)
	for {
		n, err := source.Read(buff)
		if err != nil && err != io.EOF {
			return err
		} else if err != nil && err == io.EOF {
			break
		}
		if len(buff) > n {
			buff = buff[:n]
		}
		if limit < int64(n) {
			buff = buff[:limit]
		}
		if _, err := destination.Write(buff); err != nil {
			return err
		}
		if err = pBar.Add(float64(len(buff))).Draw("Completed :percent%; Left: :left bytes"); err != nil {
			log.Println(err)
		}
		limit -= int64(n)
		if limit <= 0 {
			break
		}
	}

	return nil
}
