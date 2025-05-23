<script setup lang="ts">
import { computed } from "vue";
import { Progress } from "@/components/ui/progress";
import { useI18n } from "vue-i18n";

const i18n = useI18n();

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
    maxLength: {
        type: Number,
        default: 30, // Default maximum character length
    },
});

const isVisible = computed(() => props.progress !== null && props.showProgress);

const truncatedMessage = computed(() => {
    if (props.message.length <= props.maxLength) {
        return props.message;
    }
    return props.message.substring(0, props.maxLength) + "...";
});
</script>

<template>
    <div
        v-if="isVisible"
        class="w-full h-28 bg-zinc-100 rounded-2xl p-4 flex flex-col justify-between border border-zinc-400"
    >
        <span
            class="text-zinc-900 text-sm text-center truncate"
            :title="message"
        >
            {{ truncatedMessage }}
        </span>
        <div class="flex flex-col gap-1">
            <Progress :infinite="true" class="w-full h-2" />
            <span class="text-xs text-zinc-500 text-center">
                {{ i18n.t("AppSidebar.downloading") }}
            </span>
        </div>
    </div>
</template>
