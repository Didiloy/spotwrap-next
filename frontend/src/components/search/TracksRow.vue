<template>
    <div class="mb-8 w-full">
        <div class="flex items-center justify-between mb-4">
            <h2
                class="relative text-lg font-bold uppercase tracking-widest text-gray-800 dark:text-gray-200 font-montserrat pb-2 after:content-[''] after:absolute after:bottom-0 after:left-0 after:w-12 after:h-[2px] after:bg-[var(--accent-color)]"
            >
                {{ i18n.t("TracksRow.title") }}
            </h2>
        </div>
        <div
            class="bg-white dark:bg-gray-800/50 p-4 rounded-xl shadow-sm w-full mx-auto overflow-hidden"
        >
            <div
                v-for="(track, index) in tracks"
                :key="track.id"
                class="group w-full grid grid-cols-[40px_1fr_auto] items-center p-3 hover:bg-gray-50 dark:hover:bg-gray-700/50 rounded-lg transition-colors duration-200 cursor-pointer"
                @click="handleTrackClick(track.id)"
            >
                <p
                    class="text-gray-500 dark:text-gray-400 text-sm font-medium text-center"
                >
                    {{ index + 1 }}
                </p>

                <!-- Track Info -->
                <div class="flex items-center gap-3 ml-2 overflow-hidden">
                    <img
                        :src="track.album.images[0]?.url"
                        loading="lazy"
                        alt="Track Cover"
                        class="w-10 h-10 rounded-md flex-shrink-0 object-cover shadow-sm"
                    />
                    <div class="min-w-0 overflow-hidden">
                        <p
                            class="font-medium truncate text-gray-900 dark:text-gray-100"
                        >
                            {{ track.name }}
                        </p>
                        <p
                            class="text-sm text-gray-500 dark:text-gray-400 truncate"
                        >
                            {{ formatArtists(track.artists) }}
                        </p>
                    </div>
                </div>

                <!-- Duration & Action -->
                <div class="flex items-center gap-4 ml-auto">
                    <p
                        class="text-sm text-gray-600 dark:text-gray-300 min-w-[50px] text-right"
                    >
                        {{ formatDuration(track.duration_ms) }}
                    </p>
                    <Button
                        @click.stop="handleArrowClick(track.id)"
                        variant="ghost"
                        size="sm"
                        class="opacity-0 group-hover:opacity-100 text-[var(--accent-color)] hover:bg-[var(--accent-color)] hover:text-white transition-opacity duration-200"
                    >
                        <ArrowRight class="w-4 h-4" />
                    </Button>
                </div>
            </div>
        </div>
    </div>
</template>

<script lang="ts" setup>
import { useI18n } from "vue-i18n";
const i18n = useI18n();
import { ArrowRight, Clock } from "lucide-vue-next";
import { Button } from "@/components/ui/button";

defineProps<{
    tracks: Array<{
        id: string;
        name: string;
        duration_ms: number;
        artists: Array<{ name: string }>;
        album: {
            images: Array<{ url: string }>;
        };
    }>;
}>();

const formatDuration = (durationMs: number): string => {
    const minutes = Math.floor(durationMs / 60000);
    const seconds = Math.floor((durationMs % 60000) / 1000);
    return `${minutes}:${seconds.toString().padStart(2, "0")}`;
};

const formatArtists = (artists: Array<{ name: string }>): string => {
    return artists.map((artist) => artist.name).join(", ");
};

const handleArrowClick = (trackId: string) => {
    console.log("Arrow clicked for track:", trackId);
};

const handleTrackClick = (trackId: string) => {
    console.log("Track clicked:", trackId);
};
</script>
