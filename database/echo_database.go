package database

import (
	"redis-in-go/interface/resp"
	"redis-in-go/resp/reply"
)

type EchoDatabase struct {
}

func NewEchoDatabase() *EchoDatabase {
	return &EchoDatabase{}
}

func (e EchoDatabase) Exec(client resp.Connection, args [][]byte) resp.Reply {
	return reply.MakeMultiBulkReply(args)
}

func (e EchoDatabase) Close() {
	panic("implement me")
}

func (e EchoDatabase) AfterClientClose(c resp.Connection) {
	panic("implement me")
}

func (e EchoDatabase) Error() error {
	panic("implement me")
}
