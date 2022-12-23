package main

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// TODO better names
	style           lipgloss.Style
	styleFirst      lipgloss.Style
	styleFirstSharp lipgloss.Style
	styleLastSharp  lipgloss.Style
	styleOmitLeft   lipgloss.Style
	styleOmitRight  lipgloss.Style
	styleSharp      lipgloss.Style
	styleReverse    lipgloss.Style
)

func init() {
	const (
		COLOR_IVORY       = "#fffff0"
		COLOR_BLACK       = "#202020"
		COLOR_WHITE_SMOKE = "#eddcc9"
		COLOR_LAVA_SMOKE  = "#5e6064"
	)

	colorWhite := lipgloss.AdaptiveColor{Light: COLOR_IVORY, Dark: COLOR_IVORY}
	colorBlack := lipgloss.AdaptiveColor{Light: COLOR_BLACK, Dark: COLOR_BLACK}
	// colorBlackBorder := lipgloss.AdaptiveColor{Light: COLOR_WHITE_SMOKE, Dark: COLOR_WHITE_SMOKE}
	colorBorder := lipgloss.AdaptiveColor{Light: COLOR_LAVA_SMOKE, Dark: COLOR_LAVA_SMOKE}

	border := lipgloss.ThickBorder()
	borderFirstWhite, borderFirstSharp := border, border
	border.BottomLeft = "┻"
	borderSharp := border
	borderSharp.BottomLeft = "┻"

	noTop := func(style lipgloss.Style) lipgloss.Style {
		return style.
			BorderTop(false).
			BorderLeft(true).
			BorderRight(false).
			BorderBottom(true)
	}

	style = noTop(lipgloss.NewStyle().
		BorderStyle(border).
		BorderForeground(colorBorder).
		BorderBackground(colorWhite).
		Background(colorWhite).
		Foreground(colorBlack).
		Align())

	styleFirst = noTop(style.Copy().Border(borderFirstWhite))

	styleSharp = noTop(style.Copy().
		BorderStyle(borderSharp).
		BorderForeground(colorBorder).
		BorderBackground(colorBlack).
		Background(colorBlack).
		Foreground(colorWhite).
		Align())

	styleFirstSharp = noTop(styleSharp.Copy().Border(borderFirstSharp))
	styleLastSharp = noTop(styleSharp.Copy().Border(borderSharp)).BorderRight(true).Align()

	omitBorder := lipgloss.HiddenBorder()
	omitBorder.Top = ""
	omitBorder.TopRight = ""
	omitBorder.TopLeft = ""
	omitBorderLeft, omitBorderRight := omitBorder, omitBorder
	omitBorderLeft.Left = style.GetBorderStyle().Left
	omitBorderLeft.BottomLeft = style.GetBorderStyle().Left
	omitBorderRight.Right = style.GetBorderStyle().Right
	omitBorderRight.BottomRight = style.GetBorderStyle().Right

	styleOmitLeft = noTop(style.Copy().Border(omitBorderLeft))
	styleOmitRight = noTop(style.Copy().Border(omitBorderRight))

	styleReverse = lipgloss.NewStyle().Reverse(true)
}
