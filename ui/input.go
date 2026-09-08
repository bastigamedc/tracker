package ui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type input struct {
	model textinput.Model
}

func newInput(value string, placeholder string, width int) *input {
	ti := textinput.New()
	ti.SetValue(value)
	ti.Placeholder = placeholder
	ti.Prompt = "> "
	ti.PromptStyle = accentStyle
	ti.Width = width
	return &input{model: ti}
}

func (i *input) Focus() {
	i.model.Focus()
}

func (i *input) Blur() {
	i.model.Blur()
}

func (i *input) Focused() bool {
	return i.model.Focused()
}

func (i *input) Value() string {
	return i.model.Value()
}

func (i *input) SetValue(v string) {
	i.model.SetValue(v)
}

func (i *input) Update(msg tea.Msg) tea.Cmd {
	updated, cmd := i.model.Update(msg)
	i.model = updated
	return cmd
}

func (i *input) View() string {
	return lipgloss.NewStyle().MarginBottom(1).Render(i.model.View())
}
