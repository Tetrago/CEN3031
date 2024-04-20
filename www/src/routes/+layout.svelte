<svelte:head>
    <title>Motmot</title>
</svelte:head>

<script>
	import Navbar from '$lib/components/Navbar.svelte';
	import { user_identifier } from './stores';
	import '../app.pcss';

	/** @type {import('./$types').LayoutData} */
	export let data;

	$: user_identifier.set(data.ident);
</script>

<div class="container mx-auto">
	<Navbar />
	<slot />
</div>