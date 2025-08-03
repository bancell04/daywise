<script lang="ts">
	import "../app.css";
    import { goto } from '$app/navigation';
	import { derived } from "svelte/store";
	import { page } from '$app/stores';

	const goToHome = () => goto('/');
    const goToLog = () => goto('/log');;
    const goToHistory = () => goto('/history');
    const goToVisualize = () => goto('/visualize');
    const goToProfile = () => goto('/profile');

  const activePage = derived(page, ($page) => {
    const path = $page.url.pathname;
    if (path.startsWith("/log")) return "log";
    if (path.startsWith("/history")) return "history";
    if (path.startsWith("/visualize")) return "visualize";
    if (path.startsWith("/profile")) return "profile";
    return "";
  });
</script>

<div class="min-h-screen flex flex-col">
	<header class="bg-orange">
		<nav class="h-20 flex items-center px-8 gap-4">
			<a href="/" on:click={goToHome}>
				<img src="/daywise.png" class="h-20" alt="daywise logo">
			</a>
			<button type="button" class="cursor-pointer text-lg font-bold" on:click={goToLog} class:active={$activePage === "log"}>Log</button>
			<button type="button" class="cursor-pointer text-lg font-bold" on:click={goToHistory} class:active={$activePage === "history"}>History</button>
			<button type="button" class="cursor-pointer text-lg font-bold" on:click={goToVisualize} class:active={$activePage === "visualize"}>Visualize</button>
			<button type="button" class="cursor-pointer text-lg font-bold" on:click={goToProfile} class:active={$activePage === "profile"}>Profile</button>
		</nav>
	</header>

	<main class="h-full w-full flex-grow">
		<slot />
	</main>

	<footer class="text-center mt-8 text-sm text-gray-500">
		© {new Date().getFullYear()} Daywise
	</footer>
</div>