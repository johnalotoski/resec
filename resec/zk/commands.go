package zk

import "github.com/seatgeek/resec/resec/state"

const (
	MasterElected    = CommandName("master_elected")
	NotMasterElected = CommandName("not_master_elected")
)

type CommandName string

type Command struct {
	name       CommandName
	redisState state.Redis
}

func (c *Command) Name() CommandName {
	return c.name
}

func (c *Command) String() string {
	return string(c.name)
}

func (c *Command) RedisState() state.Redis {
	return c.redisState
}

func NewCommand(cmd CommandName, redisState state.Redis) Command {
	return Command{
		name:       cmd,
		redisState: redisState,
	}
}
