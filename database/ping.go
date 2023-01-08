package database

import (
	"redis-in-go/interface/resp"
	"redis-in-go/resp/reply"
)

func Ping(db *DB, args [][]byte) resp.Reply {
	return reply.MakePongReply()
}

func init() {
	RegisterCommand("ping", Ping, 1)
}
