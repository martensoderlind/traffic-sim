package tools

type Tool interface {
	Click(mouseX, mouseY float64) error
	Cancel()
}

type DragTool interface {
	Tool
	StartDrag(mouseX, mouseY float64)
	UpdateDrag(mouseX, mouseY float64) error
	EndDrag()
}
