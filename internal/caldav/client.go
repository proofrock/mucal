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

package caldav

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/emersion/go-ical"
	"github.com/emersion/go-webdav/caldav"
	"github.com/mano/mucal/internal/config"
)

// Client wraps the CalDAV client with calendar configuration
type Client struct {
	httpClient   *http.Client
	caldavClient *caldav.Client
	calendar     *config.Calendar
	timezone     *time.Location
}

// NewClient creates a new CalDAV client for the given calendar
func NewClient(cal *config.Calendar, tz *time.Location) (*Client, error) {
	password, err := cal.GetPassword()
	if err != nil {
		return nil, fmt.Errorf("failed to get password for calendar %s: %w", cal.Name, err)
	}

	// Create HTTP client with Basic Auth
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &basicAuthTransport{
			Username: cal.UserID,
			Password: password,
		},
	}

	// Create CalDAV client
	caldavClient, err := caldav.NewClient(httpClient, cal.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to create CalDAV client for %s: %w", cal.Name, err)
	}

	return &Client{
		httpClient:   httpClient,
		caldavClient: caldavClient,
		calendar:     cal,
		timezone:     tz,
	}, nil
}

// FetchEvents fetches calendar events within the given time range
func (c *Client) FetchEvents(start, end time.Time) ([]*Event, error) {
	// Query for calendar objects within the date range
	query := &caldav.CalendarQuery{
		CompRequest: caldav.CalendarCompRequest{
			Name:  "VCALENDAR",
			Props: []string{"VERSION"},
			Comps: []caldav.CalendarCompRequest{
				{
					Name: "VEVENT",
					Props: []string{
						"UID",
						"SUMMARY",
						"DESCRIPTION",
						"LOCATION",
						"DTSTART",
						"DTEND",
						"DURATION",
						"RRULE",
						"EXDATE",
						"RECURRENCE-ID",
					},
				},
			},
		},
		CompFilter: caldav.CompFilter{
			Name: "VCALENDAR",
			Comps: []caldav.CompFilter{
				{
					Name:  "VEVENT",
					Start: start,
					End:   end,
				},
			},
		},
	}

	// Fetch calendar objects
	ctx := context.Background()
	objects, err := c.caldavClient.QueryCalendar(ctx, "", query)
	if err != nil {
		return nil, fmt.Errorf("failed to query calendar %s: %w", c.calendar.Name, err)
	}

	// Parse events from calendar objects
	var events []*Event
	for _, obj := range objects {
		parsedEvents, err := c.parseCalendarObject(&obj, start, end)
		if err != nil {
			// Log error but continue processing other events
			fmt.Fprintf(os.Stderr, "Error parsing calendar object in %s: %v\n", c.calendar.Name, err)
			continue
		}
		events = append(events, parsedEvents...)
	}

	// Sort events
	sort.Sort(Events(events))

	return events, nil
}

// parseCalendarObject parses a CalDAV calendar object into events
func (c *Client) parseCalendarObject(obj *caldav.CalendarObject, queryStart, queryEnd time.Time) ([]*Event, error) {
	// obj.Data is already an *ical.Calendar
	cal := obj.Data
	if cal == nil {
		return nil, fmt.Errorf("calendar object has no data")
	}

	var events []*Event

	// Process each VEVENT
	for _, comp := range cal.Children {
		if comp.Name != "VEVENT" {
			continue
		}

		event, err := c.parseEvent(comp, queryStart, queryEnd)
		if err != nil {
			// Log error but continue
			fmt.Fprintf(os.Stderr, "Error parsing event in %s: %v\n", c.calendar.Name, err)
			continue
		}

		if event != nil {
			events = append(events, event...)
		}
	}

	return events, nil
}

