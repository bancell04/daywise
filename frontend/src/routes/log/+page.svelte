<script lang="ts">
    import { onMount } from 'svelte';
    import { MoveLeft, MoveRight, Trash } from 'lucide-svelte';
    import type { Task, Category } from '$lib/types';

    let date : Date
    let formTask : Task = { title: "", category: -1 }
    let currTaskId : number
    let categories : Category[] = [];
    let tasks : Task[] = [];
	let error : string = "";

    let currTime: Date = new Date();

    function startClock(): ReturnType<typeof setInterval> {
        const id = setInterval(() => {
            currTime = new Date();
        }, 2000);
        return id;
    }

    function stopClock(id: ReturnType<typeof setInterval>): void {
        clearInterval(id);
    }

    const clockId = startClock();

	onMount(() => {
        const now = new Date();
        date = now;
        fetchLogsByDate(date);
        fetchUserCategories();

        window.addEventListener("beforeunload", () => stopClock(clockId));

        return () => {
            stopClock(clockId);
            window.removeEventListener("beforeunload", () => stopClock(clockId));
        };
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

    async function fetchLogsByDate(date : Date) {
        try {
            const start = new Date(date);
            start.setHours(0, 0, 0, 0);

            const end = new Date(date);
            end.setHours(23, 59, 59, 999);

            const startIso = start.toISOString();
            const endIso = end.toISOString();

            const res = await fetch(`http://localhost:8080/tasks/${startIso}/to/${endIso}`);
            console.log('http://localhost:8080/tasks/' + startIso + "/to/" + endIso);
            tasks = await res.json();

			if (!res.ok) throw new Error('Failed to fetch tasks');
        } catch (err) {
            if (err instanceof Error) {
                error = err.message;
            } else {
                error = 'An unexpected error occurred';
            }
			console.error(err); 
        }
    }

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

    async function handleTaskSubmit(event : any) {
        event.preventDefault();

        const task : Task = { 
            id: formTask.id,
            title: formTask.title,
            category: formTask.category,
            start: new Date(formTask.start!).toISOString(),
            end: new Date(formTask.end!).toISOString()
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
        fetchLogsByDate(date);
    }

	function formatDateTimeLocal(date: Date): string {
	    // get local date/time in `YYYY-MM-DDTHH:mm` format
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
        fetchLogsByDate(date);
    }


    // interval in minutes
    let interval = 15
    let numIntervals = 1440 / interval

    function formatTimelineSlot(i: number): string {
        const minutesFromMidnight = i * interval;
        const hours24 = Math.floor(minutesFromMidnight / 60);
        const minutes = minutesFromMidnight % 60;

        const hours12 = hours24 % 12 || 12; 
        const ampm = hours24 < 12 ? "AM" : "PM";

        const hoursStr = String(hours12).padStart(2, "0");
        const minutesStr = String(minutes).padStart(2, "0");

        return `${hoursStr}:${minutesStr} ${ampm}`;
    }

    function getTaskHeightProperty(task: Task): string {
        let startDate = new Date(task.start!)
        let taskStartMinutes = startDate.getHours() * 60 + startDate.getMinutes()
        if (task.end != null) {
            let endDate = new Date(task.end)
            let taskEndMinutes = endDate.getHours() * 60 + endDate.getMinutes()
            let durationIntervals = (taskEndMinutes - taskStartMinutes) / interval

            // each slot height is 2rem
            return (durationIntervals * 2).toString()
        } else {
            // return start to NOW.
            return ""
        }
    }

    function getTaskTopProperty(start: Date): string {
        let minutesFromMidnight = (start.getHours() * 60) + start.getMinutes() + (start.getSeconds() / 60)

        let intervalsFromMidnight = minutesFromMidnight / interval
        // each interval height is 2rem
        return (intervalsFromMidnight * 2).toString()
    }

    function autofillTaskData(task: Task) {
        currTaskId = task.id!
        const titleInput = document.getElementById("title") as HTMLInputElement | null;
        const categoryInput = document.getElementById("category") as HTMLInputElement | null;
        const startInput = document.getElementById("start") as HTMLInputElement | null;
        const endInput = document.getElementById("end") as HTMLInputElement | null;

        if (task.id) {
            formTask.id = task.id
        }

        titleInput!.value = task.title
        formTask.title = task.title
        const category = categories.filter(c => c.id === task.category)[0] || null;
        categoryInput!.value = category ? category!.id!.toString() : "";
        formTask.category = category ? category!.id! : 0
        startInput!.value = toDateTimeLocalString(task.start!);
        formTask.start = toDateTimeLocalString(task.start!)
        endInput!.value = toDateTimeLocalString(task.end!);
        formTask.end = toDateTimeLocalString(task.end!)
    }

    function toDateTimeLocalString(dateString: string): string {
        const date = new Date(dateString);
        const pad = (num: number) => num.toString().padStart(2, "0");

        const year = date.getFullYear();
        const month = pad(date.getMonth() + 1);
        const day = pad(date.getDate());
        const hours = pad(date.getHours());
        const minutes = pad(date.getMinutes());

        return `${year}-${month}-${day}T${hours}:${minutes}`;
    }

    function hexToRgba(hex: string, alpha: number) {
        if (hex != null) {
            const r = parseInt(hex.slice(1, 3), 16);
            const g = parseInt(hex.slice(3, 5), 16);
            const b = parseInt(hex.slice(5, 7), 16);
            return `rgba(${r}, ${g}, ${b}, ${alpha})`;
        }
    }

    async function deleteTask(task : Task) {
        if (task.id) {
            try {
                const res = await fetch(`http://localhost:8080/task/${task.id}`, {
                    method: 'DELETE'
                });

                if (!res.ok) {
                    throw new Error(`Failed to delete task: ${res.status} ${res.statusText}`);
                }

                fetchLogsByDate(date);
            } catch (error) {
                console.error(error);
                return null;
            }
        } else {
            tasks = tasks.filter(t => 
                !(t.title === formTask.title && t.start === formTask.start && t.end === formTask.end)
            );
        }

        formTask = {
            title: "",
            category: -1,
            start: undefined,
            end: undefined
        }
    }

    async function createNewTask(i : number) {
        await fetchLogsByDate(date);

        const minutesAfterMidnight = interval * i;
        let startDate = new Date(date);
        startDate.setHours(Math.floor(minutesAfterMidnight / 60));
        startDate.setMinutes(Math.floor(minutesAfterMidnight % 60));

        let task : Task = {
            title : "",
            category : -1,
            start : formatDateTimeLocal(startDate),
            end : formatDateTimeLocal(new Date(startDate.getTime() + interval * 60 * 1000))
        }
        tasks = [...tasks, task];
        formTask = task;
    }

    function reRenderTasks() {
        tasks = [...tasks]
    }
</script>


<div class="min-h-175 w-full flex flex-col items-center bg-gray-100">
    <p>{currTime}</p>
    {#if date}
        <div class="flex flex-row items-center justify-between min-w-175 pt-3 mb-4">
            <button class="cursor-pointer" on:click={() => incrementDate(-1)}>
                <MoveLeft size={36} color="black" class="mr-8" />
            </button>
            <h1 class="text-5xl">{formatDate(date)}</h1>
            <button class="cursor-pointer" disabled={formatDate(date) === "Today"} on:click={() => incrementDate(1)}>
                <MoveRight size={36} color="black" class="ml-8" />
            </button>
        </div>

        <div class="flex flex-row items-center align-">
            <div class="relative max-h-[36rem] overflow-y-auto rounded-lg border border-gray-300">
            {#each {length: numIntervals} as _, i}
                <div
                role="button"
                tabindex="0"
                on:keydown={(e) => e.key === 'Enter' && createNewTask(i)}
                class="flex items-center border border-gray-300 h-[2rem] w-200 px-2 z-3"
                class:bg-gray-50={i % 2 === 0}
                class:bg-gray-100={i % 2 !== 0}
                on:click={() => createNewTask(i)}
                >
                <span class="text-sm font-medium">{formatTimelineSlot(i)}</span>
                </div>
            {/each}

            {#if tasks && tasks.length > 0}
                {#each tasks as task}

                <div
                role="button"
                tabindex="0"
                class="flex justify-center items-center absolute border left-[8rem] rounded-lg w-150 bg-white z-5"
                style={`
                    top: ${getTaskTopProperty(new Date(task.start!))}rem;
                    height: ${getTaskHeightProperty(task)}rem;
                    border-color: ${categories.find(c => c.id === task.category)?.color};
                    background-color: ${hexToRgba(categories.find(c => c.id === task.category)?.color!, 0.15)};
                `}
                on:click={() => autofillTaskData(task)}
                on:keydown={(e) => e.key === 'Enter' && autofillTaskData(task)}
                aria-label="task"
                >
                {task.title}
                </div>
                {/each}
            {/if}
                
                <div 
                class="absolute rounded-lg w-full bg-[#7dc4d9] h-[4px]"
                style:top={`${getTaskTopProperty(currTime)}rem`}
                >
                </div>
            </div>

            <div class="ml-8">
                <form on:submit={handleTaskSubmit} class="w-full max-w-md bg-white p-6 rounded-lg shadow">
                    <div class="mb-4">
                        <label class="block text-gray-700 text-sm font-semibold mb-2" for="title">Title</label>
                        <input id="title" placeholder="Task Title" required on:change={reRenderTasks} on:input={reRenderTasks} bind:value={formTask.title} type="text" class="w-full px-4 py-2 border rounded-md" />
                    </div>

                    <div class="mb-4">
                        <label class="block text-gray-700 text-sm font-semibold mb-2" for="category">Category</label>
                        <select id="category"  on:change={reRenderTasks} bind:value={formTask.category} required class="w-full px-4 py-2 border rounded-md">
                            <option value="" disabled selected>Select a category</option>
                            {#each categories as c}
                                <option value={c.id}>{c.name}</option>
                            {/each}
                        </select>
                    </div>

                    <div class="mb-4">
                        <label class="block text-gray-700 text-sm font-semibold mb-2" for="start">Start</label>
                        <div class="flex flex-row">
                            <input id="start" required  on:change={reRenderTasks} bind:value={formTask.start} type="datetime-local" class="w-full mr-2 px-4 py-2 border rounded-md" />
                            <button type="button" on:click={setStartTime}>Now</button>
                        </div>
                    </div>

                    <div class="mb-6">
                        <label class="block text-gray-700 text-sm font-semibold mb-2" for="end">End</label>
                        <div class="flex flex-row ">
                            <input id="end" required  on:change={reRenderTasks} bind:value={formTask.end} type="datetime-local" class="w-full px-4 py-2 mr-2 border rounded-md" />
                            <button type="button" on:click={setEndTime}>Now</button>
                        </div>
                    </div>

                    <div class="flex flex-row">
                        <button type="submit" disabled={formTask.title.length === 0 && formTask.category === -1} class="w-full bg-blue-500 hover:bg-blue-600 text-white font-semibold py-2 px-4 rounded-md">
                            Submit
                        </button>
                        <button type="button" class="cursor-pointer py-2 ml-2" disabled={formTask.title.length === 0 && formTask.category === -1} on:click={() => deleteTask(formTask)}>
                            <Trash/>
                        </button>
                    </div>
                </form>
            </div>
        </div>
    {/if}
</div>

