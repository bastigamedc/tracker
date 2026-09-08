package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) companiesView() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("🏢  Firmen verwalten"))
	b.WriteString("\n")

	if m.addingCompany {
		b.WriteString(labelStyle("Neue Firma:"))
		b.WriteString(m.newCompanyInput.View())
		b.WriteString(helpStyle.Render("[Enter] Speichern   [esc] Abbrechen"))
		return b.String()
	}

	if len(m.companies) == 0 {
		b.WriteString(dimStyle.Render("Noch keine Firmen vorhanden. Füge eine hinzu mit [a]."))
		b.WriteString("\n")
	} else {
		b.WriteString(labelStyle("Vorhandene Firmen:"))
		for i, c := range m.companies {
			if i == m.companySel {
				b.WriteString(selectedMenuItemStyle.Render("▸ " + c))
			} else {
				b.WriteString(menuItemStyle.Render("  " + c))
			}
			b.WriteString("\n")
		}
	}

	b.WriteString(helpStyle.Render("[↑/↓] Auswählen   [a] Hinzufügen   [x] Entfernen   [esc] Menü"))
	return b.String()
}

func (m *Model) handleCompaniesKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.addingCompany {
		switch msg.String() {
		case "enter":
			name := strings.TrimSpace(m.newCompanyInput.Value())
			if name == "" {
				return m, nil
			}
			if err := m.store.EnsureCompany(name); err == nil {
				m.loadCompanies()
			} else {
				m.addingCompany = false
			}
		case "esc":
			m.addingCompany = false
		default:
			cmd := m.newCompanyInput.Update(msg)
			return m, cmd
		}
		return m, nil
	}

	switch msg.String() {
	case "esc":
		m.screen = ScreenMenu
	case "up", "k":
		if len(m.companies) > 0 {
			m.companySel = (m.companySel + len(m.companies) - 1) % len(m.companies)
		}
	case "down", "j":
		if len(m.companies) > 0 {
			m.companySel = (m.companySel + 1) % len(m.companies)
		}
	case "a":
		m.addingCompany = true
		m.newCompanyInput = newInput("", "Firmenname", 24)
		m.newCompanyInput.Focus()
	case "x", "delete":
		if len(m.companies) == 0 || m.companySel < 0 || m.companySel >= len(m.companies) {
			return m, nil
		}
		name := m.companies[m.companySel]
		_ = m.store.RemoveCompany(name)
		m.loadCompanies()
	}
	return m, nil
}
