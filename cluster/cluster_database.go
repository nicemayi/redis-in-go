package cluster

import (
	"context"
	"redis-in-go/config"
	database2 "redis-in-go/database"
	"redis-in-go/interface/database"
	"redis-in-go/interface/resp"
	"redis-in-go/lib/consistenthash"
	"redis-in-go/lib/logger"
	"redis-in-go/resp/reply"
	"strings"

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

type CmdFunc func(cluster *ClusterDatabase, c resp.Connection, cmdArgs [][]byte) resp.Reply

var router = makeRouter()

func (c *ClusterDatabase) Exec(client resp.Connection, args [][]byte) (result resp.Reply) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(err)
			result = reply.UnknownErrReply{}
		}
	}()

	cmdName := strings.ToLower(string(args[0]))
	cmdFunc, ok := router[cmdName]
	if !ok {
		return reply.MakeStandardErrReply("not support cmd")
	}
	result = cmdFunc(c, client, args)

	return
}

func (c *ClusterDatabase) Close() {
	c.db.Close()
}

func (cluster *ClusterDatabase) AfterClientClose(d resp.Connection) {
	cluster.db.AfterClientClose(d)
}
