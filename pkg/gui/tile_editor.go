package gui

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/zrealm-msx/zrealm/pkg/models"
)

// PenTool define a ferramenta ativa de desenho no grid de pixels.
type PenTool int

const (
	PenToggle PenTool = iota // Alterna o pixel
	PenDraw                  // Pinta 1 (Foreground)
	PenErase                 // Pinta 0 (Background)
)

// TileCanvas implementa o grid interativo de edição 8x8 pixels para o modo V9938 SCREEN 4.
type TileCanvas struct {
	widget.BaseWidget
	patternBytes []byte
	colorBytes   []byte
	tool         PenTool
	onChanged    func()

	rects    [8][8]*canvas.Rectangle
	borders  [8][8]*canvas.Rectangle
	cellSize float32
}

// NewTileCanvas inicializa o canvas de desenho 8x8.
func NewTileCanvas(pattern, color []byte, onChanged func()) *TileCanvas {
	tc := &TileCanvas{
		tool:      PenToggle,
		onChanged: onChanged,
		cellSize:  26,
	}
	tc.SetData(pattern, color)
	tc.ExtendBaseWidget(tc)
	return tc
}

// SetData atualiza os dados de padrão (8 bytes) e cor (8 bytes) exibidos no canvas.
func (tc *TileCanvas) SetData(pattern, col []byte) {
	if len(pattern) == 8 {
		tc.patternBytes = make([]byte, 8)
		copy(tc.patternBytes, pattern)
	} else {
		tc.patternBytes = make([]byte, 8)
	}

	if len(col) == 8 {
		tc.colorBytes = make([]byte, 8)
		copy(tc.colorBytes, col)
	} else {
		// Padrão MSX 2: Fg=Branco (15), Bg=Preto (1) -> 0xF1
		tc.colorBytes = []byte{0xF1, 0xF1, 0xF1, 0xF1, 0xF1, 0xF1, 0xF1, 0xF1}
	}
	tc.Refresh()
}

// PatternBytes retorna uma cópia dos 8 bytes de padrão atuais.
func (tc *TileCanvas) PatternBytes() []byte {
	res := make([]byte, 8)
	copy(res, tc.patternBytes)
	return res
}

// ColorBytes retorna uma cópia dos 8 bytes de atributos de cor atuais.
func (tc *TileCanvas) ColorBytes() []byte {
	res := make([]byte, 8)
	copy(res, tc.colorBytes)
	return res
}

// SetTool altera a ferramenta de desenho ativa (Toggle, Draw ou Erase).
func (tc *TileCanvas) SetTool(tool PenTool) {
	tc.tool = tool
}

func (tc *TileCanvas) applyPixel(x, y int) {
	if x < 0 || x >= 8 || y < 0 || y >= 8 {
		return
	}

	shift := uint(7 - x)
	curBit := (tc.patternBytes[y] >> shift) & 1

	switch tc.tool {
	case PenToggle:
		if curBit == 1 {
			tc.patternBytes[y] &^= (1 << shift)
		} else {
			tc.patternBytes[y] |= (1 << shift)
		}
	case PenDraw:
		tc.patternBytes[y] |= (1 << shift)
	case PenErase:
		tc.patternBytes[y] &^= (1 << shift)
	}

	tc.Refresh()
	if tc.onChanged != nil {
		tc.onChanged()
	}
}

// Tapped processa cliques individuais no grid 8x8.
func (tc *TileCanvas) Tapped(e *fyne.PointEvent) {
	size := tc.Size()
	side := min(size.Width, size.Height)
	cs := side / 8
	if cs <= 0 {
		return
	}
	x := int(e.Position.X / cs)
	y := int(e.Position.Y / cs)
	tc.applyPixel(x, y)
}

// Dragged processa pintura por arrasto de mouse contínuo.
func (tc *TileCanvas) Dragged(e *fyne.DragEvent) {
	size := tc.Size()
	side := min(size.Width, size.Height)
	cs := side / 8
	if cs <= 0 {
		return
	}
	x := int(e.Position.X / cs)
	y := int(e.Position.Y / cs)
	tc.applyPixel(x, y)
}

