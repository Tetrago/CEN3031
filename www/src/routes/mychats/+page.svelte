<script>
     import { BASE_API_PATH } from "$lib/env";
     import { onMount } from 'svelte';
     import { goto } from '$app/navigation';
   
     /** @type {import('./$types').PageData}*/
     export let data;
   
     let groups = data.groups;
     let groupColors = [];
   
      // Predefined dark color palette
     const darkColors = ['#270a02', '#1c2a36', '#161827', '#20142a', '#171123', '#0a2404', '#2f260e', '#200918'];

     onMount(() => {
          groupColors = groups.map((_, i) => darkColors[i % darkColors.length]);
     });

</script>



<style>
     .grid-container {
       display: grid;
       grid-template-columns: repeat(auto-fit, minmax(300px, 1fr)); 
       gap: 30px; 
       justify-content: center; 
       padding: 30px; 
     }
   
     .card {
       width: 100%; 
       margin:  auto; 
     }

     .dark-btn {
          background-color: #03070a73; 
          color: #e6e6e6; 
          border: none; 
     }

</style>
   
<div class="grid-container">
     {#each groups as {name}, i}
          <div class="card bg-base-200 shadow-xl" style="background-color: {groupColors[i]};">
               <div class="card-body">
                    <h1 class="card-title">{name}</h1>
                    <div class="card-actions justify-end">
                         <button class="btn btn-primary dark-btn" on:click={() => goto(`/mychats/${name}`)}>Enter chat</button>                   
                     </div>
               </div>
          </div>
     {/each}
</div>
   

