package zk

import (
	"log"
	"sync"
	"time"

	zkapi "github.com/go-zookeeper/zk"
)

func NewConnection() (*Manager, error) {
	log.Printf("zk: NewConnection")
	conn, _, err := zkapi.Connect([]string{"zk1:2181", "zk2:2181", "zk3:2181"}, 1*time.Second)
	if err != nil {
		return nil, err
	}
	log.Printf("zk: Connected")

	return &Manager{
		zkConn:    conn,
		commandCh: make(chan Command, 1),
	}, nil
}

type Manager struct {
	mu sync.Mutex

	zkConn    *zkapi.Conn
	commandCh chan Command
	// state
	registeredPath string // set if master
}

func (m *Manager) GetCommandWriter() chan Command {
	return m.commandCh
}

func (m *Manager) CommandRunner() {
	log.Printf("zk: CommandRunner start: %p", m)
	for {
		select {
		case command := <-m.commandCh:
			switch command.Name() {
			case MasterElected:
				log.Printf("zk: %v", command.RedisState().Info)
				m.registerMaster()
			case NotMasterElected:
				log.Printf("zk: %v", command.RedisState().Info)
				m.deregisterMaster()
			}
		}
	}
	log.Println("zk: CommandRunner finished")
}

func (m *Manager) start() {

}

func (m *Manager) cleanup() {
	m.deregisterMaster()
}

// register master node to zookeeper
func (m *Manager) registerMaster() {
	m.mu.Lock()
	defer m.mu.Unlock()

	log.Printf("zk: registerMaster: %p %s", m, m.registeredPath)
	if m.registeredPath != "" {
		log.Println("zk master already registered")
		return
	}

	p1, err1 := m.zkConn.Create("/gfredis", []byte{}, 0, zkapi.WorldACL(zkapi.PermAll))
	p2, err2 := m.zkConn.Create("/gfredis/prod", []byte{}, 0, zkapi.WorldACL(zkapi.PermAll))
	p3, err3 := m.zkConn.Create("/gfredis/prod/question", []byte{}, 0, zkapi.WorldACL(zkapi.PermAll))
	p4, err4 := m.zkConn.Create("/gfredis/prod/question/nodes", []byte{}, 0, zkapi.WorldACL(zkapi.PermAll))

	log.Println(p1, err1)
	log.Println(p2, err2)
	log.Println(p3, err3)
	log.Println(p4, err4)

	path, err := m.zkConn.Create("/gfredis/prod/question/nodes/member_", []byte("data"), zkapi.FlagSequence|zkapi.FlagEphemeral, zkapi.WorldACL(zkapi.PermAll))
	handleError(err)
	log.Printf("zk master registered: %s", path)

	m.registeredPath = path
}

// deregister master node from zookeeper
func (m *Manager) deregisterMaster() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.registeredPath == "" {
		log.Println("zk master not registered, nothing to deregister")
		return
	}

	m.zkConn.Delete(m.registeredPath, 0)
	log.Printf("zk master deregistered: %s", m.registeredPath)

	m.registeredPath = ""
}

func handleError(err error) {
	log.Println(err)
	//TODO
}
