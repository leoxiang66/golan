<script>
  import { page_idx, isDark } from "../stores";
  import { MessageCircle, GithubIcon, MessageSquare, MessageCircleCode, Info } from "@lucide/svelte";

  // 每次 isDark 变化，都同步写到 localStorage，顺带把 data-theme / class 切过去
  isDark.subscribe((dark) => {
    localStorage.setItem("theme", dark ? "dark" : "light");
    // 如果你用的是 Tailwind CSS 原生 dark 模式：
    if (dark) {
      document.documentElement.classList.add("dark");
    } else {
      document.documentElement.classList.remove("dark");
    }
    // 如果你用的是 DaisyUI 主题切换（data-theme）：
    // document.documentElement.setAttribute('data-theme', dark ? 'synthwave' : 'light');
  });
</script>

<div class="flex">
  <div class="ml-2 flex flex-col items-center h-full w-[70px]">
    <div class="p-6 pt-10">
      <div class="pl-1 w-16">
        <img src="/appicon.png" alt="avatar" class="w-14" />
      </div>
    </div>



    <div class="w-full m-1">
      <!-- svelte-ignore a11y_click_events_have_key_events -->
      <!-- svelte-ignore a11y_no_static_element_interactions -->
      <div
        class="sidebar-item "
        onclick={() => page_idx.set(1)}
      >
        <MessageCircle class={$page_idx == 1? "w-8 h-8 fill-current stroke-[1.5px]" : "w-8 h-8 stroke-[0.5px]"}/>
      </div>
    </div>


    <!-- <div class="w-full m-1">
      <div
        class="sidebar-item "
        onclick={() => page_idx.set(2)}
      >
        <MessageCircleCode class={$page_idx == 2? "w-8 h-8 fill-current stroke-[1.5px]" : "w-8 h-8 stroke-[0.5px]"}/>
      </div>
    </div> -->



    

    <div class="w-full mt-auto">
      <!-- svelte-ignore a11y_click_events_have_key_events -->
      <!-- svelte-ignore a11y_no_static_element_interactions -->
      <div
        class="sidebar-item"
      >
        <Info />
      </div>
    </div>

    <a href="#" class="w-full m-2">
      <button class="sidebar-item">
        <GithubIcon class="w-[20px]" />
      </button>
    </a>


    <!-- <label class="border border-black dark:border-white toggle text-base-content mb-6">
      <input type="checkbox" value="synthwave" bind:checked={$isDark}/>

      <svg
        aria-label="sun"
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 0 24 24"
        ><g
          stroke-linejoin="round"
          stroke-linecap="round"
          stroke-width="2"
          fill="none"
          stroke="currentColor"
          ><circle cx="12" cy="12" r="4"></circle><path d="M12 2v2"></path><path
            d="M12 20v2"
          ></path><path d="m4.93 4.93 1.41 1.41"></path><path
            d="m17.66 17.66 1.41 1.41"
          ></path><path d="M2 12h2"></path><path d="M20 12h2"></path><path
            d="m6.34 17.66-1.41 1.41"
          ></path><path d="m19.07 4.93-1.41 1.41"></path></g
        ></svg
      >

      <svg
        aria-label="moon"
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 0 24 24"
        ><g
          stroke-linejoin="round"
          stroke-linecap="round"
          stroke-width="2"
          fill="none"
          stroke="currentColor"
          ><path d="M12 3a6 6 0 0 0 9 9 9 9 0 1 1-9-9Z"></path></g
        ></svg
      >
    </label> -->
  </div>

  <div class="border-l border-gray-100 h-full"></div>
</div>
