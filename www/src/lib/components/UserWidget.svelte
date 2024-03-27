<script>
    import { BASE_API_PATH } from '$lib/env';
    import { user_identifier } from '../../routes/stores';

    /** @type HTMLDialogElement */
    let signInModal;

    $: failed = false;

    let formEmail = "";
    let formPassword = "";

    async function signIn() {
        const res = await fetch('/auth/login', {
            method: 'post',
            body: JSON.stringify({
                email: formEmail,
                password: formPassword
            })
        });

        if(res.ok) {
            const data = await res.json();

            failed = false;
            signInModal.close();
            user_identifier.set(data.ident);
        } else {
            failed = true;
        }
    }

    async function signOut() {
        await fetch('/auth/logout', {
            method: 'post'
        });

        user_identifier.set("");
    }
</script>

{#if $user_identifier !== ""}
    <div class="dropdown dropdown-bottom dropdown-end">
        <div tabindex="0" class="btn btn-ghost btn-circle avatar" role="button">
            <div class="rounded-full w-10">
                <img alt="Avatar" src={`${BASE_API_PATH}/user/profile_picture/${$user_identifier}`} />
            </div>
        </div>
        <ul tabindex="0" class="dropdown-content z-[1] menu p-2 shadow bg-base-300 rounded-box w-52">
            <li><button on:click={signOut}>Sign out</button></li>
        </ul>
    </div>
{:else}
    <button on:click={() => signInModal.showModal()} class="btn rounded-full font-bold">Sign in</button>
{/if}

<dialog bind:this={signInModal} id="modal_sign_in" class="modal modal-bottom sm:modal-middle">
    <div class="modal-box">
        <form method="dialog">
            <button class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2">✕</button>
        </form>
        <h3 class="font-bold text-lg">Sign in</h3>
        <label class={`input input-bordered flex items-center gap-2 mt-2 ${failed ? "input-error" : ""}`}>
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" fill="currentColor" class="w-4 h-4 opacity-70"><path d="M2.5 3A1.5 1.5 0 0 0 1 4.5v.793c.026.009.051.02.076.032L7.674 8.51c.206.1.446.1.652 0l6.598-3.185A.755.755 0 0 1 15 5.293V4.5A1.5 1.5 0 0 0 13.5 3h-11Z" /><path d="M15 6.954 8.978 9.86a2.25 2.25 0 0 1-1.956 0L1 6.954V11.5A1.5 1.5 0 0 0 2.5 13h11a1.5 1.5 0 0 0 1.5-1.5V6.954Z" /></svg>
            <input bind:value={formEmail} type="text" name="email" class="grow border-none focus:ring-0" placeholder="Email" />
        </label>
        <label class={`input input-bordered flex items-center gap-2 mt-2 ${failed ? "input-error" : ""}`}>
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" fill="currentColor" class="w-4 h-4 opacity-70"><path fill-rule="evenodd" d="M14 6a4 4 0 0 1-4.899 3.899l-1.955 1.955a.5.5 0 0 1-.353.146H5v1.5a.5.5 0 0 1-.5.5h-2a.5.5 0 0 1-.5-.5v-2.293a.5.5 0 0 1 .146-.353l3.955-3.955A4 4 0 1 1 14 6Zm-4-2a.75.75 0 0 0 0 1.5.5.5 0 0 1 .5.5.75.75 0 0 0 1.5 0 2 2 0 0 0-2-2Z" clip-rule="evenodd" /></svg>
            <input bind:value={formPassword} type="password" name="password" class="grow border-none focus:ring-0" placeholder="Password" />
        </label>
        <div class="flex justify-end mt-2">
            <button on:click={signIn} class="btn btn-primary">Done</button>
        </div>
    </div>
</dialog>