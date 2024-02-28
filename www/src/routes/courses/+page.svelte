<script>
    import { BASE_API_PATH } from '$lib/env';

    /** @type {import('./$types').PageData} */
    export let data;

    let departmentField = "";
    let codeField = "";
    let groups = data.groups;

    $: coursesPromise = fetchCourses(departmentField);

    /**
	 * @param {string} dep
	 */
    async function fetchCourses(dep) {
        if(dep.match(/^[A-Za-z]{3}$/) === null) {
            return [];
        }

        const res = await fetch(`${BASE_API_PATH}/course/department/${dep.toUpperCase()}`, {
            method: 'get',
            mode: 'cors'
        })
        return await res.json();
    }

    async function fetchGroups() {
        const res = await fetch(`${BASE_API_PATH}/user/groups`, {
            method: 'get',
            mode: 'cors',
            credentials: 'include'
        });

        return await res.json();
    }

    /**
     * @param {string} label
     */
    async function join(label) {
        const parts = label.split(' ')

        const res = await fetch(`${BASE_API_PATH}/course/group/${parts[0]}/${parts[1]}`)
        const data = await res.json()

        await fetch(`${BASE_API_PATH}/user/join`, {
            method: 'post',
            mode: 'cors',
            credentials: 'include',
            body: JSON.stringify({
                group_id: data
            })
        });

        groups = await fetchGroups();
    }

    /**
     * @param {number} id
     */
    async function leave(id) {
        await fetch(`${BASE_API_PATH}/user/leave`, {
            method: 'post',
            mode: 'cors',
            credentials: 'include',
            body: JSON.stringify({
                group_id: id
            })
        })

        groups = await fetchGroups();
    }
</script>

<div class="join flex justify-center">
    <input type="text" bind:value={departmentField} maxlength="3" class="join-item input input-bordered w-full max-w-40 text-right" placeholder="CEN" />
    <input type="text" bind:value={codeField} maxlength="5" class="join-item input input-bordered w-full max-w-80" placeholder="3031" />
</div>

{#await coursesPromise}
    <progress class="progress w-full"></progress>
{:then courses}
    {#if courses.length == 0}
        <div role="alert" class="mt-5 alert">
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="stroke-info shrink-0 w-6 h-6"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
            <span>It's pretty empty in here</span>
        </div>
    {:else}
        {#each courses.filter(({ label }) => label.includes(codeField, 4)) as { label, name }, i}
            <div class="mt-5 card bg-base-200 shadow-xl cursor-default select-none">
                <div class="card-body">
                    <h2 class="card-title">{label}</h2>
                    <p>{name}</p>

                    {#if groups !== undefined}
                        {#if groups.some(d => d.name === label)}
                            <div class="absolute inset-y-0 right-5 flex items-center h-full">
                                <button on:click={() => leave(groups.find(d => d.name === label).group_id)} class="btn btn-neutral rounded-full w-32">Leave</button>
                            </div>
                        {:else}
                            <div class="absolute inset-y-0 right-5 flex items-center h-full">
                                <button on:click={() => join(label)} class="btn btn-outline rounded-full w-32">Join</button>
                            </div>
                        {/if}
                    {/if}
                </div>
            </div>
        {/each}
    {/if}
{:catch}
    <div role="alert" class="mt-5 alert alert-error">
        <svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
        <span>Something went wrong</span>
    </div>
{/await}