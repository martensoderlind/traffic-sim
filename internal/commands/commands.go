package commands

import "traffic-sim/internal/world"

type Command interface {
	Execute(w *world.World) error
}

type CommandExecutor struct {
	world *world.World
}

func NewCommandExecutor(w *world.World) *CommandExecutor {
	return &CommandExecutor{world: w}
}

func (e *CommandExecutor) Execute(cmd Command) error {
	return cmd.Execute(e.world)
}