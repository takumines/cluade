package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/takumines/cluade/internal/system"
	"github.com/takumines/cluade/internal/terminal"
)

// RunStartupScreen displays the welcome TUI, waits, then runs the joke animation.
func RunStartupScreen() {
	rendered := renderTUI()
	fmt.Print(rendered)

	tuiLines := strings.Split(rendered, "\n")
	tuiStartRow, logoRowOffset, logoColOffset := findLogoPosition(tuiLines)

	time.Sleep(3 * time.Second)

	if tuiStartRow > 0 && logoRowOffset > 0 {
		runFallingAnimation(tuiStartRow, logoRowOffset, logoColOffset, len(tuiLines))
	} else {
		fmt.Print(renderJokeFallback())
	}
}

func renderTUI() string {
	username := system.GetUsername()
	currentDir := system.GetCurrentDir()

	panels := buildPanels(username, currentDir)
	bordered := borderStyle.Render(panels)

	title := "Cluade Code " + version()
	bordered = embedBorderTitle(bordered, title)

	return lipgloss.NewStyle().Padding(1, 0).Render(bordered)
}

// findLogoPosition scans the rendered TUI to locate the first logo line.
// Returns (tuiStartRow, logoRowOffset, logoColOffset).
func findLogoPosition(tuiLines []string) (int, int, int) {
	firstLogoLine := LogoLines()[0]
	totalLines := len(tuiLines)

	cursorRow := terminal.GetCursorRow()
	if cursorRow <= 0 {
		return 0, 0, 0
	}

	tuiStartRow := cursorRow - totalLines + 1
	if tuiStartRow < 1 {
		tuiStartRow = 1
	}

	logoRowOffset := 0
	logoColOffset := 0
	for i, line := range tuiLines {
		plain := terminal.StripAnsi(line)
		if plainIdx := strings.Index(plain, firstLogoLine); plainIdx >= 0 {
			logoRowOffset = i
			logoColOffset = plainIdx
			break
		}
	}

	return tuiStartRow, logoRowOffset, logoColOffset
}
