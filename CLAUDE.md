# μCal Development Documentation

This document details the development process, technical decisions, and bug fixes for μCal, a minimal read-only CalDAV calendar viewer built with Claude Code.

## Project Overview

μCal is a web application that displays calendar events from one or more CalDAV servers in a clean, responsive interface. It was built as a minimal viable product with the following core requirements:

- Read-only calendar view (no editing)
- Svelte 5 frontend with Bootstrap 5
- Go backend with embedded frontend
- Docker deployment with GitHub Actions CI/CD
- Support for recurring events (RRULE)
- Timezone handling
- Multiple calendar support with color coding

## Architecture

### Technology Stack

**Backend:**
- Go 1.23+ (using `golang:alpine` for latest stable)
- `github.com/emersion/go-webdav` - CalDAV protocol client
- `github.com/emersion/go-ical` - iCalendar parsing
- `github.com/teambition/rrule-go` - Recurring rule expansion
- `gopkg.in/yaml.v3` - Configuration parsing
- Standard library HTTP server (no external framework)

**Frontend:**
- Svelte 5 (v5.46.3) with runes ($state, $derived, $effect)
- Bootstrap 5.3.8 for UI components
- TypeScript for type safety
- date-fns 4.1.0 for date manipulation
- Vite 7.3.1 for build tooling
- `node:alpine` for builds (latest LTS)

**Deployment:**
- Multi-stage Dockerfile (Node build → Go build → Alpine runtime)
- GitHub Actions for automated releases
- Version injection via build-time ldflags

### Project Structure

```
mucal/
├── cmd/mucal/               # Application entry point
│   └── main.go
├── internal/
│   ├── api/                 # REST API handlers
│   │   ├── handler.go       # Endpoint implementations
│   │   └── middleware.go    # Logging, CORS, recovery
│   ├── caldav/              # CalDAV integration
│   │   ├── client.go        # CalDAV client with auth
│   │   ├── event.go         # Event structures and sorting
│   │   └── recurring.go     # RRULE expansion logic
│   ├── config/              # Configuration management
│   │   └── config.go        # YAML parsing, validation
│   └── version/             # Version information
│       └── version.go
├── web/                     # Svelte 5 frontend
│   ├── public/
│   │   └── favicon.svg      # Custom SVG favicon
│   ├── src/
│   │   ├── lib/
│   │   │   ├── components/  # Svelte components
│   │   │   │   ├── DayColumn.svelte
│   │   │   │   ├── ErrorPopup.svelte
│   │   │   │   ├── EventCard.svelte
│   │   │   │   ├── Header.svelte
│   │   │   │   ├── MonthCalendar.svelte
│   │   │   │   └── WeekView.svelte
│   │   │   ├── services/
│   │   │   │   └── api.ts   # API client
│   │   │   ├── stores/
│   │   │   │   ├── calendar.svelte.ts  # State management
│   │   │   │   └── error.svelte.ts
│   │   │   ├── types.ts     # TypeScript interfaces
│   │   │   └── utils/
│   │   │       └── date.ts  # Date utilities
│   │   ├── app.css          # Global styles
│   │   ├── App.svelte       # Root component
│   │   └── main.ts          # Entry point
│   ├── index.html
│   ├── package.json
│   └── vite.config.ts
├── embed.go                 # Go embed directive
├── Dockerfile               # Multi-stage build
├── Makefile                 # Build automation
├── .github/workflows/
│   └── release.yml          # CI/CD pipeline
├── config.example.yaml      # Example configuration
├── .gitignore
├── .dockerignore
├── .editorconfig
└── README.md
```

## Implementation Timeline

### Phase 1: Backend Foundation

1. **Go Module Initialization**
   - Created `go.mod` with module path
   - Added dependencies for CalDAV, iCalendar, RRULE, and YAML

