<script>
    import { BASE_API_PATH } from '$lib/env';

    // This 'data' variable contains the data we fetched in page.server.js
    // data.groups contains the groups the API returned (or undefined if
    // we aren't logged in)
    /** @type {import('./$types').PageData} */
    export let data;

    // These variables are bound to input text boxes later in this file
    let departmentField = "";
    let codeField = "";

    // This grabs the groups we queryed from the backend and stores them locally.
    // They'll be changed on the client side when we join and leave groups, and the
    // data variable shouldn't be changed.
    let groups = data.groups;

    // The $: syntax is special to Svelte. In this situation, it tells Svelte "hey,
    // if the departmentField variable ever changes, update coursesPromise too".
    // If this was a regular let statement, nothing would happen when the user
    // edits the course department search bar on the website, which would suck.
    $: coursesPromise = fetchCourses(departmentField);

    /**
     * Fetches all courses at UF with the given 3 digit department code
     * 
	 * @param {string} dep
	 */
    async function fetchCourses(dep) {
        // This is a RegEx statement that makes sure the department code is valid.
        // This keeps the client from querying the backend over values it knows are invalid.
        if(dep.match(/^[A-Za-z]{3}$/) === null) {
            return [];
        }

        // Ask the backend for all UF courses with the given department code.
        const res = await fetch(`${BASE_API_PATH}/course/department/${dep.toUpperCase()}`, {
            // Ordinarily, this second argument could be admitted, but we need CORS so we
            // have to do this bullshit.
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
     * Joins the given group label (Ex. "COP 4600")
     *
     * @param {string} label
     */
    async function join(label) {
        // parts[0] = "COP"
        // parts[1] = "4600"
        const parts = label.split(' ')

        // /course/group/COP/4600
        // This gets the ID of the COP 4600 group
        const res = await fetch(`${BASE_API_PATH}/course/group/${parts[0]}/${parts[1]}`)
        const data = await res.json()

        await fetch(`${BASE_API_PATH}/user/join`, {
            // Here we need to make a POST request instead of a GET request. The reason
            // being that GET requests can't send data, and that wouldn't make sense anyways.
            // POST requests are used to tell a server to do something. Here we're telling
            // the server to join a group and specifying the group id.
            method: 'post',
            mode: 'cors',
            credentials: 'include',
            
            // The backend debug webpage will show you what arguments each POST request takes
            body: JSON.stringify({
                group_id: data
            })
        });

        // Update the page with the new list of groups
        groups = await fetchGroups();
    }

    /**
     * Leaves the given group by it's ID
     *
     * @param {number} id
     */
    async function leave(id) {
        // Since we have a list of which groups we're in, we know what the group ID is,
        // so there's no point in querying the server for the ID of the course.

        await fetch(`${BASE_API_PATH}/user/leave`, {
            method: 'post',
            mode: 'cors',
            credentials: 'include',
            body: JSON.stringify({
                group_id: id
            })
        })

        // Update the page with the new list of groups
        groups = await fetchGroups();
    }
</script>

<!-- These classes are part of TailwindCSS and daisyUI (look at the README) -->
<div class="join flex justify-center">
    <!-- The bind:value tag binds the value of the input field to the variables we specified above -->
    <input type="text" bind:value={departmentField} maxlength="3" class="join-item input input-bordered w-full max-w-40 text-right" placeholder="CEN" />
    <input type="text" bind:value={codeField} maxlength="5" class="join-item input input-bordered w-full max-w-80" placeholder="3031" />
</div>

<!--
    This is special Svelte syntax that is used to show loading bars while fetching information.
    Remeber in the +page.server.js file where I mentioned why async and await exist? This allows
    us to provide dummy visuals (a loading bar) until the information we requested is available.
-->
{#await coursesPromise}
    <progress class="progress w-full"></progress>
{:then courses}
    <!-- Once the information we requested is available, it's stored in the variable 'courses', as we specified. -->

    <!--
        Special Svelte thing again. Here, we're displaying a special message if there aren't any groups.
        If the page just showed up as being blank, most people would interpret that as an error of some kind.
    -->
    {#if courses.length == 0}
        <div role="alert" class="mt-5 alert">
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="stroke-info shrink-0 w-6 h-6"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
            <span>It's pretty empty in here</span>
        </div>
    {:else}
        <!--
            This is a for loop in svelte that adds one of those div blocks for every course available.
            We also filter the list with all courses that match the course code input box.
        -->
        {#each courses.filter(({ label }) => label.includes(codeField, 4)) as { label, name }, i}
            <!-- The HTML is taking straight from daisyUI's page on cards -->
            <div class="mt-5 card bg-base-200 shadow-xl cursor-default select-none">
                <div class="card-body">
                    <!-- The {label} here inserts the course label into the HTML -->
                    <h2 class="card-title">{label}</h2>
                    <p>{name}</p>

                    <!-- This section checks if a user is logged in. If so, then we want to display Join/Leave buttons for each course -->
                    {#if groups !== undefined}
                        <!-- Crude check to see if the user is already in this group -->
                        {#if groups.some(d => d.name === label)}
                            <div class="absolute inset-y-0 right-5 flex items-center h-full">
                                <!--
                                    The on:click attribute takes a *function* to call when the button is clicked.
                                    That's why the JavaScript labmda syntax "() =>" is used instead of a regular statement
                                -->
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
    <!--
        If for some reason the request failed, we want to let the user know that something went wrong.
        One of those "good practice" things since this likely won't happen a lot in our situation.
    -->

    <div role="alert" class="mt-5 alert alert-error">
        <svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
        <span>Something went wrong</span>
    </div>
{/await}