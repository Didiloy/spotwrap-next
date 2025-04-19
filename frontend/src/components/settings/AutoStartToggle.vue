<script setup lang="ts">
import { useAutoStartStore } from "@/store/autostart";
import { Switch } from "@/components/ui/switch";
import { onMounted } from "vue";
import { useI18n } from "vue-i18n";

const i18n = useI18n();
const autoStartStore = useAutoStartStore();

onMounted(async () => {
    await autoStartStore.checkStatus();
});
</script>

<template>
    <div class="flex items-center justify-between p-0 w-full">
        <div>
            <h3 class="font-medium">{{ i18n.t("Settings.auto_start") }}</h3>
            <p class="text-sm text-muted-foreground">
                {{ i18n.t("Settings.auto_start_description") }}
            </p>
        </div>
        <Switch
            v-model="autoStartStore.isEnabled"
            @click="autoStartStore.toggleAutoStart"
        />
    </div>
</template>
