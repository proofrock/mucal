# μCal

A minimal, read-only CalDAV calendar viewer with a responsive web interface.

## Quick Start

```bash
# Clone and setup
git clone https://github.com/mano/mucal.git
cd mucal

# Copy example config
cp config.example.yaml config.yaml
# Edit config.yaml with your CalDAV credentials

# Build and run
make install-deps
make build
make run
```

Visit http://localhost:8080 to see your calendars!

## Features

- **Read-only calendar view** - Display events from one or more CalDAV calendars
- **Week view** - Events grouped by day, displayed vertically for easy scrolling
- **Smart past day hiding** - In current week, past days are hidden by default to focus on today and future (expandable with one click)
- **Collapsible month calendar** - Toggle on-demand to select different weeks
- **Quick navigation** - "Today" button to instantly jump to current week
- **Recurring events** - Full support for RRULE with proper timezone handling
- **Color coding** - Different colors for different calendars
- **Current event highlighting** - Ongoing events are subtly highlighted
- **Responsive design** - Single-column layout works perfectly on desktop and mobile
- **Auto-refresh** - Configurable automatic event refresh
- **Timezone support** - Display events in your configured timezone
- **iCalendar compliance** - Proper handling of escaped characters in event text

## Installation

### Docker (Recommended)

```bash
docker run -d \
  -p 8080:8080 \
  -v $(pwd)/config.yaml:/config/config.yaml:ro \
  -v $(pwd)/secrets:/secrets:ro \
  ghcr.io/mano/mucal:latest
```

### Docker Compose

```yaml
version: '3.8'
services:
  mucal:
    image: ghcr.io/mano/mucal:latest
    ports:
      - "8080:8080"
    volumes:
      - ./config.yaml:/config/config.yaml:ro
      - ./secrets:/secrets:ro
    restart: unless-stopped
```

### Build from Source

Requirements:
- Go 1.23 or higher
- Node.js 20 or higher
- npm
- make

```bash
# Clone the repository
git clone https://github.com/mano/mucal.git
cd mucal

# Install dependencies and build
make install-deps
make build

# Run
make run
```

#### Makefile Targets

The project includes a comprehensive Makefile for building and testing:

```bash
make help              # Show all available targets
make build             # Build complete application
make test              # Run tests
make run               # Build and run locally
make clean             # Remove build artifacts
make cleanup           # Remove all artifacts, caches, and dependencies
make update            # Update all dependencies to latest versions
make docker-build      # Build Docker image
make docker-run        # Run application in Docker
```

For development:
```bash
make dev-frontend      # Run frontend dev server (Vite)
make dev-backend       # Run backend in dev mode
```

## Configuration

Create a `config.yaml` file:

```yaml
# Timezone for displaying events
time_zone: "Europe/Rome"

# Auto-refresh interval in seconds
auto_refresh: 60

# List of CalDAV calendars
calendars:
  - name: "Personal"
    url: "https://calendar.example.com/caldav/personal"
    user_id: "your-username"
    password_file: "/secrets/personal.txt"
    color: "#4ECDC4"

  - name: "Work"
    url: "https://calendar.example.com/caldav/work"
    user_id: "your-username"
    password_file: "/secrets/work.txt"
    color: "#FF6B6B"
```

### Password Files

For security, passwords are stored in separate files (one password per file):

```bash
# Create secrets directory
mkdir secrets

# Create password files
echo -n "your-password" > secrets/personal.txt
echo -n "your-password" > secrets/work.txt

# Secure the files
chmod 600 secrets/*.txt
```

## Usage

### Running the Application

By default, μCal looks for `config.yaml` in the current directory:
```bash
./mucal
```

You can specify a different config file in two ways:
```bash
# Using the -config flag
./mucal -config /path/to/my-config.yaml

# Or as a positional argument
./mucal /path/to/my-config.yaml
```

### Using the Calendar

1. Access the application at `http://localhost:8080`
2. Click "Show Calendar" button to open the month calendar
3. Click any day in the calendar to jump to that week
4. Use Previous/Next buttons or "Today" button to navigate
5. Events automatically refresh based on the configured interval

### Navigation

- **Show/Hide Calendar**: Toggle button at the top to show/hide the month calendar
- **Month calendar**: Click any day to jump to that week (when calendar is visible)
- **Week view**: Use Previous/Next buttons to navigate between weeks
- **Today button**: Instantly jump to the current week
- **Past days in current week**: Hidden by default to focus on today and future; click "Show N past days" button to expand
- **Past/future weeks**: All days are always visible when viewing non-current weeks
- **Single scroll**: Scroll through all days and events together in one continuous view
- **Current day**: Highlighted in blue
- **Current events**: Ongoing events have a subtle background highlight
- **Event colors**: Left border stripe shows the calendar color
- **Event details**: Hover over events to see full text for truncated summaries

## API Endpoints

μCal provides a REST API:

- `GET /api/health` - Health check and version
- `GET /api/config` - Application configuration (sanitized)
- `GET /api/events?start=YYYY-MM-DD&end=YYYY-MM-DD` - Events for date range
- `GET /api/events/month?year=YYYY&month=MM` - Days with events

## Architecture

- **Backend**: Go with embedded frontend
- **Frontend**: Svelte 5 with Bootstrap 5
- **CalDAV**: Direct connection, no database required
- **Port**: Fixed to 8080 (HTTP only)

## Development

### Run Development Server

Using Makefile (recommended):
```bash
# Run backend in dev mode
make dev-backend

# Run frontend dev server (in another terminal)
make dev-frontend
```

Or manually:
```bash
# Backend
go run ./cmd/mucal -config config.yaml

# Frontend (in another terminal)
cd web && npm run dev
```

### Build Docker Image

```bash
make docker-build

# Or manually
docker build -t mucal:latest .
```

### Update Dependencies

Update all dependencies to their latest versions:
```bash
make update
```

This will update both Go modules and npm packages.

### Cleanup

Remove build artifacts:
```bash
make clean
```

Complete cleanup (including dependencies and caches):
```bash
make cleanup
```

## CI/CD

Releases are automated via GitHub Actions:

1. Create and push a semantic version tag:
   ```bash
   git tag v0.1.0
   git push origin v0.1.0
   ```

2. GitHub Actions will:
   - Build the Docker image
   - Push to ghcr.io with version tag and `latest`
   - Version is injected into the binary and displayed in the UI

## Requirements

- CalDAV server (local or remote)
- Basic authentication support
- Docker (for containerized deployment)

## Limitations

- Read-only (no editing, adding, or deleting events)
- No reminders or notifications
- Basic authentication only (no OAuth)
- HTTP only (use a reverse proxy for HTTPS)
- No persistent caching (fetches from CalDAV on each request)

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

Copyright 2026 Mano

## Troubleshooting

### Events not showing

- Check CalDAV server URL and credentials
- Verify password files contain only the password (no newlines)
- Check logs with `docker logs <container-id>`

### Connection errors

- Ensure CalDAV server is accessible from the container
- Check firewall settings
- Verify URL format includes full path to calendar

### Timezone issues

- Use IANA timezone names (e.g., "Europe/Rome", "America/New_York")
- Check available timezones: `docker run --rm alpine cat /usr/share/zoneinfo/zone.tab`

## Version

To check the version, visit `http://localhost:8080/api/health` or look at the header in the UI.
