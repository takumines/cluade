package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/takumines/cluade/internal/terminal"
)

const (
	leftPanelWidth     = 27
	rightPanelMinWidth = 20
	borderTitleOffset  = 3 // number of dashes between ╭ and the title
)

func version() string {
	now := time.Now()
	return fmt.Sprintf("v%d.%d.%d", now.Year()%100, int(now.Month()), now.Day())
}

func buildLeftPanel(username, currentDir string) string {
	welcome := "Welcome back!"
	if username != "" && len(username) <= 20 {
		welcome = fmt.Sprintf("Welcome back %s!", username)
	}

	var logoRendered []string
	for _, line := range LogoLines() {
		centered := lipgloss.NewStyle().Width(leftPanelWidth).Align(lipgloss.Center).Render(
			logoStyle.Render(line),
		)
		logoRendered = append(logoRendered, centered)
	}

	modelLine := dimStyle.Render("Opus 1000 · Cluade Max")
	cwdLine := dimStyle.Render(currentDir)

	var b strings.Builder
	b.WriteString(boldWhiteStyle.Render(welcome))
	b.WriteString("\n\n\n")
	for i, line := range logoRendered {
		b.WriteString(line)
		if i < len(logoRendered)-1 {
			b.WriteString("\n")
		}
	}
	b.WriteString("\n")
	b.WriteString(modelLine)
	b.WriteString("\n")
	b.WriteString(cwdLine)

	return lipgloss.NewStyle().Width(leftPanelWidth).Align(lipgloss.Center).Render(b.String())
}

func buildRightPanel(innerWidth int) string {
	var b strings.Builder

	b.WriteString(boldWhiteStyle.Render("Tips for not misspelling"))
	b.WriteString("\n")
	b.WriteString(mutedStyle.Render("Run ") + boldWhiteStyle.Render("/spellcheck") + mutedStyle.Render(" to verify you typed 'claude' not 'cluade'"))
	b.WriteString("\n")

	sep := dimStyle.Render(strings.Repeat("─", innerWidth))
	b.WriteString(sep)
	b.WriteString("\n")

	b.WriteString(boldWhiteStyle.Render("Recent typos"))
	b.WriteString("\n")
	b.WriteString(dimStyle.Render("cluade, cluade, cluade, cluade, cluade..."))

	return lipgloss.NewStyle().Width(innerWidth).Render(b.String())
}

func buildPanels(username, currentDir string) string {
	left := buildLeftPanel(username, currentDir)

	termWidth := terminal.GetWidth()
	// border(2) + paddingX(2) = 4, divider with spaces(" │ ") = 3
	innerWidth := termWidth - 4
	rightWidth := innerWidth - leftPanelWidth - 3
	if rightWidth < rightPanelMinWidth {
		rightWidth = rightPanelMinWidth
	}
	right := buildRightPanel(rightWidth)

	leftHeight := lipgloss.Height(left)
	rightHeight := lipgloss.Height(right)
	maxHeight := leftHeight
	if rightHeight > maxHeight {
		maxHeight = rightHeight
	}
	if maxHeight < 1 {
		maxHeight = 1
	}

	dividerLines := make([]string, maxHeight)
	for i := range dividerLines {
		dividerLines[i] = " " + dimStyle.Render("│") + " "
	}
	divider := strings.Join(dividerLines, "\n")

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		left,
		divider,
		right,
	)
}

// embedBorderTitle replaces the top border line with a title embedded in it.
// Produces: ╭─── Cluade Code v1.0.0 ───...───╮
func embedBorderTitle(rendered, title string) string {
	lines := strings.Split(rendered, "\n")
	if len(lines) == 0 {
		return rendered
	}

	topLine := lines[0]
	totalVisibleWidth := lipgloss.Width(topLine)
	if totalVisibleWidth < 4 {
		return rendered
	}

	titleContent := " " + title + " "
	titleLen := lipgloss.Width(titleContent)

	remainingDashes := totalVisibleWidth - 2 - borderTitleOffset - titleLen
	if remainingDashes < 0 {
		remainingDashes = 0
	}

	borderColorStyle := lipgloss.NewStyle().Foreground(colorClaude)
	newTop := borderColorStyle.Render("╭"+strings.Repeat("─", borderTitleOffset)) +
		borderColorStyle.Render(titleContent) +
		borderColorStyle.Render(strings.Repeat("─", remainingDashes)+"╮")

	lines[0] = newTop

	return strings.Join(lines, "\n")
}
