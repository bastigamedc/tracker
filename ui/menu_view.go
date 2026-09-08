package ui

import (
	"strings"
)

func (m *Model) menuView() string {
	var b strings.Builder
	b.WriteString(bigTitleStyle.Render("⏱  WorkTrack"))
	b.WriteString(subtitleStyle.Render("Zeiterfassung für mehrere Firmen"))
	b.WriteString("\n\n")

	items := m.mainMenuItems()
	for i, item := range items {
		if i == m.menuIndex {
			b.WriteString(selectedMenuItemStyle.Render("▸ " + item))
		} else {
			b.WriteString(menuItemStyle.Render("  " + item))
		}
		b.WriteString("\n")
	}

	b.WriteString(dividerStyle.Render("─"))
	b.WriteString(helpStyle.Render("[↑/↓] Auswählen   [Enter] Öffnen   [q] Beenden"))
	return b.String()
}
