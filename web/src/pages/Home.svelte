<script lang="ts">
  import { push } from "svelte-spa-router";
</script>

<div class="text-center py-12">
  <h1 class="text-4xl font-bold mb-4">Welcome to DepGit</h1>
  <p class="text-xl text-text-muted mb-8">A modern Git repository management system</p>

  <div class="flex justify-center space-x-4">
    <button
      class="bg-primary hover:bg-primary-hover text-white px-6 py-3 rounded-lg shadow-md"
      on:click={() => push('/repositories')}
    >
      Browse Repositories
    </button>
    <button
      class="bg-secondary hover:bg-secondary-hover text-white px-6 py-3 rounded-lg shadow-md"
      on:click={() => push('/users')}
    >
      Manage Users
    </button>
  </div>

  <div class="mt-16 grid grid-cols-1 md:grid-cols-3 gap-8">
    <div class="bg-surface p-6 rounded-lg shadow-md">
      <h2 class="text-xl font-semibold mb-4">Git Repository Hosting</h2>
      <p class="text-text-muted">Host your Git repositories with secure SSH access and fine-grained permissions.</p>
    </div>
    <div class="bg-surface p-6 rounded-lg shadow-md">
      <h2 class="text-xl font-semibold mb-4">User Management</h2>
      <p class="text-text-muted">Create and manage users with different access levels and SSH key authentication.</p>
    </div>
    <div class="bg-surface p-6 rounded-lg shadow-md">
      <h2 class="text-xl font-semibold mb-4">Access Control</h2>
      <p class="text-text-muted">Control repository access with customizable roles and permissions.</p>
    </div>
  </div>
</div>
