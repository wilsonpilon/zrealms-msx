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

// SpriteCanvas implementa o grid interativo de 16x16 pixels do Modo 2 do V9938.
type SpriteCanvas struct {
	widget.BaseWidget
	sprite    *models.Sprite
	tool      PenTool
	onChanged func()

	rects    [16][16]*canvas.Rectangle
	borders  [16][16]*canvas.Rectangle
	cellSize float32
}

// NewSpriteCanvas inicializa o canvas de desenho 16x16.
func NewSpriteCanvas(s *models.Sprite, onChanged func()) *SpriteCanvas {
	sc := &SpriteCanvas{
		tool:      PenToggle,
		onChanged: onChanged,
		cellSize:  18,
	}
	sc.SetSprite(s)
	sc.ExtendBaseWidget(sc)
	return sc
}

// SetSprite associa o sprite atual a ser desenhado no canvas.
func (sc *SpriteCanvas) SetSprite(s *models.Sprite) {
	if s == nil {
		sc.sprite = models.NewSprite("Temporario")
	} else {
		sc.sprite = s
	}
	sc.Refresh()
}

// SetTool define a ferramenta de pintura ativa (Toggle, Draw ou Erase).
func (sc *SpriteCanvas) SetTool(t PenTool) {
	sc.tool = t
}

func (sc *SpriteCanvas) applyPixel(x, y int) {
	if sc.sprite == nil || x < 0 || x >= 16 || y < 0 || y >= 16 {
		return
	}

	cur, err := sc.sprite.GetPixel(x, y)
	if err != nil {
		return
	}

	var newVal bool
	switch sc.tool {
	case PenToggle:
		newVal = !cur
	case PenDraw:
		newVal = true
	case PenErase:
		newVal = false
	}

	_ = sc.sprite.SetPixel(x, y, newVal)
	sc.Refresh()
	if sc.onChanged != nil {
		sc.onChanged()
	}
}

// Tapped processa cliques individuais no grid 16x16.
func (sc *SpriteCanvas) Tapped(e *fyne.PointEvent) {
	size := sc.Size()
	side := min(size.Width, size.Height)
	cs := side / 16
	if cs <= 0 {
		return
	}
	x := int(e.Position.X / cs)
	y := int(e.Position.Y / cs)
	sc.applyPixel(x, y)
}

// Dragged processa pintura com arrasto contínuo do mouse.
func (sc *SpriteCanvas) Dragged(e *fyne.DragEvent) {
	size := sc.Size()
	side := min(size.Width, size.Height)
	cs := side / 16
	if cs <= 0 {
		return
	}
	x := int(e.Position.X / cs)
	y := int(e.Position.Y / cs)
	sc.applyPixel(x, y)
}

func (sc *SpriteCanvas) DragEnd() {}

func (sc *SpriteCanvas) CreateRenderer() fyne.WidgetRenderer {
	r := &spriteCanvasRenderer{canvas: sc}
	normalBorder := color.NRGBA{R: 0x2A, G: 0x34, B: 0x47, A: 0x88}
	quadrantBorder := color.NRGBA{R: 0x00, G: 0xE5, B: 0xFF, A: 0xAA} // Cyan para o centro

	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
			b := canvas.NewRectangle(color.Transparent)
			// Destaca as divisórias dos blocos de 8x8 (quadrantes centrais)
			if x == 7 || y == 7 {
				b.StrokeColor = quadrantBorder
				b.StrokeWidth = 1.5
			} else {
				b.StrokeColor = normalBorder
				b.StrokeWidth = 0.5
			}
			sc.borders[y][x] = b

			rect := canvas.NewRectangle(color.Black)
			sc.rects[y][x] = rect
		}
	}
	return r
}

type spriteCanvasRenderer struct {
	canvas *SpriteCanvas
}

func (r *spriteCanvasRenderer) Destroy() {}

func (r *spriteCanvasRenderer) MinSize() fyne.Size {
	return fyne.NewSize(288, 288)
}

func (r *spriteCanvasRenderer) Layout(size fyne.Size) {
	side := min(size.Width, size.Height)
	cs := side / 16
	r.canvas.cellSize = cs

	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
			pos := fyne.NewPos(float32(x)*cs, float32(y)*cs)
			sz := fyne.NewSize(cs, cs)

			r.canvas.rects[y][x].Move(pos)
			r.canvas.rects[y][x].Resize(sz)

			r.canvas.borders[y][x].Move(pos)
			r.canvas.borders[y][x].Resize(sz)
		}
	}
}