// parseEvent parses a single VEVENT component
func (c *Client) parseEvent(comp *ical.Component, queryStart, queryEnd time.Time) ([]*Event, error) {
	// Extract basic properties
	uid := comp.Props.Get("UID")
	if uid == nil {
		return nil, fmt.Errorf("event missing UID")
	}

	summary := ""
	if prop := comp.Props.Get("SUMMARY"); prop != nil {
		summary = unescapeICalText(prop.Value)
	}

	description := ""
	if prop := comp.Props.Get("DESCRIPTION"); prop != nil {
		description = unescapeICalText(prop.Value)
	}

	location := ""
	if prop := comp.Props.Get("LOCATION"); prop != nil {
		location = unescapeICalText(prop.Value)
	}

	// Parse start time
	dtstart := comp.Props.Get("DTSTART")
	if dtstart == nil {
		return nil, fmt.Errorf("event missing DTSTART")
	}

	startTime, allDay, err := c.parseDateTime(dtstart)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DTSTART: %w", err)
	}

	// Parse end time
	var endTime time.Time
	dtend := comp.Props.Get("DTEND")
	duration := comp.Props.Get("DURATION")

	if dtend != nil {
		endTime, _, err = c.parseDateTime(dtend)
		if err != nil {
			return nil, fmt.Errorf("failed to parse DTEND: %w", err)
		}
	} else if duration != nil {
		// Parse duration
		dur, err := parseDuration(duration.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to parse DURATION: %w", err)
		}
		endTime = startTime.Add(dur)
	} else {
		// Default: all-day events are 1 day, timed events are 0 duration
		if allDay {
			endTime = startTime.Add(24 * time.Hour)
		} else {
			endTime = startTime
		}
	}

	// Check if event has recurrence rule
	rrule := comp.Props.Get("RRULE")
	if rrule != nil {
		// Recurring event - expand it
		return c.expandRecurringEvent(comp, uid.Value, summary, description, location,
			startTime, endTime, allDay, queryStart, queryEnd)
	}

	// Single event
	event := &Event{
		UID:          uid.Value,
		Summary:      summary,
		Description:  description,
		Location:     location,
		Start:        startTime,
		End:          endTime,
		AllDay:       allDay,
		CalendarName: c.calendar.Name,
		CalendarColor: c.calendar.Color,
		IsRecurring:  false,
	}

	return []*Event{event}, nil
}

// parseDateTime parses an iCalendar date/time property
func (c *Client) parseDateTime(prop *ical.Prop) (time.Time, bool, error) {
	// Check if it's a DATE (all-day) or DATE-TIME
	valueType := prop.Params.Get("VALUE")
	isDate := valueType == "DATE"

	var t time.Time
	var err error

	if isDate {
		// Parse as DATE (YYYYMMDD)
		t, err = time.ParseInLocation("20060102", prop.Value, c.timezone)
		if err != nil {
			return time.Time{}, false, err
		}
		return t, true, nil
	}

	// Parse as DATE-TIME
	// Check for TZID parameter
	tzid := prop.Params.Get("TZID")
	if tzid != "" {
		// Parse with specific timezone
		loc, err := time.LoadLocation(tzid)
		if err != nil {
			// Fallback to configured timezone
			loc = c.timezone
		}

		// Parse timestamp (YYYYMMDDTHHMMSS or YYYYMMDDTHHMMSSZ)
		t, err = time.ParseInLocation("20060102T150405", prop.Value, loc)
		if err != nil {
			// Try with Z suffix
			t, err = time.ParseInLocation("20060102T150405Z", prop.Value, time.UTC)
			if err != nil {
				return time.Time{}, false, err
			}
		}
	} else {
		// No timezone specified, try UTC or local
		value := prop.Value
		if len(value) > 0 && value[len(value)-1] == 'Z' {
			t, err = time.Parse("20060102T150405Z", value)
		} else {
			t, err = time.ParseInLocation("20060102T150405", value, c.timezone)
		}
		if err != nil {
			return time.Time{}, false, err
		}
	}

	// Convert to configured timezone
	t = t.In(c.timezone)

	return t, false, nil
}

// GetCalendarName returns the name of the calendar
func (c *Client) GetCalendarName() string {
	return c.calendar.Name
}

// unescapeICalText unescapes iCalendar TEXT values according to RFC 5545
// Handles: \, -> comma, \; -> semicolon, \n or \N -> newline, \\ -> backslash
func unescapeICalText(s string) string {
	var result strings.Builder
	result.Grow(len(s))

	for i := 0; i < len(s); i++ {
		if s[i] == '\\' && i+1 < len(s) {
			// Handle escape sequences
			next := s[i+1]
			switch next {
			case ',', ';', '\\':
				result.WriteByte(next)
				i++ // Skip the next character
			case 'n', 'N':
				result.WriteByte('\n')
				i++ // Skip the next character
			default:
				// Unknown escape sequence, keep the backslash
				result.WriteByte('\\')
			}
		} else {
			result.WriteByte(s[i])
		}
	}

	return result.String()
}

// basicAuthTransport is an http.RoundTripper that adds Basic Authentication
type basicAuthTransport struct {
	Username string
	Password string
}

func (t *basicAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.SetBasicAuth(t.Username, t.Password)
	return http.DefaultTransport.RoundTrip(req)
}
