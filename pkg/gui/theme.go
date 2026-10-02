package gui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// RetroDarkTheme implementa um tema escuro customizado inspirado na estética do MSX 2 e V9938.
type RetroDarkTheme struct{}

var _ fyne.Theme = (*RetroDarkTheme)(nil)

// Cores da Paleta Retrô MSX 2
var (
	colorMSXBgDark       = color.NRGBA{R: 0x12, G: 0x15, B: 0x1E, A: 0xFF} // Fundo profundo (#12151E)
	colorMSXBgSurface    = color.NRGBA{R: 0x1A, G: 0x20, B: 0x2C, A: 0xFF} // Superfície de cartões (#1A202C)
	colorMSXBgInput      = color.NRGBA{R: 0x14, G: 0x18, B: 0x24, A: 0xFF} // Fundo de inputs (#141824)
	colorMSXSelection    = color.NRGBA{R: 0x22, G: 0x33, B: 0x4A, A: 0xFF} // Seleção (#22334A)
	colorMSXCyan         = color.NRGBA{R: 0x00, G: 0xE5, B: 0xFF, A: 0xFF} // MSX Cyan (#00E5FF)
	colorMSXGreen        = color.NRGBA{R: 0x34, G: 0xD3, B: 0x99, A: 0xFF} // MSX Light Green (#34D399)
	colorMSXAmber        = color.NRGBA{R: 0xFB, G: 0xBF, B: 0x24, A: 0xFF} // MSX Light Yellow (#FBBF24)
	colorMSXRed          = color.NRGBA{R: 0xF8, G: 0x71, B: 0x71, A: 0xFF} // MSX Medium Red (#F87171)
	colorMSXTextLight    = color.NRGBA{R: 0xF1, G: 0xF5, B: 0xF9, A: 0xFF} // Texto primário (#F1F5F9)
	colorMSXTextMuted    = color.NRGBA{R: 0x94, G: 0xA3, B: 0xB8, A: 0xFF} // Texto secundário (#94A3B8)
	colorMSXBorder       = color.NRGBA{R: 0x33, G: 0x41, B: 0x55, A: 0xFF} // Bordas (#334155)
	colorMSXScrollbar    = color.NRGBA{R: 0x47, G: 0x55, B: 0x69, A: 0x80} // Scrollbar translúcida
)

func (t *RetroDarkTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return colorMSXBgDark
	case theme.ColorNameInputBackground:
		return colorMSXBgInput
	case theme.ColorNameMenuBackground, theme.ColorNameOverlayBackground:
		return colorMSXBgSurface
	case theme.ColorNameForeground:
		return colorMSXTextLight
	case theme.ColorNamePlaceHolder:
		return colorMSXTextMuted
	case theme.ColorNamePrimary:
		return colorMSXCyan
	case theme.ColorNameFocus:
		return colorMSXCyan
	case theme.ColorNameSelection:
		return colorMSXSelection
	case theme.ColorNameButton:
		return colorMSXBgSurface
	case theme.ColorNameHover:
		return color.NRGBA{R: 0x2D, G: 0x37, B: 0x48, A: 0xFF}
	case theme.ColorNamePressed:
		return color.NRGBA{R: 0x1E, G: 0x29, B: 0x3B, A: 0xFF}
	case theme.ColorNameError:
		return colorMSXRed
	case theme.ColorNameSuccess:
		return colorMSXGreen
	case theme.ColorNameWarning:
		return colorMSXAmber
	case theme.ColorNameShadow:
		return color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x88}
	case theme.ColorNameScrollBar:
		return colorMSXScrollbar
	case theme.ColorNameSeparator:
		return colorMSXBorder
	default:
		return theme.DefaultTheme().Color(name, theme.VariantDark)
	}
}

func (t *RetroDarkTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (t *RetroDarkTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (t *RetroDarkTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNamePadding:
		return 6
	case theme.SizeNameInlineIcon:
		return 18
	case theme.SizeNameScrollBar:
		return 10
	case theme.SizeNameText:
		return 13
	default:
		return theme.DefaultTheme().Size(name)
	}
}

// NewRetroDarkTheme instancia o tema retrô dark customizado.
func NewRetroDarkTheme() fyne.Theme {
	return &RetroDarkTheme{}
}
