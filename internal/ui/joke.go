package ui

import (
	"fmt"
	"strings"
	"time"
)

const (
	jokeStartCol = 4 // left edge of left panel: border(1) + padding(1) + margin(2)
)

// jokeMessage defines a single line of the joke with its display delay.
type jokeMessage struct {
	text  string
	delay time.Duration
}

func jokeMessages() []jokeMessage {
	return []jokeMessage{
		{logoStyle.Render("Typo detected!"), 0},
		{mutedStyle.Render("cluade") + " → " + greenStyle.Render("claude"), 300 * time.Millisecond},
		{dimStyle.Render("Check your spelling ;)"), 300 * time.Millisecond},
	}
}

// drawJokeOverlay writes joke messages directly to the terminal at the logo position.
func drawJokeOverlay(tuiStartRow, logoRowOffset, tuiTotalLines int) {
	row := tuiStartRow + logoRowOffset
	msgs := jokeMessages()

	for i, msg := range msgs {
		time.Sleep(msg.delay)
		fmt.Printf("\033[%d;%dH%s", row+i, jokeStartCol, msg.text)
	}

	// Move cursor below the TUI
	fmt.Printf("\033[%d;1H\n", tuiStartRow+tuiTotalLines)
}

// renderJokeFallback returns a plain-text version of the joke for non-interactive terminals.
func renderJokeFallback() string {
	msgs := jokeMessages()

	var out strings.Builder
	out.WriteString("\n")
	for _, msg := range msgs {
		out.WriteString("  ")
		out.WriteString(msg.text)
		out.WriteString("\n")
	}
	out.WriteString("\n")
	return out.String()
}
