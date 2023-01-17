package cluster

import (
	"context"
	"errors"
	"redis-in-go/interface/resp"
	"redis-in-go/lib/utils"
	"redis-in-go/resp/client"
	"redis-in-go/resp/reply"
	"strconv"
)

func (cluster *ClusterDatabase) getPeerClient(peer string) (*client.Client, error) {
	pool, ok := cluster.peerconnection[peer]
	if !ok {
		return nil, errors.New("connection not found")
	}

	object, err := pool.BorrowObject(context.Background())
	if err != nil {
		return nil, err
	}

	client, ok := object.(*client.Client)
	if !ok {
		return nil, errors.New("wrong type")
	}

	return client, err
}

func (cluster *ClusterDatabase) returnPeerClient(peer string, peerClient *client.Client) error {
	pool, ok := cluster.peerconnection[peer]
	if !ok {
		return errors.New("connection not found")
	}
	return pool.ReturnObject(context.Background(), peerClient)
}

func (cluster *ClusterDatabase) relay(peer string, c resp.Connection, args [][]byte) resp.Reply {
	if peer == cluster.self {
		return cluster.db.Exec(c, args)
	}

	peerClient, err := cluster.getPeerClient(peer)
	if err != nil {
		return reply.MakeStandardErrReply(err.Error())
	}

	defer func() {
		cluster.returnPeerClient(peer, peerClient)
	}()

	peerClient.Send(utils.ToCmdLine("SELECT", strconv.Itoa(c.GetDBIndex())))
	return peerClient.Send(args)
}

func (cluster *ClusterDatabase) broadcase(c resp.Connection, args [][]byte) map[string]resp.Reply {
	results := make(map[string]resp.Reply)
	for _, node := range cluster.nodes {
		result := cluster.relay(node, c, args)
		results[node] = result
	}

	return results
}
