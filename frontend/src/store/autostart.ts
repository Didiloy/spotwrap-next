import { defineStore } from "pinia";
import {
  Enable,
  Disable,
  IsEnabled,
} from "../../wailsjs/go/autostart/AutoStart";

export const useAutoStartStore = defineStore("autostart", {
  state: () => ({
    isEnabled: false,
  }),
  actions: {
    async checkStatus() {
      this.isEnabled = await IsEnabled();
    },
    async toggleAutoStart() {
      this.isEnabled ? await Disable() : await Enable();
    },
  },
});
