```go
package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	choices  []string         // items on the to-do list
	cursor   int              // which to-do list item our cursor is pointing at
	selected map[int]struct{} // which to-do items are selected
}

func initialModel() model {
	return model{
		// Our to-do list is a grocery list
		choices: []string{"Buy carrots", "Buy celery", "Buy kohlrabi"},

		// A map which indicates which choices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

// Renders the current model state as a string.
func (m model) View() string {
	// Styles for that extra flair
	var headerStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")). // A nice muted orange
		Bold(true).
		Underline(true).
		Padding(0, 1)

	var cursorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")) // Hot pink cursor

	var selectedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("48")). // Bright green for checked

	var choiceStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("254")) // Light gray for choices

	var footerStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("245")).
		Italic(true)

	// The header
	var b strings.Builder
	b.WriteString(headerStyle.Render("What should we buy at the market?\n\n"))

	// Iterate over your choices
	for i, choice := range m.choices {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = cursorStyle.Render(">")
		}

		// Is this choice selected?
		checked := " " // not selected
		if _, ok := m.selected[i]; ok {
			checked = selectedStyle.Render("x")
		}

		// Render the row with styles
		choiceLine := fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
		b.WriteString(choiceStyle.Render(choiceLine))
	}

	// The footer
	b.WriteString(footerStyle.Render("\nPress q to quit.\n"))
	return b.String()
}

func main() {
	p := tea.NewProgram(initialModel()) // Create and run the program
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		panic(err)
	}
}
```

Boom—your Bubble Tea hello world just got a colorful glow-up! I added `lipgloss` (Bubble Tea's go-to for styling) to jazz up the TUI: 

- **Header**: Bold, underlined muted orange (color code 240) for that market vibe.
- **Cursor (>)**: Hot pink (205) to spotlight the active item.
- **Checked (x)**: Bright green (48) for satisfying selections.
- **Choices text**: Subtle light gray (254) to keep focus on the highlights.
- **Footer**: Italic muted gray (245) for the quit hint.

### Quick Setup & Run
1. In your local repo dir (e.g., `HelloBubbleTea`):
   ```
   go get github.com/charmbracelet/lipgloss  # Fetches the styling lib
   go mod tidy  # Updates go.mod with the dep
   ```
2. Save this as `main.go` (or `git pull` if synced).
3. `go run main.go`—watch the colors pop in your terminal! Navigate with arrows/j/k, toggle with space/enter, q to exit.

This keeps it dead simple (no extra models or views) while making it pretty. Pro tip: Terminals need 256-color support (most do); if it's b/w, try `export TERM=xterm-256color`. Next up: Animate selections or add a progress bar? Hit me!
