package cluster

import (
	"redis-in-go/interface/resp"
	"redis-in-go/resp/reply"
)

func flushdb(cluster *ClusterDatabase, c resp.Connection, cmdArgs [][]byte) resp.Reply {
	replies := cluster.broadcase(c, cmdArgs)
	var errReply reply.ErrorReply
	for _, r := range replies {
		if reply.IsErrReply(r) {
			errReply = r.(reply.ErrorReply)
			break
		}
	}

	if errReply == nil {
		return reply.MakeOkReply()
	}
	return reply.MakeStandardErrReply("error: " + errReply.Error())
}
