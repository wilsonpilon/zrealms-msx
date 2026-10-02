package gui

import (
	"image"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"

	"github.com/zrealm-msx/zrealm/pkg/models"
)

// RoomTool define a ferramenta ativa de edição no mapa da sala.
type RoomTool int

const (
	ToolBrush RoomTool = iota // Carimbo de tile
	ToolFill                  // Balde de tinta (Flood Fill)
	ToolEyedropper            // Conta-gotas (seleciona o tile clicado)
	ToolEraser                // Borracha (limpa para o tile 0)
)

// RoomCanvas implementa o canvas interativo da sala 32x18 (256x144 pixels).
type RoomCanvas struct {
	widget.BaseWidget
	room          *models.Room
	tilesMap      map[int]*models.Tile
	entities      []*models.Entity
	activeTool    RoomTool
	selectedTile  byte
	onTilePicked  func(tileIdx byte)
	onChanged     func()
	onHoverCoords func(x, y int, tileIdx byte)

	rasterImage *canvas.Image
}

// NewRoomCanvas cria o viewport gráfico da sala 32x18.
func NewRoomCanvas(room *models.Room, tiles []*models.Tile, entities []*models.Entity) *RoomCanvas {
	rc := &RoomCanvas{
		room:         room,
		tilesMap:     make(map[int]*models.Tile),
		entities:     entities,
		activeTool:   ToolBrush,
		selectedTile: 0,
	}
	for _, t := range tiles {
		rc.tilesMap[t.TileIndex] = t
	}

	rc.rasterImage = canvas.NewImageFromImage(rc.renderRoomImage())
	rc.rasterImage.SetMinSize(fyne.NewSize(512, 288)) // Escala 2x nativa
	rc.rasterImage.ScaleMode = canvas.ImageScalePixels

	rc.ExtendBaseWidget(rc)
	return rc
}

// SetRoom associa a sala ativa ao canvas.
func (rc *RoomCanvas) SetRoom(room *models.Room, entities []*models.Entity) {
	rc.room = room
	rc.entities = entities
	rc.updateImage()
}

// SetEntities atualiza a lista de entidades exibidas como marcadores na sala.
func (rc *RoomCanvas) SetEntities(entities []*models.Entity) {
	rc.entities = entities
	rc.updateImage()
}

// GetSelectedTile retorna o tile atualmente selecionado para carimbo.
func (rc *RoomCanvas) GetSelectedTile() byte {
	return rc.selectedTile
}

// GetActiveTool retorna a ferramenta ativa do canvas.
func (rc *RoomCanvas) GetActiveTool() RoomTool {
	return rc.activeTool
}

// RefreshCanvas força a atualização gráfica da matriz.
func (rc *RoomCanvas) RefreshCanvas() {
	rc.updateImage()
}

// Clear preenche a sala inteira com o tile especificado.
func (rc *RoomCanvas) Clear(tile byte) {
	if rc.room == nil {
		return
	}
	ClearRoomMatrix(rc.room.TileMatrix, tile)
	rc.updateImage()
	if rc.onChanged != nil {
		rc.onChanged()
	}
}

// FillBorder preenche as 4 bordas da sala com o tile especificado.
func (rc *RoomCanvas) FillBorder(tile byte) {
	if rc.room == nil {
		return
	}
	FillBorderRoom(rc.room.TileMatrix, tile)
	rc.updateImage()
	if rc.onChanged != nil {
		rc.onChanged()
	}
}

// SetTiles atualiza o catálogo de padrões e cores do tileset para renderização.
func (rc *RoomCanvas) SetTiles(tiles []*models.Tile) {
	rc.tilesMap = make(map[int]*models.Tile)
	for _, t := range tiles {
		rc.tilesMap[t.TileIndex] = t
	}
	rc.updateImage()
}

// SetTool altera a ferramenta ativa de pintura.
func (rc *RoomCanvas) SetTool(tool RoomTool) {
	rc.activeTool = tool
}

// SetSelectedTile define o tile a ser carimbado com o pincel/balde.
func (rc *RoomCanvas) SetSelectedTile(idx byte) {
	rc.selectedTile = idx
}

// OnTilePicked registra callback para quando o conta-gotas for usado.
func (rc *RoomCanvas) OnTilePicked(fn func(tileIdx byte)) {
	rc.onTilePicked = fn
}

// OnChanged registra callback disparado ao alterar a matriz da sala.
func (rc *RoomCanvas) OnChanged(fn func()) {
	rc.onChanged = fn
}

// OnHoverCoords registra callback para exibir coordenadas atuais no mouse.
func (rc *RoomCanvas) OnHoverCoords(fn func(x, y int, tileIdx byte)) {
	rc.onHoverCoords = fn
}

