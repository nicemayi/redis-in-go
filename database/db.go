package database

import (
	"redis-in-go/datastructure/dict"
	"redis-in-go/interface/resp"
)

type DB struct {
	index int
	data  dict.Dict
}

func makeDB() *DB {
	db := &DB{
		data: dict.MakeSyncDict(),
	}
	return db
}

type ExecFunc func(db *DB, args [][]byte) resp.Reply

type CmdLine = [][]byte
