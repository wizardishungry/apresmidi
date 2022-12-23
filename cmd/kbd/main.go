package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

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

func (m model) View() string {
	// The header
	s := "What should we buy at the market?\n\n"

	// Iterate over our choices
	for i, choice := range m.choices {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if _, ok := m.selected[i]; ok {
			checked = "x" // selected!
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	// The footer
	s += "\nPress q to quit.\n"

	const (
		maxNote    = 127
		initOctave = -2
	)

	type note struct {
		name  string
		sharp bool
	}

	notes := []note{
		{"C", false},
		{"C", true},
		{"D", false},
		{"D", true},
		{"E", false},
		{"F", false},
		{"F", true},
		{"G", false},
		{"G", true},
		{"A", false},
		{"A", true},
		{"B", false},
	}

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

	var style = noTop(lipgloss.NewStyle().
		BorderStyle(border).
		BorderForeground(colorBorder).
		BorderBackground(colorWhite).
		Background(colorWhite).
		Foreground(colorBlack).
		Align())

	var styleFirst = noTop(style.Copy().Border(borderFirstWhite))

	styleSharp := noTop(style.Copy().
		BorderStyle(borderSharp).
		BorderForeground(colorBorder).
		BorderBackground(colorBlack).
		Background(colorBlack).
		Foreground(colorWhite).
		Align())

	var styleFirstSharp = noTop(styleSharp.Copy().Border(borderFirstSharp))
	var styleLastSharp = noTop(styleSharp.Copy().Border(borderSharp)).BorderRight(true).Align()

	omitBorder := lipgloss.HiddenBorder()
	omitBorder.Top = ""
	omitBorder.TopRight = ""
	omitBorder.TopLeft = ""
	omitBorderLeft, omitBorderRight := omitBorder, omitBorder
	omitBorderLeft.Left = style.GetBorderStyle().Left
	omitBorderLeft.BottomLeft = style.GetBorderStyle().Left
	omitBorderRight.Right = style.GetBorderStyle().Right
	omitBorderRight.BottomRight = style.GetBorderStyle().Right

	styleOmitLeft := noTop(style.Copy().Border(omitBorderLeft))
	styleOmitRight := noTop(style.Copy().Border(omitBorderRight))

	var white, black []string

	for i := 0; i < maxNote; i++ {
		style := &style
		if i == 0 {
			style = &styleFirst
		}
		noteInOct := i % len(notes)
		note := notes[noteInOct]

		lastNoteInOct := (i - 1)
		if lastNoteInOct < 0 {
			lastNoteInOct += len(notes)
		}
		lastNoteInOct %= len(notes)
		lastNote := notes[lastNoteInOct]

		lastLastNoteInOct := (i - 2)
		if lastLastNoteInOct < 0 {
			lastLastNoteInOct += len(notes)
		}
		lastLastNoteInOct %= len(notes)
		lastLastNote := notes[lastLastNoteInOct]

		nextNoteInOct := (i + 1) % len(notes)
		nextNote := notes[nextNoteInOct]

		nextNextNoteInOct := (i + 2) % len(notes)
		nextNextNote := notes[nextNextNoteInOct]

		s := note.name
		if note.sharp {
			s += "#"
			// s += "♯"
		} else {
			s += " "
		}

		const pad = ""

		if note.sharp {
			styleSharp := &styleSharp
			if !lastLastNote.sharp {
				styleSharp = &styleFirstSharp
			}
			if !nextNextNote.sharp {
				styleSharp = &styleLastSharp
			}

			black = append(black, styleSharp.Render(s))
			// white = append(white, pad)
		} else {
			if !lastNote.sharp {
				black = append(black, styleOmitLeft.Render(pad))
			}
			if !nextNote.sharp {
				black = append(black, styleOmitRight.Render(pad))
			}
			white = append(white, style.Render(s))
			// black = append(black, pad)
		}

	}

	strBlack := lipgloss.JoinHorizontal(lipgloss.Bottom, black...)
	// strBlack = ""
	strWhite := lipgloss.JoinHorizontal(lipgloss.Bottom, white...)
	str := lipgloss.JoinVertical(lipgloss.Left, strBlack, strWhite)
	s += str + "\n"

	// Send the UI for rendering
	return s
}
