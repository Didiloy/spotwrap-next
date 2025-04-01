<template>
    <div class="artist-detail h-full overflow-y-auto">
        <!-- Hero Section -->
        <div class="relative w-full" :style="heroSectionStyle">
            <div
                class="container mx-auto px-6 py-12 flex flex-col md:flex-row items-start gap-8"
            >
                <div
                    class="relative w-48 h-48 md:w-64 md:h-64 lg:w-80 lg:h-80 flex-shrink-0"
                >
                    <img
                        :src="artist?.images[0]?.url"
                        alt="Artist image"
                        class="w-full h-full object-cover rounded-lg shadow-2xl"
                    />
                    <div
                        class="absolute inset-0 rounded-lg ring-1 ring-inset ring-white/10"
                    ></div>
                </div>

                <!-- Artist Info -->
                <div class="flex-1 space-y-4 text-white pt-4">
                    <div>
                        <h1 class="text-4xl font-bold mt-1 mb-2">
                            {{ artist?.name }}
                        </h1>
                        <div
                            class="flex flex-wrap items-center gap-x-4 gap-y-2 text-gray-300"
                        >
                            <span class="text-sm">
                                {{ formatNumber(artist?.followers?.total) }}
                                {{ i18n.t("ArtistDetails.followers") }}
                            </span>
                            <span class="text-sm">
                                {{ artist?.genres?.slice(0, 3).join(", ") }}
                            </span>
                        </div>
                    </div>
                    <!-- Action Button -->
                    <div class="flex flex-col items-start gap-4 pt-4">
                        <Button
                            @click="toggleFollow"
                            :variant="isFollowing ? 'default' : 'outline'"
                            class="px-8 py-3 text-lg font-bold rounded-full transition-colors"
                            :class="{
                                'bg-purple-600 hover:bg-purple-700':
                                    isFollowing,
                                'bg-white/10 hover:bg-white/20': !isFollowing,
                            }"
                        >
                            {{
                                isFollowing
                                    ? i18n.t("ArtistDetails.following")
                                    : i18n.t("ArtistDetails.follow")
                            }}
                        </Button>
                    </div>
                </div>
            </div>
        </div>

        <!-- Content Section -->
        <div class="p-6 space-y-8">
            <!-- Albums Section -->
            <div v-if="albums.length">
                <AlbumsRow :albums="albums" />
            </div>

            <!-- Popular Tracks -->
            <TracksRow
                v-if="topTracks.length"
                :tracks="topTracks"
                :title="i18n.t('ArtistDetails.popular_tracks')"
            />

            <!-- Singles Section -->
            <div v-if="singles.length">
                <AlbumsRow
                    :albums="singles"
                    :title="i18n.t('ArtistDetails.singles')"
                />
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import { useRoute } from "vue-router";
import { Button } from "@/components/ui/button";
import { useToast } from "@/components/ui/toast/use-toast";
import { GetArtist, GetDominantColor } from "../../../wailsjs/go/main/App";
import {
    AddArtist,
    RemoveArtist,
    GetArtistsFromDB,
} from "../../../wailsjs/go/database/Database";
import AlbumsRow from "@/components/search/AlbumsRow.vue";
import TracksRow from "@/components/search/TracksRow.vue";
import { useI18n } from "vue-i18n";

const route = useRoute();
const artist = ref<any>(null);
const artistData = ref<any>({});
const isFollowing = ref(false);
const i18n = useI18n();
const { toast } = useToast();
const dominantColors = ref<string[]>([]);

const topTracks = computed(() => {
    return (
        artistData.value.top_tracks?.map((track: any) => ({
            ...track,
            album: {
                images: track.album?.images || artist.value?.images,
            },
        })) || []
    );
});

const albums = computed(() => {
    return (
        artistData.value.albums?.filter(
            (a: any) => a.album_group === "album",
        ) || []
    );
});

const singles = computed(() => {
    return (
        artistData.value.albums?.filter(
            (a: any) => a.album_group === "single",
        ) || []
    );
});

const heroSectionStyle = computed(() => {
    if (dominantColors.value.length >= 2) {
        return {
            background: `linear-gradient(135deg, ${dominantColors.value[0]} 0%, ${dominantColors.value[1]} 100%)`,
        };
    }
    return {
        background: "linear-gradient(135deg, #4f46e5 0%, #1e40af 100%)",
    };
});

// Helper functions
const formatNumber = (num: number) => {
    return new Intl.NumberFormat().format(num);
};

const toggleFollow = async () => {
    if (isFollowing.value) {
        const success = await RemoveArtist(artist.value.id);
        if (!success) {
            console.error("Failed to remove artist");
            toast({
                title: i18n.t("ArtistDetails.error_title"),
                description: i18n.t("ArtistDetails.error_unfollowing"),
                variant: "destructive",
            });
            return;
        }
    } else {
        const success = await AddArtist(artist.value.id);
        if (!success) {
            console.error("Failed to add artist");
            toast({
                title: i18n.t("ArtistDetails.error_title"),
                description: i18n.t("ArtistDetails.error_following"),
                variant: "destructive",
            });
            return;
        }
    }
    isFollowing.value = !isFollowing.value;
};

onMounted(async () => {
    const artistId = route.params.id as string;
    const data = await getArtistDetails(artistId);
    const subbed_artists = await GetArtistsFromDB();
    let isSubbed = false;
    isSubbed = subbed_artists.some(
        (artist: { SpotifyID: string; LastChecked: Date }) =>
            artist.SpotifyID === artistId,
    );
    artistData.value = data;
    artist.value = data.artist;
    isFollowing.value = isSubbed;

    if (artist.value?.images?.[0]?.url) {
        dominantColors.value = await GetDominantColor(
            artist.value.images[0].url,
        );
    }
});

const getArtistDetails = async (id: string) => {
    try {
        const artistData = await GetArtist(id);
        return artistData;
    } catch (error) {
        console.error("Error fetching artist details:", error);
        toast({
            title: i18n.t("TrackDetails.error_title"),
            description: i18n.t("TrackDetails.error_getting"),
            variant: "destructive",
        });
        return {};
    }
};
</script>

<style scoped>
.artist-detail {
    scrollbar-width: none; /* Firefox */
}
.artist-detail::-webkit-scrollbar {
    display: none; /* Chrome/Safari */
}
</style>
