# Blood Pressure Journal — Implementation Plan

## Milestone 1: Project Foundation

Set up the core project structure, storage layer, and frontend scaffold.

### 1.1 Project Structure

- [ ] Define the target directory layout:
  ```
  cmd/
    serve.go              # `bp serve` command
  internal/
    config/
      config.go           # app config (port, db path, no-open)
      paths.go            # platform-specific data dir resolution
    domain/
      models.go           # Session, Reading, Settings types
      validation.go       # BP/pulse validation rules
      indicators.go       # traffic-light BP status logic
    storage/
      store.go            # storage interface
      toml.go             # TOML file-based implementation
    server/
      server.go           # HTTP server setup
      routes.go           # route registration
      handlers/
        dashboard.go
        measurements.go
        history.go
        reports.go
        settings.go
    report/
      averages.go         # hierarchical averaging engine
      summary.go          # report summary statistics
      report.go           # report generation orchestration
    version/
      version.go          # (exists)
  frontend/               # SvelteJS SPA
    src/
    package.json
    vite.config.js
  ```
- [ ] Remove the placeholder `cmd/hello.go` command
- [ ] Remove `urfave/cli/v2` dependency
- [ ] Rewrite `main.go` to use stdlib `flag` package:
  - Parse subcommand (`serve`) via `os.Args`
  - Use `flag.NewFlagSet("serve", ...)` for serve-specific flags
  - Print usage/help on unknown command or `--help`

### 1.2 Domain Types

- [ ] Define Go types: `MeasurementSession`, `Reading`, `Settings` (as specified in the idea)
- [ ] Implement validation rules:
  - Systolic: 70–250
  - Diastolic: 40–150
  - Pulse: 30–220
  - Warn on extreme values, but allow saving
- [ ] Implement traffic-light BP indicator logic:
  - Green: SYS < 120 AND DIA < 80
  - Yellow: SYS 120–139 OR DIA 80–89
  - Red: SYS >= 140 OR DIA >= 90
  - Dark red: SYS >= 180 OR DIA >= 120
  - Always use the worse category between SYS and DIA

### 1.3 Storage Layer

- [ ] Implement platform-specific data directory resolution:
  - Linux: `~/.local/share/bp-journal/`
  - macOS: `~/Library/Application Support/bp-journal/`
  - Fallback: `~/.local/share/bp-journal/`
- [ ] Auto-create data directory on first run
- [ ] Design TOML file structure (monthly files):
  ```
  ~/.local/share/bp-journal/
    settings.toml
    sessions/
      2026-05.toml        # all sessions for May 2026
      2026-06.toml        # all sessions for June 2026
  ```
- [ ] Implement storage interface:
  ```go
  type Store interface {
      CreateSession(session *MeasurementSession) error
      GetSession(id string) (*MeasurementSession, error)
      UpdateSession(session *MeasurementSession) error
      DeleteSession(id string) error
      ListSessions(from, to time.Time) ([]MeasurementSession, error)

      GetSettings() (*Settings, error)
      SaveSettings(settings *Settings) error
  }
  ```
- [ ] Implement TOML-based store with monthly file partitioning:
  - Derive filename from session's `measured_at` month (`YYYY-MM.toml`)
  - `CreateSession`: load month file, append session, write back
  - `ListSessions`: determine which month files overlap the date range, load only those
  - `UpdateSession` / `DeleteSession`: locate the correct month file by session date, modify in place
  - Handle edge case: session moved to a different month on edit (delete from old file, insert into new)
- [ ] Write unit tests for storage CRUD operations

### 1.4 Frontend Scaffold

- [ ] Initialize SvelteJS project in `frontend/`
- [ ] Set up Tailwind CSS
- [ ] Configure Vite to build to `frontend/dist/`
- [ ] Set up `go:embed` to embed `frontend/dist/` into the Go binary
- [ ] Create base layout component (nav bar with links to all pages)
- [ ] Set up SvelteJS router with placeholder pages:
  - `/` — Dashboard
  - `/measurements/new` — Add measurement
  - `/history` — History
  - `/reports` — Reports
  - `/settings` — Settings

---

## Milestone 2: CLI & HTTP Server

Wire up the `bp serve` command and API server.

### 2.1 Serve Command

- [ ] Implement `cmd/serve.go` using `flag.NewFlagSet`:
  - `--port` flag (default: 7391)
  - `--data` flag (default: platform data dir)
  - `--no-open` flag (default: false)
  - `--debug` flag
- [ ] Terminal output on start:
  ```
  Blood Pressure Journal is running.

  Local UI: http://127.0.0.1:7391
  Data:     ~/.local/share/bp-journal/
  
  Press Ctrl+C to stop.
  ```
- [ ] Auto-open browser (unless `--no-open`), using `open` on macOS, `xdg-open` on Linux
- [ ] Graceful shutdown on SIGINT/SIGTERM

### 2.2 HTTP Server & API Routes

- [ ] Set up `net/http` server bound to `127.0.0.1` only
- [ ] Serve embedded SvelteJS SPA for all non-API routes (SPA fallback to `index.html`)
- [ ] Register JSON API routes:
  ```
  GET    /api/sessions              — list sessions (with date range query params)
  POST   /api/sessions              — create session
  GET    /api/sessions/{id}         — get session
  PUT    /api/sessions/{id}         — update session
  DELETE /api/sessions/{id}         — delete session

  GET    /api/reports/preview       — generate report data (query params: from, to, sections)

  GET    /api/settings              — get settings
  PUT    /api/settings              — update settings
  ```
- [ ] Implement request/response JSON serialization
- [ ] Add basic error handling middleware (consistent JSON error responses)

---

## Milestone 3: Measurement CRUD (Core Feature)

