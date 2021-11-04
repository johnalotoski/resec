package zk

import (
	"strings"
	"time"

	zkapi "github.com/go-zookeeper/zk"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func NewConnection(c *cli.Context) (*Manager, error) {
	logger := log.WithField("system", "zk")

	zkServers := c.String("zk-servers")
	if zkServers == "" {
		logger.Printf("No zookeeper configured")
		return &Manager{
			zkConn:  nil,
			eventCh: make(chan Event, 1),
			logger:  logger,
		}, nil
	}

	conn, _, err := zkapi.Connect(strings.Split(zkServers, ","), 1*time.Second)
	if err != nil {
		return nil, err
	}
	logger.Printf("zk: Connected")
	redisHost, redisPort, err := parseAddr(c.String("redis-addr"))
	if err != nil {
		return nil, err
	}

	return &Manager{
		zkConn:     conn,
		eventCh:    make(chan Event, 1),
		logger:     logger,
		redisHost:  redisHost,
		redisPort:  redisPort,
		zkBasePath: c.String("zk-base-path"),
	}, nil
}
