<script>
  import Sidebar from "../components/Sidebar.svelte";
  import { peers } from "../stores";
  import { InviteSocket, Greet } from "../../wailsjs/go/main/App";
  import { EventsOn } from "../../wailsjs/runtime/runtime";
  let focusedUser = $state(-1);
  let chatting = $state(false);
  let inviting = $state(false);

  EventsOn("lan:socket_accepted", () => {
    inviting = false;
    chatting = true;
  });
</script>

<div
  class="flex w-full h-full bg-white dark:bg-zinc-800 text-zinc-800 dark:text-white"
>
  <Sidebar />

  <div class="flex flex-col overflow-y-auto">
    {#each $peers as p, idx}
      <!-- svelte-ignore a11y_click_events_have_key_events -->
      <!-- svelte-ignore a11y_no_static_element_interactions -->
      <div
        class=" flex items-center hover:bg-slate-100"
        class:bg-slate-100={focusedUser == idx}
        onclick={() => (focusedUser = idx)}
      >
        <div class="w-auto py-2 my-4 px-2">
          <div class="avatar avatar-placeholder">
            <div class="bg-neutral text-neutral-content w-12 rounded-full">
              <span class="text-sm">user{idx}</span>
            </div>
          </div>
        </div>
        <p class="text-zinc-500 text-sm">{p}</p>
      </div>
    {/each}
  </div>

  <div class="border-l border-gray-100 h-full"></div>

  {#if chatting}
    1
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
