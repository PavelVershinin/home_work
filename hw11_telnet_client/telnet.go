package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

var ErrorConnectionIsNil = errors.New("connection is nil")

type TelnetClient interface {
	Connect() error
	Send() error
	Receive() error
	Close() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &Client{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

type Client struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func (c *Client) Connect() error {
	var err error
	c.conn, err = net.DialTimeout("tcp", c.address, c.timeout)
	return err
}
func (c Client) Send() error {
	if _, err := io.Copy(c.conn, c.in); err != nil {
		return fmt.Errorf("can not receive: %w", err)
	}
	return nil
}
func (c Client) Receive() error {
	if _, err := io.Copy(c.out, c.conn); err != nil {
		return fmt.Errorf("can not receive: %w", err)
	}
	return nil
}
func (c *Client) Close() error {
	if c.conn == nil {
		return ErrorConnectionIsNil
	}
	return c.conn.Close()
}
