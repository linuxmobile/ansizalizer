package app

import (
	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/Zebbeni/ansizalizer/app/adapt"
	"github.com/Zebbeni/ansizalizer/app/process"
	"github.com/Zebbeni/ansizalizer/io"
)

func (m Model) handleStartRenderMsg() (Model, tea.Cmd) {
	m.viewer.WaitingOnRender = true
	return m, m.processRenderCmd
}

func (m Model) handleFinishRenderMsg(msg io.FinishRenderMsg) (Model, tea.Cmd) {
	// cut out early if the finished render is for a previously selected image
	if msg.FilePath != m.controls.FileBrowser.ActiveFile {
		return m, nil
	}

	var cmd tea.Cmd
	m.viewer, cmd = m.viewer.Update(msg)
	return m, cmd
}

func (m Model) processRenderCmd() tea.Msg {
	imgString := process.RenderImageFile(m.controls.Settings, m.controls.FileBrowser.ActiveFile)
	return io.FinishRenderMsg{FilePath: m.controls.FileBrowser.ActiveFile, ImgString: imgString}
}

func (m Model) handleStartAdaptingMsg() (Model, tea.Cmd) {
	return m, m.processAdaptingCmd
}

func (m Model) handleFinishAdaptingMsg(msg io.FinishAdaptingMsg) (Model, tea.Cmd) {
	m.controls.Settings.Colors.Adaptive.Palette = msg.Palette
	return m, tea.Batch(io.StartRenderCmd, io.BuildDisplayCmd("Rendering..."))
}

func (m Model) processAdaptingCmd() tea.Msg {
	palette := adapt.GeneratePalette(m.controls.Settings.Colors.Adaptive, m.controls.FileBrowser.ActiveFile)
	return io.FinishAdaptingMsg{Palette: palette}
}

func (m Model) handleControlsUpdate(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.controls, cmd = m.controls.Update(msg)
	return m, cmd
}

func (m Model) handleDisplayMsg(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.display, cmd = m.display.Update(msg)
	return m, cmd
}

func (m Model) handleCopy() (Model, tea.Cmd) {
	if err := clipboard.WriteAll(m.viewer.View()); err != nil {
		return m, io.BuildDisplayCmd("Error copying to clipboard")
		// we should have a place in the UI where we display errors or processing messages,
		// and package our desired message to the user in a specific command
	}
	return m, io.BuildDisplayCmd("Copied text to clipboard")
}