func (rc *RoomCanvas) applyTool(x, y int) {
	if rc.room == nil || x < 0 || x >= models.RoomWidth || y < 0 || y >= models.RoomHeight {
		return
	}

	switch rc.activeTool {
	case ToolBrush:
		_ = rc.room.SetTile(x, y, rc.selectedTile)
	case ToolEraser:
		_ = rc.room.SetTile(x, y, 0)
	case ToolFill:
		FloodFillRoom(rc.room.TileMatrix, x, y, rc.selectedTile)
	case ToolEyedropper:
		cur, _ := rc.room.GetTile(x, y)
		rc.selectedTile = cur
		if rc.onTilePicked != nil {
			rc.onTilePicked(cur)
		}
		return
	}

	rc.updateImage()
	if rc.onChanged != nil {
		rc.onChanged()
	}
}

// Tapped processa cliques individuais no grid da sala.
func (rc *RoomCanvas) Tapped(e *fyne.PointEvent) {
	sz := rc.Size()
	if sz.Width <= 0 || sz.Height <= 0 {
		return
	}
	x := int(e.Position.X / (sz.Width / float32(models.RoomWidth)))
	y := int(e.Position.Y / (sz.Height / float32(models.RoomHeight)))
	rc.applyTool(x, y)
	if rc.onHoverCoords != nil && rc.room != nil {
		cur, _ := rc.room.GetTile(x, y)
		rc.onHoverCoords(x, y, cur)
	}
}

// Dragged processa pintura por arrasto de mouse.
func (rc *RoomCanvas) Dragged(e *fyne.DragEvent) {
	if rc.activeTool == ToolFill || rc.activeTool == ToolEyedropper {
		return // Balde e conta-gotas não são contínuos
	}
	sz := rc.Size()
	if sz.Width <= 0 || sz.Height <= 0 {
		return
	}
	x := int(e.Position.X / (sz.Width / float32(models.RoomWidth)))
	y := int(e.Position.Y / (sz.Height / float32(models.RoomHeight)))
	rc.applyTool(x, y)
	if rc.onHoverCoords != nil && rc.room != nil {
		cur, _ := rc.room.GetTile(x, y)
		rc.onHoverCoords(x, y, cur)
	}
}

func (rc *RoomCanvas) DragEnd() {}

func (rc *RoomCanvas) updateImage() {
	if rc.rasterImage == nil {
		return
	}
	img := rc.renderRoomImage()
	rc.rasterImage.Image = img
	rc.rasterImage.Refresh()
}

func (rc *RoomCanvas) renderRoomImage() image.Image {
	img := image.NewNRGBA(image.Rect(0, 0, models.RoomWidth*models.TileWidth, models.RoomHeight*models.TileHeight))
	if rc.room == nil {
		return img
	}

	// 1. Renderiza os 32 x 18 tiles
	for ry := 0; ry < models.RoomHeight; ry++ {
		for rx := 0; rx < models.RoomWidth; rx++ {
			tileIdx, _ := rc.room.GetTile(rx, ry)
			t := rc.tilesMap[int(tileIdx)]

			var patBytes, colBytes []byte
			if t != nil && len(t.PatternBytes) == 8 && len(t.ColorBytes) == 8 {
				patBytes = t.PatternBytes
				colBytes = t.ColorBytes
			} else {
				patBytes = make([]byte, 8)
				colBytes = []byte{0xF1, 0xF1, 0xF1, 0xF1, 0xF1, 0xF1, 0xF1, 0xF1}
			}

			// Pinta os 8x8 pixels do tile
			for py := 0; py < 8; py++ {
				cb := colBytes[py]
				fg := GetMSXColor((cb >> 4) & 0x0F)
				bg := GetMSXColor(cb & 0x0F)
				pb := patBytes[py]

				for px := 0; px < 8; px++ {
					bit := (pb >> uint(7-px)) & 1
					destX := rx*8 + px
					destY := ry*8 + py
					if bit == 1 {
						img.SetNRGBA(destX, destY, fg)
					} else {
						img.SetNRGBA(destX, destY, bg)
					}
				}
			}
		}
	}

	// 2. Renderiza marcador visual das entidades posicionadas na sala
	entityBorderCol := color.NRGBA{R: 0xFF, G: 0xE6, B: 0x00, A: 0xFF} // Amarelo/Dourado
	entityCenterCol := color.NRGBA{R: 0xFF, G: 0x55, B: 0x55, A: 0xDD} // Vermelho

	for _, ent := range rc.entities {
		if ent.PosX < 0 || ent.PosX >= models.RoomWidth || ent.PosY < 0 || ent.PosY >= models.RoomHeight {
			continue
		}
		startX := ent.PosX * 8
		startY := ent.PosY * 8

		// Desenha borda de destaque ao redor do tile da entidade
		for i := 0; i < 8; i++ {
			img.SetNRGBA(startX+i, startY, entityBorderCol)
			img.SetNRGBA(startX+i, startY+7, entityBorderCol)
			img.SetNRGBA(startX, startY+i, entityBorderCol)
			img.SetNRGBA(startX+7, startY+i, entityBorderCol)
		}
		// Desenha ponto central
		for cy := 2; cy <= 5; cy++ {
			for cx := 2; cx <= 5; cx++ {
				img.SetNRGBA(startX+cx, startY+cy, entityCenterCol)
			}
		}
	}

	return img
}

func (rc *RoomCanvas) CreateRenderer() fyne.WidgetRenderer {
	return &roomCanvasRenderer{
		canvas: rc,
		img:    rc.rasterImage,
	}
}

