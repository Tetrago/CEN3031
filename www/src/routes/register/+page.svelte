<!-- <script>
    import { BASE_API_PATH } from '$lib/env';
    
    let id = "";
    let pw = "";
    async function get_info(){
        const res = await fetch(`${BASE_API_PATH}/user/register`, {
            // Ordinarily, this second argument could be admitted, but we need CORS so we
            // have to do this bullshit.
            method: 'post',
            mode: 'cors',
            body: JSON.stringify({
                display_name: " ",
                email: id,
                password: pw
            })
        })
    }
</script>


<div class = "p-5 m-5"> 
    <input bind:value={id} type="text" placeholder="Type here" class="input input-bordered w-full max-w-xs" />
</div>
<div class = "p-5 m-5"> 
    <input bind:value={pw} type="text" placeholder="Type here" class="input input-bordered w-full max-w-xs" />
</div>


<button class="btn btn-wide" on:click={get_info}>Sign up</button> -->
<script>
    import { BASE_API_PATH } from '$lib/env';
    import { goto } from '$app/navigation';

    let name = "";
    let id = "";
    let pw1 = "";
    let pw2 = "";
    let errorMessage = "";

    async function get_info(){
        if(pw1 !== pw2){
            errorMessage = "Passwords do not match.";
            return;
        }

        if (!id.endsWith('@ufl.edu')) {
            errorMessage = "Please use a UF email address ending with @ufl.edu.";
            return;
        }

        const res = await fetch(`${BASE_API_PATH}/user/register`, {
            method: 'post',
            mode: 'cors',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                display_name: name, 
                email: id,
                password: pw1
            })
        });

        if (res.ok) {
            goto('/');
        } else {
            errorMessage = "An error occurred. Please try again later.";
        }
    }
</script>

<div class="w-full flex justify-center">
    <div class="flex flex-col card w-1/2 shadow-xl bg-base-200">
        <div class="card-body">
            <p class="text-2xl">Create your account</p>
            <form class="flex flex-col" on:submit={get_info}>
                <div>
                    <div class="label">
                        <div class="label-text">Name</div>
                    </div>
                    <input type="text" class={`w-full input input-bordered ${errorMessage !== "" ? "input-error" : ""}`} bind:value={name} required>
                </div>
                <div>
                    <div class="label">
                        <div class="label-text">Email</div>
                    </div>
                    <input type="email" class={`w-full input input-bordered ${errorMessage !== "" ? "input-error" : ""}`} placeholder="@ufl.edu" bind:value={id} required>
                </div>
                <div>
                    <div class="label">
                        <div class="label-text">Password</div>
                    </div>
                    <input type="password" class={`w-full input input-bordered ${errorMessage !== "" ? "input-error" : ""}`} bind:value={pw1} required>
                </div>
                <div>
                    <div class="label">
                        <div class="label-text">Confirm password</div>
                    </div>
                    <input type="password" class={`w-full input input-bordered ${errorMessage !== "" ? "input-error" : ""}`} bind:value={pw2} required>
                </div>
                <button class="btn btn-primary mt-5" type="submit">Create account</button>
            </form>
        </div>
    </div>
</div>

{#if errorMessage !== ""}
    <div class="toast toast-end">
        <div class="alert alert-error">
            <span>{errorMessage}</span>
        </div>
    </div>
{/if}