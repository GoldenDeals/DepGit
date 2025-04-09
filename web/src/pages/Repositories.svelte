<script lang="ts">
  import AddRepositoryModal from "../components/repositories/AddRepositoryModal.svelte";
  import ViewRepositoryModal from "../components/repositories/ViewRepositoryModal.svelte";
  import DeleteRepositoryModal from "../components/repositories/DeleteRepositoryModal.svelte";

  // Define repository type
  type Repository = {
    id: number;
    name: string;
    description: string;
    owner: string;
  };

  // Mock data for repositories
  let repositories: Repository[] = [
    { id: 1, name: 'project-alpha', description: 'Main project repository', owner: 'admin' },
    { id: 2, name: 'backend-api', description: 'Backend API service', owner: 'developer1' },
    { id: 3, name: 'frontend-app', description: 'Frontend application', owner: 'developer2' },
  ];

  // Modal states
  let isAddModalOpen = false;
  let isViewModalOpen = false;
  let isDeleteModalOpen = false;
  let selectedRepository: Repository | null = null;

  // Modal handlers
  function openAddModal() {
    isAddModalOpen = true;
  }

  function openViewModal(repo: Repository) {
    selectedRepository = repo;
    isViewModalOpen = true;
  }

  function openDeleteModal(repo: Repository) {
    selectedRepository = repo;
    isDeleteModalOpen = true;
  }

  function closeModals() {
    isAddModalOpen = false;
    isViewModalOpen = false;
    isDeleteModalOpen = false;
  }

  function handleAddRepository(repoData: { name: string, description: string }) {
    // In a real app, this would make an API call
    const newRepo: Repository = {
      id: repositories.length + 1,
      name: repoData.name,
      description: repoData.description,
      owner: 'current-user' // In a real app, this would be the current user
    };
    repositories = [...repositories, newRepo];
  }

  function handleDeleteRepository() {
    // In a real app, this would make an API call
    if (selectedRepository) {
      repositories = repositories.filter(repo => repo.id !== selectedRepository.id);
      closeModals();
    }
  }
</script>

<div>
  <div class="flex justify-between items-center mb-6">
    <h1 class="text-2xl font-bold">Repositories</h1>
    <button
      class="bg-accent hover:bg-accent-hover text-white px-4 py-2 rounded"
      on:click={openAddModal}
    >
      New Repository
    </button>
  </div>

  <div class="bg-surface rounded-lg shadow overflow-hidden">
    <table class="min-w-full divide-y divide-border">
      <thead class="bg-surface-hover">
        <tr>
          <th class="px-6 py-3 text-left text-xs font-medium text-text-muted uppercase tracking-wider">Name</th>
          <th class="px-6 py-3 text-left text-xs font-medium text-text-muted uppercase tracking-wider">Description</th>
          <th class="px-6 py-3 text-left text-xs font-medium text-text-muted uppercase tracking-wider">Owner</th>
          <th class="px-6 py-3 text-left text-xs font-medium text-text-muted uppercase tracking-wider">Actions</th>
        </tr>
      </thead>
      <tbody class="bg-surface divide-y divide-border">
        {#each repositories as repo}
          <tr>
            <td class="px-6 py-4 whitespace-nowrap">
              <div class="text-sm font-medium text-primary">{repo.name}</div>
            </td>
            <td class="px-6 py-4">
              <div class="text-sm text-text">{repo.description}</div>
            </td>
            <td class="px-6 py-4 whitespace-nowrap">
              <div class="text-sm text-text">{repo.owner}</div>
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
              <button
                class="text-primary hover:text-primary-hover mr-3"
                on:click={() => openViewModal(repo)}
              >
                View
              </button>
              <button
                class="text-error hover:text-error"
                on:click={() => openDeleteModal(repo)}
              >
                Delete
              </button>
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
  </div>
</div>

<!-- Modals -->
<AddRepositoryModal
  isOpen={isAddModalOpen}
  onClose={closeModals}
  onAdd={handleAddRepository}
/>

<ViewRepositoryModal
  isOpen={isViewModalOpen}
  onClose={closeModals}
  repository={selectedRepository}
/>

<DeleteRepositoryModal
  isOpen={isDeleteModalOpen}
  onClose={closeModals}
  onDelete={handleDeleteRepository}
  repository={selectedRepository}
/>
