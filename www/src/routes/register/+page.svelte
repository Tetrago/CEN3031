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

        const data = await res.json();
        if (res.ok) {
            errorMessage = "";
            goto('/home');
        } else {
            errorMessage = data.message || "An error occurred. Please try again later.";
        }
    }
</script>

<div class="container">
    <h2>Create Your Account</h2>
    <form on:submit={get_info}>
        <div class="form-group">
            <label for="name">Name</label>
            <input type="text" id="name" bind:value={name} required>
        </div>
        <div class="form-group">
            <label for="email">Email Address (Must be UF Email)</label>
            <input type="email" id="email" bind:value={id} required>
        </div>
        <div class="form-group">
            <label for="password">Password</label>
            <input type="password" id="password" bind:value={pw1} required>
        </div>
        <div class="form-group">
            <label for="confirmPassword">Confirm Password</label>
            <input type="password" id="confirmPassword" bind:value={pw2} required>
        </div>
        <button type="submit">Create account</button>
    </form>
    {#if errorMessage}
    <p class="error">{errorMessage}</p>
    {/if}
</div>


<style>
    .container {
        max-width: 400px;
        margin: 50px auto;
        padding: 20px;
        border: 1px solid #ccc;
        border-radius: 5px;
    }
    .form-group {
        margin-bottom: 15px;
    }
    label {
        display: block;
        font-weight: bold;
    }
    h2{
        font-size: 24px;
        font-weight: bold;
        margin-bottom: 30px;
    }
    input[type="text"],
    input[type="email"],
    input[type="password"] {
        width: 100%;
        padding: 10px;
        border: 1px solid #ccc;
        border-radius: 5px;
    }
    button {
        display: block;
        width: 100%;
        padding: 10px;
        background-color: #007bff;
        color: #fff;
        border: none;
        border-radius: 5px;
        cursor: pointer;
        margin-top: 30px;
    }
    button:hover {
        background-color: #0056b3;
    }
    .error {
        color: red;
        margin-bottom: 10px;
    }
</style>