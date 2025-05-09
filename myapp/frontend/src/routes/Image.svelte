<script>
  import Sidebar from "../components/Sidebar.svelte";
  import { peers } from "../stores";
  import { InviteSocket, Greet, NotifyBackend } from "../../wailsjs/go/app/App";
  import { EventsOn } from "../../wailsjs/runtime/runtime.js";
  import { tick } from "svelte";

  let chatContainer;
  let dialogRef;
  let focusedUser = $state("nil"); //todo: 改成ID
  let chatting = $state(new Map([["self", true]]));

  // 更新对象状态
  // 通过重新赋值来强制触发 Svelte 的更新
  function updateChatting(user, status) {
    chatting.set(user, status);
    chatting = new Map(chatting); // 重新创建 Map 引用
  }

  let inviting = $state(false);
  let receiveInvite = $state("nil");
  let showFullChat = $state(-1);

  async function handleInvite() {
    inviting = true;
    const result = await InviteSocket(focusedUser, 30);

    if (result) {
      console.log("guest accepted");

      // Update chatting state
      updateChatting(focusedUser, true);
      const tmp = focusedUser;
      focusedUser = "nil";
      focusedUser = tmp;

      // Ensure UI updates after state changes
      console.log("chatting after update", chatting);
    } else {
      console.log("guest rejected");
      updateChatting(focusedUser, false);
      focusedUser = "nil"; // Reset focusedUser
      console.log("chatting after reject", chatting);
    }

    // Finish invitation process
    inviting = false;
  }

  let chatHistory = $state([
    [0, "hi, how are you"],
    [1, "I'm fine, thanks! How about you?"],
    [0, "Doing great, just working on that Svelte project. Doing great..."],
    [1, "Nice! Let me know if you need any help."],
    [0, "Sure thing, thanks 😊"],
  ]);

  let newMessage = $state("");

  // Handle 'Enter' key press
  function handleKeydown(event) {
    if (event.key === "Enter") {
      event.preventDefault();
      const text = newMessage.trim();
      if (!text) return;
      chatHistory = [...chatHistory, [0, text]]; // add to chat history
      newMessage = ""; // clear input field
      requestAnimationFrame(scrollToBottom);
    }
  }

  function scrollToBottom() {
    if (chatContainer) {
      chatContainer.scrollTop = chatContainer.scrollHeight;
    }
  }

  $effect(() => {
    if (chatting.has(focusedUser) && chatting.get(focusedUser)) {
      console.log("进入chatting");

      scrollToBottom();
    }

    if (focusedUser != "nil") {
      console.log(focusedUser);
    }
  });

  EventsOn("lan:receive_invite", (...args) => {
    if (args.length === 0) return; // ignore empty calls
    const data = args[0];
    receiveInvite = data;
  });

  EventsOn("lan:conn_closed", (...args) => {
    if (args.length === 0) return; // ignore empty calls
    const id = args[0];

    focusedUser = "nil";
    updateChatting(id, false);
    inviting = false;
    receiveInvite = "nil";
  });
</script>

<div
  class="flex w-full h-full bg-white dark:bg-zinc-800 text-zinc-800 dark:text-white"
