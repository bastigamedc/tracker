package ui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/bastigamedc/tracker/store"
)

func newTestModel(t *testing.T) (*Model, *store.Store) {
	t.Helper()
	dir := t.TempDir()
	home := filepath.Join(dir, "home")
	if err := os.MkdirAll(home, 0o755); err != nil {
		t.Fatal(err)
	}
	st := store.NewStoreWithDir(home)
	m := New(st)
	return m, st
}

func updateWithKey(m *Model, s string) {
	msg := tea.KeyMsg{}
	switch s {
	case "enter":
		msg.Type = tea.KeyEnter
	case "esc":
		msg.Type = tea.KeyEsc
	case "up":
		msg.Type = tea.KeyUp
	case "down":
		msg.Type = tea.KeyDown
	case "left":
		msg.Type = tea.KeyLeft
	case "right":
		msg.Type = tea.KeyRight
	case "delete":
		msg.Type = tea.KeyDelete
	default:
		msg.Type = tea.KeyRunes
		msg.Runes = []rune(s)
	}
	m.handleKey(msg)
}

func TestMenuNavigation(t *testing.T) {
	m, _ := newTestModel(t)
	updateWithKey(m, "down")
	if m.menuIndex != 1 {
		t.Fatalf("menuIndex = %d, want 1", m.menuIndex)
	}
	view := m.menuView()
	if !strings.Contains(view, "Heutige Zusammenfassung") {
		t.Error("menu missing Today entry")
	}
	updateWithKey(m, "enter")
	if m.screen != ScreenToday {
		t.Fatalf("screen = %d, want ScreenToday", m.screen)
	}
}

func TestFullEntryFlow(t *testing.T) {
	m, st := newTestModel(t)

	m.beginEntryForm()
	if m.entry == nil {
		t.Fatal("entry form not created")
	}
	if !m.entry.companyTyping {
		t.Fatal("expected typing mode with no companies")
	}

	for _, r := range "acme" {
		updateWithKey(m, string(r))
	}
	if got := m.entry.companyInput.Value(); got != "acme" {
		t.Fatalf("typed company = %q, want acme", got)
	}
	updateWithKey(m, "enter")
	if m.entry.step != entryStepDate {
		t.Fatalf("step = %d, want entryStepDate", m.entry.step)
	}

	updateWithKey(m, "enter")
	if m.entry.step != entryStepDuration {
		t.Fatalf("step = %d, want duration", m.entry.step)
	}

	for _, r := range "2h30m" {
		updateWithKey(m, string(r))
	}
	updateWithKey(m, "enter")
	if m.entry.step != entryStepNote {
		t.Fatalf("step = %d, want note", m.entry.step)
	}

	for _, r := range "API" {
		updateWithKey(m, string(r))
	}
	updateWithKey(m, "enter")
	if m.entry.step != entryStepConfirm {
		t.Fatalf("step = %d, want confirm", m.entry.step)
	}

	updateWithKey(m, "enter")
	if m.entry.step != entryStepDone {
		t.Fatalf("step = %d, want done", m.entry.step)
	}

	companies, _ := st.LoadCompanies()
	if len(companies) != 1 || companies[0] != "acme" {
		t.Fatalf("companies = %v", companies)
	}
	now := time.Now()
	year, week := store.ISOWeek(now)
	wf, _ := st.LoadWeek(year, week)
	if len(wf.Entries) != 1 {
		t.Fatalf("entries = %d, want 1", len(wf.Entries))
	}
	e := wf.Entries[0]
	if e.Company != "acme" || e.Duration != "2h30m" || e.Note != "API" {
		t.Fatalf("entry = %+v", e)
	}
}

func TestTodayAfterEntry(t *testing.T) {
	m, st := newTestModel(t)
	now := time.Now()
	year, week := store.ISOWeek(now)
	st.AddEntry(store.Entry{ID: "1", Company: "acme", Date: now.Format("2006-01-02"), Duration: "1h", Note: "x"})
	st.AddEntry(store.Entry{ID: "2", Company: "acme", Date: "2020-01-01", Duration: "2h", Note: "old"})
	_ = year
	_ = week

	m.loadToday()
	m.screen = ScreenToday
	if len(m.todayEntries) != 1 {
		t.Fatalf("today entries = %d, want 1", len(m.todayEntries))
	}
	if m.todayEntries[0].Note != "x" {
		t.Fatalf("note = %q", m.todayEntries[0].Note)
	}

	m.todaySelected = 0
	updateWithKey(m, "d")
	year, week = store.ISOWeek(now)
	wf, _ := st.LoadWeek(year, week)
	if len(wf.Entries) != 0 {
		t.Fatalf("entries after delete = %d, want 0", len(wf.Entries))
	}
}

func TestReport(t *testing.T) {
	m, st := newTestModel(t)
	now := time.Now()
	st.AddEntry(store.Entry{ID: "1", Company: "acme", Date: now.Format("2006-01-02"), Duration: "1h30m", Note: ""})
	st.AddEntry(store.Entry{ID: "2", Company: "acme", Date: now.Format("2006-01-02"), Duration: "30m", Note: ""})
	st.AddEntry(store.Entry{ID: "3", Company: "globex", Date: "2020-01-01", Duration: "5h", Note: ""})

	m.reportYear = now.Year()
	m.loadReportMonth(now.Month())
	if len(m.reportEntries) != 2 {
		t.Fatalf("report entries = %d, want 2", len(m.reportEntries))
	}

	view := m.reportView()
	if !strings.Contains(view, "2h") {
		t.Error("report missing aggregated 2h total for acme")
	}
	if !strings.Contains(view, "acme") {
		t.Error("report missing acme")
	}
	if strings.Contains(view, "globex") {
		t.Error("report should not contain globex from 2020")
	}
	if !strings.Contains(view, "Gesamt:") {
		t.Error("report missing total")
	}
}

func TestCompaniesFlow(t *testing.T) {
	m, st := newTestModel(t)
	m.loadCompanies()
	m.screen = ScreenCompanies
	updateWithKey(m, "a")
	if !m.addingCompany {
		t.Fatal("expected adding mode")
	}
	for _, r := range "initech" {
		updateWithKey(m, string(r))
	}
	updateWithKey(m, "enter")
	companies, _ := st.LoadCompanies()
	if len(companies) != 1 || companies[0] != "initech" {
		t.Fatalf("companies = %v", companies)
	}

	st.EnsureCompany("acme")
	m.loadCompanies()
	m.companySel = 0
	updateWithKey(m, "x")
	companies, _ = st.LoadCompanies()
	if len(companies) != 1 || companies[0] != "acme" {
		t.Fatalf("after remove = %v", companies)
	}
}
