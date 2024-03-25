<script>
     import { onMount, onDestroy } from 'svelte';
     import { page } from '$app/stores'; 
	import { BASE_API_PATH } from '$lib/env';
     import { BASE_WS_PATH } from '$lib/env';

    /** @type {import('./$types').PageData} */
     export let data;

     let chatHistory = data.post.chatHistory;
     let groupId = data.post.id;

     let message='';
     $: messages = []
     let oldMessages = chatHistory.map(item => item).reverse();
     let socket; 


     onMount(() => {
        socket = new WebSocket(`${BASE_WS_PATH}/ws/${groupId}`);

        socket.onmessage = (event) => {
            const messageData = JSON.parse(event.data);
            const newMessage = messageData;
            messages = [...messages, newMessage]
          };

        return () => {
            socket.close(); // 
        };
    });

     async function sendMessage() {
          if (message.trim() !== '') {
               console.log("Sending message:", message);
               socket.send(message);
               messages = [...messages, message]
               message = ''; // Clear the message input after sending
          }
     }

     async function handleKeyup(event) {
          if (event.key === 'Enter' && !event.shiftKey && !event.ctrlKey) {
               event.preventDefault(); // Prevent the default action to stop from inserting a newline (if applicable)
               sendMessage();
          }
     }

     async function fetchUser(ident){
               const res = await fetch(`${BASE_API_PATH}/user/get/${ident}`, {
               method: 'get',
               mode: 'cors',
               credentials: 'include'
          });

          return await res.json()
     }

     async function fetch_profile_pic(user_ident){
          const res = await fetch(`${BASE_API_PATH}/user/profile_picture/${user_ident}`, {
               method: 'get', 
               mode: 'cors',
               credentials: 'include'
          })
          return res
     }
     


</script>
   



<div class="card">
     {#each oldMessages as message} 
     <div class="flex flex-col h-100">
          <div class="flex-grow overflow-auto">
               <div class="chat chat-start">
                    <div class="chat-image avatar">
                         <div class="w-10 rounded-full">
                              <img alt="Tailwind CSS chat bubble component" src="https://daisyui.com/images/stock/photo-1534528741775-53994a69daeb.jpg" />
                         </div>
                    </div>
                    <div class="chat-bubble">{message.contents}</div>
               </div>
               
          </div>
     </div>
     {/each}

     {#each messages as mess} 
     <div class="flex flex-col h-100">
          <div class="flex-grow overflow-auto">
               <div class="chat chat-start">
                    <div class="chat-image avatar">
                         <div class="w-10 rounded-full">
                              <img alt="Tailwind CSS chat bubble component" src="https://daisyui.com/images/stock/photo-1534528741775-53994a69daeb.jpg" />
                         </div>
                    </div>
                    <div class="chat-bubble">{mess}</div>
               </div>
               
          </div>
     </div>
     {/each}
</div>

<div class="p-4">
     <input       type="text"
     placeholder="Type here"
     class="input input-bordered input-primary w-full max-w-m"
     bind:value={message}
     on:keyup={handleKeyup} />
</div>