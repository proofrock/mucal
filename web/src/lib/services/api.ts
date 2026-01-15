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

// API service layer for μCal

import type { Config, Event, APIError } from '../types';
import { format } from 'date-fns';

const API_BASE = '/api';

// Fetch wrapper with error handling
async function fetchAPI<T>(url: string): Promise<T> {
  const response = await fetch(url);

  if (!response.ok) {
    const error: APIError = await response.json().catch(() => ({
      error: `HTTP ${response.status}: ${response.statusText}`,
      code: response.status,
    }));
    throw new Error(error.error);
  }

  return response.json();
}

// Fetch application configuration
export async function fetchConfig(): Promise<Config> {
  return fetchAPI<Config>(`${API_BASE}/config`);
}

// Fetch events for a date range
export async function fetchEvents(start: Date, end: Date): Promise<Event[]> {
  const startStr = format(start, 'yyyy-MM-dd');
  const endStr = format(end, 'yyyy-MM-dd');
  const response = await fetchAPI<{ events: Event[] }>(
    `${API_BASE}/events?start=${startStr}&end=${endStr}`
  );
  return response.events;
}

// Fetch days with events for a month
export async function fetchMonthEventDays(
  year: number,
  month: number
): Promise<number[]> {
  const response = await fetchAPI<{ days: number[] }>(
    `${API_BASE}/events/month?year=${year}&month=${month}`
  );
  return response.days;
}

// Fetch health/version
export async function fetchHealth(): Promise<{ status: string; version: string }> {
  return fetchAPI<{ status: string; version: string }>(`${API_BASE}/health`);
}
