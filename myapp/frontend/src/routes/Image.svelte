<script>
  import Sidebar from "../components/Sidebar.svelte";
  import { peers } from "../stores";
  import { InviteSocket, Greet } from "../../wailsjs/go/main/App";
  import { EventsOn } from "../../wailsjs/runtime/runtime";
  import { onMount } from "svelte";
  // svelte-ignore non_reactive_update
  let chatContainer;
  let dialogRef;
  let focusedUser = $state(-1);
  let chatting = $state(false);
  let inviting = $state(false);
  let showFullChat = $state(-1);

  // 0:self 1:the other
  let chatHistory = $state([
    [0, "hi, how are you"],
    [1, "I'm fine, thanks! How about you?"],
    [
      0,
      "Doing great, just working on that Svelte project. Doing great, just working on that Svelte project. Doing great, just working on that Svelte project.Doing great, just working on that Svelte project. Doing great, just working on that Svelte project. Doing great, just working on that Svelte project.Doing great, just working on that Svelte project. Doing great, just working on that Svelte project. Doing great, just working on that Svelte project.Doing great, just working on that Svelte project. Doing great, just working on that Svelte project. Doing great, just working on that Svelte project.Doing great, just working on that Svelte project. Doing great, just working on that Svelte project. Doing great, just working on that Svelte project.Doing great, just working on that Svelte project. Doing great, just working on that Svelte project. Doing great, just working on that Svelte project.Doing great, just working on that Svelte project. Doing great, just working on that Svelte project. Doing great, just working on that Svelte project.Doing great, just working on that Svelte project. Doing great, just working on that Svelte project. Doing great, just working on that Svelte project.Doing great, just working on that Svelte project. Doing great, just working on that Svelte project. Doing great, just working on that Svelte project.Doing great, just working on that Svelte project. Doing great, just working on that Svelte project. Doing great, just working on that Svelte project.Doing great, just working on that Svelte project. Doing great, just working on that Svelte project. Doing great, just working on that Svelte project.Doing great, just working on that Svelte project. Doing great, just working on that Svelte project. Doing great, just working on that Svelte project.Doing great, just working on that Svelte project. Doing great, just working on that Svelte project. Doing great, just working on that Svelte project.Doing great, just working on that Svelte project. Doing great, just working on that Svelte project. Doing great, just working on that Svelte project.Doing great, just working on that Svelte project. Doing great, just working on that Svelte project. Doing great, just working on that Svelte project. ",
    ],
    [1, "Nice! Let me know if you need any help."],
    [0, "Sure thing, thanks 😊"],
  ]);

  // 绑定给 textarea 的输入内容
  let newMessage = $state("");

  // 处理 Enter 发送
  function handleKeydown(event) {
    if (event.key === "Enter") {
      event.preventDefault();
      const text = newMessage.trim();
      if (!text) return;
      // 添加到历史
      chatHistory = [...chatHistory, [0, text]];
      // 清空输入
      newMessage = "";
      // 滚到底部
      // 下一帧保证 DOM 更新完
      requestAnimationFrame(scrollToBottom);
    }
  }

  // 把滚动条推到底部
  function scrollToBottom() {
    if (chatContainer) {
      chatContainer.scrollTop = chatContainer.scrollHeight;
    }
  }

  $effect(() => {
    if (chatting) {
      scrollToBottom();
    }
  });

  // onMount(() => {
  //   // requestAnimationFrame
  //   scrollToBottom();
  // });

  EventsOn("lan:socket_accepted", () => {
    inviting = false;
    chatting = true;
  });

  EventsOn("lan:conn_closed", () => {
    focusedUser = -1;
    chatting = false;
    inviting = false;
  });
</script>

<div
  class="flex w-full h-full bg-white dark:bg-zinc-800 text-zinc-800 dark:text-white"
>
  <Sidebar />

  <div class="flex flex-col overflow-y-auto w-[300px]">
    {#each $peers as p, idx}
      <!-- svelte-ignore a11y_click_events_have_key_events -->
      <!-- svelte-ignore a11y_no_static_element_interactions -->
      <div
        class=" flex items-center hover:bg-slate-100"
        class:bg-slate-100={focusedUser == idx}
        onclick={() => {
          if (idx != 0) {
            focusedUser = idx;
          }
        }}
      >
        <div class="w-auto py-2 my-4 px-2">
          <div
            class="avatar avatar-placeholder"
            class:avatar-online={idx == 0 || (chatting && focusedUser == idx)}
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

  {#if chatting}
    <div class="flex flex-col justify-start h-full w-full">
      <div
        class="h-[60%] overflow-y-auto overflow-x-hidden scroll-smooth"
        bind:this={chatContainer}
      >
        {#each chatHistory as data, idx}
          <div
            class="chat"
            class:chat-start={data[0] == 1}
            class:chat-end={data[0] == 0}
          >
            <div class="chat-bubble mx-2 my-6 bg-[#EFEFEF] text-black">
              <div class="dropdown">
                <div
                  tabindex="0"
                  role="button"
                  class="
                  inline-block /* 或者 block，都能让宽度限制生效 */
                  max-w-[300px] /* 限制最大宽度，根据 UI 调整 */
                  truncate /* overflow-hidden whitespace-nowrap text-overflow-ellipsis */
                  font-normal font-sans
                  bg-transparent border-0 m-1
                  hover:bg-transparent hover:border-0 shadow-none
                  text-sm
                "
                >
                  {data[1]}
                </div>
                <!-- svelte-ignore a11y_no_noninteractive_tabindex -->
                <ul
                  tabindex="0"
                  class="dropdown-content menu bg-base-100 rounded-box z-10 w-auto p-2 shadow-sm"
                >
                  <li><button>Copy</button></li>
                  <li>
                    <button
                      onclick={() => {
                        showFullChat = idx;
                        dialogRef.showModal();
                      }}>Show Full Text</button
                    >
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
  {:else if focusedUser != -1 && $peers.length > 0}
    <div class="block m-auto text-lg">
      <div class="flex flex-col items-center">
        <div class="m-2">
          Invite <span
            class="code text-sm text-red-400 bg-slate-100 px-2 py-1 rounded-md"
            >{$peers[focusedUser]}</span
          > for socket connection?
        </div>
        <div class="flex">
          <button
            class="btn m-2"
            onclick={async () => {
              inviting = true;
              const result = await InviteSocket($peers[focusedUser], 30);
              if (result) {
                chatting = true;
              } else {
                chatting = false;
                inviting = false;
                focusedUser = -1;
              }
            }}>Yes</button
          >
          <button class="btn m-2" onclick={() => (focusedUser = -1)}
            >Cancel</button
          >
        </div>
      </div>
    </div>
  {/if}
</div>

<dialog bind:this={dialogRef} class="modal">
  <div class="modal-box">
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
