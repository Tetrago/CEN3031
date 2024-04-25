<script>
     import { onMount, onDestroy } from 'svelte';
     import { page } from '$app/stores'; 
	import { BASE_API_PATH } from '$lib/env';
     import { BASE_WS_PATH } from '$lib/env';
     import { user_identifier } from '../../stores'
     import { restrictedWordsRegex } from './filter';
     import SearchButton from '$lib/components/SearchButton.svelte';

    /** @type {import('./$types').PageData} */
     export let data;

     let chatHistory = data.post.chatHistory;
     let groupId = data.post.id;

     let message='';
     $: messages = []
     let oldMessages = chatHistory.map(item => item).reverse();
     let socket; 



     onMount(() => {
          fetchDisplayNamesForOldMessages();
          socket = new WebSocket(`${BASE_WS_PATH}/ws/${groupId}`);

          socket.onmessage = async (event) => {
               const messageData = JSON.parse(event.data);
               if(data.blocked.includes(messageData.user_ident)) return;
               const display_name = await fetchUser(messageData.user_ident);
               const newMessage = { ...messageData, display_name: display_name };
               messages = [...messages, newMessage];
          };

          return () => {
               socket.close();
          };
     });

     async function fetchDisplayNamesForOldMessages() {
          const promises = oldMessages.map(async (msg) => {
               const display_name = await fetchUser(msg.user_ident);
               return { ...msg, display_name }; 
          });
          oldMessages = await Promise.all(promises);
     }


     async function sendMessage() {
          if (message.trim() !== '') {
               message = message.replace(restrictedWordsRegex, match => '*' .repeat(match.length));
               socket.send(message);
               const display_name = await fetchUser($user_identifier);
               messages = [...messages, {user_ident: $user_identifier, contents: message, display_name: display_name}]
               console.log(messages)
               message = ''; // Clear the message input after sending
          }
     }

     async function handleKeyup(event) {
          if (event.key === 'Enter' && !event.shiftKey && !event.ctrlKey) {
               event.preventDefault(); 
               sendMessage();
          }
     }

     async function fetchUser(ident){
          const res = await fetch(`${BASE_API_PATH}/user/get/${ident}`, {
          method: 'get',
          mode: 'cors',
          credentials: 'include'
          });

          const data = await res.json();
          return data.display_name;
     }

     async function blockUser(ident) {
          await fetch(`${BASE_API_PATH}/user/block`, {
               method: 'post',
               mode: 'cors',
               credentials: 'include',
               body: JSON.stringify({ ident })
          });

          oldMessages = oldMessages.filter(m => m.user_ident != ident);
          messages = messages.filter(m => m.user_ident != ident);
     }


</script>
   
<style>

   
     .messages-container {
       height: calc(75vh - 10px); 
       overflow-y: auto; 
       margin-bottom: -2px; 
     }
   
     .text-entry-box {
       position: fixed;
       bottom: 3%;
       width: 65%; 
       left: 50%;
       transform: translateX(-50%); 
       z-index: 1000; 
     }

     .search-button {
       position: fixed;
       bottom: 2.75%;
       left:10%;
     }

     .chat-bub{
       background-color: rgb(50, 50, 50);
     }

     .username{
       position: sticky;
     }

     .prof-pic{
          padding-left: 5px;
     }

</style>


<div class="card bg-base-200">
     <div class="messages-container">
          {#each oldMessages as message} 
          {#if $user_identifier === message.user_ident}
          <div class="flex flex-col h-100 mb-5">
               <div class="flex-grow overflow-auto">
                    <div class="flex items-center ">
                         <div class="avatar prof-pic">
                              <div class="w-10 rounded-full ">
                                   <img alt="User Profile" src={`${BASE_API_PATH}/user/profile_picture/${message.user_ident}`} />
                              </div>
                         </div>

                         <div class="flex-col ml-3">
                              <div class="font-bold chat-header username">{message.display_name}</div>
                              <div class="chat ">
                                   <div class="chat-bubble chat-bub max-w-xl break-words">{message.contents}</div>
                              </div>
                         </div>

                    </div>
               </div>
          </div>
          {:else}
          <div class="flex flex-col h-100 mb-5 dropdown dropdown-bottom dropdown-begin">
               <div tabindex="0" class="flex-grow overflow-auto">
                    <div class="flex items-center ">
                         <div tabindex="0"class="avatar prof-pic">
                              <div class="w-10 rounded-full ">
                                   <img alt="User Profile" src={`${BASE_API_PATH}/user/profile_picture/${message.user_ident}`} />
                              </div>
                         </div>

                         <div class="flex-col ml-3">
                              <div class="font-bold chat-header username">{message.display_name}</div>
                              <div class="chat ">
                                   <div class="chat-bubble chat-bub max-w-xl break-words">{message.contents}</div>
                              </div>
                         </div>

                    </div>
               </div>

               <ul tabindex="0" class="dropdown-content z-[1] menu p-2 shadow bg-base-300 rounded-box w-52">
                   <li><button on:click={() => blockUser(message.user_ident)}>Block</button></li>
               </ul>
          </div>
          {/if}
          {/each}

          {#each messages as mess} 
          <div class="flex flex-col h-100 mb-5 dropdown dropdown-bottom dropdown-begin">
               <div tabindex="0" class="flex-grow overflow-auto">
                    <div class="flex items-center">
                         <div class="dropdown dropdown-bottom dropdown-end">
                              <div tabindex="0" class="btn btn-ghost btn-circle avatar" role="button">
                                   <div class="w-10 rounded-full">
                                        <img alt="User Profile" src={`${BASE_API_PATH}/user/profile_picture/${mess.user_ident}`} />
                                   </div>
                              </div>
                              <ul tabindex="0" class="dropdown-content z-[1] menu p-2 shadow bg-base-300 rounded-box w-52">
                                  <li><button on:click={() => blockUser(message.user_ident)}>Block</button></li>
                              </ul>
                         </div>
                         <div class="flex flex-col ml-2">
                              <div class="font-bold chat-header username">{mess.display_name}</div>
                              <div class="chat">
                                   <div class="chat-bubble chat-bub max-w-xl break-words">{mess.contents}</div>
                              </div>
                         </div>
                    </div>
               </div>

               {#if $user_identifier !== mess.user_ident}
               <ul tabindex="0" class="dropdown-content z-[1] menu p-2 shadow bg-base-300 rounded-box w-52">
                   <li><button on:click={() => blockUser(mess.user_ident)}>Block</button></li>
               </ul>
               {/if}
          </div>
          {/each}
     </div>

</div>

<div class="flex flex-col sm:flex-row items-center space-y-2 sm:space-y-0 sm:space-x-2">
     <div class="text-entry-box">
          <input type="text"
               placeholder="Type here"
               class="input input-bordered input-primary w-full max-w-m"
               bind:value={message}
               on:keyup={handleKeyup} 
          />
     </div>
     
     <div class="search-button">
          <SearchButton {oldMessages} {messages}/>
     </div>

</div>