>
  {#if receiveInvite !== "nil"}
    <div class="flex flex-col justify-center items-center w-full">
      <p class="text-xl font-mono">
        Received invitation for socket connection from <span
          class="code text-sm text-red-400 bg-slate-100 px-2 py-1 rounded-md"
        >
          {receiveInvite}
        </span>, accept?
      </p>

      <div class="flex">
        <button
          class="btn m-2"
          onclick={() => {
            NotifyBackend(true);
            updateChatting(receiveInvite, true); // Assume `chatting` is a Map
            inviting = false;
            focusedUser = receiveInvite;
            receiveInvite = "nil"; // Clear the invite
          }}
        >
          Yes
        </button>

        <button
          class="btn m-2"
          onclick={() => {
            NotifyBackend(false);
            receiveInvite = "nil";
            inviting = false;
            focusedUser = "nil"; // Reset focusedUser
          }}
        >
          No
        </button>
      </div>
    </div>
  {:else}
    <Sidebar />

    <div class="flex flex-col overflow-y-auto w-[300px]">
      {#each $peers as p, idx}
        <div
          class="flex items-center hover:bg-slate-100"
          class:bg-slate-100={focusedUser === p}
          onclick={() => {
            focusedUser = p;
          }}
        >
          <div class="w-auto py-2 my-4 px-2">
            <div
              class="avatar avatar-placeholder"
              class:avatar-online={chatting.has(p) && chatting.get(p)}
            >
              <div class="bg-neutral text-neutral-content w-12 rounded-full">
                <span class="text-sm">user{idx}</span>
              </div>
            </div>
          </div>
          <p class="text-blue-400 px-2 py-1 rounded-md text-sm mr-2">{p}</p>
        </div>
      {/each}
    </div>

    <div class="border-l border-gray-100 h-full"></div>

    {#if focusedUser !== "nil"}
      {#if chatting.has(focusedUser.toString()) && chatting.get(focusedUser.toString())}
        <div class="flex flex-col justify-start h-full w-full">
          <div
            class="h-[60%] overflow-y-auto overflow-x-hidden scroll-smooth"
            bind:this={chatContainer}
          >
            {#each chatHistory as data, idx}
              <div
                class="chat"
                class:chat-start={data[0] === 1}
                class:chat-end={data[0] === 0}
              >
                <div class="chat-bubble mx-2 my-6 bg-[#EFEFEF] text-black">
                  <div class="dropdown">
                    <div
                      tabindex="0"
                      role="button"
                      class="inline-block max-w-[300px] truncate font-normal font-sans bg-transparent border-0 m-1 hover:bg-transparent hover:border-0 shadow-none text-sm"
                    >
                      {data[1]}
                    </div>
                    <ul
                      tabindex="0"
                      class="dropdown-content menu bg-white rounded-box z-10 w-auto p-2 shadow-sm"
                    >
                      <li class="hover:bg-zinc-200 rounded-md">
                        <button>Copy</button>
                      </li>
                      <li class="hover:bg-zinc-200 rounded-md">
                        <button
                          onclick={() => {
                            showFullChat = idx;
                            dialogRef.showModal();
                          }}
                        >
                          Show Full Text
                        </button>
                      </li>
                    </ul>
                  </div>
                </div>
              </div>
            {/each}
          </div>

          <div class="border-t border-gray-100 w-full"></div>

          <div class="h-[38%]">
            <textarea
              bind:value={newMessage}
              onkeydown={handleKeydown}
              class="w-full h-full textarea focus:outline-none border-none focus:ring-0 bg-transparent"
              placeholder="Type your message here"
            ></textarea>
          </div>
        </div>
      {:else if inviting}
        <div class="block m-auto">
          <span class="loading loading-dots loading-lg"></span>
        </div>
      {:else}
        <div class="block m-auto text-lg">
          <div class="flex flex-col items-center">
            <div class="m-2">
              Invite <span
                class="code text-sm text-red-400 bg-slate-100 px-2 py-1 rounded-md"
                >{focusedUser}</span
              > for socket connection?
            </div>
            <div class="flex">
              <button
                class="btn m-2"
                onclick={async () => {
                  await handleInvite();
                }}>Yes</button
              >
              <button class="btn m-2" onclick={() => (focusedUser = "nil")}
                >Cancel</button
              >
            </div>
          </div>
        </div>
      {/if}
    {/if}
  {/if}
</div>

<dialog bind:this={dialogRef} class="modal">
  <div class="modal-box bg-white text-black">
    {#if showFullChat != -1}
      <p class="py-4">{chatHistory[showFullChat][1]}</p>
    {/if}
    <div class="flex justify-end items-center">
      <button class="btn m-2">Copy</button>
      <form method="dialog">
        <button class="btn m-2">Close</button>
      </form>
    </div>
  </div>
</dialog>
