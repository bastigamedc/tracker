package ui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/google/uuid"

	"github.com/bastigamedc/tracker/store"
)

const (
	entryStepCompany = iota
	entryStepDate
	entryStepDuration
	entryStepNote
	entryStepConfirm
	entryStepDone
)

func labelStyle(s string) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("#8a8f98")).Bold(true).MarginTop(1).Render(s) + "\n"
}

func (m *Model) entryView() string {
	f := m.entry
	if f == nil {
		return ""
	}

	var b strings.Builder
	b.WriteString(titleStyle.Render("📝  Neue Zeit eintragen"))
	b.WriteString("\n")

	switch f.step {
	case entryStepCompany:
		if f.companyTyping {
			b.WriteString(labelStyle("Firma (neu):"))
			b.WriteString(f.companyInput.View())
			if f.errorText != "" {
				b.WriteString(errorStyle.Render(f.errorText))
			}
			b.WriteString(helpStyle.Render("[Enter] Weiter   [esc] Zurück"))
		} else {
			b.WriteString(labelStyle("Firma auswählen:"))
			companies, _ := m.store.LoadCompanies()
			for i, c := range companies {
				if i == f.companyIdx {
					b.WriteString(selectedRowStyle.Render("▸ " + c))
				} else {
					b.WriteString(menuItemStyle.Render("  " + c))
				}
				b.WriteString("\n")
			}
			if f.companyIdx == len(companies) {
				b.WriteString(selectedRowStyle.Render("▸ ✚ Neue Firma eingeben"))
			} else {
				b.WriteString(menuItemStyle.Render("  ✚ Neue Firma eingeben"))
			}
			b.WriteString("\n")
			b.WriteString(helpStyle.Render("[↑/↓] Auswählen   [Enter] Bestätigen   [esc] Menü"))
		}

	case entryStepDate:
		b.WriteString(labelStyle("Datum:"))
		b.WriteString(f.dateInput.View())
		if f.errorText != "" {
			b.WriteString(errorStyle.Render(f.errorText))
		}
		b.WriteString(helpStyle.Render("[Enter] Weiter   [esc] Zurück"))

	case entryStepDuration:
		b.WriteString(labelStyle("Dauer:"))
		b.WriteString(f.durationInput.View())
		if f.errorText != "" {
			b.WriteString(errorStyle.Render(f.errorText))
		}
		b.WriteString(helpStyle.Render("[Enter] Weiter   [esc] Zurück"))

	case entryStepNote:
		b.WriteString(labelStyle("Notiz (optional):"))
		b.WriteString(f.noteInput.View())
		b.WriteString(helpStyle.Render("[Enter] Weiter   [esc] Zurück"))

	case entryStepConfirm:
		b.WriteString(accentStyle.Render("Eintrag speichern?"))
		b.WriteString("\n\n")
		b.WriteString(confirmRow("Firma", m.entryCompanyName()))
		b.WriteString("\n")
		b.WriteString(confirmRow("Datum", m.entryDateParsed().Format("02.01.2006")))
		b.WriteString("\n")
		b.WriteString(confirmRow("Dauer", m.entryDurationDisplay()))
		if n := strings.TrimSpace(f.noteInput.Value()); n != "" {
			b.WriteString("\n")
			b.WriteString(confirmRow("Notiz", n))
		}
		if f.errorText != "" {
			b.WriteString(errorStyle.Render(f.errorText))
		}
		b.WriteString("\n")
		b.WriteString(helpStyle.Render("[Enter] Speichern   [esc] Bearbeiten"))

	case entryStepDone:
		b.WriteString(successStyle.Render("✓ Zeit wurde gespeichert"))
		b.WriteString("\n\n")
		b.WriteString(f.successMessage)
		b.WriteString("\n")
		b.WriteString(helpStyle.Render("[Enter] Zum Menü"))
	}

	return b.String()
}

func confirmRow(key string, value string) string {
	return lipgloss.NewStyle().Width(14).Foreground(lipgloss.Color("#8a8f98")).Render(key) +
		accentStyle.Render(value)
}

func (m *Model) entryCompanyName() string {
	f := m.entry
	if f == nil {
		return ""
	}
	if f.companyTyping {
		return strings.TrimSpace(f.companyText)
	}
	companies, _ := m.store.LoadCompanies()
	if f.companyIdx >= 0 && f.companyIdx < len(companies) {
		return companies[f.companyIdx]
	}
	return ""
}

func (m *Model) entryDateParsed() time.Time {
	f := m.entry
	if f == nil {
		return time.Now()
	}
	raw := strings.TrimSpace(f.dateInput.Value())
	if raw == "" {
		raw = time.Now().Format("02.01.2006")
	}
	for _, layout := range []string{"02.01.2006", "2006-01-02", "02.01.06"} {
		if t, err := time.ParseInLocation(layout, raw, time.Local); err == nil {
			return t
		}
	}
	return time.Now()
}

func (m *Model) entryDurationMinutes() int {
	f := m.entry
	if f == nil {
		return 0
	}
	minutes, err := store.ParseDuration(f.durationInput.Value())
	if err != nil {
		return 0
	}
	return minutes
}

