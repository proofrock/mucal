<script lang="ts">
  import { calendarStore } from '../stores/calendar.svelte';
  import { getMonthGrid, formatDate, isSameDayAs, isDateToday, previousMonth, nextMonth } from '../utils/date';
  import { onMount } from 'svelte';

  let currentMonth = $state(new Date());

  const monthGrid = $derived(getMonthGrid(currentMonth));
  const monthName = $derived(formatDate(currentMonth, 'MMMM yyyy'));

  // Load event days when month changes
  $effect(() => {
    const year = currentMonth.getFullYear();
    const month = currentMonth.getMonth() + 1;
    calendarStore.loadMonthEventDays(year, month);
  });

  function handleDayClick(date: Date) {
    calendarStore.selectWeek(date);
  }

  function goToPreviousMonth() {
    currentMonth = previousMonth(currentMonth);
  }

  function goToNextMonth() {
    currentMonth = nextMonth(currentMonth);
  }

  function hasEvents(date: Date): boolean {
    const day = date.getDate();
    return calendarStore.monthEventDays.includes(day);
  }

  function isInCurrentMonth(date: Date): boolean {
    return date.getMonth() === currentMonth.getMonth();
  }
</script>

<div class="month-calendar">
  <div class="month-header">
    <button class="btn btn-sm btn-outline-secondary" aria-label="Previous month" onclick={goToPreviousMonth}>
      <i class="bi bi-chevron-left"></i>
    </button>
    <h6 class="mb-0">{monthName}</h6>
    <button class="btn btn-sm btn-outline-secondary" aria-label="Next month" onclick={goToNextMonth}>
      <i class="bi bi-chevron-right"></i>
    </button>
  </div>

  <table class="table table-sm table-bordered">
    <thead>
      <tr>
        <th>Mon</th>
        <th>Tue</th>
        <th>Wed</th>
        <th>Thu</th>
        <th>Fri</th>
        <th>Sat</th>
        <th>Sun</th>
      </tr>
    </thead>
    <tbody>
      {#each monthGrid as week}
        <tr>
          {#each week as day}
            <td
              class="day-cell"
              class:other-month={!isInCurrentMonth(day)}
              class:today={isDateToday(day)}
              class:has-events={hasEvents(day) && isInCurrentMonth(day)}
              onclick={() => handleDayClick(day)}
            >
              <div class="day-number">{day.getDate()}</div>
              {#if hasEvents(day) && isInCurrentMonth(day)}
                <div class="event-dot"></div>
              {/if}
            </td>
          {/each}
        </tr>
      {/each}
    </tbody>
  </table>
</div>

<style>
  .month-calendar {
    background: white;
    border-radius: 8px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    padding: 1rem;
    width: 100%;
  }

  .month-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1rem;
  }

  .table {
    margin-bottom: 0;
    width: 100%;
    table-layout: fixed;
  }

  .table thead th {
    text-align: center;
    font-size: 0.75rem;
    font-weight: 600;
    padding: 0.5rem 0.25rem;
    background-color: #f8f9fa;
  }

  .day-cell {
    text-align: center;
    padding: 0.5rem;
    cursor: pointer;
    position: relative;
    transition: background-color 0.2s;
    vertical-align: middle;
    height: 50px;
  }

  .day-cell:hover {
    background-color: #f8f9fa;
  }

  .day-cell.other-month {
    color: #adb5bd;
  }

  .day-cell.today {
    background-color: #e7f3ff;
    font-weight: 700;
  }

  .day-cell.has-events .day-number {
    font-weight: 600;
  }

  .day-number {
    font-size: 0.9rem;
  }

  .event-dot {
    width: 6px;
    height: 6px;
    background-color: #0d6efd;
    border-radius: 50%;
    position: absolute;
    bottom: 8px;
    left: 50%;
    transform: translateX(-50%);
  }

  @media (max-width: 768px) {
    .month-calendar {
      padding: 0.75rem;
    }

    .table thead th {
      font-size: 0.7rem;
      padding: 0.4rem 0.2rem;
    }

    .day-cell {
      padding: 0.4rem;
      height: 45px;
    }

    .day-number {
      font-size: 0.8rem;
    }
  }
</style>
