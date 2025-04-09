<script lang="ts">
  import Modal from "../Modal.svelte";

  export let isOpen: boolean = false;
  export let onClose: () => void;
  export let repository: {
    id: number,
    name: string,
    description: string,
    owner: string
  } | null = null;

  // Clone URL would typically come from the backend
  $: cloneUrl = repository ? `ssh://git@example.com/${repository.name}.git` : '';

  let copySuccess = false;

  function copyToClipboard() {
    if (!cloneUrl) return;

    navigator.clipboard.writeText(cloneUrl)
      .then(() => {
        copySuccess = true;
        setTimeout(() => {
          copySuccess = false;
        }, 2000);
      })
      .catch(err => {
        console.error('Failed to copy text: ', err);
      });
  }
</script>

<Modal {isOpen} title="Repository Details" {onClose} maxWidth="max-w-2xl">
  {#if repository}
    <div class="space-y-6">
      <div>
        <h3 class="text-lg font-medium text-gray-900">{repository.name}</h3>
        <p class="mt-1 text-sm text-gray-500">{repository.description || 'No description provided'}</p>
      </div>

      <div class="border-t border-gray-200 pt-4">
        <dl class="divide-y divide-gray-200">
          <div class="py-3 flex justify-between">
            <dt class="text-sm font-medium text-gray-500">Owner</dt>
            <dd class="text-sm text-gray-900">{repository.owner}</dd>
          </div>

          <div class="py-3">
            <dt class="text-sm font-medium text-gray-500 mb-1">Clone URL</dt>
            <dd class="mt-1 flex">
              <input
                type="text"
                readonly
                value={cloneUrl}
                class="flex-1 px-3 py-2 border border-gray-300 rounded-md shadow-sm bg-gray-50 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
              />
              <button
                type="button"
                class="ml-2 px-3 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
                on:click={copyToClipboard}
              >
                {copySuccess ? 'Copied!' : 'Copy'}
              </button>
            </dd>
          </div>
        </dl>
      </div>

      <div class="border-t border-gray-200 pt-4">
        <h4 class="text-sm font-medium text-gray-500 mb-2">Usage Instructions</h4>
        <div class="bg-gray-50 rounded-md p-3">
          <pre class="text-xs text-gray-700 overflow-x-auto">
# Clone this repository
git clone {cloneUrl}

# Navigate to the repository directory
cd {repository.name}

# Create a new branch
git checkout -b feature/your-feature-name

# Make changes and commit them
git add .
git commit -m "Your commit message"

# Push changes to the remote repository
git push origin feature/your-feature-name</pre>
        </div>
      </div>

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