func (tc *TileCanvas) DragEnd() {}

func (tc *TileCanvas) CreateRenderer() fyne.WidgetRenderer {
	r := &tileCanvasRenderer{canvas: tc}
	borderCol := color.NRGBA{R: 0x33, G: 0x41, B: 0x55, A: 0xFF}

	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			b := canvas.NewRectangle(color.Transparent)
			b.StrokeColor = borderCol
			b.StrokeWidth = 1
			tc.borders[y][x] = b

			rect := canvas.NewRectangle(color.Black)
			tc.rects[y][x] = rect
		}
	}
	return r
}

type tileCanvasRenderer struct {
	canvas *TileCanvas
}

func (r *tileCanvasRenderer) Destroy() {}

func (r *tileCanvasRenderer) MinSize() fyne.Size {
	return fyne.NewSize(220, 220)
}

func (r *tileCanvasRenderer) Layout(size fyne.Size) {
	side := min(size.Width, size.Height)
	cs := side / 8
	r.canvas.cellSize = cs

	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			pos := fyne.NewPos(float32(x)*cs, float32(y)*cs)
			sz := fyne.NewSize(cs, cs)

			r.canvas.rects[y][x].Move(pos)
			r.canvas.rects[y][x].Resize(sz)

			r.canvas.borders[y][x].Move(pos)
			r.canvas.borders[y][x].Resize(sz)
		}
	}
}

func (r *tileCanvasRenderer) Refresh() {
	for y := 0; y < 8; y++ {
		colByte := r.canvas.colorBytes[y]
		fgIdx := (colByte >> 4) & 0x0F
		bgIdx := colByte & 0x0F

		fgColor := GetMSXColor(fgIdx)
		bgColor := GetMSXColor(bgIdx)

		patByte := r.canvas.patternBytes[y]
		for x := 0; x < 8; x++ {
			bit := (patByte >> uint(7-x)) & 1
			if bit == 1 {
				r.canvas.rects[y][x].FillColor = fgColor
			} else {
				r.canvas.rects[y][x].FillColor = bgColor
			}
			r.canvas.rects[y][x].Refresh()
			r.canvas.borders[y][x].Refresh()
		}
	}
}

func (r *tileCanvasRenderer) Objects() []fyne.CanvasObject {
	objs := make([]fyne.CanvasObject, 0, 128)
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			objs = append(objs, r.canvas.rects[y][x])
			objs = append(objs, r.canvas.borders[y][x])
		}
	}
	return objs
}

// -------------------------------------------------------------
// Componente de Edição Completa de Tile (TileEditorWidget)
// -------------------------------------------------------------

// TileEditorWidget reúne o grid de pixels, controles de cor por linha, preview e persistência.
type TileEditorWidget struct {
	state          *ProjectState
	win            fyne.Window
	tile           *models.Tile
	canvas         *TileCanvas
	previewImage   *canvas.Image
	collisionRadio *widget.RadioGroup
	lblStatus      *widget.Label
	container      *fyne.Container

	fgButtons [8]*widget.Button
	bgButtons [8]*widget.Button

	onSaved func()
}

