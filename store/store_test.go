package store

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestParseDuration(t *testing.T) {
	cases := map[string]int{
		"2h30m":  150,
		"2h":     120,
		"45m":    45,
		"1h 15m": 75,
		"2H30M":  150,
	}
	for in, want := range cases {
		got, err := ParseDuration(in)
		if err != nil {
			t.Errorf("ParseDuration(%q) error: %v", in, err)
			continue
		}
		if got != want {
			t.Errorf("ParseDuration(%q) = %d, want %d", in, got, want)
		}
	}
	bad := []string{"abc", "", "2x", "0m"}
	for _, in := range bad {
		if _, err := ParseDuration(in); err == nil {
			t.Errorf("ParseDuration(%q) should fail", in)
		}
	}
}

func TestFormatDuration(t *testing.T) {
	cases := map[int]string{
		150: "2h30m",
		120: "2h",
		45:  "45m",
		0:   "0m",
	}
	for in, want := range cases {
		if got := FormatDuration(in); got != want {
			t.Errorf("FormatDuration(%d) = %q, want %q", in, got, want)
		}
	}
}

func TestStoreAddAndLoadWeek(t *testing.T) {
	dir := t.TempDir()
	old := os.Getenv("HOME")
	os.Setenv("HOME", dir)
	defer os.Setenv("HOME", old)

	st, err := NewStore()
	if err != nil {
		t.Fatal(err)
	}

	e := Entry{
		ID:       "test-1",
		Company:  "acme",
		Date:     "2026-09-08",
		Duration: "2h30m",
		Note:     "API-Refactoring",
	}
	if err := st.AddEntry(e); err != nil {
		t.Fatal(err)
	}

	year, week := ISOWeek(time.Date(2026, time.September, 8, 12, 0, 0, 0, time.Local))
	if year != 2026 || week != 37 {
		t.Fatalf("ISOWeek = %d/%d, want 2026/37", year, week)
	}

	wf, err := st.LoadWeek(year, week)
	if err != nil {
		t.Fatal(err)
	}
	if len(wf.Entries) != 1 {
		t.Fatalf("entries = %d, want 1", len(wf.Entries))
	}
	if wf.Entries[0].Company != "acme" {
		t.Errorf("company = %q, want acme", wf.Entries[0].Company)
	}
}

func TestStoreCompanies(t *testing.T) {
	dir := t.TempDir()
	old := os.Getenv("HOME")
	os.Setenv("HOME", dir)
	defer os.Setenv("HOME", old)

	st, _ := NewStore()
	if err := st.EnsureCompany("acme"); err != nil {
		t.Fatal(err)
	}
	if err := st.EnsureCompany("globex"); err != nil {
		t.Fatal(err)
	}
	if err := st.EnsureCompany("acme"); err != nil {
		t.Fatal(err)
	}
	companies, _ := st.LoadCompanies()
	if len(companies) != 2 {
		t.Fatalf("companies = %v, want 2", companies)
	}
	if err := st.RemoveCompany("acme"); err != nil {
		t.Fatal(err)
	}
	companies, _ = st.LoadCompanies()
	if len(companies) != 1 || companies[0] != "globex" {
		t.Fatalf("after remove: %v", companies)
	}
}

func TestStoreDeleteEntry(t *testing.T) {
	dir := t.TempDir()
	old := os.Getenv("HOME")
	os.Setenv("HOME", dir)
	defer os.Setenv("HOME", old)

	st, _ := NewStore()
	st.AddEntry(Entry{ID: "a", Company: "acme", Date: "2026-09-08", Duration: "1h"})
	st.AddEntry(Entry{ID: "b", Company: "globex", Date: "2026-09-08", Duration: "30m"})

	if err := st.DeleteEntry("a"); err != nil {
		t.Fatal(err)
	}

	year, week := ISOWeek(time.Date(2026, time.September, 8, 12, 0, 0, 0, time.Local))
	wf, _ := st.LoadWeek(year, week)
	if len(wf.Entries) != 1 || wf.Entries[0].ID != "b" {
		t.Fatalf("entries after delete = %+v", wf.Entries)
	}
}

func TestEntriesInRange(t *testing.T) {
	dir := t.TempDir()
	old := os.Getenv("HOME")
	os.Setenv("HOME", dir)
	defer os.Setenv("HOME", old)

	// Force store to use temp dir for data
	home, _ := os.UserHomeDir()
	st := &Store{dir: filepath.Join(home, dirName)}

	st.AddEntry(Entry{ID: "a", Company: "acme", Date: "2026-09-01", Duration: "1h"})
	st.AddEntry(Entry{ID: "b", Company: "globex", Date: "2026-09-15", Duration: "1h"})
	st.AddEntry(Entry{ID: "c", Company: "acme", Date: "2026-10-01", Duration: "1h"})

	start := time.Date(2026, time.September, 1, 0, 0, 0, 0, time.Local)
	end := time.Date(2026, time.September, 30, 23, 59, 59, 0, time.Local)
	entries, err := st.EntriesInRange(start, end)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("range entries = %d, want 2", len(entries))
	}
}
