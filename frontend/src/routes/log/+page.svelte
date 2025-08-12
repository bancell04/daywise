<script lang="ts">
    import { onMount } from 'svelte';
    import { MoveLeft, MoveRight } from 'lucide-svelte';
    import type { Task } from '$lib/types';

    let date : Date
    let title : string
    let category : string
    let start : Date
    let end : Date

    let tasks : Task[] = [];
	let error : string = "";

	onMount(async () => {
		try {
            const now = new Date();
            date = now;
            const formattedDate = now.toISOString().split('T')[0];
            const res = await fetch(`http://localhost:8080/tasks/${formattedDate}`);
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

    function formatDate(date : Date) {
        const today = new Date();

        if (
            date.getDate() === today.getDate() &&
            date.getMonth() === today.getMonth() &&
            date.getFullYear() === today.getFullYear()
        ) {
            return "Today";
        }

        // Format like "Monday Aug 11, 2025"
        return date.toLocaleDateString("en-US", {
            weekday: "long",
            month: "short",
            day: "numeric",
            year: "numeric"
        });
    }

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

	function formatDateTimeLocal(date: Date): string {
	    // Get local date/time in `YYYY-MM-DDTHH:mm` format
		const pad = (n: number) => n.toString().padStart(2, "0");
		return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())}T${pad(date.getHours())}:${pad(date.getMinutes())}`;
	}

	function setStartTime() {
		const now = new Date();
		const input = document.getElementById("start") as HTMLInputElement | null;
		if (input) {
			input.value = formatDateTimeLocal(now);
		}
	}

	function setEndTime() {
		const now = new Date();
		const input = document.getElementById("end") as HTMLInputElement | null;
		if (input) {
			input.value = formatDateTimeLocal(now);
		}
	}

    function incrementDate(direction : number) {
        const newDate = new Date(date);

        newDate.setDate(newDate.getDate() + direction);

        date = newDate;

        console.log(date)
    }

</script>


<div class="min-h-175 w-full flex flex-col justify-center items-center bg-gray-100">
    {#if date}
        <div class="flex flex-row items-center justify-center mb-8">
            <button class="cursor-pointer" on:click={() => incrementDate(-1)}>
                <MoveLeft size={36} color="black" class="mr-8" />
            </button>
            <h1 class="text-5xl">{formatDate(date)}</h1>
            <button class="cursor-pointer" disabled={formatDate(date) === "Today"} on:click={() => incrementDate(1)}>
                <MoveRight size={36} color="black" class="ml-8" />
            </button>
        </div>
    {/if}
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
            <div class="flex flex-row">
                <input id="start" required bind:value={start} type="datetime-local" class="w-full mr-2 px-4 py-2 border rounded-md" />
                <button type="button" on:click={setStartTime}>Now</button>
            </div>
        </div>

        <div class="mb-6">
            <label class="block text-gray-700 text-sm font-semibold mb-2" for="end">End</label>
            <div class="flex flex-row ">
                <input id="end" required bind:value={end} type="datetime-local" class="w-full px-4 py-2 mr-2 border rounded-md" />
                <button type="button" on:click={setEndTime}>Now</button>
            </div>
        </div>

        <button type="submit" class="w-full bg-blue-500 hover:bg-blue-600 text-white font-semibold py-2 px-4 rounded-md">
            Submit
        </button>
  </form>
</div>

