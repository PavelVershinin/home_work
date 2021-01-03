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
	log.SetFlags(log.Lshortfile | log.Ltime)
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "--timeout=10s")
}

func main() {
	flag.Parse()

	if flag.NArg() < 2 {
		log.Fatal("You must past the host and port")
	}

	hostPort := net.JoinHostPort(flag.Arg(0), flag.Arg(1))
	client := NewTelnetClient(hostPort, timeout, ioutil.NopCloser(os.Stdin), os.Stdout)
	if err := client.Connect(); err != nil {
		log.Fatalf("Connection error: %s", err)
	}
	defer client.Close()

	log.Printf("...Connected to %s\n", hostPort)

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
