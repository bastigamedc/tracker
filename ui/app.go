package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/bastigamedc/tracker/store"
)

type Screen int

const (
	ScreenMenu Screen = iota
	ScreenEntry
	ScreenToday
	ScreenReport
	ScreenCompanies
)

type entryForm struct {
	step           int
	companyIdx     int
	companyTyping  bool
	companyText    string
	companyInput   *input
	dateInput      *input
	durationInput  *input
	noteInput      *input
	successMessage string
	errorText      string
}

type Model struct {
	store  *store.Store
	width  int
	height int
	screen Screen

	menuIndex int

	entry *entryForm

	todayEntries  []store.Entry
	todaySelected int
	todayWeek     int
	todayYear     int

	reportEntries []store.Entry
	reportMonth   time.Month
	reportYear    int

	companies       []string
	companySel      int
	addingCompany   bool
	newCompanyInput *input
}

func New(store *store.Store) *Model {
	return &Model{
		store:  store,
		screen: ScreenMenu,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		return m.handleKey(msg)
	default:
		return m, nil
	}
}

func (m *Model) View() string {
	var content string
	switch m.screen {
	case ScreenMenu:
		content = m.menuView()
	case ScreenEntry:
		content = m.entryView()
	case ScreenToday:
		content = m.todayView()
	case ScreenReport:
		content = m.reportView()
	case ScreenCompanies:
		content = m.companiesView()
	}
	return baseStyle.Render(content)
}

func (m *Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.screen {
	case ScreenMenu:
		return m.handleMenuKey(msg)
	case ScreenEntry:
		return m.handleEntryKey(msg)
	case ScreenToday:
		return m.handleTodayKey(msg)
	case ScreenReport:
		return m.handleReportKey(msg)
	case ScreenCompanies:
		return m.handleCompaniesKey(msg)
	}
	return m, nil
}

func (m *Model) quitCmd() tea.Cmd {
	return tea.Quit
}

func (m *Model) mainMenuItems() []string {
	return []string{
		"Neue Zeit eintragen",
		"Heutige Zusammenfassung",
		"Monatsreport",
		"Firmen verwalten",
	}
}

func (m *Model) confirmMenuSelection() (tea.Model, tea.Cmd) {
	switch m.menuIndex {
	case 0:
		m.beginEntryForm()
	case 1:
		m.loadToday()
		m.screen = ScreenToday
	case 2:
		m.loadReport()
		m.screen = ScreenReport
	case 3:
		m.loadCompanies()
		m.screen = ScreenCompanies
	}
	return m, nil
}

func (m *Model) handleMenuKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, m.quitCmd()
	case "up", "k":
		m.menuIndex = (m.menuIndex + len(m.mainMenuItems()) - 1) % len(m.mainMenuItems())
	case "down", "j":
		m.menuIndex = (m.menuIndex + 1) % len(m.mainMenuItems())
	case "enter":
		return m.confirmMenuSelection()
	}
	return m, nil
}

func (m *Model) beginEntryForm() {
	companies, _ := m.store.LoadCompanies()
	now := time.Now()
	m.entry = &entryForm{
		step:          0,
		companyIdx:    0,
		companyTyping: len(companies) == 0,
		companyInput:  newInput("", "Firmenname", 24),
		dateInput:     newInput(now.Format("02.01.2006"), "02.01.2026 oder leer für heute", 22),
		durationInput: newInput("", "z.B. 2h30m", 16),
		noteInput:     newInput("", "optional", 44),
	}
	if len(companies) == 0 {
		m.entry.companyInput.Focus()
	}
	m.screen = ScreenEntry
}

func (m *Model) loadToday() {
	now := time.Now()
	year, week := store.ISOWeek(now)
	entries, _ := m.store.AllEntries(year, week)
	date := now.Format("2006-01-02")
	var today []store.Entry
	for _, e := range entries {
		if e.Date == date {
			today = append(today, e)
		}
	}
	m.todayEntries = today
	m.todaySelected = 0
	m.todayWeek = week
	m.todayYear = year
}

func (m *Model) loadReport() {
	now := time.Now()
	m.reportYear = now.Year()
	m.loadReportMonth(now.Month())
}

func (m *Model) loadReportMonth(month time.Month) {
	entries, _ := m.store.EntriesForYear(m.reportYear)
	var filtered []store.Entry
	for _, e := range entries {
		t, err := time.Parse("2006-01-02", e.Date)
		if err != nil {
			continue
		}
		if t.Month() == month && t.Year() == m.reportYear {
			filtered = append(filtered, e)
		}
	}
	m.reportEntries = filtered
	m.reportMonth = month
}

func (m *Model) loadCompanies() {
	companies, _ := m.store.LoadCompanies()
	m.companies = companies
	m.companySel = 0
	m.addingCompany = false
}

func monthName(m time.Month) string {
	names := [...]string{
		"Januar", "Februar", "März", "April", "Mai", "Juni",
		"Juli", "August", "September", "Oktober", "November", "Dezember",
	}
	return names[m-1]
}
