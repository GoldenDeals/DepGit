<script lang="ts">
  import Modal from "../Modal.svelte";
  
  export let isOpen: boolean = false;
  export let onClose: () => void;
  export let onAdd: (userData: { username: string, email: string, role: string, password: string }) => void;
  
  let username: string = '';
  let email: string = '';
  let role: string = 'Developer';
  let password: string = '';
  let confirmPassword: string = '';
  let error: string = '';
  
  const roles = ['Administrator', 'Developer', 'Viewer'];
  
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
    
    if (!password) {
      error = 'Password is required';
      return;
    }
    
    if (password !== confirmPassword) {
      error = 'Passwords do not match';
      return;
    }
    
    // Submit form
    onAdd({ username, email, role, password });
    
    // Reset form
    resetForm();
    
    // Close modal
    onClose();
  }
  
  function resetForm() {
    username = '';
    email = '';
    role = 'Developer';
    password = '';
    confirmPassword = '';
    error = '';
  }
  
  // Reset form when modal is opened
  $: if (isOpen) {
    resetForm();
  }
</script>

<Modal {isOpen} title="Add New User" {onClose}>
  <form on:submit|preventDefault={handleSubmit} class="space-y-4">
    {#if error}
      <div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">
        {error}
      </div>
    {/if}
    
    <div>
      <label for="username" class="block text-sm font-medium text-gray-700 mb-1">
        Username *
      </label>
      <input
        type="text"
        id="username"
        bind:value={username}
        class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
        placeholder="e.g. johndoe"
      />
    </div>
    
    <div>
      <label for="email" class="block text-sm font-medium text-gray-700 mb-1">
        Email *
      </label>
      <input
        type="email"
        id="email"
        bind:value={email}
        class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
        placeholder="e.g. john@example.com"
      />
    </div>
    
    <div>
      <label for="role" class="block text-sm font-medium text-gray-700 mb-1">
        Role *
      </label>
      <select
        id="role"
        bind:value={role}
        class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
      >
        {#each roles as roleOption}
          <option value={roleOption}>{roleOption}</option>
        {/each}
      </select>
    </div>
    
    <div>
      <label for="password" class="block text-sm font-medium text-gray-700 mb-1">
        Password *
      </label>
      <input
        type="password"
        id="password"
        bind:value={password}
        class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
      />
    </div>
    
    <div>
      <label for="confirmPassword" class="block text-sm font-medium text-gray-700 mb-1">
        Confirm Password *
      </label>
      <input
        type="password"
        id="confirmPassword"
        bind:value={confirmPassword}
        class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
      />
    </div>
    
    <div class="flex justify-end space-x-3 pt-4">
      <button
        type="button"
        class="px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
        on:click={onClose}
      >
        Cancel
      </button>
      <button
        type="submit"
        class="px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
      >
        Create User
      </button>
    </div>
  </form>
</Modal>
