<script>
  import { page_idx,isDark } from "../stores";
  import { Layers, Cuboid, ChartLine, GithubIcon } from "@lucide/svelte";

  // 每次 isDark 变化，都同步写到 localStorage，顺带把 data-theme / class 切过去
  isDark.subscribe((dark) => {
    localStorage.setItem('theme', dark ? 'dark' : 'light');
    // 如果你用的是 Tailwind CSS 原生 dark 模式：
    if (dark) {
      document.documentElement.classList.add('dark');
    } else {
      document.documentElement.classList.remove('dark');
    }
    // 如果你用的是 DaisyUI 主题切换（data-theme）：
    // document.documentElement.setAttribute('data-theme', dark ? 'synthwave' : 'light');
  });

</script>

<div class="flex">
  <div class="ml-2 flex flex-col items-center h-full w-[70px]">
    <div class="p-6 pt-10">
      <div class="pl-1 w-16">
        <img src="/box.png" alt="avatar" class="w-14" />
      </div>
    </div>

    <!-- <div class="w-full m-1 flex justify-center tooltip" data-tip="Status">
      <button
        class="sidebar-item hover:bg-green-500 {$page_idx === 0
          ? 'bg-green-600 hover:bg-green-600'
          : ''}"
        on:click={() => page_idx.set(0)}
      >
        <ChartLine />
      </button>
    </div> -->

    <div class="w-full m-1 flex justify-center tooltip" data-tip="格式转换">
      <button
        class="sidebar-item hover:bg-green-500 {$page_idx === 1
          ? 'bg-green-600 hover:bg-green-600'
          : ''}"
        on:click={() => page_idx.set(1)}
      >
        <Layers />
      </button>
    </div>

    <!-- <div class="w-full m-1 flex justify-center tooltip" data-tip="Containers">
      <button
        class="sidebar-item hover:bg-green-500 {$page_idx === 2
          ? 'bg-green-600 hover:bg-green-600'
          : ''}"
        on:click={() => page_idx.set(2)}
      >
        <Cuboid />
      </button>
    </div> -->

    <a href="#" class="w-full mb-6 mt-auto flex justify-center">
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
