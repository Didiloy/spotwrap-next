<script setup lang="ts">
import { cn } from "@/lib/utils";
import {
    ProgressRoot,
    ProgressIndicator,
    type ProgressRootProps,
} from "reka-ui";
import { type HTMLAttributes } from "vue";

const props = withDefaults(
    defineProps<
        ProgressRootProps & {
            class?: HTMLAttributes["class"];
            infinite?: boolean;
        }
    >(),
    {
        modelValue: 0,
        infinite: false,
    },
);
</script>

<template>
    <ProgressRoot
        :class="
            cn(
                'relative h-2 w-full overflow-hidden rounded-full bg-zinc-900/20 dark:bg-zinc-50/20',
                props.class,
            )
        "
    >
        <ProgressIndicator
            v-if="!props.infinite"
            class="h-full w-full flex-1 bg-zinc-900 transition-all dark:bg-zinc-50"
            :style="`transform: translateX(-${100 - (props.modelValue ?? 0)}%);`"
        />
        <ProgressIndicator
            v-else
            class="h-full w-full bg-zinc-900 animate-infinite-progress"
        />
    </ProgressRoot>
</template>

<style scoped>
@keyframes infinite-progress {
    0% {
        transform: translateX(-100%);
    }
    100% {
        transform: translateX(100%);
    }
}

.animate-infinite-progress {
    animation: infinite-progress 1.5s infinite linear;
}
</style>
