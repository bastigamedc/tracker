package ui

import (
	"fmt"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/bastigamedc/tracker/store"
)

func (m *Model) reportView() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("📈  Monatsreport"))
	b.WriteString(subtitleStyle.Render(fmt.Sprintf("%s %d",
		monthName(m.reportMonth), m.reportYear)))
	b.WriteString("\n")

	if len(m.reportEntries) == 0 {
		b.WriteString(dimStyle.Render("Keine Einträge in diesem Monat."))
		b.WriteString("\n")
	} else {
		type companyTotal struct {
			company string
			minutes int
			count   int
		}
		aggregate := map[string]*companyTotal{}
		var order []string
		for _, e := range m.reportEntries {
			if _, ok := aggregate[e.Company]; !ok {
				aggregate[e.Company] = &companyTotal{company: e.Company}
				order = append(order, e.Company)
			}
			aggregate[e.Company].minutes += store.MinutesFromDuration(e.Duration)
			aggregate[e.Company].count++
		}
		sort.Strings(order)

		header := fmt.Sprintf("%-14s %-10s %s",
			lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#8a8f98")).Render("FIRMA"),
			lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#8a8f98")).Render("DAUER"),
			lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#8a8f98")).Render("EINTRÄGE"),
		)
		b.WriteString(header)
		b.WriteString("\n")

		grandTotal := 0
		for _, name := range order {
			ct := aggregate[name]
			grandTotal += ct.minutes
			b.WriteString(fmt.Sprintf("%-14s %-10s %d",
				companyStyle.Render(name),
				accentStyle.Render(store.FormatDuration(ct.minutes)),
				ct.count,
			))
			b.WriteString("\n")
		}

		b.WriteString(dividerStyle.Render("─"))
		b.WriteString("\n")
		b.WriteString(fmt.Sprintf("%-14s %s",
			"Gesamt:",
			accentStyle.Render(store.FormatDuration(grandTotal)),
		))
		b.WriteString("\n\n")

		workdays := len(m.reportEntries)
		_ = workdays
		b.WriteString(dimStyle.Render(fmt.Sprintf("Tage mit Einträgen: %d", len(uniqueDates(m.reportEntries)))))
		b.WriteString("\n")
		b.WriteString(helpStyle.Render("[←/→] Monat wechseln   [esc] Menü"))
		return b.String()
	}

	b.WriteString(helpStyle.Render("[←/→] Monat wechseln   [esc] Menü"))
	return b.String()
}

func uniqueDates(entries []store.Entry) []string {
	seen := map[string]bool{}
	var out []string
	for _, e := range entries {
		if !seen[e.Date] {
			seen[e.Date] = true
			out = append(out, e.Date)
		}
	}
	return out
}

func (m *Model) handleReportKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.screen = ScreenMenu
	case "left", "h":
		prev := m.reportMonth - 1
		if prev < time.January {
			prev = time.December
			m.reportYear--
		}
		m.loadReportMonth(prev)
	case "right", "l":
		next := m.reportMonth + 1
		if next > time.December {
			next = time.January
			m.reportYear++
		}
		m.loadReportMonth(next)
	default:
		return m, nil
	}
	return m, nil
}
