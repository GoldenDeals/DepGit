<script lang="ts">
  import Modal from "../Modal.svelte";
  
  export let isOpen: boolean = false;
  export let onClose: () => void;
  export let onDelete: () => void;
  export let user: { 
    id: number, 
    username: string 
  } | null = null;
  
  let confirmUsername = '';
  let error = '';
  
  function handleDelete() {
    if (!user) return;
    
    if (confirmUsername !== user.username) {
      error = 'Username does not match';
      return;
    }
    
    error = '';
    confirmUsername = '';
    onDelete();
  }
  
  // Reset form when modal is closed or opened
  $: if (isOpen) {
    confirmUsername = '';
    error = '';
  }
</script>

<Modal {isOpen} title="Delete User" {onClose}>
  {#if user}
    <div class="space-y-4">
      <div class="bg-yellow-50 border-l-4 border-yellow-400 p-4">
        <div class="flex">
          <div class="flex-shrink-0">
            <svg class="h-5 w-5 text-yellow-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
            </svg>
          </div>
          <div class="ml-3">
            <p class="text-sm text-yellow-700">
              <strong>Warning:</strong> This action cannot be undone. This will permanently delete the user <strong>{user.username}</strong> and all associated data.
            </p>
          </div>
        </div>
      </div>
      
      <p class="text-sm text-gray-500">
        Please type <strong>{user.username}</strong> to confirm.
      </p>
      
      {#if error}
        <div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">
          {error}
        </div>
      {/if}
      
      <div>
        <input
          type="text"
          id="confirm"
          bind:value={confirmUsername}
          class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-red-500 focus:border-red-500"
          placeholder={user.username}
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
          type="button"
          class="px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-red-600 hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500"
          on:click={handleDelete}
        >
          Delete User
        </button>
      </div>
    </div>
  {/if}
</Modal>
