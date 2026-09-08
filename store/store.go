package store

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	dirName       = ".worktrack"
	companiesFile = "companies.json"
	dataDirName   = "data"
)

type Entry struct {
	ID       string `json:"id"`
	Company  string `json:"company"`
	Date     string `json:"date"`
	Duration string `json:"duration"`
	Note     string `json:"note"`
}

type WeekFile struct {
	Week    int     `json:"week"`
	Year    int     `json:"year"`
	Entries []Entry `json:"entries"`
}

type CompaniesFile struct {
	Companies []string `json:"companies"`
}

type Store struct {
	dir string
}

func NewStore() (*Store, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("home dir: %w", err)
	}
	return &Store{dir: filepath.Join(home, dirName)}, nil
}

func NewStoreWithDir(dir string) *Store {
	return &Store{dir: dir}
}

func (s *Store) RootDir() string {
	return s.dir
}

func (s *Store) dataRoot() string {
	return filepath.Join(s.dir, dataDirName)
}

func (s *Store) weekPath(year int, week int) string {
	return filepath.Join(s.dataRoot(), fmt.Sprintf("%d", year), fmt.Sprintf("week-%d.json", week))
}

func ISOWeek(t time.Time) (int, int) {
	return t.ISOWeek()
}

func (s *Store) LoadCompanies() ([]string, error) {
	path := filepath.Join(s.dir, companiesFile)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}
	var cf CompaniesFile
	if err := decode(data, &cf); err != nil {
		return nil, err
	}
	return cf.Companies, nil
}

func (s *Store) SaveCompanies(companies []string) error {
	if err := os.MkdirAll(s.dir, 0o755); err != nil {
		return err
	}
	path := filepath.Join(s.dir, companiesFile)
	cf := CompaniesFile{Companies: companies}
	data, err := encode(cf)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func (s *Store) EnsureCompany(name string) error {
	companies, err := s.LoadCompanies()
	if err != nil {
		return err
	}
	for _, c := range companies {
		if c == name {
			return nil
		}
	}
	companies = append(companies, name)
	return s.SaveCompanies(companies)
}

func (s *Store) RemoveCompany(name string) error {
	companies, err := s.LoadCompanies()
	if err != nil {
		return err
	}
	filtered := companies[:0]
	for _, c := range companies {
		if c != name {
			filtered = append(filtered, c)
		}
	}
	return s.SaveCompanies(filtered)
}

func (s *Store) LoadWeek(year int, week int) (*WeekFile, error) {
	path := s.weekPath(year, week)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &WeekFile{Week: week, Year: year}, nil
		}
		return nil, err
	}
	var wf WeekFile
	if err := decode(data, &wf); err != nil {
		return nil, err
	}
	return &wf, nil
}

func (s *Store) SaveWeek(wf *WeekFile) error {
	dir := filepath.Dir(s.weekPath(wf.Year, wf.Week))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	path := s.weekPath(wf.Year, wf.Week)
	data, err := encode(wf)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func (s *Store) AddEntry(entry Entry) error {
	t, err := time.Parse("2006-01-02", entry.Date)
	if err != nil {
		return fmt.Errorf("date: %w", err)
	}
	year, week := ISOWeek(t)
	wf, err := s.LoadWeek(year, week)
	if err != nil {
		return err
	}
	wf.Entries = append(wf.Entries, entry)
	return s.SaveWeek(wf)
}

func (s *Store) DeleteEntry(id string) error {
	years, err := s.listYears()
	if err != nil {
		return err
	}
	for _, year := range years {
		weeks, err := s.listWeeks(year)
		if err != nil {
			return err
		}
		for _, week := range weeks {
			wf, err := s.LoadWeek(year, week)
			if err != nil {
				return err
			}
			filtered := wf.Entries[:0]
			found := false
			for _, e := range wf.Entries {
				if e.ID == id {
					found = true
					continue
				}
				filtered = append(filtered, e)
			}
			if found {
				wf.Entries = filtered
				return s.SaveWeek(wf)
			}
		}
	}
	return nil
}

func (s *Store) listYears() ([]int, error) {
	entries, err := os.ReadDir(s.dataRoot())
	if err != nil {
		if os.IsNotExist(err) {
			return []int{}, nil
		}
		return nil, err
	}
	var years []int
	for _, e := range entries {
		if e.IsDir() {
			var y int
			if _, err := fmt.Sscanf(e.Name(), "%d", &y); err == nil {
				years = append(years, y)
			}
		}
	}
	return years, nil
}

func (s *Store) listWeeks(year int) ([]int, error) {
	dir := filepath.Join(s.dataRoot(), fmt.Sprintf("%d", year))
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return []int{}, nil
		}
		return nil, err
	}
	var weeks []int
	for _, e := range entries {
		var w int
		if _, err := fmt.Sscanf(e.Name(), "week-%d.json", &w); err == nil {
			weeks = append(weeks, w)
		}
	}
	return weeks, nil
}

func (s *Store) AllEntries(year int, week int) ([]Entry, error) {
	wf, err := s.LoadWeek(year, week)
	if err != nil {
		return nil, err
	}
	return wf.Entries, nil
}

func (s *Store) EntriesForYear(year int) ([]Entry, error) {
	weeks, err := s.listWeeks(year)
	if err != nil {
		return nil, err
	}
	var all []Entry
	for _, w := range weeks {
		wf, err := s.LoadWeek(year, w)
		if err != nil {
			return nil, err
		}
		all = append(all, wf.Entries...)
	}
	return all, nil
}

func (s *Store) EntriesInRange(start time.Time, end time.Time) ([]Entry, error) {
	var all []Entry
	years, err := s.listYears()
	if err != nil {
		return nil, err
	}
	for _, year := range years {
		entries, err := s.EntriesForYear(year)
		if err != nil {
			return nil, err
		}
		for _, e := range entries {
			t, err := time.Parse("2006-01-02", e.Date)
			if err != nil {
				continue
			}
			if (t.Equal(start) || t.After(start)) && (t.Equal(end) || t.Before(end)) {
				all = append(all, e)
			}
		}
	}
	return all, nil
}
