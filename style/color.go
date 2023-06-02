package style

import "github.com/charmbracelet/lipgloss"

var (
	NormalColor1   = lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#aaaaaa"}
	NormalColor2   = lipgloss.AdaptiveColor{Light: "#3a3a3a", Dark: "#888888"}
	SelectedColor1 = lipgloss.AdaptiveColor{Light: "#444444", Dark: "#ffffff"}
	SelectedColor2 = lipgloss.AdaptiveColor{Light: "#666666", Dark: "#dddddd"}
	DimmedColor1   = lipgloss.AdaptiveColor{Light: "#999999", Dark: "#777777"}
	DimmedColor2   = lipgloss.AdaptiveColor{Light: "#aaaaaa", Dark: "#666666"}

	NormalTitle     = lipgloss.NewStyle().Foreground(NormalColor1)
	NormalParagraph = lipgloss.NewStyle().Foreground(NormalColor2)

	SelectedTitle     = lipgloss.NewStyle().Foreground(SelectedColor1)
	SelectedParagraph = lipgloss.NewStyle().Foreground(SelectedColor2)

	DimmedTitle     = lipgloss.NewStyle().Foreground(DimmedColor1)
	DimmedParagraph = lipgloss.NewStyle().Foreground(DimmedColor2)
)
