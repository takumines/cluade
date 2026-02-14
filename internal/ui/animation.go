package ui

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

const (
	gravity         = 0.4
	frameInterval   = 30 * time.Millisecond
	fallStaggerStep = 2 // frames of delay between each character's start
	clearPause      = 200 * time.Millisecond
)

// fallingChar represents a single character that participates in the falling animation.
type fallingChar struct {
	char       string
	charWidth  int
	startRow   int
	startCol   int
	currentY   float64
	velocity   float64
	maxRow     int
	landed     bool
	startDelay int // frames to wait before this character begins falling
}

// runFallingAnimation makes each logo character fall with gravity,
// then clears them and displays the joke overlay.
func runFallingAnimation(tuiStartRow, logoRowOffset, logoColOffset, tuiTotalLines int) {
	logoLines := LogoLines()
	logoHeight := len(logoLines)

	// Characters fall only within the logo area to avoid overwriting content below.
	maxFallRow := tuiStartRow + logoRowOffset + logoHeight - 1

	var chars []fallingChar
	delayCounter := 0

	for rowIdx, line := range logoLines {
		col := 0
		for _, r := range line {
			ch := string(r)
			w := lipgloss.Width(ch)
			if ch != " " {
				absRow := tuiStartRow + logoRowOffset + rowIdx
				absCol := logoColOffset + col + 1 // terminal columns are 1-based

				chars = append(chars, fallingChar{
					char:       ch,
					charWidth:  w,
					startRow:   absRow,
					startCol:   absCol,
					currentY:   float64(absRow),
					velocity:   0,
					maxRow:     maxFallRow,
					landed:     false,
					startDelay: delayCounter * fallStaggerStep,
				})
				delayCounter++
			}
			col += w
		}
	}

	fmt.Print("\033[?25l")         // hide cursor
	defer fmt.Print("\033[?25h")   // restore cursor on exit

	frame := 0
	for {
		allLanded := true
		for i := range chars {
			c := &chars[i]
			if c.landed {
				continue
			}
			if frame < c.startDelay {
				allLanded = false
				continue
			}

			allLanded = false

			oldRow := int(math.Round(c.currentY))
			fmt.Printf("\033[%d;%dH%s", oldRow, c.startCol, strings.Repeat(" ", c.charWidth))

			c.velocity += gravity
			c.currentY += c.velocity

			newRow := int(math.Round(c.currentY))
			if newRow >= c.maxRow {
				newRow = c.maxRow
				c.landed = true
			}

			fmt.Printf("\033[%d;%dH%s", newRow, c.startCol, logoStyle.Render(c.char))
		}

		if allLanded {
			break
		}

		time.Sleep(frameInterval)
		frame++
	}

	// Clear landed characters
	time.Sleep(clearPause)
	for _, c := range chars {
		fmt.Printf("\033[%d;%dH%s", c.maxRow, c.startCol, strings.Repeat(" ", c.charWidth))
	}

	drawJokeOverlay(tuiStartRow, logoRowOffset, tuiTotalLines)
}