2. **Configuration Module** (`internal/config/config.go`)
   - YAML parsing with `gopkg.in/yaml.v3`
   - Password file reading (separate files for security)
   - Timezone validation using `time.LoadLocation()`
   - Sanitized config for frontend (strips credentials)

3. **CalDAV Client** (`internal/caldav/client.go`)
   - Custom HTTP transport for Basic Authentication
   - Calendar event fetching with date ranges
   - Timezone conversion (all times to configured timezone)
   - All-day event detection (DTSTART with VALUE=DATE)
   - Parallel calendar fetching with goroutines

4. **Recurring Events** (`internal/caldav/recurring.go`)
   - RRULE parsing and expansion using `rrule-go`
   - EXDATE handling (excluded dates)
   - Date vs DateTime handling for all-day events
   - Duration parsing (ISO 8601 format)
   - 1000 occurrence limit to prevent infinite expansion
   - 30-day buffer for occurrence generation

5. **Event Structures** (`internal/caldav/event.go`)
   - Event struct with JSON tags
   - Custom sorting: all-day first, then chronological, then alphabetical
   - Support for recurring event flags

### Phase 2: REST API

6. **API Handlers** (`internal/api/handler.go`)
   - `GET /api/health` - Health check with version
   - `GET /api/config` - Sanitized configuration
   - `GET /api/events` - Events for date range with parallel fetching
   - `GET /api/events/month` - Days with events for month calendar
   - Error handling with JSON responses

7. **Middleware** (`internal/api/middleware.go`)
   - Request logging to stderr
   - Panic recovery
   - CORS headers for development

8. **Main Application** (`cmd/mucal/main.go`)
   - CLI flag parsing with `-config` flag
   - Positional argument support for config path
   - HTTP server setup on port 8080
   - Graceful shutdown handling
   - Version injection via ldflags

### Phase 3: Frontend Foundation

9. **Svelte + Vite Setup**
   - Initialized with Vite and Svelte 5
   - TypeScript configuration
   - Bootstrap 5 integration via CDN

10. **State Management** (`web/src/lib/stores/`)
    - Calendar store using Svelte 5 runes
    - `$state` for reactive variables
    - `$derived` for computed values
    - `$effect` for auto-refresh functionality
    - Error store for toast notifications

11. **API Service** (`web/src/lib/services/api.ts`)
    - Fetch wrappers for all endpoints
    - Error handling with JSON parsing
    - Type-safe responses

12. **Root Layout** (`web/src/App.svelte`)
    - Single-column centered layout (max-width 800px)
    - Collapsible month calendar with toggle button
    - Week view below calendar
    - Loading states
    - Error popup integration

### Phase 4: UI Components

13. **EventCard** (`web/src/lib/components/EventCard.svelte`)
    - Compact two-line format
    - First line: time/calendar name
    - Second line: summary with location
    - Text truncation with ellipsis
    - Hover tooltip for full details
    - Left border color stripe
    - Current event highlighting

14. **DayColumn** (`web/src/lib/components/DayColumn.svelte`)
    - Day section card with header
    - Day name, date, and event count
    - Today highlighting (blue background)
    - Vertical event list
    - "No events" message

15. **WeekView** (`web/src/lib/components/WeekView.svelte`)
    - Gradient header (purple)
    - Previous/Today/Next navigation
    - Week date range display
    - Scrollable days list (single scroll window)
    - Loading spinner
    - Responsive button text (hidden on mobile)

16. **MonthCalendar** (`web/src/lib/components/MonthCalendar.svelte`)
    - Month grid (Monday-Sunday)
    - Previous/Next month navigation
    - Day click to select week
    - Days with events marked (blue dot)
    - Today highlighting
    - Other month days grayed out

17. **Header** (`web/src/lib/components/Header.svelte`)
    - App name (μCal)
    - Version badge
    - Shadow styling
    - Responsive font sizes

18. **ErrorPopup** (`web/src/lib/components/ErrorPopup.svelte`)
    - Bootstrap toast notification
    - Auto-dismiss after 5 seconds
    - Manual close button
    - Error message display