type roomCanvasRenderer struct {
	canvas *RoomCanvas
	img    *canvas.Image
}

func (r *roomCanvasRenderer) Destroy() {}

func (r *roomCanvasRenderer) MinSize() fyne.Size {
	return fyne.NewSize(512, 288) // 32x18 tiles com 16x16 pixels de visualização
}

func (r *roomCanvasRenderer) Layout(size fyne.Size) {
	r.img.Move(fyne.NewPos(0, 0))
	r.img.Resize(size)
}

func (r *roomCanvasRenderer) Refresh() {
	r.canvas.updateImage()
}

func (r *roomCanvasRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.img}
}

// -------------------------------------------------------------
// Algoritmo de Preenchimento (Flood Fill) para Matriz 32x18
// -------------------------------------------------------------

// FloodFillRoom preenche áreas contíguas do mesmo tipo de tile com o novo tile.
func FloodFillRoom(matrix []byte, startX, startY int, fillTile byte) {
	if len(matrix) != models.RoomMatrixSize {
		return
	}
	startIdx := startY*models.RoomWidth + startX
	target := matrix[startIdx]
	if target == fillTile {
		return
	}

	type point struct{ x, y int }
	queue := []point{{startX, startY}}
	matrix[startIdx] = fillTile

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		neighbors := [][2]int{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}
		for _, n := range neighbors {
			nx := curr.x + n[0]
			ny := curr.y + n[1]
			if nx >= 0 && nx < models.RoomWidth && ny >= 0 && ny < models.RoomHeight {
				idx := ny*models.RoomWidth + nx
				if matrix[idx] == target {
					matrix[idx] = fillTile
					queue = append(queue, point{nx, ny})
				}
			}
		}
	}
}

// FillBorderRoom preenche as 4 bordas da sala 32x18 com um tile especificado (ex: paredes externas).
func FillBorderRoom(matrix []byte, borderTile byte) {
	if len(matrix) != models.RoomMatrixSize {
		return
	}
	// Linhas superior (y = 0) e inferior (y = 17)
	for x := 0; x < models.RoomWidth; x++ {
		matrix[x] = borderTile
		matrix[(models.RoomHeight-1)*models.RoomWidth+x] = borderTile
	}
	// Colunas esquerda (x = 0) e direita (x = 31)
	for y := 0; y < models.RoomHeight; y++ {
		matrix[y*models.RoomWidth] = borderTile
		matrix[y*models.RoomWidth+(models.RoomWidth-1)] = borderTile
	}
}

// ClearRoomMatrix preenche todos os 576 bytes da matriz com o tile especificado.
func ClearRoomMatrix(matrix []byte, fillTile byte) {
	if len(matrix) != models.RoomMatrixSize {
		return
	}
	for i := range matrix {
		matrix[i] = fillTile
	}
}

// AutoConnectRooms conecta automaticamente as salas adjacentes de acordo com suas coordenadas WorldX e WorldY.
// Uma sala em (x, y) conecta-se com:
// - Norte: (x, y-1)
// - Sul:   (x, y+1)
// - Leste: (x+1, y)
// - Oeste: (x-1, y)
// Retorna o total de conexões criadas ou modificadas.
func AutoConnectRooms(rooms []*models.Room) int {
	type coord struct{ x, y int }
	lookup := make(map[coord]*models.Room)
	for _, r := range rooms {
		lookup[coord{r.WorldX, r.WorldY}] = r
	}

	changes := 0
	for _, r := range rooms {
		// Norte: (x, y-1)
		if n, ok := lookup[coord{r.WorldX, r.WorldY - 1}]; ok {
			if r.NorthRoomID == nil || *r.NorthRoomID != n.ID {
				nid := n.ID
				r.NorthRoomID = &nid
				changes++
			}
		} else if r.NorthRoomID != nil {
			r.NorthRoomID = nil
			changes++
		}

		// Sul: (x, y+1)
		if s, ok := lookup[coord{r.WorldX, r.WorldY + 1}]; ok {
			if r.SouthRoomID == nil || *r.SouthRoomID != s.ID {
				sid := s.ID
				r.SouthRoomID = &sid
				changes++
			}
		} else if r.SouthRoomID != nil {
			r.SouthRoomID = nil
			changes++
		}

		// Leste: (x+1, y)
		if e, ok := lookup[coord{r.WorldX + 1, r.WorldY}]; ok {
			if r.EastRoomID == nil || *r.EastRoomID != e.ID {
				eid := e.ID
				r.EastRoomID = &eid
				changes++
			}
		} else if r.EastRoomID != nil {
			r.EastRoomID = nil
			changes++
		}

		// Oeste: (x-1, y)
		if w, ok := lookup[coord{r.WorldX - 1, r.WorldY}]; ok {
			if r.WestRoomID == nil || *r.WestRoomID != w.ID {
				wid := w.ID
				r.WestRoomID = &wid
				changes++
			}
		} else if r.WestRoomID != nil {
			r.WestRoomID = nil
			changes++
		}
	}

	return changes
}
