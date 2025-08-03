<script lang="ts">
    const secret = ""
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
        } catch (err) {
            console.error('Reset error', err);
        }
    }
</script>
<h1>
    Profile
</h1>
<button on:click={resetDatabase}>
    Reset Database
</button>