### Phase 5: Integration

19. **State-API Integration**
    - Connected stores to API endpoints
    - Implemented auto-refresh with `$effect`
    - Week navigation methods
    - Event filtering by day
    - Month event days loading

20. **Error Handling**
    - Try-catch blocks around API calls
    - Error store updates
    - User-friendly error messages
    - Toast notifications

21. **Current Event Detection**
    - Check if current time is between event start and end
    - Subtle background highlight
    - Real-time updates

22. **Responsive Design**
    - Mobile-first approach
    - Media queries for tablet and desktop
    - Single-column layout for all viewports
    - Touch-friendly controls
    - No horizontal overflow

### Phase 6: Embedding & Docker

23. **Frontend Build Integration**
    - Vite production build to `web/dist/`
    - Asset optimization and minification
    - CSS and JS bundling

24. **Go Embed** (`embed.go`)
    - `//go:embed web/dist/*` directive
    - Static file serving from embedded FS
    - SPA routing (serve index.html for all paths)

25. **Multi-stage Dockerfile**
    - Stage 1: Node build (frontend)
    - Stage 2: Go build (backend with embedded frontend)
    - Stage 3: Alpine runtime (minimal final image)
    - Version injection via build args

26. **Local Testing**
    - Docker build and run tests
    - Volume mounts for config and secrets
    - Port mapping (8080:8080)

### Phase 7: CI/CD

27. **GitHub Actions Workflow** (`.github/workflows/release.yml`)
    - Trigger on tag push (v*)
    - Multi-stage Docker build
    - Version extraction from git tag
    - Push to ghcr.io with version tag and latest
    - Automated semantic versioning

28. **Version Injection**
    - Build-time ldflags for Go binary
    - Version displayed in UI header
    - Version in /api/health endpoint

### Phase 8: Polish & Documentation

29. **Configuration Examples** (`config.example.yaml`)
    - Sample calendars with different colors
    - Timezone configuration
    - Auto-refresh settings
    - Password file references

30. **Documentation**
    - Comprehensive README.md
    - Installation instructions (Docker, source)
    - Configuration guide
    - Usage examples
    - API documentation
    - Troubleshooting guide

31. **Build Automation** (`Makefile`)
    - 20+ targets with colored output
    - `help` - Show all targets
    - `build` - Full build
    - `run` - Build and run
    - `cleanup` - Remove all artifacts and dependencies
    - `update` - Update all dependencies
    - `docker-build` / `docker-run` - Docker operations
    - `dev-frontend` / `dev-backend` - Development servers

32. **Final Testing**
    - End-to-end testing with real CalDAV server
    - Recurring event validation
    - Multiple calendars
    - Timezone handling
    - Error scenarios

## Major Technical Challenges & Solutions

### 1. Recurring Events with Timezones

**Problem:** All-day recurring events were showing one day off when using RRULE expansion.

**Root Cause:**
- Backend was using DATETIME format with UTC for all events
- `rrule-go` interpreted all-day events with timestamps, causing timezone shifts
- Frontend used `toISOString()` which converted to UTC, shifting dates

**Solution:**
- Backend: Use DATE format (YYYYMMDD) for all-day events in RRULE
- Frontend: Compare dates using local date parts instead of UTC strings

**Code Changes:**
```go
// recurring.go - Use DATE format for all-day events
if allDay {
    dtstart = startTime.Format("20060102")  // DATE
} else {
    dtstart = startTime.UTC().Format("20060102T150405Z")  // DATETIME
}
```

```typescript
// calendar.svelte.ts - Use local date parts
getEventsForDay(date: Date): Event[] {
  const year = date.getFullYear();
  const month = date.getMonth();
  const day = date.getDate();

  return this.events.filter(event => {
    const eventDate = new Date(event.start);
    return eventDate.getFullYear() === year &&
           eventDate.getMonth() === month &&
           eventDate.getDate() === day;
  });
}
```