// NewTileEditorWidget constrói o painel de edição visual completo de um tile do MSX 2.
func NewTileEditorWidget(state *ProjectState, win fyne.Window, onSaved func()) *TileEditorWidget {
	te := &TileEditorWidget{
		state:     state,
		win:       win,
		lblStatus: widget.NewLabel("Nenhum tile selecionado."),
		onSaved:   onSaved,
	}

	te.canvas = NewTileCanvas(nil, nil, func() {
		te.updatePreview()
		te.lblStatus.SetText("Alterações não salvas no tile.")
	})

	// Preview 64x64 pixels do tile atual
	te.previewImage = canvas.NewImageFromImage(te.renderImage())
	te.previewImage.SetMinSize(fyne.NewSize(64, 64))
	te.previewImage.ScaleMode = canvas.ImageScalePixels

	// Seletor de Física/Colisão
	te.collisionRadio = widget.NewRadioGroup([]string{
		"Passável (0 - Livre)",
		"Sólido (1 - Parede)",
		"Água (2 - Rio)",
		"Dano (3 - Lava/Espinhos)",
		"Gatilho (4 - Evento)",
	}, func(selected string) {
		if te.tile != nil {
			te.tile.CollisionType = collisionTypeFromString(selected)
			te.lblStatus.SetText("Colisão alterada (salvamento pendente).")
		}
	})
	te.collisionRadio.Horizontal = false

	// Botões de Ferramentas de Pintura
	btnToggle := widget.NewButtonWithIcon("Alternar", theme.ContentUndoIcon(), func() {
		te.canvas.SetTool(PenToggle)
	})
	btnDraw := widget.NewButtonWithIcon("Lápis (Fg)", theme.DocumentCreateIcon(), func() {
		te.canvas.SetTool(PenDraw)
	})
	btnErase := widget.NewButtonWithIcon("Borracha (Bg)", theme.ContentClearIcon(), func() {
		te.canvas.SetTool(PenErase)
	})

	toolsBox := container.NewHBox(btnToggle, btnDraw, btnErase)

	// Botões de Transformações Rápidas
	btnRotate := widget.NewButtonWithIcon("Girar 90°", theme.ViewRefreshIcon(), func() {
		te.canvas.patternBytes = rotatePattern90(te.canvas.patternBytes)
		te.canvas.Refresh()
		te.updatePreview()
		te.lblStatus.SetText("Tile girado em 90°.")
	})
	btnFlipH := widget.NewButton("Inverter H", func() {
		te.canvas.patternBytes = flipHorizontal(te.canvas.patternBytes)
		te.canvas.Refresh()
		te.updatePreview()
		te.lblStatus.SetText("Espelhado horizontalmente.")
	})
	btnFlipV := widget.NewButton("Inverter V", func() {
		te.canvas.patternBytes = flipVertical(te.canvas.patternBytes)
		te.canvas.Refresh()
		te.updatePreview()
		te.lblStatus.SetText("Espelhado verticalmente.")
	})
	btnClear := widget.NewButton("Limpar", func() {
		te.canvas.patternBytes = make([]byte, 8)
		te.canvas.Refresh()
		te.updatePreview()
		te.lblStatus.SetText("Padrão limpo (todos pixels 0).")
	})
	btnFill := widget.NewButton("Preencher", func() {
		te.canvas.patternBytes = []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
		te.canvas.Refresh()
		te.updatePreview()
		te.lblStatus.SetText("Padrão preenchido (todos pixels 1).")
	})

	transformBox := container.NewGridWithColumns(3,
		btnRotate, btnFlipH, btnFlipV,
		btnClear, btnFill,
	)

	// Seletor de Cores por Scanline (8 Linhas do V9938)
	colorRows := container.NewVBox()
	for line := 0; line < 8; line++ {
		l := line
		btnFg := widget.NewButton(fmt.Sprintf("Fg: %d", (te.canvas.colorBytes[l]>>4)&0x0F), func() {
			te.pickColor(l, true)
		})
		btnBg := widget.NewButton(fmt.Sprintf("Bg: %d", te.canvas.colorBytes[l]&0x0F), func() {
			te.pickColor(l, false)
		})
		te.fgButtons[l] = btnFg
		te.bgButtons[l] = btnBg

		row := container.NewHBox(
			widget.NewLabel(fmt.Sprintf("Linha %d:", l)),
			btnFg,
			btnBg,
		)
		colorRows.Add(row)
	}

	btnCopyColorToAll := widget.NewButton("Aplicar Cores da Linha 0 a Todas", func() {
		baseCol := te.canvas.colorBytes[0]
		for i := 0; i < 8; i++ {
			te.canvas.colorBytes[i] = baseCol
		}
		te.updateColorButtons()
		te.canvas.Refresh()
		te.updatePreview()
		te.lblStatus.SetText("Cores aplicadas a todas as 8 linhas.")
	})

	colorCard := widget.NewCard(
		"Cores V9938 (Foreground / Background)",
		"Cada linha de 8 pixels possui 1 cor de frente e 1 cor de fundo",
		container.NewVBox(colorRows, btnCopyColorToAll),
	)

	// Botão Principal de Salvamento
	btnSave := widget.NewButtonWithIcon("Salvar Tile no Projeto SQLite", theme.DocumentSaveIcon(), func() {
		te.SaveCurrentTile()
	})
	btnSave.Importance = widget.HighImportance

	// Montagem do Layout
	leftColumn := container.NewVBox(
		widget.NewLabelWithStyle("Grid Pixel-a-Pixel (8x8)", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		toolsBox,
		container.NewCenter(te.canvas),
		widget.NewSeparator(),
		widget.NewLabelWithStyle("Transformações do Padrão", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		transformBox,
	)

	rightColumn := container.NewVBox(
		container.NewHBox(
			container.NewVBox(
				widget.NewLabelWithStyle("Preview (64x64)", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				container.NewCenter(te.previewImage),
			),
			container.NewVBox(
				widget.NewLabelWithStyle("Propriedades Físicas", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				te.collisionRadio,
			),
		),
		widget.NewSeparator(),
		colorCard,
		widget.NewSeparator(),
		btnSave,
		te.lblStatus,
	)

	split := container.NewHSplit(leftColumn, rightColumn)
	split.Offset = 0.45

	te.container = container.NewPadded(split)
	return te
}

// Container retorna o container visual do editor de tiles.
func (te *TileEditorWidget) Container() fyne.CanvasObject {
	return te.container
}

// LoadTile carrega um tile existente para o editor.
func (te *TileEditorWidget) LoadTile(t *models.Tile) {
	te.tile = t
	if t == nil {
		te.canvas.SetData(nil, nil)
		te.lblStatus.SetText("Nenhum tile selecionado.")
		return
	}

	te.canvas.SetData(t.PatternBytes, t.ColorBytes)
	te.updateColorButtons()
	te.updatePreview()

	// Seleciona o radio correto
	switch t.CollisionType {
	case models.CollisionPassable:
		te.collisionRadio.SetSelected("Passável (0 - Livre)")
	case models.CollisionSolid:
		te.collisionRadio.SetSelected("Sólido (1 - Parede)")
	case models.CollisionWater:
		te.collisionRadio.SetSelected("Água (2 - Rio)")
	case models.CollisionDamage:
		te.collisionRadio.SetSelected("Dano (3 - Lava/Espinhos)")
	case models.CollisionTrigger:
		te.collisionRadio.SetSelected("Gatilho (4 - Evento)")
	}

	te.lblStatus.SetText(fmt.Sprintf("Editando Tile #%d (Índice VRAM: %d)", t.ID, t.TileIndex))
}

func (te *TileEditorWidget) updateColorButtons() {
	for l := 0; l < 8; l++ {
		fg := (te.canvas.colorBytes[l] >> 4) & 0x0F
		bg := te.canvas.colorBytes[l] & 0x0F
		if te.fgButtons[l] != nil {
			te.fgButtons[l].SetText(fmt.Sprintf("Fg: %d (%s)", fg, MSXPalette[fg].Name))
		}
		if te.bgButtons[l] != nil {
			te.bgButtons[l].SetText(fmt.Sprintf("Bg: %d (%s)", bg, MSXPalette[bg].Name))
		}
	}
}

func (te *TileEditorWidget) pickColor(line int, isForeground bool) {
	swatches := make([]fyne.CanvasObject, 16)
	var winDialog dialog.Dialog

	for i := 0; i < 16; i++ {
		colorIdx := uint8(i)
		btn := widget.NewButton(fmt.Sprintf("%d: %s", colorIdx, MSXPalette[colorIdx].Name), func() {
			curCol := te.canvas.colorBytes[line]
			if isForeground {
				te.canvas.colorBytes[line] = (colorIdx << 4) | (curCol & 0x0F)
			} else {
				te.canvas.colorBytes[line] = (curCol & 0xF0) | (colorIdx & 0x0F)
			}
			te.updateColorButtons()
			te.canvas.Refresh()
			te.updatePreview()
			te.lblStatus.SetText(fmt.Sprintf("Cor da linha %d atualizada para %s.", line, MSXPalette[colorIdx].Name))
			if winDialog != nil {
				winDialog.Hide()
			}
		})
		swatches[i] = btn
	}

	grid := container.NewGridWithColumns(2, swatches...)
	title := fmt.Sprintf("Escolha a Cor de Frente (Foreground) - Linha %d", line)
	if !isForeground {
		title = fmt.Sprintf("Escolha a Cor de Fundo (Background) - Linha %d", line)
	}

	winDialog = dialog.NewCustom(title, "Cancelar", container.NewPadded(grid), te.win)
	winDialog.Resize(fyne.NewSize(380, 420))
	winDialog.Show()
}

func (te *TileEditorWidget) updatePreview() {
	img := te.renderImage()
	te.previewImage.Image = img
	te.previewImage.Refresh()
}

func (te *TileEditorWidget) renderImage() image.Image {
	img := image.NewNRGBA(image.Rect(0, 0, 8, 8))
	for y := 0; y < 8; y++ {
		colByte := te.canvas.colorBytes[y]
		fg := GetMSXColor((colByte >> 4) & 0x0F)
		bg := GetMSXColor(colByte & 0x0F)
		pat := te.canvas.patternBytes[y]

		for x := 0; x < 8; x++ {
			bit := (pat >> uint(7-x)) & 1
			if bit == 1 {
				img.SetNRGBA(x, y, fg)
			} else {
				img.SetNRGBA(x, y, bg)
			}
		}
	}
	return img
}

// SaveCurrentTile grava os dados de padrão, cores e colisão no banco SQLite.
func (te *TileEditorWidget) SaveCurrentTile() {
	if te.tile == nil {
		dialog.ShowInformation("Aviso", "Nenhum tile selecionado para salvar.", te.win)
		return
	}
	proj := te.state.Current()
	if proj == nil {
		dialog.ShowInformation("Aviso", "Nenhum projeto aberto.", te.win)
		return
	}

	te.tile.PatternBytes = te.canvas.PatternBytes()
	te.tile.ColorBytes = te.canvas.ColorBytes()

	err := proj.Storage().Tilesets.SaveTile(context.Background(), te.tile)
	if err != nil {
		dialog.ShowError(fmt.Errorf("falha ao salvar tile: %w", err), te.win)
		return
	}

	te.state.NotifyDataChanged()
	te.lblStatus.SetText(fmt.Sprintf("Tile #%d salvo com sucesso no banco de dados!", te.tile.ID))
	if te.onSaved != nil {
		te.onSaved()
	}
}

// -------------------------------------------------------------
// Funções Utilitárias de Transformação
// -------------------------------------------------------------

func rotatePattern90(p []byte) []byte {
	res := make([]byte, 8)
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			bit := (p[y] >> uint(7-x)) & 1
			if bit == 1 {
				res[x] |= 1 << uint(y)
			}
		}
	}
	return res
}

func flipHorizontal(p []byte) []byte {
	res := make([]byte, 8)
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			bit := (p[y] >> uint(7-x)) & 1
			if bit == 1 {
				res[y] |= 1 << uint(x)
			}
		}
	}
	return res
}

func flipVertical(p []byte) []byte {
	res := make([]byte, 8)
	for y := 0; y < 8; y++ {
		res[7-y] = p[y]
	}
	return res
}

func collisionTypeFromString(s string) models.CollisionType {
	switch {
	case strings.HasPrefix(s, "Passável") || strings.HasPrefix(s, "P"):
		return models.CollisionPassable
	case strings.HasPrefix(s, "Sólido") || strings.HasPrefix(s, "S"):
		return models.CollisionSolid
	case strings.HasPrefix(s, "Água") || strings.HasPrefix(s, "Agua") || strings.HasPrefix(s, "A"):
		return models.CollisionWater
	case strings.HasPrefix(s, "Dano") || strings.HasPrefix(s, "D"):
		return models.CollisionDamage
	case strings.HasPrefix(s, "Gatilho") || strings.HasPrefix(s, "G"):
		return models.CollisionTrigger
	default:
		return models.CollisionPassable
	}
}
