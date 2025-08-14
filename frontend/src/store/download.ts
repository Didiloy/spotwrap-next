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
      if (!downloadMessages.value.includes(message)) {
        downloadMessages.value = [...downloadMessages.value, message];
      }

      if (message === "Done") {
        isDownloading.value = false;
      } else if (!isDownloading.value && message.includes("Downloading")) {
        isDownloading.value = true;
      }
    });
  }

  function clearMessages() {
    downloadMessages.value = [];
    isDownloading.value = false;
  }

  return { downloadMessages, isDownloading, setupEventListener, clearMessages };
});
