package gui

import "image/color"

// MSXColor representa uma entrada na paleta oficial de 16 cores do MSX (V9938 / TMS9918).
type MSXColor struct {
	Index uint8
	Name  string
	Hex   string
	Color color.NRGBA
}

// MSXPalette contém as 16 cores padrão do hardware V9938 no modo Graphic 3 / SCREEN 4.
var MSXPalette = [16]MSXColor{
	{Index: 0, Name: "Transparente", Hex: "#000000", Color: color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x00}},
	{Index: 1, Name: "Preto", Hex: "#000000", Color: color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}},
	{Index: 2, Name: "Verde Médio", Hex: "#21C842", Color: color.NRGBA{R: 0x21, G: 0xC8, B: 0x42, A: 0xFF}},
	{Index: 3, Name: "Verde Claro", Hex: "#5CE471", Color: color.NRGBA{R: 0x5C, G: 0xE4, B: 0x71, A: 0xFF}},
	{Index: 4, Name: "Azul Escuro", Hex: "#5455ED", Color: color.NRGBA{R: 0x54, G: 0x55, B: 0xED, A: 0xFF}},
	{Index: 5, Name: "Azul Claro", Hex: "#7D70FF", Color: color.NRGBA{R: 0x7D, G: 0x70, B: 0xFF, A: 0xFF}},
	{Index: 6, Name: "Vermelho Escuro", Hex: "#D05441", Color: color.NRGBA{R: 0xD0, G: 0x54, B: 0x41, A: 0xFF}},
	{Index: 7, Name: "Ciano", Hex: "#42EBF5", Color: color.NRGBA{R: 0x42, G: 0xEB, B: 0xF5, A: 0xFF}},
	{Index: 8, Name: "Vermelho Médio", Hex: "#FF5555", Color: color.NRGBA{R: 0xFF, G: 0x55, B: 0x55, A: 0xFF}},
	{Index: 9, Name: "Vermelho Claro", Hex: "#FF7978", Color: color.NRGBA{R: 0xFF, G: 0x79, B: 0x78, A: 0xFF}},
	{Index: 10, Name: "Amarelo Escuro", Hex: "#D4C154", Color: color.NRGBA{R: 0xD4, G: 0xC1, B: 0x54, A: 0xFF}},
	{Index: 11, Name: "Amarelo Claro", Hex: "#E6CE80", Color: color.NRGBA{R: 0xE6, G: 0xCE, B: 0x80, A: 0xFF}},
	{Index: 12, Name: "Verde Escuro", Hex: "#21B03B", Color: color.NRGBA{R: 0x21, G: 0xB0, B: 0x3B, A: 0xFF}},
	{Index: 13, Name: "Magenta", Hex: "#C95BBA", Color: color.NRGBA{R: 0xC9, G: 0x5B, B: 0xBA, A: 0xFF}},
	{Index: 14, Name: "Cinza", Hex: "#CCCCCC", Color: color.NRGBA{R: 0xCC, G: 0xCC, B: 0xCC, A: 0xFF}},
	{Index: 15, Name: "Branco", Hex: "#FFFFFF", Color: color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}},
}

// GetMSXColor retorna a cor correspondente ao índice de 0 a 15 da paleta do V9938.
func GetMSXColor(idx uint8) color.NRGBA {
	if idx > 15 {
		idx = 15
	}
	// Se for índice 0 (transparente), na edição em tela exibe preto profundo com alpha 255 para visualização
	if idx == 0 {
		return color.NRGBA{R: 0x08, G: 0x08, B: 0x08, A: 0xFF}
	}
	return MSXPalette[idx].Color
}
