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

// Type definitions for μCal

export interface Event {
  uid: string;
  summary: string;
  description: string;
  location: string;
  start: string; // ISO 8601 timestamp
  end: string; // ISO 8601 timestamp
  allDay: boolean;
  calendarName: string;
  calendarColor: string;
  isRecurring: boolean;
}

export interface Calendar {
  name: string;
  color: string;
}

export interface Config {
  timezone: string;
  autoRefresh: number; // seconds
  calendars: Calendar[];
}

export interface APIError {
  error: string;
  code: number;
}
