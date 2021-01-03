package main

import (
	"flag"
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

	log.Printf("Connected to %s\n", hostPort)

	receiveErrorCh := make(chan struct{}, 1)
	sendErrorCh := make(chan struct{}, 1)
	sigkillCh := make(chan os.Signal, 1)
	signal.Notify(sigkillCh, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		for {
			if err := client.Receive(); err != nil {
				log.Println(err)
				receiveErrorCh <- struct{}{}
				return
			}
		}
	}()

	go func() {
		for {
			if err := client.Send(); err != nil {
				log.Println(err)
				sendErrorCh <- struct{}{}
				return
			}
		}
	}()

	select {
	case <-receiveErrorCh:
	case <-sendErrorCh:
	case <-sigkillCh:
	}

	log.Println("\nBye")
}