### 2. iCalendar Text Escaping

**Problem:** Event summaries showed literal backslashes before commas (e.g., `fondi\, banca`).

**Root Cause:** iCalendar format escapes special characters per RFC 5545:
- `\,` → comma
- `\;` → semicolon
- `\n` → newline
- `\\` → backslash

**Solution:** Implement proper unescaping for TEXT values.

**Code Changes:**
```go
// client.go - Unescape function
func unescapeICalText(s string) string {
    var result strings.Builder
    for i := 0; i < len(s); i++ {
        if s[i] == '\\' && i+1 < len(s) {
            next := s[i+1]
            switch next {
            case ',', ';', '\\':
                result.WriteByte(next)
                i++
            case 'n', 'N':
                result.WriteByte('\n')
                i++
            default:
                result.WriteByte('\\')
            }
        } else {
            result.WriteByte(s[i])
        }
    }
    return result.String()
}

// Apply to SUMMARY, DESCRIPTION, LOCATION
summary = unescapeICalText(prop.Value)
```

### 3. Horizontal Overflow on Mobile

**Problem:** UI broke on mobile with content overflowing horizontally.

**Root Cause:**
- Default CSS allowed elements to exceed viewport width
- No max-width constraints on containers
- Month calendar table used auto layout

**Solution:**
- Added `overflow-x: hidden` globally
- Set `max-width: 100%` on containers
- Used `table-layout: fixed` for calendar
- Ensured all containers respect viewport width

**Code Changes:**
```css
/* app.css - Prevent overflow */
html, body {
  max-width: 100%;
  overflow-x: hidden;
}

#app {
  width: 100%;
  max-width: 100vw;
  overflow-x: hidden;
}

.container, .container-fluid {
  max-width: 100%;
  overflow-x: hidden;
}

/* MonthCalendar.svelte - Fixed table layout */
.table {
  width: 100%;
  table-layout: fixed;
}
```

### 4. Accessibility Warnings

**Problem:** Svelte build warnings about buttons without accessible labels.

**Root Cause:** Icon-only buttons need `aria-label` for screen readers.

**Solution:** Added `aria-label` to all icon-only buttons.

**Code Changes:**
```svelte
<button aria-label="Close error message">
  <i class="bi bi-x"></i>
</button>

<button aria-label="Previous month">
  <i class="bi bi-chevron-left"></i>
</button>
```

### 5. Config Path Flexibility

**Problem:** Default config path `/config/config.yaml` only worked in Docker.

**Root Cause:** Hardcoded path for Docker deployment didn't work for local development.

**Solution:**
- Changed default to `config.yaml` (current directory)
- Added positional argument support
- Docker CMD still uses `/config/config.yaml`

**Code Changes:**
```go
// main.go
configPath := flag.String("config", "config.yaml", "Path to configuration file")
flag.Parse()

// Support positional argument
if len(flag.Args()) > 0 {
    *configPath = flag.Args()[0]
}
```

### 6. UI Layout Redesign

**Problem:** Initial 7-column grid layout didn't work well and broke on mobile.

**User Feedback:** "doesn't work at all", wanted single-column vertical layout.

**Solution:** Complete redesign to single-column layout:
- Month calendar at top (collapsible)
- Vertical stack of day sections below
- Single scroll window for entire page
- Compact two-line event cards

**Major Changes:**
- `App.svelte` - Single column with toggle button
- `WeekView.svelte` - Vertical stack, removed separate scroll
- `DayColumn.svelte` - Day section card with header
- `EventCard.svelte` - Two-line compact format

### 7. Favicon Design

**Problem:** First favicon had blue calendar on blue background.

**User Feedback:** "blue on blue...? What were you thinking?"

**Solution:**
- Changed background from blue circle to subtle gray gradient
- Kept calendar icon with blue header
- Added μ symbol in center