func (m *Model) entryDurationDisplay() string {
	return store.FormatDuration(m.entryDurationMinutes())
}

func (m *Model) handleEntryKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	f := m.entry
	if f == nil {
		m.screen = ScreenMenu
		return m, nil
	}

	switch f.step {
	case entryStepCompany:
		return m.handleEntryCompanyStep(msg)
	case entryStepDate:
		return m.handleEntryDateStep(msg)
	case entryStepDuration:
		return m.handleEntryDurationStep(msg)
	case entryStepNote:
		return m.handleEntryNoteStep(msg)
	case entryStepConfirm:
		return m.handleEntryConfirmStep(msg)
	case entryStepDone:
		if msg.String() == "enter" || msg.String() == "esc" {
			m.entry = nil
			m.screen = ScreenMenu
		}
		return m, nil
	}
	return m, nil
}

func (m *Model) handleEntryCompanyStep(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	f := m.entry

	if f.companyTyping {
		switch msg.String() {
		case "enter":
			name := strings.TrimSpace(f.companyInput.Value())
			if name == "" {
				f.errorText = "Firmenname darf nicht leer sein"
				return m, nil
			}
			f.errorText = ""
			f.companyText = name
			f.step = entryStepDate
			f.dateInput.Focus()
			return m, nil
		case "esc":
			f.errorText = ""
			companies, _ := m.store.LoadCompanies()
			if len(companies) > 0 {
				f.companyTyping = false
			} else {
				m.entry = nil
				m.screen = ScreenMenu
			}
			return m, nil
		default:
			cmd := f.companyInput.Update(msg)
			return m, cmd
		}
	}

	companies, _ := m.store.LoadCompanies()
	count := len(companies) + 1
	switch msg.String() {
	case "up", "k":
		f.companyIdx = (f.companyIdx + count - 1) % count
	case "down", "j":
		f.companyIdx = (f.companyIdx + 1) % count
	case "enter":
		if f.companyIdx == len(companies) {
			f.companyTyping = true
			f.companyInput.Focus()
		} else {
			f.step = entryStepDate
			f.dateInput.Focus()
		}
	case "esc":
		m.entry = nil
		m.screen = ScreenMenu
	default:
		return m, nil
	}
	return m, nil
}

func (m *Model) handleEntryDateStep(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	f := m.entry
	switch msg.String() {
	case "enter":
		raw := strings.TrimSpace(f.dateInput.Value())
		valid := raw == ""
		for _, layout := range []string{"02.01.2006", "2006-01-02", "02.01.06"} {
			if _, err := time.ParseInLocation(layout, raw, time.Local); err == nil {
				valid = true
				break
			}
		}
		if !valid {
			f.errorText = "Ungültiges Datum - erwartet z.B. 02.01.2026"
			return m, nil
		}
		f.errorText = ""
		f.step = entryStepDuration
		f.durationInput.Focus()
		return m, nil
	case "esc":
		f.errorText = ""
		f.step = entryStepCompany
		return m, nil
	default:
		cmd := f.dateInput.Update(msg)
		return m, cmd
	}
}

func (m *Model) handleEntryDurationStep(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	f := m.entry
	switch msg.String() {
	case "enter":
		if _, err := store.ParseDuration(f.durationInput.Value()); err != nil {
			f.errorText = err.Error()
			return m, nil
		}
		f.errorText = ""
		f.step = entryStepNote
		f.noteInput.Focus()
		return m, nil
	case "esc":
		f.errorText = ""
		f.step = entryStepDate
		return m, nil
	default:
		cmd := f.durationInput.Update(msg)
		return m, cmd
	}
}

func (m *Model) handleEntryNoteStep(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	f := m.entry
	switch msg.String() {
	case "enter":
		f.step = entryStepConfirm
		return m, nil
	case "esc":
		f.step = entryStepDuration
		return m, nil
	default:
		cmd := f.noteInput.Update(msg)
		return m, cmd
	}
}

func (m *Model) handleEntryConfirmStep(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	f := m.entry
	switch msg.String() {
	case "enter":
		company := m.entryCompanyName()
		if company == "" {
			f.errorText = "Keine gültige Firma"
			return m, nil
		}
		date := m.entryDateParsed()
		duration := m.entryDurationDisplay()
		note := strings.TrimSpace(f.noteInput.Value())

		if err := m.store.EnsureCompany(company); err != nil {
			f.errorText = "Fehler beim Speichern: " + err.Error()
			return m, nil
		}
		entry := store.Entry{
			ID:       uuid.NewString(),
			Company:  company,
			Date:     date.Format("2006-01-02"),
			Duration: duration,
			Note:     note,
		}
		if err := m.store.AddEntry(entry); err != nil {
			f.errorText = "Fehler beim Speichern: " + err.Error()
			return m, nil
		}
		f.successMessage = fmt.Sprintf("%s   %s   %s",
			accentStyle.Render(company),
			accentStyle.Render(duration),
			dimStyle.Render(note),
		)
		f.step = entryStepDone
		return m, nil
	case "esc":
		f.errorText = ""
		f.step = entryStepNote
		return m, nil
	default:
		return m, nil
	}
}
