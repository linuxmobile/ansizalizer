package lospec

import (
	"fmt"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/lipgloss"

	"github.com/Zebbeni/ansizalizer/style"
)

var (
	stateNames = map[State]string{
		CountForm:        "Colors",
		TagForm:          "Tag",
		FilterAny:        "Any",
		FilterExact:      "Exact",
		FilterMax:        "Max",
		FilterMin:        "Min",
		SortAlphabetical: "A-Z",
		SortDownloads:    "Downloads",
		SortNewest:       "Newest",
	}

	filterOrder = []State{FilterAny, FilterExact, FilterMax, FilterMin}
	sortOrder   = []State{SortAlphabetical, SortDownloads, SortNewest}

	activeColor = lipgloss.Color("#aaaaaa")
	focusColor  = lipgloss.Color("#ffffff")
	normalColor = lipgloss.Color("#555555")
	titleStyle  = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888"))
)

func (m Model) drawInputs() string {
	prompt, placeholder := m.getInputColors(CountForm)

	m.countInput.PromptStyle = m.countInput.PromptStyle.Copy().Foreground(prompt)
	m.countInput.PlaceholderStyle = m.countInput.PlaceholderStyle.Copy().Foreground(placeholder)
	if m.countInput.Focused() == false {
		m.countInput.Placeholder = fmt.Sprintf("%4s", m.countInput.Value())
	} else {
		m.countInput.Placeholder = "    "
	}
	if m.countInput.Focused() {
		m.countInput.Cursor.SetMode(cursor.CursorBlink)
	} else {
		m.countInput.Cursor.SetMode(cursor.CursorHide)
	}

	prompt, placeholder = m.getInputColors(TagForm)
	m.tagInput.PromptStyle = m.countInput.PromptStyle.Copy().Foreground(prompt)
	m.tagInput.PlaceholderStyle = m.countInput.PlaceholderStyle.Copy().Foreground(placeholder)
	if m.tagInput.Focused() == false {
		m.tagInput.Placeholder = m.tagInput.Value()
	} else {
		m.tagInput.Placeholder = "    "
	}
	if m.tagInput.Focused() {
		m.tagInput.Cursor.SetMode(cursor.CursorBlink)
	} else {
		m.tagInput.Cursor.SetMode(cursor.CursorHide)
	}

	countForm := m.countInput.View()
	tagForm := m.tagInput.View()
	return lipgloss.JoinHorizontal(lipgloss.Left, countForm, tagForm)
}

func (m Model) drawFilterButtons() string {
	buttons := make([]string, len(filterOrder))
	for i, filter := range filterOrder {
		buttonStyle := style.NormalButton
		if filter == m.focus {
			if m.IsActive {
				buttonStyle = style.FocusButton
			} else {
				buttonStyle = style.ActiveButton
			}
		}
		buttons[i] = buttonStyle.Render(stateNames[filter])
	}

	return lipgloss.JoinHorizontal(lipgloss.Left, buttons...)
}

func (m Model) getInputColors(state State) (lipgloss.Color, lipgloss.Color) {
	if m.IsActive {
		if m.focus == state {
			return focusColor, focusColor
		} else if m.active == state {
			return activeColor, activeColor
		}
	}
	return normalColor, normalColor
}
