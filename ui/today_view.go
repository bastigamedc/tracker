package ui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/bastigamedc/tracker/store"
)

func (m *Model) todayView() string {
	var b strings.Builder

	now := time.Now()
	weekday := deutscherWochentag(now)
	date := now.Format("02.01.2006")
	b.WriteString(titleStyle.Render("📊  Heutige Zusammenfassung"))
	b.WriteString(subtitleStyle.Render(fmt.Sprintf("%s, %s  ·  KW %02d/%d",
		weekday, date, m.todayWeek, m.todayYear)))
	b.WriteString("\n")

	if len(m.todayEntries) == 0 {
		b.WriteString(dimStyle.Render("Heute wurde noch keine Zeit erfasst."))
		b.WriteString("\n")
	} else {
		header := fmt.Sprintf("%-12s %-10s %s",
			lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#8a8f98")).Render("FIRMA"),
			lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#8a8f98")).Render("DAUER"),
			lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#8a8f98")).Render("NOTIZ"),
		)
		b.WriteString(header)
		b.WriteString("\n")

		total := 0
		for i, e := range m.todayEntries {
			duration := store.MinutesFromDuration(e.Duration)
			total += duration
			line := fmt.Sprintf("%-12s %-10s %s",
				companyStyle.Render(e.Company),
				accentStyle.Render(e.Duration),
				dimStyle.Render(e.Note),
			)
			if i == m.todaySelected {
				b.WriteString(selectedMenuItemStyle.Render("▸ " + line))
			} else {
				b.WriteString(menuItemStyle.Render("  " + line))
			}
			b.WriteString("\n")
		}

		b.WriteString(dividerStyle.Render("─"))
		b.WriteString("\n")
		b.WriteString(fmt.Sprintf("%-12s %-10s",
			"Gesamt:",
			accentStyle.Render(store.FormatDuration(total)),
		))
		b.WriteString("\n")
	}

	if len(m.todayEntries) > 0 {
		b.WriteString(helpStyle.Render("[↑/↓] Auswählen   [d] Löschen   [esc] Menü"))
	} else {
		b.WriteString(helpStyle.Render("[esc] Menü"))
	}
	return b.String()
}

func deutscherWochentag(t time.Time) string {
	days := [...]string{"Sonntag", "Montag", "Dienstag", "Mittwoch", "Donnerstag", "Freitag", "Samstag"}
	return days[int(t.Weekday())]
}

func (m *Model) handleTodayKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.screen = ScreenMenu
	case "up", "k":
		if len(m.todayEntries) > 0 {
			m.todaySelected = (m.todaySelected + len(m.todayEntries) - 1) % len(m.todayEntries)
		}
	case "down", "j":
		if len(m.todayEntries) > 0 {
			m.todaySelected = (m.todaySelected + 1) % len(m.todayEntries)
		}
	case "d", "delete":
		if len(m.todayEntries) == 0 {
			return m, nil
		}
		id := m.todayEntries[m.todaySelected].ID
		_ = m.store.DeleteEntry(id)
		m.loadToday()
	}
	return m, nil
}
