package aof

import (
	"os"
	"redis-in-go/database"
)

type CmdLine = [][]byte

type payload struct {
	cmdLine CmdLine
	dbIndex int
}

type AofHandler struct {
	database    database.Database
	aofChan     chan *payload
	aofFile     *os.File
	aofFilename string
	currentDB   int
}

// NewAofHandler

// Add payload(set k v) -> aofChan

// handleAof payload(set k ) <- aofChan

// LoadAof
