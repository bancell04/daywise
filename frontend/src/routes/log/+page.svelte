<script lang="ts">
    import { onMount } from 'svelte';
    import { MoveLeft, MoveRight } from 'lucide-svelte';
    import type { Task } from '$lib/types';

    let date : Date
    let title : string = ""
    let category : string = ""
    let start : Date
    let end : Date

    let tasks : Task[] = [];
	let error : string = "";

	onMount(async () => {
        const now = new Date();
        date = now;
        fetchLogsByDate(date);
	});

    async function fetchLogsByDate(date : Date) {
        try {
            const formattedDate = date.toISOString().split('T')[0];
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
        let startDate = new Date(task.start)
        let taskStartMinutes = startDate.getHours() * 60 + startDate.getMinutes()
        if (task.end != null) {
            let endDate = new Date(task.end)
            let taskEndMinutes = endDate.getHours() * 60 + endDate.getMinutes()
            console.log("Start Minutes: " + taskStartMinutes)
            console.log("End Minutes: " + taskEndMinutes)
            let durationIntervals = (taskEndMinutes - taskStartMinutes) / interval
            // each slot height is 2rem
            console.log("Height: " + durationIntervals * 2 + "rem")
            return (durationIntervals * 2).toString()
        } else {
            // return start to NOW.
            return ""
        }
    }

    function getTaskTopProperty(task: Task): string {
        let startDate = new Date(task.start)
        let minutesFromMidnight = (startDate.getHours() * 60) + startDate.getMinutes()

        let intervalsFromMidnight = minutesFromMidnight / interval
        // each interval height is 2rem
        return (intervalsFromMidnight * 2).toString()
    }
</script>


<div class="min-h-175 w-full flex flex-col items-center bg-gray-100">
    <h1 class="mb-8 pb-2 text-7xl font-bold bg-gradient-to-r from-[#7dc4d9] to-[#e1db7f] bg-clip-text text-transparent font-bold">Log</h1>
    {#if date}
        <div class="flex flex-row items-center justify-center mb-4">
            <button class="cursor-pointer" on:click={() => incrementDate(-1)}>
                <MoveLeft size={36} color="black" class="mr-8" />
            </button>
            <h1 class="text-5xl">{formatDate(date)}</h1>
            <button class="cursor-pointer" disabled={formatDate(date) === "Today"} on:click={() => incrementDate(1)}>
                <MoveRight size={36} color="black" class="ml-8" />
            </button>
        </div>


        <div class="timeline relative">
            {#each {length: numIntervals}, i}
                <div class="flex items-center border border-gray-300 h-[2rem] w-200 px-2">
                    <span class="text-sm font-medium">{formatTimelineSlot(i)}</span>
                </div>
            {/each}

            {#if tasks && tasks.length > 0}
                {#each tasks as task}
                    <div class="absolute border border-red-300 left-[8rem] w-150" style={`top: ${getTaskTopProperty(task)}rem; height: ${getTaskHeightProperty(task)}rem;`}>
                        {task.title}
                    </div>
                {/each}
            {/if}
        </div>
    {/if}
</div>

