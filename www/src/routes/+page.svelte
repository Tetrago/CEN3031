<script>
	import { BASE_API_PATH } from '$lib/env';
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import logo from '$lib/assets/combinedLogo.jpg';

	/**
	 * @type {any[]}
	 */
	let popularGroups = [];

	// Function to fetch popular groups
	async function fetchPopularGroups() {
			try {
					const response = await fetch ('${Base_API_PATH}/group/popular/5');
					if (response.ok) {
							popularGroups = await response.json();
					} else {
							console.error('Failed to fetch popular groups:', response.statusText);
					}
			} catch (error) {
					console.error('Error fetching popular groups:', error);
			}
	}

	// Fetch popular groups on component mount
	onMount(fetchPopularGroups);
</script>

<div class="container mx-auto">
	<!-- Logo and Name -->
	<div class="flex justify-center items-center mt-8">
		<img src={ logo } alt="Logo" class="h-64 w-64 mr-2" />
	</div>

	<div class="flex justify-center items-center mt-8">
			<h1 class="text-3xl font-semibold">Welcome to MotMot!</h1>
	</div>

	<p class="flex justify-center items-center mt-4 text-lg">What is it, you ask? It's a bird, of course.</p>
	<p class="flex justify-center items-center mt-4 text-lg">It's also our all-encompassing chat app that lets you interact with fellow students taking the same UF class as you.</p>
	<p class="flex justify-center items-center mt-4 text-lg"><i> Ready to get started?</i></p>

	<ul class="mr-1 px-1">
		<li><a href="/courses" class="flex justify-center items-center rounded-full font-bold text-lg"><i>Click to search and join courses...</i></a></li>
	</ul>

	<!-- Popular Groups Section -->
	<h2 class="flex justify-center items-center mt-8 text-2xl font-semibold ml-2">Your most popular groups:</h2>
	<ul class="mt-4">
		{#each popularGroups as group}
				<li>{group.name}</li>
		{/each}
	</ul>
</div>
