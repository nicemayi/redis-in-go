package tcp

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"redis-in-go/lib/logger"
	"redis-in-go/lib/sync/atomic"
	"redis-in-go/lib/sync/wait"
	"sync"
	"time"
)

type EchoClient struct {
	Conn    net.Conn
	Waiting wait.Wait
}

func (e *EchoClient) Close() error {
	e.Waiting.WaitWithTimeout(10 * time.Second)
	e.Conn.Close()

	return nil
}

type EchoHandler struct {
	activeConn sync.Map
	closing    atomic.Boolean
}

func MakeHandler() *EchoHandler {
	return &EchoHandler{}
}

func (handler *EchoHandler) Handle(ctx context.Context, conn net.Conn) {
	if handler.closing.Get() {
		conn.Close()
	}

	client := &EchoClient{
		Conn: conn,
	}

	handler.activeConn.Store(client, struct{}{})
	reader := bufio.NewReader(conn)

	for {
		msg, err := reader.ReadString('\n')

		if err != nil {
			if err == io.EOF {
				logger.Info(fmt.Sprintf("Connection close with msg: %s\n", err.Error()))
				handler.activeConn.Delete(client)
			} else {
				logger.Warn(err)
			}

			return
		}

		logger.Info(fmt.Sprintf("received mssage: %s\n", msg))

		client.Waiting.Add(1)
		b := []byte(msg)
		conn.Write(b)
		client.Waiting.Done()
	}
}

func (handler *EchoHandler) Close() error {
	logger.Info("handler shutting down!")
	handler.closing.Set(true)

	handler.activeConn.Range(func(key, value interface{}) bool {
		client := key.(*EchoClient)
		client.Conn.Close()
		return true
	})

	return nil
}
