<script>

     import {BASE_API_PATH} from '$lib/env'
     
     export let oldMessages;
     export let messages;
     
     /** @type HTMLDialogElement */
     let searchModal;

     $: chatHistory = oldMessages.concat(messages)
     console.log(chatHistory)
     let searchInput = '';
     let searchResults = [];

     $: {
          searchResults = chatHistory
               .filter(message => 
                    message.contents.toLowerCase().includes(searchInput.toLowerCase())
               )
               .slice(0, 5);
     }


     function escapeRegExp(string) {
        return string.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
     }

    function highlightMatches(text, query) {
        if (!query) {
            return [{ text, match: false }];
        }
        const escapedQuery = escapeRegExp(query);
        const regex = new RegExp(`(${escapedQuery})`, 'gi');
        const parts = text.split(regex);
        return parts.map(part => ({
            text: part,
            match: regex.test(part),
        }));
     }


</script>


<button on:click={() => searchModal.showModal()} class="btn btn-circle">
     <svg xmlns="http://www.w3.org/2000/svg" class="h-10 w-10" viewBox="-17, 45, 300, 150">
          <g fill="#5f65aa" fill-rule="nonzero" stroke="none" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" stroke-miterlimit="10" stroke-dasharray="" stroke-dashoffset="0" font-family="none" font-weight="none" font-size="none" text-anchor="none" style="mix-blend-mode: normal"><g transform="scale(2,2)"><path d="M52.34961,14.40039c-9.725,0 -19.44961,3.69961 -26.84961,11.09961c-14.8,14.8 -14.8,38.89922 0,53.69922c7.4,7.4 17.10039,11.10156 26.90039,11.10156c9.8,0 19.50039,-3.70156 26.90039,-11.10156c14.7,-14.8 14.69844,-38.89922 -0.10156,-53.69922c-7.4,-7.4 -17.12461,-11.09961 -26.84961,-11.09961zM52.30078,20.30078c8.2,0 16.39961,3.09844 22.59961,9.39844c12.5,12.5 12.49961,32.80078 0.09961,45.30078c-12.5,12.5 -32.80078,12.5 -45.30078,0c-12.5,-12.5 -12.5,-32.80078 0,-45.30078c6.2,-6.2 14.40156,-9.39844 22.60156,-9.39844zM52.30078,26.30078c-6.9,0 -13.40078,2.69922 -18.30078,7.69922c-4.7,4.7 -7.29961,10.80039 -7.59961,17.40039c-0.1,1.7 1.20039,2.99961 2.90039,3.09961h0.09961c1.6,0 2.9,-1.30039 3,-2.90039c0.2,-5.1 2.29883,-9.80039 5.79883,-13.40039c3.8,-3.8 8.80156,-5.89844 14.10156,-5.89844c1.7,0 3,-1.3 3,-3c0,-1.7 -1.3,-3 -3,-3zM35,64c-1.65685,0 -3,1.34315 -3,3c0,1.65685 1.34315,3 3,3c1.65685,0 3,-1.34315 3,-3c0,-1.65685 -1.34315,-3 -3,-3zM83.36328,80.5c-0.7625,0 -1.5125,0.30039 -2.0625,0.90039c-1.2,1.2 -1.2,3.09922 0,4.19922l2.5,2.5c-0.6,1.2 -0.90039,2.50039 -0.90039,3.90039c0,2.4 0.89961,4.70039 2.59961,6.40039l12.80078,12.59961c1.8,1.8 4.09844,2.69922 6.39844,2.69922c2.3,0 4.60039,-0.89961 6.40039,-2.59961c3.5,-3.5 3.5,-9.19922 0,-12.69922l-12.79883,-12.80078c-1.7,-1.7 -4.00039,-2.59961 -6.40039,-2.59961c-1.4,0 -2.70039,0.30039 -3.90039,0.90039l-2.5,-2.5c-0.6,-0.6 -1.37422,-0.90039 -2.13672,-0.90039zM91.90039,88.90039c0.8,0 1.59961,0.30039 2.09961,0.90039l12.69922,12.69922c1.2,1.2 1.2,3.09922 0,4.19922c-1.2,1.2 -3.09922,1.2 -4.19922,0l-12.69922,-12.59961c-0.6,-0.6 -0.90039,-1.39922 -0.90039,-2.19922c0,-0.8 0.30039,-1.59961 0.90039,-2.09961c0.6,-0.6 1.29961,-0.90039 2.09961,-0.90039z"></path></g></g>
     </svg>
</button>

<dialog bind:this={searchModal} id="modal search" class="modal">

         <div class="modal-box" style="width: 600px; height: 400px;">
               <div class="p-4 sticky top-0 z-10 bg-base-100"> 
                    <input type="text"
                              placeholder="Search..."
                              class="input input-bordered input-primary w-full"
                              bind:value={searchInput}
                    />
               </div>
          
               <div class="overflow-auto break-words p-2 bg-base-200 rounded-box max-w-lg"> 
                    {#each searchResults as message}
                         <div class="p-2 pb-1 hover:bg-base-300 rounded">
                              {#each highlightMatches(message.contents, searchInput) as part}
                                   <span class={part.match ? 'text-primary' : ''}>
                                        {part.text}
                                   </span>
                              {/each}
                         </div>
                    {/each}
               </div>
         </div>

         <form method="dialog" class="modal-backdrop">
               <button>close</button>
        </form>

 </dialog>