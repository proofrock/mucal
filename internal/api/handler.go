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

package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/mano/mucal/internal/caldav"
	"github.com/mano/mucal/internal/config"
	"github.com/mano/mucal/internal/version"
)

// Handler handles HTTP requests for the API
type Handler struct {
	config   *config.Config
	clients  []*caldav.Client
	timezone *time.Location
	version  string
}

// NewHandler creates a new API handler
func NewHandler(cfg *config.Config) (*Handler, error) {
	// Get timezone
	tz, err := cfg.GetLocation()
	if err != nil {
		return nil, fmt.Errorf("failed to load timezone: %w", err)
	}

	// Create CalDAV clients for each calendar
	var clients []*caldav.Client
	for i := range cfg.Calendars {
		client, err := caldav.NewClient(&cfg.Calendars[i], tz)
		if err != nil {
			return nil, fmt.Errorf("failed to create client for calendar %s: %w", cfg.Calendars[i].Name, err)
		}
		clients = append(clients, client)
	}

	return &Handler{
		config:   cfg,
		clients:  clients,
		timezone: tz,
		version:  version.Version,
	}, nil
}

// Health handles the health check endpoint
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":  "ok",
		"version": h.version,
	}
	writeJSON(w, http.StatusOK, response)
}

// GetConfig handles the config endpoint (sanitized, no credentials)
func (h *Handler) GetConfig(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.config.Sanitize())
}

// GetEvents handles the events endpoint
func (h *Handler) GetEvents(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")

	if startStr == "" || endStr == "" {
		writeError(w, http.StatusBadRequest, "start and end query parameters are required (format: YYYY-MM-DD)")
		return
	}

	// Parse dates
	start, err := time.ParseInLocation("2006-01-02", startStr, h.timezone)
	if err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("invalid start date format: %v", err))
		return
	}

	end, err := time.ParseInLocation("2006-01-02", endStr, h.timezone)
	if err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("invalid end date format: %v", err))
		return
	}

	// Add one day to end to make it inclusive
	end = end.Add(24 * time.Hour)

	// Fetch events from all calendars in parallel
	var (
		allEvents []*caldav.Event
		mu        sync.Mutex
		wg        sync.WaitGroup
		errs      []error
	)

	for _, client := range h.clients {
		wg.Add(1)
		go func(c *caldav.Client) {
			defer wg.Done()

			events, err := c.FetchEvents(start, end)
			if err != nil {
				mu.Lock()
				errs = append(errs, fmt.Errorf("calendar %s: %w", c.GetCalendarName(), err))
				mu.Unlock()
				return
			}

			mu.Lock()
			allEvents = append(allEvents, events...)
			mu.Unlock()
		}(client)
	}

	wg.Wait()

	// If all calendars failed, return error
	if len(errs) > 0 && len(allEvents) == 0 {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("failed to fetch events: %v", errs))
		return
	}

	// If some calendars failed, log but continue
	if len(errs) > 0 {
		for _, err := range errs {
			fmt.Fprintf(os.Stderr, "Error fetching events: %v\n", err)
		}
	}

	// Sort events
	sort.Sort(caldav.Events(allEvents))

	response := map[string]interface{}{
		"events": allEvents,
	}
	writeJSON(w, http.StatusOK, response)
}

// GetEventsMonth handles the month events endpoint
// Returns days that have events in the specified month
func (h *Handler) GetEventsMonth(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	yearStr := r.URL.Query().Get("year")
	monthStr := r.URL.Query().Get("month")

	if yearStr == "" || monthStr == "" {
		writeError(w, http.StatusBadRequest, "year and month query parameters are required")
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid year format")
		return
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid month format")
		return
	}

	if month < 1 || month > 12 {
		writeError(w, http.StatusBadRequest, "month must be between 1 and 12")
		return
	}

	// Calculate start and end of month
	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, h.timezone)
	end := start.AddDate(0, 1, 0) // First day of next month

	// Fetch events from all calendars in parallel
	var (
		allEvents []*caldav.Event
		mu        sync.Mutex
		wg        sync.WaitGroup
	)

	for _, client := range h.clients {
		wg.Add(1)
		go func(c *caldav.Client) {
			defer wg.Done()

			events, err := c.FetchEvents(start, end)
			if err != nil {
				// Log but continue
				fmt.Fprintf(os.Stderr, "Error fetching events for month view: %v\n", err)
				return
			}

			mu.Lock()
			allEvents = append(allEvents, events...)
			mu.Unlock()
		}(client)
	}

	wg.Wait()

	// Extract unique days
	daysSet := make(map[int]bool)
	for _, event := range allEvents {
		// Get the day of month for the event start
		day := event.Start.In(h.timezone).Day()
		daysSet[day] = true

		// If event spans multiple days, mark all days
		if !event.AllDay {
			eventEnd := event.End.In(h.timezone)
			for d := event.Start.In(h.timezone); d.Before(eventEnd) && d.Month() == time.Month(month); d = d.Add(24 * time.Hour) {
				daysSet[d.Day()] = true
			}
		}
	}

	// Convert to sorted slice
	days := make([]int, 0, len(daysSet))
	for day := range daysSet {
		days = append(days, day)
	}

	response := map[string]interface{}{
		"days": days,
	}
	writeJSON(w, http.StatusOK, response)
}

// writeJSON writes a JSON response
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding JSON response: %v\n", err)
	}
}

// writeError writes a JSON error response
func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	response := map[string]interface{}{
		"error": message,
		"code":  status,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding error response: %v\n", err)
	}
}
