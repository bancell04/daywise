<script lang="ts">
	import { onMount } from 'svelte';
    import type { Task } from '$lib/types';

	let tasks : Task[] = [];
	let error : string = "";

	onMount(async () => {
		try {
			const res = await fetch('http://localhost:8080/tasks');
			if (!res.ok) throw new Error('Failed to fetch tasks');
			tasks = await res.json();
		} catch (err) {
            if (err instanceof Error) {
                error = err.message;
            } else {
                error = 'An unexpected error occurred';
            }
			console.error(err);
		}
	});
</script>

<h1>
    History
</h1>
{#if error}
    <p class="error">{error}</p>
{:else if tasks.length === 0}
    <p>No tasks found.</p>
{:else}
    <ul>
        {#each tasks as task}
            <li>
				<strong>{task.title}</strong> - {task.category} <br />
				{task.start ? new Date(task.start).toLocaleString() : 'No start'} -
				{task.end ? new Date(task.end).toLocaleString() : 'No end'}
            </li>
        {/each}
    </ul>
{/if}