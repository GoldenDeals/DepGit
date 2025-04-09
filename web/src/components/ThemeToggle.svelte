<script lang="ts">
  import { onMount } from 'svelte';

  // Theme state
  let isDarkMode = false;

  // Initialize theme based on user preference or system preference
  onMount(() => {
    // Check for saved theme preference
    const savedTheme = localStorage.getItem('theme');

    if (savedTheme === 'dark') {
      enableDarkMode();
    } else if (savedTheme === 'light') {
      enableLightMode();
    } else {
      // Check system preference
      if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
        enableDarkMode();
      } else {
        enableLightMode();
      }
    }

    // Listen for system preference changes
    window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', (e) => {
      if (!localStorage.getItem('theme')) {
        if (e.matches) {
          enableDarkMode();
        } else {
          enableLightMode();
        }
      }
    });
  });

  function toggleTheme() {
    if (isDarkMode) {
      enableLightMode();
    } else {
      enableDarkMode();
    }
  }

  function enableDarkMode() {
    document.documentElement.classList.add('dark');
    localStorage.setItem('theme', 'dark');
    isDarkMode = true;
  }

  function enableLightMode() {
    document.documentElement.classList.remove('dark');
    localStorage.setItem('theme', 'light');
    isDarkMode = false;
  }
</script>

<button
  type="button"
  class="p-2 rounded-full bg-surface text-text hover:bg-surface-hover focus:outline-none focus:ring-2 focus:ring-primary transition-colors duration-200 dark:bg-primary dark:text-white dark:hover:bg-primary-light"
  on:click={toggleTheme}
  aria-label={isDarkMode ? "Switch to light mode" : "Switch to dark mode"}
  title={isDarkMode ? "Switch to light mode" : "Switch to dark mode"}
>
  {#if isDarkMode}
    <!-- Sun icon for light mode -->
    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
    </svg>
  {:else}
    <!-- Moon icon for dark mode -->
    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
    </svg>
  {/if}
</button>
