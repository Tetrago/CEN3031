<script>
    import { onMount } from 'svelte';
    import { BASE_API_PATH } from '$lib/env';
    import { goto } from '$app/navigation';
    import { user_identifier } from '../stores';


    let profilePicture, fileinput, files;
    let editing = false;
    let name = "";
    let email = "";
    let bio = "";

    let new_name = "";
    let new_email = "";
    let new_bio = "";

    let errorMessage = "";
    let p_errorMessage = "";
    let currentPassword = "";
    let newPassword = "";
    let confirmPassword = "";

    const formData  = new FormData();
    $: if(files){

        formData.append('request', files[0]);
        profilePicture = URL.createObjectURL(files[0]);
    }

	
    profilePicture = `${BASE_API_PATH}/user/profile_picture/${$user_identifier}`;
    onMount(async () => {
        const res_post_info = await fetch(`${BASE_API_PATH}/user/get/${$user_identifier}`, {
            method: 'get',
            mode: 'cors',
            credentials: 'include'
        });

        let temp = await res_post_info.json();
        name = temp.display_name;
        email = temp.email;
        bio = temp.bio;

        new_name = name;
        new_email = email;
        new_bio = bio;
    })

    function edit(){
        errorMessage = "";
        editing = true;
    }
    async function save(){
        errorMessage = "";

        // Check if the entered values are different from the current values
        if (name !== new_name) { // Display Name
            name = new_name;
            const res = await fetch(`${BASE_API_PATH}/user/display_name`, {
                method: 'post',
                mode: 'cors',
                credentials: 'include',
                body: JSON.stringify({
                    display_name: name
                })
            });

            if (!res.ok){
                errorMessage = "Failed to change name";
                return;
            }
        }
        if (email !== new_email) { // Email
            email = new_email;
            const res = await fetch(`${BASE_API_PATH}/user/email`, {
                method: 'post',
                mode: 'cors',
                credentials: 'include',
                body: JSON.stringify({
                    email: email
                })
            });
            if (!res.ok){
                errorMessage = "Failed to change email";
                return;
            }
        }
        if (bio !== new_bio) { // Bio
            bio = new_bio;

            const res = await fetch(`${BASE_API_PATH}/user/bio`, {
                method: 'post',
                mode: 'cors',
                credentials: 'include',
                body: JSON.stringify({
                    bio: bio
                })
            });

            if (!res.ok){
                errorMessage = "Failed to change bio";
                return;
            }
        }

        // Check if the input email ends with '@ufl.edu'
        if (!email.endsWith('@ufl.edu')) {
            errorMessage = "Email must end with '@ufl.edu'";
            return;
        }

    
        if(files){
           const res_post_image = await fetch(`${BASE_API_PATH}/user/profile_picture`, {
            method: 'post',
            mode: 'cors',
            credentials: 'include',
            body: formData
            });

            if(!res_post_image.ok){
                profilePicture = `${BASE_API_PATH}/user/profile_picture/${$user_identifier}`;
                errorMessage = "Please check your image."
            } 
        }
        
        editing = false;
    }

    async function changePassword(){
        if(newPassword !== confirmPassword){
            p_errorMessage = "Passwords do not match.";
            return;
        }

        if (!currentPassword || !newPassword || !confirmPassword) {
            p_errorMessage = "Please fill out all password fields.";
            return;
        }

        const res = await fetch(`${BASE_API_PATH}/user/password`, {
            method: 'post',
            mode: 'cors',
            credentials: 'include',
            body: JSON.stringify({
                "new": newPassword,
                "previous": currentPassword
            })
        });

        if (res.ok) {
            p_errorMessage = "";
            alert("Password changed successfully. You will be logged out.");
            logoutAndRedirect();
        } else {
            p_errorMessage = "Failed to change password. Please check your current password.";
        } 
    }

    function logoutAndRedirect() {
        user_identifier.set("");
        goto('/');
    }
</script>

<div class="flex flex-col items-center w-full">
    <!-- Profile Settings Header -->
    <div class="text-3xl font-bold mt-5">
      My Profile
    </div>
  
    <!-- Profile Content -->
    <div class="card w-full max-w-4xl bg-base-200 shadow-xl flex-row flex mt-5 items-center justify-between">

        <div class="flex flex-col items-center m-5">
            <img class="upload w-20 h-20 mt-2" alt="Avatar" src={profilePicture} />
            {#if editing}
                <div class="file-input-container mt-4">
                    <input accept="image/png, image/jpeg" bind:files id="avatar" name="avatar" type="file" style="display: none;" />
                    <label for="avatar" class="btn variant-soft">Choose File</label>
                </div>
            {/if}
        </div>

        <div class="form-control m-5">
            <div class="font-bold">Name</div>
            {#if editing}
                <input type="text" bind:value={new_name} class="input input-bordered w-full max-w-xl">
            {:else}
                <div>{name}</div>
            {/if}
            <div class="font-bold">Email</div>
            {#if editing}
                <input type="email" bind:value={new_email} class="input input-bordered w-full max-w-xl">
            {:else}
                <div>{email}</div>
            {/if}
            <div class="font-bold">Biography</div>
            {#if editing}
                <textarea bind:value={new_bio} class="textarea textarea-bordered w-full max-w-xl h-32"></textarea>
            {:else}
                <div>{bio}</div>
            {/if}
            <!-- Show error message -->
            {#if errorMessage}
                <div class="text-red-500">{errorMessage}</div>
            {/if}
        </div>
        
        <!-- Edit Profile Button -->
        <div class="m-5 self-start">
            {#if editing}
                <button on:click={save} class="btn btn-sm btn-success">Save</button>
            {:else}
                <button on:click={edit} class="btn btn-sm btn-neutral">Edit Profile</button>
            {/if}
        </div>
    </div>

    <!-- Change Password Section -->
    <div class="card w-full max-w-4xl bg-base-200 shadow-xl flex-row flex mt-5 items-center justify-between">
        <!-- Information Section -->
        <div class="form-control m-5" style="margin-left: 20px;">
            <div class="font-bold">Current Password</div>
            <input type="password" bind:value={currentPassword} class="input input-bordered password-input">
            <div class="font-bold">New Password</div>
            <input type="password" bind:value={newPassword} class="input input-bordered password-input">
            <div class="font-bold">Confirm New Password</div>
            <input type="password" bind:value={confirmPassword} class="input input-bordered password-input">
            <div class="text-red-500">{p_errorMessage}</div> <!-- Error message -->
        </div>
        <div class="m-5 self-start">
            <div on:click={changePassword} class="btn btn-sm btn-neutral">Change Password</div>
        </div>
    </div>
</div>


<style>
    .password-input {
        width: 500px; 
    }
    .file-input-container {
        display: flex;
        align-items: center;
        gap: 0.5rem; /* 버튼과 사진 사이의 간격 조절 */
    }
</style>
