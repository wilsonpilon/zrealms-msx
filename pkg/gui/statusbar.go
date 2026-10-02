package gui

import (
	"context"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/zrealm-msx/zrealm/pkg/project"
	"github.com/zrealm-msx/zrealm/pkg/version"
)

// StatusBar representa a barra de status inferior da IDE Z-Realm.
type StatusBar struct {
	state      *ProjectState
	lblProject *widget.Label
	lblMemory  *widget.Label
	lblVersion *widget.Label
	iconStatus *widget.Icon
	container  *fyne.Container
}

// NewStatusBar inicializa a barra de status inferior.
func NewStatusBar(state *ProjectState) *StatusBar {
	sb := &StatusBar{
		state:      state,
		lblProject: widget.NewLabel("Projeto: [Nenhum Aberto]"),
		lblMemory:  widget.NewLabel("Mapper: 0 seg (0 KB)"),
		lblVersion: widget.NewLabel(fmt.Sprintf("Z-Realm %s", version.String())),
		iconStatus: widget.NewIcon(theme.RadioButtonIcon()),
	}

	state.OnProjectLoaded(func(_ *project.Project, path string) {
		sb.lblProject.SetText(fmt.Sprintf("Projeto: %s", path))
		sb.updateMemoryEstimate()
	})

	state.OnProjectClosed(func() {
		sb.lblProject.SetText("Projeto: [Nenhum Aberto]")
		sb.lblMemory.SetText("Mapper: 0 seg (0 KB)")
	})

	state.OnModifiedChanged(func(mod bool) {
		p := sb.state.FilePath()
		if p == "" {
			p = "[Nenhum Aberto]"
		}
		if mod {
			sb.lblProject.SetText(fmt.Sprintf("Projeto: %s *", p))
		} else {
			sb.lblProject.SetText(fmt.Sprintf("Projeto: %s", p))
		}
	})

	state.OnDataChanged(func() {
		sb.updateMemoryEstimate()
	})

	sep1 := widget.NewSeparator()
	sep2 := widget.NewSeparator()

	sb.container = container.NewHBox(
		sb.iconStatus,
		sb.lblProject,
		sep1,
		sb.lblMemory,
		sep2,
		sb.lblVersion,
	)

	return sb
}

func (sb *StatusBar) updateMemoryEstimate() {
	proj := sb.state.Current()
	if proj == nil {
		sb.lblMemory.SetText("Mapper: 0 seg (0 KB)")
		return
	}
	rooms, _ := proj.Storage().Rooms.ListRooms(context.Background())
	// Estimativa: 576 bytes por sala + 2048 bytes por tileset
	totalBytes := len(rooms)*576 + 2048
	segments := (totalBytes + 16383) / 16384
	if segments == 0 {
		segments = 1
	}
	sb.lblMemory.SetText(fmt.Sprintf("Mapper: %d seg (%d KB / %d salas)", segments, segments*16, len(rooms)))
}

// Container retorna o objeto visual da barra de status.
func (sb *StatusBar) Container() fyne.CanvasObject {
	return container.NewVBox(widget.NewSeparator(), sb.container)
}
