<script>
	import { goto } from '$app/navigation';
    import logo from '$lib/assets/logo.png';
	import { user_identifier } from '../../routes/stores';
	import UserWidget from './UserWidget.svelte';
</script>

<div class="navbar bg-base-200 rounded-full my-5">
    <div class="navbar-start">
        <div class="btn btn-ghost btn-circle avatar" on:click={() => goto("/")}>
            <div class="w-10 rounded-full">
                <img alt="Motmot" src={logo}/>
            </div>
        </div>

        <p class="ml-1 cursor-default text-2xl font-bold select-none">Motmot</p>
    </div>

    <div class="navbar-middle lg:flex">
        <ul class="menu menu-horizontal px-1 mr-1">
            <li><a href="/courses" class="rounded-full font-bold">Courses</a></li>
        </ul>
        {#if $user_identifier !== ""}
            <ul class="menu menu-horizontal px-1 mr-1">
                <li><a href="/mychats" class="rounded-full font-bold">My Chats</a></li>
            </ul>
        {/if}
    </div>

    <div class="navbar-end">
        <UserWidget />
    </div>
</div>