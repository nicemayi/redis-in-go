package cluster

import (
	"redis-in-go/interface/resp"
	"redis-in-go/resp/reply"
)

func Del(cluster *ClusterDatabase, c resp.Connection, cmdArgs [][]byte) resp.Reply {
	replies := cluster.broadcase(c, cmdArgs)
	var errReply reply.ErrorReply
	var deleted int64 = 0
	for _, r := range replies {
		if reply.IsErrReply(r) {
			errReply = r.(reply.ErrorReply)
			break
		}

		intReply, ok := r.(*reply.IntReply)
		if !ok {
			errReply = reply.MakeStandardErrReply("error")
		}

		deleted += intReply.Code
	}

	if errReply == nil {
		return reply.MakeIntReply(deleted)
	}
	return reply.MakeStandardErrReply("error: " + errReply.Error())
}
