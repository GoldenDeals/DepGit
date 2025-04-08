<script lang="ts">
  import Modal from "../Modal.svelte";

  export let isOpen: boolean = false;
  export let onClose: () => void;
  export let onLoginClick: () => void;
  export let onResetPassword: (data: { token: string, password: string }) => void;
  export let token: string = '';

  let password: string = '';
  let confirmPassword: string = '';
  let error: string = '';
  let isLoading: boolean = false;
  let isSuccess: boolean = false;

  function handleSubmit() {
    // Reset states
    error = '';

    // Validate form
    if (!password) {
      error = 'Password is required';
      return;
    }

    if (password.length < 8) {
      error = 'Password must be at least 8 characters long';
      return;
    }

    if (password !== confirmPassword) {
      error = 'Passwords do not match';
      return;
    }

    // Show loading state
    isLoading = true;

    // In a real app, this would be an async call
    setTimeout(() => {
      // Submit form
      onResetPassword({ token, password });

      // Show success message
      isSuccess = true;

      // Hide loading state
      isLoading = false;
    }, 1000);
  }

  function resetForm() {
    password = '';
    confirmPassword = '';
    error = '';
    isSuccess = false;
  }

  // Reset form when modal is opened
  $: if (isOpen) {
    resetForm();
  }

  function handleLoginClick() {
    resetForm();
    onClose();
    onLoginClick();
  }
</script>

<Modal {isOpen} title="Reset your password" {onClose}>
  <div class="space-y-4">
    {#if isSuccess}
      <div class="bg-success bg-opacity-10 border border-success text-success px-4 py-3 rounded">
        <p class="font-medium">Password reset successful!</p>
        <p class="text-sm mt-1">
          Your password has been reset successfully. You can now log in with your new password.
        </p>
      </div>

      <div class="flex justify-end space-x-3 pt-4">
        <button
          type="button"
          class="px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-primary hover:bg-primary-hover focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary"
          on:click={handleLoginClick}
        >
          Go to login
        </button>
      </div>
    {:else}
      <p class="text-sm text-text-muted">
        Create a new password for your account. Password must be at least 8 characters long.
      </p>

      <form on:submit|preventDefault={handleSubmit} class="space-y-4">
        {#if error}
          <div class="bg-error bg-opacity-10 border border-error text-error px-4 py-3 rounded">
            {error}
          </div>
        {/if}

        <div>
          <label for="password" class="block text-sm font-medium mb-1">
            New password
          </label>
          <input
            type="password"
            id="password"
            bind:value={password}
            class="w-full px-3 py-2 border border-border rounded-md shadow-sm focus:outline-none focus:ring-primary focus:border-primary bg-background"
            placeholder="Enter your new password"
            autocomplete="new-password"
            disabled={isLoading}
          />
        </div>

        <div>
          <label for="confirm-password" class="block text-sm font-medium mb-1">
            Confirm new password
          </label>
          <input
            type="password"
            id="confirm-password"
            bind:value={confirmPassword}
            class="w-full px-3 py-2 border border-border rounded-md shadow-sm focus:outline-none focus:ring-primary focus:border-primary bg-background"
            placeholder="Confirm your new password"
            autocomplete="new-password"
            disabled={isLoading}
          />
        </div>

        <div class="flex justify-end space-x-3 pt-4">
          <button
            type="button"
            class="px-4 py-2 border border-border rounded-md shadow-sm text-sm font-medium text-text bg-surface hover:bg-surface-hover focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary"
            on:click={onClose}
            disabled={isLoading}
          >
            Cancel
          </button>
          <button
            type="submit"
            class="px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-primary hover:bg-primary-hover focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary disabled:opacity-50 disabled:cursor-not-allowed"
            disabled={isLoading}
          >
            {#if isLoading}
              <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white inline-block" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              Resetting...
            {:else}
              Reset password
            {/if}
          </button>
        </div>
      </form>
    {/if}
  </div>
</Modal>
