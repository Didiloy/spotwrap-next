<script setup lang="ts">
import { useSettingsStore } from "@/store/settings";
import { Switch } from "@/components/ui/switch";
import { useI18n } from "vue-i18n";

const { t } = useI18n();
const settingsStore = useSettingsStore();
</script>

<template>
  <div class="flex items-center justify-between p-0 w-full">
    <div>
      <h3 class="">{{ t("Settings.auto_download") }}</h3>
      <p class="text-sm text-muted-foreground">
        {{ t("Settings.auto_download_description") }}
      </p>
    </div>
    <Switch
      :model-value="settingsStore.autoDownloadNewReleases"
      @update:model-value="settingsStore.toggleAutoDownloadNewReleases"
    />
  </div>
</template> 