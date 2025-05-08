<!-- App.svelte -->
<script>
  import Home from "./routes/Home.svelte";
  import Image from "./routes/Image.svelte";
  import Container from "./routes/Container.svelte";
  import {page_idx } from "./stores";
</script>

<div class="w-screen h-screen overflow-hidden  ">
  {#if $page_idx === 0}
    <Home />
    {:else if $page_idx === 1}
    <Image />
    {:else if $page_idx === 2}
    <Container />
  {/if}
</div>
