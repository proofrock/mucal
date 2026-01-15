<script lang="ts">
  import { calendarStore } from '../stores/calendar.svelte';
  import DayColumn from './DayColumn.svelte';
  import { formatDateRange, isDateToday } from '../utils/date';

  let showPastDays = $state(false);

  const weekRange = $derived(
    formatDateRange(calendarStore.selectedWeekStart, calendarStore.weekEnd)
  );

  // Check if we're viewing the current week
  const isCurrentWeek = $derived(() => {
    const today = new Date();
    today.setHours(0, 0, 0, 0);
    const weekStart = new Date(calendarStore.selectedWeekStart);
    weekStart.setHours(0, 0, 0, 0);
    const weekEnd = new Date(calendarStore.weekEnd);
    weekEnd.setHours(0, 0, 0, 0);
    return today >= weekStart && today <= weekEnd;
  });

  // Split days into past and today/future (only for current week)
  const pastDays = $derived(
    calendarStore.weekDays.filter(day => {
      const today = new Date();
      today.setHours(0, 0, 0, 0);
      const dayDate = new Date(day);
      dayDate.setHours(0, 0, 0, 0);
      return dayDate < today;
    })
  );

  const todayAndFutureDays = $derived(
    calendarStore.weekDays.filter(day => {
      const today = new Date();
      today.setHours(0, 0, 0, 0);
      const dayDate = new Date(day);
      dayDate.setHours(0, 0, 0, 0);
      return dayDate >= today;
    })
  );

  // Reset showPastDays when week changes
  $effect(() => {
    calendarStore.selectedWeekStart;
    showPastDays = false;
  });
</script>

<div class="week-view">
  <div class="week-header">
    <button
      class="btn btn-sm btn-primary"
      aria-label="Previous week"
      onclick={() => calendarStore.previousWeek()}
    >
      <i class="bi bi-chevron-left"></i>
      <span>Previous</span>
    </button>
    <div class="week-header-center">
      <h5 class="week-range">{weekRange}</h5>
      <button
        class="btn btn-sm btn-outline-light"
        aria-label="Go to current week"
        onclick={() => calendarStore.goToCurrentWeek()}
      >
        <i class="bi bi-calendar-check"></i>
        <span>Today</span>
      </button>
    </div>
    <button
      class="btn btn-sm btn-primary"
      aria-label="Next week"
      onclick={() => calendarStore.nextWeek()}
    >
      <span>Next</span>
      <i class="bi bi-chevron-right"></i>
    </button>
  </div>

  {#if calendarStore.loading}
    <div class="loading-container">
      <div class="spinner-border text-primary" role="status">
        <span class="visually-hidden">Loading...</span>
      </div>
      <p class="loading-text">Loading events...</p>
    </div>
  {:else}
    <div class="days-list">
      {#if isCurrentWeek()}
        <!-- Current week: hide past days by default -->
        {#if !showPastDays && pastDays.length > 0}
          <div class="show-past-button-container">
            <button
              class="btn btn-sm btn-outline-secondary"
              onclick={() => showPastDays = true}
              aria-label="Show past days"
            >
              <i class="bi bi-chevron-down"></i>
              Show {pastDays.length} past {pastDays.length === 1 ? 'day' : 'days'}
            </button>
          </div>
        {/if}

        {#if showPastDays}
          {#each pastDays as day}
            <DayColumn date={day} events={calendarStore.getEventsForDay(day)} />
          {/each}
        {/if}

        {#each todayAndFutureDays as day}
          <DayColumn date={day} events={calendarStore.getEventsForDay(day)} />
        {/each}
      {:else}
        <!-- Past or future week: show all days -->
        {#each calendarStore.weekDays as day}
          <DayColumn date={day} events={calendarStore.getEventsForDay(day)} />
        {/each}
      {/if}
    </div>
  {/if}
</div>

<style>
  .week-view {
    background: white;
    border-radius: 8px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    overflow: hidden;
  }

  .week-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1rem;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    gap: 0.75rem;
  }

  .week-header-center {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.5rem;
    flex: 1;
  }

  .week-range {
    margin: 0;
    font-size: 1.1rem;
    font-weight: 600;
    text-align: center;
  }

  .week-header button {
    white-space: nowrap;
    display: flex;
    align-items: center;
    gap: 0.25rem;
  }

  .loading-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 3rem 1rem;
    gap: 1rem;
  }

  .loading-text {
    margin: 0;
    color: #6c757d;
    font-size: 0.9rem;
  }

  .days-list {
    padding: 1rem;
  }

  .show-past-button-container {
    display: flex;
    justify-content: center;
    padding: 1rem 0;
    margin-bottom: 0.5rem;
  }

  .show-past-button-container button {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  @media (max-width: 768px) {
    .week-header {
      padding: 0.75rem;
      gap: 0.5rem;
    }

    .week-range {
      font-size: 0.95rem;
    }

    .week-header button {
      font-size: 0.85rem;
      padding: 0.375rem 0.5rem;
    }

    .days-list {
      padding: 0.75rem;
    }
  }

  @media (max-width: 576px) {
    .week-header button span {
      display: none;
    }

    .week-header button {
      padding: 0.5rem;
    }

    .week-range {
      font-size: 0.85rem;
    }
  }
</style>
