import { writable } from "svelte/store";
import { EventsOn } from "../wailsjs/runtime/runtime";

export const page_idx = writable(1);
// 用一个 store 来保存当前主题状态（true = 深色 / synthwave，false = 浅色 / 默认 light）
export const isDark = writable(false);

export const peers = writable([]);

EventsOn("lan:peers", (...args) => {
  if (args.length === 0) return; // ignore empty calls
  const data = args[0];
  if (!Array.isArray(data)) return; // make sure it's really an array

  //   console.log("peers array:", data);
  //   console.log("number of peers:", data.length);
  data.sort();
  peers.set(data);
});
