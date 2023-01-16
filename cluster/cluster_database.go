package cluster

import (
	"context"
	"redis-in-go/config"
	database2 "redis-in-go/database"
	"redis-in-go/interface/database"
	"redis-in-go/interface/resp"
	"redis-in-go/lib/consistenthash"

	pool "github.com/jolestar/go-commons-pool"
)

type ClusterDatabase struct {
	self           string
	nodes          []string
	peerPicker     *consistenthash.NodeMap
	peerconnection map[string]*pool.ObjectPool
	db             database.Database
}

func MakeClusterDatabase() *ClusterDatabase {
	cluster := &ClusterDatabase{
		self:           config.Properties.Self,
		db:             database2.NewStandaloneDatabase(),
		peerPicker:     consistenthash.NewNodeMap(nil),
		peerconnection: make(map[string]*pool.ObjectPool),
	}

	nodes := make([]string, 0, len(config.Properties.Peers)+1)
	for _, peer := range config.Properties.Peers {
		nodes = append(nodes, peer)
	}
	nodes = append(nodes, config.Properties.Self)
	cluster.peerPicker.AddNode(nodes...)

	ctx := context.Background()
	for _, peer := range config.Properties.Peers {
		pool.NewObjectPoolWithDefaultConfig(ctx, &connectionFactory{
			Peer: peer,
		})
	}

	cluster.nodes = nodes
	return cluster
}

func (c *ClusterDatabase) Exec(client resp.Connection, args [][]byte) resp.Reply {
	panic("implement me")
}

func (c *ClusterDatabase) Close() {
	panic("implement me")
}

func (c *ClusterDatabase) AfterClientClose(d resp.Connection) {
	panic("implement me")
}
