<script lang="ts">
  import { link } from "svelte-spa-router";
  import LoginModal from "./auth/LoginModal.svelte";
  import RegisterModal from "./auth/RegisterModal.svelte";
  import ForgotPasswordModal from "./auth/ForgotPasswordModal.svelte";
  import ResetPasswordModal from "./auth/ResetPasswordModal.svelte";
  import ThemeToggle from "./ThemeToggle.svelte";

  // Define user type
  type User = {
    username: string;
  };

  // Auth state
  let isLoggedIn = false;
  let currentUser: User | null = null;
  let isLoginModalOpen = false;
  let isRegisterModalOpen = false;
  let isForgotPasswordModalOpen = false;
  let isResetPasswordModalOpen = false;
  let resetToken = 'demo-token-123'; // In a real app, this would come from URL params

  function openLoginModal() {
    isLoginModalOpen = true;
  }

  function openRegisterModal() {
    isRegisterModalOpen = true;
  }

  function openForgotPasswordModal() {
    isForgotPasswordModalOpen = true;
  }

  function openResetPasswordModal() {
    isResetPasswordModalOpen = true;
  }

  function closeModals() {
    isLoginModalOpen = false;
    isRegisterModalOpen = false;
    isForgotPasswordModalOpen = false;
    isResetPasswordModalOpen = false;
  }

  function handleLogin(credentials: { username: string, password: string }) {
    console.log('Login with:', credentials);
    // In a real app, this would make an API call
    isLoggedIn = true;
    currentUser = { username: credentials.username };
  }

  function handleRegister(userData: { username: string, email: string, password: string }) {
    console.log('Register with:', userData);
    // In a real app, this would make an API call
    isLoggedIn = true;
    currentUser = { username: userData.username };
  }

  function handlePasswordReset(email: string) {
    console.log('Password reset requested for:', email);
    // In a real app, this would make an API call to send a reset email

    // For demo purposes, we'll open the reset password modal
    // In a real app, this would be handled via a URL with a token
    setTimeout(() => {
      openResetPasswordModal();
    }, 2000);
  }

  function handleResetPassword(data: { token: string, password: string }) {
    console.log('Reset password with token:', data.token, 'and new password');
    // In a real app, this would make an API call to reset the password
  }

  function handleLogout() {
    // In a real app, this would make an API call
    isLoggedIn = false;
    currentUser = null;
  }
</script>

<nav class="bg-primary text-white shadow-md">
  <div class="container mx-auto px-4">
    <div class="flex justify-between items-center py-4">
      <div class="flex items-center space-x-4">
        <a href="/" use:link class="text-xl font-bold">DepGit</a>
        <div class="hidden md:flex space-x-4">
          <a href="/repositories" use:link class="hover:text-primary-light">Repositories</a>
          <a href="/users" use:link class="hover:text-primary-light">Users</a>
        </div>
      </div>
      <div class="flex items-center space-x-2">
        <ThemeToggle />
        {#if isLoggedIn}
          <div class="flex items-center space-x-4">
            <span class="text-sm">Welcome, {currentUser?.username || 'User'}</span>
            <button
              class="bg-primary-hover hover:bg-primary-light px-4 py-2 rounded text-white"
              on:click={handleLogout}
            >
              Logout
            </button>
          </div>
        {:else}
          <div class="flex items-center space-x-2">
            <button
              class="text-white hover:text-primary-light px-3 py-2"
              on:click={openLoginModal}
            >
              Log in
            </button>
            <button
              class="bg-primary-hover hover:bg-primary-light px-4 py-2 rounded text-white"
              on:click={openRegisterModal}
            >
              Sign up
            </button>
          </div>
        {/if}
      </div>
    </div>
  </div>
</nav>

<!-- Auth Modals -->
<LoginModal
  isOpen={isLoginModalOpen}
  onClose={closeModals}
  onLogin={handleLogin}
  onRegisterClick={openRegisterModal}
  onForgotPasswordClick={openForgotPasswordModal}
/>

<RegisterModal
  isOpen={isRegisterModalOpen}
  onClose={closeModals}
  onRegister={handleRegister}
  onLoginClick={openLoginModal}
/>

<ForgotPasswordModal
  isOpen={isForgotPasswordModalOpen}
  onClose={closeModals}
  onRequestReset={handlePasswordReset}
  onLoginClick={openLoginModal}
/>

<ResetPasswordModal
  isOpen={isResetPasswordModalOpen}
  onClose={closeModals}
  onResetPassword={handleResetPassword}
  onLoginClick={openLoginModal}
  token={resetToken}
/>
