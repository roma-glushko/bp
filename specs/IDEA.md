# Blood Pressure Journal — Technical Plan

## Goal

Build a local-first CLI application that starts a private local web UI for logging blood pressure readings and generating printable doctor reports.

The CLI is only an app launcher. All user interaction happens through the browser UI.

```bash
bp serve
The app starts a local server:
http://127.0.0.1:7391
```

The user can then:

- Add blood pressure measurement sessions
- Record multiple readings per session
- View history
- Edit/delete records
- Generate PDF reports for a selected time range
- Print/export reports for doctors

# Core Product Requirements

1. Local-first
   Data is stored locally on the user’s machine.
   No account system.
   No cloud sync in the MVP.
   Server binds to 127.0.0.1 only.
   Data is stored as TOML files for human readability and easy sync via git or Dropbox.

Default data location:

```
~/.local/share/bp-journal/
```

On macOS, optionally:

```
~/Library/Application Support/bp-journal/
```

2. CLI Scope

The CLI should only start the local UI server.

Required command:

```bash
bp serve
```

Optional flags:

```
bp serve --port 7391
bp serve --data ~/.local/share/bp-journal/
bp serve --no-open
```

## Recommended Tech Stack

### Language: 

Use Go.

Reasons:

- Easy to ship as a single binary
- Good standard library
- Simple local HTTP server
- Easy file I/O for TOML-based storage
- Good fit for CLI + local web app

No CLI commands are needed for adding readings or exporting reports.

### CLI

Use: stdlib flags

Only one real command is needed:

```
serve
```

### HTTP Server

Use Go standard library:

```
net/http
```

The frontend is a SvelteJS SPA embedded in the Go binary via `go:embed`.

The Go server serves the SPA static files and exposes a JSON API under `/api/`.

### Database

Use TOML files.

Reason:

- easy to read as just text file
- easy to sync via git or dropbox

### Frontend

SvelteJS 5 SPA with Tailwind CSS v4, using hash-based routing (`/#/history`).

Built with Vite, output embedded into Go binary via `go:embed`.

### PDF Generation

Use just print-CSS styles and ctrl + P type of trigger to open print window in browser and save the page as PDF

## Main User Flow

Start app

```
bp serve
```

Terminal output:

```
Blood Pressure Journal is running.

Local UI: http://127.0.0.1:7391
Data:     ~/.local/share/bp-journal/

Press Ctrl+C to stop.

The browser opens automatically unless --no-open is passed.
```

## UI Pages

```
/
Dashboard

/measurements/new
Add measurement session

/history
View, edit, delete measurement sessions

/reports
Generate doctor report

/settings
Configure defaults
```

## Measurement Model

The app should model a measurement session, not only individual readings.

A session represents one measurement moment, for example:

```
May 7, 09:00, morning
Try 1: 126 / 86, pulse 78
Try 2: 119 / 83, pulse 81
```

The app calculates a session average from the readings.

## Measurement Form

The /measurements/new page should contain:

```
Date:        [2026-05-12]
Time:        [09:00]
Period:      [Morning / Evening / Custom]
Arm:         [Left / Right]
Position:    [Sitting / Lying / Standing]

Readings:

#   SYS   DIA   Pulse
1   [126] [86]  [78]
2   [119] [83]  [81]

[+ Add another reading]

Notes:
[________________________________]

[Save Session]
```

After saving, show:

```
Saved.

Session average:
123 / 85, pulse 80

This session has 2 readings.
```

### Average Calculation Rules

Use a clear hierarchical averaging model.

#### Reading

One physical measurement attempt.

Example:

```
126 / 86, pulse 78
```

#### Session Average

Average of all readings inside one measurement session.

Example:

```
Morning session:
Try 1: 126 / 86, pulse 78
Try 2: 119 / 83, pulse 81

Session average:
123 / 85, pulse 80
```

#### Day Average

Average of session averages for the same calendar day.

Do not average all raw readings directly.

Reason:

A day with more repeated attempts should not dominate the day/week/month statistics.

#### Week Average

Average of day averages within the same ISO week.

#### Month Average

Average of day averages within the same calendar month.

#### Example Averaging

Raw data:

```
May 7 morning
Try 1: 126 / 86, pulse 78
Try 2: 119 / 83, pulse 81

May 7 evening
Try 1: 119 / 78, pulse 75
Try 2: 119 / 80, pulse 78
```

Session averages:

```
May 7 morning: 123 / 85, pulse 80
May 7 evening: 119 / 79, pulse 77
```

Day average:

```
May 7: 121 / 82, pulse 79
```

For display, round averages to whole numbers.

Internally, keep decimal precision.

## Data Storage

All data is stored as TOML files in the data directory.

### Directory Structure

```
~/.local/share/bp-journal/
  settings.toml
  sessions/
    2026-05.toml
    2026-06.toml
```

### Monthly Session File Format

Each month has a single TOML file containing all measurement sessions for that month.

Example: `sessions/2026-05.toml`

```toml
[[sessions]]
id = "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
measured_at = 2026-05-07T09:00:00Z
period = "morning"
arm = "left"
position = "sitting"
cuff_location = ""
device_name = ""
symptoms = []
context_notes = []
notes = "felt well rested"
created_at = 2026-05-07T09:05:00Z
updated_at = 2026-05-07T09:05:00Z

  [[sessions.readings]]
  id = "r1-uuid"
  reading_no = 1
  systolic = 126
  diastolic = 86
  pulse = 78
  created_at = 2026-05-07T09:05:00Z

  [[sessions.readings]]
  id = "r2-uuid"
  reading_no = 2
  systolic = 119
  diastolic = 83
  pulse = 81
  created_at = 2026-05-07T09:05:00Z

[[sessions]]
id = "b2c3d4e5-f6a7-8901-bcde-f12345678901"
measured_at = 2026-05-07T22:00:00Z
period = "evening"
arm = "left"
position = "sitting"
notes = ""
created_at = 2026-05-07T22:05:00Z
updated_at = 2026-05-07T22:05:00Z

  [[sessions.readings]]
  id = "r3-uuid"
  reading_no = 1
  systolic = 119
  diastolic = 78
  pulse = 75
  created_at = 2026-05-07T22:05:00Z

  [[sessions.readings]]
  id = "r4-uuid"
  reading_no = 2
  systolic = 119
  diastolic = 80
  pulse = 78
  created_at = 2026-05-07T22:05:00Z
```

### Settings File Format

```toml
[settings]
patient_name = "Roman Hlushko"
default_arm = "left"
default_position = "sitting"
device_name = ""
updated_at = 2026-05-07T09:05:00Z
```

### Validation Rules

Systolic

Allowed range: 70–250

Diastolic

Allowed range: 40–150

Pulse

Allowed range: 30–220

These are validation ranges, not medical diagnosis rules.

The user should be able to save unusual values, but the UI should ask for confirmation if values look extreme.


## API Routes

All API routes return JSON. The SPA is served as static files at the root.

```
GET    /api/sessions              — list sessions (query: from, to)
POST   /api/sessions              — create session
GET    /api/sessions/{id}         — get session
PUT    /api/sessions/{id}         — update session
DELETE /api/sessions/{id}         — delete session

GET    /api/reports/preview       — generate report data (query: from, to, sections)

GET    /api/settings              — get settings
PUT    /api/settings              — update settings
```

## Report Generation

The report page should allow selecting a time range.

```
Date from: [2026-05-01]
Date to:   [2026-05-31]

Include:
[x] Summary
[x] Monthly averages
[x] Weekly averages
[x] Daily averages
[x] Session averages
[x] Raw measurements
[x] Notes

[Download PDF]
```

## PDF Report Structure

The PDF should be printable and doctor-friendly.

```
Blood Pressure Report

Patient: Roman Hlushko
Period: 2026-05-01 — 2026-05-31
Generated: 2026-05-12

1. Summary

2. Monthly Averages

3. Weekly Averages

4. Daily Averages

5. Session Averages

6. Raw Measurements

7. Symptoms and Notes

Footer:
This report summarizes self-measured blood pressure records and is not a medical diagnosis.
```

## Report Summary Section

Include:

Days measured
Number of sessions
Number of raw readings

Overall average
Morning average
Evening average

Highest systolic reading
Highest diastolic reading
Lowest systolic reading
Lowest diastolic reading

Average pulse
Highest pulse
Lowest pulse

