<script lang="ts">
  import Modal from "../Modal.svelte";

  export let isOpen: boolean = false;
  export let onClose: () => void;
  export let onLogin: (credentials: { username: string, password: string }) => void;
  export let onRegisterClick: () => void;
  export let onForgotPasswordClick: () => void;

  let username: string = '';
  let password: string = '';
  let rememberMe: boolean = false;
  let error: string = '';
  let isLoading: boolean = false;

  function handleSubmit() {
    // Validate form
    if (!username.trim()) {
      error = 'Username is required';
      return;
    }

    if (!password) {
      error = 'Password is required';
      return;
    }

    // Show loading state
    isLoading = true;

    // In a real app, this would be an async call
    setTimeout(() => {
      // Submit form
      onLogin({ username, password });

      // Reset form
      resetForm();

      // Hide loading state
      isLoading = false;

      // Close modal
      onClose();
    }, 500);
  }

  function resetForm() {
    username = '';
    password = '';
    rememberMe = false;
    error = '';
  }

  // Reset form when modal is opened
  $: if (isOpen) {
    resetForm();
  }

  function handleRegisterClick() {
    onClose();
    onRegisterClick();
  }
</script>

<Modal {isOpen} title="Log in to DepGit" {onClose}>
  <form on:submit|preventDefault={handleSubmit} class="space-y-4">
    {#if error}
      <div class="bg-error bg-opacity-10 border border-error text-error px-4 py-3 rounded">
        {error}
      </div>
    {/if}

    <div>
      <label for="username" class="block text-sm font-medium mb-1">
        Username
      </label>
      <input
        type="text"
        id="username"
        bind:value={username}
        class="w-full px-3 py-2 border border-border rounded-md shadow-sm focus:outline-none focus:ring-primary focus:border-primary bg-background"
        placeholder="Enter your username"
        autocomplete="username"
        disabled={isLoading}
      />
    </div>

    <div>
      <label for="password" class="block text-sm font-medium mb-1">
        Password
      </label>
      <input
        type="password"
        id="password"
        bind:value={password}
        class="w-full px-3 py-2 border border-border rounded-md shadow-sm focus:outline-none focus:ring-primary focus:border-primary bg-background"
        placeholder="Enter your password"
        autocomplete="current-password"
        disabled={isLoading}
      />
    </div>

    <div class="flex items-center justify-between">
      <div class="flex items-center">
        <input
          id="remember-me"
          type="checkbox"
          bind:checked={rememberMe}
          class="h-4 w-4 text-primary focus:ring-primary border-border rounded"
          disabled={isLoading}
        />
        <label for="remember-me" class="ml-2 block text-sm text-text-muted">
          Remember me
        </label>
      </div>

      <div class="text-sm">
        <button
          type="button"
          class="font-medium text-primary hover:text-primary-hover"
          on:click={onForgotPasswordClick}
        >
          Forgot your password?
        </button>
      </div>
    </div>

    <div>
      <button
        type="submit"
        class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-primary hover:bg-primary-hover focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary disabled:opacity-50 disabled:cursor-not-allowed"
        disabled={isLoading}
      >
        {#if isLoading}
          <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          Logging in...
        {:else}
          Log in
        {/if}
      </button>
    </div>

    <div class="text-center mt-4">
      <p class="text-sm text-text-muted">
        Don't have an account?
        <button
          type="button"
          class="font-medium text-primary hover:text-primary-hover"
          on:click={handleRegisterClick}
        >
          Register now
        </button>
      </p>
    </div>
  </form>
</Modal>
