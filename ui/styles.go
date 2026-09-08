package ui

import "github.com/charmbracelet/lipgloss"

var (
	baseStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7a5af8"))
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7a5af8")).
			Bold(true).
			MarginBottom(1)
	subtitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8a8f98")).
			MarginBottom(1)
	bigTitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7a5af8")).
			Bold(true).
			MarginBottom(1)
	menuItemStyle = lipgloss.NewStyle().
			PaddingLeft(1)
	selectedMenuItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#7a5af8")).
				Bold(true).
				PaddingLeft(1).
				Border(lipgloss.ThickBorder()).
				BorderForeground(lipgloss.Color("#7a5af8"))
	selectedRowStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#7a5af8")).
				Bold(true).
				PaddingLeft(1).
				Border(lipgloss.ThickBorder()).
				BorderTop(false).
				BorderRight(false).
				BorderBottom(false).
				BorderLeft(true).
				BorderLeftForeground(lipgloss.Color("#7a5af8"))
	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8a8f98")).
			MarginTop(1)
	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f84343")).
			Bold(true).
			MarginTop(1)
	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#3db562")).
			Bold(true).
			MarginTop(1)
	companyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6fb6ff")).
			Bold(true)
	accentStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7a5af8")).
			Bold(true)
	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8a8f98"))
)

var (
	weekStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8a8f98")).
			Bold(true)
	dividerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#3d4048")).
			MarginTop(1)
)
