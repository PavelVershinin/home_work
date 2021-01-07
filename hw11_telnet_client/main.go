package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var timeout time.Duration

func init() {
	log.SetPrefix("...")
	log.SetFlags(log.Lmsgprefix)
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "--timeout=10s")
}

func main() {
	flag.Parse()

	if flag.NArg() < 2 {
		log.Fatal("You must past the host and port")
	}

	run(
		net.JoinHostPort(flag.Arg(0), flag.Arg(1)),
		timeout,
		ioutil.NopCloser(os.Stdin),
		os.Stdout,
	)
}

func run(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) {
	client := NewTelnetClient(address, timeout, in, out)
	if err := client.Connect(); err != nil {
		log.Fatalf("Connection error: %s", err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			log.Println(err)
		}
	}()

	log.Printf("Connected to %s\n", address)

	errorCh := make(chan error, 1)
	sigkillCh := make(chan os.Signal, 1)
	signal.Notify(sigkillCh, syscall.SIGTERM, syscall.SIGINT)

	go receive(client, errorCh)
	go send(client, errorCh)

	select {
	case <-sigkillCh:
	case err := <-errorCh:
		log.Println(err)
	}
}

func receive(client TelnetClient, errorCh chan error) {
	if err := client.Receive(); err != nil {
		errorCh <- err
	}
}

func send(client TelnetClient, errorCh chan error) {
	if err := client.Send(); err != nil {
		errorCh <- err
	}
	errorCh <- io.EOF
}
