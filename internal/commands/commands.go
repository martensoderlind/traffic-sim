package commands

import "traffic-sim/internal/world"

type Command interface {
	Execute(w *world.World) error
}

type ExecuteWithLocking interface {
	ExecuteUnlocked(w *world.World) error
}

type ExecuteWithReadLocking interface {
    ExecuteReadUnlocked(w *world.World) error
}

type CommandExecutor struct {
	world *world.World
}

func NewCommandExecutor(w *world.World) *CommandExecutor {
	return &CommandExecutor{world: w}
}

func (e *CommandExecutor) Execute(cmd Command) error {
	if readCmd, ok := cmd.(ExecuteWithReadLocking); ok {
		e.world.Mu.RLock()
		defer e.world.Mu.RUnlock()
		return readCmd.ExecuteReadUnlocked(e.world)
	}

	if lockingCmd, ok := cmd.(ExecuteWithLocking); ok {
		e.world.Mu.Lock()
		defer e.world.Mu.Unlock()
		return lockingCmd.ExecuteUnlocked(e.world)
	}

	return cmd.Execute(e.world)
}