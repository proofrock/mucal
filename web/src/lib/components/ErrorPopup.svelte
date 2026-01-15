<script lang="ts">
  import { errorStore } from '../stores/error.svelte';

  function handleClose() {
    errorStore.clearError();
  }
</script>

{#if errorStore.show}
  <div class="position-fixed top-0 end-0 p-3" style="z-index: 11">
    <div class="toast show" role="alert">
      <div class="toast-header bg-danger text-white">
        <strong class="me-auto">Error</strong>
        <button
          type="button"
          class="btn-close btn-close-white"
          aria-label="Close error message"
          onclick={handleClose}
        ></button>
      </div>
      <div class="toast-body">
        {errorStore.message}
      </div>
    </div>
  </div>
{/if}

<style>
  .toast {
    min-width: 300px;
  }
</style>
