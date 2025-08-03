<script lang="ts">
    import type { Task } from '$lib/types';

    async function logTask(task: Task): Promise<any> {
        try {
            const res = await fetch('http://localhost:8080/task/', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(task)
            });

            if (!res.ok) {
                throw new Error(`Failed to log task: ${res.status} ${res.statusText}`);
            }

            return await res.json();
        } catch (error) {
            console.error(error);
            return null;
        }
    }

    function handleLogUpload() {
        const exampleTask: Task = {
            title: "Example Task",
            category: "Test Category",
            start: new Date().toISOString(),
            end: new Date().toISOString()
        }
        logTask(exampleTask);
    }
</script>
<h1>
    Log
</h1>

<button aria-label="log task button" on:click={handleLogUpload}>
    Log task here
</button>