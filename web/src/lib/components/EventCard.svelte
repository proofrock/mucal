<script lang="ts">
  import type { Event } from '../types';
  import { formatTime, parseISODate, isEventCurrent } from '../utils/date';

  interface Props {
    event: Event;
  }

  let { event }: Props = $props();

  const startTime = $derived(parseISODate(event.start));
  const endTime = $derived(parseISODate(event.end));
  const isCurrent = $derived(isEventCurrent(event.start, event.end));

  const timeDisplay = $derived(
    event.allDay
      ? 'All day'
      : `${formatTime(startTime)} - ${formatTime(endTime)}`
  );
</script>

<div
  class="event-item"
  class:current-event={isCurrent && !event.allDay}
  style="border-left: 4px solid {event.calendarColor}"
  title={event.summary + (event.location ? '\n' + event.location : '') + (event.description ? '\n' + event.description : '')}
>
  <div class="event-meta">
    <span class="event-time">{timeDisplay}</span>
    <span class="event-calendar" style="color: {event.calendarColor}">
      {event.calendarName}
    </span>
  </div>
  <div class="event-summary">
    {event.summary || '(No title)'}
    {#if event.location}
      <span class="event-location">📍 {event.location}</span>
    {/if}
  </div>
</div>

<style>
  .event-item {
    background: white;
    padding: 0.5rem 0.75rem;
    border-radius: 4px;
    margin-bottom: 0.5rem;
    cursor: pointer;
    transition: all 0.2s;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  }

  .event-item:hover {
    box-shadow: 0 2px 6px rgba(0, 0, 0, 0.15);
    transform: translateX(2px);
  }

  .current-event {
    background-color: #fffbf0;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.15);
  }

  .event-meta {
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-size: 0.75rem;
    margin-bottom: 0.25rem;
    gap: 0.5rem;
  }

  .event-time {
    font-weight: 600;
    color: #495057;
  }

  .event-calendar {
    font-weight: 500;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .event-summary {
    font-size: 0.9rem;
    color: #212529;
    font-weight: 500;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .event-location {
    font-size: 0.8rem;
    color: #6c757d;
    font-weight: normal;
    margin-left: 0.5rem;
  }

  @media (max-width: 576px) {
    .event-item {
      padding: 0.4rem 0.6rem;
    }

    .event-meta {
      font-size: 0.7rem;
    }

    .event-summary {
      font-size: 0.85rem;
    }
  }
</style>
