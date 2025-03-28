<template>
    <div class="mb-6">
        <h2 class="text-lg font-semibold mb-3">
            {{ i18n.t("TracksRow.title") }}
        </h2>
        <div class="bg-white p-4 rounded-lg h-fit overflow-y-auto">
            <div
                v-for="track in tracks"
                :key="track.id"
                class="w-full max-w-full grid grid-cols-[380px_1fr_auto] items-center p-2 border-b cursor-pointer"
            >
                <div class="w-full flex items-center gap-3 mr-4">
                    <img
                        :src="track.album.images[0]?.url"
                        alt="Track Cover"
                        class="w-10 h-10 rounded-md"
                    />
                    <div class="min-w-0 w-full">
                        <p class="font-medium truncate">{{ track.name }}</p>
                        <p
                            class="w-full text-sm text-gray-500 truncate"
                            :title="
                                track.artists
                                    .map((artist: Artist) => artist.name)
                                    .join(', ')
                            "
                        >
                            {{
                                track.artists
                                    .map((artist: Artist) => artist.name)
                                    .join(", ")
                            }}
                        </p>
                    </div>
                </div>
                <div
                    class="text-center flex-grow flex flex-row items-center justify-center"
                >
                    <Clock class="w-5 h-5 inline-block mr-1" />
                    <p class="text-sm text-gray-600">
                        {{ formatDuration(track.duration_ms) }}
                    </p>
                </div>
                <div>
                    <Button
                        @click.stop="handleArrowClick(track.id)"
                        variant="ghost"
                        class="text-[var(--accent-color)] hover:text-white"
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

defineProps<{ tracks: any[] }>();

interface Artist {
    name: string;
}

interface Track {
    id: string;
    name: string;
    duration_ms: number;
    artists: Artist[];
    album: {
        images: { url: string }[];
    };
}
// Function to format duration (milliseconds to mm:ss)
const formatDuration = (durationMs: number): string => {
    const minutes = Math.floor(durationMs / 60000);
    const seconds = Math.floor((durationMs % 60000) / 1000);
    return `${minutes}:${seconds.toString().padStart(2, "0")}`;
};

const handleArrowClick = (trackId: string) => {
    console.log("Arrow clicked for track:", trackId);
};
</script>