func (r *spriteCanvasRenderer) Refresh() {
	if r.canvas.sprite == nil {
		return
	}

	bg1 := color.NRGBA{R: 0x14, G: 0x1A, B: 0x24, A: 0xFF}
	bg2 := color.NRGBA{R: 0x1C, G: 0x24, B: 0x32, A: 0xFF}

	for y := 0; y < 16; y++ {
		colorIdx := uint8(0x0F)
		if y < len(r.canvas.sprite.ColorBytes) {
			colorIdx = r.canvas.sprite.ColorBytes[y] & 0x0F
		}
		fgColor := GetMSXColor(colorIdx)

		for x := 0; x < 16; x++ {
			on, _ := r.canvas.sprite.GetPixel(x, y)
			if on {
				r.canvas.rects[y][x].FillColor = fgColor
			} else {
				// Fundo xadrez sutil para indicar transparência no MSX 2
				if (x+y)%2 == 0 {
					r.canvas.rects[y][x].FillColor = bg1
				} else {
					r.canvas.rects[y][x].FillColor = bg2
				}
			}
			r.canvas.rects[y][x].Refresh()
			r.canvas.borders[y][x].Refresh()
		}
	}
}

func (r *spriteCanvasRenderer) Objects() []fyne.CanvasObject {
	objs := make([]fyne.CanvasObject, 0, 512)
	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
			objs = append(objs, r.canvas.rects[y][x])
			objs = append(objs, r.canvas.borders[y][x])
		}
	}
	return objs
}

// -------------------------------------------------------------
// Componente Completo de Edição de Sprite (SpriteEditorWidget)
// -------------------------------------------------------------

// SpriteEditorWidget gerencia o editor de sprites Modo 2 com preview, cores e ferramentas.
type SpriteEditorWidget struct {
	state        *ProjectState
	win          fyne.Window
	sprite       *models.Sprite
	canvas       *SpriteCanvas
	preview1x    *canvas.Image
	preview4x    *canvas.Image
	entryName    *widget.Entry
	lblStatus    *widget.Label
	container    *fyne.Container
	colorButtons [16]*widget.Button

	onSaved func()
}

