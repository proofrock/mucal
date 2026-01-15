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
	"time"
)

// Event represents a calendar event
type Event struct {
	UID          string    `json:"uid"`
	Summary      string    `json:"summary"`
	Description  string    `json:"description"`
	Location     string    `json:"location"`
	Start        time.Time `json:"start"`
	End          time.Time `json:"end"`
	AllDay       bool      `json:"allDay"`
	CalendarName string    `json:"calendarName"`
	CalendarColor string   `json:"calendarColor"`
	IsRecurring  bool      `json:"isRecurring"`
}

// Events is a slice of Event pointers with sorting capabilities
type Events []*Event

func (e Events) Len() int      { return len(e) }
func (e Events) Swap(i, j int) { e[i], e[j] = e[j], e[i] }

// Less implements sort.Interface for Events
// All-day events come first, then sorted by start time
func (e Events) Less(i, j int) bool {
	// All-day events come first
	if e[i].AllDay != e[j].AllDay {
		return e[i].AllDay
	}

	// Then by start time
	if !e[i].Start.Equal(e[j].Start) {
		return e[i].Start.Before(e[j].Start)
	}

	// Finally by summary (alphabetically)
	return e[i].Summary < e[j].Summary
}
