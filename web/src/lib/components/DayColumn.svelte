<script lang="ts">
  import type { Event } from '../types';
  import EventCard from './EventCard.svelte';
  import { formatDate, isDateToday } from '../utils/date';

  interface Props {
    date: Date;
    events: Event[];
  }

  let { date, events }: Props = $props();

  const isToday = $derived(isDateToday(date));
  const dayName = $derived(formatDate(date, 'EEEE'));
  const dayNumber = $derived(formatDate(date, 'd'));
  const monthName = $derived(formatDate(date, 'MMMM'));

  // Sort events: all-day first, then by start time
  const sortedEvents = $derived(
    [...events].sort((a, b) => {
      if (a.allDay !== b.allDay) return a.allDay ? -1 : 1;
      return new Date(a.start).getTime() - new Date(b.start).getTime();
    })
  );
</script>

<div class="day-section" class:today={isToday}>
  <div class="day-header">
    <div class="day-info">
      <span class="day-name">{dayName}</span>
      <span class="day-date">{dayNumber} {monthName}</span>
    </div>
    {#if sortedEvents.length > 0}
      <span class="event-count">{sortedEvents.length} event{sortedEvents.length !== 1 ? 's' : ''}</span>
    {/if}
  </div>

  <div class="day-events">
    {#if sortedEvents.length === 0}
      <div class="no-events">No events</div>
    {:else}
      {#each sortedEvents as event (event.uid)}
        <EventCard {event} />
      {/each}
    {/if}
  </div>
</div>

<style>
  .day-section {
    margin-bottom: 1.5rem;
    background: #f8f9fa;
    border-radius: 8px;
    overflow: hidden;
  }

  .day-section.today {
    background: #e7f3ff;
    box-shadow: 0 2px 8px rgba(13, 110, 253, 0.2);
  }

  .day-header {
    padding: 0.75rem 1rem;
    background: white;
    border-bottom: 2px solid #dee2e6;
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .day-section.today .day-header {
    background: #0d6efd;
    color: white;
    border-bottom-color: #0d6efd;
  }

  .day-info {
    display: flex;
    flex-direction: column;
  }

  .day-name {
    font-size: 0.85rem;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }

  .day-date {
    font-size: 0.75rem;
    opacity: 0.8;
    margin-top: 0.1rem;
  }

  .event-count {
    font-size: 0.75rem;
    font-weight: 600;
    padding: 0.25rem 0.5rem;
    background: rgba(0, 0, 0, 0.1);
    border-radius: 12px;
  }

  .day-section.today .event-count {
    background: rgba(255, 255, 255, 0.2);
  }

  .day-events {
    padding: 0.75rem;
  }

  .no-events {
    text-align: center;
    color: #6c757d;
    font-size: 0.85rem;
    padding: 1rem;
    font-style: italic;
  }

  @media (max-width: 576px) {
    .day-section {
      margin-bottom: 1rem;
      border-radius: 6px;
    }

    .day-header {
      padding: 0.6rem 0.75rem;
    }

    .day-name {
      font-size: 0.8rem;
    }

    .day-date {
      font-size: 0.7rem;
    }

    .day-events {
      padding: 0.5rem;
    }
  }
</style>
