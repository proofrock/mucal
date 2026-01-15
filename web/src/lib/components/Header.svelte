<script lang="ts">
  import { calendarStore } from '../stores/calendar.svelte';
</script>

<nav class="navbar navbar-dark bg-primary">
  <div class="container-fluid px-3">
    <span class="navbar-brand mb-0 h1">μCal</span>
    {#if calendarStore.version}
      <span class="badge bg-light text-primary">v{calendarStore.version}</span>
    {/if}
  </div>
</nav>

<style>
  .navbar {
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  }

  .navbar-brand {
    font-size: 1.25rem;
    font-weight: 700;
  }

  @media (max-width: 576px) {
    .navbar-brand {
      font-size: 1.1rem;
    }

    .badge {
      font-size: 0.7rem;
    }
  }
</style>
