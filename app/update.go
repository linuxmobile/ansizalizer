package app

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/Zebbeni/ansizalizer/app/adapt"
	"github.com/Zebbeni/ansizalizer/app/process"
	"github.com/Zebbeni/ansizalizer/event"
)

func (m Model) handleStartRenderMsg() (Model, tea.Cmd) {
	m.viewer.WaitingOnRender = true
	return m, m.processRenderCmd
}

func (m Model) handleFinishRenderMsg(msg event.FinishRenderMsg) (Model, tea.Cmd) {
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
	colorsString := "true color"
	if m.controls.Settings.Colors.IsLimited() {
		palette := m.controls.Settings.Colors.GetCurrentPalette()
		colorsString = palette.Title()
	}
	return event.FinishRenderMsg{FilePath: m.controls.FileBrowser.ActiveFile, ImgString: imgString, ColorsString: colorsString}
}

func (m Model) handleStartAdaptingMsg() (Model, tea.Cmd) {
	filename := m.controls.FileBrowser.ActiveFilename()
	message := fmt.Sprintf("generating palette from %s...", filename)
	return m, tea.Batch(event.BuildDisplayCmd(message), m.processAdaptingCmd)
}

func (m Model) handleFinishAdaptingMsg(msg event.FinishAdaptingMsg) (Model, tea.Cmd) {
	m.controls.Settings.Colors.Adapter = m.controls.Settings.Colors.Adapter.SetPalette(msg.Colors, msg.Name)
	return m, tea.Batch(event.StartRenderCmd, event.BuildDisplayCmd("rendering..."))
}

type Foo struct {
	Bar string
}

func (m Model) handleLospecRequestMsg(msg event.LospecRequestMsg) (Model, tea.Cmd) {
	// make url request
	r, err := http.Get(msg.URL)
	if err != nil {
		return m, event.BuildDisplayCmd("error making lospec request")
	}
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return m, event.BuildDisplayCmd("error reading lospec response")
	}

	// parse json and populate LospecResponseMsg
	data := new(event.LospecData)
	err = json.Unmarshal(body, &data)
	if err != nil {
		return m, event.BuildDisplayCmd("error decoding lospec request")
	}

	// build Data Cmd
	return m, event.BuildLospecResponseCmd(event.LospecResponseMsg{
		ID:   msg.ID,
		Page: msg.Page,
		Data: *data,
	})
}

func (m Model) handleLospecResponseMsg(msg event.LospecResponseMsg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.controls.Settings.Colors.Lospec, cmd = m.controls.Settings.Colors.Lospec.Update(msg)
	return m, cmd
}

func (m Model) processAdaptingCmd() tea.Msg {
	colors, name := adapt.GeneratePalette(m.controls.Settings.Colors.Adapter, m.controls.FileBrowser.ActiveFile)
	return event.FinishAdaptingMsg{
		Name:   name,
		Colors: colors,
	}
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
		return m, event.BuildDisplayCmd("Error copying to clipboard")
		// we should have a place in the UI where we display errors or processing messages,
		// and package our desired event to the user in a specific command
	}
	filename := m.controls.FileBrowser.ActiveFilename()
	name := strings.Split(filename, ".")[0] // strip extension
	return m, event.BuildDisplayCmd(fmt.Sprintf("copied %s to clipboard", name))
}

func (m Model) handleSave() (Model, tea.Cmd) {
	var filename string
	return m, event.BuildDisplayCmd(fmt.Sprintf("saved to %s", filename))
}
