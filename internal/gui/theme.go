//go:build gui

package gui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// App-specific color names, resolved by appTheme like fyne's own.
const (
	colorSurface       fyne.ThemeColorName = "al.surface"
	colorSurfaceBorder fyne.ThemeColorName = "al.surfaceBorder"
	colorMuted         fyne.ThemeColorName = "al.muted"
	colorSidebar       fyne.ThemeColorName = "al.sidebar"
	colorInfo          fyne.ThemeColorName = "al.info"
)

// palette maps each color name to its {dark, light} value.
var palette = map[fyne.ThemeColorName][2]color.NRGBA{
	theme.ColorNameBackground:          {rgb(0x0B0D12, 0xFF), rgb(0xF5F6FA, 0xFF)},
	theme.ColorNameHeaderBackground:    {rgb(0x12151C, 0xFF), rgb(0xFFFFFF, 0xFF)},
	theme.ColorNameMenuBackground:      {rgb(0x161A23, 0xFF), rgb(0xFFFFFF, 0xFF)},
	theme.ColorNameOverlayBackground:   {rgb(0x12151C, 0xFF), rgb(0xFFFFFF, 0xFF)},
	theme.ColorNameInputBackground:     {rgb(0x161A23, 0xFF), rgb(0xFFFFFF, 0xFF)},
	theme.ColorNameInputBorder:         {rgb(0x2A3040, 0xFF), rgb(0xD5D9E3, 0xFF)},
	theme.ColorNameButton:              {rgb(0x1C2130, 0xFF), rgb(0xECEEF4, 0xFF)},
	theme.ColorNameDisabledButton:      {rgb(0x151922, 0xFF), rgb(0xF0F1F5, 0xFF)},
	theme.ColorNameForeground:          {rgb(0xE6E8EE, 0xFF), rgb(0x161922, 0xFF)},
	theme.ColorNameDisabled:            {rgb(0x5B6170, 0xFF), rgb(0xA3A8B5, 0xFF)},
	theme.ColorNamePlaceHolder:         {rgb(0x7A8193, 0xFF), rgb(0x8A90A0, 0xFF)},
	theme.ColorNamePrimary:             {rgb(0x7C6CFF, 0xFF), rgb(0x5B4BFF, 0xFF)},
	theme.ColorNameForegroundOnPrimary: {rgb(0xFFFFFF, 0xFF), rgb(0xFFFFFF, 0xFF)},
	theme.ColorNameFocus:               {rgb(0x7C6CFF, 0x99), rgb(0x5B4BFF, 0x80)},
	theme.ColorNameHover:               {rgb(0xFFFFFF, 0x0F), rgb(0x000000, 0x0A)},
	theme.ColorNamePressed:             {rgb(0xFFFFFF, 0x1A), rgb(0x000000, 0x14)},
	theme.ColorNameSelection:           {rgb(0x7C6CFF, 0x40), rgb(0x5B4BFF, 0x2E)},
	theme.ColorNameSeparator:           {rgb(0x1F2430, 0xFF), rgb(0xE3E6ED, 0xFF)},
	theme.ColorNameScrollBar:           {rgb(0xFFFFFF, 0x33), rgb(0x000000, 0x33)},
	theme.ColorNameScrollBarBackground: {rgb(0x000000, 0x00), rgb(0x000000, 0x00)},
	theme.ColorNameShadow:              {rgb(0x000000, 0x66), rgb(0x000000, 0x22)},
	theme.ColorNameHyperlink:           {rgb(0x9D92FF, 0xFF), rgb(0x4A3AF0, 0xFF)},
	theme.ColorNameSuccess:             {rgb(0x22C55E, 0xFF), rgb(0x16A34A, 0xFF)},
	theme.ColorNameWarning:             {rgb(0xF59E0B, 0xFF), rgb(0xD97706, 0xFF)},
	theme.ColorNameError:               {rgb(0xF43F5E, 0xFF), rgb(0xE11D48, 0xFF)},
	colorSurface:                       {rgb(0x12151C, 0xFF), rgb(0xFFFFFF, 0xFF)},
	colorSurfaceBorder:                 {rgb(0x222836, 0xFF), rgb(0xE3E6ED, 0xFF)},
	colorMuted:                         {rgb(0x8B92A5, 0xFF), rgb(0x6B7182, 0xFF)},
	colorSidebar:                       {rgb(0x0F1218, 0xFF), rgb(0xFFFFFF, 0xFF)},
	colorInfo:                          {rgb(0x38BDF8, 0xFF), rgb(0x0284C7, 0xFF)},
}

func rgb(v uint32, alpha uint8) color.NRGBA {
	return color.NRGBA{R: uint8(v >> 16), G: uint8(v >> 8), B: uint8(v), A: alpha} //nolint:gosec // byte extraction
}

// appTheme is the app's palette and geometry over fyne's default fonts and icons. The variant is
// the app's own (sidebar toggle), not the desktop's: under sudo the root session has no portal to
// report one.
type appTheme struct {
	variant fyne.ThemeVariant
}

func (t *appTheme) Color(name fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	if c, ok := palette[name]; ok {
		if t.variant == theme.VariantLight {
			return c[1]
		}
		return c[0]
	}
	return theme.DefaultTheme().Color(name, t.variant)
}

func (t *appTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (t *appTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (t *appTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNameInputRadius, theme.SizeNameButtonRadius, theme.SizeNameSelectionRadius:
		return 8
	case theme.SizeNameCardRadius, theme.SizeNameDialogRadius, theme.SizeNamePopupRadius:
		return 12
	case theme.SizeNameMenuRadius:
		return 10
	case theme.SizeNameScrollBarRadius:
		return 4
	case theme.SizeNameHeadingText:
		return 26
	case theme.SizeNameSubHeadingText:
		return 17
	case theme.SizeNameCaptionText:
		return 12
	case theme.SizeNameSeparatorThickness:
		return 1
	case theme.SizeNameInputBorder:
		return 1
	case theme.SizeNamePadding:
		return 6
	}
	return theme.DefaultTheme().Size(name)
}

// themeColor resolves name against the running app's theme.
func themeColor(name fyne.ThemeColorName) color.Color {
	return fyne.CurrentApp().Settings().Theme().Color(name, fyne.CurrentApp().Settings().ThemeVariant())
}

// tint is c at alpha a (0-255), for pill and icon-tile backgrounds.
func tint(c color.Color, a uint8) color.Color {
	n := color.NRGBAModel.Convert(c).(color.NRGBA)
	n.A = a
	return n
}
