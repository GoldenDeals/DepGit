<script lang="ts">
  import Modal from "../Modal.svelte";

  export let isOpen: boolean = false;
  export let onClose: () => void;
  export let onLoginClick: () => void;
  export let onRequestReset: (email: string) => void;

  let email: string = '';
  let error: string = '';
  let isLoading: boolean = false;
  let isSuccess: boolean = false;

  function handleSubmit() {
    // Reset states
    error = '';
    isSuccess = false;

    // Validate form
    if (!email.trim()) {
      error = 'Email is required';
      return;
    }

    if (!isValidEmail(email)) {
      error = 'Please enter a valid email address';
      return;
    }

    // Show loading state
    isLoading = true;

    // In a real app, this would be an async call
    setTimeout(() => {
      // Submit form
      onRequestReset(email);

      // Show success message
      isSuccess = true;

      // Hide loading state
      isLoading = false;
    }, 1000);
  }

  function isValidEmail(email: string): boolean {
    const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return re.test(email);
  }

  function resetForm() {
    email = '';
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
        <p class="font-medium">Password reset email sent!</p>
        <p class="text-sm mt-1">
          We've sent an email to <strong>{email}</strong> with instructions to reset your password.
          Please check your inbox and follow the link in the email.
        </p>
      </div>

      <div class="flex justify-end space-x-3 pt-4">
        <button
          type="button"
          class="px-4 py-2 border border-border rounded-md shadow-sm text-sm font-medium text-text bg-surface hover:bg-surface-hover focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary"
          on:click={onClose}
        >
          Close
        </button>
        <button
          type="button"
          class="px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-primary hover:bg-primary-hover focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary"
          on:click={handleLoginClick}
        >
          Return to login
        </button>
      </div>
    {:else}
      <p class="text-sm text-text-muted">
        Enter the email address associated with your account, and we'll send you a link to reset your password.
      </p>

      <form on:submit|preventDefault={handleSubmit} class="space-y-4">
        {#if error}
          <div class="bg-error bg-opacity-10 border border-error text-error px-4 py-3 rounded">
            {error}
          </div>
        {/if}

        <div>
          <label for="email" class="block text-sm font-medium mb-1">
            Email address
          </label>
          <input
            type="email"
            id="email"
            bind:value={email}
            class="w-full px-3 py-2 border border-border rounded-md shadow-sm focus:outline-none focus:ring-primary focus:border-primary bg-background"
            placeholder="you@example.com"
            autocomplete="email"
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
              Sending...
            {:else}
              Send reset link
            {/if}
          </button>
        </div>
      </form>

      <div class="text-center mt-4">
        <p class="text-sm text-text-muted">
          Remember your password?
          <button
            type="button"
            class="font-medium text-primary hover:text-primary-hover"
            on:click={handleLoginClick}
          >
            Log in
          </button>
        </p>
      </div>
    {/if}
  </div>
</Modal>
