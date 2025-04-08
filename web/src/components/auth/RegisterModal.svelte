<script lang="ts">
  import Modal from "../Modal.svelte";

  export let isOpen: boolean = false;
  export let onClose: () => void;
  export let onRegister: (userData: {
    username: string,
    email: string,
    password: string
  }) => void;
  export let onLoginClick: () => void;

  let username: string = '';
  let email: string = '';
  let password: string = '';
  let confirmPassword: string = '';
  let agreeTerms: boolean = false;
  let error: string = '';
  let isLoading: boolean = false;

  function handleSubmit() {
    // Validate form
    if (!username.trim()) {
      error = 'Username is required';
      return;
    }

    if (!email.trim()) {
      error = 'Email is required';
      return;
    }

    if (!isValidEmail(email)) {
      error = 'Please enter a valid email address';
      return;
    }

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

    if (!agreeTerms) {
      error = 'You must agree to the terms and conditions';
      return;
    }

    // Show loading state
    isLoading = true;

    // In a real app, this would be an async call
    setTimeout(() => {
      // Submit form
      onRegister({ username, email, password });

      // Reset form
      resetForm();

      // Hide loading state
      isLoading = false;

      // Close modal
      onClose();
    }, 500);
  }

  function isValidEmail(email: string): boolean {
    const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return re.test(email);
  }

  function resetForm() {
    username = '';
    email = '';
    password = '';
    confirmPassword = '';
    agreeTerms = false;
    error = '';
  }

  // Reset form when modal is opened
  $: if (isOpen) {
    resetForm();
  }

  function handleLoginClick() {
    onClose();
    onLoginClick();
  }
</script>

<Modal {isOpen} title="Create an account" {onClose} maxWidth="max-w-lg">
  <form on:submit|preventDefault={handleSubmit} class="space-y-4">
    {#if error}
      <div class="bg-red bg-opacity-10 border border-red text-red px-4 py-3 rounded">
        {error}
      </div>
    {/if}

    <div>
      <label for="username" class="block text-sm font-medium mb-1">
        Username *
      </label>
      <input
        type="text"
        id="username"
        bind:value={username}
        class="w-full px-3 py-2 border border-overlay0 rounded-md shadow-sm focus:outline-none focus:ring-blue focus:border-blue bg-base"
        placeholder="Choose a username"
        autocomplete="username"
        disabled={isLoading}
      />
      <p class="mt-1 text-xs text-subtext0">
        This will be your unique identifier on DepGit.
      </p>
    </div>

    <div>
      <label for="email" class="block text-sm font-medium mb-1">
        Email address *
      </label>
      <input
        type="email"
        id="email"
        bind:value={email}
        class="w-full px-3 py-2 border border-overlay0 rounded-md shadow-sm focus:outline-none focus:ring-blue focus:border-blue bg-base"
        placeholder="you@example.com"
        autocomplete="email"
        disabled={isLoading}
      />
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div>
        <label for="password" class="block text-sm font-medium mb-1">
          Password *
        </label>
        <input
          type="password"
          id="password"
          bind:value={password}
          class="w-full px-3 py-2 border border-overlay0 rounded-md shadow-sm focus:outline-none focus:ring-blue focus:border-blue bg-base"
          placeholder="Create a password"
          autocomplete="new-password"
          disabled={isLoading}
        />
        <p class="mt-1 text-xs text-subtext0">
          Must be at least 8 characters.
        </p>
      </div>

      <div>
        <label for="confirm-password" class="block text-sm font-medium mb-1">
          Confirm password *
        </label>
        <input
          type="password"
          id="confirm-password"
          bind:value={confirmPassword}
          class="w-full px-3 py-2 border border-overlay0 rounded-md shadow-sm focus:outline-none focus:ring-blue focus:border-blue bg-base"
          placeholder="Confirm your password"
          autocomplete="new-password"
          disabled={isLoading}
        />
      </div>
    </div>

    <div class="flex items-center">
      <input
        id="agree-terms"
        type="checkbox"
        bind:checked={agreeTerms}
        class="h-4 w-4 text-blue focus:ring-blue border-overlay0 rounded"
        disabled={isLoading}
      />
      <label for="agree-terms" class="ml-2 block text-sm text-subtext0">
        I agree to the <a href="/terms" class="text-blue hover:text-sapphire">Terms of Service</a> and <a href="/privacy" class="text-blue hover:text-sapphire">Privacy Policy</a>
      </label>
    </div>

    <div>
      <button
        type="submit"
        class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue hover:bg-sapphire focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue disabled:opacity-50 disabled:cursor-not-allowed"
        disabled={isLoading}
      >
        {#if isLoading}
          <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          Creating account...
        {:else}
          Create account
        {/if}
      </button>
    </div>

    <div class="text-center mt-4">
      <p class="text-sm text-subtext0">
        Already have an account?
        <button
          type="button"
          class="font-medium text-blue hover:text-sapphire"
          on:click={handleLoginClick}
        >
          Log in
        </button>
      </p>
    </div>
  </form>
</Modal>
