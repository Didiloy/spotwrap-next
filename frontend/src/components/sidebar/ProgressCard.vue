<script setup lang="ts">
import { computed } from "vue";
import { Progress } from "@/components/ui/progress";

const props = defineProps({
    message: {
        type: String,
        required: true,
    },
    progress: {
        type: Number,
        default: null, // Null means no progress to show
    },
    showProgress: {
        type: Boolean,
        default: true,
    },
});

const isVisible = computed(() => props.progress !== null && props.showProgress);
</script>

<template>
    <div
        v-if="isVisible"
        class="w-full h-24 bg-zinc-100 rounded-2xl p-4 flex flex-col justify-between border border-zinc-400"
    >
        <span class="text-zinc-900 text-sm text-center">{{ message }}</span>
        <Progress :infinite="true" class="w-full h-2" />
    </div>
</template>

<style scoped></style>
