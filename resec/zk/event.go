package zk

import "github.com/seatgeek/resec/resec/state"

// zk has to get events from
const (
	NodeAsMasterElected    = EventName("node_as_master_elected")
	NodeNotAsMasterElected = EventName("node_not_as_master_elected")
)

type EventName string

type Event struct {
	name       EventName
	redisState state.Redis
}

func (c *Event) Name() EventName {
	return c.name
}

func (c *Event) String() string {
	return string(c.name)
}

func (c *Event) RedisState() state.Redis {
	return c.redisState
}

func NewEvent(event EventName, redisState state.Redis) Event {
	return Event{
		name:       event,
		redisState: redisState,
	}
}
