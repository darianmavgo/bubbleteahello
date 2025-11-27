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
		Foreground(lipgloss.Color("240")).
		Bold(true).
		Underline(true).
		Padding(0, 1)

	var cursorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("205"))

	var selectedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("48"))

	var choiceStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("254"))

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

		// Render the row with styles applied to parts
		cursorPart := cursor
		checkedPart := fmt.Sprintf("[%s] ", checked)
		choicePart := choiceStyle.Render(choice)
		b.WriteString(fmt.Sprintf("%s%s%s\n", cursorPart, checkedPart, choicePart))
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
