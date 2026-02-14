package ui

import "github.com/charmbracelet/lipgloss"

var (
	// Claude theme color: rgb(215,119,87)
	colorClaude    = lipgloss.Color("#D77757")
	colorDimGray   = lipgloss.Color("#666666")
	colorLightGray = lipgloss.Color("#999999")
	colorWhite     = lipgloss.Color("#FFFFFF")
	colorGreen     = lipgloss.Color("#5AF78E")

	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorClaude).
			PaddingLeft(1).PaddingRight(1).
			PaddingTop(1).PaddingBottom(1)

	boldWhiteStyle = lipgloss.NewStyle().
			Foreground(colorWhite).
			Bold(true)

	logoStyle = lipgloss.NewStyle().
			Foreground(colorClaude).
			Bold(true)

	mutedStyle = lipgloss.NewStyle().
			Foreground(colorLightGray)

	dimStyle = lipgloss.NewStyle().
			Foreground(colorDimGray)

	greenStyle = lipgloss.NewStyle().
			Foreground(colorGreen).
			Bold(true)
)
