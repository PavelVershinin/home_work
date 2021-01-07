package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	t.Run("log test", func(t *testing.T) {
		wg := sync.WaitGroup{}
		stdOut := &bytes.Buffer{}
		log.SetOutput(stdOut)
		log.SetPrefix("...")
		log.SetFlags(log.Lmsgprefix)

		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		pReader, pWriter := io.Pipe()

		wg.Add(2)
		go func() {
			run(
				l.Addr().String(),
				time.Second*10,
				pReader,
				nil,
			)
			wg.Done()
		}()

		go func() {
			require.NoError(t, pWriter.Close())
			wg.Done()
		}()

		wg.Wait()

		expected := fmt.Sprintf("...Connected to %s\n...EOF\n", l.Addr().String())
		require.Equal(t, expected, stdOut.String())
	})
}
