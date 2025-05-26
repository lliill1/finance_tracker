package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type customTheme struct {
	fyne.Theme
	isDark bool
}

func newCustomTheme(dark bool) *customTheme {
	return &customTheme{
		Theme:  theme.DefaultTheme(),
		isDark: dark,
	}
}

func (t *customTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if t.isDark {
		return theme.DarkTheme().Color(name, variant)
	}

	switch name {
	case theme.ColorNameBackground:
		return color.NRGBA{R: 255, G: 240, B: 245, A: 255} // Светло-розовый фон
	case theme.ColorNameForeground:
		return color.NRGBA{R: 75, G: 0, B: 130, A: 255} // Темно-фиолетовый текст
	case theme.ColorNamePrimary:
		return color.NRGBA{R: 219, G: 112, B: 147, A: 255} // Розовый для акцентов
	case theme.ColorNameHover:
		return color.NRGBA{R: 255, G: 192, B: 203, A: 255} // Светло-розовый при наведении
	case theme.ColorNamePressed:
		return color.NRGBA{R: 199, G: 21, B: 133, A: 255} // Темно-розовый при нажатии
	case theme.ColorNameScrollBar:
		return color.NRGBA{R: 255, G: 182, B: 193, A: 255} // Розовый для скроллбара
	case theme.ColorNameShadow:
		return color.NRGBA{R: 255, G: 192, B: 203, A: 128} // Полупрозрачный розовый для теней
	default:
		return theme.LightTheme().Color(name, variant)
	}
}

func (t *customTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

func (t *customTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (t *customTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
} 