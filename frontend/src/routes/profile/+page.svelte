<script lang="ts">
    import type { Task, Category } from '$lib/types';
    import { onMount } from 'svelte';
    import { Plus, Trash } from 'lucide-svelte';

    const secret = ""
    var message = ""
    async function resetDatabase(): Promise<any> {
        try {
            const res = await fetch('http://localhost:8080/db-reset', {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                    'X-Admin-Secret': secret
                }
            });
            if (!res.ok) throw new Error('Reset failed');
            message = res.statusText
        } catch (err) {
            console.error('Reset error', err);
        }
    }

    let title : string = ""
    let category : number
    let start : string
    let end : string

    let categories : Category[] = []

    onMount(async () => {
        fetchUserCategories();
	});

    async function fetchUserCategories() {
        try {
            const res = await fetch('http://localhost:8080/categories', {
                method: 'GET',
                headers: { 'Content-Type': 'application/json' },
            });

            if (!res.ok) {
                throw new Error(`Failed to fetch categories: ${res.status} ${res.statusText}`);
            }
            categories = await res.json();

        } catch (error) {
            console.error(error);
            return null;
        }
    }

    async function handleTaskSubmit(event : any) {
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
            console.log("START: " + formatDateTimeLocal(now))
			input.value = formatDateTimeLocal(now);
            start = input.value;
		}
	}

	function setEndTime() {
		const now = new Date();
		const input = document.getElementById("end") as HTMLInputElement | null;
		if (input) {
            console.log("END: " + now.toDateString())
			input.value = formatDateTimeLocal(now);
            end = input.value;
		}
	}

	function addCategory() {
		categories = [...categories, { name: '', color: '#000000' }];
	}

    function removeCategory(category: Category) {
        categories = categories.filter(c => c !== category);
    }

    async function submitCategories() {
        try {
            const res = await fetch('http://localhost:8080/categories', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(categories)
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
<h1>
    Profile
</h1>
<button on:click={resetDatabase}>
    Reset Database
</button>
{#if message}
    <h2 class="error">{message}</h2>
{/if}

<form on:submit={handleTaskSubmit} class="w-full max-w-md bg-white p-6 rounded-lg shadow">
        <div class="mb-4">
            <label class="block text-gray-700 text-sm font-semibold mb-2" for="title">Title</label>
            <input id="title" required bind:value={title} type="text" class="w-full px-4 py-2 border rounded-md" />
        </div>

        <div class="mb-4">
            <label class="block text-gray-700 text-sm font-semibold mb-2" for="category">Category</label>
            <select id="category" bind:value={category} required class="w-full px-4 py-2 border rounded-md">
                <option value="" disabled selected>Select a category</option>
                {#each categories as c}
                    <option value={c.id}>{c.name}</option>
                {/each}
            </select>
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
  
  <div class="border w-[500px] p-4 flex flex-col items-center justify-center">
    <h1 class="text-5xl pb-4">Categories</h1>
    {#if categories && categories.length > 0}
        {#each categories as category}
            <div class="border rounded-lg p-2 min-w-[150px] flex flex-row items-center">
                <input type="text" placeholder="Category Name" bind:value={category.name}>
                <input class="px-4" type="color" bind:value={category.color}>
                <button class="cursor-pointer py-2" on:click={() => removeCategory(category)}>
                    <Trash/>
                </button>
            </div>
        {/each}
    {:else}
        <p>No categories found.</p>
    {/if}
    
    <button class="cursor-pointer py-2" on:click={addCategory}>
        <Plus/>
    </button>

    <button class="w-[150px] bg-blue-500 hover:bg-blue-600 text-white font-semibold py-2 px-4 rounded-md" on:click={submitCategories}>
        Save
    </button>
  </div>
