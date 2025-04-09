<script lang="ts">
  import Modal from "../Modal.svelte";
  
  export let isOpen: boolean = false;
  export let onClose: () => void;
  export let onAdd: (repo: { name: string, description: string }) => void;
  
  let name: string = '';
  let description: string = '';
  let error: string = '';
  
  function handleSubmit() {
    // Validate form
    if (!name.trim()) {
      error = 'Repository name is required';
      return;
    }
    
    // Submit form
    onAdd({ name, description });
    
    // Reset form
    name = '';
    description = '';
    error = '';
    
    // Close modal
    onClose();
  }
</script>

<Modal {isOpen} title="Add New Repository" {onClose}>
  <form on:submit|preventDefault={handleSubmit} class="space-y-4">
    {#if error}
      <div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">
        {error}
      </div>
    {/if}
    
    <div>
      <label for="name" class="block text-sm font-medium text-gray-700 mb-1">
        Repository Name *
      </label>
      <input
        type="text"
        id="name"
        bind:value={name}
        class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
        placeholder="e.g. my-project"
      />
    </div>
    
    <div>
      <label for="description" class="block text-sm font-medium text-gray-700 mb-1">
        Description
      </label>
      <textarea
        id="description"
        bind:value={description}
        rows="3"
        class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
        placeholder="Brief description of the repository"
      ></textarea>
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
        Create Repository
      </button>
    </div>
  </form>
</Modal>
