package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type ConnectionParams struct {
	Address string
	Timeout time.Duration
	In      io.ReadCloser
	Out     io.Writer
	Conn    net.Conn
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &ConnectionParams{
		Address: address,
		Timeout: timeout,
		In:      in,
		Out:     out,
	}
}

func (c *ConnectionParams) Connect() error {
	conn, err := net.DialTimeout("tcp", c.Address, c.Timeout)
	if err != nil {
		return fmt.Errorf("connection error: %w", err)
	}

	c.Conn = conn
	return nil
}

func (c *ConnectionParams) Close() error {
	if c.Conn != nil {
		if err := c.Conn.Close(); err != nil {
			return fmt.Errorf("close error: %w", err)
		}
	}

	return nil
}

func (c *ConnectionParams) Send() error {
	_, err := io.Copy(c.Conn, c.In)
	if err != nil {
		return fmt.Errorf("send error: %w", err)
	}

	return nil
}

func (c *ConnectionParams) Receive() error {
	_, err := io.Copy(c.Out, c.Conn)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return fmt.Errorf("receive error: %w", err)
	}

	return nil
}
