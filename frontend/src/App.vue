<script lang="ts" setup>
import { onMounted, ref } from "vue";
import AppSidebar from "@/components/sidebar/AppSidebar.vue";
import { SidebarProvider, SidebarTrigger } from "@/components/ui/sidebar";
import Toaster from "@/components/ui/toast/Toaster.vue";
import SpotifyCredentialsModal from "@/components/settings/SpotifyCredentialsModal.vue";
import { useSettingsStore } from "@/store/settings";

const settingsStore = useSettingsStore();
const showMainTrigger = ref(true); // Default to showing the trigger

onMounted(async () => {
  // Check if Spotify credentials are valid
  await settingsStore.loadSpotifyCredentials();
  const hasValidCredentials = await settingsStore.checkCredentialsValidity();
  
  if (!hasValidCredentials) {
    settingsStore.showCredentialsModal = true;
  }
});
</script>

<template>
    <Toaster />
    <SidebarProvider class="bg-zinc-100 h-screen flex" @update:open="showMainTrigger = !$event">
        <AppSidebar class="flex-shrink-0" />
        <main class="relative flex-grow h-full overflow-hidden p-4">
            <div v-if="showMainTrigger" class="fixed top-4 left-4 z-50">
                <SidebarTrigger />
            </div>
            <router-view class="w-full h-full" />
        </main>
    </SidebarProvider>
    
    <SpotifyCredentialsModal v-model:open="settingsStore.showCredentialsModal" />
</template>

<style></style>
