<script lang="ts">
  import Modal from "../Modal.svelte";
  
  export let isOpen: boolean = false;
  export let onClose: () => void;
  export let user: { 
    id: number, 
    username: string 
  } | null = null;
  
  // Mock SSH keys data
  let sshKeys = [
    { id: 1, name: 'Work Laptop', fingerprint: 'SHA256:ZmL5t8fGgG3T5K7d8J9Ls9f2K8dH3Jd9s8F7g6T5K4d3', added: '2023-01-15' },
    { id: 2, name: 'Home Desktop', fingerprint: 'SHA256:9sD8f7G6h5J4k3L2m1N0p9O8i7U6y5T4r3E2w1Q0', added: '2023-02-20' }
  ];
  
  let newKeyName: string = '';
  let newKeyValue: string = '';
  let error: string = '';
  let addingKey: boolean = false;
  
  function toggleAddKey() {
    addingKey = !addingKey;
    if (addingKey) {
      newKeyName = '';
      newKeyValue = '';
      error = '';
    }
  }
  
  function addKey() {
    if (!newKeyName.trim()) {
      error = 'Key name is required';
      return;
    }
    
    if (!newKeyValue.trim()) {
      error = 'SSH key is required';
      return;
    }
    
    // Validate SSH key format (basic validation)
    if (!newKeyValue.startsWith('ssh-rsa ') && !newKeyValue.startsWith('ssh-ed25519 ')) {
      error = 'Invalid SSH key format. Key should start with ssh-rsa or ssh-ed25519';
      return;
    }
    
    // In a real app, this would make an API call to add the key
    const newKey = {
      id: sshKeys.length + 1,
      name: newKeyName,
      fingerprint: `SHA256:${Math.random().toString(36).substring(2, 15)}${Math.random().toString(36).substring(2, 15)}`,
      added: new Date().toISOString().split('T')[0]
    };
    
    sshKeys = [...sshKeys, newKey];
    
    // Reset form
    newKeyName = '';
    newKeyValue = '';
    error = '';
    addingKey = false;
  }
  
  function deleteKey(id: number) {
    // In a real app, this would make an API call to delete the key
    sshKeys = sshKeys.filter(key => key.id !== id);
  }
  
  // Reset state when modal is opened
  $: if (isOpen) {
    addingKey = false;
    error = '';
  }
</script>

<Modal {isOpen} title={user ? `SSH Keys for ${user.username}` : 'SSH Keys'} {onClose} maxWidth="max-w-2xl">
  {#if user}
    <div class="space-y-6">
      <div class="flex justify-between items-center">
        <h3 class="text-lg font-medium text-gray-900">Public Keys</h3>
        {#if !addingKey}
          <button
            type="button"
            class="px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
            on:click={toggleAddKey}
          >
            Add SSH Key
          </button>
        {/if}
      </div>
      
      {#if addingKey}
        <div class="bg-gray-50 p-4 rounded-md border border-gray-200">
          <h4 class="text-md font-medium text-gray-700 mb-3">Add New SSH Key</h4>
          
          {#if error}
            <div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">
              {error}
            </div>
          {/if}
          
          <div class="space-y-4">
            <div>
              <label for="keyName" class="block text-sm font-medium text-gray-700 mb-1">
                Key Name *
              </label>
              <input
                type="text"
                id="keyName"
                bind:value={newKeyName}
                class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                placeholder="e.g. Work Laptop"
              />
            </div>
            
            <div>
              <label for="sshKey" class="block text-sm font-medium text-gray-700 mb-1">
                SSH Public Key *
              </label>
              <textarea
                id="sshKey"
                bind:value={newKeyValue}
                rows="3"
                class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                placeholder="ssh-rsa AAAAB3NzaC1yc2E..."
              ></textarea>
              <p class="mt-1 text-xs text-gray-500">
                Paste your public SSH key, usually found in ~/.ssh/id_rsa.pub or ~/.ssh/id_ed25519.pub
              </p>
            </div>
            
            <div class="flex justify-end space-x-3">
              <button
                type="button"
                class="px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
                on:click={toggleAddKey}
              >
                Cancel
              </button>
              <button
                type="button"
                class="px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
                on:click={addKey}
              >
                Add Key
              </button>
            </div>
          </div>
        </div>
      {/if}
      
      {#if sshKeys.length === 0}
        <div class="text-center py-6 text-gray-500">
          No SSH keys found. Add a key to allow Git access.
        </div>
      {:else}
        <div class="bg-white rounded-lg shadow overflow-hidden">
          <table class="min-w-full divide-y divide-gray-200">
            <thead class="bg-gray-50">
              <tr>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Name</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Fingerprint</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Added</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
              {#each sshKeys as key}
                <tr>
                  <td class="px-6 py-4 whitespace-nowrap">
                    <div class="text-sm font-medium text-gray-900">{key.name}</div>
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap">
                    <div class="text-sm text-gray-500 font-mono">{key.fingerprint}</div>
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap">
                    <div class="text-sm text-gray-500">{key.added}</div>
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
                    <button 
                      class="text-red-600 hover:text-red-900"
                      on:click={() => deleteKey(key.id)}
                    >
                      Delete
                    </button>
                  </td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      {/if}
      
      <div class="flex justify-end pt-4">
        <button
          type="button"
          class="px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
          on:click={onClose}
        >
          Close
        </button>
      </div>
    </div>
  {/if}
</Modal>
