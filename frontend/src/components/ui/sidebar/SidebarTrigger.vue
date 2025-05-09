<script setup lang="ts">
import type { HTMLAttributes } from "vue";
import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import { PanelLeft, PanelRightClose } from "lucide-vue-next";
import { useSidebar } from "./utils";

const props = defineProps<{
    class?: HTMLAttributes["class"];
    variant?: "default" | "sidebar";
}>();

const { toggleSidebar, state } = useSidebar();
</script>

<template>
    <Button
        data-sidebar="trigger"
        :variant="variant === 'sidebar' ? 'ghost' : 'outline'"
        size="icon"
        :class="cn('h-8 w-8 rounded-md transition-all', 
            state === 'expanded' ? 'data-[state=expanded]' : 'data-[state=collapsed]',
            props.class)"
        @click="toggleSidebar"
    >
        <PanelLeft v-if="state === 'collapsed'" class="h-4 w-4" />
        <PanelRightClose v-else class="h-4 w-4" />
        <span class="sr-only">Toggle Sidebar</span>
    </Button>
</template>