**Code:**
```svg
<circle cx="64" cy="64" r="60" fill="url(#bg)"/>
<defs>
  <linearGradient id="bg">
    <stop offset="0%" stop-color="#f8f9fa"/>
    <stop offset="100%" stop-color="#e9ecef"/>
  </linearGradient>
</defs>
```

## Key Design Decisions

### 1. No Database
- Events fetched directly from CalDAV on each request
- No persistent caching
- Simpler deployment, no migration concerns
- Acceptable for personal/small-team use

### 2. Embedded Frontend
- Frontend built and embedded in Go binary
- Single binary deployment
- No need to serve static files separately
- Simplified Docker image

### 3. Port 8080 Fixed
- No configuration for port
- Use reverse proxy (nginx, Traefik) for HTTPS and custom ports
- Simplifies configuration

### 4. Password Files
- Passwords in separate files, not in config YAML
- Better for secrets management (Docker secrets, Kubernetes secrets)
- Prevents accidental commit of credentials
- Files should be chmod 600

### 5. Svelte 5 Runes
- Modern reactive primitives (`$state`, `$derived`, `$effect`)
- Simpler than stores for this use case
- Better TypeScript support
- Clearer reactivity model

### 6. Bootstrap via CDN
- No build-time CSS processing
- Faster frontend builds
- Standard UI components
- Smaller bundle size

### 7. Single Scroll Window
- Page-level scrolling instead of component-level
- Better mobile experience
- More natural navigation
- Simpler CSS

### 8. Makefile for Build Automation
- Cross-platform build commands
- Colored output for clarity
- Consistent development workflow
- Easy dependency management

### 9. Context-Aware Past Day Hiding
- Hide past days only in current week
- Show all days in past/future weeks
- Expandable with one click when needed
- Focuses user attention on relevant events
- Prevents unnecessary scrolling on weekends

## API Design

### Endpoint: GET /api/health

**Response:**
```json
{
  "status": "ok",
  "version": "v0.1.0"
}
```

**Purpose:** Health checks and version display.

### Endpoint: GET /api/config

**Response:**
```json
{
  "time_zone": "Europe/Rome",
  "auto_refresh": 60,
  "calendars": [
    {
      "name": "Personal",
      "color": "#4ECDC4"
    }
  ]
}
```

**Purpose:** Frontend configuration (credentials stripped).

### Endpoint: GET /api/events

**Query Parameters:**
- `start` - Start date (YYYY-MM-DD)
- `end` - End date (YYYY-MM-DD)

**Response:**
```json
[
  {
    "uid": "event-123",
    "summary": "Meeting",
    "description": "Team sync",
    "location": "Office",
    "start": "2026-01-15T10:00:00+01:00",
    "end": "2026-01-15T11:00:00+01:00",
    "allDay": false,
    "calendarName": "Work",
    "calendarColor": "#FF6B6B",
    "isRecurring": false
  }
]
```

**Purpose:** Get events for date range with recurring expansion.

### Endpoint: GET /api/events/month

**Query Parameters:**
- `year` - Year (YYYY)
- `month` - Month (1-12)

**Response:**
```json
[1, 5, 10, 15, 20, 25, 30]
```

**Purpose:** Days with events for month calendar marking.

## Security Considerations

1. **Basic Auth Only**
   - No OAuth support
   - Passwords in separate files
   - Files should be chmod 600

2. **No Input Validation**
   - Trusts CalDAV server data
   - No user input (read-only)
   - Safe for private deployments

3. **HTTP Only**
   - No TLS termination
   - Use reverse proxy for HTTPS
   - Production should have SSL

4. **No Authentication**
   - Application has no login
   - Should be behind firewall or VPN
   - Consider adding basic auth in reverse proxy

## Performance Characteristics

- **Event Fetching:** Parallel goroutines for multiple calendars
- **RRULE Expansion:** Limited to 1000 occurrences
- **Timezone Conversion:** Done once during parsing
- **Auto-refresh:** Configurable interval (default 60s)
- **Frontend Bundle:** ~22KB gzipped JS, ~32KB gzipped CSS
- **Docker Image:** ~30MB (Alpine-based)

