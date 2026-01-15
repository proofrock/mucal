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

// Calendar store using Svelte 5 runes

import type { Config, Event } from '../types';
import { fetchConfig, fetchEvents, fetchMonthEventDays, fetchHealth } from '../services/api';
import { getWeekStart, getWeekEnd, getWeekDays } from '../utils/date';
import { errorStore } from './error.svelte';
import { addDays } from 'date-fns';

export class CalendarStore {
  // State
  currentDate = $state(new Date());
  selectedWeekStart = $state(getWeekStart(new Date()));
  events = $state<Event[]>([]);
  monthEventDays = $state<number[]>([]);
  loading = $state(false);
  config = $state<Config | null>(null);
  version = $state<string>('');

  // Derived state
  weekEnd = $derived(getWeekEnd(this.selectedWeekStart));
  weekDays = $derived(getWeekDays(this.selectedWeekStart));

  // Initialize: load config and initial events
  async init() {
    try {
      // Load config
      this.config = await fetchConfig();

      // Load version
      const health = await fetchHealth();
      this.version = health.version;

      // Load initial events
      await this.loadEvents();
    } catch (error) {
      errorStore.showError(
        error instanceof Error ? error.message : 'Failed to initialize'
      );
    }
  }

  // Load events for the selected week
  async loadEvents() {
    if (this.loading) return;

    this.loading = true;
    try {
      this.events = await fetchEvents(this.selectedWeekStart, this.weekEnd);
    } catch (error) {
      errorStore.showError(
        error instanceof Error ? error.message : 'Failed to load events'
      );
    } finally {
      this.loading = false;
    }
  }

  // Load month event days for the month calendar
  async loadMonthEventDays(year: number, month: number) {
    try {
      this.monthEventDays = await fetchMonthEventDays(year, month);
    } catch (error) {
      // Silently fail for month markers
      console.error('Failed to load month event days:', error);
    }
  }

  // Select a new week
  selectWeek(date: Date) {
    this.selectedWeekStart = getWeekStart(date);
    this.loadEvents();
  }

  // Navigate to next week
  nextWeek() {
    this.selectedWeekStart = addDays(this.selectedWeekStart, 7);
    this.loadEvents();
  }

  // Navigate to previous week
  previousWeek() {
    this.selectedWeekStart = addDays(this.selectedWeekStart, -7);
    this.loadEvents();
  }

  // Go to current week
  goToCurrentWeek() {
    this.selectedWeekStart = getWeekStart(new Date());
    this.loadEvents();
  }

  // Get events for a specific day
  getEventsForDay(date: Date): Event[] {
    // Use local date parts to avoid timezone conversion issues
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
}

export const calendarStore = new CalendarStore();
