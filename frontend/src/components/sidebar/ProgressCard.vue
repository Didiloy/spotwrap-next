<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from "vue";
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
});

const isVisible = computed(() => props.progress !== null && props.showProgress);

const containerRef = ref<HTMLElement | null>(null);
const textRef = ref<HTMLElement | null>(null);
const shouldScroll = ref(false);
const durationSec = ref(0);
const translateTo = ref("0px");

const SPEED_PX_PER_SEC = 50; // Lower is slower

function updateScrollMetrics() {
    const container = containerRef.value;
    const text = textRef.value;
    if (!container || !text) return;
    const containerWidth = container.clientWidth;
    const textWidth = text.scrollWidth;
    if (textWidth > containerWidth) {
        shouldScroll.value = true;
        const distance = textWidth - containerWidth;
        durationSec.value = Math.max(4, distance / SPEED_PX_PER_SEC);
        translateTo.value = `-${distance}px`;
    } else {
        shouldScroll.value = false;
        durationSec.value = 0;
        translateTo.value = "0px";
    }
}

onMounted(() => {
    nextTick(updateScrollMetrics);
    window.addEventListener("resize", updateScrollMetrics);
});

onBeforeUnmount(() => {
    window.removeEventListener("resize", updateScrollMetrics);
});

watch(
    () => props.message,
    async () => {
        await nextTick();
        updateScrollMetrics();
    }
);
</script>

<template>
    <div
        v-if="isVisible"
        class="w-full h-28 bg-zinc-100 rounded-2xl p-4 flex flex-col justify-between border border-zinc-400"
    >
        <div class="text-zinc-900 text-sm text-center" :title="message">
            <div ref="containerRef" class="marquee-container">
                <div
                    ref="textRef"
                    class="marquee-content"
                    :class="{ scrolling: shouldScroll }"
                    :style="shouldScroll ? { '--toX': translateTo, '--duration': durationSec + 's' } as any : undefined"
                >
                    {{ message }}
                </div>
            </div>
        </div>
        <div class="flex flex-col gap-1">
            <Progress :infinite="true" class="w-full h-2" />
            <span class="text-xs text-zinc-500 text-center">
                {{ i18n.t("AppSidebar.downloading") }}
            </span>
        </div>
    </div>
</template>

<style scoped>
.marquee-container {
    overflow: hidden;
    white-space: nowrap;
}

.marquee-content {
    display: inline-block;
}

.scrolling {
    animation-name: marqueeAlt;
    animation-duration: var(--duration);
    animation-timing-function: linear;
    animation-iteration-count: infinite;
    animation-direction: alternate;
}

@keyframes marqueeAlt {
    from {
        transform: translateX(0);
    }
    to {
        transform: translateX(var(--toX));
    }
}
</style>