// NewSpriteEditorWidget cria o painel de edição visual do sprite de 16x16.
func NewSpriteEditorWidget(state *ProjectState, win fyne.Window, onSaved func()) *SpriteEditorWidget {
	se := &SpriteEditorWidget{
		state:     state,
		win:       win,
		entryName: widget.NewEntry(),
		lblStatus: widget.NewLabel("Nenhum sprite selecionado."),
		onSaved:   onSaved,
	}

	se.canvas = NewSpriteCanvas(nil, func() {
		se.updatePreviews()
		se.lblStatus.SetText("Alterações no sprite não salvas.")
	})

	// Previews em 1x (16x16) e 4x (64x64)
	se.preview1x = canvas.NewImageFromImage(se.renderImage())
	se.preview1x.SetMinSize(fyne.NewSize(16, 16))
	se.preview1x.ScaleMode = canvas.ImageScalePixels

	se.preview4x = canvas.NewImageFromImage(se.renderImage())
	se.preview4x.SetMinSize(fyne.NewSize(64, 64))
	se.preview4x.ScaleMode = canvas.ImageScalePixels

	// Ferramentas de Pintura
	btnToggle := widget.NewButtonWithIcon("Alternar", theme.ContentUndoIcon(), func() {
		se.canvas.SetTool(PenToggle)
	})
	btnDraw := widget.NewButtonWithIcon("Lápis (Pintar)", theme.DocumentCreateIcon(), func() {
		se.canvas.SetTool(PenDraw)
	})
	btnErase := widget.NewButtonWithIcon("Borracha (Apagar)", theme.ContentClearIcon(), func() {
		se.canvas.SetTool(PenErase)
	})
	toolsBox := container.NewHBox(btnToggle, btnDraw, btnErase)

	// Transformações do Sprite
	btnFlipH := widget.NewButton("Inverter H", func() {
		if se.sprite != nil {
			FlipSpriteHorizontal(se.sprite)
			se.canvas.Refresh()
			se.updatePreviews()
			se.lblStatus.SetText("Sprite espelhado horizontalmente.")
		}
	})
	btnFlipV := widget.NewButton("Inverter V", func() {
		if se.sprite != nil {
			FlipSpriteVertical(se.sprite)
			se.canvas.Refresh()
			se.updatePreviews()
			se.lblStatus.SetText("Sprite espelhado verticalmente.")
		}
	})
	btnRotate := widget.NewButtonWithIcon("Girar 90°", theme.ViewRefreshIcon(), func() {
		if se.sprite != nil {
			RotateSprite90(se.sprite)
			se.canvas.Refresh()
			se.updatePreviews()
			se.lblStatus.SetText("Sprite girado em 90°.")
		}
	})
	btnClear := widget.NewButton("Limpar", func() {
		if se.sprite != nil {
			ClearSprite(se.sprite)
			se.canvas.Refresh()
			se.updatePreviews()
			se.lblStatus.SetText("Sprite limpo (todos pixels transparentes).")
		}
	})
	btnFill := widget.NewButton("Preencher", func() {
		if se.sprite != nil {
			FillSprite(se.sprite)
			se.canvas.Refresh()
			se.updatePreviews()
			se.lblStatus.SetText("Sprite preenchido totalmente.")
		}
	})

	transformGrid := container.NewGridWithColumns(3,
		btnRotate, btnFlipH, btnFlipV,
		btnClear, btnFill,
	)

	// Shift direcional (1 pixel)
	btnUp := widget.NewButton("▲ Cima", func() { se.shift(0, -1) })
	btnDown := widget.NewButton("▼ Baixo", func() { se.shift(0, 1) })
	btnLeft := widget.NewButton("◄ Esq", func() { se.shift(-1, 0) })
	btnRight := widget.NewButton("► Dir", func() { se.shift(1, 0) })
	shiftBox := container.NewHBox(
		widget.NewLabel("Deslocar:"),
		btnLeft, btnRight, btnUp, btnDown,
	)

	// Seletor de Cores por Scanline (16 Linhas Modo 2)
	colorCol1 := container.NewVBox()
	colorCol2 := container.NewVBox()

	for line := 0; line < 16; line++ {
		l := line
		btnColor := widget.NewButton(fmt.Sprintf("L%02d: 15 (Branco)", l), func() {
			se.pickColor(l)
		})
		se.colorButtons[l] = btnColor

		if line < 8 {
			colorCol1.Add(btnColor)
		} else {
			colorCol2.Add(btnColor)
		}
	}

	btnApplyColorAll := widget.NewButton("Aplicar Cor da Linha 0 em Todas", func() {
		if se.sprite == nil || len(se.sprite.ColorBytes) == 0 {
			return
		}
		baseCol := se.sprite.ColorBytes[0]
		for i := range se.sprite.ColorBytes {
			se.sprite.ColorBytes[i] = baseCol
		}
		se.updateColorButtons()
		se.canvas.Refresh()
		se.updatePreviews()
		se.lblStatus.SetText("Cor aplicada a todas as 16 scanlines do Modo 2.")
	})

	colorsSplit := container.NewGridWithColumns(2, colorCol1, colorCol2)
	colorCard := widget.NewCard(
		"Cores por Scanline (Modo 2 do V9938)",
		"Cada uma das 16 linhas pode ter 1 cor independente",
		container.NewVBox(colorsSplit, btnApplyColorAll),
	)

	// Botões de Ação do Sprite
	btnSave := widget.NewButtonWithIcon("Salvar Sprite", theme.DocumentSaveIcon(), func() {
		se.SaveCurrentSprite()
	})
	btnSave.Importance = widget.HighImportance

	btnDuplicate := widget.NewButtonWithIcon("Duplicar Quadro (Novo Frame)", theme.ContentCopyIcon(), func() {
		se.duplicateCurrentSprite()
	})

	// Layout do Editor
	leftCol := container.NewVBox(
		widget.NewLabelWithStyle("Grid Interativo (16x16 pixels)", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		toolsBox,
		container.NewCenter(se.canvas),
		widget.NewSeparator(),
		shiftBox,
		transformGrid,
	)

	previewCard := widget.NewCard("Pré-Visualização", "MSX 2 Real", container.NewHBox(
		container.NewVBox(widget.NewLabel("1x (16x16)"), container.NewCenter(se.preview1x)),
		container.NewVBox(widget.NewLabel("4x (64x64)"), container.NewCenter(se.preview4x)),
	))

	nameForm := widget.NewFormItem("Nome do Sprite:", se.entryName)

	rightCol := container.NewVBox(
		widget.NewForm(nameForm),
		previewCard,
		widget.NewSeparator(),
		colorCard,
		widget.NewSeparator(),
		container.NewHBox(btnSave, btnDuplicate),
		se.lblStatus,
	)

	split := container.NewHSplit(leftCol, rightCol)
	split.Offset = 0.50

	se.container = container.NewPadded(split)
	return se
}

// Container retorna a visualização principal do editor.
func (se *SpriteEditorWidget) Container() fyne.CanvasObject {
	return se.container
}

// LoadSprite carrega um sprite para visualização e edição.
func (se *SpriteEditorWidget) LoadSprite(s *models.Sprite) {
	se.sprite = s
	if s == nil {
		se.canvas.SetSprite(nil)
		se.entryName.SetText("")
		se.lblStatus.SetText("Nenhum sprite selecionado.")
		return
	}

	se.entryName.SetText(s.Name)
	se.canvas.SetSprite(s)
	se.updateColorButtons()
	se.updatePreviews()
	se.lblStatus.SetText(fmt.Sprintf("Editando Sprite #%d: %s", s.ID, s.Name))
}

func (se *SpriteEditorWidget) updateColorButtons() {
	if se.sprite == nil {
		return
	}
	for l := 0; l < 16; l++ {
		colorIdx := uint8(0x0F)
		if l < len(se.sprite.ColorBytes) {
			colorIdx = se.sprite.ColorBytes[l] & 0x0F
		}
		if se.colorButtons[l] != nil {
			se.colorButtons[l].SetText(fmt.Sprintf("L%02d: %d (%s)", l, colorIdx, MSXPalette[colorIdx].Name))
		}
	}
}

func (se *SpriteEditorWidget) pickColor(line int) {
	swatches := make([]fyne.CanvasObject, 16)
	var winDialog dialog.Dialog

	for i := 0; i < 16; i++ {
		colorIdx := uint8(i)
		btn := widget.NewButton(fmt.Sprintf("%d: %s", colorIdx, MSXPalette[colorIdx].Name), func() {
			if se.sprite != nil && line < len(se.sprite.ColorBytes) {
				se.sprite.ColorBytes[line] = colorIdx
				se.updateColorButtons()
				se.canvas.Refresh()
				se.updatePreviews()
				se.lblStatus.SetText(fmt.Sprintf("Scanline %d definida para %s.", line, MSXPalette[colorIdx].Name))
			}
			if winDialog != nil {
				winDialog.Hide()
			}
		})
		swatches[i] = btn
	}

	grid := container.NewGridWithColumns(2, swatches...)
	winDialog = dialog.NewCustom(fmt.Sprintf("Escolha a Cor da Scanline %d", line), "Cancelar", container.NewPadded(grid), se.win)
	winDialog.Resize(fyne.NewSize(380, 420))
	winDialog.Show()
}

func (se *SpriteEditorWidget) shift(dx, dy int) {
	if se.sprite == nil {
		return
	}
	ShiftSprite(se.sprite, dx, dy)
	se.canvas.Refresh()
	se.updatePreviews()
	se.lblStatus.SetText(fmt.Sprintf("Deslocado em X=%d, Y=%d", dx, dy))
}

func (se *SpriteEditorWidget) updatePreviews() {
	img := se.renderImage()
	se.preview1x.Image = img
	se.preview1x.Refresh()
	se.preview4x.Image = img
	se.preview4x.Refresh()
}

func (se *SpriteEditorWidget) renderImage() image.Image {
	img := image.NewNRGBA(image.Rect(0, 0, 16, 16))
	if se.sprite == nil {
		return img
	}

	for y := 0; y < 16; y++ {
		colorIdx := uint8(0x0F)
		if y < len(se.sprite.ColorBytes) {
			colorIdx = se.sprite.ColorBytes[y] & 0x0F
		}
		fg := GetMSXColor(colorIdx)

		for x := 0; x < 16; x++ {
			on, _ := se.sprite.GetPixel(x, y)
			if on {
				img.SetNRGBA(x, y, fg)
			} else {
				img.SetNRGBA(x, y, color.NRGBA{0, 0, 0, 0}) // Transparente no preview
			}
		}
	}
	return img
}

// SaveCurrentSprite persiste as alterações do sprite no SQLite.
func (se *SpriteEditorWidget) SaveCurrentSprite() {
	if se.sprite == nil {
		dialog.ShowInformation("Aviso", "Nenhum sprite selecionado para salvar.", se.win)
		return
	}
	proj := se.state.Current()
	if proj == nil {
		dialog.ShowInformation("Aviso", "Nenhum projeto aberto.", se.win)
		return
	}

	se.sprite.Name = se.entryName.Text
	if strings.TrimSpace(se.sprite.Name) == "" {
		se.sprite.Name = "Sprite_Sem_Nome"
	}

	var err error
	if se.sprite.ID > 0 {
		err = proj.Storage().Sprites.UpdateSprite(context.Background(), se.sprite)
	} else {
		err = proj.Storage().Sprites.CreateSprite(context.Background(), se.sprite)
	}

	if err != nil {
		dialog.ShowError(fmt.Errorf("falha ao salvar sprite: %w", err), se.win)
		return
	}

	se.state.NotifyDataChanged()
	se.lblStatus.SetText(fmt.Sprintf("Sprite '%s' (#%d) salvo com sucesso!", se.sprite.Name, se.sprite.ID))
	if se.onSaved != nil {
		se.onSaved()
	}
}

func (se *SpriteEditorWidget) duplicateCurrentSprite() {
	if se.sprite == nil {
		dialog.ShowInformation("Aviso", "Selecione um sprite para duplicar como próximo quadro.", se.win)
		return
	}
	proj := se.state.Current()
	if proj == nil {
		return
	}

	clone := &models.Sprite{
		Name:         fmt.Sprintf("%s_Frame", se.sprite.Name),
		PatternBytes: make([]byte, models.SpritePatternSize),
		ColorBytes:   make([]byte, models.SpriteColorSize),
	}
	copy(clone.PatternBytes, se.sprite.PatternBytes)
	copy(clone.ColorBytes, se.sprite.ColorBytes)

	err := proj.Storage().Sprites.CreateSprite(context.Background(), clone)
	if err != nil {
		dialog.ShowError(fmt.Errorf("falha ao duplicar sprite: %w", err), se.win)
		return
	}

	se.state.NotifyDataChanged()
	dialog.ShowInformation("Quadro Duplicado", fmt.Sprintf("Novo quadro de animação criado: '%s' (#%d)!", clone.Name, clone.ID), se.win)
	if se.onSaved != nil {
		se.onSaved()
	}
}

// -------------------------------------------------------------
// Funções Utilitárias de Transformação de Sprites (16x16 Modo 2)
// -------------------------------------------------------------

// FlipSpriteHorizontal espelha o padrão horizontalmente (colunas 0..15).
func FlipSpriteHorizontal(s *models.Sprite) {
	for y := 0; y < 16; y++ {
		for x := 0; x < 8; x++ {
			left, _ := s.GetPixel(x, y)
			right, _ := s.GetPixel(15-x, y)
			_ = s.SetPixel(x, y, right)
			_ = s.SetPixel(15-x, y, left)
		}
	}
}

// FlipSpriteVertical espelha o padrão e as cores das scanlines verticalmente.
func FlipSpriteVertical(s *models.Sprite) {
	for y := 0; y < 8; y++ {
		for x := 0; x < 16; x++ {
			top, _ := s.GetPixel(x, y)
			bottom, _ := s.GetPixel(x, 15-y)
			_ = s.SetPixel(x, y, bottom)
			_ = s.SetPixel(x, 15-y, top)
		}
		s.ColorBytes[y], s.ColorBytes[15-y] = s.ColorBytes[15-y], s.ColorBytes[y]
	}
}

// RotateSprite90 rotaciona o padrão do sprite 90° no sentido horário.
func RotateSprite90(s *models.Sprite) {
	temp := models.NewSprite("")
	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
			pixel, _ := s.GetPixel(x, y)
			_ = temp.SetPixel(15-y, x, pixel)
		}
	}
	copy(s.PatternBytes, temp.PatternBytes)
}

// ShiftSprite translada os pixels em dx e dy.
func ShiftSprite(s *models.Sprite, dx, dy int) {
	temp := models.NewSprite("")
	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
			srcX := x - dx
			srcY := y - dy
			if srcX >= 0 && srcX < 16 && srcY >= 0 && srcY < 16 {
				pixel, _ := s.GetPixel(srcX, srcY)
				_ = temp.SetPixel(x, y, pixel)
			}
		}
	}
	copy(s.PatternBytes, temp.PatternBytes)
}

// ClearSprite limpa todos os pixels do sprite.
func ClearSprite(s *models.Sprite) {
	for i := range s.PatternBytes {
		s.PatternBytes[i] = 0x00
	}
}

// FillSprite preenche todos os pixels do sprite.
func FillSprite(s *models.Sprite) {
	for i := range s.PatternBytes {
		s.PatternBytes[i] = 0xFF
	}
}
