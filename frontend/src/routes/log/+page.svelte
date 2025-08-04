<script lang="ts">
    import type { Task } from '$lib/types';

    let title : string
    let category : string
    let start : Date
    let end : Date

    async function handleTaskSubmit(event) {
        event.preventDefault();
        
        const task : Task = { 
            title,
            category,
            start: new Date(start).toISOString(),
            end: new Date(end).toISOString()
        }

        try {
            const res = await fetch('http://localhost:8080/task', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(task)
            });

            if (!res.ok) {
                throw new Error(`Failed to log task: ${res.status} ${res.statusText}`);
            }
        } catch (error) {
            console.error(error);
            return null;
        }
    }
</script>


<div class="min-h-175 w-full flex flex-col justify-center items-center bg-gray-100">
    <h1 class="mb-2 pb-1 text-7xl font-bold bg-gradient-to-r from-[#7dc4d9] to-[#e1db7f] bg-clip-text text-transparent font-bold">Upload A Task</h1>
    <form on:submit={handleTaskSubmit} class="w-full max-w-md bg-white p-6 rounded-lg shadow">
        <div class="mb-4">
            <label class="block text-gray-700 text-sm font-semibold mb-2" for="title">Title</label>
            <input id="title" required bind:value={title} type="text" class="w-full px-4 py-2 border rounded-md" />
        </div>

        <div class="mb-4">
            <label class="block text-gray-700 text-sm font-semibold mb-2" for="category">Category</label>
            <input id="category" required bind:value={category} type="text" class="w-full px-4 py-2 border rounded-md" />
        </div>

        <div class="mb-4">
            <label class="block text-gray-700 text-sm font-semibold mb-2" for="start">Start</label>
            <input id="start" required bind:value={start} type="datetime-local" class="w-full px-4 py-2 border rounded-md" />
        </div>

        <div class="mb-6">
            <label class="block text-gray-700 text-sm font-semibold mb-2" for="end">End</label>
            <div class="flex flex-row">
                <input id="end" required bind:value={end} type="datetime-local" class="w-full px-4 py-2 border rounded-md" />
                <button type="button" on:click={setEndTime}>Now</button>
            </div>
        </div>

        <button type="submit" class="w-full bg-blue-500 hover:bg-blue-600 text-white font-semibold py-2 px-4 rounded-md">
            Submit
        </button>
  </form>
</div>