The primary user flow: adding, viewing, editing, and deleting blood pressure sessions.

### 3.1 API Handlers

- [ ] `POST /api/sessions` — create a measurement session with 1+ readings
  - Validate all fields
  - Generate UUID for session and each reading
  - Return created session with computed session average
- [ ] `GET /api/sessions` — list sessions
  - Support `from` and `to` query params for date filtering
  - Return sessions ordered by `measured_at` descending
  - Include computed session averages and traffic-light status
- [ ] `GET /api/sessions/{id}` — get single session with all readings
- [ ] `PUT /api/sessions/{id}` — update session (replace readings entirely)
- [ ] `DELETE /api/sessions/{id}` — delete session (remove from monthly TOML file)

### 3.2 Measurement Form UI

- [ ] Build the `/measurements/new` page:
  - Date picker (default: today)
  - Time picker (default: now, rounded to nearest 5 min)
  - Period selector: Morning / Evening / Custom
  - Arm selector: Left / Right
  - Position selector: Sitting / Lying / Standing
  - Dynamic readings list:
    - Start with 1 reading row
    - "Add another reading" button
    - Remove reading button (if > 1 reading)
    - Fields per reading: SYS, DIA, Pulse
  - Notes textarea
  - Traffic-light indicator shown live as values are entered
  - Validation warnings for extreme values (confirmation before save)
- [ ] "Save Session" button → POST to API
- [ ] Success screen after save:
  - "Saved." confirmation
  - Session average with traffic-light indicator
  - Reading count
  - Button to add another session or go to history

### 3.3 Edit Flow

- [ ] Build the edit page (reuses measurement form component)
- [ ] Pre-populate form with existing session data
- [ ] PUT to API on save
- [ ] Redirect to history after successful edit

---

## Milestone 4: Dashboard & History

### 4.1 Dashboard (`/`)

- [ ] Show today's sessions (if any) with averages and traffic-light
- [ ] Show quick stats:
  - Last reading (value + time)
  - Today's average (if multiple sessions)
  - This week's average
- [ ] Quick-add button linking to `/measurements/new`
- [ ] Empty state for new users (welcome message + prompt to add first reading)

### 4.2 History Page (`/history`)

- [ ] Paginated/scrollable list of measurement sessions
- [ ] Each row shows: date, time, period, session average BP, pulse, traffic-light, reading count
- [ ] Date range filter (from/to)
- [ ] Click to expand → show individual readings
- [ ] Edit button → navigate to edit page
- [ ] Delete button → confirmation dialog → DELETE API call
- [ ] Group sessions by day with day averages shown as section headers

---

## Milestone 5: Reporting Engine & Report UI

### 5.1 Hierarchical Averaging Engine

- [ ] Implement `report/averages.go`:
  - `SessionAverage(session)` — average of readings in a session
  - `DayAverages(sessions)` — average of session averages per calendar day
  - `WeekAverages(dayAvgs)` — average of day averages per ISO week
  - `MonthAverages(dayAvgs)` — average of day averages per calendar month
  - Morning/evening split averages
- [ ] Round displayed averages to whole numbers; keep internal decimal precision
- [ ] Write unit tests for the averaging logic (use the spec example as a test case)

### 5.2 Report Summary

- [ ] Implement `report/summary.go`:
  - Days measured, session count, raw reading count
  - Overall / morning / evening averages
  - Highest/lowest systolic, diastolic, pulse
- [ ] Traffic-light status for all averages

### 5.3 Report UI (`/reports`)

- [ ] Report configuration form:
  - Date range picker (from / to)
  - Section toggles (checkboxes):
    - Summary
    - Monthly averages
    - Weekly averages
    - Daily averages
    - Session averages
    - Raw measurements
    - Notes
- [ ] "Preview Report" button → fetch report data from API → render in-page
- [ ] Report preview with all selected sections rendered as formatted tables
- [ ] Print-optimized CSS (`@media print`):
  - Clean, doctor-friendly layout
  - Header with patient name, date range, generation date
  - Footer disclaimer: "This report summarizes self-measured blood pressure records and is not a medical diagnosis."
  - Traffic-light indicators rendered as text labels in print
- [ ] "Print / Save as PDF" button → triggers `window.print()`

---

## Milestone 6: Settings & Polish

### 6.1 Settings Page (`/settings`)

- [ ] Configurable defaults:
  - Patient name (used in report header)
  - Default arm (Left / Right)
  - Default position (Sitting / Lying / Standing)
  - Default device name
- [ ] Settings persisted in `settings.toml`
- [ ] Settings applied as defaults in measurement form (user can override per session)

### 6.2 UI Polish

- [ ] Responsive design (usable on phone browser for bedside logging)
- [ ] Consistent color scheme and typography
- [ ] Loading states for API calls
- [ ] Error toast notifications
- [ ] Keyboard navigation support in measurement form (tab through fields)
- [ ] Auto-focus first SYS field when adding a new reading

### 6.3 Build & Distribution

- [ ] Update Makefile with frontend build step (`make build` runs both `npm run build` and `go build`)
- [ ] Update `.goreleaser.yml` for cross-platform binary releases
- [ ] Update Dockerfile if applicable
- [ ] Update README.md with usage instructions
- [ ] Test on macOS and Linux

---

## Implementation Order

The recommended sequence:

```
Milestone 1 (Foundation)  ←  start here
    ↓
Milestone 2 (CLI + Server)
    ↓
Milestone 3 (Measurement CRUD)  ←  first usable feature
    ↓
Milestone 4 (Dashboard + History)
    ↓
Milestone 5 (Reports)
    ↓
Milestone 6 (Settings + Polish)
```

Each milestone produces a working (if incomplete) application.
After Milestone 3, the app is usable for daily logging.
After Milestone 5, it covers the full spec.
