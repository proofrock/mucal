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

// Date utility functions

import {
  startOfWeek,
  endOfWeek,
  eachDayOfInterval,
  format,
  startOfMonth,
  endOfMonth,
  eachWeekOfInterval,
  isSameDay,
  isToday,
  addWeeks,
  subWeeks,
  addMonths,
  subMonths,
  parseISO,
  isWithinInterval,
} from 'date-fns';

// Get the start of the week (Monday) for a given date
export function getWeekStart(date: Date): Date {
  return startOfWeek(date, { weekStartsOn: 1 });
}

// Get the end of the week (Sunday) for a given date
export function getWeekEnd(date: Date): Date {
  return endOfWeek(date, { weekStartsOn: 1 });
}

// Get an array of all days in a week
export function getWeekDays(weekStart: Date): Date[] {
  return eachDayOfInterval({
    start: weekStart,
    end: getWeekEnd(weekStart),
  });
}

// Format date for display
export function formatDate(date: Date, formatStr: string = 'MMM d, yyyy'): string {
  return format(date, formatStr);
}

// Format time for display
export function formatTime(date: Date): string {
  return format(date, 'HH:mm');
}

// Check if two dates are the same day
export function isSameDayAs(date1: Date, date2: Date): boolean {
  return isSameDay(date1, date2);
}

// Check if date is today
export function isDateToday(date: Date): boolean {
  return isToday(date);
}

// Navigate weeks
export function nextWeek(date: Date): Date {
  return addWeeks(date, 1);
}

export function previousWeek(date: Date): Date {
  return subWeeks(date, 1);
}

// Navigate months
export function nextMonth(date: Date): Date {
  return addMonths(date, 1);
}

export function previousMonth(date: Date): Date {
  return subMonths(date, 1);
}

// Get month calendar grid
export function getMonthGrid(date: Date): Date[][] {
  const monthStart = startOfMonth(date);
  const monthEnd = endOfMonth(date);

  const weeks = eachWeekOfInterval(
    { start: monthStart, end: monthEnd },
    { weekStartsOn: 1 }
  );

  return weeks.map(weekStart => getWeekDays(weekStart));
}

// Parse ISO date string to Date
export function parseISODate(dateStr: string): Date {
  return parseISO(dateStr);
}

// Check if event is currently happening
export function isEventCurrent(start: string, end: string): boolean {
  const now = new Date();
  return isWithinInterval(now, {
    start: parseISO(start),
    end: parseISO(end),
  });
}

// Format date range for display
export function formatDateRange(start: Date, end: Date): string {
  return `${formatDate(start, 'MMM d')} - ${formatDate(end, 'MMM d, yyyy')}`;
}
