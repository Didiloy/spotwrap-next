import { defineStore } from "pinia";
import { ref } from "vue";
import { EventsOn } from "../../wailsjs/runtime/runtime";

export interface DownloadState {
  downloadMessages: string[];
  isDownloading: boolean;
}

export const useDownloadStore = defineStore("download", () => {
  const downloadMessages = ref<string[]>([]);
  const isDownloading = ref(false);

  function setupEventListener() {
    EventsOn("update_in_download", (message: string) => {
      // Create a new array reference to ensure reactivity
      downloadMessages.value = [...downloadMessages.value, message];

      // Update downloading status
      if (message === "Done") {
        isDownloading.value = false;
      } else if (!isDownloading.value && message.includes("Downloading")) {
        isDownloading.value = true;
      }

      console.log("Download update:", message);
    });
  }

  function clearMessages() {
    downloadMessages.value = [];
    isDownloading.value = false;
  }

  return { downloadMessages, isDownloading, setupEventListener, clearMessages };
});
