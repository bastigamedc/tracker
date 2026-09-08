package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/bastigamedc/tracker/store"
	"github.com/bastigamedc/tracker/ui"
)

var version = "dev" // set via ldflags at release time

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--version", "-v":
			fmt.Printf("worktrack %s\n", version)
			os.Exit(0)
		case "--help", "-h":
			fmt.Println("worktrack - Zeiterfassung für mehrere Firmen")
			fmt.Println()
			fmt.Println("Aufruf: worktrack")
			fmt.Println()
			fmt.Println("Optionen:")
			fmt.Println("  -v, --version   Version anzeigen")
			fmt.Println("  -h, --help      Hilfe anzeigen")
			os.Exit(0)
		}
	}

	st, err := store.NewStore()
	if err != nil {
		fmt.Fprintln(os.Stderr, "WorkTrack:", err)
		os.Exit(1)
	}

	if err := os.MkdirAll(st.RootDir(), 0o755); err != nil {
		fmt.Fprintln(os.Stderr, "WorkTrack:", err)
		os.Exit(1)
	}

	p := tea.NewProgram(ui.New(st), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "WorkTrack:", err)
		os.Exit(1)
	}
}