### Report Tables

Monthly averages

Month      Avg BP     Avg Pulse   Days Measured   Sessions
2026-05    120/80     78          14              22

Weekly averages

Week        Avg BP     Avg Pulse   Days Measured   Sessions
2026-W19    121/81     79          5               8
2026-W20    119/78     77          4               6

Daily averages

Date         Avg BP     Avg Pulse   Sessions   Notes
2026-05-07   121/82     79          2          headache
2026-05-08   118/78     81          1          congestion

Session averages

Date         Time    Period    Avg BP     Avg Pulse   Attempts
2026-05-07   09:00   Morning   123/85     80          2
2026-05-07   22:00   Evening   119/79     77          2

Raw measurements
Date         Time    Period    Attempt   SYS   DIA   Pulse
2026-05-07   09:00   Morning   1         126   86    78
2026-05-07   09:00   Morning   2         119   83    81

### Core Go Types

```
type Reading struct {
ID         string
SessionID  string
AttemptNo  int
Systolic   int
Diastolic  int
Pulse      *int
CreatedAt  time.Time
}

type MeasurementSession struct {
ID           string
MeasuredAt   time.Time
Period       string

    Arm          string
    Position     string
    CuffLocation string
    DeviceName   string

    Symptoms     []string
    ContextNotes []string
    Notes        string

    Readings     []Reading

    CreatedAt    time.Time
    UpdatedAt    time.Time
}
```

### Report Types

```
type ReportRange struct {
From time.Time
To   time.Time
}

type SessionAverage struct {
SessionID    string
MeasuredAt   time.Time
LocalDate    string
Period       string

    AvgSystolic  float64
    AvgDiastolic float64
    AvgPulse     *float64

    Attempts     int
    Notes        string
    Symptoms     []string
}

type DayAverage struct {
Date         string
AvgSystolic  float64
AvgDiastolic float64
AvgPulse     *float64
Sessions     int
}

type WeekAverage struct {
Year         int
Week        int
AvgSystolic  float64
AvgDiastolic float64
AvgPulse     *float64
DaysMeasured int
Sessions     int
}

type MonthAverage struct {
Year         int
Month        int
AvgSystolic  float64
AvgDiastolic float64
AvgPulse     *float64
DaysMeasured int
Sessions     int
}

type ReportSummary struct {
DaysMeasured int
Sessions     int
RawReadings  int

    OverallAvgSystolic  float64
    OverallAvgDiastolic float64
    OverallAvgPulse     *float64

    MorningAvgSystolic  *float64
    MorningAvgDiastolic *float64
    MorningAvgPulse     *float64

    EveningAvgSystolic  *float64
    EveningAvgDiastolic *float64
    EveningAvgPulse     *float64

    HighestSystolic  int
    HighestDiastolic int
    LowestSystolic   int
    LowestDiastolic  int

    HighestPulse *int
    LowestPulse  *int
}

type BPReport struct {
Range       ReportRange
GeneratedAt time.Time

    Summary     ReportSummary

    Sessions    []SessionAverage
    Days        []DayAverage
    Weeks       []WeekAverage
    Months      []MonthAverage
    RawReadings []Reading
}
```

## Traffic-Light Blood Pressure Indicators

The UI should show a simple traffic-light status for each reading and average.

This should appear in:

- Measurement form after values are entered
- Saved session summary
- History table
- Dashboard
- Report preview
- PDF report

The traffic-light status is only a visual aid. It is not a diagnosis.

---

## Default BP Indicator Rules

Use this simplified default model:

| Status | Color | Meaning | Rule |
|---|---|---|---|
| Great | Green | Blood pressure looks good | SYS < 120 and DIA < 80 |
| Elevated | Yellow | Blood pressure is elevated | SYS 120–139 or DIA 80–89 |
| Too high | Red | Blood pressure is high | SYS >= 140 or DIA >= 90 |
| Very high | Dark red | Potentially urgent reading | SYS >= 180 or DIA >= 120 |

The app should always use the worse category between systolic and diastolic.

Example:

```text
119 / 78 -> Green
126 / 76 -> Yellow
118 / 86 -> Yellow
142 / 82 -> Red
135 / 92 -> Red
181 / 95 -> Dark red
150 / 121 -> Dark red
