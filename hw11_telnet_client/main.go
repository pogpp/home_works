package main

import (
	"bufio"
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var duration time.Duration

func init() {
	flag.DurationVar(&duration, "timeout", 10, "default timeout duration")
}

func main() {
	flag.Parse()
	port := os.Args[len(os.Args)-1]
	host := os.Args[len(os.Args)-2]

	in := &bytes.Buffer{}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	defer stop()

	client := NewTelnetClient(host+":"+port, duration, io.NopCloser(in), os.Stdout)
	defer client.Close()

	err := client.Connect()
	if err != nil {
		fmt.Printf("connection failed: %s\n", err)
		return
	}

	var wg sync.WaitGroup

	wg.Add(2)
	go sendByTCP(ctx, stop, &wg, in, client)
	go connectionHandler(ctx, &wg, client)

	wg.Wait()
}

func sendByTCP(ctx context.Context, stop context.CancelFunc, wg *sync.WaitGroup, in *bytes.Buffer, t TelnetClient) {
	go func() {
		<-ctx.Done()
		err := t.Close()
		if err != nil {
			fmt.Printf("failed to colse connection %s", err)
		}
		wg.Done()
	}()

	reader := bufio.NewReader(os.Stdin)
	for {
		resp, err := reader.ReadString('\n')
		if err != nil {
			stop()
			return
		}

		in.WriteString(resp)
		err = t.Send()
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func connectionHandler(ctx context.Context, wg *sync.WaitGroup, t TelnetClient) {
	errCh := make(chan error)
	defer func() {
		wg.Done()
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case errCh <- t.Receive():
			err := <-errCh
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}
