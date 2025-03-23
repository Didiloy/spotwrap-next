<script lang="ts" setup>
import AppSidebar from "@/components/AppSidebar.vue";
import { SidebarProvider, SidebarTrigger } from "@/components/ui/sidebar";
import { Button } from "@/components/ui/button";
function handleClick() {
    console.log("Button clicked!");
}
</script>

<template>
    <SidebarProvider>
        <AppSidebar />
        <main class="bg-zinc-100 w-full">
            <SidebarTrigger />
            <!-- <router-view /> -->
            <div
                class="flex flex-col items-center justify-center h-full w-full bg-zinc-100"
            >
                <img
                    alt="Wails logo"
                    src="./assets/images/logo-universal.png"
                    class="w-1/2 h-1/2 object-contain"
                />
                <Button @click="handleClick" class="mt-4">Click me!</Button>
            </div>
        </main>
    </SidebarProvider>
</template>

<style></style>
