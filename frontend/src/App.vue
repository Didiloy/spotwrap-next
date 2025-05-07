<script lang="ts" setup>
import { onMounted } from "vue";
import AppSidebar from "@/components/sidebar/AppSidebar.vue";
import { SidebarProvider, SidebarTrigger } from "@/components/ui/sidebar";
import Toaster from "@/components/ui/toast/Toaster.vue";
import SpotifyCredentialsModal from "@/components/settings/SpotifyCredentialsModal.vue";
import { useSettingsStore } from "@/store/settings";

const settingsStore = useSettingsStore();

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
    <SidebarProvider class="bg-zinc-100 h-screen flex">
        <AppSidebar class="flex-shrink-0" />
        <main class="relative flex-grow h-full overflow-hidden">
            <div class="absolute top-0 left-0 z-10">
                <SidebarTrigger />
            </div>
            <router-view class="w-full h-full" />
        </main>
    </SidebarProvider>
    
    <SpotifyCredentialsModal v-model:open="settingsStore.showCredentialsModal" />
</template>

<style></style>
