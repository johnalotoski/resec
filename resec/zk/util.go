package zk

import (
	"fmt"
	"path"
	"strconv"
	"strings"

	"github.com/go-zookeeper/zk"
)

// creates path in zookeeper recursively if not exists
func zkEnsurePath(zkConn *zk.Conn, zkPath string) error {
	if b, _, err := zkConn.Exists(zkPath); err != nil {
		return err
	} else if b {
		return nil
	}

	parent := path.Dir(zkPath)
	if err := zkEnsurePath(zkConn, parent); err != nil {
		return err
	}

	_, err := zkConn.Create(zkPath, []byte{}, 0, zkPermission)
	return err
}

// parses and breaks address into host and port
func parseAddr(addr string) (host string, port int, err error) {
	sArr := strings.Split(addr, ":")

	if len(sArr) == 1 {
		host = sArr[0]
		port = defaultRedisPort
	} else if len(sArr) == 2 {
		host = sArr[0]
		portNumber, err := strconv.ParseInt(sArr[1], 10, 64)
		if err != nil {
			err = fmt.Errorf("Error parsing redis addr's port: %s", sArr[1])
		} else {
			port = int(portNumber)
		}
	} else {
		err = fmt.Errorf("Unknown format of REDIS_ADDR: %s", addr)
	}

	return
}
