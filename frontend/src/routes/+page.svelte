<script>
    import { goto } from '$app/navigation';

    const goToLog = () => goto('/log');
    const goToHistory = () => goto('/history');
    const goToVisualize = () => goto('/visualize');
    const goToProfile = () => goto('/profile');
</script>
<h1>Welcome to Daywise</h1>
<div class="flex flex-col ">
    <button aria-label="navigate to log" on:click={goToLog}>
        Log
    </button>
    <button aria-label="navigate to log" on:click={goToHistory}>
        History
    </button>
    <button aria-label="navigate to visualize" on:click={goToVisualize}>
        Visualize
    </button>
    <button aria-label="navigate to profile" on:click={goToProfile}>
        Profile
    </button>
</div>
