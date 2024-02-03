package cmd

import "github.com/charmbracelet/lipgloss"

var (
	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	// special    = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(highlight)
	docStyle   = lipgloss.NewStyle().Padding(1, 2, 1, 2)
	helpStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render
	padding    = 2
	maxWidth   = 80

	/*
		infoStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.NormalBorder()).
				BorderTop(true).
				BorderForeground(subtle)
	*/

	listHeader = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(subtle).
			MarginRight(2).
			Render

	ratingRow  = lipgloss.NewStyle().Width(50)
	ratingKey  = lipgloss.NewStyle().AlignHorizontal(lipgloss.Right).Width(20).PaddingRight(2)
	ratingItem = lipgloss.NewStyle().AlignHorizontal(lipgloss.Right).Width(20)

	listItem      = lipgloss.NewStyle().PaddingLeft(2).Render
	listItemMajor = lipgloss.NewStyle().PaddingLeft(2).Bold(true).Render
)
