package terminal

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/term"
)

const (
	defaultWidth       = 100
	minWidth           = 40
	cursorQueryTimeout = 500 * time.Millisecond
)

// GetWidth returns the current terminal width, falling back to defaultWidth.
func GetWidth() int {
	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || w < minWidth {
		return defaultWidth
	}
	return w
}

// GetCursorRow queries the terminal for the current cursor row using ANSI DSR.
// Returns 0 if the query fails or times out.
func GetCursorRow() int {
	fd := int(os.Stdin.Fd())

	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return 0
	}
	defer term.Restore(fd, oldState)

	// Send Device Status Report request
	fmt.Print("\033[6n")

	type result struct {
		row int
	}

	ch := make(chan result, 1)
	go func() {
		reader := bufio.NewReader(os.Stdin)

		// Skip until ESC
		for {
			b, err := reader.ReadByte()
			if err != nil {
				ch <- result{0}
				return
			}
			if b == '\033' {
				break
			}
		}

		b, err := reader.ReadByte()
		if err != nil || b != '[' {
			ch <- result{0}
			return
		}

		// Parse row number from "row;colR" response
		row := 0
		for {
			b, err = reader.ReadByte()
			if err != nil {
				ch <- result{0}
				return
			}
			if b == ';' {
				break
			}
			if b >= '0' && b <= '9' {
				row = row*10 + int(b-'0')
			}
		}

		// Consume remaining bytes until 'R'
		for {
			b, err = reader.ReadByte()
			if err != nil {
				ch <- result{0}
				return
			}
			if b == 'R' {
				break
			}
		}

		ch <- result{row}
	}()

	select {
	case r := <-ch:
		return r.row
	case <-time.After(cursorQueryTimeout):
		return 0
	}
}

// StripAnsi removes ANSI escape sequences from a string.
func StripAnsi(s string) string {
	var result strings.Builder
	inEscape := false
	for i := 0; i < len(s); i++ {
		if s[i] == '\033' {
			inEscape = true
			continue
		}
		if inEscape {
			if (s[i] >= 'A' && s[i] <= 'Z') || (s[i] >= 'a' && s[i] <= 'z') {
				inEscape = false
			}
			continue
		}
		result.WriteByte(s[i])
	}
	return result.String()
}
