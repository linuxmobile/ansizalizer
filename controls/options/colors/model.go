package colors

import (
	"image/color"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/makeworld-the-better-one/dither/v2"

	"github.com/Zebbeni/ansizalizer/controls/options/colors/adaptive"
	"github.com/Zebbeni/ansizalizer/controls/options/colors/limited"
)

type State int

// None consists of a few different components that are shown or hidden
// depending on which toggles have been set on / off. The Model state indicates
// which component is currently focused. From top to bottom the components are:

// 1) Limited (on/off)
// 2) Palette (Name) (if Limited) -> [Enter] displays Palette menu
// 3) Dithering (on/off) (if Limited)
// 4) Serpentine (on/off) (if Dithering)
// 5) Matrix (Name) (if Dithering) -> [Enter] displays to Matrix menu

// These can all be part of a single list, but we need to onSelect the list items

const (
	TrueColor State = iota
	Adaptive
	Palette
	PalettedControls
	AdaptiveControls
)

type Model struct {
	selected State
	focus    State // the component taking input
	controls State

	Adaptive adaptive.Model
	Palette  limited.Model

	ShouldClose      bool
	ShouldDeactivate bool

	IsActive bool
}

func New() Model {
	m := Model{
		selected:         TrueColor,
		focus:            TrueColor,
		controls:         TrueColor,
		Adaptive:         adaptive.New(),
		Palette:          limited.New(),
		ShouldClose:      false,
		ShouldDeactivate: false,
		IsActive:         false,
	}
	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch m.focus {
	case AdaptiveControls:
		return m.handleAdaptiveUpdate(msg)
	case PalettedControls:
		return m.handlePaletteUpdate(msg)
	}
	return m.handleMenuUpdate(msg)
}

func (m Model) View() string {
	buttons := m.drawButtons()
	if m.IsActive == false {
		return buttons
	}

	var controls string
	switch m.controls {
	case Adaptive:
		controls = m.Adaptive.View()
	case Palette:
		controls = m.Palette.View()
	}
	if len(controls) == 0 {
		return buttons
	}

	return lipgloss.JoinVertical(lipgloss.Top, buttons, controls)
}

func (m Model) IsLimited() bool {
	return m.selected != TrueColor
}

func (m Model) IsDithered() bool {
	return false
}

func (m Model) IsSerpentine() bool {
	return true
}

func (m Model) Matrix() dither.ErrorDiffusionMatrix {
	return dither.FloydSteinberg
}

func (m Model) IsAdaptive() bool {
	return m.selected == Adaptive
}

func (m Model) IsPaletted() bool {
	return m.selected == Palette
}

func (m Model) GetCurrentPalette() color.Palette {
	if m.selected == Palette {
		return m.Palette.GetCurrent()
	}
	return m.Adaptive.Palette
}
