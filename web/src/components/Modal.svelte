<script lang="ts">
  import { onMount } from 'svelte';

  export let isOpen: boolean = false;
  export let title: string = '';
  export let onClose: () => void;
  export let maxWidth: string = 'max-w-md';

  // Handle escape key press to close the modal
  function handleEscape(event: KeyboardEvent) {
    if (event.key === 'Escape') {
      onClose();
    }
  }

  // Handle backdrop click to close the modal
  function handleBackdropClick() {
    onClose();
  }

  // Focus trap management
  let modalElement: HTMLDivElement;
  let previouslyFocusedElement: Element | null = null;

  onMount(() => {
    return () => {
      // Restore focus when component is destroyed
      if (previouslyFocusedElement && 'focus' in previouslyFocusedElement) {
        (previouslyFocusedElement as HTMLElement).focus();
      }
    };
  });

  // When modal opens, save currently focused element and focus the modal
  $: if (isOpen && typeof document !== 'undefined') {
    previouslyFocusedElement = document.activeElement;
    // Focus the modal in the next tick after it's rendered
    setTimeout(() => {
      if (modalElement) {
        modalElement.focus();
      }
    }, 0);
  }
</script>

<svelte:window on:keydown={handleEscape} />

{#if isOpen}
  <!-- Backdrop with flex container -->
  <div class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-center justify-center p-4">
    <!-- Modal dialog -->
    <div
      bind:this={modalElement}
      role="dialog"
      aria-modal="true"
      aria-labelledby="dialog-title"
      class={`bg-surface0 rounded-lg shadow-xl w-full ${maxWidth} overflow-hidden dark:bg-surface1 dark:shadow-2xl`}
      tabindex="-1"
    >
      <!-- Header with title and close button -->
      <div class="flex justify-between items-center border-b border-overlay0 p-4 dark:border-surface2">
        <h2 id="dialog-title" class="text-xl font-semibold dark:text-white">{title}</h2>
        <button
          class="text-subtext0 hover:text-text focus:outline-none"
          on:click={onClose}
          aria-label="Close dialog"
        >
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
          </svg>
        </button>
      </div>
      <!-- Content area -->
      <div class="p-4 dark:text-white">
        <slot />
      </div>
    </div>

    <!-- Invisible button to capture clicks outside the modal -->
    <button
      class="fixed inset-0 w-full h-full opacity-0 cursor-default"
      on:click={handleBackdropClick}
      tabindex="-1"
      aria-hidden="true"
      style="z-index: -1;"
    ></button>
  </div>
{/if}

<!-- No default export needed, Svelte components are automatically exported -->
