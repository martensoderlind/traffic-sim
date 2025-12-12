package commands

import (
	"traffic-sim/internal/persistence"
	"traffic-sim/internal/world"
)

type SaveWorldCommand struct{}

func (c *SaveWorldCommand) ExecuteReadUnlocked(w *world.World) error {
	saveData := persistence.SerializeWorld(w)
	return persistence.SaveToFile(saveData)
}

func (c *SaveWorldCommand) Execute(w *world.World) error {
	return c.ExecuteReadUnlocked(w)
}