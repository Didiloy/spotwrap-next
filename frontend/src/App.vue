<script lang="ts" setup>
import AppSidebar from "@/components/AppSidebar.vue";
import { SidebarProvider, SidebarTrigger } from "@/components/ui/sidebar";
</script>

<template>
    <SidebarProvider class="bg-zinc-100 h-screen flex">
        <AppSidebar class="flex-shrink-0" />
        <main class="flex-grow h-full overflow-hidden">
            <SidebarTrigger />
            <router-view class="w-full h-full" />
        </main>
    </SidebarProvider>
</template>

<style></style>
