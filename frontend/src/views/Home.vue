<template>
    <div class="artist-updates h-full overflow-y-auto p-6">
        <!-- Header -->
        <div class="mb-8">
            <h1 class="text-3xl font-bold text-zinc-900">
                {{ $t("Home.title") }}
            </h1>
            <p class="text-gray-400 mt-2">
                {{ $t("Home.subtitle") }}
            </p>
        </div>

        <!-- Timeline -->
        <div class="relative">
            <!-- Timeline line -->
            <div
                class="absolute left-3 top-0 h-full w-0.5 bg-gradient-to-b from-purple-500/30 via-purple-500/50 to-purple-500/30"
            ></div>

            <!-- Timeline items -->
            <div class="space-y-8 pl-10">
                <div
                    v-for="(artist, index) in timelineItems"
                    :key="index"
                    class="relative group"
                >
                    <!-- Timeline dot -->
                    <div
                        class="absolute -left-10 top-1/2 transform -translate-y-1/2 w-6 h-6 rounded-full bg-gradient-to-br from-purple-600 to-purple-400 border-4 border-gray-900 z-10"
                    ></div>

                    <!-- Content card -->
                    <div
                        class="backdrop-blur-sm rounded-xl p-5 border border-gray-700/50 transition-all duration-300 group-hover:border-purple-500/50"
                        :style="
                            getCardStyle(artist.album.images[0]?.url, index)
                        "
                    >
                        <div class="flex flex-col md:flex-row gap-5">
                            <!-- Album cover -->
                            <div
                                class="relative flex-shrink-0 w-24 h-24 md:w-32 md:h-32"
                            >
                                <img
                                    :src="artist.album.images[0].url"
                                    :alt="artist.album.name"
                                    class="w-full h-full object-cover rounded-lg shadow-lg"
                                />
                                <div
                                    class="absolute inset-0 rounded-lg ring-1 ring-inset ring-white/10"
                                ></div>
                            </div>

                            <!-- Content -->
                            <div class="flex-1">
                                <div
                                    class="flex items-start justify-between gap-4"
                                >
                                    <div>
                                        <h3
                                            class="text-xl font-bold text-white"
                                        >
                                            {{ artist.album.name }}
                                        </h3>
                                        <p class="text-purple-400 mt-1">
                                            {{ artist.artist.name }}
                                        </p>
                                    </div>
                                    <span
                                        class="text-xs px-2 py-1 rounded-full bg-white/10 text-gray-300"
                                    >
                                        {{
                                            formatReleaseDate(
                                                artist.album.release_date,
                                            )
                                        }}
                                    </span>
                                </div>

                                <div class="mt-3 flex flex-wrap gap-2">
                                    <span
                                        class="text-xs px-3 py-1 rounded-full bg-white/10 text-gray-300"
                                    >
                                        {{
                                            artist.album.album_type.toUpperCase()
                                        }}
                                    </span>
                                    <span
                                        v-if="artist.album.total_tracks"
                                        class="text-xs px-3 py-1 rounded-full bg-white/10 text-gray-300"
                                    >
                                        {{ artist.album.total_tracks }} tracks
                                    </span>
                                </div>

                                <div class="mt-4 flex items-center gap-3">
                                    <Button
                                        size="sm"
                                        class="rounded-full bg-purple-600 hover:bg-purple-700 transition-colors"
                                    >
                                        View Album
                                    </Button>
                                    <Button
                                        size="sm"
                                        variant="outline"
                                        class="rounded-full bg-white/10 hover:bg-white/20 transition-colors"
                                    >
                                        Download
                                    </Button>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script lang="ts" setup>
import { ref, onMounted } from "vue";
import { Button } from "@/components/ui/button";
import { GetArtist, GetDominantColor } from "../../wailsjs/go/main/App";
import { GetArtistsFromDB } from "../../wailsjs/go/database/Database";

interface TimelineItem {
    artist: {
        id: string;
        name: string;
        images?: Array<{ url: string }>;
    };
    album: {
        id: string;
        name: string;
        album_type: string;
        release_date: string;
        total_tracks: number;
        images: Array<{ url: string }>;
    };
    type: "album" | "single";
    date: Date;
    dominantColors?: string[];
}

const timelineItems = ref<(TimelineItem & { dominantColors?: string[] })[]>([]);

const formatReleaseDate = (dateString: string) => {
    const options: Intl.DateTimeFormatOptions = {
        year: "numeric",
        month: "short",
        day: "numeric",
    };
    return new Date(dateString).toLocaleDateString(undefined, options);
};

const getCardStyle = (imageUrl: string, index: number) => {
    const colors = timelineItems.value[index]?.dominantColors;
    if (colors && colors.length >= 2) {
        return {
            background: `linear-gradient(135deg, ${colors[0]} 0%, ${colors[1]} 100%)`,
        };
    }
    return {
        background: "linear-gradient(135deg, #4f46e5 0%, #1e40af 100%)",
    };
};

onMounted(async () => {
    const artists = await GetArtistsFromDB();
    const allAlbums: TimelineItem[] = [];

    for (const artist of artists) {
        const artistData = await GetArtist(artist.SpotifyID);

        if (artistData.albums) {
            for (const album of artistData.albums) {
                const dominantColors = album.images?.[0]?.url
                    ? await GetDominantColor(album.images[0].url)
                    : [];

                allAlbums.push({
                    artist: {
                        id: artistData.artist.id,
                        name: artistData.artist.name,
                        images: artistData.artist.images,
                    },
                    album: {
                        id: album.id,
                        name: album.name,
                        album_type: album.album_type,
                        release_date: album.release_date,
                        total_tracks: album.total_tracks,
                        images: album.images,
                    },
                    type: album.album_type === "album" ? "album" : "single",
                    date: new Date(album.release_date),
                    dominantColors,
                });
            }
        }
    }

    timelineItems.value = allAlbums
        .sort((a, b) => b.date.getTime() - a.date.getTime())
        .slice(0, 20);
});
</script>

<style scoped>
.artist-updates {
    scrollbar-width: none; /* Firefox */
}
.artist-updates::-webkit-scrollbar {
    display: none; /* Chrome/Safari */
}

/* Animation for timeline items */
@keyframes fadeIn {
    from {
        opacity: 0;
        transform: translateY(10px);
    }
    to {
        opacity: 1;
        transform: translateY(0);
    }
}

.relative.group {
    animation: fadeIn 0.3s ease-out forwards;
    opacity: 0;
}

/* Apply staggered animations */
.relative.group:nth-child(1) {
    animation-delay: 0.1s;
}
.relative.group:nth-child(2) {
    animation-delay: 0.2s;
}
.relative.group:nth-child(3) {
    animation-delay: 0.3s;
}
.relative.group:nth-child(4) {
    animation-delay: 0.4s;
}
.relative.group:nth-child(5) {
    animation-delay: 0.5s;
}
</style>
