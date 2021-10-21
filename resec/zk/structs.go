package zk

import (
	"sync"

	zkapi "github.com/go-zookeeper/zk"
	log "github.com/sirupsen/logrus"
)

type Manager struct {
	mu      sync.Mutex       // lock
	zkConn  *zkapi.Conn      // zookeeper connection
	eventCh chan Event       // event channel
	stopCh  chan interface{} // stop channel
	logger  *log.Entry       // logger

	redisHost      string // redis host announced to zookeeper
	redisPort      int    // redis port announced to zookeeper
	zkBasePath     string // redis base path announced to zookeeper
	registeredPath string // is filled if redis is master and is registered to zookeeper, otherwise empty
}
