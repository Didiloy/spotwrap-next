<script lang="ts" setup>
import AppSidebar from "@/components/sidebar/AppSidebar.vue";
import { SidebarProvider, SidebarTrigger } from "@/components/ui/sidebar";
import Toaster from "@/components/ui/toast/Toaster.vue";
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
</template>

<style></style>
