package source

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/Zebbeni/ansizalizer/controls/browser"
	"github.com/Zebbeni/ansizalizer/event"
)

type Direction int

const (
	Left Direction = iota
	Right
	Up
	Down
)

var (
	navMap = map[Direction]map[State]State{
		Right: {ExpFile: ExpDirectory},
		Left:  {ExpDirectory: ExpFile},
		Down:  {ExpFile: Input, ExpDirectory: Input},
		Up:    {Input: ExpFile},
	}
	fileExts = []string{".png", ".jpg"}
)

func (m Model) handleEsc() (Model, tea.Cmd) {
	m.ShouldClose = true
	m.IsActive = false
	return m, nil
}

func (m Model) handleNav(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch {
	case key.Matches(msg, event.KeyMap.Right):
		if next, hasNext := navMap[Right][m.focus]; hasNext {
			m.focus = next
		}
	case key.Matches(msg, event.KeyMap.Left):
		if next, hasNext := navMap[Left][m.focus]; hasNext {
			m.focus = next
		}
	case key.Matches(msg, event.KeyMap.Down):
		if next, hasNext := navMap[Down][m.focus]; hasNext {
			m.focus = next
		}
	case key.Matches(msg, event.KeyMap.Up):
		if next, hasNext := navMap[Up][m.focus]; hasNext {
			m.focus = next
		} else {
			m.ShouldClose = true
		}
	}
	return m, nil
}

func (m Model) handleEnter() (Model, tea.Cmd) {
	switch m.focus {
	case ExpFile:
		m.focus = Browser
		m.doExportDirectory = false
		m.Browser = browser.New(fileExts, m.width)
	case ExpDirectory:
		m.focus = Browser
		m.doExportDirectory = true
		m.Browser = browser.New(nil, m.width)
	case SubDirYes:
		m.doExportSubDirectories = true
	case SubDirsNo:
		m.doExportSubDirectories = false
	case Input:
		m.focus = Browser
	}
	return m, nil
}

func (m Model) handleSrcBrowserUpdate(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.Browser, cmd = m.Browser.Update(msg)
	if m.doExportDirectory {
		m.selectedDir = m.Browser.SelectedDir
	} else {
		m.selectedFile = m.Browser.ActiveFile
	}

	if m.Browser.ShouldClose {
		m.focus = Input
		m.Browser.ShouldClose = false
	}
	return m, cmd
}
