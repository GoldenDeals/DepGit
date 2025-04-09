<script lang="ts">
  import AddUserModal from "../components/users/AddUserModal.svelte";
  import EditUserModal from "../components/users/EditUserModal.svelte";
  import SSHKeysModal from "../components/users/SSHKeysModal.svelte";
  import DeleteUserModal from "../components/users/DeleteUserModal.svelte";

  // Define user type
  type User = {
    id: number;
    username: string;
    email: string;
    role: string;
  };

  // Mock data for users
  let users: User[] = [
    { id: 1, username: 'admin', email: 'admin@example.com', role: 'Administrator' },
    { id: 2, username: 'developer1', email: 'dev1@example.com', role: 'Developer' },
    { id: 3, username: 'developer2', email: 'dev2@example.com', role: 'Developer' },
  ];

  // Modal states
  let isAddModalOpen = false;
  let isEditModalOpen = false;
  let isSSHKeysModalOpen = false;
  let isDeleteModalOpen = false;
  let selectedUser: User | null = null;

  // Modal handlers
  function openAddModal() {
    isAddModalOpen = true;
  }

  function openEditModal(user: User) {
    selectedUser = user;
    isEditModalOpen = true;
  }

  function openSSHKeysModal(user: User) {
    selectedUser = user;
    isSSHKeysModalOpen = true;
  }

  function openDeleteModal(user: User) {
    selectedUser = user;
    isDeleteModalOpen = true;
  }

  function closeModals() {
    isAddModalOpen = false;
    isEditModalOpen = false;
    isSSHKeysModalOpen = false;
    isDeleteModalOpen = false;
  }

  function handleAddUser(userData: { username: string, email: string, role: string, password: string }) {
    // In a real app, this would make an API call
    const newUser: User = {
      id: users.length + 1,
      username: userData.username,
      email: userData.email,
      role: userData.role
    };
    users = [...users, newUser];
  }

  function handleEditUser(userData: { id: number, username: string, email: string, role: string }) {
    // In a real app, this would make an API call
    users = users.map(user =>
      user.id === userData.id ? { ...userData } : user
    );
  }

  function handleDeleteUser() {
    // In a real app, this would make an API call
    if (selectedUser) {
      users = users.filter(user => user.id !== selectedUser.id);
      closeModals();
    }
  }
</script>

<div>
  <div class="flex justify-between items-center mb-6">
    <h1 class="text-2xl font-bold">Users</h1>
    <button
      class="bg-peach hover:bg-maroon text-white px-4 py-2 rounded"
      on:click={openAddModal}
    >
      New User
    </button>
  </div>

  <div class="bg-surface0 rounded-lg shadow overflow-hidden dark:shadow-lg">
    <table class="min-w-full divide-y divide-overlay0">
      <thead class="bg-surface1 dark:bg-surface2">
        <tr>
          <th class="px-6 py-3 text-left text-xs font-medium text-subtext0 uppercase tracking-wider dark:text-white dark:font-semibold">Username</th>
          <th class="px-6 py-3 text-left text-xs font-medium text-subtext0 uppercase tracking-wider dark:text-white dark:font-semibold">Email</th>
          <th class="px-6 py-3 text-left text-xs font-medium text-subtext0 uppercase tracking-wider dark:text-white dark:font-semibold">Role</th>
          <th class="px-6 py-3 text-left text-xs font-medium text-subtext0 uppercase tracking-wider dark:text-white dark:font-semibold">Actions</th>
        </tr>
      </thead>
      <tbody class="bg-surface0 divide-y divide-overlay0 dark:divide-surface2">
        {#each users as user}
          <tr>
            <td class="px-6 py-4 whitespace-nowrap">
              <div class="text-sm font-medium dark:text-white">{user.username}</div>
            </td>
            <td class="px-6 py-4 whitespace-nowrap">
              <div class="text-sm dark:text-white">{user.email}</div>
            </td>
            <td class="px-6 py-4 whitespace-nowrap">
              <span class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-green bg-opacity-20 text-green dark:bg-green dark:bg-opacity-30 dark:text-white">
                {user.role}
              </span>
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
              <button
                class="text-blue hover:text-sapphire mr-3 dark:text-lavender dark:hover:text-blue"
                on:click={() => openEditModal(user)}
              >
                Edit
              </button>
              <button
                class="text-blue hover:text-sapphire mr-3 dark:text-lavender dark:hover:text-blue"
                on:click={() => openSSHKeysModal(user)}
              >
                SSH Keys
              </button>
              <button
                class="text-red hover:text-maroon dark:text-red dark:hover:text-peach"
                on:click={() => openDeleteModal(user)}
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
<AddUserModal
  isOpen={isAddModalOpen}
  onClose={closeModals}
  onAdd={handleAddUser}
/>

<EditUserModal
  isOpen={isEditModalOpen}
  onClose={closeModals}
  onSave={handleEditUser}
  user={selectedUser}
/>

<SSHKeysModal
  isOpen={isSSHKeysModalOpen}
  onClose={closeModals}
  user={selectedUser}
/>

<DeleteUserModal
  isOpen={isDeleteModalOpen}
  onClose={closeModals}
  onDelete={handleDeleteUser}
  user={selectedUser}
/>
