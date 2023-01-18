package cluster

import (
	"redis-in-go/interface/resp"
	"redis-in-go/resp/reply"
)

func Rename(cluster *ClusterDatabase, c resp.Connection, cmdArgs [][]byte) resp.Reply {
	if len(cmdArgs) != 3 {
		return reply.MakeStandardErrReply("ERR Wrong number args")
	}

	src := string(cmdArgs[1])
	dest := string(cmdArgs[2])

	srcPeer := cluster.peerPicker.PickNode(src)
	destPeer := cluster.peerPicker.PickNode(dest)

	if srcPeer != destPeer {
		return reply.MakeStandardErrReply("ERR rename must within one peer")
	}

	return cluster.relay(srcPeer, c, cmdArgs)
}
