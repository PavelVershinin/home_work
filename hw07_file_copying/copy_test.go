package main

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	sep     = string(os.PathSeparator)
	testDir = "testdata"
)

func TestCopy(t *testing.T) {
	sourceFiles, err := ioutil.ReadDir(testDir)
	require.NoError(t, err)

	tmp, err := ioutil.TempFile(os.TempDir(), "*.tmp")
	require.NoError(t, err)

	defer func() {
		if err := os.Remove(tmp.Name()); err != nil {
			log.Println(err)
		}
	}()

	require.NoError(t, tmp.Close())

	t.Run("copy empty file", func(t *testing.T) {
		err = Copy(testDir+sep+"input_empty.txt", tmp.Name(), 0, 10000)
		require.EqualError(t, err, ErrUnsupportedFile.Error())
	})

	t.Run("offset exceeds file size", func(t *testing.T) {
		for _, f := range sourceFiles {
			if !f.IsDir() && f.Name() != "." && f.Name() != ".." && f.Size() > 0 {
				err := Copy(testDir+sep+f.Name(), tmp.Name(), f.Size()+1, f.Size())
				require.EqualError(t, err, ErrOffsetExceedsFileSize.Error())
			}
		}
	})

	t.Run("offset equal file size", func(t *testing.T) {
		for _, f := range sourceFiles {
			if !f.IsDir() && f.Name() != "." && f.Name() != ".." && f.Size() > 0 {
				err := Copy(testDir+sep+f.Name(), tmp.Name(), f.Size(), f.Size())
				require.NoError(t, err)
			}
		}
	})

	t.Run("offset equal file size minus one byte", func(t *testing.T) {
		for _, f := range sourceFiles {
			if !f.IsDir() && f.Name() != "." && f.Name() != ".." && f.Size() > 0 {
				err := Copy(testDir+sep+f.Name(), tmp.Name(), f.Size()-1, f.Size())
				require.NoError(t, err)

				b, err := ioutil.ReadFile(tmp.Name())
				require.NoError(t, err)
				require.Len(t, b, 1)
			}
		}
	})
}
