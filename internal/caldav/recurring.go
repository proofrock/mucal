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
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/emersion/go-ical"
	"github.com/teambition/rrule-go"
)

// expandRecurringEvent expands a recurring event based on its RRULE
func (c *Client) expandRecurringEvent(comp *ical.Component, uid, summary, description, location string,
	startTime, endTime time.Time, allDay bool, queryStart, queryEnd time.Time) ([]*Event, error) {

	rruleProp := comp.Props.Get("RRULE")
	if rruleProp == nil {
		return nil, fmt.Errorf("no RRULE found for recurring event")
	}

	// Parse RRULE
	// For all-day events, use DATE format; for timed events, use DATETIME format in UTC
	var dtstart string
	if allDay {
		// Use DATE format (YYYYMMDD) for all-day events to avoid timezone shifts
		dtstart = startTime.Format("20060102")
	} else {
		// Use DATETIME format in UTC for timed events
		dtstart = startTime.UTC().Format("20060102T150405Z")
	}

	rruleStr := "DTSTART:" + dtstart + "\nRRULE:" + rruleProp.Value

	rOption, err := rrule.StrToROption(rruleStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RRULE: %w", err)
	}

	// Create RRule from options
	rule, err := rrule.NewRRule(*rOption)
	if err != nil {
		return nil, fmt.Errorf("failed to create RRule: %w", err)
	}

	// Create RRule set
	rset := &rrule.Set{}
	rset.RRule(rule)

	// Handle EXDATE (excluded dates)
	// Try to get EXDATE properties
	exdateProp := comp.Props.Get("EXDATE")
	if exdateProp != nil {
		exdates := strings.Split(exdateProp.Value, ",")
		for _, exdateStr := range exdates {
			exdate, err := parseICalDate(exdateStr, c.timezone)
			if err != nil {
				// Log but continue
				fmt.Fprintf(os.Stderr, "Failed to parse EXDATE: %v\n", err)
				continue
			}
			rset.ExDate(exdate)
		}
	}

	// Generate occurrences within the query range (with some buffer)
	// Add buffer to catch events that might span across the boundary
	bufferStart := queryStart.Add(-30 * 24 * time.Hour)
	bufferEnd := queryEnd.Add(30 * 24 * time.Hour)

	occurrences := rset.Between(bufferStart, bufferEnd, true)

	// Limit occurrences to prevent infinite expansion
	const maxOccurrences = 1000
	if len(occurrences) > maxOccurrences {
		occurrences = occurrences[:maxOccurrences]
	}

	// Calculate event duration
	duration := endTime.Sub(startTime)

	var events []*Event

	// Create event for each occurrence
	for _, occurrence := range occurrences {
		// Convert to configured timezone
		occStart := occurrence.In(c.timezone)
		occEnd := occStart.Add(duration)

		// Filter to only include events that overlap with query range
		if occEnd.Before(queryStart) || occStart.After(queryEnd) {
			continue
		}

		event := &Event{
			UID:          uid + "_" + occStart.Format("20060102T150405"),
			Summary:      summary,
			Description:  description,
			Location:     location,
			Start:        occStart,
			End:          occEnd,
			AllDay:       allDay,
			CalendarName: c.calendar.Name,
			CalendarColor: c.calendar.Color,
			IsRecurring:  true,
		}

		events = append(events, event)
	}

	return events, nil
}

// parseDuration parses an iCalendar DURATION value
// Format: P[n]W[n]D[T[n]H[n]M[n]S]
func parseDuration(value string) (time.Duration, error) {
	if !strings.HasPrefix(value, "P") && !strings.HasPrefix(value, "-P") {
		return 0, fmt.Errorf("invalid duration format: %s", value)
	}

	negative := strings.HasPrefix(value, "-")
	if negative {
		value = value[2:] // Remove "-P"
	} else {
		value = value[1:] // Remove "P"
	}

	var duration time.Duration

	// Split by T for date and time parts
	parts := strings.Split(value, "T")
	datePart := parts[0]
	timePart := ""
	if len(parts) > 1 {
		timePart = parts[1]
	}

	// Parse date part
	if strings.Contains(datePart, "W") {
		// Weeks
		weeks, err := strconv.Atoi(strings.TrimSuffix(datePart, "W"))
		if err != nil {
			return 0, fmt.Errorf("invalid weeks in duration: %s", datePart)
		}
		duration += time.Duration(weeks) * 7 * 24 * time.Hour
	} else {
		// Days
		if strings.Contains(datePart, "D") {
			days, err := strconv.Atoi(strings.TrimSuffix(datePart, "D"))
			if err != nil {
				return 0, fmt.Errorf("invalid days in duration: %s", datePart)
			}
			duration += time.Duration(days) * 24 * time.Hour
		}
	}

	// Parse time part
	if timePart != "" {
		remaining := timePart

		// Hours
		if strings.Contains(remaining, "H") {
			idx := strings.Index(remaining, "H")
			hours, err := strconv.Atoi(remaining[:idx])
			if err != nil {
				return 0, fmt.Errorf("invalid hours in duration: %s", remaining[:idx])
			}
			duration += time.Duration(hours) * time.Hour
			remaining = remaining[idx+1:]
		}

		// Minutes
		if strings.Contains(remaining, "M") {
			idx := strings.Index(remaining, "M")
			minutes, err := strconv.Atoi(remaining[:idx])
			if err != nil {
				return 0, fmt.Errorf("invalid minutes in duration: %s", remaining[:idx])
			}
			duration += time.Duration(minutes) * time.Minute
			remaining = remaining[idx+1:]
		}

		// Seconds
		if strings.Contains(remaining, "S") {
			idx := strings.Index(remaining, "S")
			seconds, err := strconv.Atoi(remaining[:idx])
			if err != nil {
				return 0, fmt.Errorf("invalid seconds in duration: %s", remaining[:idx])
			}
			duration += time.Duration(seconds) * time.Second
		}
	}

	if negative {
		duration = -duration
	}

	return duration, nil
}

// parseICalDate parses an iCalendar date/datetime string
func parseICalDate(value string, tz *time.Location) (time.Time, error) {
	// Try different formats
	formats := []string{
		"20060102T150405Z",
		"20060102T150405",
		"20060102",
	}

	for _, format := range formats {
		t, err := time.ParseInLocation(format, value, tz)
		if err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", value)
}
