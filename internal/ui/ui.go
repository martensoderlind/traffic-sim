package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type UIManager struct {
	buttons []*Button
	labels  []*Label
}

func NewUIManager() *UIManager {
	return &UIManager{
		buttons: make([]*Button, 0),
		labels:  make([]*Label, 0),
	}
}

func (ui *UIManager) AddButton(btn *Button) {
	ui.buttons = append(ui.buttons, btn)
}

func (ui *UIManager) AddLabel(lbl *Label) {
	ui.labels = append(ui.labels, lbl)
}

func (ui *UIManager) Update(mouseX, mouseY int, clicked bool) {
	for _, btn := range ui.buttons {
		btn.Update(mouseX, mouseY, clicked)
	}
}

func (ui *UIManager) Draw(screen *ebiten.Image) {
	for _, btn := range ui.buttons {
		btn.Draw(screen)
	}
	
	for _, lbl := range ui.labels {
		lbl.Draw(screen)
	}
}

func (ui *UIManager) Clear() {
	ui.buttons = make([]*Button, 0)
	ui.labels = make([]*Label, 0)
}
