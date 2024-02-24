<svelte:head>
    <title>Motmot</title>
</svelte:head>

<script>
	import Navbar from '$lib/components/Navbar.svelte';
	import '../app.pcss';
</script>

<div class="container mx-auto">
	<Navbar />
	<slot />
</div>
