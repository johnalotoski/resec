package zk

import (
	"encoding/json"
	"fmt"

	zkapi "github.com/go-zookeeper/zk"
)

const defaultRedisPort = 6379

var zkPermission = zkapi.WorldACL(zkapi.PermAll)

func (m *Manager) GetEventWriter() chan Event {
	return m.eventCh
}

func (m *Manager) EventRunner() {
	m.logger.Info("EventRunner start")
	for {
		select {
		case <-m.stopCh:
			m.logger.Info("Shutting down zookeeper EventRunner")
			m.deregisterMaster()
			return
		case command := <-m.eventCh:
			switch command.Name() {
			case NodeAsMasterElected:
				m.logger.Debugf("zk: %v", command.RedisState().Info)
				m.logger.Infof("REDIS_ADDR: %s:%d", m.redisHost, m.redisPort)
				m.registerMaster()
			case NodeNotAsMasterElected:
				m.logger.Debugf("zk: %v", command.RedisState().Info)
				m.deregisterMaster()
			}
		}
	}
}

// register master node to zookeeper
func (m *Manager) registerMaster() {
	if m.zkConn == nil { // no zookeeper configured
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.logger.Infof("registerMaster: %p %s", m, m.registeredPath)
	if m.registeredPath != "" {
		m.logger.Info("zk master already registered")
		return
	}

	zkPath := m.zkBasePath
	err := zkEnsurePath(m.zkConn, zkPath)
	if err != nil {
		m.logger.Errorf("Failed to ensure path in zk: %v", err)
		return
	}

	data, err := json.Marshal(map[string]interface{}{
		"serviceEndpoint": map[string]interface{}{
			"host": m.redisHost,
			"port": m.redisPort,
		},
		"additionalEndpoints": map[string]interface{}{},
		"status":              "ALIVE",
	})
	if err != nil {
		m.logger.Errorf("Failed to ensure path in zk: %v", err)
		return
	}

	memberPath := fmt.Sprintf("%s/member_", zkPath)
	path, err := m.zkConn.Create(memberPath, []byte(data), zkapi.FlagSequence|zkapi.FlagEphemeral, zkPermission)
	if err != nil {
		m.logger.Errorf("Failed creating member %s in zookeeper: %v", memberPath, err)
		return
	}

	m.logger.Infof("Registered Master: %s", path)
	m.registeredPath = path
}

// deregister master node from zookeeper
func (m *Manager) deregisterMaster() {
	if m.zkConn == nil { // no zookeeper configured
		return
	}

	var err error

	m.mu.Lock()
	defer m.mu.Unlock()

	if m.registeredPath == "" {
		m.logger.Debug("Master not registered, nothing to deregister")
		return
	}

	err = m.zkConn.Delete(m.registeredPath, 0)
	if err != nil {
		m.logger.Errorf("Failed to delete path %s", m.registeredPath)
		return
	}

	m.logger.Infof("zk master deregistered: %s", m.registeredPath)
	m.registeredPath = ""
}
