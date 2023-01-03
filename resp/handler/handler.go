package handler

import (
	"context"
	"io"
	"net"
	"redis-in-go/database"
	databaseface "redis-in-go/interface/database"
	"redis-in-go/interface/resp"
	"redis-in-go/lib/logger"
	"redis-in-go/lib/sync/atomic"
	"redis-in-go/resp/connection"
	"redis-in-go/resp/parser"
	"redis-in-go/resp/reply"
	"strings"
	"sync"
)

var (
	unknownErrReplyBytes = []byte("-ERR unknown\r\n")
)

func (r *RespHandler) Exec(client resp.Connection, args [][]byte) resp.Reply {
	panic("not implemented") // TODO: Implement
}

func (r *RespHandler) Close() error {
	logger.Info("handler shutting down")
	r.closing.Set(true)
	r.activeConn.Range(
		func(key interface{}, value interface{}) bool {
			client := key.(*connection.Connection)
			client.Close()
			return true
		},
	)
	r.db.Close()
	return nil
}

func (r *RespHandler) closeClient(client *connection.Connection) {
	client.Close()
	r.db.AfterClientClose(client)
	r.activeConn.Delete(client)
}

func (r *RespHandler) Handle(ctx context.Context, conn net.Conn) {
	if r.closing.Get() {
		conn.Close()
	}

	client := connection.NewConn(conn)
	r.activeConn.Store(client, struct{}{})
	ch := parser.ParseStream(conn)
	for payload := range ch {
		if payload.Err != nil {
			if payload.Err == io.EOF || payload.Err == io.ErrUnexpectedEOF || strings.Contains(payload.Err.Error(), "use of closed network connection") {
				r.closeClient(client)
				logger.Info("connection closed: " + client.RemoteAddr().String())
				return
			}
			// protocal error
			errReply := reply.MakeStandardErrReply(payload.Err.Error())
			err := client.Write(errReply.ToBytes())
			if err != nil {
				r.closeClient(client)
				logger.Info("connection closed" + client.RemoteAddr().String())
				return
			}
			continue
		}

		// exec
		if payload.Data == nil {
			continue
		}
		reply, ok := payload.Data.(*reply.MultiBulkReply)
		if !ok {
			logger.Error("require multi bulk reply")
			continue
		}
		result := r.db.Exec(client, reply.Args)
		if result != nil {
			client.Write((result.ToBytes()))
		} else {
			client.Write(unknownErrReplyBytes)
		}
	}
}

func (r *RespHandler) AfterClientClose(c resp.Connection) {
	panic("not implemented") // TODO: Implement
}

type RespHandler struct {
	activeConn sync.Map
	db         databaseface.Database
	closing    atomic.Boolean
}

func MakeHandler() *RespHandler {
	var db databaseface.Database
	db = database.NewEchoDatabase()
	return &RespHandler{
		db: db,
	}
}