## Future Enhancements

Potential improvements not implemented in MVP:

1. **Caching Layer**
   - In-memory cache with TTL
   - Redis for distributed caching
   - Reduce CalDAV server load

2. **OAuth Support**
   - OAuth 2.0 flow
   - Token storage
   - Broader server compatibility

3. **Notification Support**
   - Browser notifications
   - Email reminders
   - Webhook integration

4. **Filter/Search**
   - Search events by text
   - Filter by calendar
   - Date range selector

5. **Day/Month Views**
   - Alternative view modes
   - Agenda view
   - Month grid view

6. **Authentication**
   - Basic auth for app
   - User accounts
   - Per-user calendar configurations

7. **Mobile Apps**
   - Progressive Web App (PWA)
   - Native iOS/Android
   - Offline support

8. **Customization**
   - Theme selection
   - Custom colors
   - Layout options

## Testing Notes

The project was tested with:
- **CalDAV Server:** Radicale (tested), also compatible with Nextcloud, Baikal
- **Recurring Events:** Weekly, monthly, yearly, with EXDATE
- **Timezones:** Europe/Rome (UTC+1), should work with any IANA timezone
- **Multiple Calendars:** Tested with 2 calendars
- **Browsers:** Chrome, Firefox, Safari (desktop and mobile)

## Dependencies

### Go Dependencies
```
github.com/emersion/go-ical v0.2.0
github.com/emersion/go-webdav v0.6.0
github.com/teambition/rrule-go v1.8.2
gopkg.in/yaml.v3 v3.0.1
```

### Frontend Dependencies
```
svelte: ^5.46.3
vite: ^7.3.1
bootstrap: ^5.3.8
date-fns: ^4.1.0
typescript: ^5.7.3
```

## Build Process

1. **Frontend Build** (Vite)
   - TypeScript compilation
   - Svelte component compilation
   - CSS processing
   - Asset optimization
   - Output to `web/dist/`

2. **Backend Build** (Go)
   - Embed frontend files
   - Compile Go code
   - Inject version via ldflags
   - Output single binary

3. **Docker Build** (Multi-stage)
   - Stage 1: Node build (frontend)
   - Stage 2: Go build (backend)
   - Stage 3: Alpine runtime (final)

## Lessons Learned

1. **Timezone Handling is Hard**
   - Always use local date parts for comparisons
   - Be careful with `toISOString()` conversions
   - Test with timezones ahead and behind UTC

2. **iCalendar is Complex**
   - Many edge cases and escape sequences
   - RFC 5545 compliance is important
   - Different servers have different quirks

3. **Responsive Design Requires Testing**
   - Test on real devices, not just browser resize
   - Horizontal overflow is easy to miss
   - Single-column layouts work better for content

4. **User Feedback is Valuable**
   - Initial design assumptions were wrong
   - Users know what they want when they see it
   - Be ready to redesign based on feedback

5. **Build Automation Saves Time**
   - Makefile for consistent commands
   - Colored output improves UX
   - Easy cleanup and updates are important

6. **Svelte 5 Runes are Great**
   - Simpler than stores for small apps
   - Clear reactivity model
   - Good TypeScript integration

## Post-MVP Improvements

After the initial MVP release, several important improvements were made based on deployment experience and user feedback:

### 1. .gitignore Fix for cmd/mucal Directory

**Problem:** The binary name `mucal` in `.gitignore` was matching any path containing "mucal", causing the `cmd/mucal/` directory to be excluded from git.

**Impact:** GitHub Actions build failed with "directory not found" because the main.go file wasn't in the repository.

**Solution:**
```gitignore
# Before
mucal

# After
/mucal
```

Changed from `mucal` to `/mucal` to match only the binary in the root directory.

### 2. Docker Image Tagging Strategy

