<template>
    <div class="mb-8 w-full">
        <div class="flex items-center justify-between mb-4">
            <h2
                class="relative text-lg font-bold uppercase tracking-widest text-gray-800 dark:text-gray-200 font-montserrat pb-2 after:content-[''] after:absolute after:bottom-0 after:left-0 after:w-12 after:h-[2px] after:bg-[var(--accent-color)]"
            >
                {{ i18n.t("TracksRow.title") }}
            </h2>
        </div>
        <div class="p-2">
            <div class="space-y-2">
                <div
                    v-for="(track, index) in tracks"
                    :key="track.id"
                    class="flex items-center p-4 hover:bg-zinc-400/50 rounded-lg transition-colors"
                    @click="handleTrackClick(track.id)"
                >
                    <img
                        v-if="track.album?.images[0]?.url"
                        :src="track.album?.images[0]?.url"
                        loading="lazy"
                        alt="Track Cover"
                        class="w-10 h-10 rounded-md flex-shrink-0 object-cover"
                    />
                    <div class="w-8 text-gray-400 text-center mr-4">
                        {{ index + 1 }}
                    </div>
                    <div class="flex-grow">
                        <div class="font-medium">{{ track.name }}</div>
                        <div class="text-sm text-gray-400">
                            {{
                                track.artists
                                    ?.map((a: any) => a.name)
                                    .join(", ")
                            }}
                        </div>
                    </div>
                    <div class="text-gray-400 text-sm">
                        {{ formatDuration(track.duration_ms) }}
                    </div>
                    <Button
                        @click="handleTrackClick(track.id)"
                        variant="ghost"
                        size="sm"
                        class="ml-4"
                    >
                        <ArrowRight class="h-4 w-4" />
                    </Button>
                    <!-- <Button
                        @click.stop="downloadTrack(track)"
                        variant="ghost"
                        size="sm"
                        class="ml-2"
                    >
                        <DownloadIcon class="h-4 w-4" />
                    </Button> -->
                </div>
            </div>
        </div>
    </div>
</template>

<script lang="ts" setup>
import { useI18n } from "vue-i18n";
const i18n = useI18n();
import { ArrowRight } from "lucide-vue-next";
import { Button } from "@/components/ui/button";
import { onMounted } from "vue";

onMounted(() => {
    console.log("TracksRow mounted");
});

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
