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

        <!-- Loading state -->
        <div
            v-if="loading"
            class="flex flex-col items-center justify-center py-12"
        >
            <div
                class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-purple-500 mb-4"
            ></div>
            <p class="text-gray-400">{{ $t("Home.loading") }}</p>
            <p class="text-gray-400 text-sm mt-2 max-w-md text-center">{{ $t("Home.loading_rate_limit") }}</p>
        </div>

        <!-- Empty state -->
        <div
            v-else-if="!loading && timelineItems.length === 0"
            class="flex flex-col items-center justify-center py-12"
        >
            <svg
                xmlns="http://www.w3.org/2000/svg"
                class="h-16 w-16 text-gray-400 mb-4"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
            >
                <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="1.5"
                    d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                />
            </svg>
            <h3 class="text-xl font-medium text-gray-300 mb-2">
                {{ $t("Home.emptyTitle") }}
            </h3>
            <p class="text-gray-400 text-center max-w-md">
                {{ $t("Home.emptySubtitle") }}
            </p>
            <Button
                class="mt-4 rounded-full bg-purple-600 hover:bg-purple-700 transition-colors"
                @click="goToSearch"
            >
                {{ $t("Home.discoverArtists") }}
            </Button>
        </div>

        <!-- Timeline -->
        <div v-else class="relative">
            <!-- Timeline line -->
            <div
                class="absolute left-3 top-0 h-full w-0.5 bg-gradient-to-b from-purple-500/30 via-purple-500/50 to-purple-500/30"
            ></div>

            <!-- Timeline items -->
            <div class="space-y-8 pl-10">
                <div
                    v-if="timelineItems.some((item) => item.isNewRelease)"
                    class="relative"
                >
                    <div
                        class="absolute -left-10 top-1/2 transform -translate-y-1/2 w-6 h-6 rounded-full bg-gradient-to-br from-green-500 to-green-300 border-4 border-gray-900 z-10"
                    ></div>
                    <div
                        class="pl-6 text-green-400 font-bold uppercase text-sm tracking-wider"
                    >
                        {{ $t("Home.newReleases") }}
                    </div>
                    <div
                        class="h-px bg-gradient-to-r from-green-500/30 to-transparent mt-2 mb-4"
                    ></div>
                </div>
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
                        <div
                            v-if="artist.isNewRelease"
                            class="absolute -top-2 -right-2 bg-green-500 text-white text-xs font-bold px-2 py-1 rounded-full z-10 shadow-md"
                        >
                            {{ $t("Home.newRelease") }}
                        </div>
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
                                        {{
                                            artist.album.total_tracks +
                                            " " +
                                            $t("Home.tracks")
                                        }}
                                    </span>
                                </div>

                                <div class="mt-4 flex items-center gap-3">
                                    <Button
                                        size="sm"
                                        class="rounded-full bg-purple-600 hover:bg-purple-700 transition-colors"
                                        @click="goToAlbum(artist.album.id)"
                                    >
                                        {{ $t("Home.viewAlbum") }}
                                    </Button>
                                    <Button
                                        v-if="artist.isNewRelease"
                                        size="sm"
                                        variant="outline"
                                        class="rounded-full text-white border-white/20 bg-white/10 transition-colors"
                                        @click="
                                            markAsSeen(artist.artist.id, index)
                                        "
                                        :loading="markingAsSeenIndex === index"
                                    >
                                        {{
                                            markingAsSeenIndex === index
                                                ? $t("Home.processing")
                                                : $t("Home.markAsSeen")
                                        }}
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
import {
    GetArtist,
    GetArtistsFromDB,
    IsANewRelease,
    AddArtist,
} from "../../wailsjs/go/main/App";
import { GetDominantColor } from "../../wailsjs/go/utils/Utils";
import { useRouter } from "vue-router";

const router = useRouter();

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
    isNewRelease?: boolean;
}

const loading = ref(true);
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
            backdropFilter: 'blur(8px)',
        };
    }
    // Fallback to default gradient if colors not loaded yet
    return {
        background: 'linear-gradient(135deg, rgba(79, 70, 229, 0.2) 0%, rgba(30, 64, 175, 0.2) 100%)',
        backdropFilter: 'blur(8px)',
    };
};

async function loadDominantColorsForItem(item: TimelineItem, index: number) {
    try {
        if (item.album.images?.[0]?.url) {
            const colors = await GetDominantColor(item.album.images[0].url);
            timelineItems.value = timelineItems.value.map((el, i) =>
                i === index ? { ...el, dominantColors: colors } : el,
            );
        }
    } catch (error) {
        console.error("Error loading dominant colors:", error);
    }
}

onMounted(async () => {
    try {
        const artists = await GetArtistsFromDB();

        if (artists.length === 0) {
            loading.value = false;
            return;
        }

        const allAlbums: TimelineItem[] = [];

        for (const artist of artists) {
            const artistData = await GetArtist(artist.SpotifyID);

            if (artistData.albums) {
                for (const album of artistData.albums) {
                    // Check if this is a new release
                    const isNewRelease = await IsANewRelease(
                        artist.SpotifyID,
                        album,
                    );

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
                        isNewRelease: isNewRelease,
                    });
                }
            }
        }

        // Sort with new releases first
        timelineItems.value = allAlbums
            .sort((a, b) => {
                // New releases come first
                if (a.isNewRelease && !b.isNewRelease) return -1;
                if (!a.isNewRelease && b.isNewRelease) return 1;
                // Then sort by date
                return b.date.getTime() - a.date.getTime();
            })
            .slice(0, 20);

        // Set loading to false immediately after setting timelineItems
        loading.value = false;

        // Load dominant colors asynchronously
        timelineItems.value.forEach((item, index) => {
            loadDominantColorsForItem(item, index);
        });
    } catch (error) {
        console.error("Error fetching artist data:", error);
        loading.value = false;
    }
});

function goToSearch() {
    router.push("/search");
}

function goToAlbum(id: string) {
    router.push(`/album/${id}`);
}

const markingAsSeenIndex = ref<number | null>(null);

async function markAsSeen(artistId: string, index: number) {
    try {
        await AddArtist(artistId);
        timelineItems.value = timelineItems.value.map((item, i) =>
            i === index ? { ...item, isNewRelease: false } : item,
        );
    } catch (error) {
        console.error("Error marking release as seen:", error);
    }
}
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