**Problem:** Initial GitHub Actions workflow generated only 3 tags without version prefix: `latest`, `0.0`, `0.0.1`.

**User Requirement:** Generate 4 semantic version tags with `v` prefix: `latest`, `v0`, `v0.0`, `v0.0.1`.

**Solution:**
```yaml
tags: |
  type=semver,pattern={{raw}}              # v0.0.1
  type=semver,pattern=v{{major}}.{{minor}} # v0.0
  type=semver,pattern=v{{major}}           # v0
  type=raw,value=latest                    # latest
```

This provides flexible Docker image pulling:
- `latest` - Always the newest version
- `v0` - Latest in major version 0
- `v0.0` - Latest in minor version 0.0
- `v0.0.1` - Specific version

### 3. Apache 2.0 License Application

**Implementation:**
- Created `LICENSE` file with full Apache 2.0 text
- Added copyright headers to all source files (Go and TypeScript)
- Updated README with license information

**Header Format:**
```go
// Copyright 2026 Mano
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
```

### 4. Smart Past Day Hiding

**User Feedback:** "Usually I don't need to see the events of past days, and if it's a Saturday I have to scroll down to see 'today'."

**Problem:** When viewing the current week on Saturday, user had to scroll past 5 days of past events to see today's events.

**Solution:** Implemented context-aware day visibility:

**For Current Week:**
- Past days are hidden by default
- "Show N past days" button appears at the top
- Clicking the button expands past days
- Focus is automatically on today and future events

**For Past/Future Weeks:**
- All days are always visible
- No button shown
- User wants to see complete historical or future weeks

**Implementation Details:**
```typescript
const isCurrentWeek = $derived(() => {
  const today = new Date();
  today.setHours(0, 0, 0, 0);
  const weekStart = new Date(calendarStore.selectedWeekStart);
  weekStart.setHours(0, 0, 0, 0);
  const weekEnd = new Date(calendarStore.weekEnd);
  weekEnd.setHours(0, 0, 0, 0);
  return today >= weekStart && today <= weekEnd;
});
```

**User Experience:**
- On Monday: 0 past days hidden, see full week
- On Wednesday: 2 past days hidden, "Show 2 past days" button appears
- On Saturday: 5 past days hidden, immediately see today and Sunday
- Navigate to last week: All 7 days visible, no hiding
- Navigate to next week: All 7 days visible, no hiding

This dramatically improves usability for the most common use case (checking today's and upcoming events) while preserving full visibility for historical review.

## Version History

- **v0.1.0** (MVP) - Initial release with core features
  - Week view with recurring events
  - Month calendar for navigation
  - Multiple CalDAV calendars
  - Docker deployment
  - GitHub Actions CI/CD

- **v0.1.1+** (Post-MVP) - Quality of life improvements
  - Fixed .gitignore to include cmd/mucal in git
  - Updated Docker tagging strategy (4 tags with v prefix)
  - Applied Apache 2.0 license to all source files
  - Smart past day hiding in current week
  - Improved UX for daily usage

## Maintenance

To maintain this project:

1. **Update Dependencies**
   ```bash
   make update
   ```

2. **Run Tests**
   ```bash
   make test
   ```

3. **Build and Test Locally**
   ```bash
   make build
   make run
   ```

4. **Create Release**
   ```bash
   git tag v0.2.0
   git push origin v0.2.0
   ```

5. **Check Logs**
   ```bash
   docker logs <container-id>
   ```

## Conclusion

μCal successfully implements a minimal, read-only CalDAV viewer with a clean interface and solid foundation for future enhancements. The project demonstrates effective use of modern web technologies (Svelte 5, Go) and development practices (CI/CD, containerization, build automation).

The key to success was iterative development with user feedback, particularly the major UI redesign that transformed the application from a complex multi-column layout to a simple, scrollable single-column interface.

The project is production-ready for personal or small-team use, with clear paths for enhancement if needed.
