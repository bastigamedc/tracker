# WorkTrack

WorkTrack ist ein Terminal-basiertes Tool zur Zeiterfassung für Remote-Arbeit bei mehreren Firmen. Die Bedienung erfolgt komplett über eine interaktive TUI (Terminal-UI) - ohne Server, ohne Cloud, ohne Abo.

## Funktionen

- **Zeit eintragen** - geführter Eintragungs-Assistent (Firma → Datum → Dauer → Notiz)
- **Tägliche Zusammenfassung** - zeigt alle Einträge des heutigen Tages, inklusive Löschen
- **Monatsreport** - Summen pro Firma und Monat
- **Firmen verwalten** - Firmen anlegen und entfernen
- **Lokale Speicherung** - alle Daten liegen als einfache JSON-Dateien auf deinem Rechner

## Installation

### macOS (Homebrew)

```bash
brew tap bastigamedc/homebrew-worktrack
brew install --cask worktrack
```

### Windows / manuell

Lade das passende Archiv von der [Releases-Seite](https://github.com/bastigamedc/tracker/releases):

- Windows: `worktrack_<version>_windows_amd64.zip` → `worktrack.exe` entpacken und ins PATH-Verzeichnis legen
- macOS / Linux: `worktrack_<version>_<os>_<arch>.tar.gz` → `worktrack` entpacken und installieren, z.B.:

```bash
tar -xzf worktrack_1.0.0_darwin_arm64.tar.gz
sudo mv worktrack /usr/local/bin/
```

Verfügbare Varianten: `linux_amd64`, `linux_arm64`, `darwin_amd64`, `darwin_arm64`, `windows_amd64`.

### Vom Quellcode bauen

Voraussetzung: [Go](https://go.dev/dl/) 1.27+

```bash
git clone https://github.com/bastigamedc/tracker.git
cd tracker
go build -o worktrack .
```

## Verwendung

Einfach `worktrack` im Terminal ausführen:

```bash
worktrack
```

Daraufhin öffnet sich das Hauptmenü:

| Menüpunkt          | Beschreibung                                  |
| ------------------ | --------------------------------------------- |
| Neue Zeit eintragen | Geführter Assistent zum Erfassen einer Arbeitszeit |
| Heutige Zusammenfassung | Alle Einträge von heute ansehen und löschen |
| Monatsreport       | Zeitsummen pro Firma für den aktuellen Monat  |
| Firmen verwalten   | Firmen anlegen oder entfernen                 |

Navigation: `↑`/`↓` Auswahl, `Enter` öffnen, `Esc` zurück, `q` (oder `Ctrl+C`) beenden.

Beim ersten Start gibt es noch keine Firmen - lege über *Firmen verwalten* zuerst eine Firma an, danach kannst du Zeiten erfassen.

## Daten & Speicherort

Alle Daten liegen unter `~/.worktrack/`:

```
~/.worktrack/
├── companies.json            # Liste der Firmen
└── data/
    └── 2026/
        └── week-37.json      # Einträge der ISO-Woche 37/2026
```

Die Dateien sind einfaches JSON und lassen sich mit jedem Editor ansehen, sichern oder auf andere Geräte übertragen:

```json
{
  "week": 42,
  "year": 2026,
  "entries": [
    {
      "id": "01J...",
      "company": "acme",
      "date": "2026-10-14",
      "duration": "1h30m",
      "note": "Bugfix Auth"
    }
  ]
}
```

- **Dauer**: Wird im Format `1h30m`, `45m`, `2h` eingegeben.
- **Datum**: Format `DD.MM.JJJJ` (z.B. `14.10.2026`).

## Befehlszeilenoptionen

```text
Aufruf: worktrack

Optionen:
  -v, --version   Version anzeigen
  -h, --help      Hilfe anzeigen
```

## Entwicklung

```bash
go build ./...    # Projekt kompilieren
go test ./...     # Tests ausführen
go vet ./...      # Statische Analyse
```

Releases werden via GitHub Actions + GoReleaser für alle Plattformen gebaut. Beim Taggen (`git tag v1.0.0 && git push origin v1.0.0`) entsteht automatisch ein Release; der Homebrew-Cask wird dabei in den Tap `bastigamedc/homebrew-worktrack` gepusht.

## Lizenz

MIT License - siehe [LICENSE](LICENSE).