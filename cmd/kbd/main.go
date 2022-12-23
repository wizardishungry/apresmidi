package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	numMidiNotes = 128
	initOctave   = -2
)

type notePressMap [numMidiNotes]*bool

func main() {
	p := tea.NewProgram(initialModel())

	go sendRandomChords(p)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func sendRandomChords(p *tea.Program) {
	const (
		MIN_NOTES = 1
		MAX_NOTES = 8
		MIN_SLEEP = 250 * time.Millisecond
		MAX_SLEEP = 2000 * time.Millisecond
	)

	var noteState [numMidiNotes]bool

	lift := make(chan [numMidiNotes]*bool)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tick := time.NewTicker(100 * time.Millisecond)
	defer tick.Stop()
	var t = true
	var f = false
	go func() {
		for i := 0; i < 5; i++ {
			time.Sleep(250 * time.Millisecond * time.Duration(i))
			go func(dur time.Duration, top int) {
				if top > numMidiNotes {
					panic("lol no")
				}
				tick := time.NewTicker(dur)
				defer tick.Stop()
				direction := false
			WALK:
				for walk := 0; walk < top; walk++ {
					var press [128]*bool

					i := walk
					if direction {
						i = top - walk - 1
					}

					press[i] = &t
					p.Send((notePressMap)(press))
					<-tick.C
					press[i] = &f
					p.Send((notePressMap)(press))
				}
				direction = !direction
				goto WALK
			}(150*time.Millisecond+5*time.Millisecond*time.Duration(i+1), 32)
		}
	}()
	return
	for {

		var press, unpress [128]*bool
		numNotes := rng.Intn(MAX_NOTES-MIN_NOTES) + MIN_NOTES
		for i := 0; i < numNotes; i++ {
			tryNote := rng.Intn(numMidiNotes)
			if noteState[tryNote] {
				continue
			}
			if press[tryNote] != nil {
				continue
			}

			noteState[tryNote] = true
			press[tryNote] = &t
			unpress[tryNote] = &f
		}

		p.Send((notePressMap)(press))
		go func(unpress [numMidiNotes]*bool) {
			sleepTime := time.Duration(rng.Int63n(int64(MAX_SLEEP-MIN_SLEEP))) + MIN_SLEEP
			time.Sleep(sleepTime)
			lift <- unpress
		}(unpress)

	SLEEP:
		for {
			select {
			case unpress := <-lift:
				p.Send((notePressMap)(unpress))
				for i, v := range unpress {
					if v != nil {
						noteState[i] = false
					}
				}
			case <-tick.C:
				break SLEEP
			}
		}

	}
}

type model struct {
	choices      []string         // items on the to-do list
	cursor       int              // which to-do list item our cursor is pointing at
	selected     map[int]struct{} // which to-do items are selected
	notePressMap [128]bool
}

func initialModel() model {
	return model{
		// Our to-do list is a grocery list
		choices: []string{"Buy carrots", "Buy celery", "Buy kohlrabi"},

		// A map which indicates which choices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		selected:     make(map[int]struct{}),
		notePressMap: [128]bool{},
	}
}
func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) updateNPM(msg notePressMap) model {

	for i, v := range msg {
		if v != nil {
			m.notePressMap[i] = *v
		}
	}

	return m
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case notePressMap:
		m = m.updateNPM(msg)

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
	s = "Press q to quit.\n"

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

	var white, black []string

	for i := 0; i < numMidiNotes; i++ {
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
		s += ""

		if m.notePressMap[i] {
			s = styleReverse.Render(s)
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
