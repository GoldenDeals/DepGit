<script lang="ts">
  import Router from "svelte-spa-router";
  import Navbar from "./components/Navbar.svelte";
  import Home from "./pages/Home.svelte";
  import Repositories from "./pages/Repositories.svelte";
  import Users from "./pages/Users.svelte";

  // Import base styles and theme CSS with integrated contrast enhancements
  import "./app.css";
  import "./theme.css";

  // Define routes
  const routes = {
    '/': Home,
    '/repositories': Repositories,
    '/users': Users
  };
</script>

<div class="min-h-screen bg-background text-text">
  <Navbar />
  <main class="container mx-auto px-4 py-8">
    <Router {routes} />
  </main>
  <footer class="bg-surface p-4 shadow-inner mt-8">
    <div class="container mx-auto text-center text-text-muted">
      <p>DepGit &copy; {new Date().getFullYear()}</p>
    </div>
  </footer>
</div>
