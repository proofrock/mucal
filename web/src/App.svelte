<script lang="ts">
  import { onMount } from 'svelte';
  import { calendarStore } from './lib/stores/calendar.svelte';
  import Header from './lib/components/Header.svelte';
  import ErrorPopup from './lib/components/ErrorPopup.svelte';
  import MonthCalendar from './lib/components/MonthCalendar.svelte';
  import WeekView from './lib/components/WeekView.svelte';

  let showCalendar = $state(false);

  // Initialize on mount
  onMount(async () => {
    await calendarStore.init();
  });

  // Auto-refresh effect
  $effect(() => {
    if (calendarStore.config && calendarStore.config.autoRefresh > 0) {
      const interval = setInterval(() => {
        calendarStore.loadEvents();
      }, calendarStore.config.autoRefresh * 1000);

      return () => clearInterval(interval);
    }
  });

  function toggleCalendar() {
    showCalendar = !showCalendar;
  }
</script>

<div class="app">
  <Header />

  <div class="container">
    {#if !calendarStore.config}
      <div class="text-center py-5">
        <div class="spinner-border" role="status">
          <span class="visually-hidden">Loading...</span>
        </div>
        <p class="mt-3">Loading configuration...</p>
      </div>
    {:else}
      <div class="app-layout">
        <!-- Calendar toggle button -->
        <div class="calendar-controls">
          <button
            class="btn btn-outline-primary btn-sm"
            onclick={toggleCalendar}
            aria-label="Toggle month calendar"
          >
            <i class="bi {showCalendar ? 'bi-calendar-x' : 'bi-calendar'}"></i>
            {showCalendar ? 'Hide' : 'Show'} Calendar
          </button>
        </div>

        <!-- Month Calendar (collapsible) -->
        {#if showCalendar}
          <div class="calendar-section">
            <MonthCalendar />
          </div>
        {/if}

        <!-- Week events list below -->
        <div class="events-section">
          <WeekView />
        </div>
      </div>
    {/if}
  </div>

  <ErrorPopup />
</div>

<style>
  .app {
    min-height: 100vh;
    background-color: #f0f2f5;
    display: flex;
    flex-direction: column;
  }

  .container {
    max-width: 800px;
    margin: 0 auto;
    padding: 0 1rem 2rem 1rem;
    width: 100%;
    flex: 1;
  }

  .app-layout {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
  }

  .calendar-controls {
    display: flex;
    justify-content: center;
    margin-top: 1rem;
  }

  .calendar-section {
    width: 100%;
  }

  .events-section {
    width: 100%;
  }

  @media (max-width: 768px) {
    .container {
      padding: 0 0.75rem 1.5rem 0.75rem;
    }

    .app-layout {
      gap: 1rem;
    }
  }

  @media (max-width: 576px) {
    .container {
      padding: 0 0.5rem 1rem 0.5rem;
    }
  }
</style>
